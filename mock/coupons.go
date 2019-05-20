package mock

import (
	"context"

	"github.com/blacklightcms/recurly"
)

var _ recurly.CouponsService = &CouponsService{}

// CouponsService manages the interactions for coupons.
type CouponsService struct {
	OnList      func(opts *recurly.PagerOptions) recurly.Pager
	ListInvoked bool

	OnGet      func(ctx context.Context, code string) (*recurly.Coupon, error)
	GetInvoked bool

	OnCreate      func(ctx context.Context, c recurly.Coupon) (*recurly.Coupon, error)
	CreateInvoked bool

	OnUpdate      func(ctx context.Context, code string, c recurly.Coupon) (*recurly.Coupon, error)
	UpdateInvoked bool

	OnRestore      func(ctx context.Context, code string, c recurly.Coupon) (*recurly.Coupon, error)
	RestoreInvoked bool

	OnDelete      func(ctx context.Context, code string) error
	DeleteInvoked bool

	OnGenerate      func(ctx context.Context, code string, n int) (recurly.Pager, error)
	GenerateInvoked bool
}

func (m *CouponsService) List(opts *recurly.PagerOptions) recurly.Pager {
	m.ListInvoked = true
	return m.OnList(opts)
}

func (m *CouponsService) Get(ctx context.Context, code string) (*recurly.Coupon, error) {
	m.GetInvoked = true
	return m.OnGet(ctx, code)
}

func (m *CouponsService) Create(ctx context.Context, c recurly.Coupon) (*recurly.Coupon, error) {
	m.CreateInvoked = true
	return m.OnCreate(ctx, c)
}

func (m *CouponsService) Update(ctx context.Context, code string, c recurly.Coupon) (*recurly.Coupon, error) {
	m.UpdateInvoked = true
	return m.OnUpdate(ctx, code, c)
}

func (m *CouponsService) Restore(ctx context.Context, code string, c recurly.Coupon) (*recurly.Coupon, error) {
	m.RestoreInvoked = true
	return m.OnRestore(ctx, code, c)
}

func (m *CouponsService) Delete(ctx context.Context, code string) error {
	m.DeleteInvoked = true
	return m.OnDelete(ctx, code)
}

func (m *CouponsService) Generate(ctx context.Context, code string, n int) (recurly.Pager, error) {
	m.GenerateInvoked = true
	return m.OnGenerate(ctx, code, n)
}
