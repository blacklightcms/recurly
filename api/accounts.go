package api

import (
	"encoding/xml"
	"fmt"

	recurly "github.com/blacklightcms/go-recurly"
)

var _ recurly.AccountsService = &AccountsService{}

// AccountsService handles communication with the accounts related methods
// of the recurly API.
type AccountsService struct {
	client *Client
}

const (
	// AccountStateActive is the status for active accounts.
	AccountStateActive = "active"

	// AccountStateClosed is the status for closed accounts.
	AccountStateClosed = "closed"
)

// List returns a list of the accounts on your site.
// https://docs.recurly.com/api/accounts#list-accounts
func (s *AccountsService) List(params recurly.Params) (*recurly.Response, []recurly.Account, error) {
	req, err := s.client.newRequest("GET", "accounts", params, nil)
	if err != nil {
		return nil, nil, err
	}

	var a struct {
		XMLName  xml.Name          `xml:"accounts"`
		Accounts []recurly.Account `xml:"account"`
	}
	resp, err := s.client.do(req, &a)

	for i := range a.Accounts {
		a.Accounts[i].BillingInfo = nil
	}

	return resp, a.Accounts, err
}

// Get returns information about a single account.
// https://docs.recurly.com/api/accounts#get-account
func (s *AccountsService) Get(code string) (*recurly.Response, *recurly.Account, error) {
	action := fmt.Sprintf("accounts/%s", code)
	req, err := s.client.newRequest("GET", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var a recurly.Account
	resp, err := s.client.do(req, &a)
	a.BillingInfo = nil

	return resp, &a, err
}

// Create will create a new account. You may optionally include billing information.
// https://docs.recurly.com/api/accounts#create-account
func (s *AccountsService) Create(a recurly.Account) (*recurly.Response, *recurly.Account, error) {
	req, err := s.client.newRequest("POST", "accounts", nil, a)
	if err != nil {
		return nil, nil, err
	}

	var dst recurly.Account
	resp, err := s.client.do(req, &dst)
	dst.BillingInfo = nil

	return resp, &dst, err
}

// Update will update an existing account.
// It's recommended to create a new account object with only the changes you
// want to make. The updated account object will be returned on success.
// https://docs.recurly.com/api/accounts#update-account
func (s *AccountsService) Update(code string, a recurly.Account) (*recurly.Response, *recurly.Account, error) {
	action := fmt.Sprintf("accounts/%s", code)
	req, err := s.client.newRequest("PUT", action, nil, a)
	if err != nil {
		return nil, nil, err
	}

	var dst recurly.Account
	resp, err := s.client.do(req, &dst)
	dst.BillingInfo = nil

	return resp, &dst, err
}

// Close marks an account as closed and cancels any active subscriptions. Any
// saved billing information will also be permanently removed from the account.
// https://docs.recurly.com/api/accounts#close-account
func (s *AccountsService) Close(code string) (*recurly.Response, error) {
	action := fmt.Sprintf("accounts/%s", code)
	req, err := s.client.newRequest("DELETE", action, nil, nil)
	if err != nil {
		return nil, err
	}

	return s.client.do(req, nil)
}

// Reopen transitions a closed account back to active.
// https://docs.recurly.com/api/accounts#reopen-account
func (s *AccountsService) Reopen(code string) (*recurly.Response, error) {
	action := fmt.Sprintf("accounts/%s/reopen", code)
	req, err := s.client.newRequest("PUT", action, nil, nil)
	if err != nil {
		return nil, err
	}

	return s.client.do(req, nil)
}

// ListNotes returns a list of the notes on an account sorted in descending order.
// https://docs.recurly.com/api/accounts#get-account-notes
func (s *AccountsService) ListNotes(code string) (*recurly.Response, []recurly.Note, error) {
	action := fmt.Sprintf("accounts/%s/notes", code)
	req, err := s.client.newRequest("GET", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var n struct {
		XMLName xml.Name       `xml:"notes"`
		Notes   []recurly.Note `xml:"note"`
	}
	resp, err := s.client.do(req, &n)

	return resp, n.Notes, err
}
