package mock

import (
	"bytes"
	"time"

	"github.com/blacklightcms/recurly"
)

var _ recurly.AccountsService = &AccountsService{}

// AccountsService represents the interactions available for accounts.
type AccountsService struct {
	OnList      func(params recurly.Params) (*recurly.Response, []recurly.Account, error)
	ListInvoked bool

	OnGet      func(code string) (*recurly.Response, *recurly.Account, error)
	GetInvoked bool

	OnCreate      func(a recurly.Account) (*recurly.Response, *recurly.Account, error)
	CreateInvoked bool

	OnUpdate      func(code string, a recurly.Account) (*recurly.Response, *recurly.Account, error)
	UpdateInvoked bool

	OnClose      func(code string) (*recurly.Response, error)
	CloseInvoked bool

	OnReopen      func(code string) (*recurly.Response, error)
	ReopenInvoked bool

	OnListNotes      func(code string) (*recurly.Response, []recurly.Note, error)
	ListNotesInvoked bool
}

func (m *AccountsService) List(params recurly.Params) (*recurly.Response, []recurly.Account, error) {
	m.ListInvoked = true
	return m.OnList(params)
}

func (m *AccountsService) Get(code string) (*recurly.Response, *recurly.Account, error) {
	m.GetInvoked = true
	return m.OnGet(code)
}

func (m *AccountsService) Create(a recurly.Account) (*recurly.Response, *recurly.Account, error) {
	m.CreateInvoked = true
	return m.OnCreate(a)
}

func (m *AccountsService) Update(code string, a recurly.Account) (*recurly.Response, *recurly.Account, error) {
	m.UpdateInvoked = true
	return m.OnUpdate(code, a)
}

func (m *AccountsService) Close(code string) (*recurly.Response, error) {
	m.CloseInvoked = true
	return m.OnClose(code)
}

func (m *AccountsService) Reopen(code string) (*recurly.Response, error) {
	m.ReopenInvoked = true
	return m.OnReopen(code)
}

func (m *AccountsService) ListNotes(code string) (*recurly.Response, []recurly.Note, error) {
	m.ListNotesInvoked = true
	return m.OnListNotes(code)
}

var _ recurly.AdjustmentsService = &AdjustmentsService{}

// AdjustmentsService represents the interactions available for adjustments.
type AdjustmentsService struct {
	OnList      func(accountCode string, params recurly.Params) (*recurly.Response, []recurly.Adjustment, error)
	ListInvoked bool

	OnGet      func(uuid string) (*recurly.Response, *recurly.Adjustment, error)
	GetInvoked bool

	OnCreate      func(accountCode string, a recurly.Adjustment) (*recurly.Response, *recurly.Adjustment, error)
	CreateInvoked bool

	OnDelete      func(uuid string) (*recurly.Response, error)
	DeleteInvoked bool
}

func (m *AdjustmentsService) List(accountCode string, params recurly.Params) (*recurly.Response, []recurly.Adjustment, error) {
	m.ListInvoked = true
	return m.OnList(accountCode, params)
}

func (m *AdjustmentsService) Get(uuid string) (*recurly.Response, *recurly.Adjustment, error) {
	m.GetInvoked = true
	return m.OnGet(uuid)
}

func (m *AdjustmentsService) Create(accountCode string, a recurly.Adjustment) (*recurly.Response, *recurly.Adjustment, error) {
	m.CreateInvoked = true
	return m.OnCreate(accountCode, a)
}

func (m *AdjustmentsService) Delete(uuid string) (*recurly.Response, error) {
	m.DeleteInvoked = true
	return m.OnDelete(uuid)
}

var _ recurly.AddOnsService = &AddOnsService{}

// AddOnsService represents the interactions available for add ons.
type AddOnsService struct {
	OnList      func(planCode string, params recurly.Params) (*recurly.Response, []recurly.AddOn, error)
	ListInvoked bool

	OnGet      func(planCode string, code string) (*recurly.Response, *recurly.AddOn, error)
	GetInvoked bool

	OnCreate      func(planCode string, a recurly.AddOn) (*recurly.Response, *recurly.AddOn, error)
	CreateInvoked bool

	OnUpdate      func(planCode string, code string, a recurly.AddOn) (*recurly.Response, *recurly.AddOn, error)
	UpdateInvoked bool

	OnDelete      func(planCode string, code string) (*recurly.Response, error)
	DeleteInvoked bool
}

func (m *AddOnsService) List(planCode string, params recurly.Params) (*recurly.Response, []recurly.AddOn, error) {
	m.ListInvoked = true
	return m.OnList(planCode, params)
}

func (m *AddOnsService) Get(planCode string, code string) (*recurly.Response, *recurly.AddOn, error) {
	m.GetInvoked = true
	return m.OnGet(planCode, code)
}

func (m *AddOnsService) Create(planCode string, a recurly.AddOn) (*recurly.Response, *recurly.AddOn, error) {
	m.CreateInvoked = true
	return m.OnCreate(planCode, a)
}

