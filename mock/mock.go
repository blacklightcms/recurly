package mock

import "github.com/blacklightcms/go-recurly/recurly"

var _ recurly.Transactions = &TransactionsService{}

// TransactionsService mocks the transaction service.
type TransactionsService struct {
	OnList      func(params recurly.Params) (*recurly.Response, []recurly.Transaction, error)
	ListInvoked bool

	OnListAccount      func(accountCode string, params recurly.Params) (*recurly.Response, []recurly.Transaction, error)
	ListAccountInvoked bool

	OnGet      func(uuid string) (*recurly.Response, *recurly.Transaction, error)
	GetInvoked bool

	OnCreate      func(trans recurly.Transaction) (*recurly.Response, *recurly.Transaction, error)
	CreateInvoked bool
}

func (m *TransactionsService) List(params recurly.Params) (*recurly.Response, []recurly.Transaction, error) {
	m.ListInvoked = true
	return m.OnList(params)
}

func (m *TransactionsService) ListAccount(accountCode string, params recurly.Params) (*recurly.Response, []recurly.Transaction, error) {
	m.ListAccountInvoked = true
	return m.OnListAccount(accountCode, params)
}

func (m *TransactionsService) Get(uuid string) (*recurly.Response, *recurly.Transaction, error) {
	m.GetInvoked = true
	return m.OnGet(uuid)
}

func (m *TransactionsService) Create(t recurly.Transaction) (*recurly.Response, *recurly.Transaction, error) {
	m.CreateInvoked = true
	return m.OnCreate(t)
}

var _ recurly.Subscriptions = &SubscriptionsService{}

// SubscriptionService mocks the subscription service.
type SubscriptionsService struct {
	OnCreate      func(sub recurly.NewSubscription) (*recurly.Response, *recurly.Subscription, error)
	CreateInvoked bool

	OnCancel      func(uuid string) (*recurly.Response, *recurly.Subscription, error)
	CancelInvoked bool
}

func (m *SubscriptionsService) Create(sub recurly.NewSubscription) (*recurly.Response, *recurly.Subscription, error) {
	m.CreateInvoked = true
	return m.OnCreate(sub)
}

func (m *SubscriptionsService) Cancel(uuid string) (*recurly.Response, *recurly.Subscription, error) {
	m.CancelInvoked = true
	return m.OnCancel(uuid)
}
