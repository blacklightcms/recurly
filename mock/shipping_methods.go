package mock

import (
	"context"

	"github.com/blacklightcms/recurly"
)

var _ recurly.ShippingMethodsService = &ShippingMethodsService{}

type ShippingMethodsService struct {
	OnList      func(opts *recurly.PagerOptions) recurly.Pager
	ListInvoked bool

	OnGet      func(ctx context.Context, code string) (*recurly.ShippingMethod, error)
	GetInvoked bool
}

func (s *ShippingMethodsService) List(opts *recurly.PagerOptions) recurly.Pager {
	s.ListInvoked = true
	return s.OnList(opts)
}

func (s *ShippingMethodsService) Get(ctx context.Context, code string) (*recurly.ShippingMethod, error) {
	s.GetInvoked = true
	return s.OnGet(ctx, code)
}
