package recurly_test

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	recurly "github.com/blacklightcms/go-recurly"
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

func TestClient_NewRequest(t *testing.T) {
	setup()
	defer teardown()

	req, err := client.NewRequest("GET", "accounts/14579", recurly.Params{"foo": "bar"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if req.URL.Path != "/v2/accounts/14579" {
		t.Fatalf("unexpected path: %s", req.URL.Path)
	} else if req.Method != "GET" {
		t.Fatalf("unexpected method: %s", req.Method)
	} else if req.Header.Get("Accept") != "application/xml" {
		t.Fatalf("unexpected Accept header: %s", req.Header.Get("Accept"))
	} else if req.Header.Get("Content-Type") != "" {
		t.Fatalf("unexpected Content-Type header: %s", req.Header.Get("Content-Type"))
	}

	query := req.URL.Query()
	for name, expected := range map[string]string{"foo": "bar"} {
		actual := query.Get(name)
		if actual != expected {
			t.Fatalf("expected '%s' to equal '%s', actual '%s'", name, expected, actual)
		}
	}

	req, err = client.NewRequest("PUT", "accounts/abc", nil, recurly.Account{Code: "abc"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if req.URL.Path != "/v2/accounts/abc" {
		t.Fatalf("unexpected path: %s", req.URL.Path)
	} else if req.Method != "PUT" {
		t.Fatalf("unexpected method: %s", req.Method)
	} else if req.Header.Get("Accept") != "application/xml" {
		t.Fatalf("unexpected Accept header: %s", req.Header.Get("Accept"))
	} else if req.Header.Get("Content-Type") != "application/xml; charset=utf-8" {
		t.Fatalf("unexpected Content-Type header: %s", req.Header.Get("Content-Type"))
	}

	expected := []byte("<account><account_code>abc</account_code></account>")
	given, _ := ioutil.ReadAll(req.Body)
	if !reflect.DeepEqual(expected, given) {
		t.Fatalf("expected string body equal to %s, given %s", expected, given)
	}

	query = req.URL.Query()
	if len(query) != 0 {
		t.Fatalf("expected %d query Params, given %d", 0, len(query))
	}
}

func TestClient_Error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(422)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
			<errors>
				<error field="model_name.field_name" symbol="not_a_number" lang="en-US">is not a number</error>
				<error field="foo.bar" symbol="not_good" lang="en-US">is not good</error>
			</errors>`)
	})

	req, err := http.NewRequest("GET", client.BaseURL+"error", nil)
	if err != nil {
		t.Fatalf("error creating request. err: %v", err)
	}

	resp, err := client.Do(req, nil)
	if err != nil {
		t.Fatalf("error making request. err: %v", err)
	} else if resp.IsOK() {
		t.Fatalf("expected response to not be ok")
	}

	expected := []recurly.Error{
		{
			XMLName: xml.Name{Local: "error"},
			Message: "is not a number",
			Field:   "model_name.field_name",
			Symbol:  "not_a_number",
		},
		{
			XMLName: xml.Name{Local: "error"},
			Message: "is not good",
			Field:   "foo.bar",
			Symbol:  "not_good",
		},
	}

	if !reflect.DeepEqual(expected, resp.Errors) {
		t.Fatalf("unexpected error: %v", resp.Errors)
	}
}

func TestClient_Unmarshal(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/account/1", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?>
<account href="https://your-subdomain.recurly.com/v2/accounts/1">
  <adjustments href="https://your-subdomain.recurly.com/v2/accounts/1/adjustments"/>
  <billing_info href="https://your-subdomain.recurly.com/v2/accounts/1/billing_info"/>
  <invoices href="https://your-subdomain.recurly.com/v2/accounts/1/invoices"/>
  <redemption href="https://your-subdomain.recurly.com/v2/accounts/1/redemption"/>
  <subscriptions href="https://your-subdomain.recurly.com/v2/accounts/1/subscriptions"/>
  <transactions href="https://your-subdomain.recurly.com/v2/accounts/1/transactions"/>
  <account_code>1</account_code>
  <state>active</state>
  <username nil="nil"></username>
  <email>verena@example.com</email>
  <first_name>Verena</first_name>
  <last_name>Example</last_name>
  <company_name></company_name>
  <vat_number nil="nil"></vat_number>
  <tax_exempt type="boolean">false</tax_exempt>
  <address>
    <address1>123 Main St.</address1>
    <address2 nil="nil"></address2>
    <city>San Francisco</city>
    <state>CA</state>
    <zip>94105</zip>
    <country>US</country>
    <phone nil="nil"></phone>
  </address>
  <accept_language nil="nil"></accept_language>
  <hosted_login_token>a92468579e9c4231a6c0031c4716c01d</hosted_login_token>
  <created_at type="datetime">2011-10-25T12:00:00Z</created_at>
</account>`)
	})

	req, err := http.NewRequest("GET", client.BaseURL+"account/1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var a recurly.Account
	resp, err := client.Do(req, &a)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatalf("expected response to be ok")
	}

	if a.FirstName != "Verena" {
		t.Fatalf("unexpected first name: %s", a.FirstName)
	} else if a.LastName != "Example" {
		t.Fatalf("unexpected last name: %s", a.LastName)
	} else if a.HostedLoginToken != "a92468579e9c4231a6c0031c4716c01d" {
		t.Fatalf("unexpected hosted login token: %s", a.HostedLoginToken)
	} else if a.Address.Address != "123 Main St." {
		t.Fatalf("unexpected address: %s", a.Address.Address)
	}
}
