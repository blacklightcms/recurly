package recurly

import "encoding/xml"

const (
	// InvoiceStateOpen is an invoice state for invoices that are open, pending
	// collection.
	InvoiceStateOpen = "open"

	// InvoiceStateCollected is an invoice state for invoices that have been
	// successfully collected.
	InvoiceStateCollected = "collected"

	// InvoiceStateFailed is an invoice state for invoices that failed to collect.
	InvoiceStateFailed = "failed"

	// InvoiceStatePastDue is an invoice state for invoices where initial collection
	// failed, but Recurly is still attempting collection.
	InvoiceStatePastDue = "past_due"
)

// Invoice is an individual invoice for an account.
// The only fields annotated with XML tags are those for posting an invoice.
// Unmarshaling an invoice is handled by the custom UnmarshalXML function.
type Invoice struct {
	XMLName               xml.Name      `xml:"invoice,omitempty"`
	AccountCode           string        `xml:"-"`
	Address               Address       `xml:"-"`
	SubscriptionUUID      string        `xml:"-"`
	OriginalInvoiceNumber int           `xml:"-"`
	UUID                  string        `xml:"-"`
	State                 string        `xml:"-"`
	InvoiceNumberPrefix   string        `xml:"-"`
	InvoiceNumber         int           `xml:"-"`
	PONumber              string        `xml:"po_number,omitempty"` // PostInvoice param
	VATNumber             string        `xml:"-"`
	SubtotalInCents       int           `xml:"-"`
	TaxInCents            int           `xml:"-"`
	TotalInCents          int           `xml:"-"`
	Currency              string        `xml:"-"`
	CreatedAt             NullTime      `xml:"-"`
	ClosedAt              NullTime      `xml:"-"`
	TaxType               string        `xml:"-"`
	TaxRegion             string        `xml:"-"`
	TaxRate               float64       `xml:"-"`
	NetTerms              NullInt       `xml:"net_terms,omitempty"`                // PostInvoice param
	CollectionMethod      string        `xml:"collection_method,omitempty"`        // PostInvoice param
	TermsAndConditions    string        `xml:"terms_and_conditions,omitempty"`     // PostInvoice param
	CustomerNotes         string        `xml:"customer_notes,omitempty"`           // PostInvoice param
	VatReverseChargeNotes string        `xml:"vat_reverse_charge_notes,omitempty"` // PostInvoice param
	LineItems             []Adjustment  `xml:"-"`
	Transactions          []Transaction `xml:"-"`
}

// UnmarshalXML unmarshals invoices and handles intermediary state during unmarshaling
// for types like href.
func (i *Invoice) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v struct {
		XMLName               xml.Name      `xml:"invoice,omitempty"`
		AccountCode           hrefString    `xml:"account,omitempty"` // Read only
		Address               Address       `xml:"address,omitempty"`
		SubscriptionUUID      hrefString    `xml:"subscription,omitempty"`
		OriginalInvoiceNumber hrefInt       `xml:"original_invoice,omitempty"` // Read only
		UUID                  string        `xml:"uuid,omitempty"`
		State                 string        `xml:"state,omitempty"`
		InvoiceNumberPrefix   string        `xml:"invoice_number_prefix,omitempty"`
		InvoiceNumber         int           `xml:"invoice_number,omitempty"`
		PONumber              string        `xml:"po_number,omitempty"`
		VATNumber             string        `xml:"vat_number,omitempty"`
		SubtotalInCents       int           `xml:"subtotal_in_cents,omitempty"`
		TaxInCents            int           `xml:"tax_in_cents,omitempty"`
		TotalInCents          int           `xml:"total_in_cents,omitempty"`
		Currency              string        `xml:"currency,omitempty"`
		CreatedAt             NullTime      `xml:"created_at,omitempty"`
		ClosedAt              NullTime      `xml:"closed_at,omitempty"`
		TaxType               string        `xml:"tax_type,omitempty"`
		TaxRegion             string        `xml:"tax_region,omitempty"`
		TaxRate               float64       `xml:"tax_rate,omitempty"`
		NetTerms              NullInt       `xml:"net_terms,omitempty"`
		CollectionMethod      string        `xml:"collection_method,omitempty"`
		LineItems             []Adjustment  `xml:"line_items>adjustment,omitempty"`
		Transactions          []Transaction `xml:"transactions>transaction,omitempty"`
	}
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	*i = Invoice{
		XMLName:               v.XMLName,
		AccountCode:           string(v.AccountCode),
		Address:               v.Address,
		SubscriptionUUID:      string(v.SubscriptionUUID),
		OriginalInvoiceNumber: int(v.OriginalInvoiceNumber),
		UUID:                v.UUID,
		State:               v.State,
		InvoiceNumberPrefix: v.InvoiceNumberPrefix,
		InvoiceNumber:       v.InvoiceNumber,
		PONumber:            v.PONumber,
		VATNumber:           v.VATNumber,
		SubtotalInCents:     v.SubtotalInCents,
		TaxInCents:          v.TaxInCents,
		TotalInCents:        v.TotalInCents,
		Currency:            v.Currency,
		CreatedAt:           v.CreatedAt,
		ClosedAt:            v.ClosedAt,
		TaxType:             v.TaxType,
		TaxRegion:           v.TaxRegion,
		TaxRate:             v.TaxRate,
		NetTerms:            v.NetTerms,
		CollectionMethod:    v.CollectionMethod,
		LineItems:           v.LineItems,
		Transactions:        v.Transactions,
	}

	return nil
}
