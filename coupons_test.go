package recurly_test

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/blacklightcms/recurly"
	"github.com/google/go-cmp/cmp"
)

// Ensure structs are encoded to XML properly.
func TestCoupons_Encoding(t *testing.T) {
	redeem, _ := time.Parse(recurly.DateTimeFormat, "2014-01-01T07:00:00Z")
	tests := []struct {
		v        recurly.Coupon
		expected string
	}{
		{
			v: recurly.Coupon{XMLName: xml.Name{Local: "coupon"}},
			expected: MustCompactString(`
				<coupon>
					<coupon_code></coupon_code>
					<name></name>
					<discount_type></discount_type>
					<plan_codes></plan_codes>
				</coupon>
			`),
		},
		{
			v: recurly.Coupon{
				XMLName:      xml.Name{Local: "coupon"},
				Code:         "special",
				Name:         "Special 10% off",
				DiscountType: "percent",
			},
			expected: MustCompactString(`
				<coupon>
					<coupon_code>special</coupon_code>
					<name>Special 10% off</name>
					<discount_type>percent</discount_type>
					<plan_codes></plan_codes>
				</coupon>
			`),
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
			expected: MustCompactString(`
				<coupon>
					<coupon_code>special</coupon_code>
					<name>Special 10% off</name>
					<redemption_resource>account</redemption_resource>
					<state>redeemable</state>
					<discount_type>percent</discount_type>
					<description>Save 10%</description>
					<plan_codes></plan_codes>
				</coupon>
			`),
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
				AppliesToAllPlans:  true,
				DiscountPercent:    recurly.NewInt(10),
			},
			expected: MustCompactString(`
				<coupon>
					<coupon_code>special</coupon_code>
					<name>Special 10% off</name>
					<redemption_resource>account</redemption_resource>
					<state>redeemable</state>
					<applies_to_all_plans>true</applies_to_all_plans>
					<discount_type>percent</discount_type>
					<description>Save 10%</description>
					<discount_percent>10</discount_percent>
					<plan_codes></plan_codes>
				</coupon>
			`),
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
			expected: MustCompactString(`
				<coupon>
					<coupon_code>special</coupon_code>
					<coupon_type>single_code</coupon_type>
					<name>Special 10% off</name>
					<discount_type>dollars</discount_type>
					<discount_in_cents>
					<USD>100</USD>
					</discount_in_cents>
					<max_redemptions>2</max_redemptions>
					<max_redemptions_per_account>1</max_redemptions_per_account>
					<plan_codes></plan_codes>
				</coupon>
			`),
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
			expected: MustCompactString(`
				<coupon>
					<coupon_code>special</coupon_code>
					<name>Special 10% off</name>
					<duration>temporal</duration>
					<discount_type></discount_type>
					<temporal_unit>day</temporal_unit>
					<temporal_amount>28</temporal_amount>
					<plan_codes></plan_codes>
				</coupon>
			`),
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
			expected: MustCompactString(`
				<coupon>
					<coupon_code>special</coupon_code>
					<name>Special 10% off</name>
					<applies_to_all_plans>true</applies_to_all_plans>
					<discount_type>percent</discount_type>
					<redeem_by_date>2014-01-01T07:00:00Z</redeem_by_date>
					<plan_codes>
						<plan_code>gold</plan_code>
						<plan_code>silver</plan_code>
					</plan_codes>
				</coupon>
			`),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Encode/%d", i), func(t *testing.T) {
			buf := new(bytes.Buffer)
			if err := xml.NewEncoder(buf).Encode(tt.v); err != nil {
				t.Fatal(err)
			} else if buf.String() != tt.expected {
				t.Fatal(buf.String())
			}
		})

		t.Run(fmt.Sprintf("Decode/%d", i), func(t *testing.T) {
			var c recurly.Coupon
			if err := xml.Unmarshal([]byte(tt.expected), &c); err != nil {
				t.Fatal(err)
			} else if diff := cmp.Diff(tt.v, c); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestCoupons_List(t *testing.T) {
	client, s := NewServer()
	defer s.Close()

	var invocations int
	s.HandleFunc("GET", "/v2/coupons", func(w http.ResponseWriter, r *http.Request) {
		invocations++
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("coupons.xml"))
	}, t)

	pager := client.Coupons.List(nil)
	for pager.Next() {
		var coupons []recurly.Coupon
		if err := pager.Fetch(context.Background(), &coupons); err != nil {
			t.Fatal(err)
		} else if !s.Invoked {
			t.Fatal("expected s to be invoked")
		} else if diff := cmp.Diff(coupons, []recurly.Coupon{*NewTestCoupon()}); diff != "" {
			t.Fatal(diff)
		}
	}
	if invocations != 1 {
		t.Fatalf("unexpected number of invocations: %d", invocations)
	}
}

func TestCoupons_Get(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		client, s := NewServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/coupons/special", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write(MustOpenFile("coupon.xml"))
		}, t)

		if coupon, err := client.Coupons.Get(context.Background(), "special"); err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(coupon, NewTestCoupon()); diff != "" {
			t.Fatal(diff)
		} else if !s.Invoked {
			t.Fatal("expected fn invocation")
		}
	})

	// Ensure a 404 returns nil values.
	t.Run("ErrNotFound", func(t *testing.T) {
		client, s := NewServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/coupons/special", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}, t)

		if coupon, err := client.Coupons.Get(context.Background(), "special"); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if coupon != nil {
			t.Fatalf("expected nil: %#v", coupon)
		}
	})
}

