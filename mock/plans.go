package mock

import (
	"context"

	"github.com/autopilot3/recurly"
)

var _ recurly.PlansService = &PlansService{}

// PlansService manages the interactions for plans.
type PlansService struct {
	OnList      func(opts *recurly.PagerOptions) recurly.Pager
	ListInvoked bool

	OnGet      func(ctx context.Context, code string) (*recurly.Plan, error)
	GetInvoked bool

	OnCreate      func(ctx context.Context, p recurly.Plan) (*recurly.Plan, error)
	CreateInvoked bool

	OnUpdate      func(ctx context.Context, code string, p recurly.Plan) (*recurly.Plan, error)
	UpdateInvoked bool

	OnDelete      func(ctx context.Context, code string) error
	DeleteInvoked bool
}

func (m *PlansService) List(opts *recurly.PagerOptions) recurly.Pager {
	m.ListInvoked = true
	return m.OnList(opts)
}

func (m *PlansService) Get(ctx context.Context, code string) (*recurly.Plan, error) {
	m.GetInvoked = true
	return m.OnGet(ctx, code)
}

func (m *PlansService) Create(ctx context.Context, p recurly.Plan) (*recurly.Plan, error) {
	m.CreateInvoked = true
	return m.OnCreate(ctx, p)
}

func (m *PlansService) Update(ctx context.Context, code string, p recurly.Plan) (*recurly.Plan, error) {
	m.UpdateInvoked = true
	return m.OnUpdate(ctx, code, p)
}

func (m *PlansService) Delete(ctx context.Context, code string) error {
	m.DeleteInvoked = true
	return m.OnDelete(ctx, code)
}
