package recurly_test

import (
	"bytes"
	"context"
	"encoding/xml"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/blacklightcms/recurly"
	"github.com/google/go-cmp/cmp"
)

// Ensure structs are encoded to XML properly.
func TestSubscriptions_NewSubscription_Encoding(t *testing.T) {
	ts, _ := time.Parse(recurly.DateTimeFormat, "2015-06-03T13:42:23.764061Z")
	tests := []struct {
		v        recurly.NewSubscription
		expected string
	}{
		// Plan code, account, and currency are required fields. They should always be present.
		{
			expected: MustCompactString(`
				<subscription>
					<plan_code></plan_code>
					<account></account>
					<currency></currency>
				</subscription>
			`),
		},
		{
			v: recurly.NewSubscription{
				PlanCode:             "gold",
				AutoRenew:            true,
				RenewalBillingCycles: recurly.NewInt(2),
				Account: recurly.Account{
					Code: "123",
					BillingInfo: &recurly.Billing{
						Token: "507c7f79bcf86cd7994f6c0e",
					},
				},
			},
			expected: MustCompactString(`
				<subscription>
					<plan_code>gold</plan_code>
					<account>
						<account_code>123</account_code>
						<billing_info>
							<token_id>507c7f79bcf86cd7994f6c0e</token_id>
						</billing_info>
					</account>
					<currency></currency>
					<renewal_billing_cycles>2</renewal_billing_cycles>
					<auto_renew>true</auto_renew>
				</subscription>
			`),
		},
		{
			v: recurly.NewSubscription{
				PlanCode: "gold",
				Currency: "USD",
				Account: recurly.Account{
					Code: "123",
				},
				SubscriptionAddOns: &[]recurly.SubscriptionAddOn{
					{
						Code:              "extra_users",
						UnitAmountInCents: recurly.NewInt(1000),
						Quantity:          2,
					},
				},
			},
			expected: MustCompactString(`
				<subscription>
					<plan_code>gold</plan_code>
					<account>
						<account_code>123</account_code>
					</account>
					<subscription_add_ons>
						<subscription_add_on>
							<add_on_code>extra_users</add_on_code>
							<unit_amount_in_cents>1000</unit_amount_in_cents>
							<quantity>2</quantity>
						</subscription_add_on>
					</subscription_add_ons>
					<currency>USD</currency>
				</subscription>
			`),
		},
		{
			v: recurly.NewSubscription{
				PlanCode: "gold",
				Currency: "USD",
				Account: recurly.Account{
					Code: "123",
				},
				CouponCode: "promo145",
			},
			expected: MustCompactString(`
				<subscription>
					<plan_code>gold</plan_code>
					<account>
						<account_code>123</account_code>
					</account>
					<coupon_code>promo145</coupon_code>
					<currency>USD</currency>
				</subscription>
			`),
		},
		{
			v: recurly.NewSubscription{
				PlanCode: "gold",
				Currency: "USD",
				Account: recurly.Account{
					Code: "123",
				},
				UnitAmountInCents: recurly.NewInt(800),
			},
			expected: MustCompactString(`
				<subscription>
					<plan_code>gold</plan_code>
					<account>
						<account_code>123</account_code>
					</account>
					<unit_amount_in_cents>800</unit_amount_in_cents>
					<currency>USD</currency>
				</subscription>
			`),
		},
		{
			v: recurly.NewSubscription{
				PlanCode: "gold",
				Currency: "USD",
				Account: recurly.Account{
					Code: "123",
				},
				Quantity: 8,
			},
			expected: MustCompactString(`
				<subscription>
					<plan_code>gold</plan_code>
					<account>
						<account_code>123</account_code>
					</account>
					<currency>USD</currency>
					<quantity>8</quantity>
				</subscription>
			`),
		},
		{
			v: recurly.NewSubscription{
				PlanCode: "gold",
				Currency: "USD",
				Account: recurly.Account{
					Code: "123",
				},
				TrialEndsAt: recurly.NewTime(ts),
			},
			expected: MustCompactString(`
				<subscription>
					<plan_code>gold</plan_code>
					<account>
						<account_code>123</account_code>
					</account>
					<currency>USD</currency>
					<trial_ends_at>2015-06-03T13:42:23Z</trial_ends_at>
				</subscription>
			`),
		},
		{
			v: recurly.NewSubscription{
				PlanCode: "gold",
				Currency: "USD",
				Account: recurly.Account{
					Code: "123",
				},
				StartsAt: recurly.NewTime(ts),
			},
			expected: MustCompactString(`
				<subscription>
					<plan_code>gold</plan_code>
					<account>
						<account_code>123</account_code>
					</account>
					<currency>USD</currency>
					<starts_at>2015-06-03T13:42:23Z</starts_at>
				</subscription>
			`),
		},
		{
			v: recurly.NewSubscription{
				PlanCode: "gold",
				Currency: "USD",
				Account: recurly.Account{
					Code: "123",
				},
				TotalBillingCycles: 24,
			},
			expected: MustCompactString(`
				<subscription>
					<plan_code>gold</plan_code>
					<account>
						<account_code>123</account_code>
					</account>
					<currency>USD</currency>
					<total_billing_cycles>24</total_billing_cycles>
				</subscription>
			`),
		},
		{
			v: recurly.NewSubscription{
				PlanCode: "gold",
				Currency: "USD",
				Account: recurly.Account{
					Code: "123",
				},
				NextBillDate: recurly.NewTime(ts),
			},
			expected: MustCompactString(`
				<subscription>
					<plan_code>gold</plan_code>
					<account>
						<account_code>123</account_code>
					</account>
					<currency>USD</currency>
					<next_bill_date>2015-06-03T13:42:23Z</next_bill_date>
				</subscription>
			`),
		},
		{
			v: recurly.NewSubscription{
				PlanCode: "gold",
				Currency: "USD",
				Account: recurly.Account{
					Code: "123",
				},
				CollectionMethod: "automatic",
			},
			expected: MustCompactString(`
				<subscription>
					<plan_code>gold</plan_code>
					<account>
						<account_code>123</account_code>
					</account>
					<currency>USD</currency>
					<collection_method>automatic</collection_method>
				</subscription>
			`),
		},
		{
			v: recurly.NewSubscription{
				PlanCode: "gold",
				Currency: "USD",
				Account: recurly.Account{
					Code: "123",
				},
				NetTerms: recurly.NewInt(30),
			},
			expected: MustCompactString(`
				<subscription>
					<plan_code>gold</plan_code>
					<account>
						<account_code>123</account_code>
					</account>
					<currency>USD</currency>
					<net_terms>30</net_terms>
				</subscription>
			`),
		},
		{
			v: recurly.NewSubscription{
				PlanCode: "gold",
				Currency: "USD",
				Account: recurly.Account{
					Code: "123",
				},
				NetTerms: recurly.NewInt(0),
			},
			expected: MustCompactString(`
				<subscription>
					<plan_code>gold</plan_code>
					<account>
						<account_code>123</account_code>
					</account>
					<currency>USD</currency>
					<net_terms>0</net_terms>
				</subscription>
			`),
		},
		{
			v: recurly.NewSubscription{
				PlanCode: "gold",
				Currency: "USD",
				Account: recurly.Account{
					Code: "123",
				},
				PONumber: "PB4532345",
			},
			expected: MustCompactString(`
				<subscription>
					<plan_code>gold</plan_code>
					<account>
						<account_code>123</account_code>
					</account>
					<currency>USD</currency>
					<po_number>PB4532345</po_number>
				</subscription>
			`),
		},
		{
			v: recurly.NewSubscription{
				PlanCode: "gold",
				Currency: "USD",
				Account: recurly.Account{
					Code: "123",
				},
				Bulk: true,
			},
			expected: MustCompactString(`
				<subscription>
					<plan_code>gold</plan_code>
					<account>
						<account_code>123</account_code>
					</account>
					<currency>USD</currency>
					<bulk>true</bulk>
				</subscription>
			`),
		},
		{
			v: recurly.NewSubscription{
				PlanCode: "gold",
				Currency: "USD",
				Account: recurly.Account{
					Code: "123",
				},
				Bulk: false,
				// Bulk of false is the zero value of a bool, so it's omitted from the XML. But that's correct because Recurly's default is false
			},
			expected: MustCompactString(`
				<subscription>
					<plan_code>gold</plan_code>
					<account>
						<account_code>123</account_code>
					</account>
					<currency>USD</currency>
				</subscription>
			`),
		},
		{
			v: recurly.NewSubscription{
				PlanCode: "gold",
				Currency: "USD",
				Account: recurly.Account{
					Code: "123",
				},
				TermsAndConditions: "foo ... bar..",
			},
			expected: MustCompactString(`
				<subscription>
					<plan_code>gold</plan_code>
					<account>
						<account_code>123</account_code>
					</account>
					<currency>USD</currency>
					<terms_and_conditions>foo ... bar..</terms_and_conditions>
				</subscription>
			`),
		},
		{
			v: recurly.NewSubscription{
				PlanCode: "gold",
				Currency: "USD",
				Account: recurly.Account{
					Code: "123",
				},
				CustomerNotes: "foo ... customer.. bar",
			},
			expected: MustCompactString(`
				<subscription>
					<plan_code>gold</plan_code>
					<account>
						<account_code>123</account_code>
					</account>
					<currency>USD</currency>
					<customer_notes>foo ... customer.. bar</customer_notes>
				</subscription>
			`),
		},
		{
			v: recurly.NewSubscription{
				PlanCode: "gold",
				Currency: "USD",
				Account: recurly.Account{
					Code: "123",
				},
				VATReverseChargeNotes: "foo ... VAT.. bar",
			},
			expected: MustCompactString(`
				<subscription>
					<plan_code>gold</plan_code>
					<account>
						<account_code>123</account_code>
					</account>
					<currency>USD</currency>
					<vat_reverse_charge_notes>foo ... VAT.. bar</vat_reverse_charge_notes>
				</subscription>
			`),
		},
		{
			v: recurly.NewSubscription{
				PlanCode: "gold",
				Currency: "USD",
				Account: recurly.Account{
					Code: "123",
				},
				BankAccountAuthorizedAt: recurly.NewTime(ts),
			},
			expected: MustCompactString(`
				<subscription>
					<plan_code>gold</plan_code>
					<account>
						<account_code>123</account_code>
					</account>
					<currency>USD</currency>
					<bank_account_authorized_at>2015-06-03T13:42:23Z</bank_account_authorized_at>
				</subscription>
			`),
		},
		{
			v: recurly.NewSubscription{
				PlanCode: "gold",
				Currency: "USD",
				Account: recurly.Account{
					Code: "123",
				},
				BankAccountAuthorizedAt: recurly.NewTime(ts),
				CustomFields: &recurly.CustomFields{
					"platform": "2.0",
					"seller":   "Recurly Merchant",
				},
			},
			expected: MustCompactString(`
				<subscription>
					<plan_code>gold</plan_code>
					<account>
						<account_code>123</account_code>
					</account>
					<currency>USD</currency>
					<bank_account_authorized_at>2015-06-03T13:42:23Z</bank_account_authorized_at>
					<custom_fields>
						<custom_field>
							<name>platform</name>
							<value>2.0</value>
						</custom_field>
						<custom_field>
							<name>seller</name>
							<value>Recurly Merchant</value>
						</custom_field>
					</custom_fields>
				</subscription>
			`),
		},
		{
			v: recurly.NewSubscription{
				PlanCode: "gold",
				Currency: "USD",
				Account: recurly.Account{
					Code: "123",
				},
				TransactionType: "moto",
			},
			expected: MustCompactString(`
				<subscription>
					<plan_code>gold</plan_code>
					<account>
						<account_code>123</account_code>
					</account>
					<currency>USD</currency>
					<transaction_type>moto</transaction_type>
				</subscription>
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

func TestSubscriptions_UpdateSubscription_Encoding(t *testing.T) {
	tests := []struct {
		v        recurly.UpdateSubscription
		expected string
	}{
		{
			expected: MustCompactString(`
				<subscription>
				</subscription>
			`),
		},
		{
			v: recurly.UpdateSubscription{
				Timeframe: "renewal",
			},
			expected: MustCompactString(`
				<subscription>
					<timeframe>renewal</timeframe>
				</subscription>
			`),
		},
		{
			v: recurly.UpdateSubscription{
				PlanCode: "new-code",
			},
			expected: MustCompactString(`
				<subscription>
					<plan_code>new-code</plan_code>
				</subscription>
			`),
		},
		{
			v: recurly.UpdateSubscription{
				CouponCode: "coupon-code",
			},
			expected: MustCompactString(`
				<subscription>
					<coupon_code>coupon-code</coupon_code>
				</subscription>
			`),
		},
		{
			v: recurly.UpdateSubscription{
				Quantity:             14,
				AutoRenew:            recurly.NewBool(true),
				RenewalBillingCycles: recurly.NewInt(2),
			},
			expected: MustCompactString(`
				<subscription>
					<quantity>14</quantity>
					<renewal_billing_cycles>2</renewal_billing_cycles>
					<auto_renew>true</auto_renew>
				</subscription>
			`),
		},
		{
			v: recurly.UpdateSubscription{
				RemainingBillingCycles: recurly.NewInt(0),
			},
			expected: MustCompactString(`
				<subscription>
					<remaining_billing_cycles>0</remaining_billing_cycles>
				</subscription>
			`),
		},
		{
			v: recurly.UpdateSubscription{
				UnitAmountInCents: recurly.NewInt(3500),
			},
			expected: MustCompactString(`
				<subscription>
					<unit_amount_in_cents>3500</unit_amount_in_cents>
				</subscription>
			`),
		},
		{
			v: recurly.UpdateSubscription{
				CollectionMethod: "manual",
			},
			expected: MustCompactString(`
				<subscription>
					<collection_method>manual</collection_method>
				</subscription>
			`),
		},
		{
			v: recurly.UpdateSubscription{
				NetTerms: recurly.NewInt(0),
			},
			expected: MustCompactString(`
				<subscription>
					<net_terms>0</net_terms>
				</subscription>
			`),
		},
		{
			v: recurly.UpdateSubscription{
				PONumber: "AB-NewPO",
			},
			expected: MustCompactString(`
				<subscription>
					<po_number>AB-NewPO</po_number>
				</subscription>
			`),
		},
		{
			v: recurly.UpdateSubscription{
				SubscriptionAddOns: &[]recurly.SubscriptionAddOn{{
					Code:              "extra_users",
					UnitAmountInCents: recurly.NewInt(1000),
					Quantity:          2,
					AddOnSource:       "plan_add_on",
				}},
			},
			expected: MustCompactString(`
				<subscription>
					<subscription_add_ons>
						<subscription_add_on>
							<add_on_code>extra_users</add_on_code>
							<unit_amount_in_cents>1000</unit_amount_in_cents>
							<quantity>2</quantity>
							<add_on_source>plan_add_on</add_on_source>
						</subscription_add_on>
					</subscription_add_ons>
				</subscription>
			`),
		},
		{
			v: recurly.UpdateSubscription{
				TransactionType: "moto",
			},
			expected: MustCompactString(`
				<subscription>
					<transaction_type>moto</transaction_type>
				</subscription>
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

func TestSubscriptions_SubscriptionNotes_Encoding(t *testing.T) {
	tests := []struct {
		v        recurly.SubscriptionNotes
		expected string
	}{
		{
			expected: MustCompactString(`
				<subscription>
          <gateway_code></gateway_code>
				</subscription>
			`),
		},
		{
			v: recurly.SubscriptionNotes{
				GatewayCode:   "test",
				CustomerNotes: "prepaid",
			},
			expected: MustCompactString(`
				<subscription>
          <customer_notes>prepaid</customer_notes>
          <gateway_code>test</gateway_code>
				</subscription>
			`),
		},
		{
			v: recurly.SubscriptionNotes{
				TermsAndConditions: "none",
			},
			expected: MustCompactString(`
				<subscription>
          <terms_and_conditions>none</terms_and_conditions>
          <gateway_code></gateway_code>
				</subscription>
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

func TestSubscriptions_List(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	var invocations int
	s.HandleFunc("GET", "/v2/subscriptions", func(w http.ResponseWriter, r *http.Request) {
		invocations++
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("subscriptions.xml"))
	}, t)

	pager := client.Subscriptions.List(nil)
	for pager.Next() {
		var subscriptions []recurly.Subscription
		if err := pager.Fetch(context.Background(), &subscriptions); err != nil {
			t.Fatal(err)
		} else if !s.Invoked {
			t.Fatal("expected s to be invoked")
		} else if diff := cmp.Diff(subscriptions, []recurly.Subscription{*NewTestSubscription()}); diff != "" {
			t.Fatal(diff)
		}
	}
	if invocations != 1 {
		t.Fatalf("unexpected number of invocations: %d", invocations)
	}
}

func TestSubscriptions_ListAccount(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	var invocations int
	s.HandleFunc("GET", "/v2/accounts/1/subscriptions", func(w http.ResponseWriter, r *http.Request) {
		invocations++
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("subscriptions.xml"))
	}, t)

	pager := client.Subscriptions.ListAccount("1", nil)
	for pager.Next() {
		var subscriptions []recurly.Subscription
		if err := pager.Fetch(context.Background(), &subscriptions); err != nil {
			t.Fatal(err)
		} else if !s.Invoked {
			t.Fatal("expected s to be invoked")
		} else if diff := cmp.Diff(subscriptions, []recurly.Subscription{*NewTestSubscription()}); diff != "" {
			t.Fatal(diff)
		}
	}
	if invocations != 1 {
		t.Fatalf("unexpected number of invocations: %d", invocations)
	}
}

func TestSubscriptions_Get(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/subscriptions/44f83d7cba354d5b84812419f923ea96", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write(MustOpenFile("subscription.xml"))
		}, t)

		if subscription, err := client.Subscriptions.Get(context.Background(), "44f83d7c-ba35-4d5b-8481-2419f923ea96"); err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(subscription, NewTestSubscription()); diff != "" {
			t.Fatal(diff)
		} else if !s.Invoked {
			t.Fatal("expected fn invocation")
		}
	})

	// Ensure a 404 returns nil values.
	t.Run("ErrNotFound", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/subscriptions/44f83d7cba354d5b84812419f923ea96", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}, t)

		if subscription, err := client.Subscriptions.Get(context.Background(), "44f83d7c-ba35-4d5b-8481-2419f923ea96"); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if subscription != nil {
			t.Fatalf("expected nil: %#v", subscription)
		}
	})
}

func TestSubscriptions_Create(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("POST", "/v2/subscriptions", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write(MustOpenFile("subscription.xml"))
	}, t)

	if subscription, err := client.Subscriptions.Create(context.Background(), recurly.NewSubscription{}); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(subscription, NewTestSubscription()); diff != "" {
		t.Fatal(diff)
	}
}

func TestSubscriptions_Preview(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("POST", "/v2/subscriptions/preview", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write(MustOpenFile("subscription.xml"))
	}, t)

	if subscription, err := client.Subscriptions.Preview(context.Background(), recurly.NewSubscription{}); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(subscription, NewTestSubscription()); diff != "" {
		t.Fatal(diff)
	}
}

func TestSubscriptions_Update(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("PUT", "/v2/subscriptions/44f83d7cba354d5b84812419f923ea96", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("subscription.xml"))
	}, t)

	if subscription, err := client.Subscriptions.Update(context.Background(), "44f83d7c-ba35-4d5b-8481-2419f923ea96", recurly.UpdateSubscription{}); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(subscription, NewTestSubscription()); diff != "" {
		t.Fatal(diff)
	}
}

func TestSubscriptions_UpdateNotes(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("PUT", "/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/notes", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("subscription.xml"))
	}, t)

	if subscription, err := client.Subscriptions.UpdateNotes(context.Background(), "44f83d7c-ba35-4d5b-8481-2419f923ea96", recurly.SubscriptionNotes{}); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(subscription, NewTestSubscription()); diff != "" {
		t.Fatal(diff)
	}
}

func TestSubscriptions_PreviewChange(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("POST", "/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/preview", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write(MustOpenFile("subscription.xml"))
	}, t)

	if subscription, err := client.Subscriptions.PreviewChange(context.Background(), "44f83d7c-ba35-4d5b-8481-2419f923ea96", recurly.UpdateSubscription{}); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(subscription, NewTestSubscription()); diff != "" {
		t.Fatal(diff)
	}
}

func TestSubscriptions_Cancel(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("PUT", "/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/cancel", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("subscription.xml"))
	}, t)

	if subscription, err := client.Subscriptions.Cancel(context.Background(), "44f83d7c-ba35-4d5b-8481-2419f923ea96"); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(subscription, NewTestSubscription()); diff != "" {
		t.Fatal(diff)
	}
}

func TestSubscriptions_Reactivate(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("PUT", "/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/reactivate", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("subscription.xml"))
	}, t)

	if subscription, err := client.Subscriptions.Reactivate(context.Background(), "44f83d7c-ba35-4d5b-8481-2419f923ea96"); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(subscription, NewTestSubscription()); diff != "" {
		t.Fatal(diff)
	}
}

func TestSubscriptions_Terminate(t *testing.T) {
	t.Run("Partial", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("PUT", "/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/terminate", func(w http.ResponseWriter, r *http.Request) {
			if v := r.URL.Query().Get("refund"); v != "partial" {
				t.Fatalf("unexpected refund type: %q", v)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(MustOpenFile("subscription.xml"))
		}, t)

		if subscription, err := client.Subscriptions.Terminate(context.Background(), "44f83d7c-ba35-4d5b-8481-2419f923ea96", "partial"); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(subscription, NewTestSubscription()); diff != "" {
			t.Fatal(diff)
		}
	})

	t.Run("Full", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("PUT", "/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/terminate", func(w http.ResponseWriter, r *http.Request) {
			if v := r.URL.Query().Get("refund"); v != "full" {
				t.Fatalf("unexpected refund type: %q", v)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(MustOpenFile("subscription.xml"))
		}, t)

		if subscription, err := client.Subscriptions.Terminate(context.Background(), "44f83d7c-ba35-4d5b-8481-2419f923ea96", "full"); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(subscription, NewTestSubscription()); diff != "" {
			t.Fatal(diff)
		}
	})

	t.Run("None", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("PUT", "/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/terminate", func(w http.ResponseWriter, r *http.Request) {
			if v := r.URL.Query().Get("refund"); v != "none" {
				t.Fatalf("unexpected refund type: %q", v)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(MustOpenFile("subscription.xml"))
		}, t)

		if subscription, err := client.Subscriptions.Terminate(context.Background(), "44f83d7c-ba35-4d5b-8481-2419f923ea96", "none"); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(subscription, NewTestSubscription()); diff != "" {
			t.Fatal(diff)
		}
	})
}

func TestSubscriptions_Pause(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("PUT", "/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/pause", func(w http.ResponseWriter, r *http.Request) {
		if str := MustReadAllString(r.Body); str != MustCompactString(`
			<subscription>
				<remaining_pause_cycles>1</remaining_pause_cycles>
			</subscription>
		`) {
			t.Fatal(str)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("subscription.xml"))
	}, t)

	if subscription, err := client.Subscriptions.Pause(context.Background(), "44f83d7c-ba35-4d5b-8481-2419f923ea96", 1); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(subscription, NewTestSubscription()); diff != "" {
		t.Fatal(diff)
	}
}

func TestSubscriptions_Postpone(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("PUT", "/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/postpone", func(w http.ResponseWriter, r *http.Request) {
		if v := r.URL.Query().Get("next_renewal_date"); v != "2015-08-27T07:00:00Z" {
			t.Fatalf("unexpected input for next_renewal_date: %q", v)
		} else if v := r.URL.Query().Get("bulk"); v != "false" {
			t.Fatalf("unexpected input for bulk: %q", v)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("subscription.xml"))
	}, t)

	if subscription, err := client.Subscriptions.Postpone(context.Background(), "44f83d7c-ba35-4d5b-8481-2419f923ea96", MustParseTime("2015-08-27T07:00:00Z"), false); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(subscription, NewTestSubscription()); diff != "" {
		t.Fatal(diff)
	}
}

func TestSubscriptions_Resume(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("PUT", "/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/resume", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("subscription.xml"))
	}, t)

	if subscription, err := client.Subscriptions.Resume(context.Background(), "44f83d7c-ba35-4d5b-8481-2419f923ea96"); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(subscription, NewTestSubscription()); diff != "" {
		t.Fatal(diff)
	}
}

func TestSubscriptions_ConvertTrial(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("PUT", "/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/convert_trial", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("subscription.xml"))
	}, t)

	if subscription, err := client.Subscriptions.ConvertTrial(context.Background(), "44f83d7c-ba35-4d5b-8481-2419f923ea96"); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(subscription, NewTestSubscription()); diff != "" {
		t.Fatal(diff)
	}
}

// Returns a Subscription corresponding to testdata/subscription.xml.
func NewTestSubscription() *recurly.Subscription {
	return &recurly.Subscription{
		XMLName: xml.Name{Local: "subscription"},
		Plan: recurly.NestedPlan{
			Code: "gold",
			Name: "Gold plan",
		},
		AccountCode:            "1",
		InvoiceNumber:          1108,
		UUID:                   "44f83d7cba354d5b84812419f923ea96", // UUID has been sanitized
		State:                  "active",
		UnitAmountInCents:      800,
		Currency:               "EUR",
		Quantity:               1,
		ActivatedAt:            recurly.NewTime(time.Date(2011, time.May, 27, 7, 0, 0, 0, time.UTC)),
		CurrentPeriodStartedAt: recurly.NewTime(time.Date(2011, time.June, 27, 7, 0, 0, 0, time.UTC)),
		CurrentPeriodEndsAt:    recurly.NewTime(time.Date(2011, time.July, 27, 7, 0, 0, 0, time.UTC)),
		TaxInCents:             72,
		TaxType:                "usst",
		TaxRegion:              "CA",
		TaxRate:                0.0875,
		NetTerms:               recurly.NewInt(0),
		CustomerNotes:          "customer_notes_test_get",
		SubscriptionAddOns: []recurly.SubscriptionAddOn{
			{
				XMLName:           xml.Name{Local: "subscription_add_on"},
				Type:              "fixed",
				Code:              "add-on-one",
				UnitAmountInCents: recurly.NewInt(1000),
				Quantity:          2,
			},
		},
		PendingSubscription: &recurly.PendingSubscription{
			XMLName: xml.Name{Local: "pending_subscription"},
			Plan: recurly.NestedPlan{
				Code: "gold",
				Name: "Gold plan",
			},
			Quantity:          1,
			UnitAmountInCents: 50000,
			SubscriptionAddOns: []recurly.SubscriptionAddOn{
				{
					XMLName:           xml.Name{Local: "subscription_add_on"},
					Type:              "fixed",
					Code:              "add-on-one",
					UnitAmountInCents: recurly.NewInt(1100),
					Quantity:          1,
				},
				{
					XMLName:  xml.Name{Local: "subscription_add_on"},
					Type:     "fixed",
					Code:     "add-on-two",
					Quantity: 1,
				},
			},
		},
		CustomFields: &recurly.CustomFields{
			"device_id":     "KIWTL-WER-ZXMRD",
			"purchase_date": "2017-01-23",
		},
		InvoiceCollection: &recurly.InvoiceCollection{
			XMLName: xml.Name{Local: "invoice_collection"},
			ChargeInvoice: &recurly.Invoice{
				XMLName:     xml.Name{Local: "invoice"},
				AccountCode: "1",
				UUID:        "43adfe52c21cbb221557a24940bcd7e5",
				State:       recurly.ChargeInvoiceStatePending,
			},
			CreditInvoices: []recurly.Invoice{
				{
					XMLName:         xml.Name{Local: "invoice"},
					AccountCode:     "1",
					UUID:            "43adfe52c21cbb221557a24940bcd7e5",
					State:           recurly.CreditInvoiceStateOpen,
					TotalInCents:    -4014,
					Currency:        "USD",
					SubtotalInCents: -4014,
					DiscountInCents: -436,
					BalanceInCents:  -2444,
					Type:            "credit",
					Origin:          "immediate_change",
				},
			},
		},
	}
}
