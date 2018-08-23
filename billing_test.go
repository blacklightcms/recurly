package recurly_test

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/launchpadcentral/recurly"
	"github.com/google/go-cmp/cmp"
)

// TestBillingEncoding ensures structs are encoded to XML properly.
// Because Recurly supports partial updates, it's important that only defined
// fields are handled properly -- including types like booleans and integers which
// have zero values that we want to send.
func TestBilling_Encoding(t *testing.T) {
	tests := []struct {
		v        recurly.Billing
		expected string
	}{
		{v: recurly.Billing{}, expected: "<billing_info></billing_info>"},
		{v: recurly.Billing{Token: "507c7f79bcf86cd7994f6c0e"}, expected: "<billing_info><token_id>507c7f79bcf86cd7994f6c0e</token_id></billing_info>"},
		{v: recurly.Billing{FirstName: "Verena", LastName: "Example"}, expected: "<billing_info><first_name>Verena</first_name><last_name>Example</last_name></billing_info>"},
		{v: recurly.Billing{Address: "123 Main St."}, expected: "<billing_info><address1>123 Main St.</address1></billing_info>"},
		{v: recurly.Billing{Address2: "Unit A"}, expected: "<billing_info><address2>Unit A</address2></billing_info>"},
		{v: recurly.Billing{City: "San Francisco"}, expected: "<billing_info><city>San Francisco</city></billing_info>"},
		{v: recurly.Billing{State: "CA"}, expected: "<billing_info><state>CA</state></billing_info>"},
		{v: recurly.Billing{Zip: "94105"}, expected: "<billing_info><zip>94105</zip></billing_info>"},
		{v: recurly.Billing{Country: "US"}, expected: "<billing_info><country>US</country></billing_info>"},
		{v: recurly.Billing{Phone: "555-555-5555"}, expected: "<billing_info><phone>555-555-5555</phone></billing_info>"},
		{v: recurly.Billing{VATNumber: "abc"}, expected: "<billing_info><vat_number>abc</vat_number></billing_info>"},
		{v: recurly.Billing{IPAddress: net.ParseIP("127.0.0.1")}, expected: "<billing_info><ip_address>127.0.0.1</ip_address></billing_info>"},
		{v: recurly.Billing{Number: 4111111111111111, Month: 5, Year: 2020, VerificationValue: 111}, expected: "<billing_info><number>4111111111111111</number><month>5</month><year>2020</year><verification_value>111</verification_value></billing_info>"},
		{v: recurly.Billing{RoutingNumber: "065400137", AccountNumber: "0123456789", AccountType: "checking"}, expected: "<billing_info><routing_number>065400137</routing_number><account_number>0123456789</account_number><account_type>checking</account_type></billing_info>"},
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

func TestBilling_Type(t *testing.T) {
	b0 := recurly.Billing{
		FirstSix: 411111,
		LastFour: "1111",
		Month:    11,
		Year:     2020,
	}

	b1 := recurly.Billing{
		NameOnAccount: "Acme, Inc",
		RoutingNumber: "123456780",
		AccountNumber: "111111111",
	}

	var b2 recurly.Billing

	if b0.Type() != "card" {
		t.Fatalf("unexpected type: %s", b0.Type())
	} else if b1.Type() != "bank" {
		t.Fatalf("unexpected type: %s", b1.Type())
	} else if b2.Type() != "" {
		t.Fatalf("unexpected type: %s", b2.Type())
	}
}

func TestBilling_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/1/billing_info", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
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

	resp, b, err := client.Billing.Get("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected get billing info to return OK")
	} else if diff := cmp.Diff(b, &recurly.Billing{
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
		LastFour:         "1111",
	}); diff != "" {
		t.Fatal(diff)
	}
}

// ACH customers may not have billing info. This asserts that nil values for
// many of the fields are safely ignored without parse errors.
func TestBilling_Get_ACH(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/1/billing_info", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
            <billing_info type="ach">
                <first_name>Verena</first_name>
                <last_name>Example</last_name>
                <address1 nil="nil"></address1>
                <address2 nil="nil"></address2>
                <city nil="nil"></city>
                <state nil="nil"></state>
                <zip nil="nil"></zip>
                <country nil="nil"></country>
                <phone nil="nil"></phone>
                <vat_number nil="nil"></vat_number>
                <account_type nil="nil"></account_type>
                <last_four nil="nil"></last_four>
                <routing_number nil="nil"></routing_number>
            </billing_info>`)
	})

	resp, b, err := client.Billing.Get("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected get billing info to return OK")
	} else if diff := cmp.Diff(b, &recurly.Billing{
		XMLName:   xml.Name{Local: "billing_info"},
		FirstName: "Verena",
		LastName:  "Example",
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestBilling_Get_ErrNotFound(t *testing.T) {
	setup()
	defer teardown()

	var invoked bool
	mux.HandleFunc("/v2/accounts/1/billing_info", func(w http.ResponseWriter, r *http.Request) {
		invoked = true
		w.WriteHeader(http.StatusNotFound)
	})

	_, billing, err := client.Billing.Get("1")
	if !invoked {
		t.Fatal("handler not invoked")
	} else if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if billing != nil {
		t.Fatalf("expected billing to be nil: %#v", billing)
	}
}

func TestBilling_Create_WithToken(t *testing.T) {
	setup()
	defer teardown()

	token := "tok-woueh7LtsBKs8sE20"
	mux.HandleFunc("/v2/accounts/1/billing_info", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if r.Method != "POST" {
			t.Fatalf("unexpected method: %s", r.Method)
		}

		var given bytes.Buffer
		given.ReadFrom(r.Body)
		expected := fmt.Sprintf("<billing_info><token_id>%s</token_id></billing_info>", token)
		if expected != given.String() {
			t.Fatalf("unexpected input: %v", given.String())
		}

		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
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

	resp, b, err := client.Billing.CreateWithToken("1", token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected creating billing info to return OK")
	} else if diff := cmp.Diff(b, &recurly.Billing{
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
		LastFour:         "1111",
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestBilling_Create_WithCC(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/1/billing_info", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("unexpected status: %s", r.Method)
		}
		var given bytes.Buffer
		given.ReadFrom(r.Body)
		expected := "<billing_info><first_name>Verena</first_name><last_name>Example</last_name><address1>123 Main St.</address1><city>San Francisco</city><state>CA</state><zip>94105</zip><country>US</country><number>4111111111111111</number><month>10</month><year>2020</year></billing_info>"
		if expected != given.String() {
			t.Fatalf("unexpected input: %v", given.String())
		}

		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><billing_info></billing_info>`)
	})

	resp, _, err := client.Billing.Create("1", recurly.Billing{
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
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected creating billing info to return OK")
	}
}

