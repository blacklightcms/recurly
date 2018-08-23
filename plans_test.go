package recurly_test

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/launchpadcentral/recurly"
	"github.com/google/go-cmp/cmp"
)

// TestPlansEncoding ensures structs are encoded to XML properly.
// Because Recurly supports partial updates, it's important that only defined
// fields are handled properly -- including types like booleans and integers which
// have zero values that we want to send.
func TestPlans_Encoding(t *testing.T) {
	tests := []struct {
		v        recurly.Plan
		expected string
	}{
		// name is a required field. It should always be present.
		{v: recurly.Plan{}, expected: "<plan><name></name></plan>"},
		{v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, Description: "abc"}, expected: "<plan><name>Gold plan</name><description>abc</description><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		{v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, AccountingCode: "gold"}, expected: "<plan><name>Gold plan</name><accounting_code>gold</accounting_code><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		{v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, IntervalUnit: "months"}, expected: "<plan><name>Gold plan</name><plan_interval_unit>months</plan_interval_unit><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		{v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, IntervalLength: 1}, expected: "<plan><name>Gold plan</name><plan_interval_length>1</plan_interval_length><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		{v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, TrialIntervalUnit: "days"}, expected: "<plan><name>Gold plan</name><trial_interval_unit>days</trial_interval_unit><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		{v: recurly.Plan{Name: "Gold plan", AutoRenew: true, UnitAmountInCents: recurly.UnitAmount{USD: 1500}, TrialIntervalLength: 10}, expected: "<plan><name>Gold plan</name><trial_interval_length>10</trial_interval_length><auto_renew>true</auto_renew><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		{v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, IntervalUnit: "months"}, expected: "<plan><name>Gold plan</name><plan_interval_unit>months</plan_interval_unit><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		{v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, SetupFeeInCents: recurly.UnitAmount{USD: 1000, EUR: 800}}, expected: "<plan><name>Gold plan</name><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents><setup_fee_in_cents><USD>1000</USD><EUR>800</EUR></setup_fee_in_cents></plan>"},
		{v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, TotalBillingCycles: recurly.NewInt(24)}, expected: "<plan><name>Gold plan</name><total_billing_cycles>24</total_billing_cycles><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		{v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, UnitName: "unit"}, expected: "<plan><name>Gold plan</name><unit_name>unit</unit_name><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		{v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, DisplayQuantity: recurly.NewBool(true)}, expected: "<plan><name>Gold plan</name><display_quantity>true</display_quantity><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		{v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, DisplayQuantity: recurly.NewBool(false)}, expected: "<plan><name>Gold plan</name><display_quantity>false</display_quantity><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		{v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, SuccessURL: "https://example.com/success"}, expected: "<plan><name>Gold plan</name><success_url>https://example.com/success</success_url><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		{v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, CancelURL: "https://example.com/cancel"}, expected: "<plan><name>Gold plan</name><cancel_url>https://example.com/cancel</cancel_url><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		{v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, TaxExempt: recurly.NewBool(true)}, expected: "<plan><name>Gold plan</name><tax_exempt>true</tax_exempt><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		{v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, TaxExempt: recurly.NewBool(false)}, expected: "<plan><name>Gold plan</name><tax_exempt>false</tax_exempt><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		{v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, TaxCode: "physical"}, expected: "<plan><name>Gold plan</name><tax_code>physical</tax_code><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
	}

	for _, tt := range tests {
		var buf bytes.Buffer
		if err := xml.NewEncoder(&buf).Encode(tt.v); err != nil {
			t.Fatalf("TestPlansEncoding Error: %s", err)
		} else if buf.String() != tt.expected {
			t.Fatalf("unexpected encoding: %s", buf.String())
		}
	}
}

