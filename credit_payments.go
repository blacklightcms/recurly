package recurly

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
)

// CreditPaymentsService manages the interactions for credit payments.
type CreditPaymentsService interface {
	// List returns a pager to paginate credit payments. PagerOptions are used to optionally
	// filter the results.
	//
	// https://dev.recurly.com/docs/list-credit-payments
	List(opts *PagerOptions) *CreditPaymentsPager

	// ListAccount returns a pager to paginate credit payments for an account.
	// PagerOptions are used to optionally filter the results.
	//
	// https://dev.recurly.com/docs/list-credit-payments
	ListAccount(accountCode string, opts *PagerOptions) *CreditPaymentsPager

	// Get retrieves a credit payment. If the credit payment does not exist,
	// a nil credit payment and nil error are returned.
	//
	// https://dev.recurly.com/docs/lookup-credit-payment
	Get(ctx context.Context, uuid string) (*CreditPayment, error)
}

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
//
// https://dev.recurly.com/docs/creditpayment-object
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
		AccountCode           href    `xml:"account"`
		OriginalInvoiceNumber hrefInt `xml:"original_invoice"`
		AppliedToInvoice      hrefInt `xml:"applied_to_invoice"`
		OriginalCreditPayment href    `xml:"original_credit_payment,omitempty"`
		RefundTransaction     href    `xml:"refund_transaction,omitempty"`
	}
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}

	*c = CreditPayment(v.creditPaymentAlias)
	c.XMLName = v.XMLName
	c.AccountCode = v.AccountCode.LastPartOfPath()
	c.OriginalInvoiceNumber = v.OriginalInvoiceNumber.LastPartOfPath()
	c.AppliedToInvoice = v.AppliedToInvoice.LastPartOfPath()
	c.OriginalCreditPaymentUUID = v.OriginalCreditPayment.LastPartOfPath()
	c.RefundTransactionUUID = v.RefundTransaction.LastPartOfPath()
	return nil
}

// CreditPaymentsPager paginates credit payments.
type CreditPaymentsPager struct {
	*pager
}

// Fetch fetches the next set of results.
func (p *CreditPaymentsPager) Fetch(ctx context.Context) ([]CreditPayment, error) {
	var dst struct {
		XMLName        xml.Name        `xml:"credit_payments"`
		CreditPayments []CreditPayment `xml:"credit_payment"`
	}
	if err := p.fetch(ctx, &dst); err != nil {
		return nil, err
	}
	return dst.CreditPayments, nil
}

// FetchAll paginates all records and returns a cumulative list.
func (p *CreditPaymentsPager) FetchAll(ctx context.Context) ([]CreditPayment, error) {
	p.setMaxPerPage()

	var all []CreditPayment
	for p.Next() {
		v, err := p.Fetch(ctx)
		if err != nil {
			return nil, err
		}
		all = append(all, v...)
	}
	return all, nil
}

var _ CreditPaymentsService = &creditInvoicesImpl{}

// creditInvoicesImpl implements CreditPaymentsService.
type creditInvoicesImpl serviceImpl

func (s *creditInvoicesImpl) List(opts *PagerOptions) *CreditPaymentsPager {
	return &CreditPaymentsPager{
		pager: s.client.newPager("GET", "/credit_payments", opts),
	}
}

func (s *creditInvoicesImpl) ListAccount(accountCode string, opts *PagerOptions) *CreditPaymentsPager {
	path := fmt.Sprintf("/accounts/%s/credit_payments", accountCode)
	return &CreditPaymentsPager{
		pager: s.client.newPager("GET", path, opts),
	}
}

func (s *creditInvoicesImpl) Get(ctx context.Context, uuid string) (*CreditPayment, error) {
	path := fmt.Sprintf("/credit_payments/%s", sanitizeUUID(uuid))
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var dst CreditPayment
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		if e, ok := err.(*ClientError); ok && e.Response.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &dst, nil
}
