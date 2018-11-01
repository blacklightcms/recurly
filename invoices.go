package recurly

import (
	"encoding/xml"
	"time"
)

// Invoice state constants.
// https://docs.recurly.com/docs/credit-invoices-release#section-invoice-attribute-changes
const (
	ChargeInvoiceStatePending    = "pending"    // previously "open"
	ChargeInvoiceStateProcessing = "processing" // ACH payments only
	ChargeInvoiceStatePastDue    = "past_due"
	ChargeInvoiceStatePaid       = "paid" // previously "collected"
	ChargeInvoiceStateFailed     = "failed"

	CreditInvoiceStateOpen       = "open"
	CreditInvoiceStateProcessing = "processing" // ACH/bank refund processing
	CreditInvoiceStateClosed     = "closed"
	CreditInvoiceStateVoided     = "voided"

	// Deprecated
	InvoiceStateOpenDeprecated      = "open"
	InvoiceStateCollectedDeprecated = "collected"
)

// Collection method constants.
const (
	CollectionMethodAutomatic = "automatic" // card on file
	CollectionMethodManual    = "manual"    // external payment method
)

// Payment method constants.
const (
	PaymentMethodCreditCard   = "credit_card"
	PaymentMethodPayPal       = "paypal"
	PaymentMethodEFT          = "eft"
	PaymentMethodWireTransfer = "wire_transfer"
	PaymentMethodMoneyOrder   = "money_order"
	PaymentMethodCheck        = "check"
	PaymentMethodOther        = "other"
)

// Invoice origin constants.
const (
	ChargeInvoiceOriginPurchase        = "purchase"
	ChargeInvoiceOriginRenewal         = "renewal"
	ChargeInvoiceOriginImmediateChange = "immediate_change"
	ChargeInvoiceOriginTermination     = "termination"

	CreditInvoiceOriginGiftCard       = "gift_card"
	CreditInvoiceOriginRefund         = "refund"
	CreditInvoiceOriginCredit         = "credit"
	CreditInvoiceOriginWriteOff       = "write_off"
	CreditInvoiceOriginExternalCredit = "external_credit"
)

// Invoice type constants.
const (
	InvoiceTypeCharge = "charge"
	InvoiceTypeCredit = "credit"
	InvoiceTypeLegacy = "legacy" // all invoices prior to change have type legacy
)

// Refund constants.
const (
	VoidRefundMethodTransactionFirst = "transaction_first"
	VoidRefundMethodCreditFirst      = "credit_first"
)

// Invoice is an individual invoice for an account.
// The only fields annotated with XML tags are those for posting an invoice.
// Unmarshaling an invoice is handled by the custom UnmarshalXML function.
type Invoice struct {
	XMLName                 xml.Name        `xml:"invoice,omitempty"`
	AccountCode             string          `xml:"-"`
	Address                 Address         `xml:"-"`
	OriginalInvoiceNumber   int             `xml:"-"`
	UUID                    string          `xml:"-"`
	State                   string          `xml:"-"`
	InvoiceNumberPrefix     string          `xml:"-"`
	InvoiceNumber           int             `xml:"-"`
	PONumber                string          `xml:"po_number,omitempty"` // PostInvoice param
	VATNumber               string          `xml:"-"`
	DiscountInCents         int             `xml:"-"`
	SubtotalInCents         int             `xml:"-"`
	TaxInCents              int             `xml:"-"`
	TotalInCents            int             `xml:"-"`
	BalanceInCents          int             `xml:"-"`
	Currency                string          `xml:"-"`
	DueOn                   NullTime        `xml:"-"`
	CreatedAt               NullTime        `xml:"-"`
	UpdatedAt               NullTime        `xml:"-"`
	AttemptNextCollectionAt NullTime        `xml:"-"`
	ClosedAt                NullTime        `xml:"-"`
	Type                    string          `xml:"-"`
	Origin                  string          `xml:"-"`
	TaxType                 string          `xml:"-"`
	TaxRegion               string          `xml:"-"`
	TaxRate                 float64         `xml:"-"`
	NetTerms                NullInt         `xml:"net_terms,omitempty"`                // PostInvoice param
	CollectionMethod        string          `xml:"collection_method,omitempty"`        // PostInvoice param
	TermsAndConditions      string          `xml:"terms_and_conditions,omitempty"`     // PostInvoice param
	CustomerNotes           string          `xml:"customer_notes,omitempty"`           // PostInvoice param
	VatReverseChargeNotes   string          `xml:"vat_reverse_charge_notes,omitempty"` // PostInvoice param
	LineItems               []Adjustment    `xml:"-"`
	Transactions            []Transaction   `xml:"-"`
	CreditPayments          []CreditPayment `xml:"-"`
}

