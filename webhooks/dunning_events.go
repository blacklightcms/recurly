package webhooks

import "github.com/blacklightcms/recurly"

// Dunning event constants.
const (
	NewDunningEvent = "new_dunning_event_notification"
)

// NewDunningEventNotification is returned for new dunning events.
type NewDunningEventNotification struct {
	Type         string `xml:"-"`
	Account      Account
	Invoice      ChargeInvoice
	Subscription recurly.Subscription
}
