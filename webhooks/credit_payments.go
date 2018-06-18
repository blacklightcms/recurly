package webhooks

import (
	"encoding/xml"

	"github.com/blacklightcms/recurly"
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
	XMLName                   xml.Name         `xml:"credit_payment,omitempty"`
	UUID                      string           `xml:"uuid,omitempty"`
	Action                    string           `xml:"action,omitempty"`
	AmountInCents             int              `xml:"amount_in_cents,omitempty"`
	OriginalInvoiceNumber     int              `xml:"original_nvoice_number,omitempty"`
	AppliedToInvoiceNumber    int              `xml:"applied_to_invoice_number,omitempty"`
	OriginalCreditPaymentUUID string           `xml:"original_credit_payment_uuid,omitempty"`
	RefundTransactionUUID     string           `xml:"refund_transaction_uuid,omitempty"`
	CreatedAt                 recurly.NullTime `xml:"created_at,omitempty"`
	VoidedAt                  recurly.NullTime `xml:"voided_at,omitempty"`
}
