package webhooks_test

import (
	"encoding/xml"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/blacklightcms/recurly"
	"github.com/blacklightcms/recurly/webhooks"
)

func MustOpenFile(name string) *os.File {
	file, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	return file
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
		t.Fatalf("unexpected type: %T, result")
	} else if !reflect.DeepEqual(n, &webhooks.ExpiredSubscriptionNotification{
		Account: recurly.Account{
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
			ActivatedAt:            recurly.NewTime(activatedTs),
			CanceledAt:             recurly.NewTime(canceledTs),
			ExpiresAt:              recurly.NewTime(expiresTs),
			CurrentPeriodStartedAt: recurly.NewTime(startedTs),
			CurrentPeriodEndsAt:    recurly.NewTime(endsTs),
		},
	}) {
		t.Fatalf("unexpected notification: %#v", n)
	}
}

func TestParse_NewInvoiceNotification(t *testing.T) {
	xmlFile := MustOpenFile("testdata/new_invoice_notification.xml")
	result, err := webhooks.Parse(xmlFile)
	if err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.NewInvoiceNotification); !ok {
		t.Fatalf("unexpected type: %T, result")
	} else if !reflect.DeepEqual(n, &webhooks.NewInvoiceNotification{
		Account: recurly.Account{
			XMLName:   xml.Name{Local: "account"},
			Code:      "1",
			Email:     "verena@example.com",
			FirstName: "Verena",
			LastName:  "Example",
		},
		Invoice: recurly.Invoice{
			XMLName:          xml.Name{Local: "invoice"},
			UUID:             "ffc64d71d4b5404e93f13aac9c63b007",
			State:            "open",
			Currency:         "USD",
			InvoiceNumber:    1000,
			TotalInCents:     1000,
			NetTerms:         recurly.NullInt{Valid: true, Int: 0},
			CollectionMethod: recurly.CollectionMethodManual,
		},
	}) {
		t.Fatalf("unexpected notification: %v", n)
	}
}

func TestParse_PastDueInvoiceNotification(t *testing.T) {
	xmlFile := MustOpenFile("testdata/past_due_invoice_notification.xml")
	result, err := webhooks.Parse(xmlFile)
	if err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.PastDueInvoiceNotification); !ok {
		t.Fatalf("unexpected type: %T, result")
	} else if !reflect.DeepEqual(n, &webhooks.PastDueInvoiceNotification{
		Account: recurly.Account{
			XMLName:     xml.Name{Local: "account"},
			Code:        "1",
			Username:    "verena",
			Email:       "verena@example.com",
			FirstName:   "Verena",
			LastName:    "Example",
			CompanyName: "Company, Inc.",
		},
		Invoice: recurly.Invoice{
			XMLName:       xml.Name{Local: "invoice"},
			UUID:          "ffc64d71d4b5404e93f13aac9c63b007",
			State:         "past_due",
			InvoiceNumber: 1000,
			TotalInCents:  1100,
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
		t.Fatalf("unexpected type: %T, result")
	} else if !reflect.DeepEqual(n, &webhooks.SuccessfulPaymentNotification{
		Account: recurly.Account{
			XMLName:     xml.Name{Local: "account"},
			Code:        "1",
			Username:    "verena",
			Email:       "verena@example.com",
			FirstName:   "Verena",
			LastName:    "Example",
			CompanyName: "Company, Inc.",
		},
		Transaction: recurly.Transaction{
			UUID:          "a5143c1d3a6f4a8287d0e2cc1d4c0427",
			InvoiceNumber: 2059,
			Action:        "purchase",
			AmountInCents: 1000,
			Status:        "success",
			Reference:     "reference",
			Source:        "subscription",
			Test:          true,
			Voidable:      recurly.NullBool{Bool: true, Valid: true},
			Refundable:    recurly.NullBool{Bool: true, Valid: true},
		},
	}) {
		t.Fatalf("unexpected notification: %#v", n)
	}
}

func TestParse_FailedPaymentNotification(t *testing.T) {
	xmlFile := MustOpenFile("testdata/failed_payment_notification.xml")
	if result, err := webhooks.Parse(xmlFile); err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*webhooks.FailedPaymentNotification); !ok {
		t.Fatalf("unexpected type: %T, result")
	} else if !reflect.DeepEqual(n, &webhooks.FailedPaymentNotification{
		Account: recurly.Account{
			XMLName:     xml.Name{Local: "account"},
			Code:        "1",
			Username:    "verena",
			Email:       "verena@example.com",
			FirstName:   "Verena",
			LastName:    "Example",
			CompanyName: "Company, Inc.",
		},
		Transaction: recurly.Transaction{
			UUID:          "a5143c1d3a6f4a8287d0e2cc1d4c0427",
			InvoiceNumber: 2059,
			Action:        "purchase",
			AmountInCents: 1000,
			Status:        "Declined",
			Reference:     "reference",
			Source:        "subscription",
			Test:          true,
			Voidable:      recurly.NullBool{Bool: false, Valid: true},
			Refundable:    recurly.NullBool{Bool: false, Valid: true},
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
		t.Fatalf("unexpected type: %T, result")
	} else if err.Error() != "unknown notification: unknown_notification" {
		t.Fatalf("unexpected error string: %s", err.Error())
	} else if e.Name() != "unknown_notification" {
		t.Fatalf("unexpected notification name: %s", e.Name())
	}
}
