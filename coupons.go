package recurly

import "encoding/xml"

// Coupon represents an individual coupon on your site.
// https://dev.recurly.com/docs/lookup-a-coupon
type Coupon struct {
	XMLName                  xml.Name    `xml:"coupon"`
	ID                       uint64      `xml:"id"`
	Code                     string      `xml:"coupon_code"`
	Type                     string      `xml:"coupon_type"`
	Name                     string      `xml:"name"`
	RedemptionResource       string      `xml:"redemption_resource"`
	State                    string      `xml:"state"`
	SingleUse                bool        `xml:"single_use"`
	AppliesToAllPlans        bool        `xml:"applies_to_all_plans"`
	Duration                 string      `xml:"duration"`
	DiscountType             string      `xml:"discount_type"`
	AppliesToNonPlanCharges  bool        `xml:"applies_to_non_plan_charges"`
	Description              string      `xml:"description,omitempty"`
	InvoiceDescription       string      `xml:"invoice_description,omitempty"`
	DiscountPercent          NullInt     `xml:"discount_percent,omitempty"`
	DiscountInCents          *UnitAmount `xml:"discount_in_cents,omitempty"`
	RedeemByDate             NullTime    `xml:"redeem_by_date,omitempty"`
	MaxRedemptions           NullInt     `xml:"max_redemptions,omitempty"`
	CreatedAt                NullTime    `xml:"created_at,omitempty"`
	UpdatedAt                NullTime    `xml:"updated_at,omitempty"`
	DeletedAt                NullTime    `xml:"deleted_at,omitempty"`
	TemporalUnit             string      `xml:"temporal_unit,omitempty"`
	TemporalAmount           NullInt     `xml:"temporal_amount,omitempty"`
	MaxRedemptionsPerAccount NullInt     `xml:"max_redemptions_per_account,omitempty"`
	UniqueCodeTemplate       string      `xml:"unique_code_template,omitempty"`
	UniqueCouponCodeCount    NullInt     `xml:"unique_coupon_codes_count,omitempty"`
	PlanCodes                []string    `xml:"plan_codes>plan_code,omitempty"`
	SubscriptionUUID         string      `xml:"subscription,omitempty"`
}
