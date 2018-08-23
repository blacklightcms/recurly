package webhooks

import "github.com/launchpadcentral/recurly"

const (
    NewDunningEvent = "new_dunning_event_notification"
)

type NewDunningEventNotification struct {
    Type string `xml:"-"`
    Account Account
    Invoice ChargeInvoice
    Subscription recurly.Subscription
}

type NewDunningEventDeprecatedNotification struct {
    Type string `xml:"-"`
    Account Account
    Invoice Invoice
    Subscription recurly.Subscription
}