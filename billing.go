package recurly

import (
	"context"
	"encoding/xml"
	"fmt"
	"net"
	"net/http"
)

// BillingService manages the interactions for billing.
type BillingService interface {
	// Get retrieves billing info for an account. If the account does not exist,
	// or the account does not have billing info, a nil billing struct and nil
	// error are returned.
	//
	// https://dev.recurly.com/docs/lookup-an-accounts-billing-info
	Get(ctx context.Context, accountCode string) (*Billing, error)

	// Create creates an account's billing info. To use the recurly.js token (recommended),
	// set b.Token and optionally b.Currency.
	// https://dev.recurly.com/docs/create-an-accounts-billing-info-token
	//
	// To create with credit card, bank info, or other, set the appropriate fields
	// on b. See the following links for specifics:
	// https://dev.recurly.com/docs/create-an-accounts-billing-info-credit-card
	// https://dev.recurly.com/docs/create-an-accounts-billing-info-bank-account
	// https://dev.recurly.com/docs/create-an-accounts-billing-info-using-external-token
	Create(ctx context.Context, accountCode string, b Billing) (*Billing, error)

	// Update updates an account's billing info. To use the recurly.js token (recommended),
	// set b.Token and optionally b.Currency.
	// https://dev.recurly.com/docs/update-an-accounts-billing-info-token
	//
	// To update with credit card, bank info, or other, set the appropriate fields
	// on b. See the following links for specifics:
	// https://dev.recurly.com/docs/update-an-accounts-billing-info-credit-card
	// https://dev.recurly.com/docs/update-an-accounts-billing-info-bank-account
	// https://dev.recurly.com/docs/update-an-accounts-billing-info-using-external-token
	Update(ctx context.Context, accountCode string, b Billing) (*Billing, error)

	// Clear removes stored billing information for an account. If the account has
	// a subscription, the renewal will go into past due unless you update the
	// billing info before the renewal occurs.
	//
	// https://dev.recurly.com/docs/clear-an-accounts-billing-info
	Clear(ctx context.Context, accountCode string) error
}

// Supported card type constants.
const (
	CardTypeAmericanExpress    = "american_express"
	CardTypeDankort            = "dankort"
	CardTypeDinersClub         = "diners_club"
	CardTypeDiscover           = "discover"
	CardTypeForbrugsforeningen = "forbrugsforeningen"
	CardTypeJCB                = "jcb"
	CardTypeLaser              = "laser"
	CardTypeMaestro            = "maestro"
	CardTypeMaster             = "master"
	CardTypeVisa               = "visa"
)

// Billing holds billing info for a single account.
type Billing struct {
	XMLName          xml.Name `xml:"billing_info"`
	FirstName        string   `xml:"first_name,omitempty"`
	LastName         string   `xml:"last_name,omitempty"`
	Company          string   `xml:"company,omitempty"`
	Address          string   `xml:"address1,omitempty"`
	Address2         string   `xml:"address2,omitempty"`
	City             string   `xml:"city,omitempty"`
	State            string   `xml:"state,omitempty"`
	Zip              string   `xml:"zip,omitempty"`
	Country          string   `xml:"country,omitempty"`
	Phone            string   `xml:"phone,omitempty"`
	VATNumber        string   `xml:"vat_number,omitempty"`
	IPAddress        net.IP   `xml:"ip_address,omitempty"`
	IPAddressCountry string   `xml:"ip_address_country,omitempty"`
	PaymentType      string   `xml:"type,attr,omitempty"`

	// Credit Card Info
	FirstSix          string `xml:"first_six,omitempty"`
	LastFour          string `xml:"last_four,omitempty"`
	CardType          string `xml:"card_type,omitempty"`
	Number            int    `xml:"number,omitempty"`
	Month             int    `xml:"month,omitempty"`
	Year              int    `xml:"year,omitempty"`
	VerificationValue int    `xml:"verification_value,omitempty"` // Create/update only

	// Paypal
	PaypalAgreementID string `xml:"paypal_billing_agreement_id,omitempty"`

	// BrainTree
	BrainTreePaymentNonce string `xml:"braintree_payment_nonce,omitempty"`

	// Amazon
	AmazonAgreementID string `xml:"amazon_billing_agreement_id,omitempty"`
	AmazonRegion      string `xml:"amazon_region,omitempty"` // 'eu', 'us', or 'uk'

	// Bank Account
	NameOnAccount string `xml:"name_on_account,omitempty"`
	RoutingNumber string `xml:"routing_number,omitempty"`
	AccountNumber string `xml:"account_number,omitempty"`
	AccountType   string `xml:"account_type,omitempty"`

	ExternalHPPType                 string `xml:"external_hpp_type,omitempty"`                     // only usable with purchases API
	Currency                        string `xml:"currency,omitempty"`                              // Create/update only
	Token                           string `xml:"token_id,omitempty"`                              // Create/update only
	ThreeDSecureActionResultTokenID string `xml:"three_d_secure_action_result_token_id,omitempty"` // Create/update only
	TransactionType                 string `xml:"transaction_type,omitempty"`                      // Create only
}

// Type returns the billing info type. Returns either  "", "bank", or an empty string.
func (b Billing) Type() string {
	if b.FirstSix != "" && b.LastFour != "" && b.Month > 0 && b.Year > 0 {
		return "card"
	} else if b.NameOnAccount != "" && b.RoutingNumber != "" && b.AccountNumber != "" {
		return "bank"
	}
	return ""
}

var _ BillingService = &billingImpl{}

// billingImpl implements BillingService.
type billingImpl serviceImpl

func (s *billingImpl) Get(ctx context.Context, accountCode string) (*Billing, error) {
	path := fmt.Sprintf("/accounts/%s/billing_info", accountCode)
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var dst Billing
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		if e, ok := err.(*ClientError); ok && e.Response.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &dst, nil
}

func (s *billingImpl) Create(ctx context.Context, accountCode string, b Billing) (*Billing, error) {
	path := fmt.Sprintf("/accounts/%s/billing_info", accountCode)
	req, err := s.client.newRequest("POST", path, b)
	if err != nil {
		return nil, err
	}

	var dst Billing
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *billingImpl) Update(ctx context.Context, accountCode string, b Billing) (*Billing, error) {
	path := fmt.Sprintf("/accounts/%s/billing_info", accountCode)
	req, err := s.client.newRequest("PUT", path, b)
	if err != nil {
		return nil, err
	}

	var dst Billing
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *billingImpl) Clear(ctx context.Context, accountCode string) error {
	path := fmt.Sprintf("/accounts/%s/billing_info", accountCode)
	req, err := s.client.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(ctx, req, nil)
	return err
}
