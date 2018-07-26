package recurly

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

var _ TransactionsService = &transactionsImpl{}

// transactionsImpl handles communication with the transactions related methods
// of the recurly API.
type transactionsImpl struct {
	client *Client
}

// List returns a list of transactions
// https://dev.recurly.com/docs/list-transactions
func (s *transactionsImpl) List(params Params) (*Response, []Transaction, error) {
	req, err := s.client.newRequest("GET", "transactions", params, nil)
	if err != nil {
		return nil, nil, err
	}

	var v struct {
		XMLName      xml.Name      `xml:"transactions"`
		Transactions []Transaction `xml:"transaction"`
	}
	resp, err := s.client.do(req, &v)

	return resp, v.Transactions, err
}

// ListAccount returns a list of transactions for an account
// https://dev.recurly.com/docs/list-accounts-transactions
func (s *transactionsImpl) ListAccount(accountCode string, params Params) (*Response, []Transaction, error) {
	action := fmt.Sprintf("accounts/%s/transactions", accountCode)
	req, err := s.client.newRequest("GET", action, params, nil)
	if err != nil {
		return nil, nil, err
	}

	var v struct {
		XMLName      xml.Name      `xml:"transactions"`
		Transactions []Transaction `xml:"transaction"`
	}
	resp, err := s.client.do(req, &v)

	return resp, v.Transactions, err
}

// Get returns account and billing information at the time the transaction was
// submitted. It may not reflect the latest account information. A
// transaction_error section may be included if the transaction failed.
// Please see transaction error codes for more details.
// https://dev.recurly.com/docs/lookup-transaction
func (s *transactionsImpl) Get(uuid string) (*Response, *Transaction, error) {
	action := fmt.Sprintf("transactions/%s", SanitizeUUID(uuid))
	req, err := s.client.newRequest("GET", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst Transaction
	resp, err := s.client.do(req, &dst)
	if err != nil || resp.StatusCode >= http.StatusBadRequest {
		return resp, nil, err
	}

	return resp, &dst, err
}

// Create creates a new transaction. The Recurly API provides a shortcut for
// creating an invoice, charge, and optionally account, and processing the
// payment immediately. When creating an account all of the required account
// attributes must be supplied. When charging an existing account only the
// account_code must be supplied.
//
// See the documentation and Transaction.MarshalXML function for a detailed field list.
// https://dev.recurly.com/docs/create-transaction
func (s *transactionsImpl) Create(t Transaction) (*Response, *Transaction, error) {
	req, err := s.client.newRequest("POST", "transactions", nil, t)
	if err != nil {
		return nil, nil, err
	}

	var dst Transaction
	resp, err := s.client.do(req, &dst)

	// If there is an error set the response transaction as the returned transaction
	// so that the caller has access to TransactionError.
	if resp.IsError() {
		if resp.transaction != nil {
			dst = *resp.transaction
		}
	}

	return resp, &dst, err
}
