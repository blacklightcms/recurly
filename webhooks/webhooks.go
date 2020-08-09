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
// If r implements the io.Closer interface, it will be closed.
//
// NOTE: It is important to validate the source of the webhook before trusting
// it came from Recurly. Please see Recurly's documentation about the IP
// addresses to expect and/or setting up HTTP Basic Authentication to verify
// the request came from Recurly's servers.
//
// https://docs.recurly.com/docs/webhooks
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

	dst, err := nameToNotification(n.XMLName.Local)
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
	case NewSubscription, UpdatedSubscription, RenewedSubscription, ExpiredSubscription, CanceledSubscription, ReactivatedAccount, PausedSubscription,
		ResumedSubscription, ScheduledPauseSubscription, ModifiedPauseSubscription, PausedRenewalSubscription, PauseCanceledSubscription:
		return &SubscriptionNotification{Type: name}, nil
	case NewChargeInvoice, ProcessingChargeInvoice, PastDueChargeInvoice, PaidChargeInvoice, FailedChargeInvoice, ReopenedChargeInvoice:
		return &ChargeInvoiceNotification{Type: name}, nil
	case NewCreditInvoice, ProcessingCreditInvoice, ClosedCreditInvoice, VoidedCreditInvoice, ReopenedCreditInvoice, OpenCreditInvoice:
		return &CreditInvoiceNotification{Type: name}, nil
	case NewCreditPayment, VoidedCreditPayment:
		return &CreditPaymentNotification{Type: name}, nil
	case SuccessfulPayment, FailedPayment, VoidPayment, SuccessfulRefund, ScheduledPayment, ProcessingPayment, TransactionStatusUpdated:
		return &PaymentNotification{Type: name}, nil
	case NewDunningEvent:
		return &NewDunningEventNotification{Type: name}, nil
	}
	return nil, ErrUnknownNotification{name: name}
}

// ErrUnknownNotification is returned when the incoming webhook does not match a
// predefined notification type. It implements the error interface.
type ErrUnknownNotification struct {
	name string
}

// Error implements the error interface.
func (e ErrUnknownNotification) Error() string {
	return fmt.Sprintf("unknown notification: %q", e.name)
}

// Name returns the name of the unknown notification.
func (e ErrUnknownNotification) Name() string {
	return e.name
}
