package recurly

import (
	"encoding/xml"
	"fmt"
)

const (
	// AccountStateActive is the status for active accounts.
	AccountStateActive = "active"

	// AccountStateClosed is the status for closed accounts.
	AccountStateClosed = "closed"
)

var _ AccountsService = &accountsImpl{}

// accountsImpl handles communication with the accounts related methods
// of the recurly API.
type accountsImpl struct {
	client *Client
}

// NewAccountsImpl returns a new instance of accountsImpl.
func NewAccountsImpl(client *Client) *accountsImpl {
	return &accountsImpl{client: client}
}

// List returns a list of the accounts on your site.
// https://docs.recurly.com/api/accounts#list-accounts
func (s *accountsImpl) List(params Params) (*Response, []Account, error) {
	req, err := s.client.NewRequest("GET", "accounts", params, nil)
	if err != nil {
		return nil, nil, err
	}

	var a struct {
		XMLName  xml.Name  `xml:"accounts"`
		Accounts []Account `xml:"account"`
	}
	resp, err := s.client.Do(req, &a)

	for i := range a.Accounts {
		a.Accounts[i].BillingInfo = nil
	}

	return resp, a.Accounts, err
}

// Get returns information about a single account.
// https://docs.recurly.com/api/accounts#get-account
func (s *accountsImpl) Get(code string) (*Response, *Account, error) {
	action := fmt.Sprintf("accounts/%s", code)
	req, err := s.client.NewRequest("GET", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var a Account
	resp, err := s.client.Do(req, &a)
	a.BillingInfo = nil

	return resp, &a, err
}

// Create will create a new account. You may optionally include billing information.
// https://docs.recurly.com/api/accounts#create-account
func (s *accountsImpl) Create(a Account) (*Response, *Account, error) {
	req, err := s.client.NewRequest("POST", "accounts", nil, a)
	if err != nil {
		return nil, nil, err
	}

	var dst Account
	resp, err := s.client.Do(req, &dst)
	dst.BillingInfo = nil

	return resp, &dst, err
}

// Update will update an existing account.
// It's recommended to create a new account object with only the changes you
// want to make. The updated account object will be returned on success.
// https://docs.recurly.com/api/accounts#update-account
func (s *accountsImpl) Update(code string, a Account) (*Response, *Account, error) {
	action := fmt.Sprintf("accounts/%s", code)
	req, err := s.client.NewRequest("PUT", action, nil, a)
	if err != nil {
		return nil, nil, err
	}

	var dst Account
	resp, err := s.client.Do(req, &dst)
	dst.BillingInfo = nil

	return resp, &dst, err
}

// Close marks an account as closed and cancels any active subscriptions. Any
// saved billing information will also be permanently removed from the account.
// https://docs.recurly.com/api/accounts#close-account
func (s *accountsImpl) Close(code string) (*Response, error) {
	action := fmt.Sprintf("accounts/%s", code)
	req, err := s.client.NewRequest("DELETE", action, nil, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Reopen transitions a closed account back to active.
// https://docs.recurly.com/api/accounts#reopen-account
func (s *accountsImpl) Reopen(code string) (*Response, error) {
	action := fmt.Sprintf("accounts/%s/reopen", code)
	req, err := s.client.NewRequest("PUT", action, nil, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// ListNotes returns a list of the notes on an account sorted in descending order.
// https://docs.recurly.com/api/accounts#get-account-notes
func (s *accountsImpl) ListNotes(code string) (*Response, []Note, error) {
	action := fmt.Sprintf("accounts/%s/notes", code)
	req, err := s.client.NewRequest("GET", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var n struct {
		XMLName xml.Name `xml:"notes"`
		Notes   []Note   `xml:"note"`
	}
	resp, err := s.client.Do(req, &n)

	return resp, n.Notes, err
}
