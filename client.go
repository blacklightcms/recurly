package recurly

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"strings"
)

const defaultBaseURL = "https://%s.recurly.com/"

// Client manages communication with the Recurly API.
type Client struct {
	// client is the HTTP Client used to communicate with the API.
	client *http.Client

	// subdomain is your account's sub domain used for authentication.
	subDomain string

	// apiKey is your account's API key used for authentication.
	apiKey string

	// BaseURL is the base url for api requests.
	BaseURL string

	// Services used for talking with different parts of the Recurly API
	Accounts          AccountsService
	Adjustments       AdjustmentsService
	Billing           BillingService
	Coupons           CouponsService
	Redemptions       RedemptionsService
	Invoices          InvoicesService
	Plans             PlansService
	AddOns            AddOnsService
	ShippingAddresses ShippingAddressesService
	Subscriptions     SubscriptionsService
	Transactions      TransactionsService
	CreditPayments    CreditPaymentsService
	Purchases         PurchasesService
}

// NewClient returns a new instance of *Client.
// apiKey should be everything after "Basic ".
func NewClient(subDomain, apiKey string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	client := &Client{
		client:    httpClient,
		subDomain: subDomain,
		apiKey:    base64.StdEncoding.EncodeToString([]byte(apiKey)),
		BaseURL:   fmt.Sprintf(defaultBaseURL, subDomain),
	}

	client.Accounts = &accountsImpl{client: client}
	client.Adjustments = &adjustmentsImpl{client: client}
	client.Billing = &billingImpl{client: client}
	client.Coupons = &couponsImpl{client: client}
	client.Redemptions = &redemptionsImpl{client: client}
	client.Invoices = &invoicesImpl{client: client}
	client.Plans = &plansImpl{client: client}
	client.AddOns = &addOnsImpl{client: client}
	client.Subscriptions = &subscriptionsImpl{client: client}
	client.ShippingAddresses = &shippingAddressesImpl{client: client}
	client.Transactions = &transactionsImpl{client: client}
	client.CreditPayments = &creditInvoicesImpl{client: client}
	client.Purchases = &purchasesImpl{client: client}

	return client
}

// newRequest creates an authenticated API request that is ready to send.
func (c *Client) newRequest(method string, action string, params Params, body interface{}) (*http.Request, error) {
	method = strings.ToUpper(method)
	endpoint := fmt.Sprintf("%sv2/%s", c.BaseURL, action)

	// Query String
	qs := url.Values{}
	for k, v := range params {
		qs.Add(k, fmt.Sprintf("%v", v))
	}

	if len(qs) > 0 {
		endpoint += "?" + qs.Encode()
	}

	// Request body
	var buf bytes.Buffer
	if body != nil {
		err := xml.NewEncoder(&buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, endpoint, &buf)
	if err != nil {
		return nil, err
	}

	// Add User-Agent tracking for Recurly statistics and potentially
	// identifying bugs or updates needed in the library.
	// https://github.com/blacklightcms/recurly/issues/41
	req.Header.Set("User-Agent", fmt.Sprintf(
		"Blacklight/2018-06-05; Go (%s) [%s-%s]",
		runtime.Version(),
		runtime.GOARCH,
		runtime.GOOS,
	))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", c.apiKey))
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("X-Api-Version", "2.13")
	if req.Method == "POST" || req.Method == "PUT" {
		req.Header.Set("Content-Type", "application/xml; charset=utf-8")
	}

	return req, err
}

// do takes a prepared API request and makes the API call to Recurly.
// It will decode the XML into a destination struct you provide as well
// as parse any validation errors that may have occurred.
// It returns a Response object that provides a wrapper around http.Response
// with some convenience methods.
func (c *Client) do(req *http.Request, v interface{}) (*Response, error) {
	req.Close = true
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := &Response{Response: resp}
	decoder := xml.NewDecoder(resp.Body)
	if response.IsError() { // Parse validation errors
		if response.StatusCode == http.StatusUnprocessableEntity {
			var ve struct {
				Errors      []Error      `xml:"error"`
				Transaction *Transaction `xml:"transaction,omitempty"`

				// At least one 422 response can return a single error instead of an array.
				// https://dev.recurly.com/docs/welcome#section-422-unprocessable-entity-responses
				Symbol      string `xml:"symbol"`
				Description string `xml:"description"`
			}

			if err = decoder.Decode(&ve); err != nil {
				return response, err
			}

			if ve.Errors == nil {
				// If the response returned single error, set error as the first error in array.
				response.Errors = []Error{{
					XMLName:     xml.Name{Local: "error"},
					Symbol:      ve.Symbol,
					Description: ve.Description,
				}}
			} else {
				response.Errors = ve.Errors
			}

			// If the response object includes a TransactionError, set the
			// transaction field on the response object and the TransactionError field.
			if ve.Transaction != nil {
				response.transaction = ve.Transaction
			}
		} else if response.IsClientError() { // Parse possible individual error message
			var ve struct {
				XMLName     xml.Name `xml:"error"`
				Symbol      string   `xml:"symbol"`
				Description string   `xml:"description"`
			}
			if err = decoder.Decode(&ve); err == io.EOF {
				return response, nil
			} else if err != nil {
				return response, err
			}

			response.Errors = []Error{
				{
					Symbol:  ve.Symbol,
					Message: ve.Description,
				},
			}
		}

		return response, nil
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = decoder.Decode(&v)
		}
	}

	return response, err
}
