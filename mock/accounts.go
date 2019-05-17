package mock

import (
	"context"

	"github.com/blacklightcms/recurly"
)

var _ recurly.AccountsService = &AccountsService{}

// AccountsService manages the interactions for accounts.
type AccountsService struct {
	OnList      func(opts *recurly.PagerOptions) *recurly.AccountsPager
	ListInvoked bool

	OnGet      func(ctx context.Context, code string) (*recurly.Account, error)
	GetInvoked bool

	OnBalance      func(ctx context.Context, code string) (*recurly.AccountBalance, error)
	BalanceInvoked bool

	OnCreate      func(ctx context.Context, a recurly.Account) (*recurly.Account, error)
	CreateInvoked bool

	OnUpdate      func(ctx context.Context, code string, a recurly.Account) (*recurly.Account, error)
	UpdateInvoked bool

	OnClose      func(ctx context.Context, code string) error
	CloseInvoked bool

	OnReopen      func(ctx context.Context, code string) error
	ReopenInvoked bool

	OnListNotes      func(code string, opts *recurly.PagerOptions) *recurly.NotesPager
	ListNotesInvoked bool
}

func (m *AccountsService) List(opts *recurly.PagerOptions) *recurly.AccountsPager {
	m.ListInvoked = true
	return m.OnList(opts)
}

func (m *AccountsService) Get(ctx context.Context, code string) (*recurly.Account, error) {
	m.GetInvoked = true
	return m.OnGet(ctx, code)
}

func (m *AccountsService) Balance(ctx context.Context, code string) (*recurly.AccountBalance, error) {
	m.BalanceInvoked = true
	return m.OnBalance(ctx, code)
}

func (m *AccountsService) Create(ctx context.Context, a recurly.Account) (*recurly.Account, error) {
	m.CreateInvoked = true
	return m.OnCreate(ctx, a)
}

func (m *AccountsService) Update(ctx context.Context, code string, a recurly.Account) (*recurly.Account, error) {
	m.UpdateInvoked = true
	return m.OnUpdate(ctx, code, a)
}

func (m *AccountsService) Close(ctx context.Context, code string) error {
	m.CloseInvoked = true
	return m.OnClose(ctx, code)
}

func (m *AccountsService) Reopen(ctx context.Context, code string) error {
	m.ReopenInvoked = true
	return m.OnReopen(ctx, code)
}

func (m *AccountsService) ListNotes(code string, opts *recurly.PagerOptions) *recurly.NotesPager {
	m.ListNotesInvoked = true
	return m.OnListNotes(code, opts)
}
