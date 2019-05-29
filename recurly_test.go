package recurly_test

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/blacklightcms/recurly"
	"github.com/google/go-cmp/cmp"
)

var dialer = &net.Dialer{
	Timeout:   100 * time.Millisecond,
	KeepAlive: 500 * time.Millisecond,
	DualStack: true,
}

// Server is a test server used for testing.
type Server struct {
	server *httptest.Server
	mux    *http.ServeMux

	Invoked bool
}

// NewServer returns an instance of *recurly.Client and *Server, with the
// client resolving to the test server.
func NewServer() (*recurly.Client, *Server) {
	s := &Server{
		mux: http.NewServeMux(),
	}
	s.server = httptest.NewTLSServer(s.mux)

	client := recurly.NewClient("test", "abc")
	client.Client = &http.Client{
		Timeout: 100 * time.Millisecond,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				u, _ := url.Parse(s.server.URL)
				return dialer.DialContext(ctx, network, u.Host)
			},
		},
	}
	return client, s
}

// HandleFunc sets up an HTTP handler where the HTTP method is asserted.
// fn is the handler, and all invocation will set s.Invoked to true.
func (s *Server) HandleFunc(method string, pattern string, fn http.HandlerFunc, t *testing.T) {
	s.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		s.Invoked = true
		if r.Method != method {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		fn(w, r)
	})
}

// Close closes the server.
func (s *Server) Close() {
	s.server.Close()
}

// MustOpenFile opens a file in the testdata directory.
func MustOpenFile(file string) []byte {
	b, err := ioutil.ReadFile("testdata/" + file)
	if err != nil {
		panic(fmt.Sprintf("error reading file %q: %#v", "testdata/"+file, err))
	}
	if strings.HasSuffix(file, ".xml") {
		return MustCompact(b)
	}
	return b
}

var rxStripXMLTags = regexp.MustCompile(`>\s+<`)

// Removes all spaces between XML tags.
func MustCompact(b []byte) []byte {
	return bytes.TrimSpace(rxStripXMLTags.ReplaceAll(b, []byte("><")))
}

// Removes all spaces between XML tags.
func MustCompactString(str string) string {
	return strings.TrimSpace(rxStripXMLTags.ReplaceAllString(str, "><"))
}

// Removes the opening <?xml ...?> tag.
func MustStripXMLTag(b []byte) []byte {
	return bytes.Replace(b, []byte(`<?xml version="1.0" encoding="UTF-8"?>`), []byte(""), 1)
}

// Opens an xml file from testdata directory, but removes all spaces between
// tags and the opening <?xml ...?> tag for simple marshaler comparisons.
func MustOpenCompactXMLFile(file string) []byte {
	return MustStripXMLTag(MustCompact(MustOpenFile(file)))
}

// MustParseTime parses a string into time.Time, panicing if there is an error.
func MustParseTime(str string) time.Time {
	t, err := time.Parse(recurly.DateTimeFormat, str)
	if err != nil {
		panic(err)
	}
	return t
}

// MustReadAll reads everything from r.
func MustReadAll(r io.ReadCloser) []byte {
	defer r.Close()

	b, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return b
}

// MustReadAllString reads everything from r.
func MustReadAllString(r io.ReadCloser) string {
	return string(MustReadAll(r))
}

// TestClient tests that requests are properly structured and all of the
// expected data is sent correctly to Recurly (e.g. api key, body, etc).
func TestClient(t *testing.T) {
	// Tests a GET method.
	t.Run("GET", func(t *testing.T) {
		client, s := NewServer()
		defer s.Close()

		timestamp := MustParseTime("2011-10-17T17:24:53Z")
		s.HandleFunc("GET", "/v2/accounts", func(w http.ResponseWriter, r *http.Request) {
			// API key should be base64 encoded.
			encoded := base64.StdEncoding.EncodeToString([]byte("abc"))

			if r.Host != "test.recurly.com" {
				t.Fatalf("unexpected host: %q", r.Host)
			} else if h := r.Header.Get("Authorization"); h != fmt.Sprintf("Basic %s", encoded) {
				t.Fatalf("unexpected Authorization header: %q", h)
			} else if h := r.Header.Get("Accept"); h != "application/xml" {
				t.Fatalf("unexpected Accept header: %q", h)
			} else if h := r.Header.Get("Content-Type"); h != "" {
				t.Fatalf("unexpected Content-Type: %q", h)
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(MustOpenFile("accounts.xml")))
		}, t)

		pager := client.Accounts.List(&recurly.PagerOptions{
			State:     "active",
			Sort:      "created_at",
			Order:     "asc",
			PerPage:   25,
			BeginTime: recurly.NewTime(timestamp),
			EndTime:   recurly.NewTime(timestamp),
		})

		pager.Next()
		var a []recurly.Account
		if err := pager.Fetch(context.Background(), &a); err != nil {
			t.Fatal(err)
		}
	})

	// Tests a POST method to ensure the request body is sent.
	t.Run("POST", func(t *testing.T) {
		client, s := NewServer()
		defer s.Close()

		s.HandleFunc("POST", "/v2/accounts", func(w http.ResponseWriter, r *http.Request) {
			if b := MustReadAll(r.Body); !bytes.Equal(b, []byte(`<account><account_code>foo</account_code></account>`)) {
				t.Fatal(string(b))
			}
			w.WriteHeader(http.StatusCreated)
			w.Write(MustOpenFile("account.xml"))
		}, t)

		if a, err := client.Accounts.Create(context.Background(), recurly.Account{Code: "foo"}); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(a, NewTestAccount()); diff != "" {
			t.Fatal(diff)
		}
	})
}

