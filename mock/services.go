package mock

import (
	"time"

	"github.com/blacklightcms/recurly"
)

var _ recurly.TransactionsService = &TransactionsService{}

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

var _ recurly.SubscriptionsService = &SubscriptionsService{}

// SubscriptionService mocks the subscription service.
type SubscriptionsService struct {
	OnList      func(params recurly.Params) (*recurly.Response, []recurly.Subscription, error)
	ListInvoked bool

	OnListAccount      func(accountCode string, params recurly.Params) (*recurly.Response, []recurly.Subscription, error)
	ListAccountInvoked bool

	OnGet      func(uuid string) (*recurly.Response, *recurly.Subscription, error)
	GetInvoked bool

	OnCreate      func(sub recurly.NewSubscription) (*recurly.Response, *recurly.Subscription, error)
	CreateInvoked bool

	OnPreview      func(sub recurly.NewSubscription) (*recurly.Response, *recurly.Subscription, error)
	PreviewInvoked bool

	OnUpdate      func(uuid string, sub recurly.UpdateSubscription) (*recurly.Response, *recurly.Subscription, error)
	UpdateInvoked bool

	OnUpdateNotes      func(uuid string, n recurly.SubscriptionNotes) (*recurly.Response, *recurly.Subscription, error)
	UpdateNotesInvoked bool

	OnPreviewChange      func(uuid string, sub recurly.UpdateSubscription) (*recurly.Response, *recurly.Subscription, error)
	PreviewChangeInvoked bool

	OnCancel      func(uuid string) (*recurly.Response, *recurly.Subscription, error)
	CancelInvoked bool

	OnReactivate      func(uuid string) (*recurly.Response, *recurly.Subscription, error)
	ReactivateInvoked bool

	OnTerminateWithPartialRefund      func(uuid string) (*recurly.Response, *recurly.Subscription, error)
	TerminateWithPartialRefundInvoked bool

	OnTerminateWithFullRefund      func(uuid string) (*recurly.Response, *recurly.Subscription, error)
	TerminateWithFullRefundInvoked bool

	OnTerminateWithoutRefund      func(uuid string) (*recurly.Response, *recurly.Subscription, error)
	TerminateWithoutRefundInvoked bool

	OnPostpone      func(uuid string, dt time.Time, bulk bool) (*recurly.Response, *recurly.Subscription, error)
	PostponeInvoked bool
}

func (m *SubscriptionsService) List(params recurly.Params) (*recurly.Response, []recurly.Subscription, error) {
	m.ListInvoked = true
	return m.OnList(params)
}

func (m *SubscriptionsService) ListAccount(accountCode string, params recurly.Params) (*recurly.Response, []recurly.Subscription, error) {
	m.ListAccountInvoked = true
	return m.OnListAccount(accountCode, params)
}

func (m *SubscriptionsService) Get(uuid string) (*recurly.Response, *recurly.Subscription, error) {
	m.GetInvoked = true
	return m.OnGet(uuid)
}

func (m *SubscriptionsService) Create(sub recurly.NewSubscription) (*recurly.Response, *recurly.Subscription, error) {
	m.CreateInvoked = true
	return m.OnCreate(sub)
}

func (m *SubscriptionsService) Preview(sub recurly.NewSubscription) (*recurly.Response, *recurly.Subscription, error) {
	m.PreviewInvoked = true
	return m.OnPreview(sub)
}

func (m *SubscriptionsService) Update(uuid string, sub recurly.UpdateSubscription) (*recurly.Response, *recurly.Subscription, error) {
	m.UpdateInvoked = true
	return m.OnUpdate(uuid, sub)
}

func (m *SubscriptionsService) UpdateNotes(uuid string, n recurly.SubscriptionNotes) (*recurly.Response, *recurly.Subscription, error) {
	m.UpdateNotesInvoked = true
	return m.OnUpdateNotes(uuid, n)
}

func (m *SubscriptionsService) PreviewChange(uuid string, sub recurly.UpdateSubscription) (*recurly.Response, *recurly.Subscription, error) {
	m.PreviewChangeInvoked = true
	return m.OnPreviewChange(uuid, sub)
}

func (m *SubscriptionsService) Cancel(uuid string) (*recurly.Response, *recurly.Subscription, error) {
	m.CancelInvoked = true
	return m.OnCancel(uuid)
}

func (m *SubscriptionsService) Reactivate(uuid string) (*recurly.Response, *recurly.Subscription, error) {
	m.ReactivateInvoked = true
	return m.OnReactivate(uuid)
}

func (m *SubscriptionsService) TerminateWithPartialRefund(uuid string) (*recurly.Response, *recurly.Subscription, error) {
	m.TerminateWithPartialRefundInvoked = true
	return m.OnTerminateWithPartialRefund(uuid)
}

func (m *SubscriptionsService) TerminateWithFullRefund(uuid string) (*recurly.Response, *recurly.Subscription, error) {
	m.TerminateWithFullRefundInvoked = true
	return m.OnTerminateWithFullRefund(uuid)
}

func (m *SubscriptionsService) TerminateWithoutRefund(uuid string) (*recurly.Response, *recurly.Subscription, error) {
	m.TerminateWithoutRefundInvoked = true
	return m.OnTerminateWithoutRefund(uuid)
}

func (m *SubscriptionsService) Postpone(uuid string, dt time.Time, bulk bool) (*recurly.Response, *recurly.Subscription, error) {
	m.PostponeInvoked = true
	return m.OnPostpone(uuid, dt, bulk)
}
