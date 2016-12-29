package recurly

import (
	"encoding/xml"

	"github.com/blacklightcms/go-recurly/types"
)

type (
	// Coupon represents an individual coupon on your site.
	Coupon struct {
		XMLName            xml.Name          `xml:"coupon"`
		Code               string            `xml:"coupon_code"`
		Name               string            `xml:"name"`
		HostedDescription  string            `xml:"hosted_description,omitempty"`
		InvoiceDescription string            `xml:"invoice_description,omitempty"`
		State              string            `xml:"state,omitempty"`
		DiscountType       string            `xml:"discount_type"`
		DiscountPercent    int               `xml:"discount_percent,omitempty"`
		DiscountInCents    int               `xml:"discount_in_cents,omitempty"`
		RedeemByDate       types.NullTime    `xml:"redeem_by_date,omitempty"`
		SingleUse          types.NullBool    `xml:"single_use,omitempty"`
		AppliesForMonths   types.NullInt     `xml:"applies_for_months,omitempty"`
		MaxRedemptions     types.NullInt     `xml:"max_redemptions,omitempty"`
		AppliesToAllPlans  types.NullBool    `xml:"applies_to_all_plans,omitempty"`
		CreatedAt          types.NullTime    `xml:"created_at,omitempty"`
		PlanCodes          *[]CouponPlanCode `xml:"plan_codes>plan_code,omitempty"`
	}

	// CouponPlanCode holds an xml array of plan_code items that this coupon
	// will work with.
	CouponPlanCode struct {
		Code string `xml:",innerxml"`
	}
)
