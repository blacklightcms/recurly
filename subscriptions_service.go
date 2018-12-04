package recurly

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
)

var _ SubscriptionsService = &subscriptionsImpl{}

// subscriptionsImpl handles communication with the subscription related methods
// of the recurly API.
type subscriptionsImpl struct {
	client *Client
}

// List returns a list of all the subscriptions.
// https://dev.recurly.com/docs/list-subscriptions
func (s *subscriptionsImpl) List(params Params) (*Response, []Subscription, error) {
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
// https://dev.recurly.com/docs/list-accounts-subscriptions
func (s *subscriptionsImpl) ListAccount(accountCode string, params Params) (*Response, []Subscription, error) {
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
// https://dev.recurly.com/docs/lookup-subscription-details
func (s *subscriptionsImpl) Get(uuid string) (*Response, *Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s", SanitizeUUID(uuid))
	req, err := s.client.newRequest("GET", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst Subscription
	resp, err := s.client.do(req, &dst)
	if err != nil || resp.StatusCode >= http.StatusBadRequest {
		return resp, nil, err
	}

	return resp, &dst, err
}

// Create creates a new subscription.
// https://dev.recurly.com/docs/create-subscription
func (s *subscriptionsImpl) Create(sub NewSubscription) (*Response, *NewSubscriptionResponse, error) {
	req, err := s.client.newRequest("POST", "subscriptions", nil, sub)
	if err != nil {
		return nil, nil, err
	}

	var dst NewSubscriptionResponse
	var subscription Subscription
	resp, err := s.client.do(req, &subscription)
	if err != nil {
		return resp, nil, err
	}

	if subscription.UUID != "" { // If subscription not present, dst.Subscription should be nil
		dst.Subscription = &subscription
	}

	if resp.transaction != nil {
		dst.Transaction = resp.transaction
	}

	return resp, &dst, err
}

// Preview returns a preview for a new subscription applied to an account.
// https://dev.recurly.com/docs/preview-subscription
func (s *subscriptionsImpl) Preview(sub NewSubscription) (*Response, *Subscription, error) {
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
// https://dev.recurly.com/docs/update-subscription
func (s *subscriptionsImpl) Update(uuid string, sub UpdateSubscription) (*Response, *Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s", SanitizeUUID(uuid))
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
// https://dev.recurly.com/docs/update-subscription-notes
func (s *subscriptionsImpl) UpdateNotes(uuid string, n SubscriptionNotes) (*Response, *Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s/notes", SanitizeUUID(uuid))
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
// https://dev.recurly.com/docs/sub-change-preview
func (s *subscriptionsImpl) PreviewChange(uuid string, sub UpdateSubscription) (*Response, *Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s/preview", SanitizeUUID(uuid))
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
// https://dev.recurly.com/docs/cancel-subscription
func (s *subscriptionsImpl) Cancel(uuid string) (*Response, *Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s/cancel", SanitizeUUID(uuid))
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
// https://dev.recurly.com/docs/reactivate-subscription
func (s *subscriptionsImpl) Reactivate(uuid string) (*Response, *Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s/reactivate", SanitizeUUID(uuid))
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
// https://dev.recurly.com/docs/terminate-subscription
func (s *subscriptionsImpl) TerminateWithPartialRefund(uuid string) (*Response, *Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s/terminate", SanitizeUUID(uuid))
	req, err := s.client.newRequest("PUT", action, Params{"refund": "partial"}, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst Subscription
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// TerminateWithFullRefund will terminate the active subscription
// immediately with a full refund.
// https://dev.recurly.com/docs/terminate-subscription
func (s *subscriptionsImpl) TerminateWithFullRefund(uuid string) (*Response, *Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s/terminate", SanitizeUUID(uuid))
	req, err := s.client.newRequest("PUT", action, Params{"refund": "full"}, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst Subscription
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// TerminateWithoutRefund will terminate the active subscription
// immediately with no refund.
// https://dev.recurly.com/docs/terminate-subscription
func (s *subscriptionsImpl) TerminateWithoutRefund(uuid string) (*Response, *Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s/terminate", SanitizeUUID(uuid))
	req, err := s.client.newRequest("PUT", action, Params{"refund": "none"}, nil)
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
// https://dev.recurly.com/docs/postpone-subscription
func (s *subscriptionsImpl) Postpone(uuid string, dt time.Time, bulk bool) (*Response, *Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s/postpone", SanitizeUUID(uuid))
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

// Pause will pause an active subscription for the specified number of billing cycles.
// The pause takes effect at the beginning of the next billing cycle.
// https://dev.recurly.com/docs/pause-subscription
func (s *subscriptionsImpl) Pause(uuid string, cycles int) (*Response, *Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s/pause", SanitizeUUID(uuid))
	type subscription struct {
		RemainingPauseCycles int `xml:"remaining_pause_cycles"`
	}
	pauseCycles := subscription{cycles}
	req, err := s.client.newRequest("PUT", action, nil, pauseCycles)

	var dst Subscription
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// Resume will immediately resume a paused subscription.
// https://dev.recurly.com/docs/resume-subscription
func (s *subscriptionsImpl) Resume(uuid string) (*Response, *Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s/resume", SanitizeUUID(uuid))
	req, err := s.client.newRequest("PUT", action, nil, nil)

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