func TestPlans_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/plans", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
		<plans type="array">
			<plan href="https://your-subdomain.recurly.com/v2/plans/gold">
				<add_ons href="https://your-subdomain.recurly.com/v2/plans/gold/add_ons"/>
				<plan_code>gold</plan_code>
				<name>Gold plan</name>
				<description nil="nil"/>
				<success_url nil="nil"/>
				<cancel_url nil="nil"/>
				<display_donation_amounts type="boolean">false</display_donation_amounts>
				<display_quantity type="boolean">false</display_quantity>
				<display_phone_number type="boolean">false</display_phone_number>
				<bypass_hosted_confirmation type="boolean">false</bypass_hosted_confirmation>
				<unit_name>unit</unit_name>
				<payment_page_tos_link nil="nil"/>
				<plan_interval_length type="integer">1</plan_interval_length>
				<plan_interval_unit>months</plan_interval_unit>
				<trial_interval_length type="integer">0</trial_interval_length>
				<trial_interval_unit>days</trial_interval_unit>
				<total_billing_cycles nil="nil"/>
				<accounting_code nil="nil"/>
				<created_at type="datetime">2015-05-29T17:38:15Z</created_at>
				<tax_exempt type="boolean">false</tax_exempt>
				<tax_code nil="nil"/>
				<unit_amount_in_cents>
					<USD type="integer">6000</USD>
					<EUR type="integer">4500</EUR>
				</unit_amount_in_cents>
				<setup_fee_in_cents>
					<USD type="integer">1000</USD>
					<EUR type="integer">800</EUR>
				</setup_fee_in_cents>
			</plan>
		</plans>`)
	})

	r, plans, err := client.Plans.List(recurly.Params{"per_page": 1})
	if err != nil {
		t.Fatalf("error occurred making API call. Err: %s", err)
	} else if r.IsError() {
		t.Fatal("expected list plans to return OK")
	} else if r.Request.URL.Query().Get("per_page") != "1" {
		t.Fatalf("expected per_page parameter of 1, given %s", r.Request.URL.Query().Get("per_page"))
	}

	ts, _ := time.Parse(recurly.DateTimeFormat, "2015-05-29T17:38:15Z")
	if diff := cmp.Diff(plans, []recurly.Plan{
		{
			XMLName: xml.Name{Local: "plan"},
			Code:    "gold",
			Name:    "Gold plan",
			DisplayDonationAmounts:   recurly.NewBool(false),
			DisplayQuantity:          recurly.NewBool(false),
			DisplayPhoneNumber:       recurly.NewBool(false),
			BypassHostedConfirmation: recurly.NewBool(false),
			UnitName:                 "unit",
			IntervalUnit:             "months",
			IntervalLength:           1,
			TrialIntervalUnit:        "days",
			TaxExempt:                recurly.NewBool(false),
			UnitAmountInCents: recurly.UnitAmount{
				USD: 6000,
				EUR: 4500,
			},
			SetupFeeInCents: recurly.UnitAmount{
				USD: 1000,
				EUR: 800,
			},
			CreatedAt: recurly.NewTime(ts),
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestPlans_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/plans/gold", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
		<plan href="https://your-subdomain.recurly.com/v2/plans/gold">
			<add_ons href="https://your-subdomain.recurly.com/v2/plans/gold/add_ons"/>
			<plan_code>gold</plan_code>
			<name>Gold plan</name>
			<description nil="nil"/>
			<success_url nil="nil"/>
			<cancel_url nil="nil"/>
			<display_donation_amounts type="boolean">false</display_donation_amounts>
			<display_quantity type="boolean">false</display_quantity>
			<display_phone_number type="boolean">false</display_phone_number>
			<bypass_hosted_confirmation type="boolean">false</bypass_hosted_confirmation>
			<unit_name>unit</unit_name>
			<payment_page_tos_link nil="nil"/>
			<plan_interval_length type="integer">1</plan_interval_length>
			<plan_interval_unit>months</plan_interval_unit>
			<trial_interval_length type="integer">0</trial_interval_length>
			<trial_interval_unit>days</trial_interval_unit>
			<total_billing_cycles nil="nil"/>
			<accounting_code nil="nil"/>
			<created_at type="datetime">2015-05-29T17:38:15Z</created_at>
			<tax_exempt type="boolean">false</tax_exempt>
			<tax_code nil="nil"/>
			<unit_amount_in_cents>
				<USD type="integer">6000</USD>
				<EUR type="integer">4500</EUR>
			</unit_amount_in_cents>
			<setup_fee_in_cents>
				<USD type="integer">1000</USD>
				<EUR type="integer">800</EUR>
			</setup_fee_in_cents>
		</plan>`)
	})

	r, plan, err := client.Plans.Get("gold")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected get plan to return OK")
	}

	ts, _ := time.Parse(recurly.DateTimeFormat, "2015-05-29T17:38:15Z")
	if diff := cmp.Diff(plan, &recurly.Plan{
		XMLName: xml.Name{Local: "plan"},
		Code:    "gold",
		Name:    "Gold plan",
		DisplayDonationAmounts:   recurly.NewBool(false),
		DisplayQuantity:          recurly.NewBool(false),
		DisplayPhoneNumber:       recurly.NewBool(false),
		BypassHostedConfirmation: recurly.NewBool(false),
		UnitName:                 "unit",
		IntervalUnit:             "months",
		IntervalLength:           1,
		TrialIntervalUnit:        "days",
		TaxExempt:                recurly.NewBool(false),
		UnitAmountInCents: recurly.UnitAmount{
			USD: 6000,
			EUR: 4500,
		},
		SetupFeeInCents: recurly.UnitAmount{
			USD: 1000,
			EUR: 800,
		},
		CreatedAt: recurly.NewTime(ts),
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestPlans_Get_ErrNotFound(t *testing.T) {
	setup()
	defer teardown()

	var invoked bool
	mux.HandleFunc("/v2/plans/gold", func(w http.ResponseWriter, r *http.Request) {
		invoked = true
		w.WriteHeader(http.StatusNotFound)
	})

	_, plan, err := client.Plans.Get("gold")
	if !invoked {
		t.Fatal("handler not invoked")
	} else if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if plan != nil {
		t.Fatalf("expected plan to be nil: %#v", plan)
	}
}

func TestPlans_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/plans", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(201)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><plan></plan>`)
	})

	r, _, err := client.Plans.Create(recurly.Plan{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected create plan to return OK")
	}
}

func TestPlans_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/plans/silver", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><plan></plan>`)
	})

	r, _, err := client.Plans.Update("silver", recurly.Plan{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected update plan to return OK")
	}
}

func TestPlans_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/plans/platinum", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(204)
	})

	r, err := client.Plans.Delete("platinum")
	if err != nil {
		t.Fatalf("error occurred making API call. Err: %s", err)
	} else if r.IsError() {
		t.Fatal("expected deleting plan to return OK")
	}
}
