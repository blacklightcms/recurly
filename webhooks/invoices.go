package webhooks

import (
	"encoding/xml"

	"github.com/launchpadcentral/recurly"
)

// Invoice notifications.
// Will be deprecated after credit invoices feature is turned on.
// https://dev.recurly.com/page/webhooks#invoice-notifications
const (
	NewInvoice        = "new_invoice_notification"
	PastDueInvoice    = "past_due_invoice_notification"
	ProcessingInvoice = "processing_invoice_notification"
	ClosedInvoice     = "closed_invoice_notification"
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
	XMLName             xml.Name         `xml:"invoice"`
	SubscriptionUUID    string           `xml:"subscription_id"`
	UUID                string           `xml:"uuid"`
	State               string           `xml:"state"`
	InvoiceNumberPrefix string           `xml:"invoice_number_prefix"`
	InvoiceNumber       int              `xml:"invoice_number"`
	PONumber            string           `xml:"po_number"`
	VATNumber           string           `xml:"vat_number"`
	TotalInCents        int              `xml:"total_in_cents"`
	Currency            string           `xml:"currency"`
	CreatedAt           recurly.NullTime `xml:"date"`
	ClosedAt            recurly.NullTime `xml:"closed_at"`
	NetTerms            recurly.NullInt  `xml:"net_terms"`
	CollectionMethod    string           `xml:"collection_method"`
}
