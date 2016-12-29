package recurly

import "encoding/xml"

// Redemption holds redeemed coupons for an account or invoice.
type Redemption struct {
	CouponCode             string
	AccountCode            string
	SingleUse              NullBool
	TotalDiscountedInCents int
	Currency               string
	State                  string
	CreatedAt              NullTime
}

// UnmarshalXML unmarshal a coupon redemption object. Minaly converts href links
// for coupons and accounts to CouponCode and AccountCodes.
func (r *Redemption) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v struct {
		XMLName                xml.Name   `xml:"redemption"`
		CouponCode             HrefString `xml:"coupon,omitempty"`
		AccountCode            HrefString `xml:"account,omitempty"`
		SingleUse              NullBool   `xml:"single_use,omitempty"`
		TotalDiscountedInCents int        `xml:"total_discounted_in_cents,omitempty"`
		Currency               string     `xml:"currency,omitempty"`
		State                  string     `xml:"state,omitempty"`
		CreatedAt              NullTime   `xml:"created_at,omitempty"`
	}
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	*r = Redemption{
		CouponCode:             string(v.CouponCode),
		AccountCode:            string(v.AccountCode),
		SingleUse:              v.SingleUse,
		TotalDiscountedInCents: v.TotalDiscountedInCents,
		Currency:               v.Currency,
		State:                  v.State,
		CreatedAt:              v.CreatedAt,
	}

	return nil
}
