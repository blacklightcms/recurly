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
		map[string]interface{}{"struct": Billing{FirstName: "Verena", LastName: "Example"}, "xml": "<billing_info><first_name>Verena</first_name><last_name>Example</last_name></billing_info>"},
		map[string]interface{}{"struct": Billing{Address: "123 Main St."}, "xml": "<billing_info><address1>123 Main St.</address1></billing_info>"},
		map[string]interface{}{"struct": Billing{Address2: "Unit A"}, "xml": "<billing_info><address2>Unit A</address2></billing_info>"},
		map[string]interface{}{"struct": Billing{City: "San Francisco"}, "xml": "<billing_info><city>San Francisco</city></billing_info>"},
		map[string]interface{}{"struct": Billing{State: "CA"}, "xml": "<billing_info><state>CA</state></billing_info>"},
		map[string]interface{}{"struct": Billing{Zip: "94105"}, "xml": "<billing_info><zip>94105</zip></billing_info>"},
		map[string]interface{}{"struct": Billing{Country: "US"}, "xml": "<billing_info><country>US</country></billing_info>"},
		map[string]interface{}{"struct": Billing{Phone: "555-555-5555"}, "xml": "<billing_info><phone>555-555-5555</phone></billing_info>"},
		map[string]interface{}{"struct": Billing{VATNumber: "abc"}, "xml": "<billing_info><vat_number>abc</vat_number></billing_info>"},
		map[string]interface{}{"struct": Billing{IPAddress: net.ParseIP("127.0.0.1")}, "xml": "<billing_info><ip_address>127.0.0.1</ip_address></billing_info>"},
		map[string]interface{}{"struct": Billing{Number: 4111111111111111, Month: 5, Year: 2020, VerificationValue: 111}, "xml": "<billing_info><number>4111111111111111</number><month>5</month><year>2020</year><verification_value>111</verification_value></billing_info>"},
		map[string]interface{}{"struct": Billing{RoutingNumber: "065400137", AccountNumber: "0123456789", AccountType: "checking"}, "xml": "<billing_info><routing_number>065400137</routing_number><account_number>0123456789</account_number><account_type>checking</account_type></billing_info>"},
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
		RoutingNumber: "123456780",
		AccountNumber: "111111111",
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
		t.Errorf("TestBillingGet Error: Error occurred making API call. Err: %s", err)
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
		t.Errorf("TestBillingCreateWithToken Error: Error occurred making API call. Err: %s", err)
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
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/1/billing_info", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("TestBillingCreateWithCC Error: Expected %s request, given %s", "POST", r.Method)
		}
		given := new(bytes.Buffer)
		given.ReadFrom(r.Body)
		expected := "<billing_info><first_name>Verena</first_name><last_name>Example</last_name><address1>123 Main St.</address1><city>San Francisco</city><state>CA</state><zip>94105</zip><country>US</country><number>4111111111111111</number><month>10</month><year>2020</year></billing_info>"
		if expected != given.String() {
			t.Errorf("TestBillingCreateWithCC Error: Expected request body of %s, given %s", expected, given.String())
		}

		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?><billing_info></billing_info>`)
	})

	r, _, err := client.Billing.Create("1", Billing{
		FirstName: "Verena",
		LastName:  "Example",
		Address:   "123 Main St.",
		City:      "San Francisco",
		State:     "CA",
		Zip:       "94105",
		Country:   "US",
		Number:    4111111111111111,
		Month:     10,
		Year:      2020,
	})

	if err != nil {
		t.Errorf("TestBillingCreateWithCC Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestBillingCreateWithCC Error: Expected creating billing info to return OK")
	}
}

func TestBillingCreateWithBankAccount(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/134/billing_info", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("TestBillingCreateWithBankAccount Error: Expected %s request, given %s", "POST", r.Method)
		}
		given := new(bytes.Buffer)
		given.ReadFrom(r.Body)
		expected := "<billing_info><first_name>Verena</first_name><last_name>Example</last_name><address1>123 Main St.</address1><city>San Francisco</city><state>CA</state><zip>94105</zip><country>US</country><name_on_account>Acme, Inc</name_on_account><routing_number>123456780</routing_number><account_number>111111111</account_number><account_type>checking</account_type></billing_info>"
		if expected != given.String() {
			t.Errorf("TestBillingCreateWithBankAccount Error: Expected request body of %s, given %s", expected, given.String())
		}

		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?><billing_info></billing_info>`)
	})

	r, _, err := client.Billing.Create("134", Billing{
		FirstName:     "Verena",
		LastName:      "Example",
		Address:       "123 Main St.",
		City:          "San Francisco",
		State:         "CA",
		Zip:           "94105",
		Country:       "US",
		NameOnAccount: "Acme, Inc",
		RoutingNumber: "123456780",
		AccountNumber: "111111111",
		AccountType:   "checking",
	})

	if err != nil {
		t.Errorf("TestBillingCreateWithBankAccount Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestBillingCreateWithBankAccount Error: Expected creating billing info to return OK")
	}
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
		t.Errorf("TestBillingUpdateWithToken Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestBillingUpdateWithToken Error: Expected updating billing info to return OK")
	}
}

