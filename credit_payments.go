package recurly

import "encoding/xml"

// Credit payment action constants.
const (
	CreditPaymentActionPayment   = "payment" // applying the credit
	CreditPaymentActionGiftCard  = "gift_card"
	CreditPaymentActionRefund    = "refund"
	CreditPaymentActionReduction = "reduction" // reducing the amount of the credit without applying it
	CreditPaymentActionWriteOff  = "write_off" // used for voiding invoices
)

// CreditPayment is a credit that has been applied to an invoice.
// This is a read-only object.
// Unmarshaling an invoice is handled by the custom UnmarshalXML function.
type CreditPayment struct {
	XMLName                   xml.Name `xml:"credit_payment"`
	AccountCode               string   `xml:"-"`
	UUID                      string   `xml:"uuid"`
	Action                    string   `xml:"action"`
	Currency                  string   `xml:"currency"`
	AmountInCents             int      `xml:"amount_in_cents"`
	OriginalInvoiceNumber     int      `xml:"-"`
	AppliedToInvoice          int      `xml:"-"`
	OriginalCreditPaymentUUID string   `xml:"-"`
	RefundTransactionUUID     string   `xml:"-"`
	CreatedAt                 NullTime `xml:"created_at"`
	UpdatedAt                 NullTime `xml:"updated_at,omitempty"`
	VoidedAt                  NullTime `xml:"voided_at,omitempty"`
}

// UnmarshalXML unmarshals invoices and handles intermediary state during unmarshaling
// for types like href.
func (c *CreditPayment) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type creditPaymentAlias CreditPayment
	var v struct {
		XMLName xml.Name `xml:"credit_payment"`
		creditPaymentAlias
		AccountCode           hrefString `xml:"account"`
		OriginalInvoiceNumber hrefInt    `xml:"original_invoice"`
		AppliedToInvoice      hrefInt    `xml:"applied_to_invoice"`
		OriginalCreditPayment hrefString `xml:"original_credit_payment,omitempty"`
		RefundTransaction     hrefString `xml:"refund_transaction,omitempty"`
	}
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	*c = CreditPayment(v.creditPaymentAlias)
	c.XMLName = v.XMLName
	c.AccountCode = string(v.AccountCode)
	c.OriginalInvoiceNumber = int(v.OriginalInvoiceNumber)
	c.AppliedToInvoice = int(v.AppliedToInvoice)
	c.OriginalCreditPaymentUUID = string(v.OriginalCreditPayment)
	c.RefundTransactionUUID = string(v.RefundTransaction)

	return nil
}