func (m *AddOnsService) Update(planCode string, code string, a recurly.AddOn) (*recurly.Response, *recurly.AddOn, error) {
	m.UpdateInvoked = true
	return m.OnUpdate(planCode, code, a)
}

func (m *AddOnsService) Delete(planCode string, code string) (*recurly.Response, error) {
	m.DeleteInvoked = true
	return m.OnDelete(planCode, code)
}

var _ recurly.BillingService = &BillingService{}

// BillingService represents the interactions available for billing.
type BillingService struct {
	OnGet      func(accountCode string) (*recurly.Response, *recurly.Billing, error)
	GetInvoked bool

	OnCreate      func(accountCode string, b recurly.Billing) (*recurly.Response, *recurly.Billing, error)
	CreateInvoked bool

	OnCreateWithToken      func(accountCode string, token string) (*recurly.Response, *recurly.Billing, error)
	CreateWithTokenInvoked bool

	OnUpdate      func(accountCode string, b recurly.Billing) (*recurly.Response, *recurly.Billing, error)
	UpdateInvoked bool

	OnUpdateWithToken      func(accountCode string, token string) (*recurly.Response, *recurly.Billing, error)
	UpdateWithTokenInvoked bool

	OnClear      func(accountCode string) (*recurly.Response, error)
	ClearInvoked bool
}

func (m *BillingService) Get(accountCode string) (*recurly.Response, *recurly.Billing, error) {
	m.GetInvoked = true
	return m.OnGet(accountCode)
}

func (m *BillingService) Create(accountCode string, b recurly.Billing) (*recurly.Response, *recurly.Billing, error) {
	m.CreateInvoked = true
	return m.OnCreate(accountCode, b)
}

func (m *BillingService) CreateWithToken(accountCode string, token string) (*recurly.Response, *recurly.Billing, error) {
	m.CreateWithTokenInvoked = true
	return m.OnCreateWithToken(accountCode, token)
}

func (m *BillingService) Update(accountCode string, b recurly.Billing) (*recurly.Response, *recurly.Billing, error) {
	m.UpdateInvoked = true
	return m.OnUpdate(accountCode, b)
}

func (m *BillingService) UpdateWithToken(accountCode string, token string) (*recurly.Response, *recurly.Billing, error) {
	m.UpdateWithTokenInvoked = true
	return m.OnUpdateWithToken(accountCode, token)
}

func (m *BillingService) Clear(accountCode string) (*recurly.Response, error) {
	m.ClearInvoked = true
	return m.OnClear(accountCode)
}

var _ recurly.CouponsService = &CouponsService{}

// CouponsService represents the interactions available for coupons.
type CouponsService struct {
	OnList      func(params recurly.Params) (*recurly.Response, []recurly.Coupon, error)
	ListInvoked bool

	OnGet      func(code string) (*recurly.Response, *recurly.Coupon, error)
	GetInvoked bool

	OnCreate      func(c recurly.Coupon) (*recurly.Response, *recurly.Coupon, error)
	CreateInvoked bool

	OnDelete      func(code string) (*recurly.Response, error)
	DeleteInvoked bool
}

func (m *CouponsService) List(params recurly.Params) (*recurly.Response, []recurly.Coupon, error) {
	m.ListInvoked = true
	return m.OnList(params)
}

func (m *CouponsService) Get(code string) (*recurly.Response, *recurly.Coupon, error) {
	m.GetInvoked = true
	return m.OnGet(code)
}

func (m *CouponsService) Create(c recurly.Coupon) (*recurly.Response, *recurly.Coupon, error) {
	m.CreateInvoked = true
	return m.OnCreate(c)
}

func (m *CouponsService) Delete(code string) (*recurly.Response, error) {
	m.DeleteInvoked = true
	return m.OnDelete(code)
}

var _ recurly.InvoicesService = &InvoicesService{}

// InvoicesService represents the interactions available for invoices.
type InvoicesService struct {
	OnList      func(params recurly.Params) (*recurly.Response, []recurly.Invoice, error)
	ListInvoked bool

	OnListAccount      func(accountCode string, params recurly.Params) (*recurly.Response, []recurly.Invoice, error)
	ListAccountInvoked bool

	OnGet      func(invoiceNumber int) (*recurly.Response, *recurly.Invoice, error)
	GetInvoked bool

	OnGetPDF      func(invoiceNumber int, language string) (*recurly.Response, *bytes.Buffer, error)
	GetPDFInvoked bool

	OnPreview      func(accountCode string) (*recurly.Response, *recurly.Invoice, error)
	PreviewInvoked bool

	OnCreate      func(accountCode string, invoice recurly.Invoice) (*recurly.Response, *recurly.Invoice, error)
	CreateInvoked bool

	OnMarkPaid      func(invoiceNumber int) (*recurly.Response, *recurly.Invoice, error)
	MarkPaidInvoked bool

	OnMarkFailed      func(invoiceNumber int) (*recurly.Response, *recurly.Invoice, error)
	MarkFailedInvoked bool
}

