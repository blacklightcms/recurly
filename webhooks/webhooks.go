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
	return parse(r, nameToNotification)
}

// ParseDeprecated parses an incoming webhook and returns the notification
// for previous webhook notifications.
// Will be deprecated after credit invoices feature is turned on.
func ParseDeprecated(r io.Reader) (interface{}, error) {
	return parse(r, nameToNotificationDeprecated)
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

// nameToNotification returns the notification interface.
func nameToNotification(name string) (interface{}, error) {
	switch name {
	case BillingInfoUpdated, NewAccount, UpdatedAccount, CanceledAccount, BillingInfoUpdateFailed:
		return &AccountNotification{Type: name}, nil
	case NewSubscription, UpdatedSubscription, RenewedSubscription, ExpiredSubscription, CanceledSubscription, ReactivatedAccount:
		return &SubscriptionNotification{Type: name}, nil
	case NewChargeInvoice, ProcessingChargeInvoice, PastDueChargeInvoice, PaidChargeInvoice, FailedChargeInvoice, ReopenedChargeInvoice:
		return &ChargeInvoiceNotification{Type: name}, nil
	case NewCreditInvoice, ProcessingCreditInvoice, ClosedCreditInvoice, VoidedCreditInvoice, ReopenedCreditInvoice, OpenCreditInvoice:
		return &CreditInvoiceNotification{Type: name}, nil
	case NewCreditPayment, VoidedCreditPayment:
		return &CreditPaymentNotification{Type: name}, nil
	case NewInvoice, PastDueInvoice, ProcessingInvoice, ClosedInvoice:
		return &InvoiceNotification{Type: name}, nil
	case SuccessfulPayment, FailedPayment, VoidPayment, SuccessfulRefund, ScheduledPayment, ProcessingPayment:
		return &PaymentNotification{Type: name}, nil
	case NewDunningEvent:
		return &NewDunningEventNotification{Type: name}, nil
	}
	return nil, ErrUnknownNotification{name: name}
}

// nameToNotificationDeprecated returns interfaces for webhooks prior to
// the credit invoices feature.
func nameToNotificationDeprecated(name string) (interface{}, error) {
	switch name {
	case BillingInfoUpdated, NewAccount, UpdatedAccount, CanceledAccount, BillingInfoUpdateFailed:
		return &AccountNotification{Type: name}, nil
	case NewSubscription, UpdatedSubscription, RenewedSubscription, ExpiredSubscription, CanceledSubscription, ReactivatedAccount:
		return &SubscriptionNotification{Type: name}, nil
	case NewInvoice, PastDueInvoice, ProcessingInvoice, ClosedInvoice:
		return &InvoiceNotification{Type: name}, nil
	case SuccessfulPayment, FailedPayment, VoidPayment, SuccessfulRefund, ScheduledPayment, ProcessingPayment:
		return &PaymentNotification{Type: name}, nil
	case NewDunningEvent:
		return &NewDunningEventDeprecatedNotification{Type: name}, nil
	}
	return nil, ErrUnknownNotification{name: name}
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
