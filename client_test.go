package recurly

import (
	"net/http"
	"testing"
)

func TestClient(t *testing.T) {
	expected := &Client{
		Client:    http.DefaultClient,
		subDomain: "foo",
		apiKey:    "bar",
		BaseURL:   "https://foo.com/",
	}

	given := NewClient("foo", "bar", nil)
	if expected.subDomain != given.subDomain {
		t.Fatalf("unexpected subdomain: %s", given.subDomain)
	} else if expected.apiKey != given.apiKey {
		t.Fatalf("unexpected api key: %s", given.apiKey)
	} else if expected.BaseURL != given.BaseURL {
		t.Fatalf("unexpected base url: %s", given.BaseURL)
	}
}