func TestBilling_Create_WithBankAccount(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/134/billing_info", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		var given bytes.Buffer
		given.ReadFrom(r.Body)
		expected := "<billing_info><first_name>Verena</first_name><last_name>Example</last_name><address1>123 Main St.</address1><city>San Francisco</city><state>CA</state><zip>94105</zip><country>US</country><name_on_account>Acme, Inc</name_on_account><routing_number>123456780</routing_number><account_number>111111111</account_number><account_type>checking</account_type></billing_info>"
		if expected != given.String() {
			t.Fatalf("unexpected input: %v", given.String())
		}

		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><billing_info></billing_info>`)
	})

	resp, _, err := client.Billing.Create("134", recurly.Billing{
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
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected creating billing info to return OK")
	}
}

func TestBilling_Update_WithToken(t *testing.T) {
	setup()
	defer teardown()

	token := "tok-UgLus845alBogoKRsiGw92vzos"
	mux.HandleFunc("/v2/accounts/abceasf/billing_info", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		var given bytes.Buffer
		given.ReadFrom(r.Body)
		expected := fmt.Sprintf("<billing_info><token_id>%s</token_id></billing_info>", token)
		if expected != given.String() {
			t.Fatalf("unexpected input: %v", given.String())
		}

		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><billing_info></billing_info>`)
	})

	resp, _, err := client.Billing.UpdateWithToken("abceasf", token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected updating billing info to return OK")
	}
}

func TestBilling_Update_InvalidToken(t *testing.T) {
	setup()
	defer teardown()

	token := "tok-UgLus845alBogoKRsiGw92vzos"
	mux.HandleFunc("/v2/accounts/abceasf/billing_info", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("unexpected method: %s", r.Method)
		}

		w.WriteHeader(404)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><error><symbol>token_invalid</symbol><description>Token is either invalid or expired</description></error>`)
	})

	resp, billing, err := client.Billing.UpdateWithToken("abceasf", token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if billing != nil {
		t.Fatalf("unexpected billing to be nil: %#v", billing)
	}

	if resp.IsOK() {
		t.Fatal("expected updating billing info with invalid token to return error")
	} else if diff := cmp.Diff(resp.Errors, []recurly.Error{
		{
			Symbol:  "token_invalid",
			Message: "Token is either invalid or expired",
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestBilling_Update_WithCC(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/1/billing_info", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		var given bytes.Buffer
		given.ReadFrom(r.Body)
		expected := "<billing_info><first_name>Verena</first_name><last_name>Example</last_name><address1>123 Main St.</address1><city>San Francisco</city><state>CA</state><zip>94105</zip><country>US</country><number>4111111111111111</number><month>10</month><year>2020</year></billing_info>"
		if expected != given.String() {
			t.Fatalf("unexpected input: %v", given.String())
		}

		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><billing_info></billing_info>`)
	})

	resp, _, err := client.Billing.Update("1", recurly.Billing{
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
		LastFour:          "1111",
		CardType:          "visa",
		PaypalAgreementID: "ppl",
		AmazonAgreementID: "asdfb",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected creating billing info to return OK")
	}
}

func TestBilling_Update_WithBankAccount(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/134/billing_info", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		var given bytes.Buffer
		given.ReadFrom(r.Body)
		expected := "<billing_info><first_name>Verena</first_name><last_name>Example</last_name><address1>123 Main St.</address1><city>San Francisco</city><state>CA</state><zip>94105</zip><country>US</country><name_on_account>Acme, Inc</name_on_account><routing_number>123456780</routing_number><account_number>111111111</account_number><account_type>checking</account_type></billing_info>"
		if expected != given.String() {
			t.Fatalf("unexpected input: %v", given.String())
		}

		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><billing_info></billing_info>`)
	})

	resp, _, err := client.Billing.Update("134", recurly.Billing{
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
		LastFour:          "1111",
		CardType:          "visa",
		PaypalAgreementID: "ppl",
		AmazonAgreementID: "asdfb",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected creating billing info to return OK")
	}
}

func TestBilling_Clear(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/account@example.com/billing_info", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(204)
	})

	resp, err := client.Billing.Clear("account@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected deleting billing_info to return OK")
	}
}
