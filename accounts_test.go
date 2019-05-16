package recurly_test

import (
	"bytes"
	"context"
	"encoding/xml"
	"net/http"
	"net/url"
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

// Test paginating accounts. Note that all pagination code uses the same underlying
// pager -- this test is the only paginated test that asserts the params being sent
// properly. It also sends a cursor to ensure Next() and Fetch() work properly.
func TestAccounts_List(t *testing.T) {
	client, s := NewServer()
	defer s.Close()

	var invocations int
	s.HandleFunc("GET", "/v2/accounts", func(w http.ResponseWriter, r *http.Request) {
		cursor := r.URL.Query().Get("cursor")
		switch invocations {
		case 0:
			if cursor != "" {
				t.Fatalf("unexpected cursor: %s", cursor)
			}
			w.Header().Set("Link", `<https://test.recurly.com/v2/accounts?cursor=1972702718353176814:A1465932489>; rel="next"`)
		case 1:
			if cursor != "1972702718353176814:A1465932489" {
				t.Fatalf("unexpected cursor: %s", cursor)
			}
		default:
			t.Fatalf("unexpected number of invocations")
		}

		query := r.URL.Query()
		query.Del("cursor") // conditionally checked above
		if diff := cmp.Diff(query, url.Values{
			"per_page":   []string{"50"},
			"sort":       []string{"created_at"},
			"order":      []string{"asc"},
			"state":      []string{"active"},
			"begin_time": []string{"2011-10-17T17:24:53Z"},
			"end_time":   []string{"2011-10-18T17:24:53Z"},
		}); diff != "" {
			t.Fatal(diff)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("accounts.xml"))
		invocations++
	}, t)

	pager := client.Accounts.List(&recurly.PagerOptions{
		PerPage:   50,
		Sort:      "created_at",
		Order:     "asc",
		State:     "active",
		BeginTime: recurly.NewTime(MustParseTime("2011-10-17T17:24:53Z")),
		EndTime:   recurly.NewTime(MustParseTime("2011-10-18T17:24:53Z")),
	})

	for pager.Next() {
		if a, err := pager.Fetch(context.Background()); err != nil {
			t.Fatal(err)
		} else if !s.Invoked {
			t.Fatal("expected s to be invoked")
		} else if diff := cmp.Diff(a, []recurly.Account{*NewTestAccount()}); diff != "" {
			t.Fatal(diff)
		}
		s.Invoked = false
	}
	if invocations != 2 {
		t.Fatalf("unexpected number of invocations: %d", invocations)
	}
}

func TestAccounts_Get(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		client, s := NewServer()
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
		client, s := NewServer()
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
	client, s := NewServer()
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
	client, s := NewServer()
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
	client, s := NewServer()
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
	client, s := NewServer()
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
	client, s := NewServer()
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
	client, s := NewServer()
	defer s.Close()

	var invocations int
	s.HandleFunc("GET", "/v2/accounts/1/notes", func(w http.ResponseWriter, r *http.Request) {
		invocations++
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("notes.xml"))
	}, t)

	pager := client.Accounts.Notes("1", nil)
	for pager.Next() {
		n, err := pager.Fetch(context.Background())
		if err != nil {
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