func (m *InvoicesService) List(params recurly.Params) (*recurly.Response, []recurly.Invoice, error) {
	m.ListInvoked = true
	return m.OnList(params)
}

func (m *InvoicesService) ListAccount(accountCode string, params recurly.Params) (*recurly.Response, []recurly.Invoice, error) {
	m.ListAccountInvoked = true
	return m.OnListAccount(accountCode, params)
}

func (m *InvoicesService) Get(invoiceNumber int) (*recurly.Response, *recurly.Invoice, error) {
	m.GetInvoked = true
	return m.OnGet(invoiceNumber)
}

func (m *InvoicesService) GetPDF(invoiceNumber int, language string) (*recurly.Response, *bytes.Buffer, error) {
	m.GetPDFInvoked = true
	return m.OnGetPDF(invoiceNumber, language)
}

func (m *InvoicesService) Preview(accountCode string) (*recurly.Response, *recurly.Invoice, error) {
	m.PreviewInvoked = true
	return m.OnPreview(accountCode)
}

func (m *InvoicesService) Create(accountCode string, invoice recurly.Invoice) (*recurly.Response, *recurly.Invoice, error) {
	m.CreateInvoked = true
	return m.OnCreate(accountCode, invoice)
}

func (m *InvoicesService) MarkPaid(invoiceNumber int) (*recurly.Response, *recurly.Invoice, error) {
	m.MarkPaidInvoked = true
	return m.OnMarkPaid(invoiceNumber)
}

func (m *InvoicesService) MarkFailed(invoiceNumber int) (*recurly.Response, *recurly.Invoice, error) {
	m.MarkFailedInvoked = true
	return m.OnMarkFailed(invoiceNumber)
}

// PlansService represents the interactions available for plans.
type PlansService struct {
	OnList      func(params recurly.Params) (*recurly.Response, []recurly.Plan, error)
	ListInvoked bool

	OnGet      func(code string) (*recurly.Response, *recurly.Plan, error)
	GetInvoked bool

	OnCreate      func(p recurly.Plan) (*recurly.Response, *recurly.Plan, error)
	CreateInvoked bool

	OnUpdate      func(code string, p recurly.Plan) (*recurly.Response, *recurly.Plan, error)
	UpdateInvoked bool

	OnDelete      func(code string) (*recurly.Response, error)
	DeleteInvoked bool
}

func (m *PlansService) List(params recurly.Params) (*recurly.Response, []recurly.Plan, error) {
	m.ListInvoked = true
	return m.OnList(params)
}

func (m *PlansService) Get(code string) (*recurly.Response, *recurly.Plan, error) {
	m.GetInvoked = true
	return m.OnGet(code)
}

func (m *PlansService) Create(p recurly.Plan) (*recurly.Response, *recurly.Plan, error) {
	m.CreateInvoked = true
	return m.OnCreate(p)
}

func (m *PlansService) Update(code string, p recurly.Plan) (*recurly.Response, *recurly.Plan, error) {
	m.UpdateInvoked = true
	return m.OnUpdate(code, p)
}

func (m *PlansService) Delete(code string) (*recurly.Response, error) {
	m.DeleteInvoked = true
	return m.OnDelete(code)
}

var _ recurly.RedemptionsService = &RedemptionsService{}

// RedemptionsService represents the interactions available for redemptions.
type RedemptionsService struct {
	OnGetForAccount      func(accountCode string) (*recurly.Response, *recurly.Redemption, error)
	GetForAccountInvoked bool

	OnGetForInvoice      func(invoiceNumber string) (*recurly.Response, *recurly.Redemption, error)
	GetForInvoiceInvoked bool

	OnRedeem      func(code string, accountCode string, currency string) (*recurly.Response, *recurly.Redemption, error)
	RedeemInvoked bool

	OnDelete      func(accountCode string) (*recurly.Response, error)
	DeleteInvoked bool
}

func (m *RedemptionsService) GetForAccount(accountCode string) (*recurly.Response, *recurly.Redemption, error) {
	m.GetForAccountInvoked = true
	return m.OnGetForAccount(accountCode)
}

func (m *RedemptionsService) GetForInvoice(invoiceNumber string) (*recurly.Response, *recurly.Redemption, error) {
	m.GetForInvoiceInvoked = true
	return m.OnGetForInvoice(invoiceNumber)
}

func (m *RedemptionsService) Redeem(code string, accountCode string, currency string) (*recurly.Response, *recurly.Redemption, error) {
	m.RedeemInvoked = true
	return m.OnRedeem(code, accountCode, currency)
}

func (m *RedemptionsService) Delete(accountCode string) (*recurly.Response, error) {
	m.DeleteInvoked = true
	return m.OnDelete(accountCode)
}

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
