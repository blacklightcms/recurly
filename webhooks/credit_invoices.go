package webhooks

import (
	"encoding/xml"

	"github.com/splice/recurly"
)

// Credit invoice notifications.
// https://dev.recurly.com/page/webhooks#charge-invoice-notifications
const (
	NewCreditInvoice        = "new_credit_invoice_notification"
	ProcessingCreditInvoice = "processing_credit_invoice_notification"
	ClosedCreditInvoice     = "closed_credit_invoice_notification"
	VoidedCreditInvoice     = "voided_credit_invoice_notification"
	ReopenedCreditInvoice   = "reopened_credit_invoice_notification"
	OpenCreditInvoice       = "open_credit_invoice_notification"
)

// CreditInvoiceNotification is returned for all credit invoice notifications.
type CreditInvoiceNotification struct {
	Type    string        `xml:"-"`
	Account Account       `xml:"account"`
	Invoice CreditInvoice `xml:"invoice"`
}

// CreditInvoice represents the credit invoice object sent in webhooks.
// https://dev.recurly.com/page/webhooks#section-credit-invoice-schema
type CreditInvoice struct {
	XMLName                       xml.Name         `xml:"invoice"`
	SubscriptionUUIDs             []string         `xml:"subscription_ids>subscription_id"`
	UUID                          string           `xml:"uuid"`
	State                         string           `xml:"state"`
	Origin                        string           `xml:"origin"`
	InvoiceNumberPrefix           string           `xml:"invoice_number_prefix"`
	InvoiceNumber                 int              `xml:"invoice_number"`
	VATNumber                     string           `xml:"vat_number"`
	BalanceInCents                int              `xml:"balance_in_cents"`
	TotalInCents                  int              `xml:"total_in_cents"`
	TaxInCents                    int              `xml:"tax_in_cents"`
	SubtotalInCents               int              `xml:"subtotal_in_cents"`
	SubTotalBeforeDiscountInCents int              `xml:"subtotal_before_discount_in_cents"`
	DiscountInCents               int              `xml:"discount_in_cents"`
	Currency                      string           `xml:"currency"`
	CreatedAt                     recurly.NullTime `xml:"created_at"`
	UpdatedAt                     recurly.NullTime `xml:"updated_at"`
	ClosedAt                      recurly.NullTime `xml:"closed_at"`
	CustomerNotes                 string           `xml:"customer_notes"`
}
