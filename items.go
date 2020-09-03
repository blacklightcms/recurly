package recurly

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
)

// ItemsService manages the interactions for items.
type ItemsService interface {
	// List returns pager to paginate items. PagerOptions are used to optionally
	// filter the results.
	//
	// https://dev.recurly.com/docs/list-items
	List(opts *PagerOptions) Pager

	// Get retrieves an item. If the item does not exist,
	// a nil item and nil error is returned.
	//
	// https://dev.recurly.com/docs/lookup-item
	Get(ctx context.Context, itemCode string) (*Item, error)

	// Create creates a new item.
	//
	// https://dev.recurly.com/docs/create-item
	Create(ctx context.Context, a Item) (*Item, error)

	// Update updates an item.
	//
	// https://dev.recurly.com/docs/update-item
	Update(ctx context.Context, itemCode string, a Item) (*Item, error)

	// Deactivates an item.
	//
	// https://dev.recurly.com/docs/delete-item
	Deactivate(ctx context.Context, itemCode string) error
}

// Item constants.
const (
	ItemStateActive = "active"
	ItemStateClosed = "closed"
)

// If you sell standard offerings or combinations of offerings to many customers,
// organizing those in a Recurly catalog provides many benefits.
// You'll experience faster charge creation, easier management of offerings, and analytics
// about your offerings across all sales channels. Because your offerings may be physical, digital,
// or service-oriented, Recurly collectively calls these "Items".
type Item struct {
	XMLName        xml.Name      `xml:"item"`
	Code           string        `xml:"item_code,omitempty"`
	Name           string        `xml:"name,omitempty"`
	Description    string        `xml:"description,omitempty"`
	ExternalSKU    string        `xml:"external_sku,omitempty"`
	AccountingCode string        `xml:"accounting_code,omitempty"`
	TaxExempt      NullBool      `xml:"tax_exempt,omitempty"`
	State          string        `xml:"state,omitempty"`
	CustomFields   *CustomFields `xml:"custom_fields,omitempty"`

	// The following are only valid with an `Avalara for Communications` integration
	AvalaraTransactionType int `xml:"avalara_transaction_type,omitempty"`
	AvalaraServiceType     int `xml:"avalara_service_type,omitempty"`

	CreatedAt NullTime `xml:"created_at,omitempty"`
	UpdatedAt NullTime `xml:"updated_at,omitempty"`
	DeletedAt NullTime `xml:"deleted_at,omitempty"`
}

var _ ItemsService = &itemsImpl{}

// ItemsImpl implements ItemsService.
type itemsImpl serviceImpl

func (s *itemsImpl) List(opts *PagerOptions) Pager {
	return s.client.newPager("GET", "/items", opts)
}

func (s *itemsImpl) Get(ctx context.Context, code string) (*Item, error) {
	path := fmt.Sprintf("/items/%s", code)
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var dst Item
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		if e, ok := err.(*ClientError); ok && e.Response.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &dst, nil
}

func (s *itemsImpl) Create(ctx context.Context, a Item) (*Item, error) {
	req, err := s.client.newRequest("POST", "/items", a)
	if err != nil {
		return nil, err
	}

	var dst Item
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, err
}

func (s *itemsImpl) Update(ctx context.Context, code string, a Item) (*Item, error) {
	path := fmt.Sprintf("/items/%s", code)
	req, err := s.client.newRequest("PUT", path, a)
	if err != nil {
		return nil, err
	}

	var dst Item
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, err
}

func (s *itemsImpl) Deactivate(ctx context.Context, code string) error {
	path := fmt.Sprintf("/items/%s", code)
	req, err := s.client.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(ctx, req, nil)
	return err
}
