package recurly

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

var _ CreditPaymentsService = &creditInvoicesImpl{}

// creditInvoicesImpl handles communication with the credit payment
// related methods of the recurly API.
type creditInvoicesImpl struct {
	client *Client
}

// List returns a list of all credit payments.
// https://dev.recurly.com/docs/list-credit-payments
func (s *creditInvoicesImpl) List(params Params) (*Response, []CreditPayment, error) {
	req, err := s.client.newRequest("GET", "credit_payments", params, nil)
	if err != nil {
		return nil, nil, err
	}

	var p struct {
		XMLName        xml.Name        `xml:"credit_payments"`
		CreditPayments []CreditPayment `xml:"credit_payment"`
	}
	resp, err := s.client.do(req, &p)

	return resp, p.CreditPayments, err
}

// ListAccount returns a list of all credit payments for an account.
// https://dev.recurly.com/docs/list-credit-payments-on-account
func (s *creditInvoicesImpl) ListAccount(accountCode string, params Params) (*Response, []CreditPayment, error) {
	action := fmt.Sprintf("accounts/%s/credit_payments", accountCode)
	req, err := s.client.newRequest("GET", action, params, nil)
	if err != nil {
		return nil, nil, err
	}

	var p struct {
		XMLName        xml.Name        `xml:"credit_payments"`
		CreditPayments []CreditPayment `xml:"credit_payment"`
	}
	resp, err := s.client.do(req, &p)

	return resp, p.CreditPayments, err
}

// Get returns detailed information about a credit payment.
// https://dev.recurly.com/docs/lookup-credit-payment
func (s *creditInvoicesImpl) Get(uuid string) (*Response, *CreditPayment, error) {
	action := fmt.Sprintf("credit_payments/%s", uuid)
	req, err := s.client.newRequest("GET", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst CreditPayment
	resp, err := s.client.do(req, &dst)
	if err != nil || resp.StatusCode >= http.StatusBadRequest {
		return resp, nil, err
	}

	return resp, &dst, err
}
