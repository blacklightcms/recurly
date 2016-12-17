package recurly

import (
	"encoding/xml"
	"fmt"
	"time"
)

type (
	// SubscriptionsInterface defines methods for the subscription related methods
	// of the recurly API.
	SubscriptionsInterface interface {
		List(params Params) (*Response, []Subscription, error)
		ListAccount(accountCode string, params Params) (*Response, []Subscription, error)
		Get(uuid string) (*Response, Subscription, error)
		Create(s NewSubscription) (*Response, Subscription, error)
		Preview(s NewSubscription) (*Response, Subscription, error)
		Update(uuid string, s UpdateSubscription) (*Response, Subscription, error)
		UpdateNotes(uuid string, n SubscriptionNotes) (*Response, Subscription, error)
		PreviewChange(uuid string, s UpdateSubscription) (*Response, Subscription, error)
		Cancel(uuid string) (*Response, Subscription, error)
		Reactivate(uuid string) (*Response, Subscription, error)
		TerminateWithPartialRefund(uuid string) (*Response, Subscription, error)
		TerminateWithFullRefund(uuid string) (*Response, Subscription, error)
		TerminateWithoutRefund(uuid string) (*Response, Subscription, error)
		Postpone(uuid string, dt time.Time, bulk bool) (*Response, Subscription, error)
	}

	// SubscriptionsService handles communication with the subscription related methods
	// of the recurly API.
	SubscriptionsService struct {
		client *Client
	}

	// Subscription represents an individual subscription.
	Subscription struct {
		XMLName                xml.Name            `xml:"subscription"`
		Plan                   nestedPlan          `xml:"plan,omitempty"`
		AccountCode            hrefString          `xml:"account"` // Read only
		InvoiceNumber          hrefString          `xml:"invoice"` // Read only
		UUID                   string              `xml:"uuid,omitempty"`
		State                  string              `xml:"state,omitempty"`
		UnitAmountInCents      int                 `xml:"unit_amount_in_cents,omitempty"`
		Currency               string              `xml:"currency,omitempty"`
		Quantity               int                 `xml:"quantity,omitempty"`
		ActivatedAt            NullTime            `xml:"activated_at,omitempty"`
		CanceledAt             NullTime            `xml:"canceled_at,omitempty"`
		ExpiresAt              NullTime            `xml:"expires_at,omitempty"`
		CurrentPeriodStartedAt NullTime            `xml:"current_period_started_at,omitempty"`
		CurrentPeriodEndsAt    NullTime            `xml:"current_period_ends_at,omitempty"`
		TrialStartedAt         NullTime            `xml:"trial_started_at,omitempty"`
		TrialEndsAt            NullTime            `xml:"trial_ends_at,omitempty"`
		TaxInCents             int                 `xml:"tax_in_cents,omitempty"`
		TaxType                string              `xml:"tax_type,omitempty"`
		TaxRegion              string              `xml:"tax_region,omitempty"`
		TaxRate                float64             `xml:"tax_rate,omitempty"`
		PONumber               string              `xml:"po_number,omitempty"`
		NetTerms               NullInt             `xml:"net_terms,omitempty"`
		SubscriptionAddOns     []SubscriptionAddOn `xml:"subscription_add_ons>subscription_add_on,omitempty"`
	}

	nestedPlan struct {
		Code string `xml:"plan_code,omitempty"`
		Name string `xml:"name,omitempty"`
	}

	// SubscriptionAddOn are add ons to subscriptions.
	// https://docs.recurly.com/api/subscriptions/subscription-add-ons
	SubscriptionAddOn struct {
		XMLName           xml.Name `xml:"subscription_add_on"`
		Code              string   `xml:"add_on_code"`
		UnitAmountInCents int      `xml:"unit_amount_in_cents"`
		Quantity          int      `xml:"quantity,omitempty"`
	}

	// NewSubscription is used to create new subscriptions.
	NewSubscription struct {
		XMLName                 xml.Name             `xml:"subscription"`
		PlanCode                string               `xml:"plan_code"`
		Account                 Account              `xml:"account"`
		SubscriptionAddOns      *[]SubscriptionAddOn `xml:"subscription_add_ons>subscription_add_on,omitempty"`
		CouponCode              string               `xml:"coupon_code,omitempty"`
		UnitAmountInCents       int                  `xml:"unit_amount_in_cents,omitempty"`
		Currency                string               `xml:"currency"`
		Quantity                int                  `xml:"quantity,omitempty"`
		TrialEndsAt             NullTime             `xml:"trial_ends_at,omitempty"`
		StartsAt                NullTime             `xml:"starts_at,omitempty"`
		TotalBillingCycles      int                  `xml:"total_billing_cycles,omitempty"`
		FirstRenewalDate        NullTime             `xml:"first_renewal_date,omitempty"`
		CollectionMethod        string               `xml:"collection_method,omitempty"`
		NetTerms                NullInt              `xml:"net_terms,omitempty"`
		PONumber                string               `xml:"po_number,omitempty"`
		Bulk                    bool                 `xml:"bulk,omitempty"`
		TermsAndConditions      string               `xml:"terms_and_conditions,omitempty"`
		CustomerNotes           string               `xml:"customer_notes,omitempty"`
		VATReverseChargeNotes   string               `xml:"vat_reverse_charge_notes,omitempty"`
		BankAccountAuthorizedAt NullTime             `xml:"bank_account_authorized_at,omitempty"`
	}

	// UpdateSubscription is used to update subscriptions
	UpdateSubscription struct {
		XMLName            xml.Name             `xml:"subscription"`
		Timeframe          string               `xml:"timeframe,omitempty"`
		PlanCode           string               `xml:"plan_code,omitempty"`
		Quantity           int                  `xml:"quantity,omitempty"`
		UnitAmountInCents  int                  `xml:"unit_amount_in_cents,omitempty"`
		CollectionMethod   string               `xml:"collection_method,omitempty"`
		NetTerms           NullInt              `xml:"net_terms,omitempty"`
		PONumber           string               `xml:"po_number,omitempty"`
		SubscriptionAddOns *[]SubscriptionAddOn `xml:"subscription_add_ons>subscription_add_on,omitempty"`
	}

	// SubscriptionNotes is used to update a subscription's notes.
	SubscriptionNotes struct {
		XMLName               xml.Name `xml:"subscription"`
		TermsAndConditions    string   `xml:"terms_and_conditions,omitempty"`
		CustomerNotes         string   `xml:"customer_notes,omitempty"`
		VATReverseChargeNotes string   `xml:"vat_reverse_charge_notes,omitempty"`
	}
)

