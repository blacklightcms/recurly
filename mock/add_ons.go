package mock

import (
	"context"

	"github.com/blacklightcms/recurly"
)

var _ recurly.AddOnsService = &AddOnsService{}

// AddOnsService manages the interactions for add ons.
type AddOnsService struct {
	OnList      func(planCode string, opts *recurly.PagerOptions) *recurly.AddOnsPager
	ListInvoked bool

	OnGet      func(ctx context.Context, planCode string, code string) (*recurly.AddOn, error)
	GetInvoked bool

	OnCreate      func(ctx context.Context, planCode string, a recurly.AddOn) (*recurly.AddOn, error)
	CreateInvoked bool

	OnUpdate      func(ctx context.Context, planCode string, code string, a recurly.AddOn) (*recurly.AddOn, error)
	UpdateInvoked bool

	OnDelete      func(ctx context.Context, planCode string, code string) error
	DeleteInvoked bool
}

func (m *AddOnsService) List(planCode string, opts *recurly.PagerOptions) *recurly.AddOnsPager {
	m.ListInvoked = true
	return m.OnList(planCode, opts)
}

func (m *AddOnsService) Get(ctx context.Context, planCode string, code string) (*recurly.AddOn, error) {
	m.GetInvoked = true
	return m.OnGet(ctx, planCode, code)
}

func (m *AddOnsService) Create(ctx context.Context, planCode string, a recurly.AddOn) (*recurly.AddOn, error) {
	m.CreateInvoked = true
	return m.OnCreate(ctx, planCode, a)
}

func (m *AddOnsService) Update(ctx context.Context, planCode string, code string, a recurly.AddOn) (*recurly.AddOn, error) {
	m.UpdateInvoked = true
	return m.OnUpdate(ctx, planCode, code, a)
}

func (m *AddOnsService) Delete(ctx context.Context, planCode string, code string) error {
	m.DeleteInvoked = true
	return m.OnDelete(ctx, planCode, code)
}
