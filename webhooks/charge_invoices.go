package webhooks

import (
	"encoding/xml"

	"github.com/autopilot3/recurly"
)

// Charge invoice notifications.
// https://dev.recurly.com/page/webhooks#charge-invoice-notifications
const (
	NewChargeInvoice                 = "new_charge_invoice_notification"
	ProcessingChargeInvoice          = "processing_charge_invoice_notification"
	PastDueInvoice                   = "past_due_invoice_notification"
	PastDueChargeInvoice             = "past_due_charge_invoice_notification"
	PaidChargeInvoice                = "paid_charge_invoice_notification"
	FailedChargeInvoice              = "failed_charge_invoice_notification"
	ReopenedChargeInvoice            = "reopened_charge_invoice_notification"
	ClosedInvoiceNotification        = "closed_invoice_notification"
	UpdatedInvoiceNotification       = "updated_invoice_notification"
	UpdatedChargeInvoiceNotification = "updated_charge_invoice_notification"
	PendingInvoiceNotification       = "pending_invoice_notification"
)

// ChargeInvoiceNotification is returned for all charge invoice notifications.
type ChargeInvoiceNotification struct {
	Type    string        `xml:"-"`
	Account Account       `xml:"account"`
	Invoice ChargeInvoice `xml:"invoice"`
}

// ChargeInvoice represents the charge invoice object sent in webhooks.
// https://dev.recurly.com/page/webhooks#section-charge-invoice-schema
type ChargeInvoice struct {
	XMLName                       xml.Name         `xml:"invoice"`
	SubscriptionUUIDs             []string         `xml:"subscription_ids>subscription_id"`
	UUID                          string           `xml:"uuid"`
	State                         string           `xml:"state"`
	Origin                        string           `xml:"origin"`
	InvoiceNumberPrefix           string           `xml:"invoice_number_prefix"`
	InvoiceNumber                 int              `xml:"invoice_number"`
	PONumber                      string           `xml:"po_number"`
	VATNumber                     string           `xml:"vat_number"`
	BalanceInCents                int              `xml:"balance_in_cents"`
	TotalInCents                  int              `xml:"total_in_cents"`
	TaxInCents                    int              `xml:"tax_in_cents"`
	DiscountInCents               int              `xml:"discount_in_cents"`
	SubtotalInCents               int              `xml:"subtotal_in_cents"`
	SubTotalBeforeDiscountInCents int              `xml:"subtotal_before_discount_in_cents"`
	Currency                      string           `xml:"currency"`
	CreatedAt                     recurly.NullTime `xml:"created_at"`
	UpdatedAt                     recurly.NullTime `xml:"updated_at"`
	ClosedAt                      recurly.NullTime `xml:"closed_at"`
	DueOn                         recurly.NullTime `xml:"due_on"`
	NetTerms                      recurly.NullInt  `xml:"net_terms"`
	CollectionMethod              string           `xml:"collection_method"`
	CustomerNotes                 string           `xml:"customer_notes"`
	TermsAndConditions            string           `xml:"terms_and_conditions"`
}
