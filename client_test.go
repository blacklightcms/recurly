package recurly_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blacklightcms/recurly"
)

var (
	// mux is the HTTP request multiplexer used with the test server
	mux *http.ServeMux

	// server is a test HTTP server used to provide mock API responses
	server *httptest.Server

	// client is the Recurly client being tested
	client *recurly.Client
)

// setup sets up a test HTTP server along with a recurly.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested
func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = recurly.NewClient("test", "abc", nil)
	client.BaseURL = server.URL + "/"
}

func teardown() {
	server.Close()
}

// Ensure a 204 No Content is handled properly.
func TestClient_NoContent(t *testing.T) {
	setup()
	defer teardown()

	var invoked bool
	mux.HandleFunc("/v2/transactions", func(w http.ResponseWriter, r *http.Request) {
		invoked = true
		if r.Method != "POST" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	resp, _, err := client.Transactions.Create(recurly.Transaction{})
	if err != nil {
		t.Fatal(err)
	} else if !invoked {
		t.Fatal("expected handler to be invoked")
	} else if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
}