// UnmarshalXML unmarshals invoices and handles intermediary state during unmarshaling
// for types like href.
func (i *Invoice) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v struct {
		XMLName xml.Name `xml:"invoice,omitempty"`
		invoiceFields
	}
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	*i = Invoice{
		XMLName:               v.XMLName,
		AccountCode:           string(v.AccountCode),
		Address:               v.Address,
		OriginalInvoiceNumber: int(v.OriginalInvoiceNumber),
		UUID:                    v.UUID,
		State:                   v.State,
		InvoiceNumberPrefix:     v.InvoiceNumberPrefix,
		InvoiceNumber:           v.InvoiceNumber,
		PONumber:                v.PONumber,
		VATNumber:               v.VATNumber,
		DiscountInCents:         v.DiscountInCents,
		SubtotalInCents:         v.SubtotalInCents,
		TaxInCents:              v.TaxInCents,
		TotalInCents:            v.TotalInCents,
		BalanceInCents:          v.BalanceInCents,
		Currency:                v.Currency,
		DueOn:                   v.DueOn,
		CreatedAt:               v.CreatedAt,
		UpdatedAt:               v.UpdatedAt,
		AttemptNextCollectionAt: v.AttemptNextCollectionAt,
		ClosedAt:                v.ClosedAt,
		Type:                    v.Type,
		Origin:                  v.Origin,
		TaxType:                 v.TaxType,
		TaxRegion:               v.TaxRegion,
		TaxRate:                 v.TaxRate,
		NetTerms:                v.NetTerms,
		CollectionMethod:        v.CollectionMethod,
		LineItems:               v.LineItems,
		Transactions:            v.Transactions,
		CreditPayments:          v.CreditPayments,
	}

	return nil
}

// InvoiceCollection is the data type returned from Preview, Post,
// MarkFailed, and inside PreviewSubscription, and PreviewSubscriptionChange.
type InvoiceCollection struct {
	XMLName        xml.Name  `xml:"invoice_collection"`
	ChargeInvoice  *Invoice  `xml:"-"`
	CreditInvoices []Invoice `xml:"-"`
}

// UnmarshalXML unmarshals invoices and handles intermediary state during unmarshaling
// for types like href.
func (i *InvoiceCollection) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v struct {
		XMLName       xml.Name `xml:"invoice_collection"`
		ChargeInvoice struct {
			XMLName xml.Name `xml:"charge_invoice,omitempty"`
			invoiceFields
		} `xml:"charge_invoice,omitempty"`
		CreditInvoices []struct {
			XMLName xml.Name `xml:"credit_invoice,omitempty"`
			invoiceFields
		} `xml:"credit_invoices>credit_invoice,omitempty"`
	}
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	chargeInvoice := v.ChargeInvoice.ToInvoice()
	creditInvoices := make([]Invoice, len(v.CreditInvoices))
	for index, ci := range v.CreditInvoices {
		creditInvoices[index] = ci.ToInvoice()
	}
	*i = InvoiceCollection{
		XMLName:        xml.Name{Local: "invoice_collection"},
		ChargeInvoice:  &chargeInvoice,
		CreditInvoices: creditInvoices,
	}

	return nil
}

