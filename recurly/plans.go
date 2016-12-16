package recurly

import (
	"encoding/xml"
	"fmt"
)

type (
	// PlansService handles communication with the plans related methods
	// of the recurly API.
	PlansService struct {
		client *Client
	}

	// Plan represents an individual plan on your site.
	Plan struct {
		XMLName                  xml.Name   `xml:"plan"`
		Code                     string     `xml:"plan_code,omitempty"`
		Name                     string     `xml:"name"`
		Description              string     `xml:"description,omitempty"`
		SuccessURL               string     `xml:"success_url,omitempty"`
		CancelURL                string     `xml:"cancel_url,omitempty"`
		DisplayDonationAmounts   NullBool   `xml:"display_donation_amounts,omitempty"`
		DisplayQuantity          NullBool   `xml:"display_quantity,omitempty"`
		DisplayPhoneNumber       NullBool   `xml:"display_phone_number,omitempty"`
		BypassHostedConfirmation NullBool   `xml:"bypass_hosted_confirmation,omitempty"`
		UnitName                 string     `xml:"unit_name,omitempty"`
		PaymentPageTOSLink       string     `xml:"payment_page_tos_link,omitempty"`
		IntervalUnit             string     `xml:"plan_interval_unit,omitempty"`
		IntervalLength           int        `xml:"plan_interval_length,omitempty"`
		TrialIntervalUnit        string     `xml:"trial_interval_unit,omitempty"`
		TrialIntervalLength      int        `xml:"trial_interval_length,omitempty"`
		TotalBillingCycles       NullInt    `xml:"total_billing_cycles,omitempty"`
		AccountingCode           string     `xml:"accounting_code,omitempty"`
		CreatedAt                NullTime   `xml:"created_at,omitempty"`
		TaxExempt                NullBool   `xml:"tax_exempt,omitempty"`
		TaxCode                  string     `xml:"tax_code,omitempty"`
		UnitAmountInCents        UnitAmount `xml:"unit_amount_in_cents"`
		SetupFeeInCents          UnitAmount `xml:"setup_fee_in_cents,omitempty"`
	}
)

// List will retrieve all your active subscription plans.
// https://docs.recurly.com/api/plans#list-plans
func (service PlansService) List(params Params) (*Response, []Plan, error) {
	req, err := service.client.newRequest("GET", "plans", params, nil)
	if err != nil {
		return nil, nil, err
	}

	var p struct {
		XMLName xml.Name `xml:"plans"`
		Plans   []Plan   `xml:"plan"`
	}
	resp, err := service.client.do(req, &p)

	return resp, p.Plans, err
}

// Get will lookup a specific plan by code.
// https://docs.recurly.com/api/plans#lookup-plan
func (service PlansService) Get(code string) (*Response, *Plan, error) {
	action := fmt.Sprintf("plans/%s", code)
	req, err := service.client.newRequest("GET", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst Plan
	resp, err := service.client.do(req, &dst)

	return resp, &dst, err
}

// Create will create a new subscription plan.
// https://docs.recurly.com/api/plans#create-plan
func (service PlansService) Create(p Plan) (*Response, *Plan, error) {
	req, err := service.client.newRequest("POST", "plans", nil, p)
	if err != nil {
		return nil, nil, err
	}

	var dst Plan
	resp, err := service.client.do(req, &dst)

	return resp, &dst, err
}

// Update will update the pricing or details for a plan. Existing subscriptions
// will remain at the previous renewal amounts.
// https://docs.recurly.com/api/plans#update-plan
func (service PlansService) Update(code string, p Plan) (*Response, *Plan, error) {
	action := fmt.Sprintf("plans/%s", code)
	req, err := service.client.newRequest("PUT", action, nil, p)
	if err != nil {
		return nil, nil, err
	}

	var dst Plan
	resp, err := service.client.do(req, &dst)

	return resp, &dst, err
}

// Delete will make a plan inactive. New accounts cannot be created on the plan.
// https://docs.recurly.com/api/plans#delete-plan
func (service PlansService) Delete(code string) (*Response, error) {
	action := fmt.Sprintf("plans/%s", code)
	req, err := service.client.newRequest("DELETE", action, nil, nil)
	if err != nil {
		return nil, err
	}

	return service.client.do(req, nil)
}
