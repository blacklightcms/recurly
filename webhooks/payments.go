package webhooks

import (
	"encoding/xml"

	"github.com/launchpadcentral/recurly"
)

// Payment notifications.
// https://dev.recurly.com/page/webhooks#payment-notifications
const (
	SuccessfulPayment = "successful_payment_notification"
	FailedPayment     = "failed_payment_notification"
	VoidPayment       = "void_payment_notification"
	SuccessfulRefund  = "successful_refund_notification"
	ScheduledPayment  = "scheduled_payment_notification"
	ProcessingPayment = "processing_payment_notification"
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
	UUID              string           `xml:"id"`
	InvoiceNumber     int              `xml:"invoice_number"`
	SubscriptionUUID  string           `xml:"subscription_id"`
	Action            string           `xml:"action"`
	PaymentMethod     string           `xml:"payment_method"`
	AmountInCents     int              `xml:"amount_in_cents"`
	Status            string           `xml:"status"`
	Message           string           `xml:"message"`
	GatewayErrorCodes string           `xml:"gateway_error_codes"`
	FailureType       string           `xml:"failure_type"`
	Reference         string           `xml:"reference"`
	Source            string           `xml:"source"`
	Test              recurly.NullBool `xml:"test"`
	Voidable          recurly.NullBool `xml:"voidable"`
	Refundable        recurly.NullBool `xml:"refundable"`
}

// Transaction constants.
const (
	TransactionFailureTypeDeclined  = "declined"
	TransactionFailureTypeDuplicate = "duplicate_transaction"
)
