package recurly_test

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/blacklightcms/recurly"
)

// TestSubscriptionsEncoding ensures structs are encoded to XML properly.
// Because Recurly supports partial updates, it's important that only defined
// fields are handled properly -- including types like booleans and integers which
// have zero values that we want to send.
func TestSubscriptions_NewSubscription_Encoding(t *testing.T) {
	ts, _ := time.Parse(recurly.DateTimeFormat, "2015-06-03T13:42:23.764061Z")
	tests := []struct {
		v        recurly.NewSubscription
		expected string
	}{
		// Plan code, account, and currency are required fields. They should always be present.
		{
			expected: "<subscription><plan_code></plan_code><account></account><currency></currency></subscription>",
		},
		{
			v: recurly.NewSubscription{
				PlanCode: "gold",
				Account: recurly.Account{
					Code: "123",
					BillingInfo: &recurly.Billing{
						Token: "507c7f79bcf86cd7994f6c0e",
					},
				},
			},
			expected: "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code><billing_info><token_id>507c7f79bcf86cd7994f6c0e</token_id></billing_info></account><currency></currency></subscription>",
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
						UnitAmountInCents: 1000,
						Quantity:          2,
					},
				},
			},
			expected: "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><subscription_add_ons><subscription_add_on><add_on_code>extra_users</add_on_code><unit_amount_in_cents>1000</unit_amount_in_cents><quantity>2</quantity></subscription_add_on></subscription_add_ons><currency>USD</currency></subscription>",
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
			expected: "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><coupon_code>promo145</coupon_code><currency>USD</currency></subscription>",
		},
		{
			v: recurly.NewSubscription{
				PlanCode: "gold",
				Currency: "USD",
				Account: recurly.Account{
					Code: "123",
				},
				UnitAmountInCents: 800,
			},
			expected: "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><unit_amount_in_cents>800</unit_amount_in_cents><currency>USD</currency></subscription>",
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
			expected: "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency><quantity>8</quantity></subscription>",
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
			expected: "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency><trial_ends_at>2015-06-03T13:42:23Z</trial_ends_at></subscription>",
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
			expected: "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency><starts_at>2015-06-03T13:42:23Z</starts_at></subscription>",
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
			expected: "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency><total_billing_cycles>24</total_billing_cycles></subscription>",
		},
		{
			v: recurly.NewSubscription{
				PlanCode: "gold",
				Currency: "USD",
				Account: recurly.Account{
					Code: "123",
				},
				FirstRenewalDate: recurly.NewTime(ts),
			},
			expected: "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency><first_renewal_date>2015-06-03T13:42:23Z</first_renewal_date></subscription>",
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
			expected: "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency><collection_method>automatic</collection_method></subscription>",
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
			expected: "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency><net_terms>30</net_terms></subscription>",
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
			expected: "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency><net_terms>0</net_terms></subscription>",
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
			expected: "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency><po_number>PB4532345</po_number></subscription>",
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
			expected: "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency><bulk>true</bulk></subscription>",
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
			expected: "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency></subscription>",
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
			expected: "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency><terms_and_conditions>foo ... bar..</terms_and_conditions></subscription>",
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
			expected: "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency><customer_notes>foo ... customer.. bar</customer_notes></subscription>",
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
			expected: "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency><vat_reverse_charge_notes>foo ... VAT.. bar</vat_reverse_charge_notes></subscription>",
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
			expected: "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency><bank_account_authorized_at>2015-06-03T13:42:23Z</bank_account_authorized_at></subscription>",
		},
	}

	for i, tt := range tests {
		var given bytes.Buffer
		if err := xml.NewEncoder(&given).Encode(tt.v); err != nil {
			t.Fatalf("(%d) unexpected encode error: %v", i, err)
		} else if tt.expected != given.String() {
			t.Fatalf("(%d) unexpected value: %s", i, given.String())
		}
	}
}

