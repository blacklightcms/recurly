package recurly_test

import (
	"context"
	"encoding/xml"
	"net/http"
	"testing"

	"github.com/autopilot3/recurly"
	"github.com/google/go-cmp/cmp"
)

func TestShippingMethods_Get(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/shipping_methods/fast_fast_fast", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write(MustOpenFile("shipping_method.xml"))
		}, t)

		if method, err := client.ShippingMethods.Get(context.Background(), "fast_fast_fast"); err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(method, NewTestShippingMethod()); diff != "" {
			t.Fatal(diff)
		} else if !s.Invoked {
			t.Fatal("expected fn invocation")
		}
	})

	// Ensure a 404 returns nil values.
	t.Run("ErrNotFound", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/shipping_methods/fast_fast_fast", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}, t)

		if method, err := client.ShippingMethods.Get(context.Background(), "fast_fast_fast"); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if method != nil {
			t.Fatalf("expected nil: %#v", method)
		}
	})
}

func TestShippingMethods_List(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	var invocations int
	s.HandleFunc("GET", "/v2/shipping_methods", func(w http.ResponseWriter, r *http.Request) {
		invocations++
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("shipping_methods.xml"))
	}, t)

	pager := client.ShippingMethods.List(nil)
	for pager.Next() {
		var methods []recurly.ShippingMethod
		if err := pager.Fetch(context.Background(), &methods); err != nil {
			t.Fatal(err)
		} else if !s.Invoked {
			t.Fatal("expected s to be invoked")
		} else if diff := cmp.Diff(methods, []recurly.ShippingMethod{*NewTestShippingMethod()}); diff != "" {
			t.Fatal(diff)
		}
	}
	if invocations != 1 {
		t.Fatalf("unexpected number of invocations: %d", invocations)
	}
}

// Returns a ShippingMethod corresponding to testdata/shipping_method.xml.
func NewTestShippingMethod() *recurly.ShippingMethod {
	return &recurly.ShippingMethod{
		XMLName:   xml.Name{Local: "shipping_method"},
		Code:      "fast_fast_fast",
		Name:      "Fast Fast Fast",
		CreatedAt: recurly.NewTime(MustParseTime("2019-05-03T23:17:16Z")),
		UpdatedAt: recurly.NewTime(MustParseTime("2019-05-03T23:17:16Z")),
	}
}
