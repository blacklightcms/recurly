package recurly

import "encoding/xml"

// Redemption holds redeemed coupons for an account or invoice.
type Redemption struct {
	XMLName                xml.Name `xml:"redemption"`
	CouponCode             string   `xml:"coupon_code,omitempty"`
	SingleUse              NullBool `xml:"single_use,omitempty"`
	TotalDiscountedInCents int      `xml:"total_discounted_in_cents,omitempty"`
	Currency               string   `xml:"currency,omitempty"`
	State                  string   `xml:"state,omitempty"`
	CreatedAt              NullTime `xml:"created_at,omitempty"`
}
