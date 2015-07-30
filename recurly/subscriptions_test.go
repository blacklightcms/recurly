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

// TestSubscriptionsEncoding ensures structs are encoded to XML properly.
// Because Recurly supports partial updates, it's important that only defined
// fields are handled properly -- including types like booleans and integers which
// have zero values that we want to send.
func TestSubscriptionsEncoding(t *testing.T) {
	ts, _ := time.Parse(datetimeFormat, "2015-06-03T13:42:23.764061Z")
	suite := []map[string]interface{}{
		// Plan code, account, and currency are required fields. They should always be present.
		map[string]interface{}{"struct": NewSubscription{}, "xml": "<subscription><plan_code></plan_code><account></account><currency></currency></subscription>"},
		map[string]interface{}{"struct": NewSubscription{
			PlanCode: "gold",
			Account: Account{
				Code: "123",
				BillingInfo: &Billing{
					Token: "507c7f79bcf86cd7994f6c0e",
				},
			},
		}, "xml": "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code><billing_info><token_id>507c7f79bcf86cd7994f6c0e</token_id></billing_info></account><currency></currency></subscription>"},
		map[string]interface{}{"struct": NewSubscription{
			PlanCode: "gold",
			Currency: "USD",
			Account: Account{
				Code: "123",
			},
			SubscriptionAddOns: &[]SubscriptionAddOn{
				SubscriptionAddOn{
					Code:              "extra_users",
					UnitAmountInCents: 1000,
					Quantity:          2,
				},
			},
		}, "xml": "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><subscription_add_ons><subscription_add_on><add_on_code>extra_users</add_on_code><unit_amount_in_cents>1000</unit_amount_in_cents><quantity>2</quantity></subscription_add_on></subscription_add_ons><currency>USD</currency></subscription>"},
		map[string]interface{}{"struct": NewSubscription{
			PlanCode: "gold",
			Currency: "USD",
			Account: Account{
				Code: "123",
			},
			CouponCode: "promo145",
		}, "xml": "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><coupon_code>promo145</coupon_code><currency>USD</currency></subscription>"},
		map[string]interface{}{"struct": NewSubscription{
			PlanCode: "gold",
			Currency: "USD",
			Account: Account{
				Code: "123",
			},
			UnitAmountInCents: 800,
		}, "xml": "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><unit_amount_in_cents>800</unit_amount_in_cents><currency>USD</currency></subscription>"},
		map[string]interface{}{"struct": NewSubscription{
			PlanCode: "gold",
			Currency: "USD",
			Account: Account{
				Code: "123",
			},
			Quantity: 8,
		}, "xml": "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency><quantity>8</quantity></subscription>"},
		map[string]interface{}{"struct": NewSubscription{
			PlanCode: "gold",
			Currency: "USD",
			Account: Account{
				Code: "123",
			},
			TrialEndsAt: NewTime(ts),
		}, "xml": "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency><trial_ends_at>2015-06-03T13:42:23Z</trial_ends_at></subscription>"},
		map[string]interface{}{"struct": NewSubscription{
			PlanCode: "gold",
			Currency: "USD",
			Account: Account{
				Code: "123",
			},
			StartsAt: NewTime(ts),
		}, "xml": "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency><starts_at>2015-06-03T13:42:23Z</starts_at></subscription>"},
		map[string]interface{}{"struct": NewSubscription{
			PlanCode: "gold",
			Currency: "USD",
			Account: Account{
				Code: "123",
			},
			TotalBillingCycles: 24,
		}, "xml": "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency><total_billing_cycles>24</total_billing_cycles></subscription>"},
		map[string]interface{}{"struct": NewSubscription{
			PlanCode: "gold",
			Currency: "USD",
			Account: Account{
				Code: "123",
			},
			FirstRenewalDate: NewTime(ts),
		}, "xml": "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency><first_renewal_date>2015-06-03T13:42:23Z</first_renewal_date></subscription>"},
		map[string]interface{}{"struct": NewSubscription{
			PlanCode: "gold",
			Currency: "USD",
			Account: Account{
				Code: "123",
			},
			CollectionMethod: "automatic",
		}, "xml": "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency><collection_method>automatic</collection_method></subscription>"},
		map[string]interface{}{"struct": NewSubscription{
			PlanCode: "gold",
			Currency: "USD",
			Account: Account{
				Code: "123",
			},
			NetTerms: NewInt(30),
		}, "xml": "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency><net_terms>30</net_terms></subscription>"},
		map[string]interface{}{"struct": NewSubscription{
			PlanCode: "gold",
			Currency: "USD",
			Account: Account{
				Code: "123",
			},
			NetTerms: NewInt(0),
		}, "xml": "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency><net_terms>0</net_terms></subscription>"},
		map[string]interface{}{"struct": NewSubscription{
			PlanCode: "gold",
			Currency: "USD",
			Account: Account{
				Code: "123",
			},
			PONumber: "PB4532345",
		}, "xml": "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency><po_number>PB4532345</po_number></subscription>"},
		map[string]interface{}{"struct": NewSubscription{
			PlanCode: "gold",
			Currency: "USD",
			Account: Account{
				Code: "123",
			},
			Bulk: true,
		}, "xml": "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency><bulk>true</bulk></subscription>"},
		map[string]interface{}{"struct": NewSubscription{
			PlanCode: "gold",
			Currency: "USD",
			Account: Account{
				Code: "123",
			},
			Bulk: false,
			// Bulk of false is the zero value of a bool, so it's omitted from the XML. But that's correct because Recurly's default is false
		}, "xml": "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency></subscription>"},
		map[string]interface{}{"struct": NewSubscription{
			PlanCode: "gold",
			Currency: "USD",
			Account: Account{
				Code: "123",
			},
			TermsAndConditions: "foo ... bar..",
		}, "xml": "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency><terms_and_conditions>foo ... bar..</terms_and_conditions></subscription>"},
		map[string]interface{}{"struct": NewSubscription{
			PlanCode: "gold",
			Currency: "USD",
			Account: Account{
				Code: "123",
			},
			CustomerNotes: "foo ... customer.. bar",
		}, "xml": "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency><customer_notes>foo ... customer.. bar</customer_notes></subscription>"},
		map[string]interface{}{"struct": NewSubscription{
			PlanCode: "gold",
			Currency: "USD",
			Account: Account{
				Code: "123",
			},
			VATReverseChargeNotes: "foo ... VAT.. bar",
		}, "xml": "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency><vat_reverse_charge_notes>foo ... VAT.. bar</vat_reverse_charge_notes></subscription>"},
		map[string]interface{}{"struct": NewSubscription{
			PlanCode: "gold",
			Currency: "USD",
			Account: Account{
				Code: "123",
			},
			BankAccountAuthorizedAt: NewTime(ts),
		}, "xml": "<subscription><plan_code>gold</plan_code><account><account_code>123</account_code></account><currency>USD</currency><bank_account_authorized_at>2015-06-03T13:42:23Z</bank_account_authorized_at></subscription>"},

		// Update Subscription Tests
		map[string]interface{}{"struct": UpdateSubscription{}, "xml": "<subscription></subscription>"},
		map[string]interface{}{"struct": UpdateSubscription{Timeframe: "renewal"}, "xml": "<subscription><timeframe>renewal</timeframe></subscription>"},
		map[string]interface{}{"struct": UpdateSubscription{PlanCode: "new-code"}, "xml": "<subscription><plan_code>new-code</plan_code></subscription>"},
		map[string]interface{}{"struct": UpdateSubscription{Quantity: 14}, "xml": "<subscription><quantity>14</quantity></subscription>"},
		map[string]interface{}{"struct": UpdateSubscription{UnitAmountInCents: 3500}, "xml": "<subscription><unit_amount_in_cents>3500</unit_amount_in_cents></subscription>"},
		map[string]interface{}{"struct": UpdateSubscription{CollectionMethod: "manual"}, "xml": "<subscription><collection_method>manual</collection_method></subscription>"},
		map[string]interface{}{"struct": UpdateSubscription{NetTerms: NewInt(0)}, "xml": "<subscription><net_terms>0</net_terms></subscription>"},
		map[string]interface{}{"struct": UpdateSubscription{PONumber: "AB-NewPO"}, "xml": "<subscription><po_number>AB-NewPO</po_number></subscription>"},
		map[string]interface{}{"struct": UpdateSubscription{SubscriptionAddOns: &[]SubscriptionAddOn{
			SubscriptionAddOn{
				Code:              "extra_users",
				UnitAmountInCents: 1000,
				Quantity:          2,
			},
		}}, "xml": "<subscription><subscription_add_ons><subscription_add_on><add_on_code>extra_users</add_on_code><unit_amount_in_cents>1000</unit_amount_in_cents><quantity>2</quantity></subscription_add_on></subscription_add_ons></subscription>"},
		map[string]interface{}{"struct": Subscription{
			SubscriptionAddOns: &[]SubscriptionAddOn{
				SubscriptionAddOn{
					Code:              "extra_users",
					UnitAmountInCents: 1000,
					Quantity:          2,
				},
			},
			PONumber: "abc-123",
			NetTerms: NewInt(23),
		}.MakeUpdate(), "xml": "<subscription><net_terms>23</net_terms><subscription_add_ons><subscription_add_on><add_on_code>extra_users</add_on_code><unit_amount_in_cents>1000</unit_amount_in_cents><quantity>2</quantity></subscription_add_on></subscription_add_ons></subscription>"},
	}

	for i, s := range suite {
		given := new(bytes.Buffer)
		err := xml.NewEncoder(given).Encode(s["struct"])
		if err != nil {
			t.Errorf("TestSubscriptionsEncoding Error (%d): %s", i, err)
		}

		if s["xml"] != given.String() {
			t.Errorf("TestSubscriptionsEncoding Error (%d): Expected %s, given %s", i, s["xml"], given.String())
		}
	}
}

