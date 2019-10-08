package mock

import (
	"bytes"
	"context"

	"github.com/blacklightcms/recurly"
)

var _ recurly.InvoicesService = &InvoicesService{}

// InvoicesService manages the interactions for invoices.
type InvoicesService struct {
	OnList      func(opts *recurly.PagerOptions) recurly.Pager
	ListInvoked bool

	OnListAccount      func(accountCode string, opts *recurly.PagerOptions) recurly.Pager
	ListAccountInvoked bool

	OnGet      func(ctx context.Context, invoiceNumber int) (*recurly.Invoice, error)
	GetInvoked bool

	OnGetPDF      func(ctx context.Context, invoiceNumber int, language string) (*bytes.Buffer, error)
	GetPDFInvoked bool

	OnPreview      func(ctx context.Context, accountCode string) (*recurly.Invoice, error)
	PreviewInvoked bool

	OnCreate      func(ctx context.Context, accountCode string, invoice recurly.Invoice) (*recurly.Invoice, error)
	CreateInvoked bool

	OnCollect      func(ctx context.Context, invoiceNumber int, collectInvoice recurly.CollectInvoice) (*recurly.Invoice, error)
	CollectInvoked bool

	OnMarkPaid      func(ctx context.Context, invoiceNumber int) (*recurly.Invoice, error)
	MarkPaidInvoked bool

	OnMarkFailed      func(ctx context.Context, invoiceNumber int) (*recurly.Invoice, error)
	MarkFailedInvoked bool

	OnRefundVoidLineItems      func(ctx context.Context, invoiceNumber int, refund recurly.InvoiceLineItemsRefund) (*recurly.Invoice, error)
	RefundVoidLineItemsInvoked bool

	OnRefundVoidOpenAmount      func(ctx context.Context, invoiceNumber int, refund recurly.InvoiceRefund) (*recurly.Invoice, error)
	RefundVoidOpenAmountInvoked bool

	OnVoidCreditInvoice      func(ctx context.Context, invoiceNumber int) (*recurly.Invoice, error)
	VoidCreditInvoiceInvoked bool

	OnRecordPayment      func(ctx context.Context, pmt recurly.OfflinePayment) (*recurly.Transaction, error)
	RecordPaymentInvoked bool
}

func (m *InvoicesService) List(opts *recurly.PagerOptions) recurly.Pager {
	m.ListInvoked = true
	return m.OnList(opts)
}

func (m *InvoicesService) ListAccount(accountCode string, opts *recurly.PagerOptions) recurly.Pager {
	m.ListAccountInvoked = true
	return m.OnListAccount(accountCode, opts)
}

func (m *InvoicesService) Get(ctx context.Context, invoiceNumber int) (*recurly.Invoice, error) {
	m.GetInvoked = true
	return m.OnGet(ctx, invoiceNumber)
}

func (m *InvoicesService) GetPDF(ctx context.Context, invoiceNumber int, language string) (*bytes.Buffer, error) {
	m.GetPDFInvoked = true
	return m.OnGetPDF(ctx, invoiceNumber, language)
}

func (m *InvoicesService) Preview(ctx context.Context, accountCode string) (*recurly.Invoice, error) {
	m.PreviewInvoked = true
	return m.OnPreview(ctx, accountCode)
}

func (m *InvoicesService) Create(ctx context.Context, accountCode string, invoice recurly.Invoice) (*recurly.Invoice, error) {
	m.CreateInvoked = true
	return m.OnCreate(ctx, accountCode, invoice)
}

func (m *InvoicesService) Collect(ctx context.Context, invoiceNumber int, collectInvoice recurly.CollectInvoice) (*recurly.Invoice, error) {
	m.CollectInvoked = true
	return m.OnCollect(ctx, invoiceNumber, collectInvoice)
}

func (m *InvoicesService) MarkPaid(ctx context.Context, invoiceNumber int) (*recurly.Invoice, error) {
	m.MarkPaidInvoked = true
	return m.OnMarkPaid(ctx, invoiceNumber)
}

func (m *InvoicesService) MarkFailed(ctx context.Context, invoiceNumber int) (*recurly.Invoice, error) {
	m.MarkFailedInvoked = true
	return m.OnMarkFailed(ctx, invoiceNumber)
}

func (m *InvoicesService) RefundVoidLineItems(ctx context.Context, invoiceNumber int, refund recurly.InvoiceLineItemsRefund) (*recurly.Invoice, error) {
	m.RefundVoidLineItemsInvoked = true
	return m.OnRefundVoidLineItems(ctx, invoiceNumber, refund)
}

func (m *InvoicesService) RefundVoidOpenAmount(ctx context.Context, invoiceNumber int, refund recurly.InvoiceRefund) (*recurly.Invoice, error) {
	m.RefundVoidOpenAmountInvoked = true
	return m.OnRefundVoidOpenAmount(ctx, invoiceNumber, refund)
}

func (m *InvoicesService) VoidCreditInvoice(ctx context.Context, invoiceNumber int) (*recurly.Invoice, error) {
	m.VoidCreditInvoiceInvoked = true
	return m.OnVoidCreditInvoice(ctx, invoiceNumber)
}

func (m *InvoicesService) RecordPayment(ctx context.Context, pmt recurly.OfflinePayment) (*recurly.Transaction, error) {
	m.RecordPaymentInvoked = true
	return m.OnRecordPayment(ctx, pmt)
}
