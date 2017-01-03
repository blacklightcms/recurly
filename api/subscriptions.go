package api

import (
	"encoding/xml"
	"fmt"
	"time"

	recurly "github.com/blacklightcms/go-recurly"
)

var _ recurly.SubscriptionsService = &SubscriptionsService{}

// SubscriptionsService handles communication with the subscription related methods
// of the recurly API.
type SubscriptionsService struct {
	client *recurly.Client
}

// List returns a list of all the subscriptions.
// https://docs.recurly.com/api/subscriptions#list-subscriptions
func (s *SubscriptionsService) List(params recurly.Params) (*recurly.Response, []recurly.Subscription, error) {
	req, err := s.client.NewRequest("GET", "subscriptions", params, nil)
	if err != nil {
		return nil, nil, err
	}

	var v struct {
		XMLName       xml.Name               `xml:"subscriptions"`
		Subscriptions []recurly.Subscription `xml:"subscription"`
	}
	resp, err := s.client.Do(req, &v)

	return resp, v.Subscriptions, err
}

// ListAccount returns a list of subscriptions for an account.
// https://docs.recurly.com/api/subscriptions#list-account-subscriptions
func (s *SubscriptionsService) ListAccount(accountCode string, params recurly.Params) (*recurly.Response, []recurly.Subscription, error) {
	action := fmt.Sprintf("accounts/%s/subscriptions", accountCode)
	req, err := s.client.NewRequest("GET", action, params, nil)
	if err != nil {
		return nil, nil, err
	}

	var v struct {
		XMLName       xml.Name               `xml:"subscriptions"`
		Subscriptions []recurly.Subscription `xml:"subscription"`
	}
	resp, err := s.client.Do(req, &v)

	return resp, v.Subscriptions, err
}

