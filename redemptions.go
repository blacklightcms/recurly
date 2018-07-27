package recurly

import "encoding/xml"

const (
	RedemptionStateActive   = "active"
	RedemptionStateInactive = "inactive"
)

// Redemption holds redeemed coupons for an account or invoice.
type Redemption struct {
	UUID                   string
	SubscriptionUUID       string
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
		XMLName                xml.Name   `xml:"redemption"`
		AccountCode            hrefString `xml:"account"`
		SubscriptionUUID       hrefString `xml:"subscription"`
		UUID                   string     `xml:"uuid"`
		SingleUse              bool       `xml:"single_use"`
		CouponCode             string     `xml:"coupon_code"`
		TotalDiscountedInCents int        `xml:"total_discounted_in_cents"`
		Currency               string     `xml:"currency,omitempty"`
		State                  string     `xml:"state,omitempty"`
		CreatedAt              NullTime   `xml:"created_at,omitempty"`
		UpdatedAt              NullTime   `xml:"updated_at,omitempty"`
	}
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	*r = Redemption{
		UUID:                   v.UUID,
		SubscriptionUUID:       string(v.SubscriptionUUID),
		AccountCode:            string(v.AccountCode),
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
