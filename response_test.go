package recurly_test

import (
	"net/http"
	"testing"

	"github.com/launchpadcentral/recurly"
)

func TestResponse_ConvenienceMethods(t *testing.T) {
	// Success response.
	resp0 := &recurly.Response{
		Response: &http.Response{
			StatusCode: http.StatusOK,
		},
	}
	if !resp0.IsOK() {
		t.Fatal("expected IsOK to be true")
	} else if resp0.IsError() {
		t.Fatal("expected IsError to be false")
	} else if resp0.IsClientError() {
		t.Fatal("expected IsClientError to be false")
	} else if resp0.IsServerError() {
		t.Fatal("expected IsServerError to be false")
	}

	// Client error.
	resp1 := &recurly.Response{
		Response: &http.Response{
			StatusCode: http.StatusForbidden,
		},
	}
	if resp1.IsOK() {
		t.Fatal("expected IsOK to be false")
	} else if !resp1.IsError() {
		t.Fatal("expected IsError to be true")
	} else if !resp1.IsClientError() {
		t.Fatal("expected IsClientError to be true")
	} else if resp0.IsServerError() {
		t.Fatal("expected IsServerError to be false")
	}

	// Server error.
	resp2 := &recurly.Response{
		Response: &http.Response{
			StatusCode: http.StatusInternalServerError,
		},
	}
	if resp2.IsOK() {
		t.Fatal("expected IsOK to be false")
	} else if !resp2.IsError() {
		t.Fatal("expected IsError to be true")
	} else if resp0.IsClientError() {
		t.Fatal("expected IsClientError to be false")
	} else if !resp2.IsServerError() {
		t.Fatal("expected IsServerError to be true")
	}
}

func TestResponse_CursorLinkParsing(t *testing.T) {
	resp0 := &recurly.Response{
		Response: &http.Response{
			StatusCode: http.StatusOK, // Prev/Next methods require resp.IsOK() == true
			Header: http.Header{
				"Link": {"<https://your-subdomain.recurly.com/v2/invoices?cursor=1827545887837797260>; rel=\"next\""},
			},
		},
	}
	if resp0.Prev() != "" {
		t.Fatalf("unexpected previous: %s", resp0.Prev())
	} else if resp0.Next() != "1827545887837797260" {
		t.Fatalf("unexpected next: %s", resp0.Next())
	}

	resp1 := &recurly.Response{
		Response: &http.Response{
			StatusCode: http.StatusOK, // Prev/Next methods require resp.IsOK() == true
			Header: http.Header{
				"Link": {"<https://your-subdomain.recurly.com/v2/invoices?state=past_due>; rel=\"start\", <https://your-subdomain.recurly.com/v2/invoices?cursor=-1325183252208393488&state=past_due>; rel=\"prev\", <https://your-subdomain.recurly.com/v2/invoices?cursor=1824642383070236054&state=past_due>; rel=\"next\""},
			},
		},
	}
	if resp1.Prev() != "-1325183252208393488" {
		t.Fatalf("unexpected previous: %s", resp1.Prev())
	} else if resp1.Next() != "1824642383070236054" {
		t.Fatalf("unexpected next: %s", resp1.Next())
	}

	resp2 := &recurly.Response{
		Response: &http.Response{
			StatusCode: http.StatusOK, // Prev/Next methods require resp.IsOK() == true
			Header: http.Header{
				"Link": {"<https://api.recurly.com/v2/accounts?cursor=1234567890&per_page=20>; rel=\"start\", <https://api.recurly.com/v2/accounts?cursor=1234566890&per_page=20>; rel=\"next\""},
			},
		},
	}
	if resp2.Prev() != "" {
		t.Fatalf("unexpected previous: %s", resp2.Prev())
	} else if resp2.Next() != "1234566890" {
		t.Fatalf("unexpected next: %s", resp2.Next())
	}
}
