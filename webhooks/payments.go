package webhooks

import (
	"encoding/xml"

	"github.com/blacklightcms/recurly"
)

// Payment notifications.
// https://dev.recurly.com/page/webhooks#payment-notifications
const (
	SuccessfulPayment = "successful_payment_notification"
	FailedPayment     = "failed_payment_notification"
	VoidPayment       = "void_payment_notification"
	SuccessfulRefund  = "successful_refund_notification"
)

// PaymentNotification is returned for all credit payment notifications.
type PaymentNotification struct {
	Type        string      `xml:"-"`
	Account     Account     `xml:"account"`
	Transaction Transaction `xml:"transaction"`
}

// Transaction represents the transaction object sent in webhooks.
type Transaction struct {
	XMLName           xml.Name         `xml:"transaction"`
	UUID              string           `xml:"id,omitempty"`
	InvoiceNumber     int              `xml:"invoice_number,omitempty"`
	SubscriptionUUID  string           `xml:"subscription_id,omitempty"`
	Action            string           `xml:"action,omitempty"`
	AmountInCents     int              `xml:"amount_in_cents,omitempty"`
	Status            string           `xml:"status,omitempty"`
	Message           string           `xml:"message,omitempty"`
	GatewayErrorCodes string           `xml:"gateway_error_codes,omitempty"`
	FailureType       string           `xml:"failure_type,omitempty"`
	Reference         string           `xml:"reference,omitempty"`
	Source            string           `xml:"source,omitempty"`
	Test              recurly.NullBool `xml:"test,omitempty"`
	Voidable          recurly.NullBool `xml:"voidable,omitempty"`
	Refundable        recurly.NullBool `xml:"refundable,omitempty"`
}

// Transaction constants.
const (
	TransactionFailureTypeDeclined  = "declined"
	TransactionFailureTypeDuplicate = "duplicate_transaction"
)
