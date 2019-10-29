package recurly_test

import (
	"bytes"
	"context"
	"encoding/xml"
	"net/http"
	"strconv"
	"testing"

	"github.com/blacklightcms/recurly"
	"github.com/google/go-cmp/cmp"
)

// Ensure structs are encoded to XML properly.
func TestAccounts_Encoding(t *testing.T) {
	tests := []struct {
		v        recurly.Account
		expected string
	}{
		{
			expected: MustCompactString(`
				<account></account>
			`),
		},
		{
			v: recurly.Account{Code: "abc"},
			expected: MustCompactString(`
				<account>
					<account_code>abc</account_code>
				</account>
			`),
		},
		{
			v: recurly.Account{State: "active"},
			expected: MustCompactString(`
				<account>
					<state>active</state>
				</account>
			`),
		},
		{
			v: recurly.Account{Email: "me@example.com"},
			expected: MustCompactString(`
				<account>
					<email>me@example.com</email>
				</account>
			`),
		},
		{
			v: recurly.Account{FirstName: "Larry"},
			expected: MustCompactString(`
				<account>
					<first_name>Larry</first_name>
				</account>
			`),
		},
		{
			v: recurly.Account{LastName: "Larrison"},
			expected: MustCompactString(`
				<account>
					<last_name>Larrison</last_name>
				</account>
			`),
		},
		{
			v: recurly.Account{FirstName: "Larry", LastName: "Larrison"},
			expected: MustCompactString(`
				<account>
					<first_name>Larry</first_name>
					<last_name>Larrison</last_name>
				</account>
			`),
		},
		{
			v: recurly.Account{CompanyName: "Acme, Inc"},
			expected: MustCompactString(`
				<account>
					<company_name>Acme, Inc</company_name>
				</account>
			`),
		},
		{
			v: recurly.Account{VATNumber: "123456789"},
			expected: MustCompactString(`
				<account>
					<vat_number>123456789</vat_number>
				</account>
			`),
		},
		{
			v: recurly.Account{TaxExempt: recurly.NewBool(true)},
			expected: MustCompactString(`
				<account>
					<tax_exempt>true</tax_exempt>
				</account>
			`),
		},
		{
			v: recurly.Account{TaxExempt: recurly.NewBool(false)},
			expected: MustCompactString(`
				<account>
					<tax_exempt>false</tax_exempt>
				</account>
			`),
		},
		{
			v: recurly.Account{AcceptLanguage: "en-US"},
			expected: MustCompactString(`
				<account>
					<accept_language>en-US</accept_language>
				</account>
			`),
		},
		{
			v: recurly.Account{PreferredLocale: "en-US"},
			expected: MustCompactString(`
				<account>
					<preferred_locale>en-US</preferred_locale>
				</account>
			`),
		},
		{
			v: recurly.Account{TransactionType: "moto"},
			expected: MustCompactString(`
				<account>
					<transaction_type>moto</transaction_type>
				</account>
			`),
		},
		{
			v: recurly.Account{FirstName: "Larry", Address: &recurly.Address{Address: "123 Main St.", City: "San Francisco", State: "CA", Zip: "94105", Country: "US"}},
			expected: MustCompactString(`
				<account>
					<first_name>Larry</first_name>
					<address>
						<address1>123 Main St.</address1>
						<city>San Francisco</city>
						<state>CA</state>
						<zip>94105</zip>
						<country>US</country>
					</address>
				</account>
			`),
		},
		{
			v: recurly.Account{Code: "test@example.com", BillingInfo: &recurly.Billing{Token: "507c7f79bcf86cd7994f6c0e"}},
			expected: MustCompactString(`
				<account>
					<account_code>test@example.com</account_code>
					<billing_info>
						<token_id>507c7f79bcf86cd7994f6c0e</token_id>
					</billing_info>
				</account>
			`),
		},
		{
			v: recurly.Account{HasPausedSubscription: recurly.NewBool(true)},
			expected: MustCompactString(`
				<account>
					<has_paused_subscription>true</has_paused_subscription>
				</account>
			`),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			buf := new(bytes.Buffer)
			if err := xml.NewEncoder(buf).Encode(tt.v); err != nil {
				t.Fatal(err)
			} else if buf.String() != tt.expected {
				t.Fatal(buf.String())
			}
		})
	}
}

// Ensure structs are encoded to XML properly.
func TestAddress_Encoding(t *testing.T) {
	tests := []struct {
		v        recurly.Address
		expected string
	}{
		{
			expected: MustCompactString(`
				<address>
				</address>
			`),
		},
		{
			v: recurly.Address{Address: "123 Main St."},
			expected: MustCompactString(`
				<address>
					<address1>123 Main St.</address1>
				</address>
			`),
		},
		{
			v: recurly.Address{Address2: "Unit A"},
			expected: MustCompactString(`
				<address>
					<address2>Unit A</address2>
				</address>
			`),
		},
		{
			v: recurly.Address{City: "San Francisco"},
			expected: MustCompactString(`
				<address>
					<city>San Francisco</city>
				</address>
			`),
		},
		{
			v: recurly.Address{State: "CA"},
			expected: MustCompactString(`
				<address>
					<state>CA</state>
				</address>
			`),
		},
		{
			v: recurly.Address{Zip: "94105"},
			expected: MustCompactString(`
				<address>
					<zip>94105</zip>
				</address>
			`),
		},
		{
			v: recurly.Address{Country: "US"},
			expected: MustCompactString(`
				<address>
					<country>US</country>
				</address>
			`),
		},
		{
			v: recurly.Address{Phone: "555-555-5555"},
			expected: MustCompactString(`
				<address>
					<phone>555-555-5555</phone>
				</address>
			`),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			buf := new(bytes.Buffer)
			if err := xml.NewEncoder(buf).Encode(tt.v); err != nil {
				t.Fatal(err)
			} else if buf.String() != tt.expected {
				t.Fatal(buf.String())
			}
		})
	}
}

