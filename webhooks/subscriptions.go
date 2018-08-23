package webhooks

import "github.com/launchpadcentral/recurly"

// Subscription notifications.
// https://dev.recurly.com/page/webhooks#subscription-notifications
const (
	NewSubscription            = "new_subscription_notification"
	UpdatedSubscription        = "updated_subscription_notification"
	RenewedSubscription        = "renewed_subscription_notification"
	ExpiredSubscription        = "expired_subscription_notification"
	CanceledSubscription       = "canceled_subscription_notification"
	PausedSubscription         = "subscription_paused_notification"
	ResumedSubscription        = "subscription_resumed_notification"
	ScheduledPauseSubscription = "scheduled_subscription_pause_notification"
	ModifiedPauseSubscription  = "subscription_pause_modified_notification"
	PausedRenewalSubscription  = "paused_subscription_renewal_notification"
	PauseCanceledSubscription  = "subscription_pause_canceled_notification"
	ReactivatedAccount         = "reactivated_account_notification"
)

// SubscriptionNotification is returned for all subscription notifications.
type SubscriptionNotification struct {
	Type         string               `xml:"-"`
	Account      Account              `xml:"account"`
	Subscription recurly.Subscription `xml:"subscription"`
}
