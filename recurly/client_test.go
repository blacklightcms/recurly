package recurly

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	expected := &Client{
		client:    http.DefaultClient,
		subDomain: "foo",
		apiKey:    "bar",
		BaseURL:   "https://foo.recurly.com/",
	}

	given := NewClient("foo", "bar", nil)
	if expected.subDomain != given.subDomain {
		t.Errorf("TestNewClient Error: Expected subDomain of %s, given %s", expected.subDomain, given.subDomain)
	}

	if expected.apiKey != given.apiKey {
		t.Errorf("TestNewClient Error: Expected apiKey of %s, given %s", expected.apiKey, given.apiKey)
	}

	if expected.BaseURL != given.BaseURL {
		t.Errorf("TestNewClient Error: Expected BaseURL of %s, given %s", expected.BaseURL, given.BaseURL)
	}
}

func TestNewRequest(t *testing.T) {
	client = NewClient("test", "abc", nil)

	req, err := client.newRequest("GET", "accounts/14579", Params{"foo": "bar"}, nil)
	if err != nil {
		t.Errorf("TestNewRequest Error: %v", err)
	}

	if req.URL.Path != "/v2/accounts/14579" {
		t.Errorf("TestNewRequest: expected path to equal %v, actual %v", "/v2/accounts/14579", req.URL.Path)
	}

	if req.Method != "GET" {
		t.Errorf("TestNewRequest: expected method to equal %s, actual %s", "GET", req.Method)
	}

	if req.Header.Get("Accept") != "application/xml" {
		t.Errorf("TestNewRequest: expected accept header of %s, given %s", "application/xml", req.Header.Get("Accept"))
	}

	if req.Header.Get("Content-Type") != "" {
		t.Errorf("TestNewRequest: expected empty content-type header, given %s", req.Header.Get("Content-Type"))
	}

	query := req.URL.Query()
	for name, expected := range map[string]string{"foo": "bar"} {
		actual := query.Get(name)
		if actual != expected {
			t.Errorf("TestNewRequest: expected '%s' to equal '%s', actual '%s'", name, expected, actual)
		}
	}

	req, err = client.newRequest("PUT", "accounts/abc", nil, Account{Code: "abc"})
	if err != nil {
		t.Errorf("TestNewRequest Error: %v", err)
	}

	if req.URL.Path != "/v2/accounts/abc" {
		t.Errorf("TestNewRequest: expected path to equal %v, actual %v", "/v2/accounts/abc", req.URL.Path)
	}

	if req.Method != "PUT" {
		t.Errorf("TestNewRequest: expected method to equal %s, actual %s", "PUT", req.Method)
	}

	if req.Header.Get("Accept") != "application/xml" {
		t.Errorf("TestNewRequest: expected accept header of %s, given %s", "application/xml", req.Header.Get("Accept"))
	}

	if req.Header.Get("Content-Type") != "application/xml; charset=utf-8" {
		t.Errorf("TestNewRequest: expected content-type header of %s, given %s", "application/xml; charset=utf-8", req.Header.Get("Content-Type"))
	}

	expected := []byte("<account><account_code>abc</account_code></account>")
	given, _ := ioutil.ReadAll(req.Body)
	if !reflect.DeepEqual(expected, given) {
		t.Errorf("TestNewRequest: expected string body equal to %s, given %s", expected, given)
	}

	query = req.URL.Query()
	if len(query) != 0 {
		t.Errorf("TestNewRequest: expected %d query params, given %d", 0, len(query))
	}
}

func TestRequestThatReturnsErrorMessages(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/error", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(422)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?>
			<errors>
				<error field="model_name.field_name" symbol="not_a_number" lang="en-US">is not a number</error>
				<error field="foo.bar" symbol="not_good" lang="en-US">is not good</error>
			</errors>`)
	})

	req, err := http.NewRequest("GET", client.BaseURL+"error", nil)
	if err != nil {
		t.Fatalf("TestRequestThatReturnsErrorMessages Error: Error creating request. err: %s", err)
	}

	resp, err := client.do(req, nil)
	if err != nil {
		t.Fatalf("TestRequestThatReturnsErrorMessages Error: Error making request. err: %s", err)
	}

	if resp.IsOK() {
		t.Errorf("TestRequestThatReturnsErrorMessages Error: Expected response to not be ok")
	}

	expected := []Error{
		Error{
			XMLName: xml.Name{Local: "error"},
			Message: "is not a number",
			Field:   "model_name.field_name",
			Symbol:  "not_a_number",
		},
		Error{
			XMLName: xml.Name{Local: "error"},
			Message: "is not good",
			Field:   "foo.bar",
			Symbol:  "not_good",
		},
	}

	if !reflect.DeepEqual(expected, resp.Errors) {
		t.Errorf("TestRequestThatReturnsErrorMessages Error: Expected errors of %#v, given %#v", expected, resp.Errors)
	}
}

func TestRequestUnmarshalsIntoStruct(t *testing.T) {
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
		t.Fatalf("TestRequestUnmarshalsIntoStruct Error: Error creating request. err: %s", err)
	}

	var a Account
	resp, err := client.do(req, &a)
	if err != nil {
		t.Fatalf("TestRequestUnmarshalsIntoStruct Error: Error making request. err: %s", err)
	}

	if resp.IsError() {
		t.Errorf("TestRequestUnmarshalsIntoStruct Error: Expected response to be ok")
	}

	expected := "Verena"
	if expected != a.FirstName {
		t.Errorf("TestRequestUnmarshalsIntoStruct Error: Expected first name to be %s, given %s", expected, a.FirstName)
	}

	expected = "Example"
	if expected != a.LastName {
		t.Errorf("TestRequestUnmarshalsIntoStruct Error: Expected first name to be %s, given %s", expected, a.LastName)
	}

	expected = "a92468579e9c4231a6c0031c4716c01d"
	if expected != a.HostedLoginToken {
		t.Errorf("TestRequestUnmarshalsIntoStruct Error: Expected first name to be %s, given %s", expected, a.HostedLoginToken)
	}

	expected = "123 Main St."
	if expected != a.Address.Address {
		t.Errorf("TestRequestUnmarshalsIntoStruct Error: Expected address1 to be %s, given %s", expected, a.Address.Address)
	}
}