func TestListSubscriptions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/subscriptions", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("TestListSubscriptions Error: Expected %s request, given %s", "GET", r.Method)
		}
		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?>
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

	r, subscriptions, err := client.Subscriptions.List(Params{"per_page": 1})
	if err != nil {
		t.Errorf("TestListSubscriptions Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestListSubscriptions Error: Expected list subcriptions to return OK")
	}

	if len(subscriptions) != 1 {
		t.Fatalf("TestListSubscriptions Error: Expected 1 subscription returned, given %d", len(subscriptions))
	}

	if r.Request.URL.Query().Get("per_page") != "1" {
		t.Errorf("TestListSubscriptions Error: Expected per_page parameter of 1, given %s", r.Request.URL.Query().Get("per_page"))
	}

	activated, _ := time.Parse(datetimeFormat, "2011-05-27T07:00:00Z")
	cpStartedAt, _ := time.Parse(datetimeFormat, "2011-06-27T07:00:00Z")
	cpEndsAt, _ := time.Parse(datetimeFormat, "2010-07-27T07:00:00Z")
	for _, given := range subscriptions {
		expected := Subscription{
			XMLName: xml.Name{Local: "subscription"},
			Plan: nestedPlan{
				Code: "gold",
				Name: "Gold plan",
			},
			Account: href{
				HREF: "https://your-subdomain.recurly.com/v2/accounts/1",
				Code: "1",
			},
			Invoice: href{
				HREF: "https://your-subdomain.recurly.com/v2/invoices/1108",
				Code: "1108",
			},
			UUID:                   "44f83d7cba354d5b84812419f923ea96",
			State:                  "active",
			UnitAmountInCents:      800,
			Currency:               "EUR",
			Quantity:               1,
			ActivatedAt:            NewTime(activated),
			CurrentPeriodStartedAt: NewTime(cpStartedAt),
			CurrentPeriodEndsAt:    NewTime(cpEndsAt),
			TaxInCents:             72,
			TaxType:                "usst",
			TaxRegion:              "CA",
			TaxRate:                0.0875,
			NetTerms:               NewInt(0),
		}

		if !reflect.DeepEqual(expected, given) {
			t.Errorf("TestListSubscriptions Error: expected subscription to equal %#v, given %#v", expected, given)
		}
	}
}

