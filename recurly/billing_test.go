package recurly

import (
	"bytes"
	"encoding/xml"
	"testing"
)

// TestBillingEncoding ensures structs are encoded to XML properly.
// Because Recurly supports partial updates, it's important that only defined
// fields are handled properly -- including types like booleans and integers which
// have zero values that we want to send.
func TestBillingEncoding(t *testing.T) {
	suite := []map[string]interface{}{
		map[string]interface{}{"struct": Billing{}, "xml": "<billing_info></billing_info>"},
		map[string]interface{}{"struct": Billing{Token: "507c7f79bcf86cd7994f6c0e"}, "xml": "<billing_info><token_id>507c7f79bcf86cd7994f6c0e</token_id></billing_info>"},

		// @todo test bank account and credit card fields when support for those updates is added.
	}

	for _, s := range suite {
		buf := new(bytes.Buffer)
		err := xml.NewEncoder(buf).Encode(s["struct"])
		if err != nil {
			t.Errorf("TestBillingEncoding Error: %s", err)
		}

		if buf.String() != s["xml"] {
			t.Errorf("TestBillingEncoding Error: Expected %s, given %s", s["xml"], buf.String())
		}
	}
}
