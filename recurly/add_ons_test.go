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
func TestAddOnsEncoding(t *testing.T) {
	suite := []map[string]interface{}{
		map[string]interface{}{"struct": AddOn{}, "xml": "<add_on></add_on>"},
		map[string]interface{}{"struct": AddOn{Code: "xyz"}, "xml": "<add_on><add_on_code>xyz</add_on_code></add_on>"},
		map[string]interface{}{"struct": AddOn{Name: "IP Addresses"}, "xml": "<add_on><name>IP Addresses</name></add_on>"},
		map[string]interface{}{"struct": AddOn{DefaultQuantity: NewInt(0)}, "xml": "<add_on><default_quantity>0</default_quantity></add_on>"},
		map[string]interface{}{"struct": AddOn{DefaultQuantity: NewInt(1)}, "xml": "<add_on><default_quantity>1</default_quantity></add_on>"},
		map[string]interface{}{"struct": AddOn{DisplayQuantityOnHostedPage: NewBool(true)}, "xml": "<add_on><display_quantity_on_hosted_page>true</display_quantity_on_hosted_page></add_on>"},
		map[string]interface{}{"struct": AddOn{DisplayQuantityOnHostedPage: NewBool(false)}, "xml": "<add_on><display_quantity_on_hosted_page>false</display_quantity_on_hosted_page></add_on>"},
		map[string]interface{}{"struct": AddOn{TaxCode: "digital"}, "xml": "<add_on><tax_code>digital</tax_code></add_on>"},
		map[string]interface{}{"struct": AddOn{UnitAmountInCents: &UnitAmount{USD: 200}}, "xml": "<add_on><unit_amount_in_cents><USD>200</USD></unit_amount_in_cents></add_on>"},
		map[string]interface{}{"struct": AddOn{AccountingCode: "abc123"}, "xml": "<add_on><accounting_code>abc123</accounting_code></add_on>"},
	}

	for _, s := range suite {
		buf := new(bytes.Buffer)
		err := xml.NewEncoder(buf).Encode(s["struct"])
		if err != nil {
			t.Errorf("TestAddOnEncoding Error: %s", err)
		}

		if buf.String() != s["xml"] {
			t.Errorf("TestAddOnEncoding Error: Expected %s, given %s", s["xml"], buf.String())
		}
	}
}

func TestAddOnsList(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/plans/gold/add_ons", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("TestAddOnsList Error: Expected %s request, given %s", "GET", r.Method)
		}
		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?>
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
		t.Errorf("TestAddOnsList Error: Error occured making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestAddOnsList Error: Expected list add ons to return OK")
	}

	if len(addOns) != 1 {
		t.Fatalf("TestAddOnsList Error: Expected 1 add on returned, given %d", len(addOns))
	}

	if r.Request.URL.Query().Get("per_page") != "1" {
		t.Errorf("TestAddOnsList Error: Expected per_page parameter of 1, given %s", r.Request.URL.Query().Get("per_page"))
	}

	ts, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2011-06-28T12:34:56Z")
	for _, given := range addOns {
		expected := AddOn{
			XMLName:                     xml.Name{Local: "add_on"},
			Code:                        "ipaddresses",
			Name:                        "IP Addresses",
			DefaultQuantity:             NewInt(1),
			DisplayQuantityOnHostedPage: NewBool(false),
			TaxCode:                     "digital",
			UnitAmountInCents:           &UnitAmount{USD: 200},
			AccountingCode:              "abc123",
			CreatedAt:                   NewTime(ts),
		}

		if !reflect.DeepEqual(expected, given) {
			t.Errorf("TestAddOnsList Error: expected add on to equal %#v, given %#v", expected, given)
		}
	}
}

func TestGetAddOn(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/plans/gold/add_ons/ipaddresses", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("TestGetAddOn Error: Expected %s request, given %s", "GET", r.Method)
		}
		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?>
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
		t.Errorf("TestGetAddOn Error: Error occured making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestGetAddOn Error: Expected get add_on to return OK")
	}

	ts, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2011-06-28T12:34:56Z")
	expected := AddOn{
		XMLName:                     xml.Name{Local: "add_on"},
		Code:                        "ipaddresses",
		Name:                        "IP Addresses",
		DefaultQuantity:             NewInt(1),
		DisplayQuantityOnHostedPage: NewBool(false),
		TaxCode:                     "digital",
		UnitAmountInCents:           &UnitAmount{USD: 200},
		AccountingCode:              "abc123",
		CreatedAt:                   NewTime(ts),
	}

	if !reflect.DeepEqual(expected, a) {
		t.Errorf("TestGetAddOn Error: expected account to equal %#v, given %#v", expected, a)
	}
}

func TestCreateAddOn(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/plans/gold/add_ons", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("TestCreateAddOn Error: Expected %s request, given %s", "POST", r.Method)
		}
		rw.WriteHeader(201)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?><add_on></add_on>`)
	})

	r, _, err := client.AddOns.Create("gold", AddOn{})
	if err != nil {
		t.Errorf("TestCreateAddOn Error: Error occured making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestCreateAddOn Error: Expected create add on to return OK")
	}
}

func TestUpdateAddOn(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/plans/gold/add_ons/ipaddresses", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("TestUpdateAddOn Error: Expected %s request, given %s", "PUT", r.Method)
		}
		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?><add_on></add_on>`)
	})

	r, _, err := client.AddOns.Update("gold", "ipaddresses", AddOn{})
	if err != nil {
		t.Errorf("TestUpdateAddOn Error: Error occured making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestUpdateAddOn Error: Expected update add on to return OK")
	}
}

func TestDeleteAddOn(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/plans/gold/add_ons/ipaddresses", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("TestDeleteAddOn Error: Expected %s request, given %s", "DELETE", r.Method)
		}
		rw.WriteHeader(204)
	})

	r, err := client.AddOns.Delete("gold", "ipaddresses")
	if err != nil {
		t.Errorf("TestDeleteAddOn Error: Error occured making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestDeleteAddOn Error: Expected deleted add on to return OK")
	}
}
