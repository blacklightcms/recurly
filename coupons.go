package recurly

import "encoding/xml"

// Coupon represents an individual coupon on your site.
type Coupon struct {
	XMLName            xml.Name          `xml:"coupon"`
	Code               string            `xml:"coupon_code"`
	Name               string            `xml:"name"`
	HostedDescription  string            `xml:"hosted_description,omitempty"`
	InvoiceDescription string            `xml:"invoice_description,omitempty"`
	State              string            `xml:"state,omitempty"`
	DiscountType       string            `xml:"discount_type"`
	DiscountPercent    int               `xml:"discount_percent,omitempty"`
	DiscountInCents    int               `xml:"discount_in_cents,omitempty"`
	RedeemByDate       NullTime          `xml:"redeem_by_date,omitempty"`
	SingleUse          NullBool          `xml:"single_use,omitempty"`
	AppliesForMonths   NullInt           `xml:"applies_for_months,omitempty"`
	MaxRedemptions     NullInt           `xml:"max_redemptions,omitempty"`
	AppliesToAllPlans  NullBool          `xml:"applies_to_all_plans,omitempty"`
	CreatedAt          NullTime          `xml:"created_at,omitempty"`
	PlanCodes          *[]CouponPlanCode `xml:"plan_codes>plan_code,omitempty"`
}

// CouponPlanCode holds an xml array of plan_code items that this coupon
// will work with.
type CouponPlanCode struct {
	Code string `xml:",innerxml"`
}
