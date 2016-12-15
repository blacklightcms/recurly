package notifications_test

import (
	"encoding/xml"
	"os"
	"reflect"
	"testing"

	"github.com/blacklightcms/go-recurly/recurly"
	"github.com/blacklightcms/go-recurly/recurly/notifications"
)

func MustOpenFile(name string) *os.File {
	file, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	return file
}

func TestParse_SuccessfulPaymentNotification(t *testing.T) {
	xmlFile := MustOpenFile("testdata/successful_payment_notification.xml")
	if result, err := notifications.Parse(xmlFile); err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*notifications.SuccessfulPaymentNotification); !ok {
		t.Fatalf("unable to reflect interface")
	} else if !reflect.DeepEqual(n, &notifications.SuccessfulPaymentNotification{
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
			XMLName:       xml.Name{Local: "transaction"},
			UUID:          "a5143c1d3a6f4a8287d0e2cc1d4c0427",
			Action:        "purchase",
			AmountInCents: 1000,
			Status:        "success",
			Reference:     "reference",
			Source:        "subscription",
			Test:          true,
			Voidable:      recurly.NullBool{Bool: true, Valid: true},
			Refundable:    recurly.NullBool{Bool: true, Valid: true},
		},
		InvoiceNumber: 2059,
	}) {
		t.Fatalf("unexpected notification: %#v", n)
	}
}

func TestParse_FailedPaymentNotification(t *testing.T) {
	xmlFile := MustOpenFile("testdata/failed_payment_notification.xml")
	if result, err := notifications.Parse(xmlFile); err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*notifications.FailedPaymentNotification); !ok {
		t.Fatalf("unable to reflect interface")
	} else if !reflect.DeepEqual(n, &notifications.FailedPaymentNotification{
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
			XMLName:       xml.Name{Local: "transaction"},
			UUID:          "a5143c1d3a6f4a8287d0e2cc1d4c0427",
			Action:        "purchase",
			AmountInCents: 1000,
			Status:        "Declined",
			Reference:     "reference",
			Source:        "subscription",
			Test:          true,
			Voidable:      recurly.NullBool{Bool: false, Valid: true},
			Refundable:    recurly.NullBool{Bool: false, Valid: true},
		},
		InvoiceNumber: 2059,
	}) {
		t.Fatalf("unexpected notification: %#v", n)
	}
}

func TestParse_PastDueInvoiceNotification(t *testing.T) {
	xmlFile := MustOpenFile("testdata/past_due_invoice_notification.xml")
	result, err := notifications.Parse(xmlFile)
	if err != nil {
		t.Fatal(err)
	} else if n, ok := result.(*notifications.PastDueInvoiceNotification); !ok {
		t.Fatalf("unable to reflect interface")
	} else if !reflect.DeepEqual(n, &notifications.PastDueInvoiceNotification{
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

func TestParse_ErrUnknownNotification(t *testing.T) {
	xmlFile := MustOpenFile("testdata/unknown_notification.xml")
	result, err := notifications.Parse(xmlFile)
	if result != nil {
		t.Fatalf("unexpected notification: %#v", result)
	} else if e, ok := err.(notifications.ErrUnknownNotification); !ok {
		t.Fatalf("unable to reflect interface")
	} else if err.Error() != "unknown notification: unknown_notification" {
		t.Fatalf("unexpected error string: %s", err.Error())
	} else if e.Name() != "unknown_notification" {
		t.Fatalf("unexpected notification name: %s", e.Name())
	}
}
