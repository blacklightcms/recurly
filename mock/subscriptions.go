package mock

import (
	"context"
	"time"

	"github.com/blacklightcms/recurly"
)

var _ recurly.SubscriptionsService = &SubscriptionsService{}

// SubscriptionsService mocks the subscription service.
type SubscriptionsService struct {
	OnList      func(opts *recurly.PagerOptions) recurly.Pager
	ListInvoked bool

	OnListAccount      func(accountCode string, opts *recurly.PagerOptions) recurly.Pager
	ListAccountInvoked bool

	OnGet      func(ctx context.Context, uuid string) (*recurly.Subscription, error)
	GetInvoked bool

	OnCreate      func(ctx context.Context, sub recurly.NewSubscription) (*recurly.Subscription, error)
	CreateInvoked bool

	OnPreview      func(ctx context.Context, sub recurly.NewSubscription) (*recurly.Subscription, error)
	PreviewInvoked bool

	OnUpdate      func(ctx context.Context, uuid string, sub recurly.UpdateSubscription) (*recurly.Subscription, error)
	UpdateInvoked bool

	OnUpdateNotes      func(ctx context.Context, uuid string, n recurly.SubscriptionNotes) (*recurly.Subscription, error)
	UpdateNotesInvoked bool

	OnPreviewChange      func(ctx context.Context, uuid string, sub recurly.UpdateSubscription) (*recurly.Subscription, error)
	PreviewChangeInvoked bool

	OnCancel      func(ctx context.Context, uuid string) (*recurly.Subscription, error)
	CancelInvoked bool

	OnReactivate      func(ctx context.Context, uuid string) (*recurly.Subscription, error)
	ReactivateInvoked bool

	OnTerminate      func(ctx context.Context, uuid string, refundType string) (*recurly.Subscription, error)
	TerminateInvoked bool

	OnPostpone      func(ctx context.Context, uuid string, dt time.Time, bulk bool) (*recurly.Subscription, error)
	PostponeInvoked bool

	OnPause      func(ctx context.Context, uuid string, cycles int) (*recurly.Subscription, error)
	PauseInvoked bool

	OnResume      func(ctx context.Context, uuid string) (*recurly.Subscription, error)
	ResumeInvoked bool

	OnConvertTrial func(ctx context.Context, uuid string) (*recurly.Subscription, error)
	ConvertTrialInvoked bool
}

func (m *SubscriptionsService) List(opts *recurly.PagerOptions) recurly.Pager {
	m.ListInvoked = true
	return m.OnList(opts)
}

func (m *SubscriptionsService) ListAccount(accountCode string, opts *recurly.PagerOptions) recurly.Pager {
	m.ListAccountInvoked = true
	return m.OnListAccount(accountCode, opts)
}

func (m *SubscriptionsService) Get(ctx context.Context, uuid string) (*recurly.Subscription, error) {
	m.GetInvoked = true
	return m.OnGet(ctx, uuid)
}

func (m *SubscriptionsService) Create(ctx context.Context, sub recurly.NewSubscription) (*recurly.Subscription, error) {
	m.CreateInvoked = true
	return m.OnCreate(ctx, sub)
}

func (m *SubscriptionsService) Preview(ctx context.Context, sub recurly.NewSubscription) (*recurly.Subscription, error) {
	m.PreviewInvoked = true
	return m.OnPreview(ctx, sub)
}

func (m *SubscriptionsService) Update(ctx context.Context, uuid string, sub recurly.UpdateSubscription) (*recurly.Subscription, error) {
	m.UpdateInvoked = true
	return m.OnUpdate(ctx, uuid, sub)
}

func (m *SubscriptionsService) UpdateNotes(ctx context.Context, uuid string, n recurly.SubscriptionNotes) (*recurly.Subscription, error) {
	m.UpdateNotesInvoked = true
	return m.OnUpdateNotes(ctx, uuid, n)
}

func (m *SubscriptionsService) PreviewChange(ctx context.Context, uuid string, sub recurly.UpdateSubscription) (*recurly.Subscription, error) {
	m.PreviewChangeInvoked = true
	return m.OnPreviewChange(ctx, uuid, sub)
}

func (m *SubscriptionsService) Cancel(ctx context.Context, uuid string) (*recurly.Subscription, error) {
	m.CancelInvoked = true
	return m.OnCancel(ctx, uuid)
}

func (m *SubscriptionsService) Reactivate(ctx context.Context, uuid string) (*recurly.Subscription, error) {
	m.ReactivateInvoked = true
	return m.OnReactivate(ctx, uuid)
}

func (m *SubscriptionsService) Terminate(ctx context.Context, uuid string, refundType string) (*recurly.Subscription, error) {
	m.TerminateInvoked = true
	return m.OnTerminate(ctx, uuid, refundType)
}

func (m *SubscriptionsService) Postpone(ctx context.Context, uuid string, dt time.Time, bulk bool) (*recurly.Subscription, error) {
	m.PostponeInvoked = true
	return m.OnPostpone(ctx, uuid, dt, bulk)
}

func (m *SubscriptionsService) Pause(ctx context.Context, uuid string, cycles int) (*recurly.Subscription, error) {
	m.PauseInvoked = true
	return m.OnPause(ctx, uuid, cycles)
}

func (m *SubscriptionsService) Resume(ctx context.Context, uuid string) (*recurly.Subscription, error) {
	m.ResumeInvoked = true
	return m.OnResume(ctx, uuid)
}

func (m *SubscriptionsService) ConvertTrial(ctx context.Context, uuid string) (*recurly.Subscription, error) {
	m.ConvertTrialInvoked = true
	return m.OnConvertTrial(ctx, uuid)
}