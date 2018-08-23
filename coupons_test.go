package recurly_test

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/launchpadcentral/recurly"
	"github.com/google/go-cmp/cmp"
)

// TestCouponsEncoding ensures structs are encoded to XML properly.
// Because Recurly supports partial updates, it's important that only defined
// fields are handled properly -- including types like booleans and integers which
// have zero values that we want to send.
func TestCoupons_Encoding(t *testing.T) {
	redeem, _ := time.Parse(recurly.DateTimeFormat, "2014-01-01T07:00:00Z")
	tests := []struct {
		v        recurly.Coupon
		expected string
	}{
		{
			v:        recurly.Coupon{XMLName: xml.Name{Local: "coupon"}},
			expected: "<coupon><id>0</id><coupon_code></coupon_code><coupon_type></coupon_type><name></name><redemption_resource></redemption_resource><state></state><single_use>false</single_use><applies_to_all_plans>false</applies_to_all_plans><duration></duration><discount_type></discount_type><applies_to_non_plan_charges>false</applies_to_non_plan_charges><plan_codes></plan_codes></coupon>",
		},
		{
			v: recurly.Coupon{
				XMLName:      xml.Name{Local: "coupon"},
				Code:         "special",
				Name:         "Special 10% off",
				DiscountType: "percent",
			},
			expected: "<coupon><id>0</id><coupon_code>special</coupon_code><coupon_type></coupon_type><name>Special 10% off</name><redemption_resource></redemption_resource><state></state><single_use>false</single_use><applies_to_all_plans>false</applies_to_all_plans><duration></duration><discount_type>percent</discount_type><applies_to_non_plan_charges>false</applies_to_non_plan_charges><plan_codes></plan_codes></coupon>",
		},
		{
			v: recurly.Coupon{
				XMLName:            xml.Name{Local: "coupon"},
				Code:               "special",
				Name:               "Special 10% off",
				State:              "redeemable",
				RedemptionResource: "account",
				Description:        "Save 10%",
				DiscountType:       "percent",
			},
			expected: "<coupon><id>0</id><coupon_code>special</coupon_code><coupon_type></coupon_type><name>Special 10% off</name><redemption_resource>account</redemption_resource><state>redeemable</state><single_use>false</single_use><applies_to_all_plans>false</applies_to_all_plans><duration></duration><discount_type>percent</discount_type><applies_to_non_plan_charges>false</applies_to_non_plan_charges><description>Save 10%</description><plan_codes></plan_codes></coupon>",
		},
		{
			v: recurly.Coupon{
				XMLName:            xml.Name{Local: "coupon"},
				Code:               "special",
				Name:               "Special 10% off",
				State:              "redeemable",
				RedemptionResource: "account",
				Description:        "Save 10%",
				DiscountType:       "percent",
				SingleUse:          true,
				AppliesToAllPlans:  true,
				DiscountPercent:    recurly.NewInt(10),
			},
			expected: "<coupon><id>0</id><coupon_code>special</coupon_code><coupon_type></coupon_type><name>Special 10% off</name><redemption_resource>account</redemption_resource><state>redeemable</state><single_use>true</single_use><applies_to_all_plans>true</applies_to_all_plans><duration></duration><discount_type>percent</discount_type><applies_to_non_plan_charges>false</applies_to_non_plan_charges><description>Save 10%</description><discount_percent>10</discount_percent><plan_codes></plan_codes></coupon>",
		},
		{
			v: recurly.Coupon{
				XMLName:                  xml.Name{Local: "coupon"},
				Code:                     "special",
				Type:                     "single_code",
				Name:                     "Special 10% off",
				DiscountType:             "dollars",
				DiscountInCents:          &recurly.UnitAmount{USD: 100},
				MaxRedemptions:           recurly.NewInt(2),
				MaxRedemptionsPerAccount: recurly.NewInt(1),
			},
			expected: "<coupon><id>0</id><coupon_code>special</coupon_code><coupon_type>single_code</coupon_type><name>Special 10% off</name><redemption_resource></redemption_resource><state></state><single_use>false</single_use><applies_to_all_plans>false</applies_to_all_plans><duration></duration><discount_type>dollars</discount_type><applies_to_non_plan_charges>false</applies_to_non_plan_charges><discount_in_cents><USD>100</USD></discount_in_cents><max_redemptions>2</max_redemptions><max_redemptions_per_account>1</max_redemptions_per_account><plan_codes></plan_codes></coupon>",
		},
		{
			v: recurly.Coupon{
				XMLName:        xml.Name{Local: "coupon"},
				Code:           "special",
				Name:           "Special 10% off",
				Duration:       "temporal",
				TemporalUnit:   "day",
				TemporalAmount: recurly.NewInt(28),
			},
			expected: "<coupon><id>0</id><coupon_code>special</coupon_code><coupon_type></coupon_type><name>Special 10% off</name><redemption_resource></redemption_resource><state></state><single_use>false</single_use><applies_to_all_plans>false</applies_to_all_plans><duration>temporal</duration><discount_type></discount_type><applies_to_non_plan_charges>false</applies_to_non_plan_charges><temporal_unit>day</temporal_unit><temporal_amount>28</temporal_amount><plan_codes></plan_codes></coupon>",
		},
		{
			v: recurly.Coupon{
				XMLName:           xml.Name{Local: "coupon"},
				Code:              "special",
				Name:              "Special 10% off",
				DiscountType:      "percent",
				AppliesToAllPlans: true,
				RedeemByDate:      recurly.NewTime(redeem),
				PlanCodes:         []string{"gold", "silver"},
			},
			expected: "<coupon><id>0</id><coupon_code>special</coupon_code><coupon_type></coupon_type><name>Special 10% off</name><redemption_resource></redemption_resource><state></state><single_use>false</single_use><applies_to_all_plans>true</applies_to_all_plans><duration></duration><discount_type>percent</discount_type><applies_to_non_plan_charges>false</applies_to_non_plan_charges><redeem_by_date>2014-01-01T07:00:00Z</redeem_by_date><plan_codes><plan_code>gold</plan_code><plan_code>silver</plan_code></plan_codes></coupon>",
		},
	}

	for _, tt := range tests {
		var buf bytes.Buffer
		if err := xml.NewEncoder(&buf).Encode(tt.v); err != nil {
			t.Fatalf("unexpected error: %v", err)
		} else if buf.String() != tt.expected {
			t.Fatalf("unexpected coupon: %v", cmp.Diff(buf.String(), tt.expected))
		}
	}

	for _, tt := range tests {
		c := recurly.Coupon{}
		if err := xml.Unmarshal([]byte(tt.expected), &c); err != nil {
			t.Fatalf("unexpected error: %v", err)
		} else if diff := cmp.Diff(tt.v, c); diff != "" {
			t.Fatalf("unexpected decode diff: %v", diff)
		}
	}
}

