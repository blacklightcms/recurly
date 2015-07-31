package recurly

import (
	"encoding/xml"
	"fmt"
	"net"
)

type (
	// TransactionsService handles communication with the transactions related methods
	// of the recurly API.
	TransactionsService struct {
		client *Client
	}

	// Transaction ...
	Transaction struct {
		XMLName         xml.Name           `xml:"transaction"`
		Invoice         href               `xml:"invoice,omitempty"`
		Subscription    href               `xml:"subscription,omitempty"`
		UUID            string             `xml:"uuid,omitempty"`
		Action          string             `xml:"action,omitempty"`
		AmountInCents   int                `xml:"amount_in_cents"`
		TaxInCents      int                `xml:"tax_in_cents,omitempty"`
		Currency        string             `xml:"currency"`
		Status          string             `xml:"status,omitempty"`
		PaymentMethod   string             `xml:"payment_method,omitempty"`
		Reference       string             `xml:"reference,omitempty"`
		Source          string             `xml:"source,omitempty"`
		Recurring       NullBool           `xml:"recurring,omitempty"`
		Test            bool               `xml:"test,omitempty"`
		Voidable        NullBool           `xml:"voidable,omitempty"`
		Refundable      NullBool           `xml:"refundable,omitempty"`
		IPAddress       net.IP             `xml:"ip_address,omitempty"`
		CVVResult       *TransactionResult `xml:"cvv_result,omitempty"`
		AVSResult       *TransactionResult `xml:"avs_result,omitempty"`
		AVSResultStreet string             `xml:"avs_result_street,omitempty"`
		AVSResultPostal string             `xml:"avs_result_postal,omitempty"`
		CreatedAt       NullTime           `xml:"created_at,omitempty"`
		Account         Account            `xml:"details>account"`
	}

	// TransactionResult holds transaction results for CVV and AVS fields.
	TransactionResult struct {
		Code    string `xml:"code,attr"`
		Message string `xml:",innerxml"`
	}
)

const (
	// TransactionStatusSuccess is the status for a successful transaction.
	TransactionStatusSuccess = "success"

	// TransactionStatusFailed is the status for a failed transaction.
	TransactionStatusFailed = "failed"

	// TransactionStatusVoid is the status for a voided transaction.
	TransactionStatusVoid = "void"
)

// List returns a list of transactions
// https://dev.recurly.com/docs/list-transactions
func (service TransactionsService) List(params Params) (*Response, []Transaction, error) {
	req, err := service.client.newRequest("GET", "transactions", params, nil)
	if err != nil {
		return nil, nil, err
	}

	var p struct {
		XMLName      xml.Name      `xml:"transactions"`
		Transactions []Transaction `xml:"transaction"`
	}
	res, err := service.client.do(req, &p)

	return res, p.Transactions, err
}

// ListAccount returns a list of transactions for an account
// https://dev.recurly.com/docs/list-accounts-transactions
func (service TransactionsService) ListAccount(accountCode string, params Params) (*Response, []Transaction, error) {
	action := fmt.Sprintf("accounts/%s/transactions", accountCode)
	req, err := service.client.newRequest("GET", action, params, nil)
	if err != nil {
		return nil, nil, err
	}

	var p struct {
		XMLName      xml.Name      `xml:"transactions"`
		Transactions []Transaction `xml:"transaction"`
	}
	res, err := service.client.do(req, &p)

	return res, p.Transactions, err
}

// Get returns account and billing information at the time the transaction was
// submitted. It may not reflect the latest account information. A
// transaction_error section may be included if the transaction failed.
// Please see transaction error codes for more details.
// https://dev.recurly.com/docs/lookup-transaction
func (service TransactionsService) Get(uuid string) (*Response, Transaction, error) {
	action := fmt.Sprintf("transactions/%s", uuid)
	req, err := service.client.newRequest("GET", action, nil, nil)
	if err != nil {
		return nil, Transaction{}, err
	}

	var a Transaction
	res, err := service.client.do(req, &a)

	return res, a, err
}

// Create creates a new transaction. The Recurly API provides a shortcut for
// creating an invoice, charge, and optionally account, and processing the
// payment immediately. When creating an account all of the required account
// attributes must be supplied. When charging an existing account only the
// account_code must be supplied.
// https://dev.recurly.com/docs/create-transaction
func (service TransactionsService) Create(a Transaction) (*Response, Transaction, error) {
	req, err := service.client.newRequest("POST", "transactions", nil, a)
	if err != nil {
		return nil, Transaction{}, err
	}

	var dest Transaction
	res, err := service.client.do(req, &dest)

	return res, dest, err
}
