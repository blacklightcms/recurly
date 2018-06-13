package webhooks_test

import (
	"encoding/xml"
	"os"
	"testing"
	"time"

	"github.com/blacklightcms/recurly"
	"github.com/blacklightcms/recurly/webhooks"
	"github.com/google/go-cmp/cmp"
	"reflect"
)

func TestParse_NewAccountNotification(t *testing.T) {
	xmlFile := MustOpenFile("testdata/new_account_notification.xml")
	result, err := webhooks.Parse(xmlFile)
	if err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.NewAccountNotification); !ok {
		t.Fatalf("unexpected type: %T, result")
	} else if !reflect.DeepEqual(n, &webhooks.NewAccountNotification{
		Account: webhooks.Account{
			XMLName:   xml.Name{Local: "account"},
			Code:      "1",
			Email:     "verena@example.com",
			FirstName: "Verena",
			LastName:  "Example",
		},
	}) {
		t.Fatalf("unexpected notification: %#v", n)
	}
}

func TestParse_UpdatedAccountNotification(t *testing.T) {
	xmlFile := MustOpenFile("testdata/updated_account_notification.xml")
	result, err := webhooks.Parse(xmlFile)
	if err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.UpdatedAccountNotification); !ok {
		t.Fatalf("unexpected type: %T, result")
	} else if !reflect.DeepEqual(n, &webhooks.UpdatedAccountNotification{
		Account: webhooks.Account{
			XMLName:   xml.Name{Local: "account"},
			Code:      "1",
			Email:     "verena@example.com",
			FirstName: "Verena",
			LastName:  "Example",
		},
	}) {
		t.Fatalf("unexpected notification: %#v", n)
	}
}

func TestParse_CanceledAccountNotification(t *testing.T) {
	xmlFile := MustOpenFile("testdata/canceled_account_notification.xml")
	result, err := webhooks.Parse(xmlFile)
	if err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.CanceledAccountNotification); !ok {
		t.Fatalf("unexpected type: %T, result")
	} else if !reflect.DeepEqual(n, &webhooks.CanceledAccountNotification{
		Account: webhooks.Account{
			XMLName:   xml.Name{Local: "account"},
			Code:      "1",
			Email:     "verena@example.com",
			FirstName: "Verena",
			LastName:  "Example",
		},
	}) {
		t.Fatalf("unexpected notification: %#v", n)
	}
}

