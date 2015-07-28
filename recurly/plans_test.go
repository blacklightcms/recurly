package recurly

import (
	"bytes"
	"encoding/xml"
	"testing"
)

// TestPlansEncoding ensures structs are encoded to XML properly.
// Because Recurly supports partial updates, it's important that only defined
// fields are handled properly -- including types like booleans and integers which
// have zero values that we want to send.
func TestPlansEncoding(t *testing.T) {
	suite := []map[string]interface{}{
		// name is a required field. It should always be present.
		map[string]interface{}{"struct": Plan{}, "xml": "<plan><name></name><unit_amount_in_cents></unit_amount_in_cents></plan>"},
		map[string]interface{}{"struct": Plan{Name: "Gold plan", UnitAmountInCents: UnitAmount{USD: 1500}, Description: "abc"}, "xml": "<plan><name>Gold plan</name><description>abc</description><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		map[string]interface{}{"struct": Plan{Name: "Gold plan", UnitAmountInCents: UnitAmount{USD: 1500}, AccountingCode: "gold"}, "xml": "<plan><name>Gold plan</name><accounting_code>gold</accounting_code><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		map[string]interface{}{"struct": Plan{Name: "Gold plan", UnitAmountInCents: UnitAmount{USD: 1500}, IntervalUnit: "months"}, "xml": "<plan><name>Gold plan</name><plan_interval_unit>months</plan_interval_unit><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		map[string]interface{}{"struct": Plan{Name: "Gold plan", UnitAmountInCents: UnitAmount{USD: 1500}, IntervalLength: 1}, "xml": "<plan><name>Gold plan</name><plan_interval_length>1</plan_interval_length><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		map[string]interface{}{"struct": Plan{Name: "Gold plan", UnitAmountInCents: UnitAmount{USD: 1500}, TrialIntervalUnit: "days"}, "xml": "<plan><name>Gold plan</name><trial_interval_unit>days</trial_interval_unit><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		map[string]interface{}{"struct": Plan{Name: "Gold plan", UnitAmountInCents: UnitAmount{USD: 1500}, TrialIntervalLength: 10}, "xml": "<plan><name>Gold plan</name><trial_interval_length>10</trial_interval_length><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		map[string]interface{}{"struct": Plan{Name: "Gold plan", UnitAmountInCents: UnitAmount{USD: 1500}, IntervalUnit: "months"}, "xml": "<plan><name>Gold plan</name><plan_interval_unit>months</plan_interval_unit><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
		map[string]interface{}{"struct": Plan{Name: "Gold plan", UnitAmountInCents: UnitAmount{USD: 1500}, SetupFeeInCents: &UnitAmount{USD: 1000, EUR: 800}}, "xml": "<plan><name>Gold plan</name><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents><setup_fee_in_cents><USD>1000</USD><EUR>800</EUR></setup_fee_in_cents></plan>"},
		map[string]interface{}{"struct": Plan{Name: "Gold plan", UnitAmountInCents: UnitAmount{USD: 1500}, TotalBillingCycles: 24}, "xml": "<plan><name>Gold plan</name><total_billing_cycles>24</total_billing_cycles><unit_amount_in_cents><USD>1500</USD></unit_amount_in_cents></plan>"},
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
