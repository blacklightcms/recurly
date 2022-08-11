package webhooks_test

import (
	"encoding/xml"
	"os"
	"testing"
	"time"

	"github.com/autopilot3/recurly"
	"github.com/autopilot3/recurly/webhooks"
	"github.com/google/go-cmp/cmp"
)

func TestParse_BillingInfoUpdatedNotification(t *testing.T) {
	result := MustParseFile("testdata/billing_info_updated_notification.xml")
	if n, ok := result.(*webhooks.AccountNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.AccountNotification{
		Type: webhooks.BillingInfoUpdated,
		Account: webhooks.Account{
			XMLName:   xml.Name{Local: "account"},
			Code:      "1",
			Email:     "verena@example.com",
			FirstName: "Verena",
			LastName:  "Example",
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_BillingInfoUpdateFailedNotification(t *testing.T) {
	result := MustParseFile("testdata/billing_info_update_failed_notification.xml")
	if n, ok := result.(*webhooks.AccountNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.AccountNotification{
		Type: webhooks.BillingInfoUpdateFailed,
		Account: webhooks.Account{
			XMLName:   xml.Name{Local: "account"},
			Code:      "1",
			Email:     "verena@example.com",
			FirstName: "Verena",
			LastName:  "Example",
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_NewAccountNotification(t *testing.T) {
	result := MustParseFile("testdata/new_account_notification.xml")
	if n, ok := result.(*webhooks.AccountNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.AccountNotification{
		Type: webhooks.NewAccount,
		Account: webhooks.Account{
			XMLName:   xml.Name{Local: "account"},
			Code:      "1",
			Email:     "verena@example.com",
			FirstName: "Verena",
			LastName:  "Example",
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_UpdatedAccountNotification(t *testing.T) {
	result := MustParseFile("testdata/updated_account_notification.xml")
	if n, ok := result.(*webhooks.AccountNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.AccountNotification{
		Type: webhooks.UpdatedAccount,
		Account: webhooks.Account{
			XMLName:   xml.Name{Local: "account"},
			Code:      "1",
			Email:     "verena@example.com",
			FirstName: "Verena",
			LastName:  "Example",
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_CanceledAccountNotification(t *testing.T) {
	result := MustParseFile("testdata/canceled_account_notification.xml")
	if n, ok := result.(*webhooks.AccountNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.AccountNotification{
		Type: webhooks.CanceledAccount,
		Account: webhooks.Account{
			XMLName:   xml.Name{Local: "account"},
			Code:      "1",
			Email:     "verena@example.com",
			FirstName: "Verena",
			LastName:  "Example",
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_ChargeInvoiceNotification(t *testing.T) {
	result := MustParseFile("testdata/charge_invoice_notification.xml")
	if n, ok := result.(*webhooks.ChargeInvoiceNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.ChargeInvoiceNotification{
		Type: webhooks.NewChargeInvoice,
		Account: webhooks.Account{
			XMLName: xml.Name{Local: "account"},
			Code:    "1234",
		},
		Invoice: webhooks.ChargeInvoice{
			XMLName:                       xml.Name{Local: "invoice"},
			UUID:                          "42feb03ce368c0e1ead35d4bfa89b82e",
			State:                         recurly.ChargeInvoiceStatePending,
			Origin:                        recurly.ChargeInvoiceOriginRenewal,
			SubscriptionUUIDs:             []string{"40b8f5e99df03b8684b99d4993b6e089"},
			InvoiceNumber:                 2405,
			Currency:                      "USD",
			BalanceInCents:                100000,
			TotalInCents:                  100000,
			SubtotalInCents:               100000,
			SubTotalBeforeDiscountInCents: 100000,
			NetTerms:                      recurly.NewInt(30),
			CollectionMethod:              recurly.CollectionMethodManual,
			CreatedAt:                     recurly.NewTime(MustParseTime("2018-02-13T16:00:04Z")),
			UpdatedAt:                     recurly.NewTime(MustParseTime("2018-02-13T16:00:04Z")),
			DueOn:                         recurly.NewTime(MustParseTime("2018-03-16T15:00:04Z")),
			CustomerNotes:                 "Thanks for your business!",
			TermsAndConditions:            "Payment can be made out to Acme, Co.",
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_CreditInvoiceNotification(t *testing.T) {
	result := MustParseFile("testdata/credit_invoice_notification.xml")
	if n, ok := result.(*webhooks.CreditInvoiceNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.CreditInvoiceNotification{
		Type: webhooks.NewCreditInvoice,
		Account: webhooks.Account{
			XMLName: xml.Name{Local: "account"},
			Code:    "1234",
		},
		Invoice: webhooks.CreditInvoice{
			XMLName:                       xml.Name{Local: "invoice"},
			UUID:                          "42fb74de65e9395eb004614144a7b91f",
			State:                         recurly.CreditInvoiceStateClosed,
			Origin:                        recurly.CreditInvoiceOriginWriteOff,
			SubscriptionUUIDs:             []string{"42fb74ba9efe4c6981c2064436a4e9cd"},
			InvoiceNumber:                 2404,
			Currency:                      "USD",
			BalanceInCents:                0,
			TotalInCents:                  -4882,
			TaxInCents:                    -382,
			SubtotalInCents:               -4500,
			SubTotalBeforeDiscountInCents: -5000,
			DiscountInCents:               -500,
			CreatedAt:                     recurly.NewTime(MustParseTime("2018-02-13T00:56:22Z")),
			UpdatedAt:                     recurly.NewTime(MustParseTime("2018-02-13T00:56:22Z")),
			ClosedAt:                      recurly.NewTime(MustParseTime("2018-02-13T00:56:22Z")),
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_CreditPaymentNotification(t *testing.T) {
	result := MustParseFile("testdata/credit_payment_notification.xml")
	if n, ok := result.(*webhooks.CreditPaymentNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.CreditPaymentNotification{
		Type: webhooks.NewCreditPayment,
		Account: webhooks.Account{
			XMLName: xml.Name{Local: "account"},
			Code:    "1234",
		},
		CreditPayment: webhooks.CreditPayment{
			XMLName:                xml.Name{Local: "credit_payment"},
			UUID:                   "42fa2a56dfeca2ace39b0e4a9198f835",
			Action:                 "payment",
			AmountInCents:          3579,
			OriginalInvoiceNumber:  2389,
			AppliedToInvoiceNumber: 2390,
			CreatedAt:              recurly.NewTime(MustParseTime("2018-02-12T18:55:20Z")),
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_NewSubscriptionNotification(t *testing.T) {
	result := MustParseFile("testdata/new_subscription_notification.xml")
	if n, ok := result.(*webhooks.SubscriptionNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.SubscriptionNotification{
		Type: webhooks.NewSubscription,
		Account: webhooks.Account{
			XMLName:   xml.Name{Local: "account"},
			Code:      "1",
			Email:     "verena@example.com",
			FirstName: "Verena",
			LastName:  "Example",
		},
		Subscription: recurly.Subscription{
			XMLName: xml.Name{Local: "subscription"},
			Plan: recurly.NestedPlan{
				Code: "bronze",
				Name: "Bronze Plan",
			},
			UUID:                   "d1b6d359a01ded71caed78eaa0fedf8e",
			State:                  "active",
			Quantity:               2,
			TotalAmountInCents:     17000,
			ActivatedAt:            recurly.NewTime(MustParseTime("2010-09-23T22:05:03Z")),
			CanceledAt:             recurly.NewTime(MustParseTime("2010-09-23T22:05:43Z")),
			ExpiresAt:              recurly.NewTime(MustParseTime("2010-09-24T22:05:03Z")),
			CurrentPeriodStartedAt: recurly.NewTime(MustParseTime("2010-09-23T22:05:03Z")),
			CurrentPeriodEndsAt:    recurly.NewTime(MustParseTime("2010-09-24T22:05:03Z")),
			CollectionMethod:       recurly.CollectionMethodAutomatic,
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_UpdatedSubscriptionNotification(t *testing.T) {
	result := MustParseFile("testdata/updated_subscription_notification.xml")
	if n, ok := result.(*webhooks.SubscriptionNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.SubscriptionNotification{
		Type: webhooks.UpdatedSubscription,
		Account: webhooks.Account{
			XMLName:   xml.Name{Local: "account"},
			Code:      "1",
			Email:     "verena@example.com",
			FirstName: "Verena",
			LastName:  "Example",
		},
		Subscription: recurly.Subscription{
			XMLName: xml.Name{Local: "subscription"},
			Plan: recurly.NestedPlan{
				Code: "1dpt",
				Name: "Subscription One",
			},
			UUID:                   "292332928954ca62fa48048be5ac98ec",
			State:                  "active",
			Quantity:               1,
			TotalAmountInCents:     200,
			ActivatedAt:            recurly.NewTime(MustParseTime("2010-09-23T22:05:03Z")),
			CanceledAt:             recurly.NewTime(MustParseTime("2010-09-23T22:05:43Z")),
			ExpiresAt:              recurly.NewTime(MustParseTime("2010-09-24T22:05:03Z")),
			CurrentPeriodStartedAt: recurly.NewTime(MustParseTime("2010-09-23T22:05:03Z")),
			CurrentPeriodEndsAt:    recurly.NewTime(MustParseTime("2010-09-24T22:05:03Z")),
			CollectionMethod:       recurly.CollectionMethodAutomatic,
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_RenewedSubscriptionNotification(t *testing.T) {
	result := MustParseFile("testdata/renewed_subscription_notification.xml")
	if n, ok := result.(*webhooks.SubscriptionNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.SubscriptionNotification{
		Type: webhooks.RenewedSubscription,
		Account: webhooks.Account{
			XMLName:     xml.Name{Local: "account"},
			Code:        "1",
			Email:       "verena@example.com",
			FirstName:   "Verena",
			LastName:    "Example",
			CompanyName: "Company, Inc.",
		},
		Subscription: recurly.Subscription{
			XMLName: xml.Name{Local: "subscription"},
			Plan: recurly.NestedPlan{
				Code: "bootstrap",
				Name: "Bootstrap",
			},
			UUID:                   "6ab458a887d38070807ebb3bed7ac1e5",
			State:                  "active",
			Quantity:               1,
			TotalAmountInCents:     9900,
			ActivatedAt:            recurly.NewTime(MustParseTime("2010-07-22T20:42:05Z")),
			CurrentPeriodStartedAt: recurly.NewTime(MustParseTime("2010-09-22T20:42:05Z")),
			CurrentPeriodEndsAt:    recurly.NewTime(MustParseTime("2010-10-22T20:42:05Z")),
			CollectionMethod:       recurly.CollectionMethodAutomatic,
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_ExpiredSubscriptionNotification(t *testing.T) {
	result := MustParseFile("testdata/expired_subscription_notification.xml")
	if n, ok := result.(*webhooks.SubscriptionNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.SubscriptionNotification{
		Type: webhooks.ExpiredSubscription,
		Account: webhooks.Account{
			XMLName:   xml.Name{Local: "account"},
			Code:      "1",
			Email:     "verena@example.com",
			FirstName: "Verena",
			LastName:  "Example",
		},
		Subscription: recurly.Subscription{
			XMLName: xml.Name{Local: "subscription"},
			Plan: recurly.NestedPlan{
				Code: "1dpt",
				Name: "Subscription One",
			},
			UUID:                   "d1b6d359a01ded71caed78eaa0fedf8e",
			State:                  "expired",
			Quantity:               1,
			TotalAmountInCents:     200,
			ActivatedAt:            recurly.NewTime(MustParseTime("2010-09-23T22:05:03Z")),
			CanceledAt:             recurly.NewTime(MustParseTime("2010-09-23T22:05:43Z")),
			ExpiresAt:              recurly.NewTime(MustParseTime("2010-09-24T22:05:03Z")),
			CurrentPeriodStartedAt: recurly.NewTime(MustParseTime("2010-09-23T22:05:03Z")),
			CurrentPeriodEndsAt:    recurly.NewTime(MustParseTime("2010-09-24T22:05:03Z")),
			CollectionMethod:       recurly.CollectionMethodAutomatic,
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_CanceledSubscriptionNotification(t *testing.T) {
	result := MustParseFile("testdata/canceled_subscription_notification.xml")
	if n, ok := result.(*webhooks.SubscriptionNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.SubscriptionNotification{
		Type: webhooks.CanceledSubscription,
		Account: webhooks.Account{
			XMLName:   xml.Name{Local: "account"},
			Code:      "1",
			Email:     "verena@example.com",
			FirstName: "Verena",
			LastName:  "Example",
		},
		Subscription: recurly.Subscription{
			XMLName: xml.Name{Local: "subscription"},
			Plan: recurly.NestedPlan{
				Code: "1dpt",
				Name: "Subscription One",
			},
			UUID:                   "dccd742f4710e78515714d275839f891",
			State:                  "canceled",
			Quantity:               1,
			TotalAmountInCents:     200,
			ActivatedAt:            recurly.NewTime(MustParseTime("2010-09-23T22:05:03Z")),
			CanceledAt:             recurly.NewTime(MustParseTime("2010-09-23T22:05:43Z")),
			ExpiresAt:              recurly.NewTime(MustParseTime("2010-09-24T22:05:03Z")),
			CurrentPeriodStartedAt: recurly.NewTime(MustParseTime("2010-09-23T22:05:03Z")),
			CurrentPeriodEndsAt:    recurly.NewTime(MustParseTime("2010-09-24T22:05:03Z")),
			CollectionMethod:       recurly.CollectionMethodAutomatic,
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_PausedSubscriptionNotification(t *testing.T) {
	result := MustParseFile("testdata/subscription_paused_notification.xml")
	if n, ok := result.(*webhooks.SubscriptionNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.SubscriptionNotification{
		Type: webhooks.PausedSubscription,
		Account: webhooks.Account{
			XMLName:   xml.Name{Local: "account"},
			Code:      "1",
			Email:     "verena@example.com",
			FirstName: "Verena",
			LastName:  "Example",
		},
		Subscription: recurly.Subscription{
			XMLName: xml.Name{Local: "subscription"},
			Plan: recurly.NestedPlan{
				Code: "daily_plan",
				Name: "daily_plan",
			},
			UUID:                   "437a818b9dba81065e444448de931842",
			State:                  "paused",
			Quantity:               10,
			TotalAmountInCents:     10000,
			ActivatedAt:            recurly.NewTime(MustParseTime("2018-03-09T17:01:59Z")),
			CurrentPeriodStartedAt: recurly.NewTime(MustParseTime("2018-03-10T22:12:08Z")),
			CurrentPeriodEndsAt:    recurly.NewTime(MustParseTime("2018-03-11T22:12:08Z")),
			PausedAt:               recurly.NewTime(MustParseTime("2018-03-10T22:12:08Z")),
			ResumeAt:               recurly.NewTime(MustParseTime("2018-03-20T22:12:08Z")),
			RemainingPauseCycles:   9,
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_ResumedSubscriptionNotification(t *testing.T) {
	result := MustParseFile("testdata/subscription_resumed_notification.xml")
	if n, ok := result.(*webhooks.SubscriptionNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.SubscriptionNotification{
		Type: webhooks.ResumedSubscription,
		Account: webhooks.Account{
			XMLName:   xml.Name{Local: "account"},
			Code:      "1",
			Email:     "verena@example.com",
			FirstName: "Verena",
			LastName:  "Example",
		},
		Subscription: recurly.Subscription{
			XMLName: xml.Name{Local: "subscription"},
			Plan: recurly.NestedPlan{
				Code: "daily_plan",
				Name: "daily_plan",
			},
			UUID:                   "437a818b9dba81065e444448de931842",
			State:                  "active",
			Quantity:               10,
			TotalAmountInCents:     10000,
			ActivatedAt:            recurly.NewTime(MustParseTime("2018-03-09T17:01:59Z")),
			CurrentPeriodStartedAt: recurly.NewTime(MustParseTime("2018-03-20T17:50:27Z")),
			CurrentPeriodEndsAt:    recurly.NewTime(MustParseTime("2018-03-21T17:50:27Z")),
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_ScheduledSubscriptionPausedNotification(t *testing.T) {
	result := MustParseFile("testdata/scheduled_subscription_pause_notification.xml")
	if n, ok := result.(*webhooks.SubscriptionNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.SubscriptionNotification{
		Type: webhooks.ScheduledPauseSubscription,
		Account: webhooks.Account{
			XMLName:   xml.Name{Local: "account"},
			Code:      "1",
			Email:     "verena@example.com",
			FirstName: "Verena",
			LastName:  "Example",
		},
		Subscription: recurly.Subscription{
			XMLName: xml.Name{Local: "subscription"},
			Plan: recurly.NestedPlan{
				Code: "daily_plan",
				Name: "daily_plan",
			},
			UUID:                   "437b9def1c442e659f90f4416086dd66",
			State:                  "active",
			Quantity:               1,
			TotalAmountInCents:     709,
			ActivatedAt:            recurly.NewTime(MustParseTime("2018-03-09T22:12:36Z")),
			CurrentPeriodStartedAt: recurly.NewTime(MustParseTime("2018-03-09T22:12:36Z")),
			CurrentPeriodEndsAt:    recurly.NewTime(MustParseTime("2019-03-09T22:12:36Z")),
			PausedAt:               recurly.NewTime(MustParseTime("2019-03-09T22:12:36Z")),
			ResumeAt:               recurly.NewTime(MustParseTime("2024-03-09T22:12:36Z")),
			RemainingPauseCycles:   5,
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_SubscriptionPauseModifiedNotification(t *testing.T) {
	result := MustParseFile("testdata/subscription_pause_modified_notification.xml")
	if n, ok := result.(*webhooks.SubscriptionNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.SubscriptionNotification{
		Type: webhooks.ModifiedPauseSubscription,
		Account: webhooks.Account{
			XMLName:   xml.Name{Local: "account"},
			Code:      "1",
			Email:     "verena@example.com",
			FirstName: "Verena",
			LastName:  "Example",
		},
		Subscription: recurly.Subscription{
			XMLName: xml.Name{Local: "subscription"},
			Plan: recurly.NestedPlan{
				Code: "daily_plan",
				Name: "daily_plan",
			},
			UUID:                   "437a818b9dba81065e444448de931842",
			State:                  "active",
			Quantity:               1,
			TotalAmountInCents:     709,
			ActivatedAt:            recurly.NewTime(MustParseTime("2018-03-09T17:01:59Z")),
			CurrentPeriodStartedAt: recurly.NewTime(MustParseTime("2018-03-09T13:33:09Z")),
			CurrentPeriodEndsAt:    recurly.NewTime(MustParseTime("2018-03-09T13:38:22Z")),
			PausedAt:               recurly.NewTime(MustParseTime("2018-03-09T13:38:22Z")),
			ResumeAt:               recurly.NewTime(MustParseTime("2023-03-09T13:38:22Z")),
			RemainingPauseCycles:   5,
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_PausedSubscriptionRenewalNotification(t *testing.T) {
	result := MustParseFile("testdata/paused_subscription_renewal_notification.xml")
	if n, ok := result.(*webhooks.SubscriptionNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.SubscriptionNotification{
		Type: webhooks.PausedRenewalSubscription,
		Account: webhooks.Account{
			XMLName:   xml.Name{Local: "account"},
			Code:      "1",
			Email:     "verena@example.com",
			FirstName: "Verena",
			LastName:  "Example",
		},
		Subscription: recurly.Subscription{
			XMLName: xml.Name{Local: "subscription"},
			Plan: recurly.NestedPlan{
				Code: "daily_plan",
				Name: "daily_plan",
			},
			UUID:                   "437a818b9dba81065e444448de931842",
			State:                  "paused",
			Quantity:               10,
			TotalAmountInCents:     10000,
			ActivatedAt:            recurly.NewTime(MustParseTime("2018-03-09T17:01:59Z")),
			CurrentPeriodStartedAt: recurly.NewTime(MustParseTime("2018-03-18T17:50:27Z")),
			CurrentPeriodEndsAt:    recurly.NewTime(MustParseTime("2018-03-19T17:50:27Z")),
			PausedAt:               recurly.NewTime(MustParseTime("2018-03-10T22:12:08Z")),
			ResumeAt:               recurly.NewTime(MustParseTime("2018-03-20T17:50:27Z")),
			RemainingPauseCycles:   1,
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}
func TestParse_SubscriptionPauseCanceledNotification(t *testing.T) {
	result := MustParseFile("testdata/subscription_pause_canceled_notification.xml")
	if n, ok := result.(*webhooks.SubscriptionNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.SubscriptionNotification{
		Type: webhooks.PauseCanceledSubscription,
		Account: webhooks.Account{
			XMLName:   xml.Name{Local: "account"},
			Code:      "1",
			Email:     "verena@example.com",
			FirstName: "Verena",
			LastName:  "Example",
		},
		Subscription: recurly.Subscription{
			XMLName: xml.Name{Local: "subscription"},
			Plan: recurly.NestedPlan{
				Code: "daily_plan",
				Name: "daily_plan",
			},
			UUID:                   "437b9def1c442e659f90f4416086dd66",
			State:                  "active",
			Quantity:               1,
			TotalAmountInCents:     2000,
			ActivatedAt:            recurly.NewTime(MustParseTime("2018-03-09T22:12:36Z")),
			CurrentPeriodStartedAt: recurly.NewTime(MustParseTime("2018-03-09T22:12:36Z")),
			CurrentPeriodEndsAt:    recurly.NewTime(MustParseTime("2019-03-09T22:12:36Z")),
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}
func TestParse_ReactivatedAccountNotification(t *testing.T) {
	result := MustParseFile("testdata/reactivated_account_notification.xml")
	if n, ok := result.(*webhooks.SubscriptionNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.SubscriptionNotification{
		Type: webhooks.ReactivatedAccount,
		Account: webhooks.Account{
			XMLName:   xml.Name{Local: "account"},
			Code:      "1",
			Email:     "verena@example.com",
			FirstName: "Verena",
			LastName:  "Example",
		},
		Subscription: recurly.Subscription{
			XMLName: xml.Name{Local: "subscription"},
			Plan: recurly.NestedPlan{
				Code: "bootstrap",
				Name: "Bootstrap",
			},
			UUID:                   "6ab458a887d38070807ebb3bed7ac1e5",
			State:                  "active",
			Quantity:               1,
			TotalAmountInCents:     9900,
			ActivatedAt:            recurly.NewTime(MustParseTime("2010-07-22T20:42:05Z")),
			CurrentPeriodStartedAt: recurly.NewTime(MustParseTime("2010-09-22T20:42:05Z")),
			CurrentPeriodEndsAt:    recurly.NewTime(MustParseTime("2010-10-22T20:42:05Z")),
			CollectionMethod:       recurly.CollectionMethodAutomatic,
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_SuccessfulPaymentNotification(t *testing.T) {
	result := MustParseFile("testdata/successful_payment_notification.xml")
	if n, ok := result.(*webhooks.PaymentNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.PaymentNotification{
		Type: webhooks.SuccessfulPayment,
		Account: webhooks.Account{
			XMLName:     xml.Name{Local: "account"},
			Code:        "1",
			Username:    "verena",
			Email:       "verena@example.com",
			FirstName:   "Verena",
			LastName:    "Example",
			CompanyName: "Company, Inc.",
		},
		Transaction: webhooks.Transaction{
			XMLName:       xml.Name{Local: "transaction"},
			UUID:          "a5143c1d3a6f4a8287d0e2cc1d4c0427",
			InvoiceNumber: 2059,
			Action:        "purchase",
			PaymentMethod: "credit_card",
			AmountInCents: 1000,
			Status:        "success",
			Message:       "Bogus Gateway: Forced success",
			Reference:     "reference",
			Source:        "subscription",
			Test:          recurly.NewBool(true),
			Voidable:      recurly.NewBool(true),
			Refundable:    recurly.NewBool(true),
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_FailedPaymentNotification(t *testing.T) {
	result := MustParseFile("testdata/failed_payment_notification.xml")
	if n, ok := result.(*webhooks.PaymentNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.PaymentNotification{
		Type: webhooks.FailedPayment,
		Account: webhooks.Account{
			XMLName:     xml.Name{Local: "account"},
			Code:        "1",
			Username:    "verena",
			Email:       "verena@example.com",
			FirstName:   "Verena",
			LastName:    "Example",
			CompanyName: "Company, Inc.",
		},
		Transaction: webhooks.Transaction{
			XMLName:          xml.Name{Local: "transaction"},
			UUID:             "a5143c1d3a6f4a8287d0e2cc1d4c0427",
			InvoiceNumber:    2059,
			SubscriptionUUID: "1974a098jhlkjasdfljkha898326881c",
			Action:           "purchase",
			PaymentMethod:    "credit_card",
			AmountInCents:    1000,
			Status:           "Declined",
			Message:          "This transaction has been declined",
			FailureType:      "Declined by the gateway",
			Reference:        "reference",
			Source:           "subscription",
			Test:             recurly.NewBool(true),
			Voidable:         recurly.NewBool(false),
			Refundable:       recurly.NewBool(false),
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_VoidPaymentNotification(t *testing.T) {
	result := MustParseFile("testdata/void_payment_notification.xml")
	if n, ok := result.(*webhooks.PaymentNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.PaymentNotification{
		Type: webhooks.VoidPayment,
		Account: webhooks.Account{
			XMLName:     xml.Name{Local: "account"},
			Code:        "1",
			Username:    "verena",
			Email:       "verena@example.com",
			FirstName:   "Verena",
			LastName:    "Example",
			CompanyName: "Company, Inc.",
		},
		Transaction: webhooks.Transaction{
			XMLName:          xml.Name{Local: "transaction"},
			UUID:             "a5143c1d3a6f4a8287d0e2cc1d4c0427",
			InvoiceNumber:    2059,
			SubscriptionUUID: "1974a098jhlkjasdfljkha898326881c",
			Action:           "purchase",
			AmountInCents:    1000,
			Status:           "void",
			Message:          "Test Gateway: Successful test transaction",
			Reference:        "reference",
			Source:           "subscription",
			Test:             recurly.NewBool(true),
			Voidable:         recurly.NewBool(true),
			Refundable:       recurly.NewBool(true),
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_SuccessfulRefundNotification(t *testing.T) {
	result := MustParseFile("testdata/successful_refund_notification.xml")
	if n, ok := result.(*webhooks.PaymentNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.PaymentNotification{
		Type: webhooks.SuccessfulRefund,
		Account: webhooks.Account{
			XMLName:     xml.Name{Local: "account"},
			Code:        "1",
			Username:    "verena",
			Email:       "verena@example.com",
			FirstName:   "Verena",
			LastName:    "Example",
			CompanyName: "Company, Inc.",
		},
		Transaction: webhooks.Transaction{
			XMLName:          xml.Name{Local: "transaction"},
			UUID:             "a5143c1d3a6f4a8287d0e2cc1d4c0427",
			InvoiceNumber:    2059,
			SubscriptionUUID: "1974a098jhlkjasdfljkha898326881c",
			Action:           "credit",
			AmountInCents:    1000,
			Status:           "success",
			Message:          "Bogus Gateway: Forced success",
			Reference:        "reference",
			Source:           "subscription",
			Test:             recurly.NewBool(true),
			Voidable:         recurly.NewBool(true),
			Refundable:       recurly.NewBool(true),
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_ProcessingPaymentNotification(t *testing.T) {
	result := MustParseFile("testdata/processing_payment_notification.xml")
	if n, ok := result.(*webhooks.PaymentNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.PaymentNotification{
		Type: webhooks.ProcessingPayment,
		Account: webhooks.Account{
			XMLName:     xml.Name{Local: "account"},
			Code:        "1",
			Username:    "verena",
			Email:       "verena@example.com",
			FirstName:   "Verena",
			LastName:    "Example",
			CompanyName: "Company, Inc.",
		},
		Transaction: webhooks.Transaction{
			XMLName:          xml.Name{Local: "transaction"},
			UUID:             "a5143c1d3a6f4a8287d0e2cc1d4c0427",
			InvoiceNumber:    2059,
			SubscriptionUUID: "1974a098jhlkjasdfljkha898326881c",
			Action:           "purchase",
			AmountInCents:    1000,
			Status:           "processing",
			Message:          "Bogus Gateway: Forced success",
			Reference:        "",
			Source:           "subscription",
			Test:             recurly.NewBool(true),
			Voidable:         recurly.NewBool(true),
			Refundable:       recurly.NewBool(true),
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_ScheduledPaymentNotification(t *testing.T) {
	result := MustParseFile("testdata/scheduled_payment_notification.xml")
	if n, ok := result.(*webhooks.PaymentNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.PaymentNotification{
		Type: webhooks.ScheduledPayment,
		Account: webhooks.Account{
			XMLName:     xml.Name{Local: "account"},
			Code:        "1",
			Username:    "verena",
			Email:       "verena@example.com",
			FirstName:   "Verena",
			LastName:    "Example",
			CompanyName: "Company, Inc.",
		},
		Transaction: webhooks.Transaction{
			XMLName:          xml.Name{Local: "transaction"},
			UUID:             "a5143c1d3a6f4a8287d0e2cc1d4c0427",
			InvoiceNumber:    2059,
			SubscriptionUUID: "1974a098jhlkjasdfljkha898326881c",
			Action:           "purchase",
			AmountInCents:    1000,
			Status:           "scheduled",
			Message:          "Bogus Gateway: Forced success",
			Reference:        "",
			Source:           "subscription",
			Test:             recurly.NewBool(true),
			Voidable:         recurly.NewBool(true),
			Refundable:       recurly.NewBool(true),
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_TransactionStatusUpdatedNotification(t *testing.T) {
	result := MustParseFile("testdata/transaction_status_updated_notification.xml")
	if n, ok := result.(*webhooks.PaymentNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.PaymentNotification{
		Type: webhooks.TransactionStatusUpdated,
		Account: webhooks.Account{
			XMLName:     xml.Name{Local: "account"},
			Code:        "1",
			Username:    "verena",
			Email:       "verena@example.com",
			FirstName:   "Verena",
			LastName:    "Example",
			CompanyName: "Company, Inc.",
		},
		Transaction: webhooks.Transaction{
			XMLName:          xml.Name{Local: "transaction"},
			UUID:             "a5143c1d3a6f4a8287d0e2cc1d4c0427",
			InvoiceNumber:    2059,
			SubscriptionUUID: "1974a098jhlkjasdfljkha898326881c",
			Action:           "purchase",
			AmountInCents:    1000,
			Status:           "void",
			Message:          "Test Gateway: Successful test transaction",
			Reference:        "reference",
			Source:           "subscription",
			Test:             recurly.NewBool(true),
			Voidable:         recurly.NewBool(true),
			Refundable:       recurly.NewBool(true),
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_NewDunningEventNotification(t *testing.T) {
	result := MustParseFile("testdata/new_dunning_event_notification.xml")
	if n, ok := result.(*webhooks.NewDunningEventNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.NewDunningEventNotification{
		Type: webhooks.NewDunningEvent,
		Account: webhooks.Account{
			XMLName:     xml.Name{Local: "account"},
			Code:        "1",
			Username:    "verena",
			Email:       "verena@example.com",
			FirstName:   "Verena",
			LastName:    "Example",
			CompanyName: "Company, Inc.",
		},
		Invoice: webhooks.ChargeInvoice{
			XMLName:                       xml.Name{Local: "invoice"},
			UUID:                          "424a9d4a2174b4f39bc776426aa19c32",
			SubscriptionUUIDs:             []string{"4110792b3b01967d854f674b7282f542"},
			State:                         "past_due",
			Origin:                        "renewal",
			Currency:                      "USD",
			CreatedAt:                     recurly.NewTime(MustParseTime("2018-01-09T16:47:43Z")),
			UpdatedAt:                     recurly.NewTime(MustParseTime("2018-02-12T16:50:23Z")),
			DueOn:                         recurly.NewTime(MustParseTime("2018-02-09T16:47:43Z")),
			BalanceInCents:                4000,
			InvoiceNumber:                 1813,
			TotalInCents:                  4500,
			SubtotalInCents:               4500,
			SubTotalBeforeDiscountInCents: 4500,
			CustomerNotes:                 "Thanks for your business!",
			NetTerms:                      recurly.NewInt(30),
			CollectionMethod:              recurly.CollectionMethodManual,
		},
		Subscription: recurly.Subscription{
			XMLName: xml.Name{Local: "subscription"},
			Plan: recurly.NestedPlan{
				Code: "gold",
				Name: "Gold",
			},
			UUID:                   "4110792b3b01967d854f674b7282f542",
			State:                  "active",
			Quantity:               1,
			TotalAmountInCents:     4500,
			ActivatedAt:            recurly.NewTime(MustParseTime("2017-11-09T16:47:30Z")),
			CurrentPeriodStartedAt: recurly.NewTime(MustParseTime("2018-02-09T16:47:30Z")),
			CurrentPeriodEndsAt:    recurly.NewTime(MustParseTime("2018-03-09T16:47:30Z")),
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_ErrUnknownNotification(t *testing.T) {
	f := MustOpenFile("testdata/unknown_notification.xml")
	defer f.Close()

	result, err := webhooks.Parse(f)
	if result != nil {
		t.Fatalf("unexpected notification: %#v", result)
	} else if e, ok := err.(webhooks.ErrUnknownNotification); !ok {
		t.Fatalf("unexpected type: %T, error", e)
	} else if err.Error() != `unknown notification: "unknown_notification"` {
		t.Fatalf("unexpected error string: %s", err.Error())
	} else if e.Name() != "unknown_notification" {
		t.Fatalf("unexpected notification name: %s", e.Name())
	}
}

func MustParseFile(file string) interface{} {
	f := MustOpenFile(file)
	result, err := webhooks.Parse(f)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	return result
}

func MustOpenFile(name string) *os.File {
	file, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	return file
}

// MustParseTime parses a string into time.Time, panicing if there is an error.
func MustParseTime(str string) time.Time {
	t, err := time.Parse(recurly.DateTimeFormat, str)
	if err != nil {
		panic(err)
	}
	return t
}