// Get returns a subscription by uuid
// https://docs.recurly.com/api/subscriptions#lookup-subscription
func (s *SubscriptionsService) Get(uuid string) (*recurly.Response, *recurly.Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s", uuid)
	req, err := s.client.NewRequest("GET", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst recurly.Subscription
	resp, err := s.client.Do(req, &dst)

	return resp, &dst, err
}

// Create creates a new subscription.
// https://docs.recurly.com/api/subscriptions#create-subscription
func (s *SubscriptionsService) Create(sub recurly.NewSubscription) (*recurly.Response, *recurly.Subscription, error) {
	req, err := s.client.NewRequest("POST", "subscriptions", nil, sub)
	if err != nil {
		return nil, nil, err
	}

	var dst recurly.Subscription
	resp, err := s.client.Do(req, &dst)

	return resp, &dst, err
}

// Preview returns a preview for a new subscription applied to an account.
// https://docs.recurly.com/api/subscriptions#preview-sub
func (s *SubscriptionsService) Preview(sub recurly.NewSubscription) (*recurly.Response, *recurly.Subscription, error) {
	req, err := s.client.NewRequest("POST", "subscriptions/preview", nil, sub)
	if err != nil {
		return nil, nil, err
	}

	var dst recurly.Subscription
	resp, err := s.client.Do(req, &dst)

	return resp, &dst, err
}

// Update requests an update to a subscription that takes place immediately or at renewal.
// Note: SubscriptionAddOns MUST be set to retain previous values. It's recommended you
// copy these over from a Subscription object, or use the data you have to recreate them
// identically. If updating SubscriptionAddOns, you should provide the entire replacement
// value. See recurly documentation for more info.
// https://docs.recurly.com/api/subscriptions#update-subscription
func (s *SubscriptionsService) Update(uuid string, sub recurly.UpdateSubscription) (*recurly.Response, *recurly.Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s", uuid)
	req, err := s.client.NewRequest("PUT", action, nil, sub)
	if err != nil {
		return nil, nil, err
	}

	var dst recurly.Subscription
	resp, err := s.client.Do(req, &dst)

	return resp, &dst, err
}

// UpdateNotes updates a subscription's invoice notes before the next renewal.
// Updating notes will not trigger the renewal.
// https://docs.recurly.com/api/subscriptions#update-subscription-notes
func (s *SubscriptionsService) UpdateNotes(uuid string, n recurly.SubscriptionNotes) (*recurly.Response, *recurly.Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s/notes", uuid)
	req, err := s.client.NewRequest("PUT", action, nil, n)
	if err != nil {
		return nil, nil, err
	}

	var dst recurly.Subscription
	resp, err := s.client.Do(req, &dst)

	return resp, &dst, err
}

// PreviewChange returns a preview for a subscription change applied to an
// account without committing a subscription change or posting an invoice.
// https://docs.recurly.com/api/subscriptions#sub-change-preview
func (s *SubscriptionsService) PreviewChange(uuid string, sub recurly.UpdateSubscription) (*recurly.Response, *recurly.Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s/preview", uuid)
	req, err := s.client.NewRequest("POST", action, nil, sub)
	if err != nil {
		return nil, nil, err
	}

	var dst recurly.Subscription
	resp, err := s.client.Do(req, &dst)

	return resp, &dst, err
}

// Cancel cancels a subscription so it remains active and then expires at the
// end of the current bill cycle.
// https://docs.recurly.com/api/subscriptions#cancel-subscription
func (s *SubscriptionsService) Cancel(uuid string) (*recurly.Response, *recurly.Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s/cancel", uuid)
	req, err := s.client.NewRequest("PUT", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst recurly.Subscription
	resp, err := s.client.Do(req, &dst)

	return resp, &dst, err
}

// Reactivate will reactivate a canceled subscription so it renews at the end
// of the current bill cycle.
// https://docs.recurly.com/api/subscriptions#reactivate-subscription
func (s *SubscriptionsService) Reactivate(uuid string) (*recurly.Response, *recurly.Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s/reactivate", uuid)
	req, err := s.client.NewRequest("PUT", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst recurly.Subscription
	resp, err := s.client.Do(req, &dst)

	return resp, &dst, err
}

// TerminateWithPartialRefund will terminate the active subscription
// immediately with a full refund.
// https://docs.recurly.com/api/subscriptions#terminate-subscription
func (s *SubscriptionsService) TerminateWithPartialRefund(uuid string) (*recurly.Response, *recurly.Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s/terminate", uuid)
	req, err := s.client.NewRequest("PUT", action, recurly.Params{"refund_type": "partial"}, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst recurly.Subscription
	resp, err := s.client.Do(req, &dst)

	return resp, &dst, err
}

// TerminateWithFullRefund will terminate the active subscription
// immediately with a full refund.
// https://docs.recurly.com/api/subscriptions#terminate-subscription
func (s *SubscriptionsService) TerminateWithFullRefund(uuid string) (*recurly.Response, *recurly.Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s/terminate", uuid)
	req, err := s.client.NewRequest("PUT", action, recurly.Params{"refund_type": "full"}, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst recurly.Subscription
	resp, err := s.client.Do(req, &dst)

	return resp, &dst, err
}

// TerminateWithoutRefund will terminate the active subscription
// immediately with no refund.
// https://docs.recurly.com/api/subscriptions#terminate-subscription
func (s *SubscriptionsService) TerminateWithoutRefund(uuid string) (*recurly.Response, *recurly.Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s/terminate", uuid)
	req, err := s.client.NewRequest("PUT", action, recurly.Params{"refund_type": "none"}, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst recurly.Subscription
	resp, err := s.client.Do(req, &dst)

	return resp, &dst, err
}

// Postpone will pause an an active subscription until the specified date.
// The subscription will not be prorated. For a subscription in a trial period,
// modifying the renewal date will modify when the trial expires.
// https://docs.recurly.com/api/subscriptions#postpone-subscription
func (s *SubscriptionsService) Postpone(uuid string, dt time.Time, bulk bool) (*recurly.Response, *recurly.Subscription, error) {
	action := fmt.Sprintf("subscriptions/%s/postpone", uuid)
	req, err := s.client.NewRequest("PUT", action, recurly.Params{
		"bulk":              bulk,
		"next_renewal_date": dt.Format(time.RFC3339),
	}, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst recurly.Subscription
	resp, err := s.client.Do(req, &dst)

	return resp, &dst, err
}

// Note: Create/Update Subscription with AddOns and Create/Update manual invoice
// are the same endpoint as Create. You just need to include additional parameters
// for each method. See the documentation here:
// https://dev.recurly.com/docs/subscription-add-ons
// https://dev.recurly.com/docs/update-subscription-with-add-ons
// https://dev.recurly.com/docs/subscriptions-for-manual-invoicing
// https://dev.recurly.com/docs/update-subscription-manual-invoice
