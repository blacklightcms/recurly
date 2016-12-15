package notifications

import (
	"encoding/xml"
	"io"
	"io/ioutil"

	"github.com/blacklightcms/go-recurly/recurly"
)

// Webhook notification constants.
const (
	SuccessfulPayment = "successful_payment_notification"
	FailedPayment     = "failed_payment_notification"
	PastDueInvoice    = "past_due_invoice_notification"
)

type notificationName struct {
	XMLName xml.Name
}

type (
	// SuccessfulPaymentNotification is sent when a payment is successful.
	SuccessfulPaymentNotification struct {
		Account     recurly.Account
		Transaction recurly.Transaction
	}

	// FailedPaymentNotification is sent when a payment fails.
	FailedPaymentNotification struct {
		Account     recurly.Account
		Transaction recurly.Transaction
	}

	// PastDueInvoiceNotification is sent when an invoice is past due.
	PastDueInvoiceNotification struct {
		Account recurly.Account
		Invoice recurly.Invoice
	}
)

// TODO respond to webhook with a 200 code within 5 minutes or Recurly will resend it.
// TODO Caller can route the notification with a switch statement on XMLName. Then caller
// will make an API call to Lookup Invoice (for which it only needs the invoice number) to
// verify that the invoice is either paid or still unpaid, then it can progress forward.

// Parse parses an incoming webhook and returns the notification.
func Parse(r io.Reader) (interface{}, error) {
	if closer, ok := r.(io.Closer); ok {
		defer closer.Close()
	}

	notification, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var n notificationName
	if err := xml.Unmarshal(notification, &n); err != nil {
		return nil, err
	}

	var dst interface{}
	switch n.XMLName.Local {
	case SuccessfulPayment:
		dst = &SuccessfulPaymentNotification{}
	case FailedPayment:
		dst = &FailedPaymentNotification{}
	case PastDueInvoice:
		dst = &PastDueInvoiceNotification{}
	}

	if err := xml.Unmarshal(notification, dst); err != nil {
		return nil, err
	}

	return dst, nil
}
