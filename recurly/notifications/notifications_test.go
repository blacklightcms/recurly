package notifications_test

import (
	"encoding/xml"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/blacklightcms/go-recurly/recurly"
	"github.com/blacklightcms/go-recurly/recurly/notifications"
)

func TestParse_SuccessfulPaymentNotification(t *testing.T) {
	xmlFile, err := os.Open("xml/successful_payment_notification.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer xmlFile.Close()

	req, err := http.NewRequest("POST", "/", xmlFile)
	if err != nil {
		t.Fatal(err)
	}

	result, err := notifications.Parse(req.Body)
	if err != nil {
		t.Fatal(err)
	}

	var n *notifications.SuccessfulPaymentNotification
	var ok bool
	switch result.(type) {
	case *notifications.SuccessfulPaymentNotification:
		n, ok = result.(*notifications.SuccessfulPaymentNotification)
		if !ok {
			t.Fatalf("unable to reflect interface")
		}
		break
	default:
		t.Fatalf("unexpected result type: %T", result)
	}

	if !reflect.DeepEqual(n, &notifications.SuccessfulPaymentNotification{
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
	xmlFile, err := os.Open("xml/failed_payment_notification.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer xmlFile.Close()

	req, err := http.NewRequest("POST", "/", xmlFile)
	if err != nil {
		t.Fatal(err)
	}

	result, err := notifications.Parse(req.Body)
	if err != nil {
		t.Fatal(err)
	}

	var n *notifications.FailedPaymentNotification
	var ok bool
	switch result.(type) {
	case *notifications.FailedPaymentNotification:
		n, ok = result.(*notifications.FailedPaymentNotification)
		if !ok {
			t.Fatalf("unable to reflect interface")
		}
		break
	default:
		t.Fatalf("unexpected result type: %T", result)
	}

	if !reflect.DeepEqual(n, &notifications.FailedPaymentNotification{
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
	xmlFile, err := os.Open("xml/past_due_invoice_notification.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer xmlFile.Close()

	req, err := http.NewRequest("POST", "/", xmlFile)
	if err != nil {
		t.Fatal(err)
	}

	result, err := notifications.Parse(req.Body)
	if err != nil {
		t.Fatal(err)
	}

	var n *notifications.PastDueInvoiceNotification
	var ok bool
	switch result.(type) {
	case *notifications.PastDueInvoiceNotification:
		n, ok = result.(*notifications.PastDueInvoiceNotification)
		if !ok {
			t.Fatalf("unable to reflect interface")
		}
		break
	default:
		t.Fatalf("unexpected result type: %T", result)
	}

	if !reflect.DeepEqual(n, &notifications.PastDueInvoiceNotification{
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
