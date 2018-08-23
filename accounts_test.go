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

// TestAccountEncoding ensures structs are encoded to XML properly.
// Because Recurly supports partial updates, it's important that only defined
// fields are handled properly -- including types like booleans and integers which
// have zero values that we want to send.
func TestAccounts_Encoding(t *testing.T) {
	tests := []struct {
		v        interface{}
		expected string
	}{
		{v: recurly.Account{}, expected: "<account></account>"},
		{v: recurly.Account{Code: "abc"}, expected: "<account><account_code>abc</account_code></account>"},
		{v: recurly.Account{State: "active"}, expected: "<account><state>active</state></account>"},
		{v: recurly.Account{Email: "me@example.com"}, expected: "<account><email>me@example.com</email></account>"},
		{v: recurly.Account{FirstName: "Larry"}, expected: "<account><first_name>Larry</first_name></account>"},
		{v: recurly.Account{LastName: "Larrison"}, expected: "<account><last_name>Larrison</last_name></account>"},
		{v: recurly.Account{FirstName: "Larry", LastName: "Larrison"}, expected: "<account><first_name>Larry</first_name><last_name>Larrison</last_name></account>"},
		{v: recurly.Account{CompanyName: "Acme, Inc"}, expected: "<account><company_name>Acme, Inc</company_name></account>"},
		{v: recurly.Account{VATNumber: "123456789"}, expected: "<account><vat_number>123456789</vat_number></account>"},
		{v: recurly.Account{TaxExempt: recurly.NewBool(true)}, expected: "<account><tax_exempt>true</tax_exempt></account>"},
		{v: recurly.Account{TaxExempt: recurly.NewBool(false)}, expected: "<account><tax_exempt>false</tax_exempt></account>"},
		{v: recurly.Account{AcceptLanguage: "en_US"}, expected: "<account><accept_language>en_US</accept_language></account>"},
		{v: recurly.Account{FirstName: "Larry", Address: recurly.Address{Address: "123 Main St.", City: "San Francisco", State: "CA", Zip: "94105", Country: "US"}}, expected: "<account><first_name>Larry</first_name><address><address1>123 Main St.</address1><city>San Francisco</city><state>CA</state><zip>94105</zip><country>US</country></address></account>"},
		{v: recurly.Account{Code: "test@example.com", BillingInfo: &recurly.Billing{Token: "507c7f79bcf86cd7994f6c0e"}}, expected: "<account><account_code>test@example.com</account_code><billing_info><token_id>507c7f79bcf86cd7994f6c0e</token_id></billing_info></account>"},
		{v: recurly.Account{HasPausedSubscription: true}, expected: "<account><has_paused_subscription>true</has_paused_subscription></account>"},
		{v: recurly.Address{}, expected: ""},
		{v: recurly.Address{Address: "123 Main St."}, expected: "<address><address1>123 Main St.</address1></address>"},
		{v: recurly.Address{Address2: "Unit A"}, expected: "<address><address2>Unit A</address2></address>"},
		{v: recurly.Address{City: "San Francisco"}, expected: "<address><city>San Francisco</city></address>"},
		{v: recurly.Address{State: "CA"}, expected: "<address><state>CA</state></address>"},
		{v: recurly.Address{Zip: "94105"}, expected: "<address><zip>94105</zip></address>"},
		{v: recurly.Address{Country: "US"}, expected: "<address><country>US</country></address>"},
		{v: recurly.Address{Phone: "555-555-5555"}, expected: "<address><phone>555-555-5555</phone></address>"},
	}

	for i, tt := range tests {
		var buf bytes.Buffer
		if err := xml.NewEncoder(&buf).Encode(tt.v); err != nil {
			t.Fatalf("TestAccountEncoding Error: %s", err)
		} else if buf.String() != tt.expected {
			t.Fatalf("(%d) unexpected value: %s", i, buf.String())
		}
	}
}

func TestAccounts_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.Header().Set("Link", `<https://your-subdomain.recurly.com/v2/accounts?cursor=1304958672>; rel="next"`)
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
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

	resp, accounts, err := client.Accounts.List(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected list accounts to return OK")
	} else if resp.Prev() != "" {
		t.Fatalf("unexpected cursor: %s", resp.Prev())
	} else if resp.Next() != "1304958672" {
		t.Fatalf("unexpected cursor: %s", resp.Next())
	}

	ts, _ := time.Parse(recurly.DateTimeFormat, "2011-10-25T12:00:00Z")
	if diff := cmp.Diff(accounts, []recurly.Account{recurly.Account{
		XMLName:   xml.Name{Local: "account"},
		Code:      "1",
		State:     "active",
		Email:     "verena@example.com",
		FirstName: "Verena",
		LastName:  "Example",
		TaxExempt: recurly.NewBool(false),
		Address: recurly.Address{
			Address: "123 Main St.",
			City:    "San Francisco",
			State:   "CA",
			Zip:     "94105",
			Country: "US",
		},
		HostedLoginToken: "a92468579e9c4231a6c0031c4716c01d",
		CreatedAt:        recurly.NewTime(ts),
	}}); diff != "" {
		t.Fatal(diff)
	}
}

