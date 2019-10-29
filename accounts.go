package recurly

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
)

// AccountsService manages the interactions for accounts.
type AccountsService interface {
	// List returns pager to paginate accounts. PagerOptions are used to optionally
	// filter the results.
	//
	// https://dev.recurly.com/docs/list-accounts
	List(opts *PagerOptions) Pager

	// Get retrieves an account. If the account does not exist,
	// a nil account and nil error is returned.
	//
	// https://dev.recurly.com/docs/get-account
	Get(ctx context.Context, accountCode string) (*Account, error)

	// Balance retrieves the balance for an account.
	//
	// https://dev.recurly.com/docs/lookup-account-balance
	Balance(ctx context.Context, accountCode string) (*AccountBalance, error)

	// Create creates a new account. You may optionally including billing information.
	//
	// https://dev.recurly.com/docs/create-an-account
	Create(ctx context.Context, a Account) (*Account, error)

	// Update updates an account. You may optionally including billing information.
	//
	// https://dev.recurly.com/docs/update-account
	Update(ctx context.Context, accountCode string, a Account) (*Account, error)

	// Close marks an account as closed and cancels any active subscriptions.
	// Paused subscriptions will be expired.
	//
	// https://dev.recurly.com/docs/close-account
	Close(ctx context.Context, accountCode string) error

	// Reopen transitions a closed account back to active.
	//
	// https://dev.recurly.com/docs/reopen-account
	Reopen(ctx context.Context, accountCode string) error

	// ListNotes returns a pager to paginate notes for an account. PagerOptions is used
	// to optionally filter the results.
	//
	// https://dev.recurly.com/docs/list-account-notes
	ListNotes(accountCode string, params *PagerOptions) Pager
}

// Account constants.
const (
	AccountStateActive = "active"
	AccountStateClosed = "closed"
)

// An Account is core to managing your customers inside of Recurly. The account object
// stores the entire Recurly history of your customer and acts as the entry point
// for working with a customer's billing information, subscription data, transactions,
// invoices and more.
// https://dev.recurly.com/docs/account-object
type Account struct {
	XMLName                 xml.Name           `xml:"account"`
	Code                    string             `xml:"account_code,omitempty"`
	State                   string             `xml:"state,omitempty"`
	Username                string             `xml:"username,omitempty"`
	Email                   string             `xml:"email,omitempty"`
	CCEmails                []string           `xml:"cc_emails,omitempty"`
	FirstName               string             `xml:"first_name,omitempty"`
	LastName                string             `xml:"last_name,omitempty"`
	BillingInfo             *Billing           `xml:"billing_info,omitempty"`
	CompanyName             string             `xml:"company_name,omitempty"`
	VATNumber               string             `xml:"vat_number,omitempty"`
	TaxExempt               NullBool           `xml:"tax_exempt,omitempty"`
	Address                 *Address           `xml:"address,omitempty"`
	ShippingAddresses       *[]ShippingAddress `xml:"shipping_addresses>shipping_address,omitempty"`
	AcceptLanguage          string             `xml:"accept_language,omitempty"`
	HostedLoginToken        string             `xml:"hosted_login_token,omitempty"`
	CreatedAt               NullTime           `xml:"created_at,omitempty"`
	UpdatedAt               NullTime           `xml:"updated_at,omitempty"`
	ClosedAt                NullTime           `xml:"closed_at,omitempty"`
	HasLiveSubscription     NullBool           `xml:"has_live_subscription,omitempty"`
	HasActiveSubscription   NullBool           `xml:"has_active_subscription,omitempty"`
	HasFutureSubscription   NullBool           `xml:"has_future_subscription,omitempty"`
	HasCanceledSubscription NullBool           `xml:"has_canceled_subscription,omitempty"`
	HasPausedSubscription   NullBool           `xml:"has_paused_subscription,omitempty"`
	HasPastDueInvoice       NullBool           `xml:"has_past_due_invoice,omitempty"`
	PreferredLocale         string             `xml:"preferred_locale,omitempty"`
	CustomFields            *CustomFields      `xml:"custom_fields,omitempty"`
	TransactionType         string             `xml:"transaction_type,omitempty"` // Create only
}

// AccountBalance is used for getting the account balance.
type AccountBalance struct {
	XMLName xml.Name   `xml:"account_balance"`
	PastDue bool       `xml:"past_due"`
	Balance UnitAmount `xml:"balance_in_cents"`
}

// Address is used for embedded addresses within other structs.
type Address struct {
	XMLName       xml.Name `xml:"address"`
	NameOnAccount string   `xml:"name_on_account,omitempty"`
	FirstName     string   `xml:"first_name,omitempty"`
	LastName      string   `xml:"last_name,omitempty"`
	Company       string   `xml:"company,omitempty"`
	Address       string   `xml:"address1,omitempty"`
	Address2      string   `xml:"address2,omitempty"`
	City          string   `xml:"city,omitempty"`
	State         string   `xml:"state,omitempty"`
	Zip           string   `xml:"zip,omitempty"`
	Country       string   `xml:"country,omitempty"`
	Phone         string   `xml:"phone,omitempty"`
}

// Note holds account notes.
type Note struct {
	XMLName   xml.Name  `xml:"note"`
	Message   string    `xml:"message,omitempty"`
	CreatedAt time.Time `xml:"created_at,omitempty"`
}

var _ AccountsService = &accountsImpl{}

// accountsImpl implements AccountsService.
type accountsImpl serviceImpl

func (s *accountsImpl) List(opts *PagerOptions) Pager {
	return s.client.newPager("GET", "/accounts", opts)
}

func (s *accountsImpl) Get(ctx context.Context, code string) (*Account, error) {
	path := fmt.Sprintf("/accounts/%s", code)
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var dst Account
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		if e, ok := err.(*ClientError); ok && e.Response.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &dst, nil
}

func (s *accountsImpl) Balance(ctx context.Context, code string) (*AccountBalance, error) {
	path := fmt.Sprintf("/accounts/%s/balance", code)
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var b AccountBalance
	if _, err := s.client.do(ctx, req, &b); err != nil {
		return nil, err
	}
	return &b, err
}

func (s *accountsImpl) Create(ctx context.Context, a Account) (*Account, error) {
	req, err := s.client.newRequest("POST", "/accounts", a)
	if err != nil {
		return nil, err
	}

	var dst Account
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, err
}

func (s *accountsImpl) Update(ctx context.Context, code string, a Account) (*Account, error) {
	path := fmt.Sprintf("/accounts/%s", code)
	req, err := s.client.newRequest("PUT", path, a)
	if err != nil {
		return nil, err
	}

	var dst Account
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, err
}

func (s *accountsImpl) Close(ctx context.Context, code string) error {
	path := fmt.Sprintf("/accounts/%s", code)
	req, err := s.client.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(ctx, req, nil)
	return err
}

func (s *accountsImpl) Reopen(ctx context.Context, code string) error {
	path := fmt.Sprintf("/accounts/%s/reopen", code)
	req, err := s.client.newRequest("PUT", path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(ctx, req, nil)
	return err
}

func (s *accountsImpl) ListNotes(accountCode string, params *PagerOptions) Pager {
	path := fmt.Sprintf("/accounts/%s/notes", accountCode)
	return s.client.newPager("GET", path, params)
}
