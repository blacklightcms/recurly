package webhooks

import (
	"encoding/xml"

	"github.com/blacklightcms/recurly"
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
type CreditInvoice struct {
	XMLName             xml.Name         `xml:"invoice,omitempty"`
	SubscriptionUUIDs   []string         `xml:"subscription_ids>subscription_id,omitempty"`
	UUID                string           `xml:"uuid,omitempty"`
	State               string           `xml:"state,omitempty"`
	Origin              string           `xml:"origin,omitempty"`
	InvoiceNumberPrefix string           `xml:"invoice_number_prefix,omitempty"`
	InvoiceNumber       int              `xml:"invoice_number,omitempty"`
	BalanceInCents      int              `xml:"balance_in_cents,omitempty"`
	TotalInCents        int              `xml:"total_in_cents,omitempty"`
	Currency            string           `xml:"currency,omitempty"`
	CreatedAt           recurly.NullTime `xml:"created_at,omitempty"`
	ClosedAt            recurly.NullTime `xml:"closed_at,omitempty"`
}