func TestAccounts_List_Pagination(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.Header().Set("Link", `<https://your-subdomain.recurly.com/v2/transactions>; rel="start",
  <https://your-subdomain.recurly.com/v2/transactions?cursor=-1318344434>; rel="prev",
<https://your-subdomain.recurly.com/v2/transactions?cursor=1318388868>; rel="next"`)
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><accounts></accounts>`)
	})

	resp, _, err := client.Accounts.List(recurly.Params{"cursor": "12345"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected list accounts to return OK")
	} else if resp.Prev() != "-1318344434" {
		t.Fatalf("unexpected cursor: %s", resp.Prev())
	} else if resp.Next() != "1318388868" {
		t.Fatalf("unexpected cursor: %s", resp.Next())
	}
}

func TestAccounts_Get(t *testing.T) {
	setup()
	defer teardown()

	var customFields = &recurly.CustomFields{
		"device_id": "KIWTL-WER-ZXMRD",
	}

	mux.HandleFunc("/v2/accounts/1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
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
			  <custom_fields type="array">
			    <custom_field>
			      <name>device_id</name>
			      <value>KIWTL-WER-ZXMRD</value>
			    </custom_field>
			  </custom_fields>
			</account>`)
	})

	resp, a, err := client.Accounts.Get("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected get accounts to return OK")
	}

	ts, _ := time.Parse(recurly.DateTimeFormat, "2011-10-25T12:00:00Z")
	if diff := cmp.Diff(a, &recurly.Account{
		XMLName:   xml.Name{Local: "account"},
		Code:      "1",
		State:     "active",
		Email:     "verena@example.com",
		FirstName: "Verena",
		LastName:  "Example",
		TaxExempt: recurly.NewBool(false),
		Address: recurly.Address{
			Address: "123 Main St.",
			City:    "San Francisco",
			State:   "CA",
			Zip:     "94105",
			Country: "US",
		},
		CustomFields:     customFields,
		HostedLoginToken: "a92468579e9c4231a6c0031c4716c01d",
		CreatedAt:        recurly.NewTime(ts),
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestAccounts_Get_ErrNotFound(t *testing.T) {
	setup()
	defer teardown()

	var invoked bool
	mux.HandleFunc("/v2/accounts/1", func(w http.ResponseWriter, r *http.Request) {
		invoked = true
		w.WriteHeader(http.StatusNotFound)
	})

	_, account, err := client.Accounts.Get("1")
	if !invoked {
		t.Fatal("handler not invoked")
	} else if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if account != nil {
		t.Fatalf("expected account to be nil: %#v", account)
	}
}

func TestAccounts_LookupAccountBalance(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/1/balance", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<account_balance href="https://your-subdomain.recurly.com/v2/accounts/1/balance">
						  <account href="https://your-subdomain.recurly.com/v2/accounts/1"/>
						  <past_due type="boolean">false</past_due>
						  <balance_in_cents>
						    <USD type="integer">3000</USD>
						    <EUR type="integer">0</EUR>
						  </balance_in_cents>
						</account_balance>`)
	})

	resp, b, err := client.Accounts.LookupAccountBalance("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected get account balance to return OK")
	}

	if diff := cmp.Diff(b, &recurly.AccountBalance{
		XMLName:     xml.Name{Local: "account_balance"},
		AccountCode: "1",
		PastDue:     false,
		Balance:     3000,
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestAccounts_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(201)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><account></account>`)
	})

	resp, _, err := client.Accounts.Create(recurly.Account{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected create account to return OK")
	}
}

func TestAccounts_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/245", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><account></account>`)
	})

	resp, _, err := client.Accounts.Update("245", recurly.Account{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected update account to return OK")
	}
}

func TestAccounts_Close(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/5322", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(204)
	})

	resp, err := client.Accounts.Close("5322")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected close account to return OK")
	}
}

func TestAccounts_Reopen(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/5322/reopen", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(204)
	})

	resp, err := client.Accounts.Reopen("5322")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected reopen account to return OK")
	}
}

func TestAccounts_ListNotes(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/abcd@example.com/notes", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
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

	resp, notes, err := client.Accounts.ListNotes("abcd@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected list notes to return OK")
	} else if diff := cmp.Diff(notes, []recurly.Note{
		{
			XMLName:   xml.Name{Local: "note"},
			Message:   "This is my second note",
			CreatedAt: time.Date(2013, time.May, 14, 18, 53, 4, 0, time.UTC),
		},
		{
			XMLName:   xml.Name{Local: "note"},
			Message:   "This is my first note",
			CreatedAt: time.Date(2013, time.May, 14, 18, 52, 50, 0, time.UTC),
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}
