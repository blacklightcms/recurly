package recurly

import (
	"encoding/xml"
	"fmt"
)

type (
	// CouponsService handles communication with the coupons related methods
	// of the recurly API.
	CouponsService struct {
		client *Client
	}

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
	CouponPlanCode struct {
		Code string `xml:",innerxml"`
	}
)

// List returns a list of all the coupons
// https://dev.recurly.com/docs/list-active-coupons
func (service CouponsService) List(params Params) (*Response, []Coupon, error) {
	req, err := service.client.newRequest("GET", "coupons", params, nil)
	if err != nil {
		return nil, nil, err
	}

	var c struct {
		XMLName xml.Name `xml:"coupons"`
		Coupons []Coupon `xml:"coupon"`
	}
	res, err := service.client.do(req, &c)

	return res, c.Coupons, err
}

// Get returns information about an active coupon.
// https://dev.recurly.com/docs/lookup-a-coupon
func (service CouponsService) Get(code string) (*Response, Coupon, error) {
	action := fmt.Sprintf("coupons/%s", code)
	req, err := service.client.newRequest("GET", action, nil, nil)
	if err != nil {
		return nil, Coupon{}, err
	}

	var a Coupon
	res, err := service.client.do(req, &a)

	return res, a, err
}

// Create a new coupon. Coupons cannot be updated after being created.
// https://dev.recurly.com/docs/create-coupon
func (service CouponsService) Create(c Coupon) (*Response, Coupon, error) {
	req, err := service.client.newRequest("POST", "coupons", nil, c)
	if err != nil {
		return nil, Coupon{}, err
	}

	var dest Coupon
	res, err := service.client.do(req, &dest)

	return res, dest, err
}

// Delete deactivates the coupon so it can no longer be redeemed.
// https://docs.recurly.com/api/plans/add-ons#delete-addon
func (service CouponsService) Delete(code string) (*Response, error) {
	action := fmt.Sprintf("coupons/%s", code)
	req, err := service.client.newRequest("DELETE", action, nil, nil)
	if err != nil {
		return nil, err
	}

	return service.client.do(req, nil)
}
