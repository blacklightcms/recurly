package webhooks

import (
	"encoding/xml"

	"github.com/blacklightcms/recurly"
)

// Charge invoice notifications.
// https://dev.recurly.com/page/webhooks#charge-invoice-notifications
const (
	NewChargeInvoice        = "new_charge_invoice_notification"
	ProcessingChargeInvoice = "processing_charge_invoice_notification"
	PastDueChargeInvoice    = "past_due_charge_invoice_notification"
	PaidChargeInvoice       = "paid_charge_invoice_notification"
	FailedChargeInvoice     = "failed_charge_invoice_notification"
	ReopenedChargeInvoice   = "reopened_charge_invoice_notification"
)

// ChargeInvoiceNotification is returned for all charge invoice notifications.
type ChargeInvoiceNotification struct {
	Type    string        `xml:"-"`
	Account Account       `xml:"account"`
	Invoice ChargeInvoice `xml:"invoice"`
}

// ChargeInvoice represents the charge invoice object sent in webhooks.
type ChargeInvoice struct {
	XMLName             xml.Name         `xml:"invoice,omitempty"`
	SubscriptionUUIDs   []string         `xml:"subscription_ids>subscription_id,omitempty"`
	UUID                string           `xml:"uuid,omitempty"`
	State               string           `xml:"state,omitempty"`
	Origin              string           `xml:"origin,omitempty"`
	InvoiceNumberPrefix string           `xml:"invoice_number_prefix,omitempty"`
	InvoiceNumber       int              `xml:"invoice_number,omitempty"`
	PONumber            string           `xml:"po_number,omitempty"`
	VATNumber           string           `xml:"vat_number,omitempty"`
	BalanceInCents      int              `xml:"balance_in_cents,omitempty"`
	TotalInCents        int              `xml:"total_in_cents,omitempty"`
	Currency            string           `xml:"currency,omitempty"`
	CreatedAt           recurly.NullTime `xml:"created_at,omitempty"`
	ClosedAt            recurly.NullTime `xml:"closed_at,omitempty"`
	NetTerms            recurly.NullInt  `xml:"net_terms,omitempty"`
	CollectionMethod    string           `xml:"collection_method,omitempty"`
}
