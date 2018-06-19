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
	Code        string   `xml:"account_code"`
	Username    string   `xml:"username"`
	Email       string   `xml:"email"`
	FirstName   string   `xml:"first_name"`
	LastName    string   `xml:"last_name"`
	CompanyName string   `xml:"company_name"`
	Phone       string   `xml:"phone"`
}