// Ensure SubscriptionsService implements interface.
var _ SubscriptionsInterface = &SubscriptionsService{}

const (
	// SubscriptionStateActive represents subscriptions that are valid for the
	// current time. This includes subscriptions in a trial period
	SubscriptionStateActive = "active"

	// SubscriptionStateCanceled are subscriptions that are valid for
	// the current time but will not renew because a cancelation was requested
	SubscriptionStateCanceled = "canceled"

	// SubscriptionStateExpired are subscriptions that have expired and are no longer valid
	SubscriptionStateExpired = "expired"

	// SubscriptionStateFuture are subscriptions that will start in the
	// future, they are not active yet
	SubscriptionStateFuture = "future"

	// SubscriptionStateInTrial are subscriptions that are active or canceled
	// and are in a trial period
	SubscriptionStateInTrial = "in_trial"

	// SubscriptionStateLive are all subscriptions that are not expired
	SubscriptionStateLive = "live"

	// SubscriptionStatePastDue are subscriptions that are active or canceled
	// and have a past-due invoice
	SubscriptionStatePastDue = "past_due"
)

// MakeUpdate creates an UpdateSubscription with values that need to be passed
// on update to be retained (meaning nil/zero values will delete that value).
// After calling MakeUpdate you should modify the struct with your updates.
// Once you're ready you can call client.Subscriptions.Update
func (s Subscription) MakeUpdate() UpdateSubscription {
	return UpdateSubscription{
		// NetTerms need to be copied over because on update they default to 0.
		// This ensures the NetTerms don't get overridden.
		NetTerms:           s.NetTerms,
		SubscriptionAddOns: &s.SubscriptionAddOns,
	}
}

// List returns a list of all the subscriptions.
// https://docs.recurly.com/api/subscriptions#list-subscriptions
func (s *SubscriptionsService) List(params Params) (*Response, []Subscription, error) {
	req, err := s.client.newRequest("GET", "subscriptions", params, nil)
	if err != nil {
		return nil, nil, err
	}

	var v struct {
		XMLName       xml.Name       `xml:"subscriptions"`
		Subscriptions []Subscription `xml:"subscription"`
	}
	resp, err := s.client.do(req, &v)

	return resp, v.Subscriptions, err
}

// ListAccount returns a list of subscriptions for an account.
// https://docs.recurly.com/api/subscriptions#list-account-subscriptions
func (s *SubscriptionsService) ListAccount(accountCode string, params Params) (*Response, []Subscription, error) {
	action := fmt.Sprintf("accounts/%s/subscriptions", accountCode)
	req, err := s.client.newRequest("GET", action, params, nil)
	if err != nil {
		return nil, nil, err
	}

	var v struct {
		XMLName       xml.Name       `xml:"subscriptions"`
		Subscriptions []Subscription `xml:"subscription"`
	}
	resp, err := s.client.do(req, &v)

	return resp, v.Subscriptions, err
}

