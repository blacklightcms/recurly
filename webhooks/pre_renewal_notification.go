package webhooks

import "github.com/autopilot3/recurly"

const (
	PreRenewal = "prerenewal_notification"
)

// NewDunningEventNotification is returned for new dunning events.
type PreRenewalNotification struct {
	Type         string `xml:"-"`
	Account      Account
	Subscription recurly.Subscription
}
