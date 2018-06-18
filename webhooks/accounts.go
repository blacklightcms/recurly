package webhooks

import "encoding/xml"

// Account notifications.
// https://dev.recurly.com/page/webhooks#account-notifications
const (
	BillingInfoUpdated = "billing_info_updated_notification"
)

// AccountNotification is returned for all account notifications.
type AccountNotification struct {
	Type    string  `xml:"-"`
	Account Account `xml:"account"`
}

// Account represents the account object sent in webhooks.
type Account struct {
	XMLName     xml.Name `xml:"account"`
	Code        string   `xml:"account_code,omitempty"`
	Username    string   `xml:"username,omitempty"`
	Email       string   `xml:"email,omitempty"`
	FirstName   string   `xml:"first_name,omitempty"`
	LastName    string   `xml:"last_name,omitempty"`
	CompanyName string   `xml:"company_name,omitempty"`
	Phone       string   `xml:"phone,omitempty"`
}
