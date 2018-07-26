package recurly

import (
	"encoding/xml"
	"fmt"
)

var _ ShippingAddressesService = &shippingAddressesImpl{}

// shippingAddressessImpl handles communication with the shipping address
// related methods of the recurly API.
type shippingAddressesImpl struct {
	client *Client
}

// ListAccount returns a list of all shipping addresses associated with an account.
func (s *shippingAddressesImpl) ListAccount(accountCode string, params Params) (*Response, []ShippingAddress, error) {
	action := fmt.Sprintf("accounts/%s/shipping_addresses", accountCode)
	req, err := s.client.newRequest("GET", action, params, nil)
	if err != nil {
		return nil, nil, err
	}

	var v struct {
		XMLName           xml.Name          `xml:"shipping_addresses"`
		ShippingAddresses []ShippingAddress `xml:"shipping_address"`
	}
	resp, err := s.client.do(req, &v)
	return resp, v.ShippingAddresses, err
}

// Create creates a new shipping address.
func (s *shippingAddressesImpl) Create(accountCode string, shippingAddress ShippingAddress) (*Response, *ShippingAddress, error) {
	action := fmt.Sprintf("accounts/%s/shipping_addresses", accountCode)
	req, err := s.client.newRequest("POST", action, nil, shippingAddress)
	if err != nil {
		return nil, nil, err
	}
	var sa ShippingAddress
	resp, err := s.client.do(req, &sa)
	return resp, &sa, err
}

// Update requests an update to an existing shipping address.
func (s *shippingAddressesImpl) Update(accountCode string, shippingAddressID int64, shippingAddress ShippingAddress) (*Response, *ShippingAddress, error) {
	action := fmt.Sprintf("accounts/%s/shipping_addresses/%d", accountCode, shippingAddressID)
	req, err := s.client.newRequest("PUT", action, nil, shippingAddress)
	if err != nil {
		return nil, nil, err
	}

	var sa ShippingAddress
	resp, err := s.client.do(req, &sa)
	return resp, &sa, err
}

// Delete removes a shipping address from an account.
func (s *shippingAddressesImpl) Delete(accountCode string, shippingAddressID int64) (*Response, error) {
	action := fmt.Sprintf("accounts/%s/shipping_addresses/%d", accountCode, shippingAddressID)
	req, err := s.client.newRequest("DELETE", action, nil, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.do(req, nil)
	return resp, err
}

// GetSubscriptions fetches the subscriptions associated with a shipping address.
func (s *shippingAddressesImpl) GetSubscriptions(accountCode string, shippingAddressID int64) (*Response, []Subscription, error) {
	action := fmt.Sprintf("accounts/%s/shipping_addresses/%d/subscriptions", accountCode, shippingAddressID)
	req, err := s.client.newRequest("GET", action, nil, nil)
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