func TestCoupons_Create(t *testing.T) {
	client, s := NewServer()
	defer s.Close()

	s.HandleFunc("POST", "/v2/coupons", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write(MustOpenFile("coupon.xml"))
	}, t)

	if coupon, err := client.Coupons.Create(context.Background(), recurly.Coupon{}); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(coupon, NewTestCoupon()); diff != "" {
		t.Fatal(diff)
	}
}

func TestCoupons_Update(t *testing.T) {
	client, s := NewServer()
	defer s.Close()

	s.HandleFunc("PUT", "/v2/coupons/special", func(w http.ResponseWriter, r *http.Request) {
		if str := MustReadAllString(r.Body); str != MustCompactString(`
			<coupon>
				<name>New Coupon Name</name>
				<description>New coupon description for the hosted pages.</description>
				<invoice_description>New coupon description for the invoice.</invoice_description>
				<redeem_by_date>2011-04-10T07:00:00Z</redeem_by_date>
				<max_redemptions>500</max_redemptions>
				<max_redemptions_per_account>1</max_redemptions_per_account>
			</coupon>
		`) {
			t.Fatal(str)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(MustOpenFile("coupon.xml"))
	}, t)

	if coupon, err := client.Coupons.Update(context.Background(), "special", recurly.Coupon{
		Name:                     "New Coupon Name",
		Description:              "New coupon description for the hosted pages.",
		InvoiceDescription:       "New coupon description for the invoice.",
		RedeemByDate:             recurly.NewTime(MustParseTime("2011-04-10T07:00:00Z")),
		MaxRedemptions:           recurly.NewInt(500),
		MaxRedemptionsPerAccount: recurly.NewInt(1),
		// Send extra coupon fields to assert they are not sent
		ID:                 1,
		State:              "redeemable",
		Type:               "bulk",
		DiscountType:       "dollars",
		DiscountInCents:    &recurly.UnitAmount{USD: 2000},
		RedemptionResource: "account",
		AppliesToAllPlans:  false,
		UniqueCodeTemplate: "'savemore'99999999",
		CreatedAt:          recurly.NewTime(MustParseTime("2011-04-10T07:00:00Z")),
		PlanCodes:          []string{"gold", "platinum"},
	}); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(coupon, NewTestCoupon()); diff != "" {
		t.Fatal(diff)
	}
}

func TestCoupons_Restore(t *testing.T) {
	t.Run("Edits", func(t *testing.T) {
		client, s := NewServer()
		defer s.Close()

		s.HandleFunc("PUT", "/v2/coupons/special/restore", func(w http.ResponseWriter, r *http.Request) {
			if str := MustReadAllString(r.Body); str != MustCompactString(`
			<coupon>
				<name>New Coupon Name</name>
				<description>New coupon description for the hosted pages.</description>
				<invoice_description>New coupon description for the invoice.</invoice_description>
				<redeem_by_date>2011-04-10T07:00:00Z</redeem_by_date>
				<max_redemptions>500</max_redemptions>
				<max_redemptions_per_account>1</max_redemptions_per_account>
			</coupon>
		`) {
				t.Fatal(str)
			}
			w.WriteHeader(http.StatusCreated)
			w.Write(MustOpenFile("coupon.xml"))
		}, t)

		if coupon, err := client.Coupons.Restore(context.Background(), "special", recurly.Coupon{
			Name:                     "New Coupon Name",
			Description:              "New coupon description for the hosted pages.",
			InvoiceDescription:       "New coupon description for the invoice.",
			RedeemByDate:             recurly.NewTime(MustParseTime("2011-04-10T07:00:00Z")),
			MaxRedemptions:           recurly.NewInt(500),
			MaxRedemptionsPerAccount: recurly.NewInt(1),
			// Send extra coupon fields to assert they are not sent
			ID:                 1,
			State:              "redeemable",
			Type:               "bulk",
			DiscountType:       "dollars",
			DiscountInCents:    &recurly.UnitAmount{USD: 2000},
			RedemptionResource: "account",
			AppliesToAllPlans:  false,
			UniqueCodeTemplate: "'savemore'99999999",
			CreatedAt:          recurly.NewTime(MustParseTime("2011-04-10T07:00:00Z")),
			PlanCodes:          []string{"gold", "platinum"},
		}); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(coupon, NewTestCoupon()); diff != "" {
			t.Fatal(diff)
		}
	})

	t.Run("NoEdits", func(t *testing.T) {
		client, s := NewServer()
		defer s.Close()

		s.HandleFunc("PUT", "/v2/coupons/special/restore", func(w http.ResponseWriter, r *http.Request) {
			if str := MustReadAllString(r.Body); str != MustCompactString(`
			<coupon>
			</coupon>
		`) {
				t.Fatal(str)
			}
			w.WriteHeader(http.StatusCreated)
			w.Write(MustOpenFile("coupon.xml"))
		}, t)

		if coupon, err := client.Coupons.Restore(context.Background(), "special", recurly.Coupon{}); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(coupon, NewTestCoupon()); diff != "" {
			t.Fatal(diff)
		}
	})
}

func TestCoupons_Delete(t *testing.T) {
	client, s := NewServer()
	defer s.Close()

	s.HandleFunc("DELETE", "/v2/coupons/special", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}, t)

	if err := client.Coupons.Delete(context.Background(), "special"); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	}
}

