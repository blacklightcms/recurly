package recurly

import (
	"encoding/xml"
	"fmt"
)

type (
	// AddOnsService handles communication with the add ons related methods
	// of the recurly API.
	AddOnsService struct {
		client *Client
	}

	// AddOn represents an individual add on linked to a plan.
	AddOn struct {
		XMLName                     xml.Name   `xml:"add_on"`
		Code                        string     `xml:"add_on_code,omitempty"`
		Name                        string     `xml:"name,omitempty"`
		DefaultQuantity             NullInt    `xml:"default_quantity,omitempty"`
		DisplayQuantityOnHostedPage NullBool   `xml:"display_quantity_on_hosted_page,omitempty"`
		TaxCode                     string     `xml:"tax_code,omitempty"`
		UnitAmountInCents           UnitAmount `xml:"unit_amount_in_cents,omitempty"`
		AccountingCode              string     `xml:"accounting_code,omitempty"`
		CreatedAt                   NullTime   `xml:"created_at,omitempty"`
	}
)

// List returns a list of add ons for a plan.
// https://docs.recurly.com/api/plans/add-ons#list-addons
func (s *AddOnsService) List(planCode string, params Params) (*Response, []AddOn, error) {
	action := fmt.Sprintf("plans/%s/add_ons", planCode)
	req, err := s.client.newRequest("GET", action, params, nil)
	if err != nil {
		return nil, nil, err
	}

	var p struct {
		XMLName xml.Name `xml:"add_ons"`
		AddOns  []AddOn  `xml:"add_on"`
	}
	resp, err := s.client.do(req, &p)

	return resp, p.AddOns, err
}

// Get returns information about an add on.
// https://docs.recurly.com/api/plans/add-ons#lookup-addon
func (s *AddOnsService) Get(planCode string, code string) (*Response, *AddOn, error) {
	action := fmt.Sprintf("plans/%s/add_ons/%s", planCode, code)
	req, err := s.client.newRequest("GET", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst AddOn
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// Create adds an add on to a plan.
// https://docs.recurly.com/api/plans/add-ons#create-addon
func (s *AddOnsService) Create(planCode string, a AddOn) (*Response, *AddOn, error) {
	action := fmt.Sprintf("plans/%s/add_ons", planCode)
	req, err := s.client.newRequest("POST", action, nil, a)
	if err != nil {
		return nil, nil, err
	}

	var dst AddOn
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// Update will update the pricing information or description for an add-on.
// Subscriptions who have already subscribed to the add-on will not receive the new pricing.
// https://docs.recurly.com/api/plans/add-ons#update-addon
func (s *AddOnsService) Update(planCode string, code string, a AddOn) (*Response, *AddOn, error) {
	action := fmt.Sprintf("plans/%s/add_ons/%s", planCode, code)
	req, err := s.client.newRequest("PUT", action, nil, a)
	if err != nil {
		return nil, nil, err
	}

	var dst AddOn
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// Delete will remove an add on from a plan.
// https://docs.recurly.com/api/plans/add-ons#delete-addon
func (s *AddOnsService) Delete(planCode string, code string) (*Response, error) {
	action := fmt.Sprintf("plans/%s/add_ons/%s", planCode, code)
	req, err := s.client.newRequest("DELETE", action, nil, nil)
	if err != nil {
		return nil, err
	}

	return s.client.do(req, nil)
}
