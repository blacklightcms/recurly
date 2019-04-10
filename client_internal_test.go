package recurly

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// TestClient_NewRequest tests the internals of recurly.client.
func TestClient_NewRequest(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	client := NewClient("test", "abc", nil)
	client.BaseURL = server.URL + "/"
	defer server.Close()

	// API key should be base64 encoded.
	encoded := base64.StdEncoding.EncodeToString([]byte("abc"))

	req, err := client.newRequest("GET", "accounts/14579", Params{"foo": "bar"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if req.URL.Path != "/v2/accounts/14579" {
		t.Fatalf("unexpected path: %s", req.URL.Path)
	} else if req.Method != "GET" {
		t.Fatalf("unexpected method: %s", req.Method)
	} else if req.Header.Get("Authorization") != fmt.Sprintf("Basic %s", encoded) { // API key should be base64 encoded
		t.Fatalf("unexpected Authorization header: %s", req.Header.Get("Authorization"))
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

	req, err = client.newRequest("PUT", "accounts/abc", nil, Account{Code: "abc"})
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
	if diff := cmp.Diff(expected, given); diff != "" {
		t.Fatal(diff)
	}

	query = req.URL.Query()
	if len(query) != 0 {
		t.Fatalf("expected %d query Params, given %d", 0, len(query))
	}
}

// TestClient_Errors tests the internals of recurly.client returning a 422
// repsonse with an array of errors.
func TestClient_Errors(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	client := NewClient("test", "abc", nil)
	client.BaseURL = server.URL + "/"
	defer server.Close()

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

	resp, err := client.do(req, nil)
	if err != nil {
		t.Fatalf("error making request. err: %v", err)
	} else if resp.IsOK() {
		t.Fatalf("expected response to not be ok")
	}

	// Transaction should be nil
	if resp.Transaction != nil {
		t.Fatal("expected transaction to be nil")
	}

	expected := []Error{
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

	if diff := cmp.Diff(expected, resp.Errors); diff != "" {
		t.Fatal(diff)
	}
}

// TestClient_Error tests the internals of recurly.client with a 422
// response with a single error.
func TestClient_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	client := NewClient("test", "abc", nil)
	client.BaseURL = server.URL + "/"
	defer server.Close()

	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(422)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
			<error>
				<symbol>simultaneous_request</symbol>
				<description>A change for subscription 3cf89f0c3fcda0b15c50134f63856d4e is already in progress.</description>
			</error>`)
	})

	req, err := http.NewRequest("GET", client.BaseURL+"error", nil)
	if err != nil {
		t.Fatalf("error creating request. err: %v", err)
	}

	resp, err := client.do(req, nil)
	if err != nil {
		t.Fatalf("error making request. err: %v", err)
	} else if resp.IsOK() {
		t.Fatalf("expected response to not be ok")
	}

	// Transaction should be nil
	if resp.Transaction != nil {
		t.Fatal("expected transaction to be nil")
	}

	expected := []Error{
		{
			XMLName:     xml.Name{Local: "error"},
			Symbol:      "simultaneous_request",
			Description: "A change for subscription 3cf89f0c3fcda0b15c50134f63856d4e is already in progress.",
		},
	}

	if diff := cmp.Diff(expected, resp.Errors); diff != "" {
		t.Fatal(diff)
	}
}