func TestSubscriptions_UpdateSubscription_Encoding(t *testing.T) {
	tests := []struct {
		v        recurly.UpdateSubscription
		expected string
	}{
		{
			expected: "<subscription></subscription>",
		},
		{
			v:        recurly.UpdateSubscription{Timeframe: "renewal"},
			expected: "<subscription><timeframe>renewal</timeframe></subscription>",
		},
		{
			v:        recurly.UpdateSubscription{PlanCode: "new-code"},
			expected: "<subscription><plan_code>new-code</plan_code></subscription>",
		},
		{
			v:        recurly.UpdateSubscription{Quantity: 14},
			expected: "<subscription><quantity>14</quantity></subscription>",
		},
		{
			v:        recurly.UpdateSubscription{UnitAmountInCents: 3500},
			expected: "<subscription><unit_amount_in_cents>3500</unit_amount_in_cents></subscription>",
		},
		{
			v:        recurly.UpdateSubscription{CollectionMethod: "manual"},
			expected: "<subscription><collection_method>manual</collection_method></subscription>",
		},
		{
			v:        recurly.UpdateSubscription{NetTerms: recurly.NewInt(0)},
			expected: "<subscription><net_terms>0</net_terms></subscription>",
		},
		{
			v:        recurly.UpdateSubscription{PONumber: "AB-NewPO"},
			expected: "<subscription><po_number>AB-NewPO</po_number></subscription>",
		},
		{
			v: recurly.UpdateSubscription{SubscriptionAddOns: &[]recurly.SubscriptionAddOn{
				{
					Code:              "extra_users",
					UnitAmountInCents: 1000,
					Quantity:          2,
				},
			}},
			expected: "<subscription><subscription_add_ons><subscription_add_on><add_on_code>extra_users</add_on_code><unit_amount_in_cents>1000</unit_amount_in_cents><quantity>2</quantity></subscription_add_on></subscription_add_ons></subscription>",
		},
		{
			v: recurly.Subscription{
				SubscriptionAddOns: []recurly.SubscriptionAddOn{
					{
						Code:              "extra_users",
						UnitAmountInCents: 1000,
						Quantity:          2,
					},
				},
				PONumber: "abc-123",
				NetTerms: recurly.NewInt(23),
			}.MakeUpdate(),
			expected: "<subscription><net_terms>23</net_terms><subscription_add_ons><subscription_add_on><add_on_code>extra_users</add_on_code><unit_amount_in_cents>1000</unit_amount_in_cents><quantity>2</quantity></subscription_add_on></subscription_add_ons></subscription>",
		},
	}
	for i, tt := range tests {
		var given bytes.Buffer
		if err := xml.NewEncoder(&given).Encode(tt.v); err != nil {
			t.Fatalf("(%d) unexpected encode error: %v", i, err)
		} else if tt.expected != given.String() {
			t.Fatalf("(%d) unexpected value: %s", i, given.String())
		}
	}
}

