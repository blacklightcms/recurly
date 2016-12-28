package mock

import "github.com/blacklightcms/go-recurly/recurly"

var _ recurly.TransactionsInterface = &TransactionService{}

// TransactionService mocks the transaction service.
type TransactionService struct {
	OnList      func(params Params) (*recurly.Response, []recurly.Transaction, error)
	ListInvoked bool

	OnListAccount      func(accountCode string, params recurly.Params) (*recurly.Response, []recurly.Transaction, error)
	ListAccountInvoked bool

	OnGet      func(uuid string) (*recurly.Response, *recurly.Transaction, error)
	GetInvoked bool

	OnCreate      func(trans recurly.Transaction) (*recurly.Response, *recurly.Transaction, error)
	CreateInvoked bool
}

func (t *TransactionService) List(params recurly.Params) (*recurly.Response, []recurly.Transaction, error) {
	t.ListInvoked = true
	return t.OnList(params)
}

func (t *TransactionService) ListAccount(accountCode string, params recurly.Params) (*recurly.Response, []recurly.Transaction, error) {
	t.ListAccountInvoked = true
	return t.OnListAccount(accountCode, params)
}

func (t *TransactionService) Get(uuid string) (*recurly.Response, *recurly.Transaction, error) {
	t.GetInvoked = true
	return t.OnGet(uuid)
}

func (t *TransactionService) Create(trans recurly.Transaction) (*recurly.Response, *recurly.Transaction, error) {
	t.CreateInvoked = true
	return t.OnCreate(trans)
}
