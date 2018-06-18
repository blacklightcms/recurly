package recurly_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/blacklightcms/recurly"
	"github.com/google/go-cmp/cmp"
)

func TestShippingAddress_ListAccount(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/1/shipping_addresses", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
			<shipping_addresses type="array">
			  <shipping_address href="https://your-subdomain.recurly.com/v2/accounts/1/shipping_addresses/2438622711411416831">
			    <account href="https://your-subdomain.recurly.com/v2/accounts/1"/>
			    <subscriptions href="https://your-subdomain.recurly.com/v2/accounts/1/shipping_addresses/2438622711411416831/subscriptions"/>
			    <id type="integer">2438622711411416831</id>
			    <nickname>Work</nickname>
			    <first_name>Verena</first_name>
			    <last_name>Example</last_name>
			    <company>Recurly Inc</company>
			    <email>verena@example.com</email>
			    <vat_number nil="nil"/>
			    <address1>123 Main St.</address1>
			    <address2>Suite 101</address2>
			    <city>San Francisco</city>
			    <state>CA</state>
			    <zip>94105</zip>
			    <country>US</country>
			    <phone>555-222-1212</phone>
			    <created_at type="datetime">2018-03-19T15:48:00Z</created_at>
			    <updated_at type="datetime">2018-03-19T15:48:00Z</updated_at>
			  </shipping_address>
			  <shipping_address href="https://your-subdomain.recurly.com/v2/accounts/1/shipping_addresses/2">
			    <account href="https://your-subdomain.recurly.com/v2/accounts/1"/>
			    <shipping_address_id>2</shipping_address_id>
			    <nickname>Home</nickname>
			    <first_name>Verena</first_name>
			    <last_name>Example</last_name>
			    <phone>555-867-5309</phone>
			    <email>verena@example.com</email>
			    <address1>123 Fourth St.</address1>
			    <address2>Apt. 101</address2>
			    <city>San Francisco</city>
			    <state>CA</state>
			    <zip>94105</zip>
			    <country>US</country>
			    <created_at type="datetime">2018-03-19T15:48:00Z</created_at>
			    <updated_at type="datetime">2018-03-19T15:48:00Z</updated_at>
			  </shipping_address>
			</shipping_addresses>`)
	})
	r, shippingAddresses, err := client.ShippingAddresses.ListAccount("1", recurly.Params{"per_page": 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected list subcriptions to return OK")
	} else if pp := r.Request.URL.Query().Get("per_page"); pp != "1" {
		t.Fatalf("unexpected per_page: %s", pp)
	}

	created, _ := time.Parse(recurly.DateTimeFormat, "2018-03-19T15:48:00Z")
	updated, _ := time.Parse(recurly.DateTimeFormat, "2018-03-19T15:48:00Z")

	if diff := cmp.Diff(shippingAddresses, []recurly.ShippingAddress{
		{
			AddressID:     "2438622711411416831",
			AccountCode:   "1",
			Subscriptions: "subscriptions",
			FirstName:     "Verena",
			LastName:      "Example",
			Company:       "Recurly Inc",
			Email:         "verena@example.com",
			Address:       "123 Main St.",
			Address2:      "Suite 101",
			City:          "San Francisco",
			State:         "CA",
			Zip:           "94105",
			Country:       "US",
			Phone:         "555-222-1212",
			CreatedAt:     recurly.NewTime(created),
			UpdatedAt:     recurly.NewTime(updated),
		},
		{
			AccountCode: "1",
			FirstName:   "Verena",
			LastName:    "Example",
			Email:       "verena@example.com",
			Address:     "123 Fourth St.",
			Address2:    "Apt. 101",
			City:        "San Francisco",
			State:       "CA",
			Zip:         "94105",
			Country:     "US",
			Phone:       "555-867-5309",
			CreatedAt:   recurly.NewTime(created),
			UpdatedAt:   recurly.NewTime(updated),
		},
	}); diff != "" {
		t.Fatal(diff)
	}

}

func TestShippingAddress_Create(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/v2/accounts/1/shipping_addresses", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(201)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><shipping_address href="https://your-subdomain.recurly.com/v2/accounts/1/shipping_addresses/2438622711411416831"></shipping_address>`)
	})

	r, _, err := client.ShippingAddresses.Create("1", recurly.ShippingAddress{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected create subscription to return OK")
	}
}

func TestShippingAddress_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/1/shipping_addresses/2438622711411416831", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(201)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><shipping_address href="https://your-subdomain.recurly.com/v2/accounts/1/shipping_addresses/2438622711411416831"></shipping_address>`)
	})

	r, _, err := client.ShippingAddresses.Update("1", "2438622711411416831", recurly.ShippingAddress{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected update shipping address to return 201 Created.")
	}
}

func TestShippingAddressDelete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/1/shipping_addresses/2438622711411416831", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(204)
	})

	r, err := client.ShippingAddresses.Delete("1", "2438622711411416831")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected delete shipping address to return 204 No content.")
	}
}