// invoiceFields is used by custom unmarshal functions.
type invoiceFields struct {
	AccountCode             hrefString      `xml:"account,omitempty"`
	Address                 Address         `xml:"address,omitempty"`
	SubscriptionUUID        hrefString      `xml:"subscription,omitempty"`
	OriginalInvoiceNumber   hrefInt         `xml:"original_invoice,omitempty"`
	UUID                    string          `xml:"uuid,omitempty"`
	State                   string          `xml:"state,omitempty"`
	InvoiceNumberPrefix     string          `xml:"invoice_number_prefix,omitempty"`
	InvoiceNumber           int             `xml:"invoice_number,omitempty"`
	PONumber                string          `xml:"po_number,omitempty"`
	VATNumber               string          `xml:"vat_number,omitempty"`
	DiscountInCents         int             `xml:"discount_in_cents,omitempty"`
	SubtotalInCents         int             `xml:"subtotal_in_cents,omitempty"`
	TaxInCents              int             `xml:"tax_in_cents,omitempty"`
	TotalInCents            int             `xml:"total_in_cents,omitempty"`
	BalanceInCents          int             `xml:"balance_in_cents,omitempty"`
	Currency                string          `xml:"currency,omitempty"`
	DueOn                   NullTime        `xml:"due_on,omitempty"`
	CreatedAt               NullTime        `xml:"created_at,omitempty"`
	UpdatedAt               NullTime        `xml:"updated_at,omitempty"`
	AttemptNextCollectionAt NullTime        `xml:"attempt_next_collection_at,omitempty"`
	ClosedAt                NullTime        `xml:"closed_at,omitempty"`
	Type                    string          `xml:"type,omitempty"`
	Origin                  string          `xml:"origin,omitempty"`
	TaxType                 string          `xml:"tax_type,omitempty"`
	TaxRegion               string          `xml:"tax_region,omitempty"`
	TaxRate                 float64         `xml:"tax_rate,omitempty"`
	NetTerms                NullInt         `xml:"net_terms,omitempty"`
	CollectionMethod        string          `xml:"collection_method,omitempty"`
	LineItems               []Adjustment    `xml:"line_items>adjustment,omitempty"`
	Transactions            []Transaction   `xml:"transactions>transaction,omitempty"`
	CreditPayments          []CreditPayment `xml:"credit_payments>credit_payment,omitempty"`
}

func (i invoiceFields) ToInvoice() Invoice {
	return Invoice{
		XMLName:               xml.Name{Local: "invoice"},
		AccountCode:           string(i.AccountCode),
		Address:               i.Address,
		OriginalInvoiceNumber: int(i.OriginalInvoiceNumber),
		UUID:                    i.UUID,
		State:                   i.State,
		InvoiceNumberPrefix:     i.InvoiceNumberPrefix,
		InvoiceNumber:           i.InvoiceNumber,
		PONumber:                i.PONumber,
		VATNumber:               i.VATNumber,
		DiscountInCents:         i.DiscountInCents,
		SubtotalInCents:         i.SubtotalInCents,
		TaxInCents:              i.TaxInCents,
		TotalInCents:            i.TotalInCents,
		BalanceInCents:          i.BalanceInCents,
		Currency:                i.Currency,
		DueOn:                   i.DueOn,
		CreatedAt:               i.CreatedAt,
		UpdatedAt:               i.UpdatedAt,
		AttemptNextCollectionAt: i.AttemptNextCollectionAt,
		ClosedAt:                i.ClosedAt,
		Type:                    i.Type,
		Origin:                  i.Origin,
		TaxType:                 i.TaxType,
		TaxRegion:               i.TaxRegion,
		TaxRate:                 i.TaxRate,
		NetTerms:                i.NetTerms,
		CollectionMethod:        i.CollectionMethod,
		LineItems:               i.LineItems,
		Transactions:            i.Transactions,
		CreditPayments:          i.CreditPayments,
	}
}

// OfflinePayment is a payment received outside the system to be recorded in Recurly.
type OfflinePayment struct {
	XMLName       xml.Name   `xml:"transaction"`
	InvoiceNumber int        `xml:"-"`
	PaymentMethod string     `xml:"payment_method"`
	CollectedAt   *time.Time `xml:"collected_at,omitempty"`
	Amount        int        `xml:"amount_in_cents,omitempty"`
	Description   string     `xml:"description,omitempty"`
}