func TestSubscriptions_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/subscriptions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
		<subscriptions type="array">
			<subscription href="https://your-subdomain.recurly.com/v2/subscriptions/44f83d7cba354d5b84812419f923ea96">
				<account href="https://your-subdomain.recurly.com/v2/accounts/1"/>
				<invoice href="https://your-subdomain.recurly.com/v2/invoices/1108"/>
				<plan href="https://your-subdomain.recurly.com/v2/plans/gold">
				  <plan_code>gold</plan_code>
				  <name>Gold plan</name>
				</plan>
				<uuid>44f83d7cba354d5b84812419f923ea96</uuid>
				<state>active</state>
				<unit_amount_in_cents type="integer">800</unit_amount_in_cents>
				<currency>EUR</currency>
				<quantity type="integer">1</quantity>
				<activated_at type="datetime">2011-05-27T07:00:00Z</activated_at>
				<canceled_at nil="nil"></canceled_at>
				<expires_at nil="nil"></expires_at>
				<current_period_started_at type="datetime">2011-06-27T07:00:00Z</current_period_started_at>
				<current_period_ends_at type="datetime">2010-07-27T07:00:00Z</current_period_ends_at>
				<trial_started_at nil="nil"></trial_started_at>
				<trial_ends_at nil="nil"></trial_ends_at>
				<tax_in_cents type="integer">72</tax_in_cents>
				<tax_type>usst</tax_type>
				<tax_region>CA</tax_region>
				<tax_rate type="float">0.0875</tax_rate>
				<po_number nil="nil"></po_number>
				<net_terms type="integer">0</net_terms>
				<subscription_add_ons type="array">
					<subscription_add_on>
						<add_on_type>fixed</add_on_type>
						<add_on_code>my_add_on</add_on_code>
						<unit_amount_in_cents type="integer">1</unit_amount_in_cents>
						<quantity type="integer">1</quantity>
					</subscription_add_on>
				</subscription_add_ons>
				<a name="cancel" href="https://your-subdomain.recurly.com/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/cancel" method="put"/>
				<a name="terminate" href="https://your-subdomain.recurly.com/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/terminate" method="put"/>
				<a name="postpone" href="https://your-subdomain.recurly.com/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/postpone" method="put"/>
			</subscription>
		</subscriptions>`)
	})

	r, subscriptions, err := client.Subscriptions.List(recurly.Params{"per_page": 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected list subcriptions to return OK")
	} else if pp := r.Request.URL.Query().Get("per_page"); pp != "1" {
		t.Fatalf("unexpected per_page: %s", pp)
	}

	activated, _ := time.Parse(recurly.DateTimeFormat, "2011-05-27T07:00:00Z")
	cpStartedAt, _ := time.Parse(recurly.DateTimeFormat, "2011-06-27T07:00:00Z")
	cpEndsAt, _ := time.Parse(recurly.DateTimeFormat, "2010-07-27T07:00:00Z")

	if !reflect.DeepEqual(subscriptions, []recurly.Subscription{
		{
			XMLName: xml.Name{Local: "subscription"},
			Plan: recurly.NestedPlan{
				Code: "gold",
				Name: "Gold plan",
			},
			AccountCode:            "1",
			InvoiceNumber:          1108,
			UUID:                   "44f83d7cba354d5b84812419f923ea96",
			State:                  "active",
			UnitAmountInCents:      800,
			Currency:               "EUR",
			Quantity:               1,
			ActivatedAt:            recurly.NewTime(activated),
			CurrentPeriodStartedAt: recurly.NewTime(cpStartedAt),
			CurrentPeriodEndsAt:    recurly.NewTime(cpEndsAt),
			TaxInCents:             72,
			TaxType:                "usst",
			TaxRegion:              "CA",
			TaxRate:                0.0875,
			NetTerms:               recurly.NewInt(0),
			SubscriptionAddOns: []recurly.SubscriptionAddOn{
				{
					XMLName:           xml.Name{Local: "subscription_add_on"},
					Type:              "fixed",
					Code:              "my_add_on",
					Quantity:          1,
					UnitAmountInCents: 1,
				},
			},
		},
	}) {
		t.Fatalf("unexpected subscriptions: %v", subscriptions)
	}
}

func TestSubscriptions_ListAccount(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/1/subscriptions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
		<subscriptions type="array">
			<subscription href="https://your-subdomain.recurly.com/v2/subscriptions/44f83d7cba354d5b84812419f923ea96">
				<account href="https://your-subdomain.recurly.com/v2/accounts/1"/>
				<invoice href="https://your-subdomain.recurly.com/v2/invoices/1108"/>
				<plan href="https://your-subdomain.recurly.com/v2/plans/gold">
				  <plan_code>gold</plan_code>
				  <name>Gold plan</name>
				</plan>
				<uuid>44f83d7cba354d5b84812419f923ea96</uuid>
				<state>active</state>
				<unit_amount_in_cents type="integer">800</unit_amount_in_cents>
				<currency>EUR</currency>
				<quantity type="integer">1</quantity>
				<activated_at type="datetime">2011-05-27T07:00:00Z</activated_at>
				<canceled_at nil="nil"></canceled_at>
				<expires_at nil="nil"></expires_at>
				<current_period_started_at type="datetime">2011-06-27T07:00:00Z</current_period_started_at>
				<current_period_ends_at type="datetime">2010-07-27T07:00:00Z</current_period_ends_at>
				<trial_started_at nil="nil"></trial_started_at>
				<trial_ends_at nil="nil"></trial_ends_at>
				<tax_in_cents type="integer">72</tax_in_cents>
				<tax_type>usst</tax_type>
				<tax_region>CA</tax_region>
				<tax_rate type="float">0.0875</tax_rate>
				<po_number nil="nil"></po_number>
				<net_terms type="integer">0</net_terms>
				<subscription_add_ons type="array">
				</subscription_add_ons>
				<a name="cancel" href="https://your-subdomain.recurly.com/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/cancel" method="put"/>
				<a name="terminate" href="https://your-subdomain.recurly.com/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/terminate" method="put"/>
				<a name="postpone" href="https://your-subdomain.recurly.com/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/postpone" method="put"/>
			</subscription>
		</subscriptions>`)
	})

	r, subscriptions, err := client.Subscriptions.ListAccount("1", recurly.Params{"per_page": 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected list subcriptions to return OK")
	} else if pp := r.Request.URL.Query().Get("per_page"); pp != "1" {
		t.Fatalf("unexpected per_page: %s", pp)
	}

	activated, _ := time.Parse(recurly.DateTimeFormat, "2011-05-27T07:00:00Z")
	cpStartedAt, _ := time.Parse(recurly.DateTimeFormat, "2011-06-27T07:00:00Z")
	cpEndsAt, _ := time.Parse(recurly.DateTimeFormat, "2010-07-27T07:00:00Z")

	if !reflect.DeepEqual(subscriptions, []recurly.Subscription{
		{
			XMLName: xml.Name{Local: "subscription"},
			Plan: recurly.NestedPlan{
				Code: "gold",
				Name: "Gold plan",
			},
			AccountCode:            "1",
			InvoiceNumber:          1108,
			UUID:                   "44f83d7cba354d5b84812419f923ea96",
			State:                  "active",
			UnitAmountInCents:      800,
			Currency:               "EUR",
			Quantity:               1,
			ActivatedAt:            recurly.NewTime(activated),
			CurrentPeriodStartedAt: recurly.NewTime(cpStartedAt),
			CurrentPeriodEndsAt:    recurly.NewTime(cpEndsAt),
			TaxInCents:             72,
			TaxType:                "usst",
			TaxRegion:              "CA",
			TaxRate:                0.0875,
			NetTerms:               recurly.NewInt(0),
		},
	}) {
		t.Fatalf("unexpected subscriptions: %v", subscriptions)
	}
}

