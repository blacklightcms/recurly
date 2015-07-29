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

// TestPlansEncoding ensures structs are encoded to XML properly.
// Because Recurly supports partial updates, it's important that only defined
// fields are handled properly -- including types like booleans and integers which
// have zero values that we want to send.
func TestPlansEncoding(t *testing.T) {
	suite := []map[string]interface{}{
		// name is a required field. It should always be present.
		map[string]interface{}{"struct": Plan{}, "xml": "<plan><name></name></plan>"},
		map[string]interface{}{"struct": Plan{Name: "Gold plan", UnitAmountInCents: UnitAmount{USD: 1500}, Description: "abc"}, "xml": "<plan><name>Gold plan</name><description>abc</description><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		map[string]interface{}{"struct": Plan{Name: "Gold plan", UnitAmountInCents: UnitAmount{USD: 1500}, AccountingCode: "gold"}, "xml": "<plan><name>Gold plan</name><accounting_code>gold</accounting_code><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		map[string]interface{}{"struct": Plan{Name: "Gold plan", UnitAmountInCents: UnitAmount{USD: 1500}, IntervalUnit: "months"}, "xml": "<plan><name>Gold plan</name><plan_interval_unit>months</plan_interval_unit><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		map[string]interface{}{"struct": Plan{Name: "Gold plan", UnitAmountInCents: UnitAmount{USD: 1500}, IntervalLength: 1}, "xml": "<plan><name>Gold plan</name><plan_interval_length>1</plan_interval_length><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		map[string]interface{}{"struct": Plan{Name: "Gold plan", UnitAmountInCents: UnitAmount{USD: 1500}, TrialIntervalUnit: "days"}, "xml": "<plan><name>Gold plan</name><trial_interval_unit>days</trial_interval_unit><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		map[string]interface{}{"struct": Plan{Name: "Gold plan", UnitAmountInCents: UnitAmount{USD: 1500}, TrialIntervalLength: 10}, "xml": "<plan><name>Gold plan</name><trial_interval_length>10</trial_interval_length><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		map[string]interface{}{"struct": Plan{Name: "Gold plan", UnitAmountInCents: UnitAmount{USD: 1500}, IntervalUnit: "months"}, "xml": "<plan><name>Gold plan</name><plan_interval_unit>months</plan_interval_unit><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		map[string]interface{}{"struct": Plan{Name: "Gold plan", UnitAmountInCents: UnitAmount{USD: 1500}, SetupFeeInCents: UnitAmount{USD: 1000, EUR: 800}}, "xml": "<plan><name>Gold plan</name><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents><setup_fee_in_cents><USD>1000</USD><EUR>800</EUR></setup_fee_in_cents></plan>"},
		map[string]interface{}{"struct": Plan{Name: "Gold plan", UnitAmountInCents: UnitAmount{USD: 1500}, TotalBillingCycles: NewInt(24)}, "xml": "<plan><name>Gold plan</name><total_billing_cycles>24</total_billing_cycles><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		map[string]interface{}{"struct": Plan{Name: "Gold plan", UnitAmountInCents: UnitAmount{USD: 1500}, UnitName: "unit"}, "xml": "<plan><name>Gold plan</name><unit_name>unit</unit_name><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		map[string]interface{}{"struct": Plan{Name: "Gold plan", UnitAmountInCents: UnitAmount{USD: 1500}, DisplayQuantity: NewBool(true)}, "xml": "<plan><name>Gold plan</name><display_quantity>true</display_quantity><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		map[string]interface{}{"struct": Plan{Name: "Gold plan", UnitAmountInCents: UnitAmount{USD: 1500}, DisplayQuantity: NewBool(false)}, "xml": "<plan><name>Gold plan</name><display_quantity>false</display_quantity><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		map[string]interface{}{"struct": Plan{Name: "Gold plan", UnitAmountInCents: UnitAmount{USD: 1500}, SuccessURL: "https://example.com/success"}, "xml": "<plan><name>Gold plan</name><success_url>https://example.com/success</success_url><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		map[string]interface{}{"struct": Plan{Name: "Gold plan", UnitAmountInCents: UnitAmount{USD: 1500}, CancelURL: "https://example.com/cancel"}, "xml": "<plan><name>Gold plan</name><cancel_url>https://example.com/cancel</cancel_url><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		map[string]interface{}{"struct": Plan{Name: "Gold plan", UnitAmountInCents: UnitAmount{USD: 1500}, TaxExempt: NewBool(true)}, "xml": "<plan><name>Gold plan</name><tax_exempt>true</tax_exempt><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		map[string]interface{}{"struct": Plan{Name: "Gold plan", UnitAmountInCents: UnitAmount{USD: 1500}, TaxExempt: NewBool(false)}, "xml": "<plan><name>Gold plan</name><tax_exempt>false</tax_exempt><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		map[string]interface{}{"struct": Plan{Name: "Gold plan", UnitAmountInCents: UnitAmount{USD: 1500}, TaxCode: "physical"}, "xml": "<plan><name>Gold plan</name><tax_code>physical</tax_code><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
	}

	for _, s := range suite {
		buf := new(bytes.Buffer)
		err := xml.NewEncoder(buf).Encode(s["struct"])
		if err != nil {
			t.Errorf("TestPlansEncoding Error: %s", err)
		}

		if buf.String() != s["xml"] {
			t.Errorf("TestPlansEncoding Error: Expected %s, given %s", s["xml"], buf.String())
		}
	}
}

func TestListPlans(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/plans", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("TestListPlans Error: Expected %s request, given %s", "GET", r.Method)
		}
		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?>
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

	r, plans, err := client.Plans.List(Params{"per_page": 1})
	if err != nil {
		t.Errorf("TestListPlans Error: Error occured making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestListPlans Error: Expected list plans to return OK")
	}

	if len(plans) != 1 {
		t.Fatalf("TestListPlans Error: Expected 1 plan returned, given %d", len(plans))
	}

	if r.Request.URL.Query().Get("per_page") != "1" {
		t.Errorf("TestListPlans Error: Expected per_page parameter of 1, given %s", r.Request.URL.Query().Get("per_page"))
	}

	ts, _ := time.Parse(datetimeFormat, "2015-05-29T17:38:15Z")
	for _, given := range plans {
		expected := Plan{
			XMLName: xml.Name{Local: "plan"},
			Code:    "gold",
			Name:    "Gold plan",
			DisplayDonationAmounts:   NewBool(false),
			DisplayQuantity:          NewBool(false),
			DisplayPhoneNumber:       NewBool(false),
			BypassHostedConfirmation: NewBool(false),
			UnitName:                 "unit",
			IntervalUnit:             "months",
			IntervalLength:           1,
			TrialIntervalUnit:        "days",
			TaxExempt:                NewBool(false),
			UnitAmountInCents: UnitAmount{
				USD: 6000,
				EUR: 4500,
			},
			SetupFeeInCents: UnitAmount{
				USD: 1000,
				EUR: 800,
			},
			CreatedAt: NewTime(ts),
		}

		if !reflect.DeepEqual(expected, given) {
			t.Errorf("TestListPlans Error: expected plan to equal %#v, given %#v", expected, given)
		}
	}
}

func TestGetPlan(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/plans/gold", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("TestGetPlan Error: Expected %s request, given %s", "GET", r.Method)
		}
		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?>
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
		t.Errorf("TestGetPlan Error: Error occured making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestGetPlan Error: Expected get plan to return OK")
	}

	ts, _ := time.Parse(datetimeFormat, "2015-05-29T17:38:15Z")
	expected := Plan{
		XMLName: xml.Name{Local: "plan"},
		Code:    "gold",
		Name:    "Gold plan",
		DisplayDonationAmounts:   NewBool(false),
		DisplayQuantity:          NewBool(false),
		DisplayPhoneNumber:       NewBool(false),
		BypassHostedConfirmation: NewBool(false),
		UnitName:                 "unit",
		IntervalUnit:             "months",
		IntervalLength:           1,
		TrialIntervalUnit:        "days",
		TaxExempt:                NewBool(false),
		UnitAmountInCents: UnitAmount{
			USD: 6000,
			EUR: 4500,
		},
		SetupFeeInCents: UnitAmount{
			USD: 1000,
			EUR: 800,
		},
		CreatedAt: NewTime(ts),
	}

	if !reflect.DeepEqual(expected, plan) {
		t.Errorf("TestGetPlan Error: expected plan to equal %#v, given %#v", expected, plan)
	}
}

func TestCreatePlan(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/plans", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("TestCreatePlan Error: Expected %s request, given %s", "POST", r.Method)
		}
		rw.WriteHeader(201)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?><plan></plan>`)
	})

	r, _, err := client.Plans.Create(Plan{})
	if err != nil {
		t.Errorf("TestCreatePlan Error: Error occured making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestCreatePlan Error: Expected create plan to return OK")
	}
}

func TestUpdatePlan(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/plans/silver", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("TestUpdatePlan Error: Expected %s request, given %s", "PUT", r.Method)
		}
		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?><plan></plan>`)
	})

	r, _, err := client.Plans.Update("silver", Plan{})
	if err != nil {
		t.Errorf("TestUpdatePlan Error: Error occured making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestUpdatePlan Error: Expected update plan to return OK")
	}
}

func TestDeletePlan(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/plans/platinum", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("TestDeletePlan Error: Expected %s request, given %s", "DELETE", r.Method)
		}
		rw.WriteHeader(204)
	})

	r, err := client.Plans.Delete("platinum")
	if err != nil {
		t.Errorf("TestDeletePlan Error: Error occured making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestDeletePlan Error: Expected deleting plan to return OK")
	}
}
