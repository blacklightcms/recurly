package webhooks

import "github.com/blacklightcms/recurly"

// Subscription notifications.
// https://dev.recurly.com/page/webhooks#subscription-notifications
const (
	NewSubscription      = "new_subscription_notification"
	UpdatedSubscription  = "updated_subscription_notification"
	RenewedSubscription  = "renewed_subscription_notification"
	ExpiredSubscription  = "expired_subscription_notification"
	CanceledSubscription = "canceled_subscription_notification"
	ReactivatedAccount   = "reactivated_account_notification"
)

// SubscriptionNotification is returned for all subscription notifications.
type SubscriptionNotification struct {
	Type         string               `xml:"-"`
	Account      Account              `xml:"account"`
	Subscription recurly.Subscription `xml:"subscription"`
}
