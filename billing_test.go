package recurly_test

import (
	"bytes"
	"context"
	"encoding/xml"
	"net"
	"net/http"
	"strconv"
	"testing"

	"github.com/blacklightcms/recurly"
	"github.com/google/go-cmp/cmp"
)

// Ensure structs are encoded to XML properly.
func TestBilling_Encoding(t *testing.T) {
	tests := []struct {
		v        recurly.Billing
		expected string
	}{
		{
			expected: MustCompactString(`
				<billing_info>
				</billing_info>
		`),
		},
		{
			v: recurly.Billing{Token: "507c7f79bcf86cd7994f6c0e"},
			expected: MustCompactString(`
				<billing_info>
					<token_id>507c7f79bcf86cd7994f6c0e</token_id>
				</billing_info>
			`),
		},
		{
			v: recurly.Billing{FirstName: "Verena", LastName: "Example"},
			expected: MustCompactString(`
				<billing_info>
					<first_name>Verena</first_name>
					<last_name>Example</last_name>
				</billing_info>
			`),
		},
		{
			v: recurly.Billing{Address: "123 Main St."},
			expected: MustCompactString(`
				<billing_info>
					<address1>123 Main St.</address1>
				</billing_info>
			`),
		},
		{
			v: recurly.Billing{Address2: "Unit A"},
			expected: MustCompactString(`
				<billing_info>
					<address2>Unit A</address2>
				</billing_info>
			`),
		},
		{
			v: recurly.Billing{City: "San Francisco"},
			expected: MustCompactString(`
				<billing_info>
					<city>San Francisco</city>
				</billing_info>
			`),
		},
		{
			v: recurly.Billing{State: "CA"},
			expected: MustCompactString(`
				<billing_info>
					<state>CA</state>
				</billing_info>
			`),
		},
		{
			v: recurly.Billing{Zip: "94105"},
			expected: MustCompactString(`
				<billing_info>
					<zip>94105</zip>
				</billing_info>
			`),
		},
		{
			v: recurly.Billing{Country: "US"},
			expected: MustCompactString(`
				<billing_info>
					<country>US</country>
				</billing_info>
			`),
		},
		{
			v: recurly.Billing{Phone: "555-555-5555"},
			expected: MustCompactString(`
				<billing_info>
					<phone>555-555-5555</phone>
				</billing_info>
			`),
		},
		{
			v: recurly.Billing{VATNumber: "abc"},
			expected: MustCompactString(`
				<billing_info>
					<vat_number>abc</vat_number>
				</billing_info>
			`),
		},
		{
			v: recurly.Billing{IPAddress: net.ParseIP("127.0.0.1")},
			expected: MustCompactString(`
				<billing_info>
					<ip_address>127.0.0.1</ip_address>
				</billing_info>
			`),
		},
		{
			v: recurly.Billing{Number: 4111111111111111, Month: 5, Year: 2020, VerificationValue: 111},
			expected: MustCompactString(`
				<billing_info>
					<number>4111111111111111</number>
					<month>5</month>
					<year>2020</year>
					<verification_value>111</verification_value>
				</billing_info>
			`),
		},
		{
			v: recurly.Billing{RoutingNumber: "065400137", AccountNumber: "0123456789", AccountType: "checking"},
			expected: MustCompactString(`
				<billing_info>
					<routing_number>065400137</routing_number>
					<account_number>0123456789</account_number>
					<account_type>checking</account_type>
				</billing_info>
			`),
		},
		{
			v: recurly.Billing{TransactionType: "moto"},
			expected: MustCompactString(`
				<billing_info>
					<transaction_type>moto</transaction_type>
				</billing_info>
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

func TestBilling_Type(t *testing.T) {
	t.Run("Card", func(t *testing.T) {
		if typ := (recurly.Billing{
			FirstSix: "011111",
			LastFour: "1111",
			Month:    11,
			Year:     2020,
		}).Type(); typ != "card" {
			t.Fatalf("unexpected type: %s", typ)
		}
	})

	t.Run("Bank", func(t *testing.T) {
		if typ := (recurly.Billing{
			NameOnAccount: "Acme, Inc",
			RoutingNumber: "123456780",
			AccountNumber: "111111111",
		}).Type(); typ != "bank" {
			t.Fatalf("unexpected type: %s", typ)
		}
	})

	t.Run("None", func(t *testing.T) {
		if typ := (recurly.Billing{}).Type(); typ != "" {
			t.Fatalf("unexpected type: %s", typ)
		}
	})
}

func TestBilling_Get(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/accounts/1/billing_info", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write(MustOpenFile("billing_info.xml"))
		}, t)

		if info, err := client.Billing.Get(context.Background(), "1"); err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(info, NewTestBillingInfo()); diff != "" {
			t.Fatal(diff)
		} else if !s.Invoked {
			t.Fatal("expected fn invocation")
		}
	})

	// ACH customers may not have billing info. This asserts that nil values for
	// many of the fields are safely ignored without parse errors.
	t.Run("ACH", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/accounts/1/billing_info", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`
				<?xml version="1.0" encoding="UTF-8"?>
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
				</billing_info>
				`,
			))
		}, t)

		if info, err := client.Billing.Get(context.Background(), "1"); err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(info, &recurly.Billing{
			XMLName:     xml.Name{Local: "billing_info"},
			FirstName:   "Verena",
			LastName:    "Example",
			PaymentType: "ach",
		}); diff != "" {
			t.Fatal(diff)
		} else if !s.Invoked {
			t.Fatal("expected fn invocation")
		}
	})

	// Ensure a 404 returns nil values.
	t.Run("ErrNotFound", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/accounts/1/billing_info", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}, t)

		if info, err := client.Billing.Get(context.Background(), "1"); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if info != nil {
			t.Fatalf("expected nil: %#v", info)
		}
	})
}

func TestBilling_Create(t *testing.T) {
	t.Run("Token", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("POST", "/v2/accounts/1/billing_info", func(w http.ResponseWriter, r *http.Request) {
			if str := MustReadAllString(r.Body); str != MustCompactString(`
			<billing_info>
				<token_id>TOKEN</token_id>
			</billing_info>
		`) {
				t.Fatal(str)
			}
			w.WriteHeader(http.StatusCreated)
			w.Write(MustOpenFile("billing_info.xml"))
		}, t)

		if info, err := client.Billing.Create(context.Background(), "1", recurly.Billing{Token: "TOKEN"}); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(info, NewTestBillingInfo()); diff != "" {
			t.Fatal(diff)
		}
	})

	t.Run("Token with 3D secure action result token", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("POST", "/v2/accounts/1/billing_info", func(w http.ResponseWriter, r *http.Request) {
			if str := MustReadAllString(r.Body); str != MustCompactString(`
			<billing_info>
				<token_id>TOKEN</token_id>
				<three_d_secure_action_result_token_id>THREE_D_SECURE_ACTION_RESULT_TOKEN</three_d_secure_action_result_token_id>
			</billing_info>
		`) {
				t.Fatal(str)
			}
			w.WriteHeader(http.StatusCreated)
			w.Write(MustOpenFile("billing_info.xml"))
		}, t)

		if info, err := client.Billing.Create(context.Background(), "1", recurly.Billing{Token: "TOKEN", ThreeDSecureActionResultTokenID: "THREE_D_SECURE_ACTION_RESULT_TOKEN"}); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(info, NewTestBillingInfo()); diff != "" {
			t.Fatal(diff)
		}
	})

	t.Run("BillingInfo", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("POST", "/v2/accounts/1/billing_info", func(w http.ResponseWriter, r *http.Request) {
			if str := MustReadAllString(r.Body); str != MustCompactString(`
			<billing_info>
				<first_name>Verena</first_name>
				<last_name>Example</last_name>
				<address1>123 Main St.</address1>
				<city>San Francisco</city>
				<state>CA</state>
				<zip>94105</zip>
				<country>US</country>
				<number>4111111111111111</number>
				<month>10</month>
				<year>2020</year>
			</billing_info>
		`) {
				t.Fatal(str)
			}
			w.WriteHeader(http.StatusCreated)
			w.Write(MustOpenFile("billing_info.xml"))
		}, t)

		if info, err := client.Billing.Create(context.Background(), "1", recurly.Billing{
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
		}); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(info, NewTestBillingInfo()); diff != "" {
			t.Fatal(diff)
		}
	})
}

func TestBilling_Update(t *testing.T) {
	t.Run("Token", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("PUT", "/v2/accounts/1/billing_info", func(w http.ResponseWriter, r *http.Request) {
			if str := MustReadAllString(r.Body); str != MustCompactString(`
				<billing_info>
					<token_id>TOKEN</token_id>
				</billing_info>
			`) {
				t.Fatal(str)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(MustOpenFile("billing_info.xml"))
		}, t)

		if info, err := client.Billing.Update(context.Background(), "1", recurly.Billing{Token: "TOKEN"}); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(info, NewTestBillingInfo()); diff != "" {
			t.Fatal(diff)
		}
	})

	t.Run("Token with 3D secure action result token", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("PUT", "/v2/accounts/1/billing_info", func(w http.ResponseWriter, r *http.Request) {
			if str := MustReadAllString(r.Body); str != MustCompactString(`
				<billing_info>
					<token_id>TOKEN</token_id>
					<three_d_secure_action_result_token_id>THREE_D_SECURE_ACTION_RESULT_TOKEN</three_d_secure_action_result_token_id>
				</billing_info>
			`) {
				t.Fatal(str)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(MustOpenFile("billing_info.xml"))
		}, t)

		if info, err := client.Billing.Update(context.Background(), "1", recurly.Billing{Token: "TOKEN", ThreeDSecureActionResultTokenID: "THREE_D_SECURE_ACTION_RESULT_TOKEN"}); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(info, NewTestBillingInfo()); diff != "" {
			t.Fatal(diff)
		}
	})

	t.Run("InvalidToken", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("PUT", "/v2/accounts/1/billing_info", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte(`
				<?xml version="1.0" encoding="UTF-8"?>
				<error>
					<symbol>token_invalid</symbol>
					<description>Token is either invalid or expired</description>
				</error>
				`,
			))
		}, t)

		if info, err := client.Billing.Update(context.Background(), "1", recurly.Billing{Token: "TOKEN"}); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if e, ok := err.(*recurly.ClientError); !ok {
			t.Fatalf("expected *recurly.ClientError, got %T: %#v", err, err)
		} else if diff := cmp.Diff(e.ValidationErrors, []recurly.ValidationError{{
			Symbol:      "token_invalid",
			Description: "Token is either invalid or expired",
		}}); diff != "" {
			t.Fatal(diff)
		} else if info != nil {
			t.Fatalf("expected info to be nil: %#v", info)
		}
	})

	t.Run("BillingInfo", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("PUT", "/v2/accounts/1/billing_info", func(w http.ResponseWriter, r *http.Request) {
			if str := MustReadAllString(r.Body); str != MustCompactString(`
			<billing_info>
				<first_name>Verena</first_name>
				<last_name>Example</last_name>
				<address1>123 Main St.</address1>
				<city>San Francisco</city>
				<state>CA</state>
				<zip>94105</zip>
				<country>US</country>
				<number>4111111111111111</number>
				<month>10</month>
				<year>2020</year>
			</billing_info>
		`) {
				t.Fatal(str)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(MustOpenFile("billing_info.xml"))
		}, t)

		if info, err := client.Billing.Update(context.Background(), "1", recurly.Billing{
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
		}); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(info, NewTestBillingInfo()); diff != "" {
			t.Fatal(diff)
		}
	})
}

func TestBilling_Clear(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("DELETE", "/v2/accounts/1/billing_info", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}, t)

	if err := client.Billing.Clear(context.Background(), "1"); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	}
}

// Returns Billing corresponding to testdata/billing_info.xml.
func NewTestBillingInfo() *recurly.Billing {
	return &recurly.Billing{
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
		FirstSix:         "411111",
		LastFour:         "1111",
		PaymentType:      "credit_card",
	}
}
