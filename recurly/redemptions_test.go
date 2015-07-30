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

// TestRedemptionEncoding ensures structs are encoded to XML properly.
func TestRedemptionsEncoding(t *testing.T) {
	suite := []map[string]interface{}{
		map[string]interface{}{"struct": Redemption{}, "xml": "<redemption></redemption>"},
	}

	for _, s := range suite {
		buf := new(bytes.Buffer)
		err := xml.NewEncoder(buf).Encode(s["struct"])
		if err != nil {
			t.Errorf("TestRedemptionsEncoding Error: %s", err)
		}

		if buf.String() != s["xml"] {
			t.Errorf("TestRedemptionsEncoding Error: Expected %s, given %s", s["xml"], buf.String())
		}
	}
}

func TestGetForAccountRedemption(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/1/redemption", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("TestGetForAccountRedemption Error: Expected %s request, given %s", "GET", r.Method)
		}
		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?>
        <redemption href="https://your-subdomain.recurly.com/v2/accounts/1/redemption">
            <coupon href="https://your-subdomain.recurly.com/v2/coupons/special"/>
            <account href="https://your-subdomain.recurly.com/v2/accounts/1"/>
            <single_use type="boolean">false</single_use>
            <total_discounted_in_cents type="integer">0</total_discounted_in_cents>
            <currency>USD</currency>
            <state>active</state>
            <created_at type="datetime">2011-06-27T12:34:56Z</created_at>
        </redemption>`)
	})

	r, a, err := client.Redemptions.GetForAccount("1")
	if err != nil {
		t.Errorf("TestGetForAccountRedemption Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestGetForAccountRedemption Error: Expected get redemption to return OK")
	}

	ts, _ := time.Parse(datetimeFormat, "2011-06-27T12:34:56Z")
	expected := Redemption{
		XMLName: xml.Name{Local: "redemption"},
		Coupon: href{
			Code: "special",
			HREF: "https://your-subdomain.recurly.com/v2/coupons/special",
		},
		Account: href{
			Code: "1",
			HREF: "https://your-subdomain.recurly.com/v2/accounts/1",
		},
		SingleUse:              NewBool(false),
		TotalDiscountedInCents: 0,
		Currency:               "USD",
		State:                  "active",
		CreatedAt:              NewTime(ts),
	}

	if !reflect.DeepEqual(expected, a) {
		t.Errorf("TestGetForAccountRedemption Error: expected account to equal %#v, given %#v", expected, a)
	}
}

func TestGetForInvoiceRedemption(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/invoices/1108/redemption", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("TestGetForInvoiceRedemption Error: Expected %s request, given %s", "GET", r.Method)
		}
		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?>
        <redemption href="https://your-subdomain.recurly.com/v2/accounts/1/redemption">
            <coupon href="https://your-subdomain.recurly.com/v2/coupons/special"/>
            <account href="https://your-subdomain.recurly.com/v2/accounts/1"/>
            <single_use type="boolean">true</single_use>
            <total_discounted_in_cents type="integer">0</total_discounted_in_cents>
            <currency>USD</currency>
            <state>inactive</state>
            <created_at type="datetime">2011-06-27T12:34:56Z</created_at>
        </redemption>`)
	})

	r, a, err := client.Redemptions.GetForInvoice("1108")
	if err != nil {
		t.Errorf("TestGetForInvoiceRedemption Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestGetForInvoiceRedemption Error: Expected get redemption to return OK")
	}

	ts, _ := time.Parse(datetimeFormat, "2011-06-27T12:34:56Z")
	expected := Redemption{
		XMLName: xml.Name{Local: "redemption"},
		Coupon: href{
			Code: "special",
			HREF: "https://your-subdomain.recurly.com/v2/coupons/special",
		},
		Account: href{
			Code: "1",
			HREF: "https://your-subdomain.recurly.com/v2/accounts/1",
		},
		SingleUse:              NewBool(true),
		TotalDiscountedInCents: 0,
		Currency:               "USD",
		State:                  "inactive",
		CreatedAt:              NewTime(ts),
	}

	if !reflect.DeepEqual(expected, a) {
		t.Errorf("TestGetForInvoiceRedemption Error: expected account to equal %#v, given %#v", expected, a)
	}
}

func TestRedeemCoupon(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/coupons/special/redeem", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("TestRedeemCoupon Error: Expected %s request, given %s", "POST", r.Method)
		}
		given := new(bytes.Buffer)
		given.ReadFrom(r.Body)
		expected := "<redemption><account_code>1</account_code><currency>USD</currency></redemption>"
		if expected != given.String() {
			t.Errorf("TestRedeemCoupon Error: Expected request body of %s, given %s", expected, given.String())
		}

		rw.WriteHeader(201)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?><redemption></redemption>`)
	})

	r, _, err := client.Redemptions.Redeem("special", "1", "USD")
	if err != nil {
		t.Errorf("TestRedeemCoupon Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestRedeemCoupon Error: Expected redeeming add on to return OK")
	}
}

func TestDeleteRedemption(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/27/redemption", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("TestRemoveRedemption Error: Expected %s request, given %s", "Delete", r.Method)
		}
		rw.WriteHeader(204)
	})

	r, err := client.Redemptions.Delete("27")
	if err != nil {
		t.Errorf("TestRemoveRedemption Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestRemoveRedemption Error: Expected delete add on to return OK")
	}
}
