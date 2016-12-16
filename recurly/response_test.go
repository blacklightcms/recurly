package recurly

import (
	"net/http"
	"net/http/httptest"
	"strconv"
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

func TestResponseConvenienceMethods(t *testing.T) {
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

	suite := []map[string]string{
		map[string]string{"endpoint": "success", "ok": "true", "error": "false", "clientError": "false", "serverError": "false"},
		map[string]string{"endpoint": "client-error", "ok": "false", "error": "true", "clientError": "true", "serverError": "false"},
		map[string]string{"endpoint": "server-error", "ok": "false", "error": "true", "clientError": "false", "serverError": "true"},
	}

	for i, s := range suite {
		req, err := http.NewRequest("GET", client.BaseURL+s["endpoint"], nil)
		if err != nil {
			t.Fatalf("TestResponse Error (%d): Error creating request for %s. err: %s", i, s["endpoint"], err)
		}

		resp, err := client.client.Do(req)
		if err != nil {
			t.Fatalf("TestResponse Error (%d): Error making request for %s. err: %s", i, s["endpoint"], err)
		}

		r := &Response{Response: resp}
		expected, _ := strconv.ParseBool(s["ok"])
		if expected != r.IsOK() {
			t.Errorf("TestResponse Error (%d): Expected ok to be %v for %s, given %v", i, expected, s["endpoint"], r.IsOK())
		}

		expected, _ = strconv.ParseBool(s["error"])
		if expected != r.IsError() {
			t.Errorf("TestResponse Error (%d): Expected error to be %v for %s, given %v", i, expected, s["endpoint"], r.IsError())
		}

		expected, _ = strconv.ParseBool(s["clientError"])
		if expected != r.IsClientError() {
			t.Errorf("TestResponse Error (%d): Expected clientError to be %v for %s, given %v", i, expected, s["endpoint"], r.IsClientError())
		}

		expected, _ = strconv.ParseBool(s["serverError"])
		if expected != r.IsServerError() {
			t.Errorf("TestResponse Error (%d): Expected serverError to be %v for %s, given %v", i, expected, s["endpoint"], r.IsServerError())
		}
	}
}

func TestPaginationLinks(t *testing.T) {
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

	suite := []map[string]string{
		{"endpoint": "/case0", "next": "1827545887837797260", "prev": ""},
		{"endpoint": "/case1", "next": "1824642383070236054", "prev": "-1325183252208393488"},
		{"endpoint": "/case2", "next": "1234566890", "prev": ""},
	}

	for i, s := range suite {
		req, err := http.NewRequest("GET", client.BaseURL+s["endpoint"], nil)
		if err != nil {
			t.Fatalf("TestPaginationLinks Error (%d): Error creating request for %s. err: %s", i, s["endpoint"], err)
		}

		resp, err := client.client.Do(req)
		if err != nil {
			t.Fatalf("TestPaginationLinks Error (%d): Error making request for %s. err: %s", i, s["endpoint"], err)
		}

		r := &Response{Response: resp}

		if r.Next() != s["next"] {
			t.Errorf("TestPaginationLinks Error (%d): Expected next cursor to be %v, got %s, given %v", i, s["next"], r.Next(), s["endpoint"])
		} else if r.Prev() != s["prev"] {
			t.Errorf("TestPaginationLinks Error (%d): Expected prev cursor to be %v, for %s, given %v", i, s["prev"], r.Prev(), s["endpoint"])
		}
	}
}