func TestParse_BillingInfoUpdatedNotification(t *testing.T) {
	xmlFile := MustOpenFile("testdata/billing_info_updated_notification.xml")
	result, err := webhooks.Parse(xmlFile)
	if err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.BillingInfoUpdatedNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.BillingInfoUpdatedNotification{
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
	xmlFile := MustOpenFile("testdata/billing_info_update_failed_notification.xml")
	result, err := webhooks.Parse(xmlFile)
	if err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.BillingInfoUpdateFailedNotification); !ok {
		t.Fatalf("unexpected type: %T, result")
	} else if !reflect.DeepEqual(n, &webhooks.BillingInfoUpdateFailedNotification{
		Account: webhooks.Account{
			XMLName:   xml.Name{Local: "account"},
			Code:      "1",
			Email:     "verena@example.com",
			FirstName: "Verena",
			LastName:  "Example",
		},
	}) {
		t.Fatalf("unexpected notification: %#v", n)
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
	} else if n, ok := result.(*webhooks.NewSubscriptionNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.NewSubscriptionNotification{
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
	} else if n, ok := result.(*webhooks.UpdatedSubscriptionNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.UpdatedSubscriptionNotification{
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
	} else if n, ok := result.(*webhooks.RenewedSubscriptionNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.RenewedSubscriptionNotification{
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
	} else if n, ok := result.(*webhooks.ExpiredSubscriptionNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.ExpiredSubscriptionNotification{
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
	} else if n, ok := result.(*webhooks.CanceledSubscriptionNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.CanceledSubscriptionNotification{
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

func TestParse_ReactivatedAccountNotification(t *testing.T) {
	activatedTs, _ := time.Parse(recurly.DateTimeFormat, "2010-07-22T20:42:05Z")
	startedTs, _ := time.Parse(recurly.DateTimeFormat, "2010-09-22T20:42:05Z")
	endsTs, _ := time.Parse(recurly.DateTimeFormat, "2010-10-22T20:42:05Z")

	xmlFile := MustOpenFile("testdata/reactivated_account_notification.xml")
	result, err := webhooks.Parse(xmlFile)
	if err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.ReactivatedAccountNotification); !ok {
		t.Fatalf("unexpected type: %T, result")
	} else if !reflect.DeepEqual(n, &webhooks.ReactivatedAccountNotification{
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
			ActivatedAt:            recurly.NewTime(activatedTs),
			CurrentPeriodStartedAt: recurly.NewTime(startedTs),
			CurrentPeriodEndsAt:    recurly.NewTime(endsTs),
		},
	}) {
		t.Fatalf("unexpected notification: %#v", n)
	}
}

func TestParse_NewInvoiceNotification(t *testing.T) {
	xmlFile := MustOpenFile("testdata/new_invoice_notification.xml")
	createdAt := time.Date(2014, 1, 1, 20, 21, 44, 0, time.UTC)
	result, err := webhooks.Parse(xmlFile)
	if err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.NewInvoiceNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.NewInvoiceNotification{
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
	} else if n, ok := result.(*webhooks.PastDueInvoiceNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.PastDueInvoiceNotification{
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

func TestParse_ProcessingInvoiceNotification(t *testing.T) {
	xmlFile := MustOpenFile("testdata/processing_invoice_notification.xml")
	createdAt := time.Date(2014, 1, 1, 20, 21, 44, 0, time.UTC)
	result, err := webhooks.Parse(xmlFile)
	if err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.ProcessingInvoiceNotification); !ok {
		t.Fatalf("unexpected type: %T, result")
	} else if !reflect.DeepEqual(n, &webhooks.ProcessingInvoiceNotification{
		Account: webhooks.Account{
			XMLName:     xml.Name{Local: "account"},
			Code:        "1",
			Username:    "",
			Email:       "verena@example.com",
			FirstName:   "Verana",
			LastName:    "Example",
			CompanyName: "",
			Phone:       "",
		},
		Invoice: webhooks.Invoice{
			XMLName:             xml.Name{Local: "invoice"},
			SubscriptionUUID:    "",
			UUID:                "ffc64d71d4b5404e93f13aac9c63b007",
			State:               "processing",
			InvoiceNumberPrefix: "",
			InvoiceNumber:       1000,
			PONumber:            "",
			VATNumber:           "",
			TotalInCents:        1000,
			Currency:            "USD",
			CreatedAt:           recurly.NullTime{Time: &createdAt},
			ClosedAt:            recurly.NullTime{},
			NetTerms:            recurly.NullInt{Int: 0, Valid: true},
			CollectionMethod:    recurly.CollectionMethodAutomatic,
		},
	}) {
		t.Fatalf("unexpected notification: %v", n)
	}
}

func TestParse_ClosedInvoiceNotification(t *testing.T) {
	xmlFile := MustOpenFile("testdata/closed_invoice_notification.xml")
	createdAt := time.Date(2014, 1, 1, 20, 20, 29, 0, time.UTC)
	closedAt := time.Date(2014, 1, 1, 20, 24, 02, 0, time.UTC)
	result, err := webhooks.Parse(xmlFile)
	if err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.ClosedInvoiceNotification); !ok {
		t.Fatalf("unexpected type: %T, result")
	} else if !reflect.DeepEqual(n, &webhooks.ClosedInvoiceNotification{
		Account: webhooks.Account{
			XMLName:     xml.Name{Local: "account"},
			Code:        "1",
			Username:    "",
			Email:       "verena@example.com",
			FirstName:   "Verana",
			LastName:    "Example",
			CompanyName: "",
			Phone:       "",
		},
		Invoice: webhooks.Invoice{
			XMLName:             xml.Name{Local: "invoice"},
			SubscriptionUUID:    "",
			UUID:                "ffc64d71d4b5404e93f13aac9c63b007",
			State:               "collected",
			InvoiceNumberPrefix: "",
			InvoiceNumber:       1000,
			PONumber:            "",
			VATNumber:           "",
			TotalInCents:        1100,
			Currency:            "USD",
			CreatedAt:           recurly.NullTime{Time: &createdAt},
			ClosedAt:            recurly.NullTime{Time: &closedAt},
			NetTerms:            recurly.NullInt{},
			CollectionMethod:    "",
		},
	}) {
		t.Fatalf("unexpected notification: %v", n)
	}
}

func TestParse_SuccessfulPaymentNotification(t *testing.T) {
	xmlFile := MustOpenFile("testdata/successful_payment_notification.xml")
	if result, err := webhooks.Parse(xmlFile); err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.SuccessfulPaymentNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.SuccessfulPaymentNotification{
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
	} else if n, ok := result.(*webhooks.FailedPaymentNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.FailedPaymentNotification{
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
	} else if n, ok := result.(*webhooks.VoidPaymentNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.VoidPaymentNotification{
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
	} else if n, ok := result.(*webhooks.SuccessfulRefundNotification); !ok {
		t.Fatalf("unexpected type: %T, result", n)
	} else if diff := cmp.Diff(n, &webhooks.SuccessfulRefundNotification{
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

func TestParse_ScheduledPaymentNotification(t *testing.T) {
	xmlFile := MustOpenFile("testdata/scheduled_payment_notification.xml")
	if result, err := webhooks.Parse(xmlFile); err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.ScheduledPaymentNotification); !ok {
		t.Fatalf("unexpected type: %T, result")
	} else if !reflect.DeepEqual(n, &webhooks.ScheduledPaymentNotification{
		Account: webhooks.Account{
			XMLName:     xml.Name{Local: "account"},
			Code:        "1",
			Username:    "verena",
			Email:       "verena@example.com",
			FirstName:   "Verena",
			LastName:    "Example",
			CompanyName: "Company, Inc.",
			Phone:       "",
		},
		Transaction: webhooks.Transaction{
			XMLName:           xml.Name{Local: "transaction"},
			UUID:              "a5143c1d3a6f4a8287d0e2cc1d4c0427",
			InvoiceNumber:     2059,
			SubscriptionUUID:  "1974a098jhlkjasdfljkha898326881c",
			Action:            "purchase",
			AmountInCents:     1000,
			Status:            "scheduled",
			Message:           "Bogus Gateway: Forced success",
			GatewayErrorCodes: "",
			FailureType:       "",
			Reference:         "",
			Source:            "subscription",
			Test:              recurly.NullBool{Bool: true, Valid: true},
			Voidable:          recurly.NullBool{Bool: true, Valid: true},
			Refundable:        recurly.NullBool{Bool: true, Valid: true},
		},
	}) {
		t.Fatalf("unexpected notification: %#v", n)
	}
}

func TestParse_ProcessingPaymentNotification(t *testing.T) {
	xmlFile := MustOpenFile("testdata/processing_payment_notification.xml")
	if result, err := webhooks.Parse(xmlFile); err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.ProcessingPaymentNotification); !ok {
		t.Fatalf("unexpected type: %T, result")
	} else if !reflect.DeepEqual(n, &webhooks.ProcessingPaymentNotification{
		Account: webhooks.Account{
			XMLName:     xml.Name{Local: "account"},
			Code:        "1",
			Username:    "verena",
			Email:       "verena@example.com",
			FirstName:   "Verena",
			LastName:    "Example",
			CompanyName: "Company, Inc.",
			Phone:       "",
		},
		Transaction: webhooks.Transaction{
			XMLName:           xml.Name{Local: "transaction"},
			UUID:              "a5143c1d3a6f4a8287d0e2cc1d4c0427",
			InvoiceNumber:     2059,
			SubscriptionUUID:  "1974a098jhlkjasdfljkha898326881c",
			Action:            "purchase",
			AmountInCents:     1000,
			Status:            "processing",
			Message:           "Bogus Gateway: Forced success",
			GatewayErrorCodes: "",
			FailureType:       "",
			Reference:         "",
			Source:            "subscription",
			Test:              recurly.NullBool{Bool: true, Valid: true},
			Voidable:          recurly.NullBool{Bool: true, Valid: true},
			Refundable:        recurly.NullBool{Bool: true, Valid: true},
		},
	}) {
		t.Fatalf("unexpected notification: %#v", n)
	}
}

func TestParse_NewDunningEventNotification(t *testing.T) {
	invoiceCreatedAt := time.Date(2016, 10, 26, 16, 00, 12, 0, time.UTC)
	invoiceClosedAt := time.Date(2016, 10, 27, 16, 00, 26, 0, time.UTC)
	subscriptionActivatedAt := time.Date(2016, 10, 26, 05, 42, 27, 0, time.UTC)
	subscriptionPeriodStart := time.Date(2016, 10, 26, 16, 00, 00, 0, time.UTC)
	subscriptionPeriodEnd := time.Date(2016, 11, 26, 16, 00, 00, 0, time.UTC)

	xmlFile := MustOpenFile("testdata/new_dunning_event_notification.xml")
	if result, err := webhooks.Parse(xmlFile); err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.NewDunningEventNotification); !ok {
		t.Fatalf("unexpected type: %T, result")
	} else if !reflect.DeepEqual(n, &webhooks.NewDunningEventNotification{
		Account: webhooks.Account{
			XMLName:     xml.Name{Local: "account"},
			Code:        "09f299492d21",
			Username:    "",
			Email:       "joseph.smith@gmail.com",
			FirstName:   "Joseph",
			LastName:    "Smith",
			CompanyName: "",
			Phone:       "3235626924",
		},
		Invoice: webhooks.Invoice{
			XMLName:             xml.Name{Local: "invoice"},
			SubscriptionUUID:    "396e4e17640ca516c2f3a84e47ae91dd",
			UUID:                "inv-7wr0r2xuawwCjO",
			State:               "failed",
			InvoiceNumberPrefix: "",
			InvoiceNumber:       781002,
			PONumber:            "",
			VATNumber:           "",
			TotalInCents:        2499,
			Currency:            "USD",
			CreatedAt:           recurly.NullTime{Time: &invoiceCreatedAt},
			ClosedAt:            recurly.NullTime{Time: &invoiceClosedAt},
			NetTerms:            recurly.NullInt{Int: 0, Valid: true},
			CollectionMethod:    "automatic",
		},
		Subscription: recurly.Subscription{
			XMLName:                xml.Name{Local: "subscription"},
			Plan:                   recurly.NestedPlan{Code: "28a3ae1fc5c00d123429", Name: "41c36e04f2d7bebc"},
			AccountCode:            "",
			InvoiceNumber:          0,
			UUID:                   "396e4e17640ca516c2f3a84e47ae91dd",
			State:                  "active",
			UnitAmountInCents:      0,
			Currency:               "",
			Quantity:               1,
			TotalAmountInCents:     2499,
			ActivatedAt:            recurly.NullTime{Time: &subscriptionActivatedAt},
			CanceledAt:             recurly.NullTime{},
			ExpiresAt:              recurly.NullTime{},
			CurrentPeriodStartedAt: recurly.NullTime{Time: &subscriptionPeriodStart},
			CurrentPeriodEndsAt:    recurly.NullTime{Time: &subscriptionPeriodEnd},
			TrialStartedAt:         recurly.NullTime{},
			TrialEndsAt:            recurly.NullTime{},
			TaxInCents:             0,
			TaxType:                "",
			TaxRegion:              "",
			TaxRate:                0,
			PONumber:               "",
			NetTerms:               recurly.NullInt{},
			SubscriptionAddOns:     nil,
			PendingSubscription:    (*recurly.PendingSubscription)(nil),
		},
		Transaction: webhooks.Transaction{
			XMLName:           xml.Name{Local: "transaction"},
			UUID:              "397083a9a871b53a3d5a4c469fa1216a",
			InvoiceNumber:     1002,
			SubscriptionUUID:  "396e4e17640ca516c2f3a84e47ae91dd",
			Action:            "purchase",
			AmountInCents:     2499,
			Status:            "declined",
			Message:           "Transaction Normal",
			GatewayErrorCodes: "00",
			FailureType:       "invalid_data",
			Reference:         "115948823",
			Source:            "subscription",
			Test:              recurly.NullBool{Bool: true, Valid: true},
			Voidable:          recurly.NullBool{Bool: false, Valid: true},
			Refundable:        recurly.NullBool{Bool: false, Valid: true},
		},
	}) {
		t.Fatalf("unexpected notification: %#v", n)
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