func TestListAccountSubscriptions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/1/subscriptions", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("TestListAccountSubscriptions Error: Expected %s request, given %s", "GET", r.Method)
		}
		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?>
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

	r, subscriptions, err := client.Subscriptions.ListForAccount("1", Params{"per_page": 1})
	if err != nil {
		t.Errorf("TestListAccountSubscriptions Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestListAccountSubscriptions Error: Expected list subcriptions to return OK")
	}

	if len(subscriptions) != 1 {
		t.Fatalf("TestListAccountSubscriptions Error: Expected 1 subscription returned, given %d", len(subscriptions))
	}

	if r.Request.URL.Query().Get("per_page") != "1" {
		t.Errorf("TestListAccountSubscriptions Error: Expected per_page parameter of 1, given %s", r.Request.URL.Query().Get("per_page"))
	}

	activated, _ := time.Parse(datetimeFormat, "2011-05-27T07:00:00Z")
	cpStartedAt, _ := time.Parse(datetimeFormat, "2011-06-27T07:00:00Z")
	cpEndsAt, _ := time.Parse(datetimeFormat, "2010-07-27T07:00:00Z")
	for _, given := range subscriptions {
		expected := Subscription{
			XMLName: xml.Name{Local: "subscription"},
			Plan: nestedPlan{
				Code: "gold",
				Name: "Gold plan",
			},
			Account: href{
				HREF: "https://your-subdomain.recurly.com/v2/accounts/1",
				Code: "1",
			},
			Invoice: href{
				HREF: "https://your-subdomain.recurly.com/v2/invoices/1108",
				Code: "1108",
			},
			UUID:                   "44f83d7cba354d5b84812419f923ea96",
			State:                  "active",
			UnitAmountInCents:      800,
			Currency:               "EUR",
			Quantity:               1,
			ActivatedAt:            NewTime(activated),
			CurrentPeriodStartedAt: NewTime(cpStartedAt),
			CurrentPeriodEndsAt:    NewTime(cpEndsAt),
			TaxInCents:             72,
			TaxType:                "usst",
			TaxRegion:              "CA",
			TaxRate:                0.0875,
			NetTerms:               NewInt(0),
		}

		if !reflect.DeepEqual(expected, given) {
			t.Errorf("TestListAccountSubscriptions Error: expected subscription to equal %#v, given %#v", expected, given)
		}
	}
}

