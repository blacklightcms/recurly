package webhooks

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
)

type notificationName struct {
	XMLName xml.Name
}

// Parse parses an incoming webhook and returns the notification.
func Parse(r io.Reader) (interface{}, error) {
	return parse(r, dst)
}

// ParseOriginal parses an incoming webhook and returns the notification
// for previous webhook notifications.
// Will be deprecated after credit invoices feature is turned on.
func ParseOriginal(r io.Reader) (interface{}, error) {
	return parse(r, origdst)
}

// parse parses an incoming webhook and returns the notification.
func parse(r io.Reader, fn func(s string) (interface{}, error)) (interface{}, error) {
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

	dst, err := fn(n.XMLName.Local)
	if err != nil {
		return nil, err
	}

	if err := xml.Unmarshal(notification, dst); err != nil {
		return nil, err
	}

	return dst, nil
}

// dst returns the notification interface.
func dst(s string) (dst interface{}, err error) {
	switch s {
	case BillingInfoUpdated:
		dst = &AccountNotification{Type: s}
	case NewSubscription, UpdatedSubscription, RenewedSubscription, ExpiredSubscription, CanceledSubscription:
		dst = &SubscriptionNotification{Type: s}
	case NewChargeInvoice, ProcessingChargeInvoice, PastDueChargeInvoice, PaidChargeInvoice, FailedChargeInvoice, ReopenedChargeInvoice:
		dst = &ChargeInvoiceNotification{Type: s}
	case NewCreditInvoice, ProcessingCreditInvoice, ClosedCreditInvoice, VoidedCreditInvoice, ReopenedCreditInvoice, OpenCreditInvoice:
		dst = &CreditInvoiceNotification{Type: s}
	case NewCreditPayment, VoidedCreditPayment:
		dst = &CreditPaymentNotification{Type: s}
	case NewInvoice, PastDueInvoice:
		dst = &InvoiceNotification{Type: s}
	case SuccessfulPayment, FailedPayment, VoidPayment, SuccessfulRefund:
		dst = &PaymentNotification{Type: s}
	default:
		return nil, ErrUnknownNotification{name: s}
	}
	return dst, nil
}

// origdst returns interfaces for webhooks prior to
// the credit invoices feature.
func origdst(s string) (dst interface{}, err error) {
	switch s {
	case BillingInfoUpdated:
		dst = &AccountNotification{Type: s}
	case NewSubscription, UpdatedSubscription, RenewedSubscription, ExpiredSubscription, CanceledSubscription:
		dst = &SubscriptionNotification{Type: s}
	case NewInvoice, PastDueInvoice:
		dst = &InvoiceNotification{Type: s}
	case SuccessfulPayment, FailedPayment, VoidPayment, SuccessfulRefund:
		dst = &PaymentNotification{Type: s}
	default:
		return nil, ErrUnknownNotification{name: s}
	}
	return dst, nil
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