func TestCoupons_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/coupons", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		io.WriteString(w, `<?xml version="1.0" encoding="UTF-8"?>
        <coupons type="array">
          <coupon href="https://your-subdomain.recurly.com/v2/coupons/special">
            <redemptions href="https://your-subdomain.recurly.com/v2/coupons/special/redemptions"/>
            <id type="integer">2151093486799579392</id>
            <coupon_code>special</coupon_code>
            <coupon_type>single_code</coupon_type>
            <name>Special 10% off</name>
            <state>redeemable</state>
            <single_use>true</single_use>
            <discount_type>percent</discount_type>
            <max_redemptions type="integer">200</max_redemptions>
            <applies_to_all_plans>false</applies_to_all_plans>
            <discount_percent type="integer">10</discount_percent>
            <redeem_by_date type="datetime">2014-01-01T07:00:00Z</redeem_by_date>
            <single_use type="boolean">true</single_use>
            <applies_for_months nil="nil"></applies_for_months>
            <max_redemptions type="integer">10</max_redemptions>
            <applies_to_all_plans type="boolean">false</applies_to_all_plans>
            <duration>single_use</duration>
            <temporal_unit nil="nil"/>
            <temporal_amount nil="nil"/>
            <redemption_resource>account</redemption_resource>
            <max_redemptions_per_account nil="nil"/>
            <created_at type="datetime">2011-04-10T07:00:00Z</created_at>
            <plan_codes type="array">
              <plan_code>gold</plan_code>
              <plan_code>platinum</plan_code>
            </plan_codes>
            <a name="redeem" href="https://your-subdomain.recurly.com/v2/coupons/special/redeem" method="post"/>
          </coupon>
        </coupons>`)
	})

	resp, coupons, err := client.Coupons.List(recurly.Params{"per_page": 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected list coupons to return OK")
	} else if pp := resp.Request.URL.Query().Get("per_page"); pp != "1" {
		t.Fatalf("unexpected per_page: %s", pp)
	}

	ts, _ := time.Parse(recurly.DateTimeFormat, "2011-04-10T07:00:00Z")
	redeem, _ := time.Parse(recurly.DateTimeFormat, "2014-01-01T07:00:00Z")
	if diff := cmp.Diff(coupons, []recurly.Coupon{
		{
			XMLName:            xml.Name{Local: "coupon"},
			ID:                 2151093486799579392,
			Code:               "special",
			Name:               "Special 10% off",
			Type:               "single_code",
			State:              "redeemable",
			RedemptionResource: "account",
			DiscountType:       "percent",
			DiscountPercent:    recurly.NewInt(10),
			RedeemByDate:       recurly.NewTime(redeem),
			SingleUse:          true,
			Duration:           "single_use",
			MaxRedemptions:     recurly.NewInt(10),
			AppliesToAllPlans:  false,
			CreatedAt:          recurly.NewTime(ts),
			PlanCodes:          []string{"gold", "platinum"},
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestCoupons_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/coupons/special", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		io.WriteString(w, `<?xml version="1.0" encoding="UTF-8"?>             
            <coupon href="https://your-subdomain.recurly.com/v2/coupons/special">
             	<redemptions href="https://your-subdomain.recurly.com/v2/coupons/special/redemptions"/>
             	<id type="integer">2151093486799579392</id>
             	<coupon_code>special</coupon_code>
             	<name>20$ off</name>
             	<state>redeemable</state>
             	<coupon_type>bulk</coupon_type>
             	<discount_type>dollars</discount_type>
             	<discount_in_cents>
             	  <USD type="integer">2000</USD>
             	</discount_in_cents>
             	<redemption_resource>account</redemption_resource>
             	<unique_code_template>'savemore'99999999</unique_code_template>
             	<redeem_by_date type="datetime">2014-01-01T07:00:00Z</redeem_by_date>
             	<max_redemptions_per_account type="integer">1</max_redemptions_per_account>
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

	resp, coupon, err := client.Coupons.Get("special")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected get coupon to return OK")
	}

	ts, _ := time.Parse(recurly.DateTimeFormat, "2011-04-10T07:00:00Z")
	redeem, _ := time.Parse(recurly.DateTimeFormat, "2014-01-01T07:00:00Z")
	if diff := cmp.Diff(coupon, &recurly.Coupon{
		XMLName:                  xml.Name{Local: "coupon"},
		ID:                       2151093486799579392,
		Code:                     "special",
		Name:                     "20$ off",
		State:                    "redeemable",
		Type:                     "bulk",
		DiscountType:             "dollars",
		DiscountInCents:          &recurly.UnitAmount{USD: 2000},
		RedeemByDate:             recurly.NewTime(redeem),
		SingleUse:                true,
		RedemptionResource:       "account",
		MaxRedemptions:           recurly.NewInt(10),
		MaxRedemptionsPerAccount: recurly.NewInt(1),
		AppliesToAllPlans:        false,
		UniqueCodeTemplate:       "'savemore'99999999",
		CreatedAt:                recurly.NewTime(ts),
		PlanCodes:                []string{"gold", "platinum"},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestCoupons_Get_ErrNotFound(t *testing.T) {
	setup()
	defer teardown()

	var invoked bool
	mux.HandleFunc("/v2/coupons/special", func(w http.ResponseWriter, r *http.Request) {
		invoked = true
		w.WriteHeader(http.StatusNotFound)
	})

	_, coupon, err := client.Coupons.Get("special")
	if !invoked {
		t.Fatal("handler not invoked")
	} else if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if coupon != nil {
		t.Fatalf("expected coupon to be nil: %#v", coupon)
	}
}

func TestCoupons_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/coupons", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(201)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><coupon></coupon>`)
	})

	resp, _, err := client.Coupons.Create(recurly.Coupon{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected create coupon to return OK")
	}
}

func TestCoupons_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/coupons/special", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(204)
	})

	resp, err := client.Coupons.Delete("special")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected deleted coupon to return OK")
	}
}
