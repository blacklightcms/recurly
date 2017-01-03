package mock

import (
	"net/http"

	recurly "github.com/blacklightcms/go-recurly"
)

// NewClient sets the unexported fields on the struct and returns a Client.
func NewClient(httpClient *http.Client) *recurly.Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	client := recurly.NewClient("a", "b", httpClient)

	client.Subscriptions = &SubscriptionsService{}
	client.Transactions = &TransactionsService{}

	return client
}
