package mock

import (
	"context"
	"github.com/blacklightcms/recurly"
)

var _ recurly.AddOnUsageService = &AddOnUsageService{}

type AddOnUsageService struct {
	OnList      func(uuid, addOnCode string, opts *recurly.PagerOptions) recurly.Pager
	ListInvoked bool

	OnGet      func(ctx context.Context, uuid, addOnCode, usageId string) (*recurly.AddOnUsage, error)
	GetInvoked bool

	OnCreate      func(ctx context.Context, uuid, addOnCode string, usage recurly.AddOnUsage) (*recurly.AddOnUsage, error)
	CreateInvoked bool

	OnUpdate      func(ctx context.Context, uuid, addOnCode string, usage recurly.AddOnUsage) (*recurly.AddOnUsage, error)
	UpdateInvoked bool

	OnDelete      func(ctx context.Context, uuid, addOnCode, usageId string) error
	DeleteInvoked bool
}

func (m *AddOnUsageService) List(uuid, addOnCode string, opts *recurly.PagerOptions) recurly.Pager {
	m.ListInvoked = true
	return m.OnList(uuid, addOnCode, opts)
}

func (m *AddOnUsageService) Get(ctx context.Context, uuid, addOnCode, usageId string) (*recurly.AddOnUsage, error) {
	m.GetInvoked = true
	return m.OnGet(ctx, uuid, addOnCode, usageId)
}

func (m *AddOnUsageService) Create(ctx context.Context, uuid, addOnCode string, usage recurly.AddOnUsage) (*recurly.AddOnUsage, error) {
	m.CreateInvoked = true
	return m.OnCreate(ctx, uuid, addOnCode, usage)
}

func (m *AddOnUsageService) Update(ctx context.Context, uuid, addOnCode, usageId string, usage recurly.AddOnUsage) (*recurly.AddOnUsage, error) {
	m.UpdateInvoked = true
	return m.OnUpdate(ctx, uuid, addOnCode, usage)
}

func (m *AddOnUsageService) Delete(ctx context.Context, uuid, addOnCode, usageId string) error {
	m.DeleteInvoked = true
	return m.OnDelete(ctx, uuid, addOnCode, usageId)
}
