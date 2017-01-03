package api

import (
	"encoding/xml"
	"fmt"

	"github.com/blacklightcms/go-recurly"
)

var _ recurly.AdjustmentsService = &AdjustmentsService{}

// AdjustmentsService handles communication with the adjustments related methods
// of the recurly API.
type AdjustmentsService struct {
	client *recurly.Client
}

// List retrieves all charges and credits issued for an account
// https://docs.recurly.com/api/adjustments#list-adjustments
func (s *AdjustmentsService) List(accountCode string, params recurly.Params) (*recurly.Response, []recurly.Adjustment, error) {
	action := fmt.Sprintf("accounts/%s/adjustments", accountCode)
	req, err := s.client.NewRequest("GET", action, params, nil)
	if err != nil {
		return nil, nil, err
	}

	var a struct {
		XMLName     xml.Name             `xml:"adjustments"`
		Adjustments []recurly.Adjustment `xml:"adjustment"`
	}
	resp, err := s.client.Do(req, &a)

	return resp, a.Adjustments, err
}

// Get returns information about a single adjustment.
// https://docs.recurly.com/api/adjustments#get-adjustments
func (s *AdjustmentsService) Get(uuid string) (*recurly.Response, *recurly.Adjustment, error) {
	action := fmt.Sprintf("adjustments/%s", uuid)
	req, err := s.client.NewRequest("GET", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst recurly.Adjustment
	resp, err := s.client.Do(req, &dst)

	return resp, &dst, err
}

// Create creates a one-time charge on an account. Charges are not invoiced or
// collected immediately. Non-invoiced charges will automatically be invoices
// when the account's subscription renews, or you trigger a collection by
// posting an invoice. Charges may be removed from an account if they have
// not been invoiced.
// https://docs.recurly.com/api/adjustments#create-adjustment
func (s *AdjustmentsService) Create(accountCode string, a recurly.Adjustment) (*recurly.Response, *recurly.Adjustment, error) {
	action := fmt.Sprintf("accounts/%s/adjustments", accountCode)
	req, err := s.client.NewRequest("POST", action, nil, a)
	if err != nil {
		return nil, nil, err
	}

	var dst recurly.Adjustment
	resp, err := s.client.Do(req, &dst)

	return resp, &dst, err
}

// Delete removes a non-invoiced adjustment from an account.
// https://docs.recurly.com/api/adjustments#delete-adjustment
func (s *AdjustmentsService) Delete(uuid string) (*recurly.Response, error) {
	action := fmt.Sprintf("adjustments/%s", uuid)
	req, err := s.client.NewRequest("DELETE", action, nil, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
