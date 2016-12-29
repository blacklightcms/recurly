package mock

import (
	"net/http"

	"github.com/blacklightcms/go-recurly/recurly"
)

var _ recurly.Transactions = &MockTransactionsService{}

// MockTransactionsService mocks the transaction service.
type MockTransactionsService struct {
	OnList      func(params recurly.Params) (*recurly.Response, []recurly.Transaction, error)
	ListInvoked bool

	OnListAccount      func(accountCode string, params recurly.Params) (*recurly.Response, []recurly.Transaction, error)
	ListAccountInvoked bool

	OnGet      func(uuid string) (*recurly.Response, *recurly.Transaction, error)
	GetInvoked bool

	OnCreate      func(trans recurly.Transaction) (*recurly.Response, *recurly.Transaction, error)
	CreateInvoked bool
}

func (m *MockTransactionsService) List(params recurly.Params) (*recurly.Response, []recurly.Transaction, error) {
	m.ListInvoked = true
	return m.OnList(params)
}

func (m *MockTransactionsService) ListAccount(accountCode string, params recurly.Params) (*recurly.Response, []recurly.Transaction, error) {
	m.ListAccountInvoked = true
	return m.OnListAccount(accountCode, params)
}

func (m *MockTransactionsService) Get(uuid string) (*recurly.Response, *recurly.Transaction, error) {
	m.GetInvoked = true
	return m.OnGet(uuid)
}

func (m *MockTransactionsService) Create(t recurly.Transaction) (*recurly.Response, *recurly.Transaction, error) {
	m.CreateInvoked = true
	return m.OnCreate(t)
}

var _ recurly.Subscriptions = &MockSubscriptionsService{}

// MockSubscriptionService mocks the subscription service.
type MockSubscriptionsService struct {
	OnCreate      func(sub recurly.NewSubscription) (*recurly.Response, *recurly.Subscription, error)
	CreateInvoked bool

	OnCancel      func(uuid string) (*recurly.Response, *recurly.Subscription, error)
	CancelInvoked bool
}

func (m *MockSubscriptionsService) Create(sub recurly.NewSubscription) (*recurly.Response, *recurly.Subscription, error) {
	m.CreateInvoked = true
	return m.OnCreate(sub)
}

func (m *MockSubscriptionsService) Cancel(uuid string) (*recurly.Response, *recurly.Subscription, error) {
	m.CancelInvoked = true
	return m.OnCancel(uuid)
}

// MockClient mocks the recurly client.
type MockClient struct {
	// client is the HTTP Client used to communicate with the API.
	client *http.Client

	// subdomain is your account's sub domain used for authentication.
	subDomain string

	// apiKey is your account's API key used for authentication.
	apiKey string

	// BaseURL is the base url for api requests.
	BaseURL string

	// Services used for talking with different parts of the Recurly API
	// Accounts      *AccountsService
	// Adjustments   *AdjustmentsService
	// Billing       *BillingService
	// Coupons       *CouponsService
	// Redemptions   *RedemptionsService
	// Invoices      *InvoicesService
	// Plans         *PlansService
	// AddOns        *AddOnsService
	Subscriptions *MockSubscriptionsService
	Transactions  *MockTransactionsService
}
