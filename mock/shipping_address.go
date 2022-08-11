package mock

import (
	"context"

	"github.com/autopilot3/recurly"
)

var _ recurly.ShippingAddressesService = &ShippingAddressesService{}

type ShippingAddressesService struct {
	OnListAccount      func(accountCode string, opts *recurly.PagerOptions) recurly.Pager
	ListAccountInvoked bool

	OnCreate      func(ctx context.Context, accountCode string, address recurly.ShippingAddress) (*recurly.ShippingAddress, error)
	CreateInvoked bool

	OnUpdate      func(ctx context.Context, accountCode string, shippingAddressID int, address recurly.ShippingAddress) (*recurly.ShippingAddress, error)
	UpdateInvoked bool

	OnDelete      func(ctx context.Context, accountCode string, shippingAddressID int) error
	DeleteInvoked bool
}

func (s *ShippingAddressesService) ListAccount(accountCode string, opts *recurly.PagerOptions) recurly.Pager {
	s.ListAccountInvoked = true
	return s.OnListAccount(accountCode, opts)
}

func (s *ShippingAddressesService) Create(ctx context.Context, accountCode string, address recurly.ShippingAddress) (*recurly.ShippingAddress, error) {
	s.CreateInvoked = true
	return s.OnCreate(ctx, accountCode, address)
}

func (s *ShippingAddressesService) Update(ctx context.Context, accountCode string, shippingAddressID int, address recurly.ShippingAddress) (*recurly.ShippingAddress, error) {
	s.UpdateInvoked = true
	return s.OnUpdate(ctx, accountCode, shippingAddressID, address)
}

func (s *ShippingAddressesService) Delete(ctx context.Context, accountCode string, shippingAddressID int) error {
	s.DeleteInvoked = true
	return s.OnDelete(ctx, accountCode, shippingAddressID)
}
