package recurly

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

var _ PlansService = &plansImpl{}

// plansImpl handles communication with the plans related methods
// of the recurly API.
type plansImpl struct {
	client *Client
}

// List will retrieve all your active subscription plans.
// https://docs.recurly.com/api/plans#list-plans
func (s *plansImpl) List(params Params) (*Response, []Plan, error) {
	req, err := s.client.newRequest("GET", "plans", params, nil)
	if err != nil {
		return nil, nil, err
	}

	var p struct {
		XMLName xml.Name `xml:"plans"`
		Plans   []Plan   `xml:"plan"`
	}
	resp, err := s.client.do(req, &p)

	return resp, p.Plans, err
}

// Get will lookup a specific plan by code.
// https://docs.recurly.com/api/plans#lookup-plan
func (s *plansImpl) Get(code string) (*Response, *Plan, error) {
	action := fmt.Sprintf("plans/%s", code)
	req, err := s.client.newRequest("GET", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst Plan
	resp, err := s.client.do(req, &dst)
	if err != nil || resp.StatusCode >= http.StatusBadRequest {
		return resp, nil, err
	}

	return resp, &dst, err
}

// Create will create a new subscription plan.
// https://docs.recurly.com/api/plans#create-plan
func (s *plansImpl) Create(p Plan) (*Response, *Plan, error) {
	req, err := s.client.newRequest("POST", "plans", nil, p)
	if err != nil {
		return nil, nil, err
	}

	var dst Plan
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// Update will update the pricing or details for a plan. Existing subscriptions
// will remain at the previous renewal amounts.
// https://docs.recurly.com/api/plans#update-plan
func (s *plansImpl) Update(code string, p Plan) (*Response, *Plan, error) {
	action := fmt.Sprintf("plans/%s", code)
	req, err := s.client.newRequest("PUT", action, nil, p)
	if err != nil {
		return nil, nil, err
	}

	var dst Plan
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// Delete will make a plan inactive. New accounts cannot be created on the plan.
// https://docs.recurly.com/api/plans#delete-plan
func (s *plansImpl) Delete(code string) (*Response, error) {
	action := fmt.Sprintf("plans/%s", code)
	req, err := s.client.newRequest("DELETE", action, nil, nil)
	if err != nil {
		return nil, err
	}

	return s.client.do(req, nil)
}