func TestCoupons_Generate(t *testing.T) {
	client, s := NewServer()
	defer s.Close()

	s.HandleFunc("POST", "/v2/coupons/special/generate", func(w http.ResponseWriter, r *http.Request) {
		if str := MustReadAllString(r.Body); str != MustCompactString(`
			<coupon>
				<number_of_unique_codes>200</number_of_unique_codes>
	  		</coupon>
		`) {
			t.Fatal(str)
		}

		w.Header().Set("Location", "https://your-subdomain.recurly.com/v2/coupons/special/unique_coupon_codes?cursor=1998184141762793924:1468970111&per_page=200")

		w.WriteHeader(http.StatusCreated)
		w.Write(MustOpenFile("coupon.xml"))
	}, t)

	// Generate unique codes.
	pager, err := client.Coupons.Generate(context.Background(), "special", 200)
	if !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	}

	// Setup handler to test pager.
	var invocations int
	s.HandleFunc("GET", "/v2/coupons/special/unique_coupon_codes", func(w http.ResponseWriter, r *http.Request) {
		if v := r.URL.Query().Get("cursor"); v != "1998184141762793924:1468970111" {
			t.Fatalf("unexpected cursor: %q", v)
		} else if v = r.URL.Query().Get("per_page"); v != "200" {
			t.Fatalf("unexpected per_page: %q", v)
		}

		invocations++
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("coupons.xml"))
	}, t)
	s.Invoked = false // reset invoked bool

	// Test pager.
	for pager.Next() {
		var coupons []recurly.Coupon
		if err := pager.Fetch(context.Background(), &coupons); err != nil {
			t.Fatal(err)
		} else if !s.Invoked {
			t.Fatal("expected s to be invoked")
		} else if diff := cmp.Diff(coupons, []recurly.Coupon{*NewTestCoupon()}); diff != "" {
			t.Fatal(diff)
		}
	}
	if invocations != 1 {
		t.Fatalf("unexpected number of invocations: %d", invocations)
	}
}

// Returns a Coupon corresponding to testdata/coupon.xml.
func NewTestCoupon() *recurly.Coupon {
	return &recurly.Coupon{
		XMLName:                  xml.Name{Local: "coupon"},
		ID:                       2151093486799579392,
		Code:                     "special",
		Name:                     "20$ off",
		State:                    "redeemable",
		Type:                     "bulk",
		DiscountType:             "dollars",
		DiscountInCents:          &recurly.UnitAmount{USD: 2000},
		RedeemByDate:             recurly.NewTime(MustParseTime("2014-01-01T07:00:00Z")),
		RedemptionResource:       "account",
		MaxRedemptions:           recurly.NewInt(10),
		MaxRedemptionsPerAccount: recurly.NewInt(1),
		AppliesToAllPlans:        false,
		UniqueCodeTemplate:       "'savemore'99999999",
		CreatedAt:                recurly.NewTime(MustParseTime("2011-04-10T07:00:00Z")),
		PlanCodes:                []string{"gold", "platinum"},
	}
}