// Ensure that client errors are handled.
func TestClient_ClientErrors(t *testing.T) {
	// 404 Not Found should return a ClientError.
	t.Run("404", func(t *testing.T) {
		// 404 response with validation errors
		t.Run("OK", func(t *testing.T) {
			client, s := NewServer()
			defer s.Close()

			s.HandleFunc("POST", "/v2/accounts", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write(MustOpenFile("error_not_found.xml"))
			}, t)

			_, err := client.Accounts.Create(context.Background(), recurly.Account{})
			if !s.Invoked {
				t.Fatal("expected invocation")
			} else if err == nil {
				t.Fatal(err)
			} else if e, ok := err.(*recurly.ClientError); !ok {
				t.Fatalf("unexpected error: %T %#v", err, err)
			} else if e.Response == nil {
				t.Fatal("expected *http.Response")
			} else if e.Response.StatusCode != http.StatusNotFound {
				t.Fatalf("unexpected status code: %d", e.Response.StatusCode)
			} else if diff := cmp.Diff(e.ValidationErrors, []recurly.ValidationError{{
				Symbol:      "not_found",
				Description: "The record could not be located.",
			}}); diff != "" {
				t.Fatal(diff)
			}
		})

		// 404 response with empty body
		t.Run("EmptyBody", func(t *testing.T) {
			client, s := NewServer()
			defer s.Close()

			s.HandleFunc("POST", "/v2/accounts", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			}, t)

			_, err := client.Accounts.Create(context.Background(), recurly.Account{})
			if !s.Invoked {
				t.Fatal("expected invocation")
			} else if err == nil {
				t.Fatal(err)
			} else if e, ok := err.(*recurly.ClientError); !ok {
				t.Fatalf("unexpected error: %T %#v", err, err)
			} else if e.Response == nil {
				t.Fatal("expected *http.Response")
			} else if e.Response.StatusCode != http.StatusNotFound {
				t.Fatalf("unexpected status code: %d", e.Response.StatusCode)
			} else if len(e.ValidationErrors) > 0 {
				t.Fatalf("unexpected validation errors: %#v", e.ValidationErrors)
			}
		})
	})

	t.Run("422", func(t *testing.T) {
		// Ensure a top-level <error> tag is properly handled.
		t.Run("SingleError", func(t *testing.T) {
			client, s := NewServer()
			defer s.Close()

			s.HandleFunc("POST", "/v2/accounts", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte(`
					<?xml version="1.0" encoding="UTF-8"?>
					<error>
						<symbol>simultaneous_request</symbol>
						<description>A change for subscription 3cf89f0c3fcda0b15c50134f63856d4e is already in progress.</description>
					</error>
				`))
			}, t)

			_, err := client.Accounts.Create(context.Background(), recurly.Account{})
			if !s.Invoked {
				t.Fatal("expected invocation")
			} else if err == nil {
				t.Fatal(err)
			} else if e, ok := err.(*recurly.ClientError); !ok {
				t.Fatalf("unexpected error: %T %#v", err, err)
			} else if e.Response == nil {
				t.Fatal("expected *http.Response")
			} else if e.Response.StatusCode != http.StatusUnprocessableEntity {
				t.Fatalf("unexpected status code: %d", e.Response.StatusCode)
			} else if diff := cmp.Diff(e.ValidationErrors, []recurly.ValidationError{{
				Symbol:      "simultaneous_request",
				Description: "A change for subscription 3cf89f0c3fcda0b15c50134f63856d4e is already in progress.",
			}}); diff != "" {
				t.Fatal(diff)
			}
		})

		// Ensure a top-level <errors> tag is properly handled.
		t.Run("MultiErrors", func(t *testing.T) {
			client, s := NewServer()
			defer s.Close()

			s.HandleFunc("POST", "/v2/accounts", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte(`
					<?xml version="1.0" encoding="UTF-8"?>
					<errors>
						<error field="model_name.field_name" symbol="not_a_number" lang="en-US">is not a number</error>
						<error field="subscription.base" symbol="already_subscribed">You already have a subscription to this plan.</error>
					</errors>
				`))
			}, t)

			_, err := client.Accounts.Create(context.Background(), recurly.Account{})
			if !s.Invoked {
				t.Fatal("expected invocation")
			} else if err == nil {
				t.Fatal(err)
			} else if e, ok := err.(*recurly.ClientError); !ok {
				t.Fatalf("unexpected error: %T %#v", err, err)
			} else if e.Response == nil {
				t.Fatal("expected *http.Response")
			} else if e.Response.StatusCode != http.StatusUnprocessableEntity {
				t.Fatalf("unexpected status code: %d", e.Response.StatusCode)
			} else if diff := cmp.Diff(e.ValidationErrors, []recurly.ValidationError{
				{
					Field:       "model_name.field_name",
					Symbol:      "not_a_number",
					Description: "is not a number",
				},
				{
					Field:       "subscription.base",
					Symbol:      "already_subscribed",
					Description: "You already have a subscription to this plan.",
				},
			}); diff != "" {
				t.Fatal(diff)
			}
		})
	})

	t.Run("Is", func(t *testing.T) {
		t.Run("SingleError", func(t *testing.T) {
			if err := (&recurly.ClientError{
				ValidationErrors: []recurly.ValidationError{{
					Symbol:      "number_of_unique_codes",
					Description: "You are limited to generating 200 at a time",
				}},
			}); !err.Is("number_of_unique_codes") {
				t.Fatal("expected true")
			} else if err.Is("not_found") {
				t.Fatal("expected false")
			}
		})

		t.Run("MultiErrors", func(t *testing.T) {
			if err := (&recurly.ClientError{
				ValidationErrors: []recurly.ValidationError{
					{
						Symbol:      "number_of_unique_codes",
						Description: "You are limited to generating 200 at a time",
					},
					{
						Symbol:      "will_not_invoice",
						Description: "No adjustments to invoice",
					},
				},
			}); !err.Is("number_of_unique_codes") {
				t.Fatal("expected true")
			} else if !err.Is("will_not_invoice") {
				t.Fatal("expected true")
			} else if err.Is("not_found") {
				t.Fatal("expected false")
			}
		})
	})
}

