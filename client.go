package recurly

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const defaultBaseURL = "https://%s.com/"

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
	Accounts      AccountsService
	Adjustments   AdjustmentsService
	Billing       BillingService
	Coupons       CouponsService
	Redemptions   RedemptionsService
	Invoices      InvoicesService
	Plans         PlansService
	AddOns        AddOnsService
	Subscriptions SubscriptionsService
	Transactions  TransactionsService
}

// NewClient returns a new instance of *Client.
func NewClient(subDomain, apiKey string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	client := &Client{
		client:    httpClient,
		subDomain: subDomain,
		apiKey:    apiKey,
		BaseURL:   fmt.Sprintf(defaultBaseURL, subDomain),
	}

	client.Accounts = NewAccountsImpl(client)
	client.Adjustments = NewAdjustmentsImpl(client)
	client.Billing = NewBillingImpl(client)
	client.Coupons = NewCouponsImpl(client)
	client.Redemptions = NewRedemptionsImpl(client)
	client.Invoices = NewInvoicesImpl(client)
	client.Plans = NewPlansImpl(client)
	client.AddOns = NewAddOnsImpl(client)
	client.Subscriptions = NewSubscriptionsImpl(client)
	client.Transactions = NewTransactionsImpl(client)

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

	req.SetBasicAuth(c.apiKey, "")
	req.Header.Set("Accept", "application/xml")
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
		if response.StatusCode == 422 {
			var ve struct {
				XMLName          xml.Name          `xml:"errors"`
				Errors           []Error           `xml:"error"`
				Transaction      *Transaction      `xml:"transaction,omitempty"`
				TransactionError *TransactionError `xml:"transaction_error,omitempty"`
			}

			if err = decoder.Decode(&ve); err != nil {
				return response, err
			}

			response.Errors = ve.Errors
			response.Transaction = ve.Transaction
			response.TransactionError = ve.TransactionError
		} else if response.IsClientError() { // Parse possible individual error message
			var ve struct {
				XMLName     xml.Name `xml:"error"`
				Symbol      string   `xml:"symbol"`
				Description string   `xml:"description"`
			}
			if err = decoder.Decode(&ve); err != nil {
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
