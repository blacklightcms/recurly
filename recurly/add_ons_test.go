package recurly

import (
	"bytes"
	"encoding/xml"
	"testing"
)

// TestAddOnEncoding ensures structs are encoded to XML properly.
// Because Recurly supports partial updates, it's important that only defined
// fields are handled properly -- including types like booleans and integers which
// have zero values that we want to send.
func TestAddOnsEncoding(t *testing.T) {
	suite := []map[string]interface{}{
		map[string]interface{}{"struct": AddOn{}, "xml": "<add_on></add_on>"},
		map[string]interface{}{"struct": AddOn{Code: "xyz"}, "xml": "<add_on><add_on_code>xyz</add_on_code></add_on>"},
		map[string]interface{}{"struct": AddOn{Name: "IP Addresses"}, "xml": "<add_on><name>IP Addresses</name></add_on>"},
		map[string]interface{}{"struct": AddOn{DefaultQuantity: NewInt(0)}, "xml": "<add_on><default_quantity>0</default_quantity></add_on>"},
		map[string]interface{}{"struct": AddOn{DefaultQuantity: NewInt(1)}, "xml": "<add_on><default_quantity>1</default_quantity></add_on>"},
		map[string]interface{}{"struct": AddOn{DisplayQuantityOnHostedPage: NewBool(true)}, "xml": "<add_on><display_quantity_on_hosted_page>true</display_quantity_on_hosted_page></add_on>"},
		map[string]interface{}{"struct": AddOn{DisplayQuantityOnHostedPage: NewBool(false)}, "xml": "<add_on><display_quantity_on_hosted_page>false</display_quantity_on_hosted_page></add_on>"},
		map[string]interface{}{"struct": AddOn{TaxCode: "digital"}, "xml": "<add_on><tax_code>digital</tax_code></add_on>"},
		map[string]interface{}{"struct": AddOn{UnitAmountInCents: &UnitAmount{USD: 200}}, "xml": "<add_on><unit_amount_in_cents><USD>200</USD></unit_amount_in_cents></add_on>"},
		map[string]interface{}{"struct": AddOn{AccountingCode: "abc123"}, "xml": "<add_on><accounting_code>abc123</accounting_code></add_on>"},
	}

	for _, s := range suite {
		buf := new(bytes.Buffer)
		err := xml.NewEncoder(buf).Encode(s["struct"])
		if err != nil {
			t.Errorf("TestAddOnEncoding Error: %s", err)
		}

		if buf.String() != s["xml"] {
			t.Errorf("TestAddOnEncoding Error: Expected %s, given %s", s["xml"], buf.String())
		}
	}
}
