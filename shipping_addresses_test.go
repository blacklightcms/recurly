package recurly_test

import (
	"context"
	"encoding/xml"
	"net/http"
	"testing"

	"github.com/blacklightcms/recurly"
	"github.com/google/go-cmp/cmp"
)

func TestShippingAddresses_ListAccount(t *testing.T) {
	client, s := NewServer()
	defer s.Close()

	var invocations int
	s.HandleFunc("GET", "/v2/accounts/1/shipping_addresses", func(w http.ResponseWriter, r *http.Request) {
		invocations++
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("shipping_addresses.xml"))
	}, t)

	pager := client.ShippingAddresses.ListAccount("1", nil)
	for pager.Next() {
		if addresses, err := pager.Fetch(context.Background()); err != nil {
			t.Fatal(err)
		} else if !s.Invoked {
			t.Fatal("expected s to be invoked")
		} else if diff := cmp.Diff(addresses, []recurly.ShippingAddress{*NewTestShippingAddress()}); diff != "" {
			t.Fatal(diff)
		}
	}
	if invocations != 1 {
		t.Fatalf("unexpected number of invocations: %d", invocations)
	}
}

func TestShippingAddresses_Create(t *testing.T) {
	client, s := NewServer()
	defer s.Close()

	s.HandleFunc("POST", "/v2/accounts/1/shipping_addresses", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write(MustOpenFile("shipping_address.xml"))
	}, t)

	if a, err := client.ShippingAddresses.Create(context.Background(), "1", recurly.ShippingAddress{}); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(a, NewTestShippingAddress()); diff != "" {
		t.Fatal(diff)
	}
}

func TestShippingAddresses_Update(t *testing.T) {
	client, s := NewServer()
	defer s.Close()

	s.HandleFunc("PUT", "/v2/accounts/1/shipping_addresses/2438622711411416831", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("shipping_address.xml"))
	}, t)

	if a, err := client.ShippingAddresses.Update(context.Background(), "1", 2438622711411416831, recurly.ShippingAddress{}); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(a, NewTestShippingAddress()); diff != "" {
		t.Fatal(diff)
	}
}

func TestShippingAddresses_Delete(t *testing.T) {
	client, s := NewServer()
	defer s.Close()

	s.HandleFunc("DELETE", "/v2/accounts/1/shipping_addresses/2438622711411416831", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}, t)

	if err := client.ShippingAddresses.Delete(context.Background(), "1", 2438622711411416831); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	}
}

// Returns a ShippingAddress corresponding to testdata/shipping_address.xml.
func NewTestShippingAddress() *recurly.ShippingAddress {
	return &recurly.ShippingAddress{
		XMLName:   xml.Name{Local: "shipping_address"},
		ID:        2438622711411416831,
		Nickname:  "Work",
		FirstName: "Verena",
		LastName:  "Example",
		Company:   "Recurly Inc",
		Email:     "verena@example.com",
		Address:   "123 Main St.",
		Address2:  "Suite 101",
		City:      "San Francisco",
		State:     "CA",
		Zip:       "94105",
		Country:   "US",
		Phone:     "555-222-1212",
		CreatedAt: recurly.NewTime(MustParseTime("2018-03-19T15:48:00Z")),
		UpdatedAt: recurly.NewTime(MustParseTime("2018-03-19T15:48:00Z")),
	}
}
