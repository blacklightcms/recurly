package mock

import (
	"net/http"

	"github.com/blacklightcms/recurly"
)

// NewClient returns a new instance of *recurly.Client with the
// services assigned to mocks.
func NewClient(httpClient *http.Client) *recurly.Client {
	client := recurly.NewClient("a", "b")
	if httpClient != nil {
		client.Client = httpClient
	}

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
	client.ShippingAddresses = &ShippingAddressesService{}
	client.Transactions = &TransactionsService{}
	client.CreditPayments = &CreditPaymentsService{}
	client.Purchases = &PurchasesService{}

	return client
}
