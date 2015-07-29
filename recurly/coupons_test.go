package recurly

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"
)

// TestCouponsEncoding ensures structs are encoded to XML properly.
// Because Recurly supports partial updates, it's important that only defined
// fields are handled properly -- including types like booleans and integers which
// have zero values that we want to send.
func TestCouponsEncoding(t *testing.T) {
	suite := []map[string]interface{}{
		map[string]interface{}{"struct": Coupon{}, "xml": "<coupon><coupon_code></coupon_code><name></name><discount_type></discount_type></coupon>"},
	}

	for _, s := range suite {
		buf := new(bytes.Buffer)
		err := xml.NewEncoder(buf).Encode(s["struct"])
		if err != nil {
			t.Errorf("TestCouponsEncoding Error: %s", err)
		}

		if buf.String() != s["xml"] {
			t.Errorf("TestCouponsEncoding Error: Expected %s, given %s", s["xml"], buf.String())
		}
	}
}

func TestCouponsList(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/coupons", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("TestCouponsList Error: Expected %s request, given %s", "GET", r.Method)
		}
		rw.WriteHeader(200)
		io.WriteString(rw, `<?xml version="1.0" encoding="UTF-8"?>
        <coupons type="array">
        	<coupon href="https://your-subdomain.recurly.com/v2/coupons/special">
        		<redemptions href="https://your-subdomain.recurly.com/v2/coupons/special/redemptions"/>
        		<coupon_code>special</coupon_code>
        		<name>Special 10% off</name>
        		<state>redeemable</state>
        		<discount_type>percent</discount_type>
        		<discount_percent type="integer">10</discount_percent>
        		<redeem_by_date type="datetime">2014-01-01T07:00:00Z</redeem_by_date>
        		<single_use type="boolean">true</single_use>
        		<applies_for_months nil="nil"></applies_for_months>
        		<max_redemptions type="integer">10</max_redemptions>
        		<applies_to_all_plans type="boolean">false</applies_to_all_plans>
        		<created_at type="datetime">2011-04-10T07:00:00Z</created_at>
        		<plan_codes type="array">
        			<plan_code>gold</plan_code>
        			<plan_code>platinum</plan_code>
        		</plan_codes>
        		<a name="redeem" href="https://your-subdomain.recurly.com/v2/coupons/special/redeem" method="post"/>
        	</coupon>
        </coupons>`)
	})

	r, coupons, err := client.Coupons.List(Params{"per_page": 1})
	if err != nil {
		t.Errorf("TestCouponsList Error: Error occured making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestCouponsList Error: Expected list coupons to return OK")
	}

	if len(coupons) != 1 {
		t.Fatalf("TestCouponsList Error: Expected 1 coupon returned, given %d", len(coupons))
	}

	if r.Request.URL.Query().Get("per_page") != "1" {
		t.Errorf("TestCouponsList Error: Expected per_page parameter of 1, given %s", r.Request.URL.Query().Get("per_page"))
	}

	ts, _ := time.Parse(datetimeFormat, "2011-04-10T07:00:00Z")
	redeem, _ := time.Parse(datetimeFormat, "2014-01-01T07:00:00Z")
	for _, given := range coupons {
		expected := Coupon{
			XMLName:           xml.Name{Local: "coupon"},
			Code:              "special",
			Name:              "Special 10% off",
			State:             "redeemable",
			DiscountType:      "percent",
			DiscountPercent:   10,
			RedeemByDate:      NewTime(redeem),
			SingleUse:         NewBool(true),
			MaxRedemptions:    NewInt(10),
			AppliesToAllPlans: NewBool(false),
			CreatedAt:         NewTime(ts),
			PlanCodes: &[]CouponPlanCode{
				CouponPlanCode{Code: "gold"},
				CouponPlanCode{Code: "platinum"},
			},
		}

		if !reflect.DeepEqual(expected, given) {
			t.Errorf("TestCouponsList Error: expected coupon to equal %#v, given %#v", expected, given)
		}
	}
}

func TestGetCoupon(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/coupons/special", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("TestGetCoupon Error: Expected %s request, given %s", "GET", r.Method)
		}
		rw.WriteHeader(200)
		io.WriteString(rw, `<?xml version="1.0" encoding="UTF-8"?>
            <coupon href="https://your-subdomain.recurly.com/v2/coupons/special">
        		<redemptions href="https://your-subdomain.recurly.com/v2/coupons/special/redemptions"/>
        		<coupon_code>special</coupon_code>
        		<name>Special 10% off</name>
        		<state>redeemable</state>
        		<discount_type>percent</discount_type>
        		<discount_percent type="integer">10</discount_percent>
        		<redeem_by_date type="datetime">2014-01-01T07:00:00Z</redeem_by_date>
        		<single_use type="boolean">true</single_use>
        		<applies_for_months nil="nil"></applies_for_months>
        		<max_redemptions type="integer">10</max_redemptions>
        		<applies_to_all_plans type="boolean">false</applies_to_all_plans>
        		<created_at type="datetime">2011-04-10T07:00:00Z</created_at>
        		<plan_codes type="array">
        			<plan_code>gold</plan_code>
        			<plan_code>platinum</plan_code>
        		</plan_codes>
        		<a name="redeem" href="https://your-subdomain.recurly.com/v2/coupons/special/redeem" method="post"/>
        	</coupon>`)
	})

	r, a, err := client.Coupons.Get("special")
	if err != nil {
		t.Errorf("TestGetCoupon Error: Error occured making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestGetCoupon Error: Expected get coupon to return OK")
	}

	ts, _ := time.Parse(datetimeFormat, "2011-04-10T07:00:00Z")
	redeem, _ := time.Parse(datetimeFormat, "2014-01-01T07:00:00Z")
	expected := Coupon{
		XMLName:           xml.Name{Local: "coupon"},
		Code:              "special",
		Name:              "Special 10% off",
		State:             "redeemable",
		DiscountType:      "percent",
		DiscountPercent:   10,
		RedeemByDate:      NewTime(redeem),
		SingleUse:         NewBool(true),
		MaxRedemptions:    NewInt(10),
		AppliesToAllPlans: NewBool(false),
		CreatedAt:         NewTime(ts),
		PlanCodes: &[]CouponPlanCode{
			CouponPlanCode{Code: "gold"},
			CouponPlanCode{Code: "platinum"},
		},
	}

	if !reflect.DeepEqual(expected, a) {
		t.Errorf("TestGetCoupon Error: expected account to equal %#v, given %#v", expected, a)
	}
}

func TestCreateCoupon(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/coupons", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("TestCreateCoupon Error: Expected %s request, given %s", "POST", r.Method)
		}
		rw.WriteHeader(201)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?><coupon></coupon>`)
	})

	r, _, err := client.Coupons.Create(Coupon{})
	if err != nil {
		t.Errorf("TestCreateCoupon Error: Error occured making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestCreateCoupon Error: Expected create coupon to return OK")
	}
}

func TestDeleteCoupon(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/coupons/special", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("TestDeleteCoupon Error: Expected %s request, given %s", "DELETE", r.Method)
		}
		rw.WriteHeader(204)
	})

	r, err := client.Coupons.Delete("special")
	if err != nil {
		t.Errorf("TestDeleteCoupon Error: Error occured making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestDeleteCoupon Error: Expected deleted coupon to return OK")
	}
}
