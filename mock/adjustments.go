package mock

import (
	"context"

	"github.com/blacklightcms/recurly"
)

var _ recurly.AdjustmentsService = &AdjustmentsService{}

// AdjustmentsService manages the interactions for adjustments.
type AdjustmentsService struct {
	OnListAccount      func(accountCode string, opts *recurly.PagerOptions) recurly.Pager
	ListAccountInvoked bool

	OnGet      func(ctx context.Context, uuid string) (*recurly.Adjustment, error)
	GetInvoked bool

	OnCreate      func(ctx context.Context, accountCode string, a recurly.Adjustment) (*recurly.Adjustment, error)
	CreateInvoked bool

	OnDelete      func(ctx context.Context, uuid string) error
	DeleteInvoked bool
}

func (m *AdjustmentsService) ListAccount(accountCode string, opts *recurly.PagerOptions) recurly.Pager {
	m.ListAccountInvoked = true
	return m.OnListAccount(accountCode, opts)
}

func (m *AdjustmentsService) Get(ctx context.Context, uuid string) (*recurly.Adjustment, error) {
	m.GetInvoked = true
	return m.OnGet(ctx, uuid)
}

func (m *AdjustmentsService) Create(ctx context.Context, accountCode string, a recurly.Adjustment) (*recurly.Adjustment, error) {
	m.CreateInvoked = true
	return m.OnCreate(ctx, accountCode, a)
}

func (m *AdjustmentsService) Delete(ctx context.Context, uuid string) error {
	m.DeleteInvoked = true
	return m.OnDelete(ctx, uuid)
}
