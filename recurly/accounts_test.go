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

// TestAccountEncoding ensures structs are encoded to XML properly.
// Because Recurly supports partial updates, it's important that only defined
// fields are handled properly -- including types like booleans and integers which
// have zero values that we want to send.
func TestAccountEncoding(t *testing.T) {
	suite := []map[string]interface{}{
		map[string]interface{}{"struct": Account{}, "xml": "<account></account>"},
		map[string]interface{}{"struct": Account{Code: "abc"}, "xml": "<account><account_code>abc</account_code></account>"},
		map[string]interface{}{"struct": Account{State: "active"}, "xml": "<account><state>active</state></account>"},
		map[string]interface{}{"struct": Account{Email: "me@example.com"}, "xml": "<account><email>me@example.com</email></account>"},
		map[string]interface{}{"struct": Account{FirstName: "Larry"}, "xml": "<account><first_name>Larry</first_name></account>"},
		map[string]interface{}{"struct": Account{LastName: "Larrison"}, "xml": "<account><last_name>Larrison</last_name></account>"},
		map[string]interface{}{"struct": Account{FirstName: "Larry", LastName: "Larrison"}, "xml": "<account><first_name>Larry</first_name><last_name>Larrison</last_name></account>"},
		map[string]interface{}{"struct": Account{CompanyName: "Acme, Inc"}, "xml": "<account><company_name>Acme, Inc</company_name></account>"},
		map[string]interface{}{"struct": Account{VATNumber: "123456789"}, "xml": "<account><vat_number>123456789</vat_number></account>"},
		map[string]interface{}{"struct": Account{TaxExempt: NewBool(true)}, "xml": "<account><tax_exempt>true</tax_exempt></account>"},
		map[string]interface{}{"struct": Account{TaxExempt: NewBool(false)}, "xml": "<account><tax_exempt>false</tax_exempt></account>"},
		map[string]interface{}{"struct": Account{AcceptLanguage: "en_US"}, "xml": "<account><accept_language>en_US</accept_language></account>"},
		map[string]interface{}{"struct": Account{FirstName: "Larry", Address: Address{Address: "123 Main St.", City: "San Francisco", State: "CA", Zip: "94105", Country: "US"}}, "xml": "<account><first_name>Larry</first_name><address><address1>123 Main St.</address1><city>San Francisco</city><state>CA</state><zip>94105</zip><country>US</country></address></account>"},
		map[string]interface{}{"struct": Account{Code: "test@example.com", BillingInfo: &Billing{Token: "507c7f79bcf86cd7994f6c0e"}}, "xml": "<account><account_code>test@example.com</account_code><billing_info><token_id>507c7f79bcf86cd7994f6c0e</token_id></billing_info></account>"},
		map[string]interface{}{"struct": Address{}, "xml": ""},
		map[string]interface{}{"struct": Address{Address: "123 Main St."}, "xml": "<address><address1>123 Main St.</address1></address>"},
		map[string]interface{}{"struct": Address{Address2: "Unit A"}, "xml": "<address><address2>Unit A</address2></address>"},
		map[string]interface{}{"struct": Address{City: "San Francisco"}, "xml": "<address><city>San Francisco</city></address>"},
		map[string]interface{}{"struct": Address{State: "CA"}, "xml": "<address><state>CA</state></address>"},
		map[string]interface{}{"struct": Address{Zip: "94105"}, "xml": "<address><zip>94105</zip></address>"},
		map[string]interface{}{"struct": Address{Country: "US"}, "xml": "<address><country>US</country></address>"},
		map[string]interface{}{"struct": Address{Phone: "555-555-5555"}, "xml": "<address><phone>555-555-5555</phone></address>"},
	}

	for _, s := range suite {
		buf := new(bytes.Buffer)
		err := xml.NewEncoder(buf).Encode(s["struct"])
		if err != nil {
			t.Errorf("TestAccountEncoding Error: %s", err)
		}

		if buf.String() != s["xml"] {
			t.Errorf("TestAccountEncoding Error: Expected %s, given %s", s["xml"], buf.String())
		}
	}
}

