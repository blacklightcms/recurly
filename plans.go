package recurly

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
)

// PlansService manages the interactions for plans.
type PlansService interface {
	// List returns a pager to paginate plans. PagerOptions are used to optionally
	// filter the results.
	//
	// https://dev.recurly.com/docs/list-plans
	List(opts *PagerOptions) Pager

	// Get retrieves a plan. If the plan does not exist,
	// a nil plan and nil error are returned.
	//
	// https://dev.recurly.com/docs/lookup-plan-details
	Get(ctx context.Context, code string) (*Plan, error)

	// Create a new subscription plan.
	//
	// https://dev.recurly.com/docs/create-plan
	Create(ctx context.Context, p Plan) (*Plan, error)

	// Update the pricing or details for a plan. Existing subscriptions will
	// remain at the previous renewal amounts.
	//
	// https://dev.recurly.com/docs/update-plan
	Update(ctx context.Context, code string, p Plan) (*Plan, error)

	// Delete makes a plan inactive. New subscriptions cannot be created
	// from inactive plans.
	//
	// https://dev.recurly.com/docs/delete-plan
	Delete(ctx context.Context, code string) error
}

// Plan tells Recurly how often and how much to charge your customers.
type Plan struct {
	XMLName                    xml.Name   `xml:"plan"`
	Code                       string     `xml:"plan_code,omitempty"`
	Name                       string     `xml:"name"`
	Description                string     `xml:"description,omitempty"`
	SuccessURL                 string     `xml:"success_url,omitempty"`
	CancelURL                  string     `xml:"cancel_url,omitempty"`
	DisplayDonationAmounts     NullBool   `xml:"display_donation_amounts,omitempty"`
	DisplayQuantity            NullBool   `xml:"display_quantity,omitempty"`
	DisplayPhoneNumber         NullBool   `xml:"display_phone_number,omitempty"`
	BypassHostedConfirmation   NullBool   `xml:"bypass_hosted_confirmation,omitempty"`
	UnitName                   string     `xml:"unit_name,omitempty"`
	PaymentPageTOSLink         string     `xml:"payment_page_tos_link,omitempty"`
	IntervalUnit               string     `xml:"plan_interval_unit,omitempty"`
	IntervalLength             int        `xml:"plan_interval_length,omitempty"`
	TrialIntervalUnit          string     `xml:"trial_interval_unit,omitempty"`
	TrialIntervalLength        int        `xml:"trial_interval_length,omitempty"`
	TotalBillingCycles         NullInt    `xml:"total_billing_cycles,omitempty"`
	AccountingCode             string     `xml:"accounting_code,omitempty"`
	CreatedAt                  NullTime   `xml:"created_at,omitempty"`
	TaxExempt                  NullBool   `xml:"tax_exempt,omitempty"`
	TaxCode                    string     `xml:"tax_code,omitempty"`
	AutoRenew                  bool       `xml:"auto_renew,omitempty"`
	UnitAmountInCents          UnitAmount `xml:"unit_amount_in_cents"`
	SetupFeeInCents            UnitAmount `xml:"setup_fee_in_cents,omitempty"`
	AllowAnyItemOnSubscription NullBool   `xml:"allow_any_item_on_subscription,omitempty"`

	// The following are only valid with an `Avalara for Communications` integration
	AvalaraTransactionType int `xml:"avalara_transaction_type,omitempty"`
	AvalaraServiceType     int `xml:"avalara_service_type,omitempty"`
}

var _ PlansService = &plansImpl{}

// plansImpl implements PlansService.
type plansImpl serviceImpl

func (s *plansImpl) List(opts *PagerOptions) Pager {
	return s.client.newPager("GET", "/plans", opts)
}

func (s *plansImpl) Get(ctx context.Context, code string) (*Plan, error) {
	path := fmt.Sprintf("/plans/%s", code)
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var dst Plan
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		if e, ok := err.(*ClientError); ok && e.Response.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &dst, nil
}

func (s *plansImpl) Create(ctx context.Context, p Plan) (*Plan, error) {
	req, err := s.client.newRequest("POST", "/plans", p)
	if err != nil {
		return nil, err
	}

	var dst Plan
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *plansImpl) Update(ctx context.Context, code string, p Plan) (*Plan, error) {
	path := fmt.Sprintf("/plans/%s", code)
	req, err := s.client.newRequest("PUT", path, p)
	if err != nil {
		return nil, err
	}

	var dst Plan
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *plansImpl) Delete(ctx context.Context, code string) error {
	path := fmt.Sprintf("/plans/%s", code)
	req, err := s.client.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(ctx, req, nil)
	return err
}
