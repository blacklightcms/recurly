package mock

import (
	"context"

	"github.com/blacklightcms/recurly"
)

var _ recurly.TransactionsService = &TransactionsService{}

// TransactionsService mocks the transaction service.
type TransactionsService struct {
	OnList      func(opts *recurly.PagerOptions) recurly.Pager
	ListInvoked bool

	OnListAccount      func(accountCode string, opts *recurly.PagerOptions) recurly.Pager
	ListAccountInvoked bool

	OnGet      func(ctx context.Context, uuid string) (*recurly.Transaction, error)
	GetInvoked bool
}

func (m *TransactionsService) List(opts *recurly.PagerOptions) recurly.Pager {
	m.ListInvoked = true
	return m.OnList(opts)
}

func (m *TransactionsService) ListAccount(accountCode string, opts *recurly.PagerOptions) recurly.Pager {
	m.ListAccountInvoked = true
	return m.OnListAccount(accountCode, opts)
}

func (m *TransactionsService) Get(ctx context.Context, uuid string) (*recurly.Transaction, error) {
	m.GetInvoked = true
	return m.OnGet(ctx, uuid)
}
