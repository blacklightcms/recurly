package api

import (
	"net/http"

	recurly "github.com/blacklightcms/go-recurly"
)

// NewClient creates a new Recurly API client.
func NewClient(subDomain string, apiKey string, httpClient *http.Client) *recurly.Client {
	c := recurly.NewClient(subDomain, apiKey, nil)

	c.Accounts = &AccountsService{client: c}
	c.Adjustments = &AdjustmentsService{client: c}
	c.Billing = &BillingService{client: c}
	c.Coupons = &CouponsService{client: c}
	c.Redemptions = &RedemptionsService{client: c}
	c.Invoices = &InvoicesService{client: c}
	c.Plans = &PlansService{client: c}
	c.AddOns = &AddOnsService{client: c}
	c.Subscriptions = &SubscriptionsService{client: c}
	c.Transactions = &TransactionsService{client: c}

	return c
}
