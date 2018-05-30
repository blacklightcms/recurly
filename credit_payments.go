package recurly

import "encoding/xml"

// Credit payment action constants.
const (
	CreditPaymentActionPayment   = "payment"
	CreditPaymentActionGiftCard  = "gift_card"
	CreditPaymentActionRefund    = "refund"
	CreditPaymentActionReduction = "reduction"
	CreditPaymentActionWriteOff  = "write_off"
)

// CreditPayment is a credit that has been applied to an invoice.
// This is a read-only object.
// Unmarshaling an invoice is handled by the custom UnmarshalXML function.
type CreditPayment struct {
	XMLName                   xml.Name `xml:"credit_payment"`
	AccountCode               string   `xml:"-"`
	UUID                      string   `xml:"-"`
	Action                    string   `xml:"-"`
	Currency                  string   `xml:"-"`
	AmountInCents             int      `xml:"-"`
	OriginalInvoiceNumber     int      `xml:"-"`
	AppliedToInvoice          int      `xml:"-"`
	OriginalCreditPaymentUUID string   `xml:"-"`
	RefundTransactionUUID     string   `xml:"-"`
	CreatedAt                 NullTime `xml:"-"`
	UpdatedAt                 NullTime `xml:"-"`
	VoidedAt                  NullTime `xml:"-"`
}

// UnmarshalXML unmarshals invoices and handles intermediary state during unmarshaling
// for types like href.
func (c *CreditPayment) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v struct {
		XMLName               xml.Name   `xml:"credit_payment"`
		AccountCode           hrefString `xml:"account"`
		UUID                  string     `xml:"uuid"`
		Action                string     `xml:"action"`
		Currency              string     `xml:"currency"`
		AmountInCents         int        `xml:"amount_in_cents"`
		OriginalInvoiceNumber hrefInt    `xml:"original_invoice"`
		AppliedToInvoice      hrefInt    `xml:"applied_to_invoice"`
		OriginalCreditPayment hrefString `xml:"original_credit_payment,omitempty"`
		RefundTransaction     hrefString `xml:"refund_transaction,omitempty"`
		CreatedAt             NullTime   `xml:"created_at"`
		UpdatedAt             NullTime   `xml:"updated_at,omitempty"`
		VoidedAt              NullTime   `xml:"voided_at,omitempty"`
	}
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	*c = CreditPayment{
		XMLName:                   v.XMLName,
		AccountCode:               string(v.AccountCode),
		UUID:                      v.UUID,
		Action:                    v.Action,
		Currency:                  v.Currency,
		AmountInCents:             v.AmountInCents,
		OriginalInvoiceNumber:     int(v.OriginalInvoiceNumber),
		AppliedToInvoice:          int(v.AppliedToInvoice),
		OriginalCreditPaymentUUID: string(v.OriginalCreditPayment),
		RefundTransactionUUID:     string(v.RefundTransaction),
		CreatedAt:                 v.CreatedAt,
		UpdatedAt:                 v.UpdatedAt,
		VoidedAt:                  v.VoidedAt,
	}

	return nil
}
