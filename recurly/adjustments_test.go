package recurly

import (
	"bytes"
	"encoding/xml"
	"testing"
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