func TestSubscriptions_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/subscriptions/44f83d7cba354d5b84812419f923ea96", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
		<subscription href="https://your-subdomain.recurly.com/v2/subscriptions/44f83d7cba354d5b84812419f923ea96">
			<account href="https://your-subdomain.recurly.com/v2/accounts/1"/>
			<invoice href="https://your-subdomain.recurly.com/v2/invoices/1108"/>
			<plan href="https://your-subdomain.recurly.com/v2/plans/gold">
			  <plan_code>gold</plan_code>
			  <name>Gold plan</name>
			</plan>
			<uuid>44f83d7cba354d5b84812419f923ea96</uuid>
			<state>active</state>
			<unit_amount_in_cents type="integer">800</unit_amount_in_cents>
			<currency>EUR</currency>
			<quantity type="integer">1</quantity>
			<activated_at type="datetime">2011-05-27T07:00:00Z</activated_at>
			<canceled_at nil="nil"></canceled_at>
			<expires_at nil="nil"></expires_at>
			<current_period_started_at type="datetime">2011-06-27T07:00:00Z</current_period_started_at>
			<current_period_ends_at type="datetime">2011-07-27T07:00:00Z</current_period_ends_at>
			<trial_started_at nil="nil"></trial_started_at>
			<trial_ends_at nil="nil"></trial_ends_at>
			<tax_in_cents type="integer">72</tax_in_cents>
			<tax_type>usst</tax_type>
			<tax_region>CA</tax_region>
			<tax_rate type="float">0.0875</tax_rate>
			<po_number nil="nil"></po_number>
			<net_terms type="integer">0</net_terms>
			<subscription_add_ons type="array">
			</subscription_add_ons>
			<a name="cancel" href="https://your-subdomain.recurly.com/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/cancel" method="put"/>
			<a name="terminate" href="https://your-subdomain.recurly.com/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/terminate" method="put"/>
			<a name="postpone" href="https://your-subdomain.recurly.com/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/postpone" method="put"/>
		</subscription>`)
	})

	r, subscription, err := client.Subscriptions.Get("44f83d7cba354d5b84812419f923ea96")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected list subcriptions to return OK")
	}

	if !reflect.DeepEqual(subscription, &recurly.Subscription{
		XMLName: xml.Name{Local: "subscription"},
		Plan: recurly.NestedPlan{
			Code: "gold",
			Name: "Gold plan",
		},
		AccountCode:            "1",
		InvoiceNumber:          1108,
		UUID:                   "44f83d7cba354d5b84812419f923ea96",
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
	}) {
		t.Fatalf("unexpected subscription: %v", subscription)
	}
}

func TestSubscriptions_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/subscriptions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(201)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><subscription></subscription>`)
	})

	r, _, err := client.Subscriptions.Create(recurly.NewSubscription{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected create subscription to return OK")
	}
}

