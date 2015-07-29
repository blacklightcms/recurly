package recurly

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net"
	"net/http"
	"reflect"
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

func TestBillingType(t *testing.T) {
	b := Billing{
		FirstSix: 411111,
		LastFour: 1111,
		Month:    11,
		Year:     2020,
	}

	b2 := Billing{
		NameOnAccount: "Acme, Inc",
		RoutingNumber: 123456780,
		AccountNumber: 111111111,
	}

	b3 := Billing{}

	if b.Type() != "card" {
		t.Errorf("TestBillingType Error: Expected card billing info to return card, given %s", b.Type())
	}

	if b2.Type() != "bank" {
		t.Errorf("TestBillingType Error: Expected bank billing info to return bank, given %s", b.Type())
	}

	if b3.Type() != "" {
		t.Errorf("TestingBillingtype Error: Expected billing info that is neither card or bank to return \"\", given %s", b3.Type())
	}
}

func TestBillingGet(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/1/billing_info", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("TestBillingGet Error: Expected %s request, given %s", "GET", r.Method)
		}
		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?>
		<billing_info href="http://api.test.host/v2/accounts/1/billing_info" type="credit_card">
			<account href="http://api.test.host/v2/accounts/1"/>
			<first_name>Verena</first_name>
			<last_name>Example</last_name>
			<company nil="nil"></company>
			<address1>123 Main St.</address1>
			<address2 nil="nil"></address2>
			<city>San Francisco</city>
			<state>CA</state>
			<zip>94105</zip>
			<country>US</country>
			<phone nil="nil"></phone>
			<vat_number>US1234567890</vat_number>
			<ip_address>127.0.0.1</ip_address>
			<ip_address_country>US</ip_address_country>
			<card_type>Visa</card_type>
			<year type="integer">2015</year>
			<month type="integer">11</month>
			<first_six>411111</first_six>
			<last_four>1111</last_four>
		</billing_info>`)
	})

	r, b, err := client.Billing.Get("1")
	if err != nil {
		t.Errorf("TestBillingGet Error: Error occured making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestBillingGet Error: Expected get billing info to return OK")
	}

	expected := Billing{
		XMLName:          xml.Name{Local: "billing_info"},
		FirstName:        "Verena",
		LastName:         "Example",
		Address:          "123 Main St.",
		City:             "San Francisco",
		State:            "CA",
		Zip:              "94105",
		Country:          "US",
		VATNumber:        "US1234567890",
		IPAddress:        net.ParseIP("127.0.0.1"),
		IPAddressCountry: "US",
		CardType:         "Visa",
		Year:             2015,
		Month:            11,
		FirstSix:         411111,
		LastFour:         1111,
	}

	if !reflect.DeepEqual(expected, b) {
		t.Errorf("TestBillingGet Error: expected billing to equal %#v, given %#v", expected, b)
	}
}

func TestBillingCreateWithToken(t *testing.T) {
	setup()
	defer teardown()

	token := "tok-woueh7LtsBKs8sE20"
	mux.HandleFunc("/v2/accounts/1/billing_info", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("TestBillingCreateWithToken Error: Expected %s request, given %s", "POST", r.Method)
		}
		given := new(bytes.Buffer)
		given.ReadFrom(r.Body)
		expected := fmt.Sprintf("<billing_info><token_id>%s</token_id></billing_info>", token)
		if expected != given.String() {
			t.Errorf("TestBillingCreateWithToken Error: Expected request body of %s, given %s", expected, given.String())
		}

		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?>
		<billing_info href="http://api.test.host/v2/accounts/1/billing_info" type="credit_card">
			<account href="http://api.test.host/v2/accounts/1"/>
			<first_name>Verena</first_name>
			<last_name>Example</last_name>
			<company nil="nil"></company>
			<address1>123 Main St.</address1>
			<address2 nil="nil"></address2>
			<city>San Francisco</city>
			<state>CA</state>
			<zip>94105</zip>
			<country>US</country>
			<phone nil="nil"></phone>
			<vat_number>US1234567890</vat_number>
			<ip_address>127.0.0.1</ip_address>
			<ip_address_country>US</ip_address_country>
			<card_type>Visa</card_type>
			<year type="integer">2015</year>
			<month type="integer">11</month>
			<first_six>411111</first_six>
			<last_four>1111</last_four>
		</billing_info>`)
	})

	r, b, err := client.Billing.CreateWithToken("1", token)
	if err != nil {
		t.Errorf("TestBillingCreateWithToken Error: Error occured making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestBillingCreateWithToken Error: Expected creating billing info to return OK")
	}

	expected := Billing{
		XMLName:          xml.Name{Local: "billing_info"},
		FirstName:        "Verena",
		LastName:         "Example",
		Address:          "123 Main St.",
		City:             "San Francisco",
		State:            "CA",
		Zip:              "94105",
		Country:          "US",
		VATNumber:        "US1234567890",
		IPAddress:        net.ParseIP("127.0.0.1"),
		IPAddressCountry: "US",
		CardType:         "Visa",
		Year:             2015,
		Month:            11,
		FirstSix:         411111,
		LastFour:         1111,
	}

	if !reflect.DeepEqual(expected, b) {
		t.Errorf("TestBillingCreateWithToken Error: expected create billing response to equal %#v, given %#v", expected, b)
	}
}

func TestBillingCreateWithCC(t *testing.T) {
	t.Skip("TestBillingCreatedWithCC Notice: Skipping test")
}

func TestBillingCreateWithBankAccount(t *testing.T) {
	t.Skip("TestBillingCreatedWithCC Notice: Skipping test")
}

func TestBillingUpdateWithToken(t *testing.T) {
	setup()
	defer teardown()

	token := "tok-UgLus845alBogoKRsiGw92vzos"
	mux.HandleFunc("/v2/accounts/abceasf/billing_info", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("TestBillingUpdateWithToken Error: Expected %s request, given %s", "PUT", r.Method)
		}
		given := new(bytes.Buffer)
		given.ReadFrom(r.Body)
		expected := fmt.Sprintf("<billing_info><token_id>%s</token_id></billing_info>", token)
		if expected != given.String() {
			t.Errorf("TestBillingUpdateWithToken Error: Expected request body of %s, given %s", expected, given.String())
		}

		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?><billing_info></billing_info>`)
	})

	r, _, err := client.Billing.UpdateWithToken("abceasf", token)
	if err != nil {
		t.Errorf("TestBillingUpdateWithToken Error: Error occured making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestBillingUpdateWithToken Error: Expected updating billing info to return OK")
	}
}

func TestBillingUpdateWithCC(t *testing.T) {
	t.Skip("TestBillingUpdatedWithCC Notice: Skipping test")
}

func TestBillingUpdateWithBankAccount(t *testing.T) {
	t.Skip("TestBillingUpdatedWithCC Notice: Skipping test")
}

func TestClearBilling(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/account@example.com/billing_info", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("TestClearBilling Error: Expected %s request, given %s", "DELETE", r.Method)
		}
		rw.WriteHeader(204)
	})

	r, err := client.Billing.Clear("account@example.com")
	if err != nil {
		t.Errorf("TestClearBilling Error: Error occured making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestClearBilling Error: Expected deleting billing_info to return OK")
	}
}