func TestGetSubscriptions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/subscriptions/44f83d7cba354d5b84812419f923ea96", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("TestGetSubscriptions Error: Expected %s request, given %s", "GET", r.Method)
		}
		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?>
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
		</subscription>`)
	})

	r, subscription, err := client.Subscriptions.Get("44f83d7cba354d5b84812419f923ea96")
	if err != nil {
		t.Errorf("TestGetSubscriptions Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestGetSubscriptions Error: Expected list subcriptions to return OK")
	}

	activated, _ := time.Parse(datetimeFormat, "2011-05-27T07:00:00Z")
	cpStartedAt, _ := time.Parse(datetimeFormat, "2011-06-27T07:00:00Z")
	cpEndsAt, _ := time.Parse(datetimeFormat, "2010-07-27T07:00:00Z")
	expected := Subscription{
		XMLName: xml.Name{Local: "subscription"},
		Plan: nestedPlan{
			Code: "gold",
			Name: "Gold plan",
		},
		Account: href{
			HREF: "https://your-subdomain.recurly.com/v2/accounts/1",
			Code: "1",
		},
		Invoice: href{
			HREF: "https://your-subdomain.recurly.com/v2/invoices/1108",
			Code: "1108",
		},
		UUID:                   "44f83d7cba354d5b84812419f923ea96",
		State:                  "active",
		UnitAmountInCents:      800,
		Currency:               "EUR",
		Quantity:               1,
		ActivatedAt:            NewTime(activated),
		CurrentPeriodStartedAt: NewTime(cpStartedAt),
		CurrentPeriodEndsAt:    NewTime(cpEndsAt),
		TaxInCents:             72,
		TaxType:                "usst",
		TaxRegion:              "CA",
		TaxRate:                0.0875,
		NetTerms:               NewInt(0),
	}

	if !reflect.DeepEqual(expected, subscription) {
		t.Errorf("TestGetSubscriptions Error: expected subscription to equal %#v, given %#v", expected, subscription)
	}
}

func TestCreateSubscription(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/subscriptions", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("TestCreateSubscription Error: Expected %s request, given %s", "POST", r.Method)
		}
		rw.WriteHeader(201)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?><subscription></subscription>`)
	})

	r, _, err := client.Subscriptions.Create(NewSubscription{})
	if err != nil {
		t.Errorf("TestCreateSubscription Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestCreateSubscription Error: Expected create subscription to return OK")
	}
}

