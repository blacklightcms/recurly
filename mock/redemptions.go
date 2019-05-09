package mock

import (
	"context"

	"github.com/blacklightcms/recurly"
)

var _ recurly.RedemptionsService = &RedemptionsService{}

// RedemptionsService manages the interactions for redemptions.
type RedemptionsService struct {
	OnListAccount      func(accountCode string, opts *recurly.PagerOptions) *recurly.RedemptionsPager
	ListAccountInvoked bool

	OnListInvoice      func(invoiceNumber int, opts *recurly.PagerOptions) *recurly.RedemptionsPager
	ListInvoiceInvoked bool

	OnListSubscription      func(uuid string, opts *recurly.PagerOptions) *recurly.RedemptionsPager
	ListSubscriptionInvoked bool

	OnRedeem      func(ctx context.Context, code string, r recurly.CouponRedemption) (*recurly.Redemption, error)
	RedeemInvoked bool

	OnDelete      func(ctx context.Context, accountCode string) error
	DeleteInvoked bool
}

func (m *RedemptionsService) ListAccount(accountCode string, opts *recurly.PagerOptions) *recurly.RedemptionsPager {
	m.ListAccountInvoked = true
	return m.OnListAccount(accountCode, opts)
}

func (m *RedemptionsService) ListInvoice(invoiceNumber int, opts *recurly.PagerOptions) *recurly.RedemptionsPager {
	m.ListInvoiceInvoked = true
	return m.OnListInvoice(invoiceNumber, opts)
}

func (m *RedemptionsService) ListSubscription(uuid string, opts *recurly.PagerOptions) *recurly.RedemptionsPager {
	m.ListSubscriptionInvoked = true
	return m.OnListSubscription(uuid, opts)
}

func (m *RedemptionsService) Redeem(ctx context.Context, code string, r recurly.CouponRedemption) (*recurly.Redemption, error) {
	m.RedeemInvoked = true
	return m.OnRedeem(ctx, code, r)
}

func (m *RedemptionsService) Delete(ctx context.Context, accountCode string) error {
	m.DeleteInvoked = true
	return m.OnDelete(ctx, accountCode)
}