// Ensure transaction errors return TransactionFailedError.
func TestClient_TransactionFailedError(t *testing.T) {
	client, s := NewServer()
	defer s.Close()

	s.HandleFunc("POST", "/v2/accounts", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(MustOpenFile("errors_transaction_failed.xml"))
	}, t)

	_, err := client.Accounts.Create(context.Background(), recurly.Account{})
	if !s.Invoked {
		t.Fatal("expected invocation")
	} else if err == nil {
		t.Fatal(err)
	} else if e, ok := err.(*recurly.TransactionFailedError); !ok {
		t.Fatalf("unexpected error: %T %#v", err, err)
	} else if e.Response == nil {
		t.Fatal("expected *http.Response")
	} else if e.Response.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("unexpected status code: %d", e.Response.StatusCode)
	} else if diff := cmp.Diff(e.TransactionError, recurly.TransactionError{
		XMLName:          xml.Name{Local: "transaction_error"},
		ErrorCode:        "fraud_security_code",
		ErrorCategory:    "fraud",
		MerchantMessage:  "The payment gateway declined the transaction because the security code (CVV) did not match.",
		CustomerMessage:  "The security code you entered does not match. Please update the CVV and try again.",
		GatewayErrorCode: "301",
	}); diff != "" {
		t.Fatal(diff)
	} else if diff := cmp.Diff(e.Transaction, NewTestTransactionFailed()); diff != "" {
		t.Fatal(diff)
	}
}

func TestClient_ServerErrors(t *testing.T) {
	client, s := NewServer()
	defer s.Close()

	s.HandleFunc("POST", "/v2/accounts", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}, t)

	_, err := client.Accounts.Create(context.Background(), recurly.Account{})
	if !s.Invoked {
		t.Fatal("expected invocation")
	} else if err == nil {
		t.Fatal(err)
	} else if e, ok := err.(*recurly.ServerError); !ok {
		t.Fatalf("unexpected error: %T %#v", err, err)
	} else if e.Response == nil {
		t.Fatal("expected *http.Response")
	} else if e.Response.StatusCode != http.StatusInternalServerError {
		t.Fatalf("unexpected status code: %d", e.Response.StatusCode)
	}
}
