package recurly

import (
	"encoding/xml"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// Response is returned for each API call.
type Response struct {
	*http.Response

	// Errors holds an array of validation errors if any occurred.
	Errors []Error

	// Transaction holds the transaction returned with a transaction error.
	// This will be populated when creating a new subscription if the payment fails.
	Transaction *Transaction

	// TransactionError holds transaction errors from your payment gateway.
	// This will only be populated when creating a new subscription,
	// updating billing information, and processing a one-time transaction.
	// https://recurly.readme.io/v2.0/page/transaction-errors
	TransactionError *TransactionError
}

var (
	// rxPaginationLink is a regex to parse prev/next links from the Link header
	rxPaginationLink = regexp.MustCompile(`<([^>]+)>;`)
)

// IsOK returns true if the request was successful.
func (r *Response) IsOK() bool {
	return r.Response.StatusCode >= 200 && r.Response.StatusCode <= 299
}

// IsError returns true if the request was not successful.
func (r *Response) IsError() bool {
	return !r.IsOK()
}

// IsClientError returns true if the request resulted in a 400-499 status code.
func (r *Response) IsClientError() bool {
	return r.Response.StatusCode >= 400 && r.Response.StatusCode <= 499
}

// IsServerError returns true if the request resulted in a 500-599 status code --
// indicating you may want to retry the request later.
func (r *Response) IsServerError() bool {
	return r.Response.StatusCode >= 500 && r.Response.StatusCode <= 599
}

// Prev returns the cursor for the previous page of paginated results. If no
// previous page exists, an empty string is returned.
func (r *Response) Prev() string {
	if !r.IsOK() || r.Header.Get("Link") == "" {
		return ""
	}

	links := strings.Split(r.Header.Get("Link"), ",")
	for _, link := range links {
		if strings.HasSuffix(link, `rel="prev"`) {
			match := rxPaginationLink.FindStringSubmatch(link)
			if len(match) == 2 {
				u, err := url.Parse(match[1])
				if err != nil {
					return ""
				}
				return u.Query().Get("cursor")
			}
		}
	}

	return ""
}

// Next returns the cursor for the next page of paginated results. If no
// next page exists, an empty string is returned.
func (r *Response) Next() string {
	if !r.IsOK() || r.Header.Get("Link") == "" {
		return ""
	}

	links := strings.Split(r.Header.Get("Link"), ",")
	for _, link := range links {
		if strings.HasSuffix(link, `rel="next"`) {
			match := rxPaginationLink.FindStringSubmatch(link)
			if len(match) == 2 {
				u, err := url.Parse(match[1])
				if err != nil {
					return ""
				}
				return u.Query().Get("cursor")
			}
		}
	}

	return ""
}

// Error is an individual validation error
type Error struct {
	XMLName xml.Name `xml:"error"`
	Message string   `xml:",innerxml"`
	Field   string   `xml:"field,attr"`
	Symbol  string   `xml:"symbol,attr"`
}

// TransactionError is an error encounted from your payment gateway that
// recurly has standardized.
// https://recurly.readme.io/v2.0/page/transaction-errors
type TransactionError struct {
	XMLName          xml.Name `xml:"transaction_error"`
	ErrorCode        string   `xml:"error_code,omitempty"`
	ErrorCategory    string   `xml:"error_category,omitempty"`
	MerchantMessage  string   `xml:"merchant_message,omitempty"`
	CustomerMessage  string   `xml:"customer_message,omitempty"`
	GatewayErrorCode string   `xml:"gateway_error_code,omitempty"`
}