// Get returns a subscription by uuid
// https://docs.recurly.com/api/subscriptions#lookup-subscription
func (s *SubscriptionsService) Get(uuid string) (*Response, *Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s", uuid)
	req, err := s.client.newRequest("GET", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst Subscription
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// Create creates a new subscription.
// https://docs.recurly.com/api/subscriptions#create-subscription
func (s *SubscriptionsService) Create(sub NewSubscription) (*Response, *Subscription, error) {
	req, err := s.client.newRequest("POST", "subscriptions", nil, sub)
	if err != nil {
		return nil, nil, err
	}

	var dst Subscription
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// Preview returns a preview for a new subscription applied to an account.
// https://docs.recurly.com/api/subscriptions#preview-sub
func (s *SubscriptionsService) Preview(sub NewSubscription) (*Response, *Subscription, error) {
	req, err := s.client.newRequest("POST", "subscriptions/preview", nil, sub)
	if err != nil {
		return nil, nil, err
	}

	var dst Subscription
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// Update requests an update to a subscription that takes place immediately or at renewal.
// Note: SubscriptionAddOns MUST be set to retain previous values. It's recommended you
// copy these over from a Subscription object, or use the data you have to recreate them
// identically. If updating SubscriptionAddOns, you should provide the entire replacement
// value. See recurly documentation for more info.
// https://docs.recurly.com/api/subscriptions#update-subscription
func (s *SubscriptionsService) Update(uuid string, sub UpdateSubscription) (*Response, *Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s", uuid)
	req, err := s.client.newRequest("PUT", action, nil, sub)
	if err != nil {
		return nil, nil, err
	}

	var dst Subscription
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// UpdateNotes updates a subscription's invoice notes before the next renewal.
// Updating notes will not trigger the renewal.
// https://docs.recurly.com/api/subscriptions#update-subscription-notes
func (s *SubscriptionsService) UpdateNotes(uuid string, n SubscriptionNotes) (*Response, *Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s/notes", uuid)
	req, err := s.client.newRequest("PUT", action, nil, n)
	if err != nil {
		return nil, nil, err
	}

	var dst Subscription
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// PreviewChange returns a preview for a subscription change applied to an
// account without committing a subscription change or posting an invoice.
// https://docs.recurly.com/api/subscriptions#sub-change-preview
func (s *SubscriptionsService) PreviewChange(uuid string, sub UpdateSubscription) (*Response, *Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s/preview", uuid)
	req, err := s.client.newRequest("POST", action, nil, sub)
	if err != nil {
		return nil, nil, err
	}

	var dst Subscription
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// Cancel cancels a subscription so it remains active and then expires at the
// end of the current bill cycle.
// https://docs.recurly.com/api/subscriptions#cancel-subscription
func (s *SubscriptionsService) Cancel(uuid string) (*Response, *Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s/cancel", uuid)
	req, err := s.client.newRequest("PUT", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst Subscription
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// Reactivate will reactivate a canceled subscription so it renews at the end
// of the current bill cycle.
// https://docs.recurly.com/api/subscriptions#reactivate-subscription
func (s *SubscriptionsService) Reactivate(uuid string) (*Response, *Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s/reactivate", uuid)
	req, err := s.client.newRequest("PUT", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst Subscription
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// TerminateWithPartialRefund will terminate the active subscription
// immediately with a full refund.
// https://docs.recurly.com/api/subscriptions#terminate-subscription
func (s *SubscriptionsService) TerminateWithPartialRefund(uuid string) (*Response, *Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s/terminate", uuid)
	req, err := s.client.newRequest("PUT", action, Params{"refund_type": "partial"}, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst Subscription
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// TerminateWithFullRefund will terminate the active subscription
// immediately with a full refund.
// https://docs.recurly.com/api/subscriptions#terminate-subscription
func (s *SubscriptionsService) TerminateWithFullRefund(uuid string) (*Response, *Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s/terminate", uuid)
	req, err := s.client.newRequest("PUT", action, Params{"refund_type": "full"}, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst Subscription
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// TerminateWithoutRefund will terminate the active subscription
// immediately with no refund.
// https://docs.recurly.com/api/subscriptions#terminate-subscription
func (s *SubscriptionsService) TerminateWithoutRefund(uuid string) (*Response, *Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s/terminate", uuid)
	req, err := s.client.newRequest("PUT", action, Params{"refund_type": "none"}, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst Subscription
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// Postpone will pause an an active subscription until the specified date.
// The subscription will not be prorated. For a subscription in a trial period,
// modifying the renewal date will modify when the trial expires.
// https://docs.recurly.com/api/subscriptions#postpone-subscription
func (s *SubscriptionsService) Postpone(uuid string, dt time.Time, bulk bool) (*Response, *Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s/postpone", uuid)
	req, err := s.client.newRequest("PUT", action, Params{
		"bulk":              bulk,
		"next_renewal_date": dt.Format(time.RFC3339),
	}, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst Subscription
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// Note: Create/Update Subscription with AddOns and Create/Update manual invoice
// are the same endpoint as Create. You just need to include additional parameters
// for each method. See the documentation here:
// https://dev.recurly.com/docs/subscription-add-ons
// https://dev.recurly.com/docs/update-subscription-with-add-ons
// https://dev.recurly.com/docs/subscriptions-for-manual-invoicing
// https://dev.recurly.com/docs/update-subscription-manual-invoice
