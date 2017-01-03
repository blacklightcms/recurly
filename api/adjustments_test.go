package api_test

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	recurly "github.com/blacklightcms/go-recurly"
)

// TestAdjustmentEncoding ensures structs are encoded to XML properly.
// Because Recurly supports partial updates, it's important that only defined
// fields are handled properly -- including types like booleans and integers which
// have zero values that we want to send.
func TestAdjustments_Encoding(t *testing.T) {
	tests := []struct {
		v        recurly.Adjustment
		expected string
	}{
		// Unit amount in cents and currency are required fields. They should always be present.
		{v: recurly.Adjustment{}, expected: "<adjustment><unit_amount_in_cents>0</unit_amount_in_cents><currency></currency></adjustment>"},
		{v: recurly.Adjustment{UnitAmountInCents: 2000, Currency: "USD"}, expected: "<adjustment><unit_amount_in_cents>2000</unit_amount_in_cents><currency>USD</currency></adjustment>"},
		{v: recurly.Adjustment{Description: "Charge for extra bandwidth", UnitAmountInCents: 2000, Currency: "USD"}, expected: "<adjustment><description>Charge for extra bandwidth</description><unit_amount_in_cents>2000</unit_amount_in_cents><currency>USD</currency></adjustment>"},
		{v: recurly.Adjustment{Quantity: 1, UnitAmountInCents: 2000, Currency: "CAD"}, expected: "<adjustment><unit_amount_in_cents>2000</unit_amount_in_cents><quantity>1</quantity><currency>CAD</currency></adjustment>"},
		{v: recurly.Adjustment{AccountingCode: "bandwidth", UnitAmountInCents: 2000, Currency: "CAD"}, expected: "<adjustment><accounting_code>bandwidth</accounting_code><unit_amount_in_cents>2000</unit_amount_in_cents><currency>CAD</currency></adjustment>"},
		{v: recurly.Adjustment{TaxExempt: recurly.NewBool(false), UnitAmountInCents: 2000, Currency: "USD"}, expected: "<adjustment><unit_amount_in_cents>2000</unit_amount_in_cents><currency>USD</currency><tax_exempt>false</tax_exempt></adjustment>"},
		{v: recurly.Adjustment{TaxCode: "digital", UnitAmountInCents: 2000, Currency: "USD"}, expected: "<adjustment><unit_amount_in_cents>2000</unit_amount_in_cents><currency>USD</currency><tax_code>digital</tax_code></adjustment>"},
	}

	for _, tt := range tests {
		var buf bytes.Buffer
		if err := xml.NewEncoder(&buf).Encode(tt.v); err != nil {
			t.Fatalf("unexpected error: %v", err)
		} else if buf.String() != tt.expected {
			t.Fatalf("unexpected value: %s", buf.String())
		}
	}
}