func TestAccountList(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("TestAccountList Error: Expected %s request, given %s", "GET", r.Method)
		}
		rw.Header().Set("Link", `<https://your-subdomain.recurly.com/v2/accounts?cursor=1304958672>; rel="next"`)
		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?>
		<accounts>
			<account href="https://your-subdomain.recurly.com/v2/accounts/1">
			  <adjustments href="https://your-subdomain.recurly.com/v2/accounts/1/adjustments"/>
			  <billing_info href="https://your-subdomain.recurly.com/v2/accounts/1/billing_info"/>
			  <invoices href="https://your-subdomain.recurly.com/v2/accounts/1/invoices"/>
			  <redemption href="https://your-subdomain.recurly.com/v2/accounts/1/redemption"/>
			  <subscriptions href="https://your-subdomain.recurly.com/v2/accounts/1/subscriptions"/>
			  <transactions href="https://your-subdomain.recurly.com/v2/accounts/1/transactions"/>
			  <account_code>1</account_code>
			  <state>active</state>
			  <username nil="nil"></username>
			  <email>verena@example.com</email>
			  <first_name>Verena</first_name>
			  <last_name>Example</last_name>
			  <company_name></company_name>
			  <vat_number nil="nil"></vat_number>
			  <tax_exempt type="boolean">false</tax_exempt>
			  <address>
			    <address1>123 Main St.</address1>
			    <address2 nil="nil"></address2>
			    <city>San Francisco</city>
			    <state>CA</state>
			    <zip>94105</zip>
			    <country>US</country>
			    <phone nil="nil"></phone>
			  </address>
			  <accept_language nil="nil"></accept_language>
			  <hosted_login_token>a92468579e9c4231a6c0031c4716c01d</hosted_login_token>
			  <created_at type="datetime">2011-10-25T12:00:00Z</created_at>
			</account>
		</accounts>`)
	})

	r, accounts, err := client.Accounts.List(nil)
	if err != nil {
		t.Errorf("TestAccountList Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestAccountList Error: Expected list accounts to return OK")
	}

	if len(accounts) != 1 {
		t.Fatalf("TestAccountList Error: Expected 1 account returned, given %d", len(accounts))
	}

	if r.Prev() != "" {
		t.Errorf("TestAccountListPagination Error: Expected prev cursor to be empty, given %s", r.Prev())
	}

	if r.Next() != "1304958672" {
		t.Errorf("TestAccountListPagination Error: Expected next cursor to equal %s, given %s", "1318388868", r.Next())
	}

	ts, _ := time.Parse(datetimeFormat, "2011-10-25T12:00:00Z")
	for _, given := range accounts {
		expected := Account{
			XMLName:   xml.Name{Local: "account"},
			Code:      "1",
			State:     "active",
			Email:     "verena@example.com",
			FirstName: "Verena",
			LastName:  "Example",
			TaxExempt: NewBool(false),
			Address: Address{
				Address: "123 Main St.",
				City:    "San Francisco",
				State:   "CA",
				Zip:     "94105",
				Country: "US",
			},
			HostedLoginToken: "a92468579e9c4231a6c0031c4716c01d",
			CreatedAt:        NewTime(ts),
		}

		if !reflect.DeepEqual(expected, given) {
			t.Errorf("TestAccountList Error: expected account to equal %#v, given %#v", expected, given)
		}
	}
}

func TestAccountListPagination(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("TestAccountList Error: Expected %s request, given %s", "GET", r.Method)
		}
		rw.Header().Set("Link", `<https://your-subdomain.recurly.com/v2/transactions>; rel="start",
  <https://your-subdomain.recurly.com/v2/transactions?cursor=-1318344434>; rel="prev",
  <https://your-subdomain.recurly.com/v2/transactions?cursor=1318388868>; rel="next"`)
		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?><accounts></accounts>`)
	})

	r, _, err := client.Accounts.List(Params{"cursor": "12345"})
	if err != nil {
		t.Errorf("TestAccountListPagination Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestAccountListPagination Error: Expected list accounts to return OK")
	}

	if r.Prev() != "-1318344434" {
		t.Errorf("TestAccountListPagination Error: Expected prev cursor to equal %s, given %s", "-1318344434", r.Prev())
	}

	if r.Next() != "1318388868" {
		t.Errorf("TestAccountListPagination Error: Expected next cursor to equal %s, given %s", "1318388868", r.Next())
	}
}

func TestGetAccount(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/1", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("TestGetAccount Error: Expected %s request, given %s", "GET", r.Method)
		}
		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?>
			<account href="https://your-subdomain.recurly.com/v2/accounts/1">
			  <adjustments href="https://your-subdomain.recurly.com/v2/accounts/1/adjustments"/>
			  <billing_info href="https://your-subdomain.recurly.com/v2/accounts/1/billing_info"/>
			  <invoices href="https://your-subdomain.recurly.com/v2/accounts/1/invoices"/>
			  <redemption href="https://your-subdomain.recurly.com/v2/accounts/1/redemption"/>
			  <subscriptions href="https://your-subdomain.recurly.com/v2/accounts/1/subscriptions"/>
			  <transactions href="https://your-subdomain.recurly.com/v2/accounts/1/transactions"/>
			  <account_code>1</account_code>
			  <state>active</state>
			  <username nil="nil"></username>
			  <email>verena@example.com</email>
			  <first_name>Verena</first_name>
			  <last_name>Example</last_name>
			  <company_name></company_name>
			  <vat_number nil="nil"></vat_number>
			  <tax_exempt type="boolean">false</tax_exempt>
			  <address>
			    <address1>123 Main St.</address1>
			    <address2 nil="nil"></address2>
			    <city>San Francisco</city>
			    <state>CA</state>
			    <zip>94105</zip>
			    <country>US</country>
			    <phone nil="nil"></phone>
			  </address>
			  <accept_language nil="nil"></accept_language>
			  <hosted_login_token>a92468579e9c4231a6c0031c4716c01d</hosted_login_token>
			  <created_at type="datetime">2011-10-25T12:00:00Z</created_at>
			</account>`)
	})

	r, a, err := client.Accounts.Get("1")
	if err != nil {
		t.Errorf("TestGetAccount Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestGetAccount Error: Expected get accounts to return OK")
	}

	ts, _ := time.Parse(datetimeFormat, "2011-10-25T12:00:00Z")
	expected := Account{
		XMLName:   xml.Name{Local: "account"},
		Code:      "1",
		State:     "active",
		Email:     "verena@example.com",
		FirstName: "Verena",
		LastName:  "Example",
		TaxExempt: NewBool(false),
		Address: Address{
			Address: "123 Main St.",
			City:    "San Francisco",
			State:   "CA",
			Zip:     "94105",
			Country: "US",
		},
		HostedLoginToken: "a92468579e9c4231a6c0031c4716c01d",
		CreatedAt:        NewTime(ts),
	}

	if !reflect.DeepEqual(expected, a) {
		t.Errorf("TestGetAccount Error: expected account to equal %#v, given %#v", expected, a)
	}
}

