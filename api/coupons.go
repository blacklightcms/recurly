package api

import (
	"encoding/xml"
	"fmt"

	recurly "github.com/blacklightcms/go-recurly"
)

var _ recurly.CouponsService = &CouponsService{}

// CouponsService handles communication with the coupons related methods
// of the recurly API.
type CouponsService struct {
	client *recurly.Client
}

// List returns a list of all the coupons on your site.
// https://dev.recurly.com/docs/list-active-coupons
func (s *CouponsService) List(params recurly.Params) (*recurly.Response, []recurly.Coupon, error) {
	req, err := s.client.NewRequest("GET", "coupons", params, nil)
	if err != nil {
		return nil, nil, err
	}

	var c struct {
		XMLName xml.Name         `xml:"coupons"`
		Coupons []recurly.Coupon `xml:"coupon"`
	}
	resp, err := s.client.Do(req, &c)

	return resp, c.Coupons, err
}

// Get returns information about an active coupon.
// https://dev.recurly.com/docs/lookup-a-coupon
func (s *CouponsService) Get(code string) (*recurly.Response, *recurly.Coupon, error) {
	action := fmt.Sprintf("coupons/%s", code)
	req, err := s.client.NewRequest("GET", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst recurly.Coupon
	resp, err := s.client.Do(req, &dst)

	return resp, &dst, err
}

// Create a new coupon. Coupons cannot be updated after being created.
// https://dev.recurly.com/docs/create-coupon
func (s *CouponsService) Create(c recurly.Coupon) (*recurly.Response, *recurly.Coupon, error) {
	req, err := s.client.NewRequest("POST", "coupons", nil, c)
	if err != nil {
		return nil, nil, err
	}

	var dst recurly.Coupon
	resp, err := s.client.Do(req, &dst)

	return resp, &dst, err
}

// Delete deactivates the coupon so it can no longer be redeemed.
// https://docs.recurly.com/api/plans/add-ons#delete-addon
func (s *CouponsService) Delete(code string) (*recurly.Response, error) {
	action := fmt.Sprintf("coupons/%s", code)
	req, err := s.client.NewRequest("DELETE", action, nil, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
