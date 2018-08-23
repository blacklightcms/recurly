package recurly_test

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/launchpadcentral/recurly"
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
			ID:          2438622711411416831,
			AccountCode: "1",
			Nickname:    "Work",
			FirstName:   "Verena",
			LastName:    "Example",
			Company:     "Recurly Inc",
			Email:       "verena@example.com",
			Address:     "123 Main St.",
			Address2:    "Suite 101",
			City:        "San Francisco",
			State:       "CA",
			Zip:         "94105",
			Country:     "US",
			Phone:       "555-222-1212",
			CreatedAt:   recurly.NewTime(created),
			UpdatedAt:   recurly.NewTime(updated),
		},
		{
			AccountCode: "1",
			FirstName:   "Verena",
			LastName:    "Example",
			Nickname:    "Home",
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
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><?xml version="1.0" encoding="UTF-8"?>
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
		  </shipping_address>`)
	})

	r, shippingAddress, err := client.ShippingAddresses.Create("1", recurly.ShippingAddress{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected create subscription to return OK")
	}

	created, _ := time.Parse(recurly.DateTimeFormat, "2018-03-19T15:48:00Z")
	updated, _ := time.Parse(recurly.DateTimeFormat, "2018-03-19T15:48:00Z")

	if diff := cmp.Diff(shippingAddress, &recurly.ShippingAddress{
		ID:          2438622711411416831,
		AccountCode: "1",
		Nickname:    "Work",
		FirstName:   "Verena",
		LastName:    "Example",
		Company:     "Recurly Inc",
		Email:       "verena@example.com",
		Address:     "123 Main St.",
		Address2:    "Suite 101",
		City:        "San Francisco",
		State:       "CA",
		Zip:         "94105",
		Country:     "US",
		Phone:       "555-222-1212",
		CreatedAt:   recurly.NewTime(created),
		UpdatedAt:   recurly.NewTime(updated),
	}); diff != "" {
		t.Fatal(diff)
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
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
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
		</shipping_address>`)
	})

	r, shippingAddress, err := client.ShippingAddresses.Update("1", 2438622711411416831, recurly.ShippingAddress{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected update shipping address to return 201 Created.")
	}
	created, _ := time.Parse(recurly.DateTimeFormat, "2018-03-19T15:48:00Z")
	updated, _ := time.Parse(recurly.DateTimeFormat, "2018-03-19T15:48:00Z")

	if diff := cmp.Diff(shippingAddress, &recurly.ShippingAddress{
		ID:          2438622711411416831,
		AccountCode: "1",
		Nickname:    "Work",
		FirstName:   "Verena",
		LastName:    "Example",
		Company:     "Recurly Inc",
		Email:       "verena@example.com",
		Address:     "123 Main St.",
		Address2:    "Suite 101",
		City:        "San Francisco",
		State:       "CA",
		Zip:         "94105",
		Country:     "US",
		Phone:       "555-222-1212",
		CreatedAt:   recurly.NewTime(created),
		UpdatedAt:   recurly.NewTime(updated),
	}); diff != "" {
		t.Fatal(diff)
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

	r, err := client.ShippingAddresses.Delete("1", 2438622711411416831)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected delete shipping address to return 204 No content.")
	}
}

func TestShippingAddress_GetSubscriptions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/1/shipping_addresses/2438622711411416831/subscriptions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
			<subscriptions type="array">
				<subscription href="https://your-subdomain.recurly.com/v2/subscriptions/44f83d7cba354d5b84812419f923ea96">
					<account href="https://your-subdomain.recurly.com/v2/accounts/1"/>
					<invoice href="https://your-subdomain.recurly.com/v2/invoices/1108"/>
					<plan href="https://your-subdomain.recurly.com/v2/plans/gold">
					  <plan_code>gold</plan_code>
					  <name>Gold plan</name>
					</plan>
					<uuid>44f83d7cba354d5b84812419f923ea96</uuid>
					<state>active</state>
					<unit_amount_in_cents type="integer">800</unit_amount_in_cents>
					<currency>EUR</currency>
					<quantity type="integer">1</quantity>
					<activated_at type="datetime">2011-05-27T07:00:00Z</activated_at>
					<canceled_at nil="nil"></canceled_at>
					<expires_at nil="nil"></expires_at>
					<current_period_started_at type="datetime">2011-06-27T07:00:00Z</current_period_started_at>
					<current_period_ends_at type="datetime">2010-07-27T07:00:00Z</current_period_ends_at>
					<trial_started_at nil="nil"></trial_started_at>
					<trial_ends_at nil="nil"></trial_ends_at>
					<tax_in_cents type="integer">72</tax_in_cents>
					<tax_type>usst</tax_type>
					<tax_region>CA</tax_region>
					<tax_rate type="float">0.0875</tax_rate>
					<po_number nil="nil"></po_number>
					<net_terms type="integer">0</net_terms>
					<subscription_add_ons type="array">
					</subscription_add_ons>
					<a name="cancel" href="https://your-subdomain.recurly.com/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/cancel" method="put"/>
					<a name="terminate" href="https://your-subdomain.recurly.com/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/terminate" method="put"/>
					<a name="postpone" href="https://your-subdomain.recurly.com/v2/subscriptions/44f83d7cba354d5b84812419f923ea96/postpone" method="put"/>
				</subscription>
			</subscriptions>`)
	})
	resp, subscriptions, err := client.ShippingAddresses.GetSubscriptions("1", 2438622711411416831)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected create subscription to return OK")
	}

	activated, _ := time.Parse(recurly.DateTimeFormat, "2011-05-27T07:00:00Z")
	cpStartedAt, _ := time.Parse(recurly.DateTimeFormat, "2011-06-27T07:00:00Z")
	cpEndsAt, _ := time.Parse(recurly.DateTimeFormat, "2010-07-27T07:00:00Z")

	if diff := cmp.Diff(subscriptions, []recurly.Subscription{
		{
			XMLName: xml.Name{Local: "subscription"},
			Plan: recurly.NestedPlan{
				Code: "gold",
				Name: "Gold plan",
			},
			AccountCode:            "1",
			InvoiceNumber:          1108,
			UUID:                   "44f83d7cba354d5b84812419f923ea96",
			State:                  "active",
			UnitAmountInCents:      800,
			Currency:               "EUR",
			Quantity:               1,
			ActivatedAt:            recurly.NewTime(activated),
			CurrentPeriodStartedAt: recurly.NewTime(cpStartedAt),
			CurrentPeriodEndsAt:    recurly.NewTime(cpEndsAt),
			TaxInCents:             72,
			TaxType:                "usst",
			TaxRegion:              "CA",
			TaxRate:                0.0875,
			NetTerms:               recurly.NewInt(0),
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}
