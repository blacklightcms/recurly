package recurly

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

var _ CouponsService = &couponsImpl{}

// couponsImpl handles communication with the coupons related methods
// of the recurly API.
type couponsImpl struct {
	client *Client
}

// List returns a list of all the coupons on your site.
// https://dev.recurly.com/docs/list-active-coupons
func (s *couponsImpl) List(params Params) (*Response, []Coupon, error) {
	req, err := s.client.newRequest("GET", "coupons", params, nil)
	if err != nil {
		return nil, nil, err
	}

	var c struct {
		XMLName xml.Name `xml:"coupons"`
		Coupons []Coupon `xml:"coupon"`
	}
	resp, err := s.client.do(req, &c)

	return resp, c.Coupons, err
}

// Get returns information about an active coupon.
// https://dev.recurly.com/docs/lookup-a-coupon
func (s *couponsImpl) Get(code string) (*Response, *Coupon, error) {
	action := fmt.Sprintf("coupons/%s", code)
	req, err := s.client.newRequest("GET", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst Coupon
	resp, err := s.client.do(req, &dst)
	if err != nil || resp.StatusCode >= http.StatusBadRequest {
		return resp, nil, err
	}

	return resp, &dst, err
}

// Create a new coupon. Coupons cannot be updated after being created.
// https://dev.recurly.com/docs/create-coupon
func (s *couponsImpl) Create(c Coupon) (*Response, *Coupon, error) {
	req, err := s.client.newRequest("POST", "coupons", nil, c)
	if err != nil {
		return nil, nil, err
	}

	var dst Coupon
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// Delete deactivates the coupon so it can no longer be redeemed.
// https://docs.recurly.com/api/plans/add-ons#delete-addon
func (s *couponsImpl) Delete(code string) (*Response, error) {
	action := fmt.Sprintf("coupons/%s", code)
	req, err := s.client.newRequest("DELETE", action, nil, nil)
	if err != nil {
		return nil, err
	}

	return s.client.do(req, nil)
}
