package recurly

import (
	"net/http"
)

var _ PurchasesService = &purchasesImpl{}

type purchasesImpl struct {
	client *Client
}

func (s *purchasesImpl) Create(p Purchase) (*Response, *InvoiceCollection, error) {
	req, err := s.client.newRequest("POST", "purchases", nil, p)
	if err != nil {
		return nil, nil, err
	}

	var dst InvoiceCollection
	resp, err := s.client.do(req, &dst)
	if err != nil || resp.StatusCode >= http.StatusBadRequest {
		return resp, nil, err
	}

	return resp, &dst, err
}

func (s *purchasesImpl) Preview(p Purchase) (*Response, *InvoiceCollection, error) {
	req, err := s.client.newRequest("POST", "purchases/preview", nil, p)
	if err != nil {
		return nil, nil, err
	}

	var dst InvoiceCollection
	resp, err := s.client.do(req, &dst)
	if err != nil || resp.StatusCode >= http.StatusBadRequest {
		return resp, nil, err
	}

	return resp, &dst, err
}
