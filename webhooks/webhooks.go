package webhooks

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/blacklightcms/recurly"
)

// Webhook notification constants.
const (
	// Subscription notifications.
	ExpiredSubscription = "expired_subscription_notification"

	// Invoice notifications.
	NewInvoice     = "new_invoice_notification"
	PastDueInvoice = "past_due_invoice_notification"

	// Payment notifications.
	SuccessfulPayment = "successful_payment_notification"
	FailedPayment     = "failed_payment_notification"
)

type notificationName struct {
	XMLName xml.Name
}

// Subscription types.
type (
	// ExpiredSubscriptionNotification is sent when a subscription is no longer valid.
	// https://dev.recurly.com/v2.4/page/webhooks#section-expired-subscription
	ExpiredSubscriptionNotification struct {
		Account      recurly.Account      `xml:"account"`
		Subscription recurly.Subscription `xml:"subscription"`
	}
)

// Invoice types.
type (
	// NewInvoiceNotification is sent when an invoice generated.
	// https://dev.recurly.com/page/webhooks#section-new-invoice
	NewInvoiceNotification struct {
		Account recurly.Account `xml:"account"`
		Invoice recurly.Invoice `xml:"invoice"`
	}

	// PastDueInvoiceNotification is sent when an invoice is past due.
	// https://dev.recurly.com/v2.4/page/webhooks#section-past-due-invoice
	PastDueInvoiceNotification struct {
		Account recurly.Account `xml:"account"`
		Invoice recurly.Invoice `xml:"invoice"`
	}
)

// Payment types.
type (
	// SuccessfulPaymentNotification is sent when a payment is successful.
	// https://dev.recurly.com/v2.4/page/webhooks#section-successful-payment
	SuccessfulPaymentNotification struct {
		Account     recurly.Account     `xml:"account"`
		Transaction recurly.Transaction `xml:"transaction"`
	}

	// FailedPaymentNotification is sent when a payment fails.
	// https://dev.recurly.com/v2.4/page/webhooks#section-failed-payment
	FailedPaymentNotification struct {
		Account     recurly.Account     `xml:"account"`
		Transaction recurly.Transaction `xml:"transaction"`
	}
)

// transactionHolder allows the uuid and invoice number fields to be set.
// The UUID field is labeled id in notifications and the invoice number
// is not included on the existing transaction struct.
type transactionHolder interface {
	setTransactionFields(id string, in int)
}

// setTransactionFields sets fields on the transaction struct.
func (n *SuccessfulPaymentNotification) setTransactionFields(id string, invoiceNumber int) {
	n.Transaction.UUID = id
	n.Transaction.InvoiceNumber = invoiceNumber
}

func (n *FailedPaymentNotification) setTransactionFields(id string, invoiceNumber int) {
	n.Transaction.UUID = id
	n.Transaction.InvoiceNumber = invoiceNumber
}

// transaction allows the transaction id and invoice number to be unmarshalled
// so they can be set on the notification struct.
type transaction struct {
	ID            string `xml:"transaction>id"`
	InvoiceNumber int    `xml:"transaction>invoice_number,omitempty"`
}

// ErrUnknownNotification is used when the incoming webhook does not match a
// predefined notification type. It implements the error interface.
type ErrUnknownNotification struct {
	name string
}

// Error implements the error interface.
func (e ErrUnknownNotification) Error() string {
	return fmt.Sprintf("unknown notification: %s", e.name)
}

// Name returns the name of the unknown notification.
func (e ErrUnknownNotification) Name() string {
	return e.name
}

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
	case ExpiredSubscription:
		dst = &ExpiredSubscriptionNotification{}
	case NewInvoice:
		dst = &NewInvoiceNotification{}
	case PastDueInvoice:
		dst = &PastDueInvoiceNotification{}
	case SuccessfulPayment:
		dst = &SuccessfulPaymentNotification{}
	case FailedPayment:
		dst = &FailedPaymentNotification{}
	default:
		return nil, ErrUnknownNotification{name: n.XMLName.Local}
	}

	if err := xml.Unmarshal(notification, dst); err != nil {
		return nil, err
	}

	if th, ok := dst.(transactionHolder); ok {
		var t transaction
		if err := xml.Unmarshal(notification, &t); err != nil {
			return nil, err
		}
		th.setTransactionFields(t.ID, t.InvoiceNumber)
	}

	return dst, nil
}
