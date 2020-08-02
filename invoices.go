package recurly

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
)

// InvoicesService manages the interactions for invoices.
type InvoicesService interface {
	// List returns a pager to paginate invoices. PagerOptions are used to optionally
	// filter the results.
	//
	// https://dev.recurly.com/docs/list-invoices
	List(opts *PagerOptions) Pager

	// ListAccount returns a pager to paginate invoices for an account. Params
	// are used to optionally filter the results.
	//
	// https://dev.recurly.com/docs/list-an-accounts-invoices
	ListAccount(accountCode string, opts *PagerOptions) Pager

	// Get retrieves an invoice. If the invoice does not exist,
	// a nil invoice and nil error are returned.
	//
	// https://dev.recurly.com/docs/lookup-invoice-details
	Get(ctx context.Context, invoiceNumber int) (*Invoice, error)

	// GetPDF retrieves an invoice as a PDF. If the invoice does not exist,
	// a nil bytes.Buffer and nil error are returned. If language is not provided
	// or not valid, English will be used.
	//
	// https://dev.recurly.com/docs/retrieve-a-pdf-invoice
	GetPDF(ctx context.Context, invoiceNumber int, language string) (*bytes.Buffer, error)

	// Preview allows you to display the invoice details, including estimated
	// tax, before you post it.
	//
	// https://dev.recurly.com/docs/preview-an-invoice
	Preview(ctx context.Context, accountCode string) (*Invoice, error)

	// Create posts an invoice with pending adjustments to the account.
	// See Recurly's documentation for details.
	//
	// https://dev.recurly.com/docs/post-an-invoice-invoice-pending-charges-on-an-acco
	Create(ctx context.Context, accountCode string, invoice Invoice) (*Invoice, error)

	// Collect allows collecting a past-due or pending invoice. This API
	// is rate limited. See Recurly's documentation for details.
	//
	// https://dev.recurly.com/docs/collect-an-invoice
	Collect(ctx context.Context, invoiceNumber int, collectInvoice CollectInvoice) (*Invoice, error)

	// MarkPaid marks an invoice as paid successfully.
	//
	// https://dev.recurly.com/docs/mark-an-invoice-as-paid-successfully
	MarkPaid(ctx context.Context, invoiceNumber int) (*Invoice, error)

	// MarkFailed marks an invoice as failed collection.
	//
	// https://dev.recurly.com/docs/mark-an-invoice-as-failed-collection
	MarkFailed(ctx context.Context, invoiceNumber int) (*Invoice, error)

	// RefundVoidLineItems allows specific invoice line items and/or quantities to
	// be refunded and generates a refund invoice. Full amount line item refunds
	//of invoices with an unsettled transaction will voide the transaction and
	// generate a void invoice. See Recurly's documentation for details.
	//
	// https://dev.recurly.com/docs/line-item-refunds
	RefundVoidLineItems(ctx context.Context, invoiceNumber int, refund InvoiceLineItemsRefund) (*Invoice, error)

	// RefundVoidOpenAmount allows custom invoice amounts to be refunded and
	// generates a refund invoice. Full open amount refunds of invoices with
	// an unsettled transaction will void the transaction and generate a
	// void invoice. See Recurly's documentation for details.
	//
	// https://dev.recurly.com/docs/open-amount-refunds
	RefundVoidOpenAmount(ctx context.Context, invoiceNumber int, refund InvoiceRefund) (*Invoice, error)

	// VoidCreditInvoice voids an open credit invoice.
	//
	// https://dev.recurly.com/docs/void-credit-invoice
	VoidCreditInvoice(ctx context.Context, invoiceNumber int) (*Invoice, error)

	// RecordPayment records an offline payment for a manual invoice, such as a check
	// or money order.
	//
	// https://dev.recurly.com/docs/enter-an-offline-payment-for-a-manual-invoice
	RecordPayment(ctx context.Context, offlinePayment OfflinePayment) (*Transaction, error)
}

// Invoice state constants.
// https://docs.recurly.com/docs/invoices#section-statuses
const (
	ChargeInvoiceStatePending    = "pending"
	ChargeInvoiceStateProcessing = "processing"
	ChargeInvoiceStatePastDue    = "past_due"
	ChargeInvoiceStatePaid       = "paid"
	ChargeInvoiceStateFailed     = "failed"

	CreditInvoiceStateOpen       = "open"
	CreditInvoiceStateProcessing = "processing"
	CreditInvoiceStateClosed     = "closed"
	CreditInvoiceStateVoided     = "voided"
)

