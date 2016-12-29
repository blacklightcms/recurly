package recurly

import (
	"encoding/xml"

	"github.com/blacklightcms/go-recurly/types"
)

type (
	// Redemption holds redeemed coupons for an account or invoice.
	Redemption struct {
		CouponCode             string
		AccountCode            string
		SingleUse              types.NullBool
		TotalDiscountedInCents int
		Currency               string
		State                  string
		CreatedAt              types.NullTime
	}
)

// UnmarshalXML unmarshal a coupon redemption object. Minaly converts href links
// for coupons and accounts to CouponCode and AccountCodes.
func (r *Redemption) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v struct {
		XMLName                xml.Name         `xml:"redemption"`
		CouponCode             types.HrefString `xml:"coupon,omitempty"`
		AccountCode            types.HrefString `xml:"account,omitempty"`
		SingleUse              types.NullBool   `xml:"single_use,omitempty"`
		TotalDiscountedInCents int              `xml:"total_discounted_in_cents,omitempty"`
		Currency               string           `xml:"currency,omitempty"`
		State                  string           `xml:"state,omitempty"`
		CreatedAt              types.NullTime   `xml:"created_at,omitempty"`
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
