package recurly

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
)

type (
	// Notification is used for webhook notifications. The type method returns
	// the notification type (i.e. "new_account_notification")
	Notification interface {
		Type() string
	}

	webhook struct {
		XMLName xml.Name
	}

	// AccountNotification ...
	AccountNotification struct {
		Account
		webhook
	}

	// InvoiceNotification ...
	InvoiceNotification struct {
		webhook
		Account Account `xml:"account"`
		Invoice Invoice `xml:"invoice"`
	}

	// SubscriptionNotification ...
	SubscriptionNotification struct {
		webhook
		Account      Account      `xml:"account"`
		Subscription Subscription `xml:"subscription"`
	}

	// PaymentNotification ...
	PaymentNotification struct {
		webhook
		Account     Account     `xml:"account"`
		Transaction Transaction `xml:"transaction"`
	}
)

// Type returns the webhook notification type.
func (w webhook) Type() string {
	return w.XMLName.Local
}

var (
	webhooks = []string{
		"new_account_notification",
		"canceled_account_notification",
		"billing_info_updated_notification",
		"reactivated_account_notification",
		"new_invoice_notification",
		"processing_invoice_notification",
		"closed_invoice_notification",
		"past_due_invoice_notification",
		"new_subscription_notification",
		"updated_subscription_notification",
		"canceled_subscription_notification",
		"expired_subscription_notification",
		"renewed_subscription_notification",
		"scheduled_payment_notification",
		"processing_payment_notification",
		"successful_payment_notification",
		"failed_payment_notification",
		"successful_refund_notification",
		"void_payment_notification",
	}
)

// HandleWebhook ...
func HandleWebhook(body io.Reader) (Notification, error) {
	n, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	for _, name := range webhooks {
		if !bytes.Contains(n, []byte("<"+name+">")) {
			continue
		}

		var dest Notification
		switch name {
		case "new_account_notification", "canceled_account_notification", "billing_info_updated_notification":
			dest = AccountNotification{}
		case "new_invoice_notification", "processing_invoice_notification", "closed_invoice_notification", "past_due_invoice_notification":
			var dest InvoiceNotification
			b := bytes.NewBuffer(n)
			xml.NewDecoder(b).Decode(&dest)

			return dest, nil
		case "reactivated_account_notification", "new_subscription_notification", "updated_subscription_notification", "canceled_subscription_notification", "expired_subscription_notification", "renewed_subscription_notification":
			dest = SubscriptionNotification{}
		case "scheduled_payment_notification", "processing_payment_notification", "successful_payment_notification", "failed_payment_notification", "successful_refund_notification", "void_payment_notification":
			dest = PaymentNotification{}
		default:
			return nil, errors.New("Unknown notification")
		}

		err := xml.Unmarshal(n, &dest)
		if err != nil {
			return nil, err
		}

		return dest, nil
	}

	return nil, errors.New("Notification not found")
}
