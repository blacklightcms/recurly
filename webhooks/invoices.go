package webhooks

import (
	"encoding/xml"

	"github.com/blacklightcms/recurly"
)

// Invoice notifications.
// Will be deprecated after credit invoices feature is turned on.
// https://dev.recurly.com/page/webhooks#invoice-notifications
const (
	NewInvoice     = "new_invoice_notification"
	PastDueInvoice = "past_due_invoice_notification"
)

// InvoiceNotification is returned for all invoice notifications.
type InvoiceNotification struct {
	Type    string  `xml:"-"`
	Account Account `xml:"account"`
	Invoice Invoice `xml:"invoice"`
}

// Invoice represents the invoice object sent in webhooks.
// After credit invoices have been turned on, these notifications will only be
// sent for TypeLegacy invoices (posted before the feature was turned on)
// then deprecated.
type Invoice struct {
	XMLName             xml.Name         `xml:"invoice,omitempty"`
	SubscriptionUUID    string           `xml:"subscription_id,omitempty"`
	UUID                string           `xml:"uuid,omitempty"`
	State               string           `xml:"state,omitempty"`
	InvoiceNumberPrefix string           `xml:"invoice_number_prefix,omitempty"`
	InvoiceNumber       int              `xml:"invoice_number,omitempty"`
	PONumber            string           `xml:"po_number,omitempty"`
	VATNumber           string           `xml:"vat_number,omitempty"`
	TotalInCents        int              `xml:"total_in_cents,omitempty"`
	Currency            string           `xml:"currency,omitempty"`
	CreatedAt           recurly.NullTime `xml:"date,omitempty"`
	ClosedAt            recurly.NullTime `xml:"closed_at,omitempty"`
	NetTerms            recurly.NullInt  `xml:"net_terms,omitempty"`
	CollectionMethod    string           `xml:"collection_method,omitempty"`
}