func TestCreateAccount(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("TestCreateAccount Error: Expected %s request, given %s", "POST", r.Method)
		}
		rw.WriteHeader(201)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?><account></account>`)
	})

	r, _, err := client.Accounts.Create(Account{})
	if err != nil {
		t.Errorf("TestCreateAccount Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestCreateAccount Error: Expected create account to return OK")
	}
}

func TestUpdateAccount(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/245", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("TestUpdateAccount Error: Expected %s request, given %s", "PUT", r.Method)
		}
		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?><account></account>`)
	})

	r, _, err := client.Accounts.Update("245", Account{})
	if err != nil {
		t.Errorf("TestUpdateAccount Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestUpdateAccount Error: Expected update account to return OK")
	}
}

func TestCloseAccount(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/5322", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("TestCloseAccount Error: Expected %s request, given %s", "DELETE", r.Method)
		}
		rw.WriteHeader(204)
	})

	r, err := client.Accounts.Close("5322")
	if err != nil {
		t.Errorf("TestCloseAccount Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestCloseAccount Error: Expected close account to return OK")
	}
}

func TestReopenAccount(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/5322/reopen", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("TestReopenAccount Error: Expected %s request, given %s", "PUT", r.Method)
		}
		rw.WriteHeader(204)
	})

	r, err := client.Accounts.Reopen("5322")
	if err != nil {
		t.Errorf("TestReopenAccount Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestReopenAccount Error: Expected reopen account to return OK")
	}
}

func TestAccountListNotes(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/abcd@example.com/notes", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("TestAccountListNotes Error: Expected %s request, given %s", "GET", r.Method)
		}
		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?>
			<notes type="array">
			  <note>
			    <account href="https://your-subdomain.recurly.com/v2/accounts/abcd@example.com"/>
			    <message>This is my second note</message>
			    <created_at type="datetime">2013-05-14T18:53:04Z</created_at>
			  </note>
			  <note>
			    <account href="https://your-subdomain.recurly.com/v2/accounts/abcd@example.com"/>
			    <message>This is my first note</message>
			    <created_at type="datetime">2013-05-14T18:52:50Z</created_at>
			  </note>
			</notes>`)
	})

	r, notes, err := client.Accounts.ListNotes("abcd@example.com")
	if err != nil {
		t.Errorf("TestAccountListNotes Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestAccountListNotes Error: Expected list notes to return OK")
	}

	if len(notes) != 2 {
		t.Fatalf("TestAccountListNotes Error: Expected 2 notes returned, given %d", len(notes))
	}

	ts1, _ := time.Parse(datetimeFormat, "2013-05-14T18:52:50Z")
	ts2, _ := time.Parse(datetimeFormat, "2013-05-14T18:53:04Z")
	expected := []Note{
		Note{
			XMLName:   xml.Name{Local: "note"},
			Message:   "This is my second note",
			CreatedAt: ts2,
		},
		Note{
			XMLName:   xml.Name{Local: "note"},
			Message:   "This is my first note",
			CreatedAt: ts1,
		},
	}

	if !reflect.DeepEqual(expected, notes) {
		t.Errorf("TestAccountListNotes Error: expected notes to equal %#v, notes %#v", expected, notes)
	}
}
