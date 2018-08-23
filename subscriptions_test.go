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

// TestSubscriptionsEncoding ensures structs are encoded to XML properly.
// Because Recurly supports partial updates, it's important that only defined
// fields are handled properly -- including types like booleans and integers which
// have zero values that we want to send.
func TestSubscriptions_NewSubscription_Encoding(t *testing.T) {
	ts, _ := time.Parse(recurly.DateTimeFormat, "2015-06-03T13:42:23.764061Z")
	var customFields = &recurly.CustomFields{
		"platform": "2.0",
		"seller":   "Recurly Merchant",
	}
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
				PlanCode:             "gold",
				AutoRenew:            true,
				RenewalBillingCycles: recurly.NullInt{Valid: true, Int: 2},
				Account: recurly.Account{
					Code: "123",
					BillingInfo: &recurly.Billing{
						Token: "507c7f79bcf86cd7994f6c0e",
					},
				},
			},
			expected: "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code><billing_info><token_id>507c7f79bcf86cd7994f6c0e</token_id></billing_info></account><currency></currency><renewal_billing_cycles>2</renewal_billing_cycles><auto_renew>true</auto_renew></subscription>",
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
		{
			v: recurly.NewSubscription{
				PlanCode: "gold",
				Currency: "USD",
				Account: recurly.Account{
					Code: "123",
				},
				BankAccountAuthorizedAt: recurly.NewTime(ts),
				CustomFields:            customFields,
			},
			expected: "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency><bank_account_authorized_at>2015-06-03T13:42:23Z</bank_account_authorized_at><custom_fields><custom_field><name>platform</name><value>2.0</value></custom_field><custom_field><name>seller</name><value>Recurly Merchant</value></custom_field></custom_fields></subscription>",
		},
	}

	for i, tt := range tests {
		var given bytes.Buffer
		if err := xml.NewEncoder(&given).Encode(tt.v); err != nil {
			t.Fatalf("(%d) unexpected encode error: %v", i, err)
		} else if diff := cmp.Diff(tt.expected, given.String()); diff != "" {
			t.Fatal(diff)
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
			v:        recurly.UpdateSubscription{Quantity: 14, AutoRenew: true, RenewalBillingCycles: recurly.NullInt{Valid: true, Int: 2}},
			expected: "<subscription><quantity>14</quantity><renewal_billing_cycles>2</renewal_billing_cycles><auto_renew>true</auto_renew></subscription>",
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
		} else if diff := cmp.Diff(tt.expected, given.String()); diff != "" {
			t.Fatal(diff)
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
				<auto_renew>true</auto_renew>
				<current_term_started_at type="datetime">2010-07-27T07:00:00Z</current_term_started_at>
				<current_term_ends_at type="datetime">2010-07-27T07:00:00Z</current_term_ends_at>
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

	if diff := cmp.Diff(subscriptions, []recurly.Subscription{
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
			CurrentTermStartedAt:   recurly.NewTime(cpEndsAt),
			CurrentTermEndsAt:      recurly.NewTime(cpEndsAt),
			TaxInCents:             72,
			TaxType:                "usst",
			TaxRegion:              "CA",
			TaxRate:                0.0875,
			AutoRenew:              true,
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
	}); diff != "" {
		t.Fatal(diff)
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

	if diff := cmp.Diff(subscriptions, []recurly.Subscription{
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
	}); diff != "" {
		t.Fatal(diff)
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

	r, subscription, err := client.Subscriptions.Get("44f83d7cb-a354d5b848-12419f923ea96") // UUID has dashes
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected list subcriptions to return OK")
	}

	if diff := cmp.Diff(subscription, &recurly.Subscription{
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
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestSubscriptions_Get_ErrNotFound(t *testing.T) {
	setup()
	defer teardown()

	var invoked bool
	mux.HandleFunc("/v2/subscriptions/44f83d7cba354d5b84812419f923ea96", func(w http.ResponseWriter, r *http.Request) {
		invoked = true
		w.WriteHeader(http.StatusNotFound)
	})

	_, subscription, err := client.Subscriptions.Get("44f83d7cba354d5b84812419f923ea96")
	if !invoked {
		t.Fatal("handler not invoked")
	} else if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if subscription != nil {
		t.Fatalf("expected subscription to be nil: %#v", subscription)
	}
}

func TestSubscriptions_Get_PendingSubscription(t *testing.T) {
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
			<pending_subscription type="subscription">
				<plan href="https://blacklighttest.recurly.com/v2/plans/integrationtestcode">
					<plan_code>gold</plan_code>
					<name>Gold plan</name>
				</plan>
				<unit_amount_in_cents type="integer">50000</unit_amount_in_cents>
				<quantity type="integer">1</quantity>
				<subscription_add_ons type="array">
					<subscription_add_on>
						<add_on_type>fixed</add_on_type>
						<add_on_code>add-on-one</add_on_code>
						<unit_amount_in_cents type="integer">1100</unit_amount_in_cents>
						<quantity type="integer">1</quantity>
					</subscription_add_on>
					<subscription_add_on>
						<add_on_type>fixed</add_on_type>
						<add_on_code>add-on-two</add_on_code>
						<unit_amount_in_cents type="integer">1300</unit_amount_in_cents>
						<quantity type="integer">1</quantity>
					</subscription_add_on>
				</subscription_add_ons>
			</pending_subscription>
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

	if diff := cmp.Diff(subscription, &recurly.Subscription{
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
					UnitAmountInCents: 1100,
					Quantity:          1,
				},
				{
					XMLName:           xml.Name{Local: "subscription_add_on"},
					Type:              "fixed",
					Code:              "add-on-two",
					UnitAmountInCents: 1300,
					Quantity:          1,
				},
			},
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestSubscriptions_Get_With_CustomFields(t *testing.T) {
	setup()
	defer teardown()

	var customFields = &recurly.CustomFields{
		"device_id":     "KIWTL-WER-ZXMRD",
		"purchase_date": "2017-01-23",
	}

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
			<custom_fields type="array">
			  <custom_field>
			    <name>device_id</name>
			    <value>KIWTL-WER-ZXMRD</value>
			  </custom_field>
			  <custom_field>
			    <name>purchase_date</name>
			    <value>2017-01-23</value>
			  </custom_field>
			</custom_fields>
			<a name="cancel" href="https://your-subdomain.recurly.com/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/cancel" method="put"/>
			<a name="terminate" href="https://your-subdomain.recurly.com/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/terminate" method="put"/>
			<a name="postpone" href="https://your-subdomain.recurly.com/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/postpone" method="put"/>
		</subscription>`)
	})

	r, subscription, err := client.Subscriptions.Get("44f83d7cb-a354d5b848-12419f923ea96") // UUID has dashes
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected list subcriptions to return OK")
	}

	if diff := cmp.Diff(subscription, &recurly.Subscription{
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
		CustomFields:           customFields,
	}); diff != "" {
		t.Fatal(diff)
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
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
			<errors>
			  <transaction_error>
			    <error_code>declined_saveable</error_code>
			    <error_category>soft</error_category>
			    <merchant_message>The transaction was declined without specific information.  Please contact your payment gateway for more details or ask the customer to contact their bank.</merchant_message>
			    <customer_message>The transaction was declined. Please use a different card or contact your bank.</customer_message>
			    <gateway_error_code nil="nil"/>
			  </transaction_error>
			  <error field="subscription.account.base" symbol="declined_saveable">The transaction was declined. Please use a different card or contact your bank.</error>
			  <transaction href="https://your-subdomain.recurly.com/v2/transactions/3c42a3ecc46a7aa602602e4033b9c2e6" type="credit_card">
			    <account href="https://your-subdomain.recurly.com/v2/accounts/1"/>
			    <subscription href="https://your-subdomain.recurly.com/v2/subscriptions/3c42a3ebabdc022739d5a646408291a6"/>
			    <uuid>3c42a3ecc46a7aa602602e4033b9c2e6</uuid>
			    <action>purchase</action>
			    <amount_in_cents type="integer">5274</amount_in_cents>
			    <currency>USD</currency>
			    <status>declined</status>
			    <payment_method>credit_card</payment_method>
			    <transaction_error>
			      <error_code>declined_saveable</error_code>
			      <error_category>soft</error_category>
			      <merchant_message>The transaction was declined without specific information.  Please contact your payment gateway for more details or ask the customer to contact their bank.</merchant_message>
			      <customer_message>The transaction was declined. Please use a different card or contact your bank.</customer_message>
			      <gateway_error_code nil="nil"/>
			    </transaction_error>
				<details>
			    </details>
			  </transaction>
			</errors>`)
	})

	r, newSubscription, err := client.Subscriptions.Create(recurly.NewSubscription{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if !r.IsError() {
		t.Fatal("expected create subscription to return OK")
	} else if newSubscription.Transaction == nil {
		t.Fatal("expected transaction to be set")
	} else if diff := cmp.Diff(newSubscription.Transaction, &recurly.Transaction{
		UUID:          "3c42a3ecc46a7aa602602e4033b9c2e6",
		Action:        "purchase",
		AmountInCents: 5274,
		Currency:      "USD",
		Status:        "declined",
		PaymentMethod: "credit_card",
		TransactionError: &recurly.TransactionError{
			XMLName:         xml.Name{Local: "transaction_error"},
			ErrorCode:       "declined_saveable",
			ErrorCategory:   "soft",
			MerchantMessage: "The transaction was declined without specific information.  Please contact your payment gateway for more details or ask the customer to contact their bank.",
			CustomerMessage: "The transaction was declined. Please use a different card or contact your bank.",
		},
	}); diff != "" {
		t.Fatal(diff)
	} else if newSubscription.Subscription != nil {
		t.Fatalf("unexpected subscription: %v", newSubscription.Subscription)
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
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><subscription>
			<invoice_collection>
				<charge_invoice href="">
				<account href="https://your-subdomain.recurly.com/v2/accounts/1"/>
				<uuid>43adfe52c21cbb221557a24940bcd7e5</uuid>
				<state>pending</state>
				</charge_invoice>
				<credit_invoices type="array">
				</credit_invoices>
			</invoice_collection>
		</subscription>`)
	})

	r, sub, err := client.Subscriptions.Preview(recurly.NewSubscription{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected preview subscription to return OK")
	} else if diff := cmp.Diff(sub, &recurly.Subscription{
		XMLName: xml.Name{Local: "subscription"},
		InvoiceCollection: &recurly.InvoiceCollection{
			XMLName: xml.Name{Local: "invoice_collection"},
			ChargeInvoice: &recurly.Invoice{
				XMLName:     xml.Name{Local: "invoice"},
				AccountCode: "1",
				UUID:        "43adfe52c21cbb221557a24940bcd7e5",
				State:       recurly.ChargeInvoiceStatePending,
			},
		},
	}); diff != "" {
		t.Fatal(diff)
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

	r, _, err := client.Subscriptions.Update("44f83d7cba-354d5b84812419-f923ea96", recurly.UpdateSubscription{}) // UUID has dashes and should be sanitized
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

	r, _, err := client.Subscriptions.UpdateNotes("44f83d7cba354d5-b8481241-9f923ea96", recurly.SubscriptionNotes{}) // UUID has dashes and should be sanitized
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected update subscription notes to return OK")
	}
}

func TestSubscriptions_PreviewChange(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/subscriptions/46652c59f3f1a1cdab9d404a6d91ddab/preview", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(201)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
			<subscription href="https://your-subdomain.recurly.com/v2/subscriptions/46652c59f3f1a1cdab9d404a6d91ddab">
				<account href="https://your-subdomain.recurly.com/v2/accounts/multi-test-one"/>
				<plan href="https://your-subdomain.recurly.com/v2/plans/multi-high">
					<plan_code>multi-high</plan_code>
					<name>Multi-Plan High End</name>
				</plan>
				<revenue_schedule_type>evenly</revenue_schedule_type>
				<uuid>46652c59f3f1a1cdab9d404a6d91ddab</uuid>
				<state>active</state>
				<unit_amount_in_cents type="integer">150000</unit_amount_in_cents>
				<currency>USD</currency>
				<quantity type="integer">1</quantity>
				<activated_at type="datetime">2018-08-01T16:45:16Z</activated_at>
				<canceled_at nil="nil"></canceled_at>
				<expires_at nil="nil"></expires_at>
				<updated_at type="datetime">2018-08-01T21:00:14Z</updated_at>
				<current_period_started_at type="datetime">2018-08-01T21:00:00Z</current_period_started_at>
				<current_period_ends_at type="datetime">2018-09-01T21:00:00Z</current_period_ends_at>
				<cost_in_cents type="integer">162000</cost_in_cents>
				<collection_method>automatic</collection_method>
				<subscription_add_ons type="array">
					<subscription_add_on>
						<add_on_type>fixed</add_on_type>
						<add_on_code>product</add_on_code>
						<unit_amount_in_cents type="integer">1500</unit_amount_in_cents>
						<quantity type="integer">8</quantity>
						<revenue_schedule_type>evenly</revenue_schedule_type>
					</subscription_add_on>
				</subscription_add_ons>
				<invoice_collection>
					<charge_invoice href="">
						<account href="https://your-subdomain.recurly.com/v2/accounts/multi-test-one"/>
						<uuid>4669ef730162b1a4db576a4722a40d8e</uuid>
						<state>pending</state>
						<total_in_cents type="integer">158094</total_in_cents>
						<currency>USD</currency>
						<attempt_next_collection_at type="datetime">2018-08-02T14:56:50Z</attempt_next_collection_at>
						<subtotal_before_discount_in_cents type="integer">158094</subtotal_before_discount_in_cents>
						<subtotal_in_cents type="integer">158094</subtotal_in_cents>
						<discount_in_cents type="integer">0</discount_in_cents>
						<due_on type="datetime">2018-08-02T14:56:50Z</due_on>
						<balance_in_cents type="integer">126865</balance_in_cents>
						<type>charge</type>
						<origin>immediate_change</origin>
						<refundable_total_in_cents type="integer">158094</refundable_total_in_cents>
						<credit_payments type="array">
							<credit_payment href="">
								<account href="https://your-subdomain.recurly.com/v2/accounts/multi-test-one"/>
								<uuid>4669ef73118cde865e2b974244887cf1</uuid>
								<action>payment</action>
								<currency>USD</currency>
								<amount_in_cents type="integer">31229</amount_in_cents>
							</credit_payment>
						</credit_payments>
						<collection_method>automatic</collection_method>
						<line_items type="array">
							<adjustment href="" type="charge">
								<account href="https://your-subdomain.recurly.com/v2/accounts/multi-test-one"/>
								<refundable_total_in_cents type="integer">146382</refundable_total_in_cents>
								<uuid>4669ef72d31522286d6e074d74869ed1</uuid>
								<state>pending</state>
								<description>Multi-Plan High End</description>
								<accounting_code nil="nil"></accounting_code>
								<product_code>multi-high</product_code>
								<origin>plan</origin>
								<unit_amount_in_cents type="integer">146382</unit_amount_in_cents>
								<quantity type="integer">1</quantity>
								<discount_in_cents type="integer">0</discount_in_cents>
								<tax_in_cents type="integer">0</tax_in_cents>
								<total_in_cents type="integer">146382</total_in_cents>
								<currency>USD</currency>
								<proration_rate type="float">0.9759</proration_rate>
								<start_date type="datetime">2018-08-02T14:56:50Z</start_date>
								<end_date type="datetime">2018-09-01T21:00:00Z</end_date>
								<created_at nil="nil"></created_at>
								<updated_at nil="nil"></updated_at>
								<revenue_schedule_type>evenly</revenue_schedule_type>
							</adjustment>
							<adjustment href="" type="charge">
								<account href="https://your-subdomain.recurly.com/v2/accounts/multi-test-one"/>
								<refundable_total_in_cents type="integer">11712</refundable_total_in_cents>
								<uuid>4669ef72d6a341370b3ca24f4eae2381</uuid>
								<state>pending</state>
								<description>desc</description>
								<accounting_code nil="nil"></accounting_code>
								<product_code>product</product_code>
								<origin>add_on</origin>
								<unit_amount_in_cents type="integer">1464</unit_amount_in_cents>
								<quantity type="integer">8</quantity>
								<total_in_cents type="integer">11712</total_in_cents>
								<currency>USD</currency>
								<proration_rate type="float">0.9759</proration_rate>
								<start_date type="datetime">2018-08-02T14:56:50Z</start_date>
								<end_date type="datetime">2018-09-01T21:00:00Z</end_date>
								<revenue_schedule_type>evenly</revenue_schedule_type>
							</adjustment>
						</line_items>
						<transactions type="array">
						</transactions>
					</charge_invoice>
					<credit_invoices type="array">
						<credit_invoice href="">
							<account href="https://your-subdomain.recurly.com/v2/accounts/multi-test-one"/>
							<uuid>4669ef72e546b859621857449e86d822</uuid>
							<state>closed</state>
							<total_in_cents type="integer">-31229</total_in_cents>
							<currency>USD</currency>
							<closed_at type="datetime">2018-08-02T14:56:50Z</closed_at>
							<subtotal_before_discount_in_cents type="integer">-31229</subtotal_before_discount_in_cents>
							<subtotal_in_cents type="integer">-31229</subtotal_in_cents>
							<type>credit</type>
							<origin>immediate_change</origin>
							<credit_payments type="array">
								<credit_payment href="">
									<account href="https://your-subdomain.recurly.com/v2/accounts/multi-test-one"/>
									<uuid>4669ef73118cde865e2b974244887cf1</uuid>
									<action>payment</action>
									<currency>USD</currency>
									<amount_in_cents type="integer">31229</amount_in_cents>
								</credit_payment>
							</credit_payments>
							<line_items type="array">
								<adjustment href="" type="credit">
									<account href="https://your-subdomain.recurly.com/v2/accounts/multi-test-one"/>
									<uuid>4669ef72bf5a190570a0fc4c37b9eef7</uuid>
									<state>pending</state>
									<description>Multi-Plan Low End</description>
									<product_code>multi-low</product_code>
									<origin>plan</origin>
									<unit_amount_in_cents type="integer">-19518</unit_amount_in_cents>
									<quantity type="integer">1</quantity>
									<quantity_remaining type="integer">1</quantity_remaining>
									<original_adjustment_uuid>466615c82c79f22fcbf024428fb69006</original_adjustment_uuid>
									<total_in_cents type="integer">-19518</total_in_cents>
									<currency>USD</currency>
									<proration_rate type="float">0.9759</proration_rate>
									<start_date type="datetime">2018-08-02T14:56:50Z</start_date>
									<end_date type="datetime">2018-09-01T21:00:00Z</end_date>
									<revenue_schedule_type>evenly</revenue_schedule_type>
									<credit_reason_code>refund</credit_reason_code>
								</adjustment>
								<adjustment href="" type="credit">
									<account href="https://your-subdomain.recurly.com/v2/accounts/multi-test-one"/>
									<uuid>4669ef72cbf2bb738fbdfb4c87bc88dc</uuid>
									<state>pending</state>
									<description>Domains</description>
									<product_code>domains</product_code>
									<origin>add_on</origin>
									<unit_amount_in_cents type="integer">-11711</unit_amount_in_cents>
									<quantity type="integer">1</quantity>
									<quantity_remaining type="integer">1</quantity_remaining>
									<original_adjustment_uuid>466615c83144d4334e855b4f5e8fcc9f</original_adjustment_uuid>
									<total_in_cents type="integer">-11711</total_in_cents>
									<currency>USD</currency>
									<proration_rate type="float">0.9759</proration_rate>
									<start_date type="datetime">2018-08-02T14:56:50Z</start_date>
									<end_date type="datetime">2018-09-01T21:00:00Z</end_date>
									<revenue_schedule_type>evenly</revenue_schedule_type>
									<credit_reason_code>refund</credit_reason_code>
								</adjustment>
							</line_items>
							<transactions type="array">
							</transactions>
						</credit_invoice>
					</credit_invoices>
				</invoice_collection>
			</subscription>`)
	})

	activatedAt := recurly.NewTimeFromString("2018-08-01T16:45:16Z")
	currentPeriodStartedAt := recurly.NewTimeFromString("2018-08-01T21:00:00Z")
	currentPeriodEndsAt := recurly.NewTimeFromString("2018-09-01T21:00:00Z")
	dueOn := recurly.NewTimeFromString("2018-08-02T14:56:50Z")
	r, sub, err := client.Subscriptions.PreviewChange("46652c59f3-f1a1cdab9d4-04a6d91ddab", recurly.UpdateSubscription{}) // UUID has dashes and should be sanitized
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected preview subscription change to return OK")
	} else if diff := cmp.Diff(sub, &recurly.Subscription{
		XMLName:     xml.Name{Local: "subscription"},
		UUID:        "46652c59f3f1a1cdab9d404a6d91ddab",
		AccountCode: "multi-test-one",
		Plan: recurly.NestedPlan{
			Code: "multi-high",
			Name: "Multi-Plan High End",
		},
		State:                  "active",
		UnitAmountInCents:      150000,
		Currency:               "USD",
		Quantity:               1,
		ActivatedAt:            activatedAt,
		CurrentPeriodStartedAt: currentPeriodStartedAt,
		CurrentPeriodEndsAt:    currentPeriodEndsAt,
		CollectionMethod:       "automatic",
		SubscriptionAddOns: []recurly.SubscriptionAddOn{
			{
				XMLName:           xml.Name{Local: "subscription_add_on"},
				Type:              "fixed",
				Code:              "product",
				UnitAmountInCents: 1500,
				Quantity:          8,
			},
		},
		InvoiceCollection: &recurly.InvoiceCollection{
			XMLName: xml.Name{Local: "invoice_collection"},
			ChargeInvoice: &recurly.Invoice{
				XMLName:                 xml.Name{Local: "invoice"},
				UUID:                    "4669ef730162b1a4db576a4722a40d8e",
				AccountCode:             "multi-test-one",
				State:                   "pending",
				TotalInCents:            158094,
				Currency:                "USD",
				AttemptNextCollectionAt: dueOn,
				SubtotalInCents:         158094,
				DueOn:                   dueOn,
				BalanceInCents:          126865,
				Type:                    "charge",
				Origin:                  "immediate_change",
				CreditPayments: []recurly.CreditPayment{{
					XMLName:       xml.Name{Local: "credit_payment"},
					AccountCode:   "multi-test-one",
					UUID:          "4669ef73118cde865e2b974244887cf1",
					Action:        "payment",
					Currency:      "USD",
					AmountInCents: 31229,
				}},
				CollectionMethod: "automatic",
				LineItems: []recurly.Adjustment{
					{
						UUID:              "4669ef72d31522286d6e074d74869ed1",
						AccountCode:       "multi-test-one",
						State:             "pending",
						Description:       "Multi-Plan High End",
						ProductCode:       "multi-high",
						Origin:            "plan",
						UnitAmountInCents: 146382,
						Quantity:          1,
						TotalInCents:      146382,
						Currency:          "USD",
						StartDate:         dueOn,
						EndDate:           currentPeriodEndsAt,
					},
					{
						UUID:              "4669ef72d6a341370b3ca24f4eae2381",
						AccountCode:       "multi-test-one",
						State:             "pending",
						Description:       "desc",
						ProductCode:       "product",
						Origin:            "add_on",
						UnitAmountInCents: 1464,
						Quantity:          8,
						TotalInCents:      11712,
						Currency:          "USD",
						StartDate:         dueOn,
						EndDate:           currentPeriodEndsAt,
					},
				},
			},
		},
	}); diff != "" {
		t.Fatal(diff)
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

	r, _, err := client.Subscriptions.Cancel("44f83d7cba-354d5b848124-19f923ea96") // UUID has dashes and should be sanitized
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

	r, _, err := client.Subscriptions.Reactivate("44f83d7cba35-4d5b8481241-9f923ea96") // UUID has dashes and should be sanitized
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

	r, _, err := client.Subscriptions.TerminateWithPartialRefund("44f83d7c-ba354d5b84812419f923ea96") // UUID has dashes and should be sanitized
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

	r, _, err := client.Subscriptions.TerminateWithFullRefund("44f83d7cba354d5b84-812419f923ea96") // UUID has dashes and should be sanitized
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

	r, _, err := client.Subscriptions.TerminateWithoutRefund("44f83d7c-ba354d5b84812419f923ea96") // UUID has dashes and should be sanitized
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

	r, _, err := client.Subscriptions.Postpone("44f83d7cba354d5b8481-2419f923ea96", ts, false) // UUID has dashes and should be sanitized
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected postpone subscription change to return OK")
	}
}

func TestSubscriptions_Pause(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/pause", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><subscription></subscription>`)
	})

	r, _, err := client.Subscriptions.Pause("44f83d7cba354d5b8481-2419f923ea96", 1) // UUID has dashes and should be sanitized
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected pause subscription change to return OK")
	}
}

func TestSubscriptions_Resume(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/resume", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><subscription></subscription>`)
	})

	r, _, err := client.Subscriptions.Resume("44f83d7cba354d5b8481-2419f923ea96") // UUID has dashes and should be sanitized
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected pause subscription change to return OK")
	}
}
