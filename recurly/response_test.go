package recurly

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	// mux is the HTTP request multiplexer used with the test server
	mux *http.ServeMux

	// server is a test HTTP server used to provide mock API responses
	server *httptest.Server

	// client is the Recurly client being tested
	client *Client
)

// setup sets up a test HTTP server along with a recurly.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested
func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// Recurly client configured to use test server
	client = NewClient("test", "abc", nil)
	client.BaseURL = server.URL + "/"
}

func teardown() {
	server.Close()
}

func TestResponse_ConvenienceMethods(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/success", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/client-error", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusForbidden)
	})

	mux.HandleFunc("/server-error", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusInternalServerError)
	})

	tests := []struct {
		endpoint  string
		ok        bool
		err       bool
		clientErr bool
		serverErr bool
	}{
		{endpoint: "success", ok: true},
		{endpoint: "client-error", ok: false, err: true, clientErr: true, serverErr: false},
		{endpoint: "server-error", ok: false, err: true, clientErr: false, serverErr: true},
	}

	for i, tt := range tests {
		req, err := http.NewRequest("GET", client.BaseURL+tt.endpoint, nil)
		if err != nil {
			t.Fatalf("(%d): error creating request for %s. err: %s", i, tt.endpoint, err)
		}

		resp, err := client.client.Do(req)
		if err != nil {
			t.Fatalf("(%d): Error making request for %s. err: %s", i, tt.endpoint, err)
		}

		r := &Response{Response: resp}
		if tt.ok != r.IsOK() {
			t.Fatalf("(%d): Expected ok to be %v for %s, given %v", i, tt.ok, tt.endpoint, r.IsOK())
		} else if tt.err != r.IsError() {
			t.Fatalf("(%d): Expected error to be %v for %s, given %v", i, tt.err, tt.endpoint, r.IsError())
		} else if tt.clientErr != r.IsClientError() {
			t.Fatalf("(%d): Expected clientError to be %v for %s, given %v", i, tt.clientErr, tt.endpoint, r.IsClientError())
		} else if tt.serverErr != r.IsServerError() {
			t.Fatalf("(%d): Expected serverError to be %v for %s, given %v", i, tt.serverErr, tt.endpoint, r.IsServerError())
		}
	}
}

func TestResponse_CursorLinkParsing(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/case0", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Link", "<https://your-subdomain.recurly.com/v2/invoices?cursor=1827545887837797260>; rel=\"next\"")
		rw.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/case1", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Link", "<https://your-subdomain.recurly.com/v2/invoices?state=past_due>; rel=\"start\", <https://your-subdomain.recurly.com/v2/invoices?cursor=-1325183252208393488&state=past_due>; rel=\"prev\", <https://your-subdomain.recurly.com/v2/invoices?cursor=1824642383070236054&state=past_due>; rel=\"next\"")
		rw.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/case2", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Link", "<https://api.recurly.com/v2/accounts?cursor=1234567890&per_page=20>; rel=\"start\", <https://api.recurly.com/v2/accounts?cursor=1234566890&per_page=20>; rel=\"next\"")
		rw.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		endpoint string
		next     string
		prev     string
	}{
		{endpoint: "/case0", next: "1827545887837797260", prev: ""},
		{endpoint: "/case1", next: "1824642383070236054", prev: "-1325183252208393488"},
		{endpoint: "/case2", next: "1234566890", prev: ""},
	}

	for i, tt := range tests {
		req, err := http.NewRequest("GET", client.BaseURL+tt.endpoint, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		resp, err := client.client.Do(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		r := &Response{Response: resp}
		if r.Next() != tt.next {
			t.Fatalf("(%d): Expected next cursor to be %v, got %s, given %v", i, tt.next, r.Next(), tt.endpoint)
		} else if r.Prev() != tt.prev {
			t.Fatalf("(%d): Expected prev cursor to be %v, for %s, given %v", i, tt.prev, r.Prev(), tt.endpoint)
		}
	}
}
