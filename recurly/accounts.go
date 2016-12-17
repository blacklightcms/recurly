package recurly

import (
	"encoding/xml"
	"fmt"
	"time"
)

type (
	// AccountsService handles communication with the accounts related methods
	// of the recurly API.
	AccountsService struct {
		client *Client
	}

	// Account represents an individual account on your site
	Account struct {
		XMLName          xml.Name `xml:"account"`
		Code             string   `xml:"account_code,omitempty"`
		State            string   `xml:"state,omitempty"`
		Username         string   `xml:"username,omitempty"`
		Email            string   `xml:"email,omitempty"`
		FirstName        string   `xml:"first_name,omitempty"`
		LastName         string   `xml:"last_name,omitempty"`
		CompanyName      string   `xml:"company_name,omitempty"`
		VATNumber        string   `xml:"vat_number,omitempty"`
		TaxExempt        NullBool `xml:"tax_exempt,omitempty"`
		BillingInfo      *Billing `xml:"billing_info,omitempty"`
		Address          Address  `xml:"address,omitempty"`
		AcceptLanguage   string   `xml:"accept_language,omitempty"`
		HostedLoginToken string   `xml:"hosted_login_token,omitempty"`
		CreatedAt        NullTime `xml:"created_at,omitempty"`
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

// MarshalXML ensures addresses marshal to nil if empty without the need
// to use pointers.
func (a Address) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if a.Address == "" && a.Address2 == "" && a.City == "" && a.State == "" && a.Zip == "" && a.Country == "" && a.Phone == "" {
		return nil
	}

	e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "address"}})
	if a.Address != "" {
		s := xml.StartElement{Name: xml.Name{Local: "address1"}}
		e.EncodeToken(s)
		e.EncodeToken(xml.CharData([]byte(a.Address)))
		e.EncodeToken(xml.EndElement{Name: s.Name})
	}

	if a.Address2 != "" {
		s := xml.StartElement{Name: xml.Name{Local: "address2"}}
		e.EncodeToken(s)
		e.EncodeToken(xml.CharData([]byte(a.Address2)))
		e.EncodeToken(xml.EndElement{Name: s.Name})
	}

	if a.City != "" {
		s := xml.StartElement{Name: xml.Name{Local: "city"}}
		e.EncodeToken(s)
		e.EncodeToken(xml.CharData([]byte(a.City)))
		e.EncodeToken(xml.EndElement{Name: s.Name})
	}

	if a.State != "" {
		s := xml.StartElement{Name: xml.Name{Local: "state"}}
		e.EncodeToken(s)
		e.EncodeToken(xml.CharData([]byte(a.State)))
		e.EncodeToken(xml.EndElement{Name: s.Name})
	}

	if a.Zip != "" {
		s := xml.StartElement{Name: xml.Name{Local: "zip"}}
		e.EncodeToken(s)
		e.EncodeToken(xml.CharData([]byte(a.Zip)))
		e.EncodeToken(xml.EndElement{Name: s.Name})
	}

	if a.Country != "" {
		s := xml.StartElement{Name: xml.Name{Local: "country"}}
		e.EncodeToken(s)
		e.EncodeToken(xml.CharData([]byte(a.Country)))
		e.EncodeToken(xml.EndElement{Name: s.Name})
	}

	if a.Phone != "" {
		s := xml.StartElement{Name: xml.Name{Local: "phone"}}
		e.EncodeToken(s)
		e.EncodeToken(xml.CharData([]byte(a.Phone)))
		e.EncodeToken(xml.EndElement{Name: s.Name})
	}

	e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "address"}})

	return nil
}

const (
	// AccountStateActive is the status for active accounts.
	AccountStateActive = "active"

	// AccountStateClosed is the status for closed accounts.
	AccountStateClosed = "closed"
)

// List returns a list of the accounts on your site.
// https://docs.recurly.com/api/accounts#list-accounts
func (s *AccountsService) List(params Params) (*Response, []Account, error) {
	req, err := s.client.newRequest("GET", "accounts", params, nil)
	if err != nil {
		return nil, nil, err
	}

	var a struct {
		XMLName  xml.Name  `xml:"accounts"`
		Accounts []Account `xml:"account"`
	}
	resp, err := s.client.do(req, &a)

	for i := range a.Accounts {
		a.Accounts[i].BillingInfo = nil
	}

	return resp, a.Accounts, err
}

// Get returns information about a single account.
// https://docs.recurly.com/api/accounts#get-account
func (s *AccountsService) Get(code string) (*Response, *Account, error) {
	action := fmt.Sprintf("accounts/%s", code)
	req, err := s.client.newRequest("GET", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var a Account
	resp, err := s.client.do(req, &a)
	a.BillingInfo = nil

	return resp, &a, err
}

// Create will create a new account. You may optionally include billing information.
// https://docs.recurly.com/api/accounts#create-account
func (s *AccountsService) Create(a Account) (*Response, *Account, error) {
	req, err := s.client.newRequest("POST", "accounts", nil, a)
	if err != nil {
		return nil, nil, err
	}

	var dst Account
	resp, err := s.client.do(req, &dst)
	dst.BillingInfo = nil

	return resp, &dst, err
}

// Update will update an existing account.
// It's recommended to create a new account object with only the changes you
// want to make. The updated account object will be returned on success.
// https://docs.recurly.com/api/accounts#update-account
func (s *AccountsService) Update(code string, a Account) (*Response, *Account, error) {
	action := fmt.Sprintf("accounts/%s", code)
	req, err := s.client.newRequest("PUT", action, nil, a)
	if err != nil {
		return nil, nil, err
	}

	var dst Account
	resp, err := s.client.do(req, &dst)
	dst.BillingInfo = nil

	return resp, &dst, err
}

// Close marks an account as closed and cancels any active subscriptions. Any
// saved billing information will also be permanently removed from the account.
// https://docs.recurly.com/api/accounts#close-account
func (s *AccountsService) Close(code string) (*Response, error) {
	action := fmt.Sprintf("accounts/%s", code)
	req, err := s.client.newRequest("DELETE", action, nil, nil)
	if err != nil {
		return nil, err
	}

	return s.client.do(req, nil)
}

// Reopen transitions a closed account back to active.
// https://docs.recurly.com/api/accounts#reopen-account
func (s *AccountsService) Reopen(code string) (*Response, error) {
	action := fmt.Sprintf("accounts/%s/reopen", code)
	req, err := s.client.newRequest("PUT", action, nil, nil)
	if err != nil {
		return nil, err
	}

	return s.client.do(req, nil)
}

// ListNotes returns a list of the notes on an account sorted in descending order.
// https://docs.recurly.com/api/accounts#get-account-notes
func (s *AccountsService) ListNotes(code string) (*Response, []Note, error) {
	action := fmt.Sprintf("accounts/%s/notes", code)
	req, err := s.client.newRequest("GET", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var n struct {
		XMLName xml.Name `xml:"notes"`
		Notes   []Note   `xml:"note"`
	}
	resp, err := s.client.do(req, &n)

	return resp, n.Notes, err
}
