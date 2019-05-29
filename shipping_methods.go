package recurly

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
)

// ShippingMethodsService manages the interactions for shipping methods.
type ShippingMethodsService interface {
	// ListAccount returns a pager to paginate available shipping methods.
	// PagerOptions are used to optionally filter the results.
	//
	// https://dev.recurly.com/docs/list-shipping-methods
	List(opts *PagerOptions) Pager

	// Get retrieves a shipping method. If the shipping method does not exist,
	// a nil shipping method and nil error are returned.
	//
	// https://dev.recurly.com/docs/lookup-shipping-method
	Get(ctx context.Context, code string) (*ShippingMethod, error)
}

// ShippingMethod holds a shipping method.
type ShippingMethod struct {
	XMLName        xml.Name `xml:"shipping_method"`
	Code           string   `xml:"code"`
	Name           string   `xml:"name"`
	AccountingCode string   `xml:"accounting_code"`
	TaxCode        string   `xml:"tax_code"`
	CreatedAt      NullTime `xml:"created_at"`
	UpdatedAt      NullTime `xml:"updated_at"`
}

var _ ShippingMethodsService = &shippingMethodsImpl{}

// shippingMethodssImpl implements ShippingMethodsService.
type shippingMethodsImpl serviceImpl

func (s *shippingMethodsImpl) List(opts *PagerOptions) Pager {
	return s.client.newPager("GET", "/shipping_methods", opts)
}

func (s *shippingMethodsImpl) Get(ctx context.Context, code string) (*ShippingMethod, error) {
	path := fmt.Sprintf("/shipping_methods/%s", code)
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var dst ShippingMethod
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		if e, ok := err.(*ClientError); ok && e.Response.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &dst, nil
}
