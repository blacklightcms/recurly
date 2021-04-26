package webhooks

import "github.com/autopilot3/recurly"

const (
	NewShippingAddress = "new_shipping_address_notification"
)

// NewDunningEventNotification is returned for new dunning events.
type NewShippingAddressNotification struct {
	Type         string `xml:"-"`
	Account      Account
	Subscription recurly.Subscription
}