func TestPreviewSubscription(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/subscriptions/preview", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("TestPreviewSubscription Error: Expected %s request, given %s", "POST", r.Method)
		}
		rw.WriteHeader(201)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?><subscription></subscription>`)
	})

	r, _, err := client.Subscriptions.Preview(NewSubscription{})
	if err != nil {
		t.Errorf("TestPreviewSubscription Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestPreviewSubscription Error: Expected preview subscription to return OK")
	}
}

func TestUpdateSubscription(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/subscriptions/44f83d7cba354d5b84812419f923ea96", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("TestUpdateSubscription Error: Expected %s request, given %s", "PUT", r.Method)
		}
		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?><subscription></subscription>`)
	})

	r, _, err := client.Subscriptions.Update("44f83d7cba354d5b84812419f923ea96", UpdateSubscription{})
	if err != nil {
		t.Errorf("TestUpdateSubscription Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestUpdateSubscription Error: Expected update subscription to return OK")
	}
}

func TestUpdateSubscriptionNotes(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/notes", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("TestUpdateSubscriptionNotes Error: Expected %s request, given %s", "PUT", r.Method)
		}
		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?><subscription></subscription>`)
	})

	r, _, err := client.Subscriptions.UpdateNotes("44f83d7cba354d5b84812419f923ea96", SubscriptionNotes{})
	if err != nil {
		t.Errorf("TestUpdateSubscriptionNotes Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestUpdateSubscriptionNotes Error: Expected update subscription notes to return OK")
	}
}

func TestPreviewSubscriptionChange(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/preview", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("TestPreviewSubscriptionChange Error: Expected %s request, given %s", "POST", r.Method)
		}
		rw.WriteHeader(201)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?><subscription></subscription>`)
	})

	r, _, err := client.Subscriptions.PreviewChange("44f83d7cba354d5b84812419f923ea96", UpdateSubscription{})
	if err != nil {
		t.Errorf("TestPreviewSubscriptionChange Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestPreviewSubscriptionChange Error: Expected preview subscription change to return OK")
	}
}

func TestCancelSubscription(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/cancel", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("TestCancelSubscription Error: Expected %s request, given %s", "PUT", r.Method)
		}
		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?><subscription></subscription>`)
	})

	r, _, err := client.Subscriptions.Cancel("44f83d7cba354d5b84812419f923ea96")
	if err != nil {
		t.Errorf("TestCancelSubscription Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestCancelSubscription Error: Expected cancel subscription change to return OK")
	}
}

func TestReactivateSubscription(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/reactivate", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("TestReactivateSubscription Error: Expected %s request, given %s", "PUT", r.Method)
		}
		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?><subscription></subscription>`)
	})

	r, _, err := client.Subscriptions.Reactivate("44f83d7cba354d5b84812419f923ea96")
	if err != nil {
		t.Errorf("TestReactivateSubscription Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestReactivateSubscription Error: Expected reactivate subscription change to return OK")
	}
}

