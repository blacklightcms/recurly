package recurly

import (
	"encoding/xml"
	"net/http"
)

type (
	// Response is returned for each API call.
	Response struct {
		*http.Response

		// Errors holds an array of validation errors if any occurred.
		Errors []Error
	}

	// Error is an individual validation error
	Error struct {
		XMLName xml.Name `xml:"error"`
		Message string   `xml:",innerxml"`
		Field   string   `xml:"field,attr"`
		Symbol  string   `xml:"symbol,attr"`
	}
)

// IsOK returns true if the request was successful.
func (r Response) IsOK() bool {
	return r.Response.StatusCode >= 200 && r.Response.StatusCode <= 299
}

// IsError returns true if the request was not successful.
func (r Response) IsError() bool {
	return !r.IsOK()
}

// IsClientError returns true if the request resulted in a 400-499 status code.
func (r Response) IsClientError() bool {
	return r.Response.StatusCode >= 400 && r.Response.StatusCode <= 499
}

// IsServerError returns true if the request resulted in a 500-599 status code --
// indicating you may want to retry the request later.
func (r Response) IsServerError() bool {
	return r.Response.StatusCode >= 500 && r.Response.StatusCode <= 599
}
