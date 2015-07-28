package recurly

import (
	"bytes"
	"encoding/xml"
	"testing"
	"time"
)

// TestSubscriptionsEncoding ensures structs are encoded to XML properly.
// Because Recurly supports partial updates, it's important that only defined
// fields are handled properly -- including types like booleans and integers which
// have zero values that we want to send.
func TestSubscriptionsEncoding(t *testing.T) {
	ts, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2015-06-03T13:42:23.764061Z")
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
