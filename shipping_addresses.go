package recurly

import (
	"context"
	"encoding/xml"
	"fmt"
)

// ShippingAddressesService manages the interactions for shipping addresses.
type ShippingAddressesService interface {
	// ListAccount returns a pager to paginate shipping addresses for an account.
	// PagerOptions are used to optionally filter the results.
	//
	// https://dev.recurly.com/docs/list-accounts-shipping-address
	ListAccount(accountCode string, opts *PagerOptions) Pager

	// Create creates a new shipping address on an existing account.
	// Note: A shipping address can also be created via Accounts.Create()
	// as well as Subscriptions.Create(). See Recurly's documentation for details.
	//
	// https://dev.recurly.com/docs/create-shipping-address-on-an-account
	Create(ctx context.Context, accountCode string, address ShippingAddress) (*ShippingAddress, error)

	// Update updates the shipping address on an account.
	//
	// https://dev.recurly.com/docs/edit-shipping-address-on-an-existing-account
	Update(ctx context.Context, accountCode string, shippingAddressID int, address ShippingAddress) (*ShippingAddress, error)

	// Delete removes an existing shipping address from ane existing account.
	//
	// https://dev.recurly.com/docs/delete-shipping-address-on-an-existing-account
	Delete(ctx context.Context, accountCode string, shippingAddressID int) error
}

// ShippingAddress represents a shipping address
type ShippingAddress struct {
	XMLName   xml.Name `xml:"shipping_address"`
	ID        int      `xml:"id,omitempty"`
	FirstName string   `xml:"first_name"`
	LastName  string   `xml:"last_name"`
	Nickname  string   `xml:"nickname,omitempty"`
	Address   string   `xml:"address1"`
	Address2  string   `xml:"address2,omitempty"`
	Company   string   `xml:"company,omitempty"`
	City      string   `xml:"city"`
	State     string   `xml:"state"`
	Zip       string   `xml:"zip"`
	Country   string   `xml:"country"`
	Phone     string   `xml:"phone,omitempty"`
	Email     string   `xml:"email,omitempty"`
	VATNumber string   `xml:"vat_number,omitempty"`
	CreatedAt NullTime `xml:"created_at,omitempty"`
	UpdatedAt NullTime `xml:"updated_at,omitempty"`
}

var _ ShippingAddressesService = &shippingAddressesImpl{}

// shippingAddressessImpl implements ShippingAddressesService.
type shippingAddressesImpl serviceImpl

func (s *shippingAddressesImpl) ListAccount(accountCode string, opts *PagerOptions) Pager {
	path := fmt.Sprintf("accounts/%s/shipping_addresses", accountCode)
	return s.client.newPager("GET", path, opts)
}

func (s *shippingAddressesImpl) Create(ctx context.Context, accountCode string, shippingAddress ShippingAddress) (*ShippingAddress, error) {
	path := fmt.Sprintf("/accounts/%s/shipping_addresses", accountCode)
	req, err := s.client.newRequest("POST", path, shippingAddress)
	if err != nil {
		return nil, err
	}

	var dst ShippingAddress
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *shippingAddressesImpl) Update(ctx context.Context, accountCode string, shippingAddressID int, shippingAddress ShippingAddress) (*ShippingAddress, error) {
	path := fmt.Sprintf("/accounts/%s/shipping_addresses/%d", accountCode, shippingAddressID)
	req, err := s.client.newRequest("PUT", path, shippingAddress)
	if err != nil {
		return nil, err
	}

	var dst ShippingAddress
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *shippingAddressesImpl) Delete(ctx context.Context, accountCode string, shippingAddressID int) error {
	path := fmt.Sprintf("/accounts/%s/shipping_addresses/%d", accountCode, shippingAddressID)
	req, err := s.client.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(ctx, req, nil)
	return err
}
