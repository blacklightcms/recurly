package webhooks_test

import (
	"encoding/xml"
	"os"
	"reflect"
	"testing"

	recurly "github.com/blacklightcms/go-recurly"
	"github.com/blacklightcms/go-recurly/webhooks"
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
