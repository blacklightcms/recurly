package recurly

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

var _ AdjustmentsService = &adjustmentsImpl{}

// adjustmentsImpl handles communication with the adjustments related methods
// of the recurly API.
type adjustmentsImpl struct {
	client *Client
}

// List retrieves all charges and credits issued for an account
// https://docs.recurly.com/api/adjustments#list-adjustments
func (s *adjustmentsImpl) List(accountCode string, params Params) (*Response, []Adjustment, error) {
	action := fmt.Sprintf("accounts/%s/adjustments", accountCode)
	req, err := s.client.newRequest("GET", action, params, nil)
	if err != nil {
		return nil, nil, err
	}

	var a struct {
		XMLName     xml.Name     `xml:"adjustments"`
		Adjustments []Adjustment `xml:"adjustment"`
	}
	resp, err := s.client.do(req, &a)

	return resp, a.Adjustments, err
}

// Get returns information about a single adjustment.
// https://docs.recurly.com/api/adjustments#get-adjustments
func (s *adjustmentsImpl) Get(uuid string) (*Response, *Adjustment, error) {
	action := fmt.Sprintf("adjustments/%s", SanitizeUUID(uuid))
	req, err := s.client.newRequest("GET", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst Adjustment
	resp, err := s.client.do(req, &dst)
	if err != nil || resp.StatusCode >= http.StatusBadRequest {
		return resp, nil, err
	}

	return resp, &dst, err
}

// Create creates a one-time charge on an account. Charges are not invoiced or
// collected immediately. Non-invoiced charges will automatically be invoices
// when the account's subscription renews, or you trigger a collection by
// posting an invoice. Charges may be removed from an account if they have
// not been invoiced.
// https://docs.recurly.com/api/adjustments#create-adjustment
func (s *adjustmentsImpl) Create(accountCode string, a Adjustment) (*Response, *Adjustment, error) {
	action := fmt.Sprintf("accounts/%s/adjustments", accountCode)
	req, err := s.client.newRequest("POST", action, nil, a)
	if err != nil {
		return nil, nil, err
	}

	var dst Adjustment
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// Delete removes a non-invoiced adjustment from an account.
// https://docs.recurly.com/api/adjustments#delete-adjustment
func (s *adjustmentsImpl) Delete(uuid string) (*Response, error) {
	action := fmt.Sprintf("adjustments/%s", SanitizeUUID(uuid))
	req, err := s.client.newRequest("DELETE", action, nil, nil)
	if err != nil {
		return nil, err
	}

	return s.client.do(req, nil)
}
