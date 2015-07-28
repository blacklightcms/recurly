package recurly

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

// TestAdjustmentEncoding ensures structs are encoded to XML properly.
// Because Recurly supports partial updates, it's important that only defined
// fields are handled properly -- including types like booleans and integers which
// have zero values that we want to send.
func TestAdjustmentsEncoding(t *testing.T) {
	suite := []map[string]interface{}{
		// Unit amount in cents and currency are required fields. They should always be present.
		map[string]interface{}{"struct": Adjustment{}, "xml": "<adjustment><unit_amount_in_cents>0</unit_amount_in_cents><currency></currency></adjustment>"},
		map[string]interface{}{"struct": Adjustment{UnitAmountInCents: 2000, Currency: "USD"}, "xml": "<adjustment><unit_amount_in_cents>2000</unit_amount_in_cents><currency>USD</currency></adjustment>"},
		map[string]interface{}{"struct": Adjustment{Description: "Charge for extra bandwidth", UnitAmountInCents: 2000, Currency: "USD"}, "xml": "<adjustment><description>Charge for extra bandwidth</description><unit_amount_in_cents>2000</unit_amount_in_cents><currency>USD</currency></adjustment>"},
		map[string]interface{}{"struct": Adjustment{Quantity: 1, UnitAmountInCents: 2000, Currency: "CAD"}, "xml": "<adjustment><unit_amount_in_cents>2000</unit_amount_in_cents><quantity>1</quantity><currency>CAD</currency></adjustment>"},
		map[string]interface{}{"struct": Adjustment{AccountingCode: "bandwidth", UnitAmountInCents: 2000, Currency: "CAD"}, "xml": "<adjustment><accounting_code>bandwidth</accounting_code><unit_amount_in_cents>2000</unit_amount_in_cents><currency>CAD</currency></adjustment>"},
		map[string]interface{}{"struct": Adjustment{TaxExempt: NewBool(false), UnitAmountInCents: 2000, Currency: "USD"}, "xml": "<adjustment><unit_amount_in_cents>2000</unit_amount_in_cents><currency>USD</currency><tax_exempt>false</tax_exempt></adjustment>"},
		map[string]interface{}{"struct": Adjustment{TaxCode: "digital", UnitAmountInCents: 2000, Currency: "USD"}, "xml": "<adjustment><unit_amount_in_cents>2000</unit_amount_in_cents><currency>USD</currency><tax_code>digital</tax_code></adjustment>"},
	}

	for _, s := range suite {
		buf := new(bytes.Buffer)
		err := xml.NewEncoder(buf).Encode(s["struct"])
		if err != nil {
			t.Errorf("TestAdjustmentEncoding Error: %s", err)
		}

		if buf.String() != s["xml"] {
			t.Errorf("TestAdjustmentEncoding Error: Expected %s, given %s", s["xml"], buf.String())
		}
	}
}

func TestAdjustmentsList(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/100/adjustments", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("TestAdjustmentsList Error: Expected %s request, given %s", "GET", r.Method)
		}
		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?>
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

	r, adjustments, err := client.Adjustments.List("100", Params{"per_page": 1})
	if err != nil {
		t.Errorf("TestAdjustmentsList Error: Error occured making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestAdjustmentsList Error: Expected list adjustments to return OK")
	}

	if len(adjustments) != 1 {
		t.Fatalf("TestAdjustmentsList Error: Expected 1 adjustment returned, given %d", len(adjustments))
	}

	if r.Request.URL.Query().Get("per_page") != "1" {
		t.Errorf("TestAdjustmentsList Error: Expected per_page parameter of 1, given %s", r.Request.URL.Query().Get("per_page"))
	}

	ts, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2011-08-31T03:30:00Z")
	for _, given := range adjustments {
		expected := Adjustment{
			XMLName: xml.Name{Local: "adjustment"},
			Account: href{
				HREF: "https://your-subdomain.recurly.com/v2/accounts/100",
				Code: "100",
			},
			Invoice: href{
				HREF: "https://your-subdomain.recurly.com/v2/invoices/1108",
				Code: "1108",
			},
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
			Taxable:                NewBool(false),
			TaxExempt:              NewBool(false),
			StartDate:              NewTime(ts),
			CreatedAt:              NewTime(ts),
		}

		if !reflect.DeepEqual(expected, given) {
			t.Errorf("TestAdjustmentsList Error: expected adjustment to equal %#v, given %#v", expected, given)
		}
	}
}

func TestGetAdjustment(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/adjustments/626db120a84102b1809909071c701c60", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("TestGetAdjustment Error: Expected %s request, given %s", "GET", r.Method)
		}
		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?>
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

	r, adjustment, err := client.Adjustments.Get("626db120a84102b1809909071c701c60")
	if err != nil {
		t.Errorf("TestGetAdjustment Error: Error occured making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestGetAdjustment Error: Expected get adjustment to return OK")
	}

	ts, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2015-02-04T23:13:07Z")
	expected := Adjustment{
		XMLName: xml.Name{Local: "adjustment"},
		Account: href{
			HREF: "https://your-subdomain.recurly.com/v2/accounts/100",
			Code: "100",
		},
		Invoice: href{
			HREF: "https://your-subdomain.recurly.com/v2/invoices/1108",
			Code: "1108",
		},
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
		Taxable:                NewBool(false),
		TaxType:                "usst",
		TaxRegion:              "CA",
		TaxRate:                0.0875,
		TaxExempt:              NewBool(false),
		TaxDetails: &[]TaxDetails{
			TaxDetails{
				XMLName:    xml.Name{Local: "tax_detail"},
				Name:       "california",
				Type:       "state",
				TaxRate:    0.065,
				TaxInCents: 130,
			},
			TaxDetails{
				XMLName:    xml.Name{Local: "tax_detail"},
				Name:       "san mateo county",
				Type:       "county",
				TaxRate:    0.01,
				TaxInCents: 20,
			},
			TaxDetails{
				XMLName:    xml.Name{Local: "tax_detail"},
				Name:       "sf municipal tax",
				Type:       "city",
				TaxRate:    0.0,
				TaxInCents: 0,
			},
			TaxDetails{
				XMLName:    xml.Name{Local: "tax_detail"},
				Type:       "special",
				TaxRate:    0.0125,
				TaxInCents: 25,
			},
		},
		StartDate: NewTime(ts),
		CreatedAt: NewTime(ts),
	}

	if !reflect.DeepEqual(expected, adjustment) {
		t.Errorf("TestGetAdjustment Error: expected adjustment to equal %#v, given %#v", expected, adjustment)
	}
}

func TestCreateAdjustment(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/1/adjustments", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("TestCreateAdjustment Error: Expected %s request, given %s", "POST", r.Method)
		}
		rw.WriteHeader(201)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?><adjustment></adjustment>`)
	})

	r, _, err := client.Adjustments.Create("1", Adjustment{})
	if err != nil {
		t.Errorf("TestCreateAdjustment Error: Error occured making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestCreateAdjustment Error: Expected create adjustment to return OK")
	}
}

func TestDeleteAdjustment(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/adjustments/945a4cb9afd64300b97b138407a51aef", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("TestDeleteAdjustment Error: Expected %s request, given %s", "DELETE", r.Method)
		}
		rw.WriteHeader(204)
	})

	r, err := client.Adjustments.Delete("945a4cb9afd64300b97b138407a51aef")
	if err != nil {
		t.Errorf("TestDeleteAdjustment Error: Error occured making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestDeleteAdjustment Error: Expected create adjustment to return OK")
	}
}
