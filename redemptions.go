package recurly

import (
	"context"
	"encoding/xml"
	"fmt"
)

// RedemptionsService manages the interactions for coupon redemptions.
type RedemptionsService interface {
	// ListAccount returns a pager to paginate redemptions for an account.
	// PagerOptions are used to optionally filter the results.
	//
	// https://dev.recurly.com/docs/coupon-redemption-object
	ListAccount(accountCode string, opts *PagerOptions) Pager

	// ListInvoice returns a pager to paginate redemptions for an invoice.
	// PagerOptions are used to optionally filter the results.
	//
	// https://dev.recurly.com/docs/lookup-a-coupon-redemption-on-an-invoice
	ListInvoice(invoiceNumber int, opts *PagerOptions) Pager

	// ListInvoice returns a pager to paginate redemptions for an invoice.
	// PagerOptions are used to optionally filter the results.
	//
	// https://dev.recurly.com/docs/lookup-a-coupon-redemption-on-a-subscription
	ListSubscription(uuid string, opts *PagerOptions) Pager

	// Redeem redeems a coupon on an existing customer's account to apply to
	// their next invoice. r.AccountCode and r.Currency are required fields.
	// Set r.SubscriptionUUID to redeem the coupon to a subscription.
	//
	// NOTE: If you want the coupon redemption to be rejected if a subscription
	// signup fails, you must redeem the coupon within the Subscriptions.Create()
	// call.
	//
	// https://dev.recurly.com/docs/redeem-a-coupon-before-or-after-a-subscription
	// https://dev.recurly.com/docs/redeem-a-coupon-before-or-after-a-subscription
	Redeem(ctx context.Context, code string, r CouponRedemption) (*Redemption, error)

	// Delete manually expires a coupon redemption on an account. Please note:
	// the coupon redemption will still count towards the "maximum redemption total"
	// of the coupon. See Recurly's documentation for details.
	//
	// https://dev.recurly.com/docs/remove-a-coupon-from-an-account
	Delete(ctx context.Context, accountCode, redemptionUUID string) error
}

// Redemptions constants.
const (
	RedemptionStateActive   = "active"
	RedemptionStateInactive = "inactive"
)

// Redemption holds redeemed coupons for an account or invoice.
//
// https://dev.recurly.com/docs/coupon-redemption-object
type Redemption struct {
	UUID                   string
	SubscriptionUUID       string // Only available if redeemed on a subscription
	AccountCode            string
	CouponCode             string
	SingleUse              bool
	TotalDiscountedInCents int
	Currency               string
	State                  string
	CreatedAt              NullTime
	UpdatedAt              NullTime
}

// UnmarshalXML unmarshal a coupon redemption object. Minaly converts href links
// for coupons and accounts to CouponCode and AccountCodes.
func (r *Redemption) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v struct {
		XMLName                xml.Name `xml:"redemption"`
		AccountCode            href     `xml:"account"`
		SubscriptionUUID       href     `xml:"subscription"`
		UUID                   string   `xml:"uuid"`
		SingleUse              bool     `xml:"single_use"`
		CouponCode             string   `xml:"coupon_code"`
		TotalDiscountedInCents int      `xml:"total_discounted_in_cents"`
		Currency               string   `xml:"currency,omitempty"`
		State                  string   `xml:"state,omitempty"`
		CreatedAt              NullTime `xml:"created_at,omitempty"`
		UpdatedAt              NullTime `xml:"updated_at,omitempty"`
	}
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	*r = Redemption{
		UUID:                   v.UUID,
		SubscriptionUUID:       v.SubscriptionUUID.LastPartOfPath(),
		AccountCode:            v.AccountCode.LastPartOfPath(),
		CouponCode:             v.CouponCode,
		SingleUse:              v.SingleUse,
		TotalDiscountedInCents: v.TotalDiscountedInCents,
		Currency:               v.Currency,
		State:                  v.State,
		CreatedAt:              v.CreatedAt,
		UpdatedAt:              v.UpdatedAt,
	}
	return nil
}

// CouponRedemption is used to redeem coupons.
type CouponRedemption struct {
	XMLName          xml.Name `xml:"redemption"`
	AccountCode      string   `xml:"account_code"`                // required
	Currency         string   `xml:"currency"`                    // required
	SubscriptionUUID string   `xml:"subscription_uuid,omitempty"` // optional, redeem to subscription
}

var _ RedemptionsService = &redemptionsImpl{}

// redemptionsImpl implements RedemptionsService.
type redemptionsImpl serviceImpl

func (s *redemptionsImpl) ListAccount(accountCode string, opts *PagerOptions) Pager {
	path := fmt.Sprintf("/accounts/%s/redemptions", accountCode)
	return s.client.newPager("GET", path, opts)
}

func (s *redemptionsImpl) ListInvoice(invoiceNumber int, opts *PagerOptions) Pager {
	path := fmt.Sprintf("/invoices/%d/redemptions", invoiceNumber)
	return s.client.newPager("GET", path, opts)
}

func (s *redemptionsImpl) ListSubscription(uuid string, opts *PagerOptions) Pager {
	path := fmt.Sprintf("/subscriptions/%s/redemptions", sanitizeUUID(uuid))
	return s.client.newPager("GET", path, opts)
}

func (s *redemptionsImpl) Redeem(ctx context.Context, code string, r CouponRedemption) (*Redemption, error) {
	if r.SubscriptionUUID != "" {
		r.SubscriptionUUID = sanitizeUUID(r.SubscriptionUUID)
	}

	path := fmt.Sprintf("/coupons/%s/redeem", code)
	req, err := s.client.newRequest("POST", path, r)
	if err != nil {
		return nil, err
	}

	var dst Redemption
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *redemptionsImpl) Delete(ctx context.Context, accountCode, redemptionUUID string) error {
	path := fmt.Sprintf("/accounts/%s/redemptions/%s", accountCode, redemptionUUID)
	req, err := s.client.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(ctx, req, nil)
	return err
}