func TestAccounts_List(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	var invocations int
	s.HandleFunc("GET", "/v2/accounts", func(w http.ResponseWriter, r *http.Request) {
		invocations++
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("accounts.xml"))
	}, t)

	pager := client.Accounts.List(nil)
	for pager.Next() {
		var a []recurly.Account
		if err := pager.Fetch(context.Background(), &a); err != nil {
			t.Fatal(err)
		} else if !s.Invoked {
			t.Fatal("expected s to be invoked")
		} else if diff := cmp.Diff(a, []recurly.Account{*NewTestAccount()}); diff != "" {
			t.Fatal(diff)
		}
	}
	if invocations != 1 {
		t.Fatalf("unexpected number of invocations: %d", invocations)
	}
}

func TestAccounts_Get(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/accounts/1", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write(MustOpenFile("account.xml"))
		}, t)

		if a, err := client.Accounts.Get(context.Background(), "1"); err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(a, NewTestAccount()); diff != "" {
			t.Fatal(diff)
		} else if !s.Invoked {
			t.Fatal("expected fn invocation")
		}
	})

	// Ensure a 404 returns nil values.
	t.Run("ErrNotFound", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/accounts/1", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}, t)

		if a, err := client.Accounts.Get(context.Background(), "1"); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if a != nil {
			t.Fatalf("expected nil account: %#v", a)
		}
	})
}

func TestAccounts_Balance(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("GET", "/v2/accounts/1/balance", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("account_balance.xml"))
	}, t)

	if balance, err := client.Accounts.Balance(context.Background(), "1"); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(balance, NewTestAccountBalance()); diff != "" {
		t.Fatal(diff)
	}
}

func TestAccounts_Create(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("POST", "/v2/accounts", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write(MustOpenFile("account.xml"))
	}, t)

	if a, err := client.Accounts.Create(context.Background(), recurly.Account{}); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(a, NewTestAccount()); diff != "" {
		t.Fatal(diff)
	}
}

func TestAccounts_Update(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("PUT", "/v2/accounts/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("account.xml"))
	}, t)

	if a, err := client.Accounts.Update(context.Background(), "1", recurly.Account{}); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(a, NewTestAccount()); diff != "" {
		t.Fatal(diff)
	}
}

func TestAccounts_Close(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("DELETE", "/v2/accounts/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}, t)

	if err := client.Accounts.Close(context.Background(), "1"); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	}
}

func TestAccounts_Reopen(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("PUT", "/v2/accounts/1/reopen", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}, t)

	if err := client.Accounts.Reopen(context.Background(), "1"); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	}
}

func TestAccounts_ListNotes(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	var invocations int
	s.HandleFunc("GET", "/v2/accounts/1/notes", func(w http.ResponseWriter, r *http.Request) {
		invocations++
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("notes.xml"))
	}, t)

	pager := client.Accounts.ListNotes("1", nil)
	for pager.Next() {
		var n []recurly.Note
		if err := pager.Fetch(context.Background(), &n); err != nil {
			t.Fatal(err)
		} else if !s.Invoked {
			t.Fatal("expected s to be invoked")
		} else if diff := cmp.Diff(n, NewTestNotes()); diff != "" {
			t.Fatal(diff)
		}
	}
	if invocations != 1 {
		t.Fatalf("unexpected number of invocations: %d", invocations)
	}
}

// Returns an account corresponding to testdata/account.xml
func NewTestAccount() *recurly.Account {
	ts := MustParseTime("2011-10-25T12:00:00Z")
	return &recurly.Account{
		XMLName:     xml.Name{Local: "account"},
		Code:        "1",
		State:       "active",
		Email:       "verena@example.com",
		FirstName:   "Verena",
		LastName:    "Example",
		TaxExempt:   recurly.NewBool(false),
		BillingInfo: NewTestBillingInfo(),
		Address: &recurly.Address{
			XMLName: xml.Name{Local: "address"},
			Address: "123 Main St.",
			City:    "San Francisco",
			State:   "CA",
			Zip:     "94105",
			Country: "US",
		},
		ShippingAddresses: &[]recurly.ShippingAddress{
			*NewTestShippingAddress(),
		},
		CustomFields: &recurly.CustomFields{
			"device_id": "KIWTL-WER-ZXMRD",
		},
		HostedLoginToken:        "a92468579e9c4231a6c0031c4716c01d",
		CreatedAt:               recurly.NewTime(ts),
		HasLiveSubscription:     recurly.NewBool(true),
		HasActiveSubscription:   recurly.NewBool(true),
		HasFutureSubscription:   recurly.NewBool(false),
		HasCanceledSubscription: recurly.NewBool(false),
		HasPastDueInvoice:       recurly.NewBool(false),
	}
}

// Returns an account balance corresponding to testdata/account_balance.xml
func NewTestAccountBalance() *recurly.AccountBalance {
	return &recurly.AccountBalance{
		XMLName: xml.Name{Local: "account_balance"},
		PastDue: false,
		Balance: recurly.UnitAmount{
			USD: 3000,
		},
	}
}

// Returns account notes corresponding to testdata/notes.xml
func NewTestNotes() []recurly.Note {
	return []recurly.Note{
		{
			XMLName:   xml.Name{Local: "note"},
			Message:   "This is my second note",
			CreatedAt: MustParseTime("2013-05-14T18:53:04Z"),
		},
		{
			XMLName:   xml.Name{Local: "note"},
			Message:   "This is my first note",
			CreatedAt: MustParseTime("2013-05-14T18:52:50Z"),
		},
	}
}
