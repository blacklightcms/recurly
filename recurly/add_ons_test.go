package recurly

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

// TestAddOnEncoding ensures structs are encoded to XML properly.
// Because Recurly supports partial updates, it's important that only defined
// fields are handled properly -- including types like booleans and integers which
// have zero values that we want to send.
func TestAddOns_Encoding(t *testing.T) {
	tests := []struct {
		v        AddOn
		expected string
	}{
		{v: AddOn{}, expected: "<add_on></add_on>"},
		{v: AddOn{Code: "xyz"}, expected: "<add_on><add_on_code>xyz</add_on_code></add_on>"},
		{v: AddOn{Name: "IP Addresses"}, expected: "<add_on><name>IP Addresses</name></add_on>"},
		{v: AddOn{DefaultQuantity: NewInt(0)}, expected: "<add_on><default_quantity>0</default_quantity></add_on>"},
		{v: AddOn{DefaultQuantity: NewInt(1)}, expected: "<add_on><default_quantity>1</default_quantity></add_on>"},
		{v: AddOn{DisplayQuantityOnHostedPage: NewBool(true)}, expected: "<add_on><display_quantity_on_hosted_page>true</display_quantity_on_hosted_page></add_on>"},
		{v: AddOn{DisplayQuantityOnHostedPage: NewBool(false)}, expected: "<add_on><display_quantity_on_hosted_page>false</display_quantity_on_hosted_page></add_on>"},
		{v: AddOn{TaxCode: "digital"}, expected: "<add_on><tax_code>digital</tax_code></add_on>"},
		{v: AddOn{UnitAmountInCents: UnitAmount{USD: 200}}, expected: "<add_on><unit_amount_in_cents><USD>200</USD></unit_amount_in_cents></add_on>"},
		{v: AddOn{AccountingCode: "abc123"}, expected: "<add_on><accounting_code>abc123</accounting_code></add_on>"},
	}

	for _, tt := range tests {
		var buf bytes.Buffer
		err := xml.NewEncoder(&buf).Encode(tt.v)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		} else if buf.String() != tt.expected {
			t.Fatalf("unexpected value: %s", buf.String())
		}
	}
}

func TestAddOns_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/plans/gold/add_ons", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
		<add_ons type="array">
			<add_on href="https://your-subdomain.recurly.com/v2/plans/gold/add_ons/ipaddresses">
				<plan href="https://your-subdomain.recurly.com/v2/plans/gold"/>
				<add_on_code>ipaddresses</add_on_code>
				<name>IP Addresses</name>
				<default_quantity type="integer">1</default_quantity>
				<display_quantity_on_hosted_page type="boolean">false</display_quantity_on_hosted_page>
				<tax_code>digital</tax_code>
				<unit_amount_in_cents>
					<USD type="integer">200</USD>
				</unit_amount_in_cents>
				<accounting_code>abc123</accounting_code>
				<created_at type="datetime">2011-06-28T12:34:56Z</created_at>
			</add_on>
		</add_ons>`)
	})

	r, addOns, err := client.AddOns.List("gold", Params{"per_page": 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected list add ons to return OK")
	} else if len(addOns) != 1 {
		t.Fatalf("unexpected length: %d", len(addOns))
	} else if pp := r.Request.URL.Query().Get("per_page"); pp != "1" {
		t.Fatalf("unexpected per_page: %s", pp)
	}

	ts, _ := time.Parse(datetimeFormat, "2011-06-28T12:34:56Z")
	for _, given := range addOns {
		expected := AddOn{
			XMLName:                     xml.Name{Local: "add_on"},
			Code:                        "ipaddresses",
			Name:                        "IP Addresses",
			DefaultQuantity:             NewInt(1),
			DisplayQuantityOnHostedPage: NewBool(false),
			TaxCode:                     "digital",
			UnitAmountInCents:           UnitAmount{USD: 200},
			AccountingCode:              "abc123",
			CreatedAt:                   NewTime(ts),
		}

		if !reflect.DeepEqual(expected, given) {
			t.Fatalf("unexpected add on: %v", given)
		}
	}
}

func TestAddOns_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/plans/gold/add_ons/ipaddresses", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
			<add_on href="https://your-subdomain.recurly.com/v2/plans/gold/add_ons/ipaddresses">
				<plan href="https://your-subdomain.recurly.com/v2/plans/gold"/>
				<add_on_code>ipaddresses</add_on_code>
				<name>IP Addresses</name>
				<default_quantity type="integer">1</default_quantity>
				<display_quantity_on_hosted_page type="boolean">false</display_quantity_on_hosted_page>
				<tax_code>digital</tax_code>
				<unit_amount_in_cents>
					<USD type="integer">200</USD>
				</unit_amount_in_cents>
				<accounting_code>abc123</accounting_code>
				<created_at type="datetime">2011-06-28T12:34:56Z</created_at>
			</add_on>`)
	})

	r, a, err := client.AddOns.Get("gold", "ipaddresses")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected get add_on to return OK")
	}

	ts, _ := time.Parse(datetimeFormat, "2011-06-28T12:34:56Z")
	expected := AddOn{
		XMLName:                     xml.Name{Local: "add_on"},
		Code:                        "ipaddresses",
		Name:                        "IP Addresses",
		DefaultQuantity:             NewInt(1),
		DisplayQuantityOnHostedPage: NewBool(false),
		TaxCode:                     "digital",
		UnitAmountInCents:           UnitAmount{USD: 200},
		AccountingCode:              "abc123",
		CreatedAt:                   NewTime(ts),
	}

	if !reflect.DeepEqual(expected, a) {
		t.Fatalf("unexpected add on: %v", a)
	}
}

func TestAddOns_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/plans/gold/add_ons", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(201)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><add_on></add_on>`)
	})

	r, _, err := client.AddOns.Create("gold", AddOn{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.StatusCode != 201 {
		t.Fatalf("unexpected response: %d", r.StatusCode)
	}
}

func TestAddOns_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/plans/gold/add_ons/ipaddresses", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><add_on></add_on>`)
	})

	r, _, err := client.AddOns.Update("gold", "ipaddresses", AddOn{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected update add on to return OK")
	}
}

func TestAddOns_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/plans/gold/add_ons/ipaddresses", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		rw.WriteHeader(204)
	})

	r, err := client.AddOns.Delete("gold", "ipaddresses")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected deleted add on to return OK")
	}
}