func TestAdjustments_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/100/adjustments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.Header().Set("Link", `<https://your-subdomain.recurly.com/v2/accounts/100/adjustments?cursor=1304958672>; rel="next"`)
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
			<adjustments type="array">
				<adjustment href="https://your-subdomain.recurly.com/v2/adjustments/626db120a84102b1809909071c701c60" type="charge">
					<account href="https://your-subdomain.recurly.com/v2/accounts/100"/>
					<invoice href="https://your-subdomain.recurly.com/v2/invoices/1108"/>
					<subscription href="https://your-subdomain.recurly.com/v2/subscriptions/17caaca1716f33572edc8146e0aaefde"/>
					<uuid>626db120a84102b1809909071c701c60</uuid>
					<state>invoiced</state>
					<description>One-time Charged Fee</description>
					<accounting_code nil="nil"/>
					<product_code>basic</product_code>
					<origin>debit</origin>
					<unit_amount_in_cents type="integer">2000</unit_amount_in_cents>
					<quantity type="integer">1</quantity>
					<original_adjustment_uuid>2cc95aa62517e56d5bec3a48afa1b3b9</original_adjustment_uuid> <!-- Only shows if adjustment is a credit created from another credit. -->
					<discount_in_cents type="integer">0</discount_in_cents>
					<tax_in_cents type="integer">180</tax_in_cents>
					<total_in_cents type="integer">2180</total_in_cents>
					<currency>USD</currency>
					<taxable type="boolean">false</taxable>
					<tax_exempt type="boolean">false</tax_exempt>
					<tax_code nil="nil"/>
					<start_date type="datetime">2011-08-31T03:30:00Z</start_date>
					<end_date nil="nil"/>
					<created_at type="datetime">2011-08-31T03:30:00Z</created_at>
				</adjustment>
			</adjustments>`)
	})

	resp, adjustments, err := client.Adjustments.List("100", recurly.Params{"per_page": 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected list adjustments to return OK")
	} else if pp := resp.Request.URL.Query().Get("per_page"); pp != "1" {
		t.Fatalf("unexpected per_page: %s", pp)
	} else if resp.Next() != "1304958672" {
		t.Fatalf("unexpected cursor: %s", resp.Next())
	}

	ts, _ := time.Parse(recurly.GetDateTimeFormat(), "2011-08-31T03:30:00Z")
	if !reflect.DeepEqual(adjustments, []recurly.Adjustment{
		{
			AccountCode:            "100",
			InvoiceNumber:          1108,
			UUID:                   "626db120a84102b1809909071c701c60",
			State:                  "invoiced",
			Description:            "One-time Charged Fee",
			ProductCode:            "basic",
			Origin:                 "debit",
			UnitAmountInCents:      2000,
			Quantity:               1,
			OriginalAdjustmentUUID: "2cc95aa62517e56d5bec3a48afa1b3b9",
			TaxInCents:             180,
			TotalInCents:           2180,
			Currency:               "USD",
			Taxable:                recurly.NewBool(false),
			TaxExempt:              recurly.NewBool(false),
			StartDate:              recurly.NewTime(ts),
			CreatedAt:              recurly.NewTime(ts),
		},
	}) {
		t.Fatalf("unexpected adjustments: %v", adjustments)
	}
}

func TestAdjustments_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/adjustments/626db120a84102b1809909071c701c60", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
			<adjustment href="https://your-subdomain.recurly.com/v2/adjustments/626db120a84102b1809909071c701c60" type="charge">
				<account href="https://your-subdomain.recurly.com/v2/accounts/100"/>
				<invoice href="https://your-subdomain.recurly.com/v2/invoices/1108"/>
				<uuid>626db120a84102b1809909071c701c60</uuid>
				<state>invoiced</state>
				<description>One-time Charged Fee</description>
				<accounting_code/>
				<product_code>basic</product_code>
				<origin>debit</origin>
				<unit_amount_in_cents type="integer">2000</unit_amount_in_cents>
				<quantity type="integer">1</quantity>
				<original_adjustment_uuid>2cc95aa62517e56d5bec3a48afa1b3b9</original_adjustment_uuid> <!-- Only shows if adjustment is a credit created from another credit. -->
				<discount_in_cents type="integer">0</discount_in_cents>
				<tax_in_cents type="integer">175</tax_in_cents>
				<total_in_cents type="integer">2175</total_in_cents>
				<currency>USD</currency>
				<taxable type="boolean">false</taxable>
				<tax_type>usst</tax_type>
				<tax_region>CA</tax_region>
				<tax_rate type="float">0.0875</tax_rate>
				<tax_exempt type="boolean">false</tax_exempt>
				<tax_details type="array">
					<tax_detail>
						<name>california</name>
						<type>state</type>
						<tax_rate type="float">0.065</tax_rate>
						<tax_in_cents type="integer">130</tax_in_cents>
					</tax_detail>
					<tax_detail>
						<name>san mateo county</name>
						<type>county</type>
						<tax_rate type="float">0.01</tax_rate>
						<tax_in_cents type="integer">20</tax_in_cents>
					</tax_detail>
					<tax_detail>
						<name>sf municipal tax</name>
						<type>city</type>
						<tax_rate type="float">0.0</tax_rate>
						<tax_in_cents type="integer">0</tax_in_cents>
					</tax_detail>
					<tax_detail>
						<name nil="nil"/>
						<type>special</type>
						<tax_rate type="float">0.0125</tax_rate>
						<tax_in_cents type="integer">25</tax_in_cents>
					</tax_detail>
				</tax_details>
				<start_date type="datetime">2015-02-04T23:13:07Z</start_date>
				<end_date nil="nil"/>
				<created_at type="datetime">2015-02-04T23:13:07Z</created_at>
			</adjustment>`)
	})

	resp, adjustment, err := client.Adjustments.Get("626db120a84102b1809909071c701c60")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected get adjustment to return OK")
	}

	ts, _ := time.Parse(recurly.GetDateTimeFormat(), "2015-02-04T23:13:07Z")
	if !reflect.DeepEqual(adjustment, &recurly.Adjustment{
		AccountCode:            "100",
		InvoiceNumber:          1108,
		UUID:                   "626db120a84102b1809909071c701c60",
		State:                  "invoiced",
		Description:            "One-time Charged Fee",
		ProductCode:            "basic",
		Origin:                 "debit",
		UnitAmountInCents:      2000,
		Quantity:               1,
		OriginalAdjustmentUUID: "2cc95aa62517e56d5bec3a48afa1b3b9",
		TaxInCents:             175,
		TotalInCents:           2175,
		Currency:               "USD",
		Taxable:                recurly.NewBool(false),
		TaxType:                "usst",
		TaxRegion:              "CA",
		TaxRate:                0.0875,
		TaxExempt:              recurly.NewBool(false),
		TaxDetails: []recurly.TaxDetail{
			{
				XMLName:    xml.Name{Local: "tax_detail"},
				Name:       "california",
				Type:       "state",
				TaxRate:    0.065,
				TaxInCents: 130,
			},
			{
				XMLName:    xml.Name{Local: "tax_detail"},
				Name:       "san mateo county",
				Type:       "county",
				TaxRate:    0.01,
				TaxInCents: 20,
			},
			{
				XMLName:    xml.Name{Local: "tax_detail"},
				Name:       "sf municipal tax",
				Type:       "city",
				TaxRate:    0.0,
				TaxInCents: 0,
			},
			{
				XMLName:    xml.Name{Local: "tax_detail"},
				Type:       "special",
				TaxRate:    0.0125,
				TaxInCents: 25,
			},
		},
		StartDate: recurly.NewTime(ts),
		CreatedAt: recurly.NewTime(ts),
	}) {
		t.Fatalf("unexpected adjustment: %v", adjustment)
	}
}

func TestAdjustments_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/1/adjustments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(201)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><adjustment></adjustment>`)
	})

	resp, _, err := client.Adjustments.Create("1", recurly.Adjustment{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.StatusCode != 201 {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
}

func TestAdjustments_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/adjustments/945a4cb9afd64300b97b138407a51aef", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(204)
	})

	resp, err := client.Adjustments.Delete("945a4cb9afd64300b97b138407a51aef")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.StatusCode != 204 {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
}
