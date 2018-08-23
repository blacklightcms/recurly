package recurly_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/launchpadcentral/recurly"
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