func TestBillingUpdateWithInvalidToken(t *testing.T) {
	setup()
	defer teardown()

	token := "tok-UgLus845alBogoKRsiGw92vzos"
	mux.HandleFunc("/v2/accounts/abceasf/billing_info", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("TestBillingUpdateWithInvalidToken Error: Expected %s request, given %s", "PUT", r.Method)
		}

		rw.WriteHeader(404)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?><error><symbol>token_invalid</symbol><description>Token is either invalid or expired</description></error>`)
	})

	r, _, err := client.Billing.UpdateWithToken("abceasf", token)
	if err != nil {
		t.Errorf("TestBillingUpdateWithInvalidToken Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsOK() {
		t.Fatal("TestBillingUpdateWithInvalidToken Error: Expected updating billing info with invalid token to return error")
	}

	if len(r.Errors) == 0 || r.Errors[0].Symbol != "token_invalid" {
		t.Errorf("TestBillingUpdateWithInvalidToken Error: Error response not parsed properly")
	}
}

func TestBillingUpdateWithCC(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/1/billing_info", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("TestBillingUpdateWithCC Error: Expected %s request, given %s", "PUT", r.Method)
		}
		given := new(bytes.Buffer)
		given.ReadFrom(r.Body)
		expected := "<billing_info><first_name>Verena</first_name><last_name>Example</last_name><address1>123 Main St.</address1><city>San Francisco</city><state>CA</state><zip>94105</zip><country>US</country><number>4111111111111111</number><month>10</month><year>2020</year></billing_info>"
		if expected != given.String() {
			t.Errorf("TestBillingUpdateWithCC Error: Expected request body of %s, given %s", expected, given.String())
		}

		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?><billing_info></billing_info>`)
	})

	r, _, err := client.Billing.Update("1", Billing{
		FirstName: "Verena",
		LastName:  "Example",
		Address:   "123 Main St.",
		City:      "San Francisco",
		State:     "CA",
		Zip:       "94105",
		Country:   "US",
		Number:    4111111111111111,
		Month:     10,
		Year:      2020,

		// Add additional fields that should be removed
		Token:             "abc",
		IPAddressCountry:  "US",
		FirstSix:          411111,
		LastFour:          1111,
		CardType:          "visa",
		PaypalAgreementID: "ppl",
		AmazonAgreementID: "asdfb",
	})

	if err != nil {
		t.Errorf("TestBillingUpdateWithCC Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestBillingUpdateWithCC Error: Expected creating billing info to return OK")
	}
}

func TestBillingUpdateWithBankAccount(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/134/billing_info", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("TestBillingUpdateWithBankAccount Error: Expected %s request, given %s", "PUT", r.Method)
		}
		given := new(bytes.Buffer)
		given.ReadFrom(r.Body)
		expected := "<billing_info><first_name>Verena</first_name><last_name>Example</last_name><address1>123 Main St.</address1><city>San Francisco</city><state>CA</state><zip>94105</zip><country>US</country><name_on_account>Acme, Inc</name_on_account><routing_number>123456780</routing_number><account_number>111111111</account_number><account_type>checking</account_type></billing_info>"
		if expected != given.String() {
			t.Errorf("TestBillingUpdateWithBankAccount Error: Expected request body of %s, given %s", expected, given.String())
		}

		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?><billing_info></billing_info>`)
	})

	r, _, err := client.Billing.Update("134", Billing{
		FirstName:     "Verena",
		LastName:      "Example",
		Address:       "123 Main St.",
		City:          "San Francisco",
		State:         "CA",
		Zip:           "94105",
		Country:       "US",
		NameOnAccount: "Acme, Inc",
		RoutingNumber: "123456780",
		AccountNumber: "111111111",
		AccountType:   "checking",

		// Add additional fields that should be removed
		Token:             "abc",
		IPAddressCountry:  "US",
		FirstSix:          111111,
		LastFour:          1111,
		CardType:          "visa",
		PaypalAgreementID: "ppl",
		AmazonAgreementID: "asdfb",
	})

	if err != nil {
		t.Errorf("TestBillingUpdateWithBankAccount Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestBillingUpdateWithBankAccount Error: Expected creating billing info to return OK")
	}
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
		t.Errorf("TestClearBilling Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestClearBilling Error: Expected deleting billing_info to return OK")
	}
}
