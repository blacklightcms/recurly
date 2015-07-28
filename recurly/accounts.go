package recurly

import (
	"encoding/xml"
	"fmt"
	"time"
)

type (
	// accountService handles all interaction with the accounts portion
	// of the recurly API.
	accountService struct {
		client *Client
	}

	// Account represents an individual account on your site
	Account struct {
		XMLName          xml.Name   `xml:"account"`
		Code             string     `xml:"account_code,omitempty"`
		State            string     `xml:"state,omitempty"`
		Username         string     `xml:"username,omitempty"`
		Email            string     `xml:"email,omitempty"`
		FirstName        string     `xml:"first_name,omitempty"`
		LastName         string     `xml:"last_name,omitempty"`
		CompanyName      string     `xml:"company_name,omitempty"`
		VATNumber        string     `xml:"vat_number,omitempty"`
		TaxExempt        NullBool   `xml:"tax_exempt,omitempty"`
		BillingInfo      *Billing   `xml:"billing_info,omitempty"`
		Address          *Address   `xml:"address,omitempty"`
		AcceptLanguage   string     `xml:"accept_language,omitempty"`
		HostedLoginToken string     `xml:"hosted_login_token,omitempty"`
		CreatedAt        *time.Time `xml:"created_at,omitempty"`
	}

	// Address is used for embedded addresses within other structs.
	Address struct {
		Address  string `xml:"address1,omitempty"`
		Address2 string `xml:"address2,omitempty"`
		City     string `xml:"city,omitempty"`
		State    string `xml:"state,omitempty"`
		Zip      string `xml:"zip,omitempty"`
		Country  string `xml:"country,omitempty"`
		Phone    string `xml:"phone,omitempty"`
	}

	// Note holds account notes.
	Note struct {
		XMLName   xml.Name  `xml:"note"`
		Message   string    `xml:"message,omitempty"`
		CreatedAt time.Time `xml:"created_at,omitempty"`
	}
)

// List returns a list of the accounts on your site.
// https://docs.recurly.com/api/accounts#list-accounts
func (as *accountService) List(params Params) (*Response, []Account, error) {
	req, err := as.client.newRequest("GET", "accounts", params, nil)
	if err != nil {
		return nil, nil, err
	}

	var a struct {
		XMLName  xml.Name  `xml:"accounts"`
		Accounts []Account `xml:"account"`
	}
	res, err := as.client.do(req, &a)

	for i, _ := range a.Accounts {
		a.Accounts[i].BillingInfo = nil
	}

	return res, a.Accounts, err
}

// Get returns information about a single account.
// https://docs.recurly.com/api/accounts#get-account
func (as accountService) Get(code string) (*Response, Account, error) {
	action := fmt.Sprintf("accounts/%s", code)
	req, err := as.client.newRequest("GET", action, nil, nil)
	if err != nil {
		return nil, Account{}, err
	}

	var a Account
	res, err := as.client.do(req, &a)

	a.BillingInfo = nil

	return res, a, err
}

// Create will create a new account. You may optionally include billing information.
// https://docs.recurly.com/api/accounts#create-account
func (as accountService) Create(a Account) (*Response, Account, error) {
	req, err := as.client.newRequest("POST", "accounts", nil, a)
	if err != nil {
		return nil, Account{}, err
	}

	var dest Account
	res, err := as.client.do(req, &dest)

	dest.BillingInfo = nil

	return res, dest, err
}

// Update will update an existing account.
// It's recommended to create a new account object with only the changes you
// want to make. The updated account object will be returned on success.
// https://docs.recurly.com/api/accounts#update-account
func (as accountService) Update(code string, a Account) (*Response, Account, error) {
	action := fmt.Sprintf("accounts/%s", code)
	req, err := as.client.newRequest("PUT", action, nil, a)
	if err != nil {
		return nil, Account{}, err
	}

	var dest Account
	res, err := as.client.do(req, &dest)

	dest.BillingInfo = nil

	return res, dest, err
}

// Close marks an account as closed and cancels any active subscriptions. Any
// saved billing information will also be permanently removed from the account.
// https://docs.recurly.com/api/accounts#close-account
func (as accountService) Close(code string) (*Response, error) {
	action := fmt.Sprintf("accounts/%s", code)
	req, err := as.client.newRequest("DELETE", action, nil, nil)
	if err != nil {
		return nil, err
	}

	return as.client.do(req, nil)
}

// Reopen transitions a closed account back to active.
// https://docs.recurly.com/api/accounts#reopen-account
func (as accountService) Reopen(code string) (*Response, error) {
	action := fmt.Sprintf("accounts/%s/reopen", code)
	req, err := as.client.newRequest("PUT", action, nil, nil)
	if err != nil {
		return nil, err
	}

	return as.client.do(req, nil)
}

// ListNotes returns a list of the notes on an account sorted in descending order.
// https://docs.recurly.com/api/accounts#get-account-notes
func (as accountService) ListNotes(code string) (*Response, []Note, error) {
	action := fmt.Sprintf("accounts/%s/notes", code)
	req, err := as.client.newRequest("GET", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var n struct {
		XMLName xml.Name `xml:"notes"`
		Notes   []Note   `xml:"note"`
	}
	res, err := as.client.do(req, &n)

	return res, n.Notes, err
}