func TestSubscriptions_Create_TransactionError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/subscriptions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(422)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
			<errors>
			  <transaction_error>
			    <error_code>fraud_security_code</error_code>
			  </transaction_error>
			  <error field="transaction.account.billing_info.verification_value" symbol="declined_bad">did not match</error>
			  <transaction href="https://your-subdomain.recurly.com/v2/transactions/3d1c6aa86e3d447eb0f3b4a6e3e074d9" type="credit_card">
			    <uuid>3d1c6aa86e3d447eb0f3b4a6e3e074d9</uuid>
			  </transaction>
			</errors>`)
	})

	r, _, err := client.Subscriptions.Create(recurly.NewSubscription{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if !r.IsError() {
		t.Fatal("expected create subscription to return OK")
	} else if r.Transaction.UUID != "3d1c6aa86e3d447eb0f3b4a6e3e074d9" {
		t.Fatalf("unexpected transaction uuid: %s", r.Transaction.UUID)
	} else if r.TransactionError.ErrorCode != "fraud_security_code" {
		t.Fatalf("unexpected transaction error code: %s", r.TransactionError.ErrorCode)
	}
}

func TestSubscriptions_Preview(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/subscriptions/preview", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(201)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><subscription></subscription>`)
	})

	r, _, err := client.Subscriptions.Preview(recurly.NewSubscription{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected preview subscription to return OK")
	}
}

func TestSubscriptions_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/subscriptions/44f83d7cba354d5b84812419f923ea96", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><subscription></subscription>`)
	})

	r, _, err := client.Subscriptions.Update("44f83d7cba354d5b84812419f923ea96", recurly.UpdateSubscription{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected update subscription to return OK")
	}
}

func TestSubscriptions_Notes(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/notes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><subscription></subscription>`)
	})

	r, _, err := client.Subscriptions.UpdateNotes("44f83d7cba354d5b84812419f923ea96", recurly.SubscriptionNotes{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected update subscription notes to return OK")
	}
}

func TestSubscriptions_Change(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/preview", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(201)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><subscription></subscription>`)
	})

	r, _, err := client.Subscriptions.PreviewChange("44f83d7cba354d5b84812419f923ea96", recurly.UpdateSubscription{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected preview subscription change to return OK")
	}
}

func TestSubscriptions_Cancel(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/cancel", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><subscription></subscription>`)
	})

	r, _, err := client.Subscriptions.Cancel("44f83d7cba354d5b84812419f923ea96")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected cancel subscription change to return OK")
	}
}

func TestSubscriptions_Reactivate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/reactivate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><subscription></subscription>`)
	})

	r, _, err := client.Subscriptions.Reactivate("44f83d7cba354d5b84812419f923ea96")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected reactivate subscription change to return OK")
	}
}

func TestSubscriptions_Terminate_PartialRefund(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/terminate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("unexpected method: %s", r.Method)
		} else if refundType := r.URL.Query().Get("refund_type"); refundType != "partial" {
			t.Fatalf("unexpected input for refund_type: %s", refundType)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><subscription></subscription>`)
	})

	r, _, err := client.Subscriptions.TerminateWithPartialRefund("44f83d7cba354d5b84812419f923ea96")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected terminate subscription with partial refund to return OK")
	}
}

func TestSubscriptions_Terminate_FullRefund(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/terminate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("unexpected method: %s", r.Method)
		} else if refundType := r.URL.Query().Get("refund_type"); refundType != "full" {
			t.Fatalf("unexpected input for refund_type: %s", refundType)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><subscription></subscription>`)
	})

	r, _, err := client.Subscriptions.TerminateWithFullRefund("44f83d7cba354d5b84812419f923ea96")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected terminate subscription with full refund to return OK")
	}
}

func TestSubscriptions_Terminate_WithoutRefund(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/terminate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("unexpected method: %s", r.Method)
		} else if refundType := r.URL.Query().Get("refund_type"); refundType != "none" {
			t.Fatalf("unexpected input for refund_type: %s", refundType)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><subscription></subscription>`)
	})

	r, _, err := client.Subscriptions.TerminateWithoutRefund("44f83d7cba354d5b84812419f923ea96")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected terminate subscription without refund to return OK")
	}
}

func TestSubscriptions_Postpone(t *testing.T) {
	setup()
	defer teardown()

	ts, _ := time.Parse(recurly.DateTimeFormat, "2015-08-27T07:00:00Z")
	mux.HandleFunc("/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/postpone", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("unexpected method: %s", r.Method)
		} else if nrd := r.URL.Query().Get("next_renewal_date"); nrd != "2015-08-27T07:00:00Z" {
			t.Fatalf("unexpected input for next_renewal_date: %s", nrd)
		} else if bulk := r.URL.Query().Get("bulk"); bulk != "false" {
			t.Fatalf("unexpected input for bulk: %s", bulk)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><subscription></subscription>`)
	})

	r, _, err := client.Subscriptions.Postpone("44f83d7cba354d5b84812419f923ea96", ts, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected postpone subscription change to return OK")
	}
}
