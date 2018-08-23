package mock

import (
	"net/http"

	"github.com/launchpadcentral/recurly"
)

// NewClient returns a new instance of *recury.Client with the
// services assigned to mocks.
func NewClient(httpClient *http.Client) *recurly.Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	client := recurly.NewClient("a", "b", httpClient)
	client.BaseURL = "https://127.0.0.1/" // Safeguard only

	// Attach mock implementations.
	client.Accounts = &AccountsService{}
	client.Adjustments = &AdjustmentsService{}
	client.Billing = &BillingService{}
	client.Coupons = &CouponsService{}
	client.Redemptions = &RedemptionsService{}
	client.Invoices = &InvoicesService{}
	client.Plans = &PlansService{}
	client.AddOns = &AddOnsService{}
	client.Subscriptions = &SubscriptionsService{}
	client.Transactions = &TransactionsService{}
	client.CreditPayments = &CreditPaymentsService{}
	client.Purchases = &PurchasesService{}

	return client
}
