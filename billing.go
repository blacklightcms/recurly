package recurly

import (
	"encoding/xml"
	"net"
)

type (
	// Billing represents billing info for a single account on your site
	Billing struct {
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

		// Credit Card Info
		FirstSix int    `xml:"first_six,omitempty"`
		LastFour int    `xml:"last_four,omitempty"`
		CardType string `xml:"card_type,omitempty"`
		Number   int    `xml:"number,omitempty"`
		Month    int    `xml:"month,omitempty"`
		Year     int    `xml:"year,omitempty"`
		// VerificationValue is only used for create/update only. A Verification
		// Value will never be returned on read.
		VerificationValue int `xml:"verification_value,omitempty"`

		// Paypal
		PaypalAgreementID string `xml:"paypal_billing_agreement_id,omitempty"`

		// Amazon
		AmazonAgreementID string `xml:"amazon_billing_agreement_id,omitempty"`

		// Bank Account
		// Note: routing numbers and account numbers may start with zeros, so need
		// to treat them as strings
		NameOnAccount string `xml:"name_on_account,omitempty"`
		RoutingNumber string `xml:"routing_number,omitempty"`
		AccountNumber string `xml:"account_number,omitempty"`
		AccountType   string `xml:"account_type,omitempty"`

		// Token is used for create/update only. A token will never be returned
		// on read.
		Token string `xml:"token_id,omitempty"`
	}
)

// Type returns the billing info type. Currently options: card, bank, ""
func (b Billing) Type() string {
	if b.FirstSix > 0 && b.LastFour > 0 && b.Month > 0 && b.Year > 0 {
		return "card"
	} else if b.NameOnAccount != "" && b.RoutingNumber != "" && b.AccountNumber != "" {
		return "bank"
	}

	return ""
}
