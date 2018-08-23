package webhooks

import (
	"encoding/xml"

	"github.com/launchpadcentral/recurly"
)

// Credit payment notifications.
// https://dev.recurly.com/page/webhooks#credit-payment-notifications
const (
	NewCreditPayment    = "new_credit_payment_notification"
	VoidedCreditPayment = "voided_credit_payment_notification"
)

// CreditPaymentNotification is returned for all credit payment notifications.
type CreditPaymentNotification struct {
	Type          string        `xml:"-"`
	Account       Account       `xml:"account"`
	CreditPayment CreditPayment `xml:"credit_payment"`
}

// CreditPayment represents the credit payment object sent in webhooks.
type CreditPayment struct {
	XMLName                   xml.Name         `xml:"credit_payment"`
	UUID                      string           `xml:"uuid"`
	Action                    string           `xml:"action"`
	AmountInCents             int              `xml:"amount_in_cents"`
	OriginalInvoiceNumber     int              `xml:"original_invoice_number"`
	AppliedToInvoiceNumber    int              `xml:"applied_to_invoice_number"`
	OriginalCreditPaymentUUID string           `xml:"original_credit_payment_uuid"`
	RefundTransactionUUID     string           `xml:"refund_transaction_uuid"`
	CreatedAt                 recurly.NullTime `xml:"created_at"`
	VoidedAt                  recurly.NullTime `xml:"voided_at"`
}