// Collection method constants.
const (
	CollectionMethodAutomatic = "automatic"
	CollectionMethodManual    = "manual"
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
	CreditInvoiceOriginGiftCard        = "gift_card"
	CreditInvoiceOriginRefund          = "refund"
	CreditInvoiceOriginCredit          = "credit"
	CreditInvoiceOriginWriteOff        = "write_off"
	CreditInvoiceOriginExternalCredit  = "external_credit"
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

// Invoice is an individual invoice for an account. Transactions are guaranteed
// to be sorted from oldest to newest by date.
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

	// TaxDetails is only available if the site has the  `Avalara for Communications` integration
	TaxDetails *[]TaxDetail `xml:"tax_details>tax_detail,omitempty"`
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

	*i = v.ToInvoice()
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
	for i := range v.CreditInvoices {
		creditInvoices[i] = v.CreditInvoices[i].ToInvoice()
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
	AccountCode             href            `xml:"account,omitempty"`
	Address                 Address         `xml:"address,omitempty"`
	SubscriptionUUID        href            `xml:"subscription,omitempty"`
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
	TaxDetails              *[]TaxDetail    `xml:"tax_details>tax_detail,omitempty"`
}

// convert to Invoice and sort transactions.
func (i invoiceFields) ToInvoice() Invoice {
	inv := Invoice{
		XMLName:                 xml.Name{Local: "invoice"},
		AccountCode:             i.AccountCode.LastPartOfPath(),
		Address:                 i.Address,
		OriginalInvoiceNumber:   i.OriginalInvoiceNumber.LastPartOfPath(),
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
		TaxDetails:              i.TaxDetails,
	}
	Transactions(inv.Transactions).Sort()
	return inv
}

// OfflinePayment is a payment received outside of Recurly (e.g. ACH/Wire).
type OfflinePayment struct {
	XMLName       xml.Name `xml:"transaction"`
	InvoiceNumber int      `xml:"-"`
	PaymentMethod string   `xml:"payment_method"`
	CollectedAt   NullTime `xml:"collected_at,omitempty"`
	Amount        int      `xml:"amount_in_cents,omitempty"`
	Description   string   `xml:"description,omitempty"`
}

// InvoiceRefund is used to refund invoices.
type InvoiceRefund struct {
	XMLName             xml.Name `xml:"invoice"`
	AmountInCents       NullInt  `xml:"amount_in_cents,omitempty"` // If left empty the remaining refundable amount will be refunded
	RefundMethod        string   `xml:"refund_method,omitempty"`
	ExternalRefund      NullBool `xml:"external_refund,omitempty"`
	CreditCustomerNotes string   `xml:"credit_customer_notes,omitempty"`
	PaymentMethod       string   `xml:"payment_method,omitempty"`
	Description         string   `xml:"description,omitempty"`
	RefundedAt          NullTime `xml:"refunded_at,omitempty"`
}

// CollectInvoice is used as the request body for collecting an invoice.
type CollectInvoice struct {
	XMLName         xml.Name `xml:"invoice"`
	TransactionType string   `xml:"transaction_type,omitempty"` // Optional transaction type. Currently accepts "moto"
	BillingInfo     *Billing `xml:"billing_info,omitempty"`
}

// InvoiceLineItemsRefund is used to refund one or more line items on an invoice.
type InvoiceLineItemsRefund struct {
	XMLName   xml.Name       `xml:"invoice"`
	LineItems []VoidLineItem `xml:"line_items>adjustment"`
	InvoiceRefund
}

// VoidLineItem is an individual line item to refund.
type VoidLineItem struct {
	XMLName  xml.Name `xml:"adjustment"`
	UUID     string   `xml:"uuid"` // Adjustment UUID
	Quantity int      `xml:"quantity"`
	Prorate  NullBool `xml:"prorate,omitempty"`
}

var _ InvoicesService = &invoicesImpl{}

// invoicesImpl implements InvoicesService.
type invoicesImpl serviceImpl

func (s *invoicesImpl) List(opts *PagerOptions) Pager {
	return s.client.newPager("GET", "/invoices", opts)
}

func (s *invoicesImpl) ListAccount(accountCode string, opts *PagerOptions) Pager {
	path := fmt.Sprintf("/accounts/%s/invoices", accountCode)
	return s.client.newPager("GET", path, opts)
}

func (s *invoicesImpl) Get(ctx context.Context, invoiceNumber int) (*Invoice, error) {
	path := fmt.Sprintf("/invoices/%d", invoiceNumber)
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var dst Invoice
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		if e, ok := err.(*ClientError); ok && e.Response.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &dst, nil
}

func (s *invoicesImpl) GetPDF(ctx context.Context, invoiceNumber int, language string) (*bytes.Buffer, error) {
	switch language {
	case "English", "Danish", "German", "Spanish", "French", "Hindi",
		"Japanese", "Dutch", "Portuguese", "Russian", "Turkish", "Chinese":
	default:
		language = "English"
	}

	path := fmt.Sprintf("/invoices/%d", invoiceNumber)
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/pdf")
	req.Header.Set("Accept-Language", language)

	b := new(bytes.Buffer)
	if _, err := s.client.do(ctx, req, b); err != nil {
		if e, ok := err.(*ClientError); ok && e.Response.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, err
	}
	return b, nil
}

func (s *invoicesImpl) Preview(ctx context.Context, accountCode string) (*Invoice, error) {
	path := fmt.Sprintf("/accounts/%s/invoices/preview", accountCode)
	req, err := s.client.newRequest("POST", path, nil)
	if err != nil {
		return nil, err
	}

	var dst InvoiceCollection
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return dst.ChargeInvoice, nil
}

func (s *invoicesImpl) Create(ctx context.Context, accountCode string, invoice Invoice) (*Invoice, error) {
	path := fmt.Sprintf("/accounts/%s/invoices", accountCode)
	req, err := s.client.newRequest("POST", path, invoice)
	if err != nil {
		return nil, err
	}

	var dst InvoiceCollection
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return dst.ChargeInvoice, nil
}

func (s *invoicesImpl) Collect(ctx context.Context, invoiceNumber int, collectInvoice CollectInvoice) (*Invoice, error) {
	path := fmt.Sprintf("/invoices/%d/collect", invoiceNumber)
	req, err := s.client.newRequest("PUT", path, collectInvoice)
	if err != nil {
		return nil, err
	}

	var dst Invoice
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *invoicesImpl) MarkPaid(ctx context.Context, invoiceNumber int) (*Invoice, error) {
	path := fmt.Sprintf("/invoices/%d/mark_successful", invoiceNumber)
	req, err := s.client.newRequest("PUT", path, nil)
	if err != nil {
		return nil, err
	}

	var dst Invoice
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *invoicesImpl) MarkFailed(ctx context.Context, invoiceNumber int) (*Invoice, error) {
	path := fmt.Sprintf("/invoices/%d/mark_failed", invoiceNumber)
	req, err := s.client.newRequest("PUT", path, nil)
	if err != nil {
		return nil, err
	}

	var dst InvoiceCollection
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return dst.ChargeInvoice, nil
}

func (s *invoicesImpl) RefundVoidLineItems(ctx context.Context, invoiceNumber int, refund InvoiceLineItemsRefund) (*Invoice, error) {
	for i := range refund.LineItems {
		refund.LineItems[i].UUID = sanitizeUUID(refund.LineItems[i].UUID)
	}
	path := fmt.Sprintf("/invoices/%d/refund", invoiceNumber)
	req, err := s.client.newRequest("POST", path, refund)
	if err != nil {
		return nil, err
	}

	var dst Invoice
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *invoicesImpl) RefundVoidOpenAmount(ctx context.Context, invoiceNumber int, refund InvoiceRefund) (*Invoice, error) {
	path := fmt.Sprintf("/invoices/%d/refund", invoiceNumber)
	req, err := s.client.newRequest("POST", path, refund)
	if err != nil {
		return nil, err
	}

	var dst Invoice
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *invoicesImpl) VoidCreditInvoice(ctx context.Context, invoiceNumber int) (*Invoice, error) {
	path := fmt.Sprintf("/invoices/%d/void", invoiceNumber)
	req, err := s.client.newRequest("PUT", path, nil)
	if err != nil {
		return nil, err
	}

	var dst Invoice
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *invoicesImpl) RecordPayment(ctx context.Context, offlinePayment OfflinePayment) (*Transaction, error) {
	path := fmt.Sprintf("/invoices/%d/transactions", offlinePayment.InvoiceNumber)
	req, err := s.client.newRequest("POST", path, offlinePayment)
	if err != nil {
		return nil, err
	}

	var dst Transaction
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, nil
}