func TestTerminateSubscriptionWithPartialRefund(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/terminate", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("TestTerminateSubscriptionWithPartialRefund Error: Expected %s request, given %s", "PUT", r.Method)
		}
		if r.URL.Query().Get("refund_type") != "partial" {
			t.Errorf("TestTerminateSubscriptionWithPartialRefund Error: Expected refund_type of partial, given %s", r.URL.Query().Get("refund_type"))
		}
		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?><subscription></subscription>`)
	})

	r, _, err := client.Subscriptions.TerminateWithPartialRefund("44f83d7cba354d5b84812419f923ea96")
	if err != nil {
		t.Errorf("TestTerminateSubscriptionWithPartialRefund Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestTerminateSubscriptionWithPartialRefund Error: Expected terminate subscription with partial refund to return OK")
	}
}

func TestTerminateSubscriptionWithFullRefund(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/terminate", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("TestTerminateSubscriptionWithFullRefund Error: Expected %s request, given %s", "PUT", r.Method)
		}
		if r.URL.Query().Get("refund_type") != "full" {
			t.Errorf("TestTerminateSubscriptionWithFullRefund Error: Expected refund_type of full, given %s", r.URL.Query().Get("refund_type"))
		}
		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?><subscription></subscription>`)
	})

	r, _, err := client.Subscriptions.TerminateWithFullRefund("44f83d7cba354d5b84812419f923ea96")
	if err != nil {
		t.Errorf("TestTerminateSubscriptionWithFullRefund Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestTerminateSubscriptionWithFullRefund Error: Expected terminate subscription with full refund to return OK")
	}
}

func TestTerminateSubscriptionWithoutRefund(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/terminate", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("TestTerminateSubscriptionWithoutRefund Error: Expected %s request, given %s", "PUT", r.Method)
		}
		if r.URL.Query().Get("refund_type") != "none" {
			t.Errorf("TestTerminateSubscriptionWithoutRefund Error: Expected refund_type of none, given %s", r.URL.Query().Get("refund_type"))
		}
		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?><subscription></subscription>`)
	})

	r, _, err := client.Subscriptions.TerminateWithoutRefund("44f83d7cba354d5b84812419f923ea96")
	if err != nil {
		t.Errorf("TestTerminateSubscriptionWithoutRefund Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestTerminateSubscriptionWithoutRefund Error: Expected terminate subscription without refund to return OK")
	}
}

func TestPostponeSubscription(t *testing.T) {
	setup()
	defer teardown()

	ts, _ := time.Parse(datetimeFormat, "2015-08-27T07:00:00Z")
	mux.HandleFunc("/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/postpone", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("TestPostponeSubscription Error: Expected %s request, given %s", "PUT", r.Method)
		}
		if r.URL.Query().Get("next_renewal_date") != "2015-08-27T07:00:00Z" {
			t.Errorf("TestPostponeSubscription Error: Expected qs param of next_renewal date equal to 2015-08-27T07:00:00Z, given %s", r.URL.Query().Get("next_renewal_date"))
		}
		if r.URL.Query().Get("bulk") != "false" {
			t.Errorf("TestPostponeSubscription Error: Expected qs param of bulk equal to false, given %s", r.URL.Query().Get("bulk"))
		}
		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?><subscription></subscription>`)
	})

	r, _, err := client.Subscriptions.Postpone("44f83d7cba354d5b84812419f923ea96", ts, false)
	if err != nil {
		t.Errorf("TestPostponeSubscription Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestPostponeSubscription Error: Expected postpone subscription change to return OK")
	}
}
