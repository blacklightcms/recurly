package mock

import (
	"context"

	"github.com/autopilot3/recurly"
)

var _ recurly.PurchasesService = &PurchasesService{}

type PurchasesService struct {
	OnCreate      func(ctx context.Context, p recurly.Purchase) (*recurly.InvoiceCollection, error)
	CreateInvoked bool

	OnPreview      func(ctx context.Context, p recurly.Purchase) (*recurly.InvoiceCollection, error)
	PreviewInvoked bool

	OnAuthorize      func(ctx context.Context, p recurly.Purchase) (*recurly.Purchase, error)
	AuthorizeInvoked bool

	OnPending      func(ctx context.Context, p recurly.Purchase) (*recurly.Purchase, error)
	PendingInvoked bool

	OnCapture      func(ctx context.Context, transactionUUID string) (*recurly.InvoiceCollection, error)
	CaptureInvoked bool

	OnCancel      func(ctx context.Context, transactionUUID string) (*recurly.InvoiceCollection, error)
	CancelInvoked bool
}

func (m *PurchasesService) Create(ctx context.Context, p recurly.Purchase) (*recurly.InvoiceCollection, error) {
	m.CreateInvoked = true
	return m.OnCreate(ctx, p)
}

func (m *PurchasesService) Preview(ctx context.Context, p recurly.Purchase) (*recurly.InvoiceCollection, error) {
	m.PreviewInvoked = true
	return m.OnPreview(ctx, p)
}

func (m *PurchasesService) Authorize(ctx context.Context, p recurly.Purchase) (*recurly.Purchase, error) {
	m.AuthorizeInvoked = true
	return m.OnAuthorize(ctx, p)
}

func (m *PurchasesService) Pending(ctx context.Context, p recurly.Purchase) (*recurly.Purchase, error) {
	m.PendingInvoked = true
	return m.OnPending(ctx, p)
}

func (m *PurchasesService) Capture(ctx context.Context, transactionUUID string) (*recurly.InvoiceCollection, error) {
	m.CaptureInvoked = true
	return m.OnCapture(ctx, transactionUUID)
}

func (m *PurchasesService) Cancel(ctx context.Context, transactionUUID string) (*recurly.InvoiceCollection, error) {
	m.CancelInvoked = true
	return m.OnCancel(ctx, transactionUUID)
}
