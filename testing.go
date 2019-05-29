package recurly

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

// TestServer is a server used for testing when mocks are not sufficient.
// This enables pointing a recurly client to this test server and simulating
// requests/responses directly.
type TestServer struct {
	server *httptest.Server
	mux    *http.ServeMux

	Invoked bool
}

// NewTestServer returns an instance of *Client and *Server, with the
// client resolving to the test server.
func NewTestServer() (*Client, *TestServer) {
	s := &TestServer{
		mux: http.NewServeMux(),
	}
	s.server = httptest.NewTLSServer(s.mux)

	client := NewClient("test", "foo")
	client.Client = &http.Client{
		Timeout: 100 * time.Millisecond,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				u, _ := url.Parse(s.server.URL)
				return (&net.Dialer{
					Timeout:   100 * time.Millisecond,
					KeepAlive: 500 * time.Millisecond,
					DualStack: true,
				}).DialContext(ctx, network, u.Host)
			},
		},
	}
	return client, s
}

// HandleFunc sets up an HTTP handler where the HTTP method is asserted.
// fn is the handler, and all invocation will set s.Invoked to true.
func (s *TestServer) HandleFunc(method string, pattern string, fn http.HandlerFunc, t *testing.T) {
	s.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		s.Invoked = true
		if r.Method != method {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		fn(w, r)
	})
}

// Close closes the server.
func (s *TestServer) Close() {
	s.server.Close()
}
