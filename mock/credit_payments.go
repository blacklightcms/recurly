package mock

import (
	"context"

	"github.com/autopilot3/recurly"
)

var _ recurly.CreditPaymentsService = &CreditPaymentsService{}

// CreditPaymentsService manages the interactions for credit payments.
type CreditPaymentsService struct {
	OnList      func(opts *recurly.PagerOptions) recurly.Pager
	ListInvoked bool

	OnListAccount      func(code string, opts *recurly.PagerOptions) recurly.Pager
	ListAccountInvoked bool

	OnGet      func(ctx context.Context, uuid string) (*recurly.CreditPayment, error)
	GetInvoked bool
}

func (m *CreditPaymentsService) List(opts *recurly.PagerOptions) recurly.Pager {
	m.ListInvoked = true
	return m.OnList(opts)
}

func (m *CreditPaymentsService) ListAccount(code string, opts *recurly.PagerOptions) recurly.Pager {
	m.ListAccountInvoked = true
	return m.OnListAccount(code, opts)
}

func (m *CreditPaymentsService) Get(ctx context.Context, uuid string) (*recurly.CreditPayment, error) {
	m.GetInvoked = true
	return m.OnGet(ctx, uuid)
}
