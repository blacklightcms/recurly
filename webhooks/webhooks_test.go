package webhooks_test

import (
	"encoding/xml"
	"os"
	"testing"
	"time"

	"github.com/blacklightcms/recurly"
	"github.com/blacklightcms/recurly/webhooks"
	"github.com/google/go-cmp/cmp"
)

func TestParse_BillingInfoUpdatedNotification(t *testing.T) {
	xmlFile := MustOpenFile("testdata/billing_info_updated_notification.xml")
	result, err := webhooks.Parse(xmlFile)
	if err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.AccountNotification); !ok {
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

func TestParse_ChargeInvoiceNotification(t *testing.T) {
	created, _ := time.Parse(recurly.DateTimeFormat, "2018-02-13T16:00:04Z")
	xmlFile := MustOpenFile("testdata/charge_invoice_notification.xml")
	result, err := webhooks.Parse(xmlFile)
	if err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.ChargeInvoiceNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.ChargeInvoiceNotification{
		Type: webhooks.NewChargeInvoice,
		Account: webhooks.Account{
			XMLName: xml.Name{Local: "account"},
			Code:    "1234",
		},
		Invoice: webhooks.ChargeInvoice{
			XMLName:           xml.Name{Local: "invoice"},
			UUID:              "42feb03ce368c0e1ead35d4bfa89b82e",
			State:             recurly.ChargeInvoiceStatePending,
			Origin:            recurly.ChargeInvoiceOriginRenewal,
			SubscriptionUUIDs: []string{"40b8f5e99df03b8684b99d4993b6e089"},
			InvoiceNumber:     2405,
			Currency:          "USD",
			BalanceInCents:    100000,
			TotalInCents:      100000,
			NetTerms:          recurly.NullInt{Int: 30, Valid: true},
			CollectionMethod:  recurly.CollectionMethodManual,
			CreatedAt:         recurly.NewTime(created),
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_CreditInvoiceNotification(t *testing.T) {
	d, _ := time.Parse(recurly.DateTimeFormat, "2018-02-13T00:56:22Z")
	xmlFile := MustOpenFile("testdata/credit_invoice_notification.xml")
	result, err := webhooks.Parse(xmlFile)
	if err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.CreditInvoiceNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.CreditInvoiceNotification{
		Type: webhooks.NewCreditInvoice,
		Account: webhooks.Account{
			XMLName: xml.Name{Local: "account"},
			Code:    "1234",
		},
		Invoice: webhooks.CreditInvoice{
			XMLName:           xml.Name{Local: "invoice"},
			UUID:              "42fb74de65e9395eb004614144a7b91f",
			State:             recurly.CreditInvoiceStateClosed,
			Origin:            recurly.CreditInvoiceOriginWriteOff,
			SubscriptionUUIDs: []string{"42fb74ba9efe4c6981c2064436a4e9cd"},
			InvoiceNumber:     2404,
			Currency:          "USD",
			BalanceInCents:    0,
			TotalInCents:      -4882,
			CreatedAt:         recurly.NewTime(d),
			ClosedAt:          recurly.NewTime(d),
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_CreditPaymentNotification(t *testing.T) {
	d, _ := time.Parse(recurly.DateTimeFormat, "2018-02-12T18:55:20Z")
	xmlFile := MustOpenFile("testdata/credit_payment_notification.xml")
	result, err := webhooks.Parse(xmlFile)
	if err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.CreditPaymentNotification); !ok {
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
			CreatedAt:              recurly.NewTime(d),
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_NewSubscriptionNotification(t *testing.T) {
	activatedTs, _ := time.Parse(recurly.DateTimeFormat, "2010-09-23T22:05:03Z")
	canceledTs, _ := time.Parse(recurly.DateTimeFormat, "2010-09-23T22:05:43Z")
	expiresTs, _ := time.Parse(recurly.DateTimeFormat, "2010-09-24T22:05:03Z")
	startedTs, _ := time.Parse(recurly.DateTimeFormat, "2010-09-23T22:05:03Z")
	endsTs, _ := time.Parse(recurly.DateTimeFormat, "2010-09-24T22:05:03Z")

	xmlFile := MustOpenFile("testdata/new_subscription_notification.xml")
	result, err := webhooks.Parse(xmlFile)
	if err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.SubscriptionNotification); !ok {
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
			ActivatedAt:            recurly.NewTime(activatedTs),
			CanceledAt:             recurly.NewTime(canceledTs),
			ExpiresAt:              recurly.NewTime(expiresTs),
			CurrentPeriodStartedAt: recurly.NewTime(startedTs),
			CurrentPeriodEndsAt:    recurly.NewTime(endsTs),
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_UpdatedSubscriptionNotification(t *testing.T) {
	activatedTs, _ := time.Parse(recurly.DateTimeFormat, "2010-09-23T22:05:03Z")
	canceledTs, _ := time.Parse(recurly.DateTimeFormat, "2010-09-23T22:05:43Z")
	expiresTs, _ := time.Parse(recurly.DateTimeFormat, "2010-09-24T22:05:03Z")
	startedTs, _ := time.Parse(recurly.DateTimeFormat, "2010-09-23T22:05:03Z")
	endsTs, _ := time.Parse(recurly.DateTimeFormat, "2010-09-24T22:05:03Z")

	xmlFile := MustOpenFile("testdata/updated_subscription_notification.xml")
	result, err := webhooks.Parse(xmlFile)
	if err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.SubscriptionNotification); !ok {
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
			ActivatedAt:            recurly.NewTime(activatedTs),
			CanceledAt:             recurly.NewTime(canceledTs),
			ExpiresAt:              recurly.NewTime(expiresTs),
			CurrentPeriodStartedAt: recurly.NewTime(startedTs),
			CurrentPeriodEndsAt:    recurly.NewTime(endsTs),
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_RenewedSubscriptionNotification(t *testing.T) {
	activatedTs, _ := time.Parse(recurly.DateTimeFormat, "2010-07-22T20:42:05Z")
	startedTs, _ := time.Parse(recurly.DateTimeFormat, "2010-09-22T20:42:05Z")
	endsTs, _ := time.Parse(recurly.DateTimeFormat, "2010-10-22T20:42:05Z")

	xmlFile := MustOpenFile("testdata/renewed_subscription_notification.xml")
	result, err := webhooks.Parse(xmlFile)
	if err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.SubscriptionNotification); !ok {
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
			ActivatedAt:            recurly.NewTime(activatedTs),
			CurrentPeriodStartedAt: recurly.NewTime(startedTs),
			CurrentPeriodEndsAt:    recurly.NewTime(endsTs),
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_ExpiredSubscriptionNotification(t *testing.T) {
	activatedTs, _ := time.Parse(recurly.DateTimeFormat, "2010-09-23T22:05:03Z")
	canceledTs, _ := time.Parse(recurly.DateTimeFormat, "2010-09-23T22:05:43Z")
	expiresTs, _ := time.Parse(recurly.DateTimeFormat, "2010-09-24T22:05:03Z")
	startedTs, _ := time.Parse(recurly.DateTimeFormat, "2010-09-23T22:05:03Z")
	endsTs, _ := time.Parse(recurly.DateTimeFormat, "2010-09-24T22:05:03Z")

	xmlFile := MustOpenFile("testdata/expired_subscription_notification.xml")
	result, err := webhooks.Parse(xmlFile)
	if err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.SubscriptionNotification); !ok {
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
			ActivatedAt:            recurly.NewTime(activatedTs),
			CanceledAt:             recurly.NewTime(canceledTs),
			ExpiresAt:              recurly.NewTime(expiresTs),
			CurrentPeriodStartedAt: recurly.NewTime(startedTs),
			CurrentPeriodEndsAt:    recurly.NewTime(endsTs),
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_CanceledSubscriptionNotification(t *testing.T) {
	activatedTs, _ := time.Parse(recurly.DateTimeFormat, "2010-09-23T22:05:03Z")
	canceledTs, _ := time.Parse(recurly.DateTimeFormat, "2010-09-23T22:05:43Z")
	expiresTs, _ := time.Parse(recurly.DateTimeFormat, "2010-09-24T22:05:03Z")
	startedTs, _ := time.Parse(recurly.DateTimeFormat, "2010-09-23T22:05:03Z")
	endsTs, _ := time.Parse(recurly.DateTimeFormat, "2010-09-24T22:05:03Z")

	xmlFile := MustOpenFile("testdata/canceled_subscription_notification.xml")
	result, err := webhooks.Parse(xmlFile)
	if err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.SubscriptionNotification); !ok {
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
			ActivatedAt:            recurly.NewTime(activatedTs),
			CanceledAt:             recurly.NewTime(canceledTs),
			ExpiresAt:              recurly.NewTime(expiresTs),
			CurrentPeriodStartedAt: recurly.NewTime(startedTs),
			CurrentPeriodEndsAt:    recurly.NewTime(endsTs),
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_NewInvoiceNotification(t *testing.T) {
	xmlFile := MustOpenFile("testdata/new_invoice_notification.xml")
	createdAt := time.Date(2014, 1, 1, 20, 21, 44, 0, time.UTC)
	result, err := webhooks.Parse(xmlFile)
	if err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.InvoiceNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.InvoiceNotification{
		Type: webhooks.NewInvoice,
		Account: webhooks.Account{
			XMLName:   xml.Name{Local: "account"},
			Code:      "1",
			Email:     "verena@example.com",
			FirstName: "Verena",
			LastName:  "Example",
		},
		Invoice: webhooks.Invoice{
			XMLName:          xml.Name{Local: "invoice"},
			UUID:             "ffc64d71d4b5404e93f13aac9c63b007",
			State:            "open",
			Currency:         "USD",
			CreatedAt:        recurly.NullTime{Time: &createdAt},
			InvoiceNumber:    1000,
			TotalInCents:     1000,
			NetTerms:         recurly.NullInt{Valid: true, Int: 0},
			CollectionMethod: recurly.CollectionMethodManual,
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_PastDueInvoiceNotification(t *testing.T) {
	xmlFile := MustOpenFile("testdata/past_due_invoice_notification.xml")
	createdAt := time.Date(2014, 1, 1, 20, 21, 44, 0, time.UTC)
	result, err := webhooks.Parse(xmlFile)
	if err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.InvoiceNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.InvoiceNotification{
		Type: webhooks.PastDueInvoice,
		Account: webhooks.Account{
			XMLName:     xml.Name{Local: "account"},
			Code:        "1",
			Username:    "verena",
			Email:       "verena@example.com",
			FirstName:   "Verena",
			LastName:    "Example",
			CompanyName: "Company, Inc.",
		},
		Invoice: webhooks.Invoice{
			XMLName:       xml.Name{Local: "invoice"},
			UUID:          "ffc64d71d4b5404e93f13aac9c63b007",
			State:         "past_due",
			CreatedAt:     recurly.NullTime{Time: &createdAt},
			InvoiceNumber: 1000,
			TotalInCents:  1100,
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_SuccessfulPaymentNotification(t *testing.T) {
	xmlFile := MustOpenFile("testdata/successful_payment_notification.xml")
	if result, err := webhooks.Parse(xmlFile); err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.PaymentNotification); !ok {
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
			AmountInCents: 1000,
			Status:        "success",
			Message:       "Bogus Gateway: Forced success",
			Reference:     "reference",
			Source:        "subscription",
			Test:          recurly.NullBool{Valid: true, Bool: true},
			Voidable:      recurly.NullBool{Valid: true, Bool: true},
			Refundable:    recurly.NullBool{Valid: true, Bool: true},
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_FailedPaymentNotification(t *testing.T) {
	xmlFile := MustOpenFile("testdata/failed_payment_notification.xml")
	if result, err := webhooks.Parse(xmlFile); err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.PaymentNotification); !ok {
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
			AmountInCents:    1000,
			Status:           "Declined",
			Message:          "This transaction has been declined",
			FailureType:      "Declined by the gateway",
			Reference:        "reference",
			Source:           "subscription",
			Test:             recurly.NullBool{Valid: true, Bool: true},
			Voidable:         recurly.NullBool{Valid: true, Bool: false},
			Refundable:       recurly.NullBool{Valid: true, Bool: false},
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_VoidPaymentNotification(t *testing.T) {
	xmlFile := MustOpenFile("testdata/void_payment_notification.xml")
	if result, err := webhooks.Parse(xmlFile); err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.PaymentNotification); !ok {
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
			Test:             recurly.NullBool{Valid: true, Bool: true},
			Voidable:         recurly.NullBool{Valid: true, Bool: true},
			Refundable:       recurly.NullBool{Valid: true, Bool: true},
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_SuccessfulRefundNotification(t *testing.T) {
	xmlFile := MustOpenFile("testdata/successful_refund_notification.xml")
	if result, err := webhooks.Parse(xmlFile); err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.PaymentNotification); !ok {
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
			Test:             recurly.NullBool{Valid: true, Bool: true},
			Voidable:         recurly.NullBool{Valid: true, Bool: true},
			Refundable:       recurly.NullBool{Valid: true, Bool: true},
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse_ErrUnknownNotification(t *testing.T) {
	xmlFile := MustOpenFile("testdata/unknown_notification.xml")
	result, err := webhooks.Parse(xmlFile)
	if result != nil {
		t.Fatalf("unexpected notification: %#v", result)
	} else if e, ok := err.(webhooks.ErrUnknownNotification); !ok {
		t.Fatalf("unexpected type: %T, error", e)
	} else if err.Error() != "unknown notification: unknown_notification" {
		t.Fatalf("unexpected error string: %s", err.Error())
	} else if e.Name() != "unknown_notification" {
		t.Fatalf("unexpected notification name: %s", e.Name())
	}
}

func MustOpenFile(name string) *os.File {
	file, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	return file
}
