package mock

import (
	"context"

	"github.com/blacklightcms/recurly"
)

var _ recurly.BillingService = &BillingService{}

// BillingService manages the interactions for billing.
type BillingService struct {
	OnGet      func(ctx context.Context, accountCode string) (*recurly.Billing, error)
	GetInvoked bool

	OnCreate      func(ctx context.Context, accountCode string, b recurly.Billing) (*recurly.Billing, error)
	CreateInvoked bool

	OnUpdate      func(ctx context.Context, accountCode string, b recurly.Billing) (*recurly.Billing, error)
	UpdateInvoked bool

	OnClear      func(ctx context.Context, accountCode string) error
	ClearInvoked bool
}

func (m *BillingService) Get(ctx context.Context, accountCode string) (*recurly.Billing, error) {
	m.GetInvoked = true
	return m.OnGet(ctx, accountCode)
}

func (m *BillingService) Create(ctx context.Context, accountCode string, b recurly.Billing) (*recurly.Billing, error) {
	m.CreateInvoked = true
	return m.OnCreate(ctx, accountCode, b)
}

func (m *BillingService) Update(ctx context.Context, accountCode string, b recurly.Billing) (*recurly.Billing, error) {
	m.UpdateInvoked = true
	return m.OnUpdate(ctx, accountCode, b)
}

func (m *BillingService) Clear(ctx context.Context, accountCode string) error {
	m.ClearInvoked = true
	return m.OnClear(ctx, accountCode)
}
