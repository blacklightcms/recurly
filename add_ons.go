package recurly

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
)

// AddOnsService manages the interactions for add-ons.
type AddOnsService interface {
	// List returns a pager to paginate add-ons for a plan. PagerOptions are used to
	// optionally filter the results.
	//
	// https://dev.recurly.com/docs/list-add-ons-for-a-plan
	List(planCode string, opts *PagerOptions) Pager

	// Get retrieves an add-on. If the add-on does not exist,
	// a nil add-on and nil error are returned.
	//
	// https://dev.recurly.com/docs/lookup-an-add-on
	Get(ctx context.Context, planCode string, code string) (*AddOn, error)

	// Create creates a new add-on to a plan.
	//
	// https://dev.recurly.com/docs/create-an-add-on
	Create(ctx context.Context, planCode string, a AddOn) (*AddOn, error)

	// Update updates the pricing information or description for an add-on.
	// Existing subscriptions with the add-on will not receive any pricing updates.
	//
	// https://dev.recurly.com/docs/update-an-add-on
	Update(ctx context.Context, planCode string, code string, a AddOn) (*AddOn, error)

	// Delete removes an add-on from a plan.
	//
	// https://dev.recurly.com/docs/delete-an-add-on
	Delete(ctx context.Context, planCode string, code string) error
}

// An AddOn is a charge billed each billing period in addition to a subscription’s
// base charge. Each plan may have one or more add-ons associated with it.
//
// https://dev.recurly.com/docs/plan-add-ons-object
type AddOn struct {
	XMLName                     xml.Name   `xml:"add_on"`
	Code                        string     `xml:"add_on_code,omitempty"`
	Name                        string     `xml:"name,omitempty"`
	DefaultQuantity             NullInt    `xml:"default_quantity,omitempty"`
	DisplayQuantityOnHostedPage NullBool   `xml:"display_quantity_on_hosted_page,omitempty"`
	TaxCode                     string     `xml:"tax_code,omitempty"`
	UnitAmountInCents           UnitAmount `xml:"unit_amount_in_cents,omitempty"`
	AccountingCode              string     `xml:"accounting_code,omitempty"`
	ExternalSKU                 string     `xml:"external_sku,omitempty"`
	ItemState                   string     `xml:"item_state,omitempty"`
	ItemCode                    string     `xml:"item_code,omitempty"`
	TierType                    string     `xml:"tier_type,omitempty"`
	Tiers                       *[]Tier    `xml:"tiers>tier,omitempty"`

	// The following are only valid with an `Avalara for Communications` integration
	AvalaraTransactionType int `xml:"avalara_transaction_type,omitempty"`
	AvalaraServiceType     int `xml:"avalara_service_type,omitempty"`

	CreatedAt NullTime `xml:"created_at,omitempty"`
}

// UnitAmount can read or write amounts in various currencies.
type UnitAmount struct {
	USD int `xml:"USD,omitempty"`
	EUR int `xml:"EUR,omitempty"`
	GBP int `xml:"GBP,omitempty"`
	CAD int `xml:"CAD,omitempty"`
	AUD int `xml:"AUD,omitempty"`
}

// MarshalXML ensures UnitAmount is not marshaled unless one or more currencies
// has a value greater than zero.
func (u UnitAmount) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if u.USD > 0 || u.EUR > 0 || u.CAD > 0 || u.GBP > 0 || u.AUD > 0 {
		type uaAlias UnitAmount
		e.EncodeElement(uaAlias(u), start)
	}
	return nil
}

// Tier is used for Quantity Based Pricing models https://docs.recurly.com/docs/billing-models#section-quantity-based
type Tier struct {
	XMLName           xml.Name   `xml:"tier"`
	UnitAmountInCents UnitAmount `xml:"unit_amount_in_cents,omitempty"`
	EndingQuantity    int        `xml:"ending_quantity,omitempty"`
}

var _ AddOnsService = &addOnsImpl{}

// addOnsImpl implements AddOnsService.
type addOnsImpl serviceImpl

func (s *addOnsImpl) List(planCode string, opts *PagerOptions) Pager {
	path := fmt.Sprintf("/plans/%s/add_ons", planCode)
	return s.client.newPager("GET", path, opts)
}

func (s *addOnsImpl) Get(ctx context.Context, planCode string, code string) (*AddOn, error) {
	path := fmt.Sprintf("/plans/%s/add_ons/%s", planCode, code)
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var dst AddOn
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		if e, ok := err.(*ClientError); ok && e.Response.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &dst, nil
}

func (s *addOnsImpl) Create(ctx context.Context, planCode string, a AddOn) (*AddOn, error) {
	path := fmt.Sprintf("/plans/%s/add_ons", planCode)
	req, err := s.client.newRequest("POST", path, a)
	if err != nil {
		return nil, err
	}

	var dst AddOn
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *addOnsImpl) Update(ctx context.Context, planCode string, code string, a AddOn) (*AddOn, error) {
	path := fmt.Sprintf("/plans/%s/add_ons/%s", planCode, code)
	req, err := s.client.newRequest("PUT", path, a)
	if err != nil {
		return nil, err
	}

	var dst AddOn
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *addOnsImpl) Delete(ctx context.Context, planCode string, code string) error {
	path := fmt.Sprintf("/plans/%s/add_ons/%s", planCode, code)
	req, err := s.client.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(ctx, req, nil)
	return err
}
