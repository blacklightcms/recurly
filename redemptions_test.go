package recurly_test

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/launchpadcentral/recurly"
	"github.com/google/go-cmp/cmp"
)

func TestRedemptions_GetForAccount(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/1/redemptions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
      <redemptions type="array">
        <redemption href="https://your-subdomain.recurly.com/v2/accounts/1/redemption">
            <coupon href="https://your-subdomain.recurly.com/v2/coupons/special"/>
            <account href="https://your-subdomain.recurly.com/v2/accounts/1"/>
            <subscription href="https://your-subdomain.recurly.com/v2/subscriptions/37bfef7a8e44cfc3817b7a43eba8a6e6" />
            <uuid>374a1c75374bd81493a3f7425db0a2b8</uuid>
             <single_use type="boolean">false</single_use>
            <total_discounted_in_cents type="integer">0</total_discounted_in_cents>
            <currency>USD</currency>
            <state>active</state>
            <coupon_code>special</coupon_code>
            <created_at type="datetime">2011-06-27T12:34:56Z</created_at>
            <updated_at type="datetime">2011-06-27T12:34:56Z</updated_at>
        </redemption>
     </redemptions>`)
	})

	r, redemptions, err := client.Redemptions.GetForAccount("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected get redemption to return OK")
	}

	ts, _ := time.Parse(recurly.DateTimeFormat, "2011-06-27T12:34:56Z")
	if diff := cmp.Diff(redemptions[0], recurly.Redemption{
		UUID:                   "374a1c75374bd81493a3f7425db0a2b8",
		SubscriptionUUID:       "37bfef7a8e44cfc3817b7a43eba8a6e6",
		AccountCode:            "1",
		CouponCode:             "special",
		SingleUse:              false,
		TotalDiscountedInCents: 0,
		Currency:               "USD",
		State:                  recurly.RedemptionStateActive,
		CreatedAt:              recurly.NewTime(ts),
		UpdatedAt:              recurly.NewTime(ts),
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestRedemptions_GetForAccount_ErrNotFound(t *testing.T) {
	setup()
	defer teardown()

	var invoked bool
	mux.HandleFunc("/v2/accounts/1/redemptions", func(w http.ResponseWriter, r *http.Request) {
		invoked = true
		w.WriteHeader(http.StatusNotFound)
	})

	_, redemptions, err := client.Redemptions.GetForAccount("1")
	if !invoked {
		t.Fatal("handler not invoked")
	} else if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if len(redemptions) != 0 {
		t.Fatalf("expect zero redemptions: %v", redemptions)
	}
}

func TestRedemptions_GetForInvoice(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/invoices/1108/redemptions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("expected %s request, given %s", "GET", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
      <redemptions type="array">
        <redemption href="https://your-subdomain.recurly.com/v2/accounts/1/redemption">
            <coupon href="https://your-subdomain.recurly.com/v2/coupons/special"/>
            <account href="https://your-subdomain.recurly.com/v2/accounts/1"/>
            <subscription href="https://your-subdomain.recurly.com/v2/subscriptions/37bfef7a8e44cfc3817b7a43eba8a6e6" />
            <uuid>374a1c75374bd81493a3f7425db0a2b8</uuid>
            <single_use type="boolean">true</single_use>
            <total_discounted_in_cents type="integer">0</total_discounted_in_cents>
            <currency>USD</currency>
            <state>inactive</state>
            <coupon_code>special</coupon_code>
            <created_at type="datetime">2011-06-27T12:34:56Z</created_at>
            <updated_at type="datetime">2011-06-27T12:34:56Z</updated_at>
        </redemption>
      </redemptions>`)
	})

	r, redemptions, err := client.Redemptions.GetForInvoice("1108")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected get redemption to return OK")
	}

	ts, _ := time.Parse(recurly.DateTimeFormat, "2011-06-27T12:34:56Z")
	if diff := cmp.Diff(redemptions[0], recurly.Redemption{
		UUID:                   "374a1c75374bd81493a3f7425db0a2b8",
		SubscriptionUUID:       "37bfef7a8e44cfc3817b7a43eba8a6e6",
		AccountCode:            "1",
		CouponCode:             "special",
		SingleUse:              true,
		TotalDiscountedInCents: 0,
		Currency:               "USD",
		State:                  recurly.RedemptionStateInactive,
		CreatedAt:              recurly.NewTime(ts),
		UpdatedAt:              recurly.NewTime(ts),
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestRedemptions_GetForInvoice_ErrNotFound(t *testing.T) {
	setup()
	defer teardown()

	var invoked bool
	mux.HandleFunc("/v2/invoices/1108/redemptions", func(w http.ResponseWriter, r *http.Request) {
		invoked = true
		w.WriteHeader(http.StatusNotFound)
	})

	_, redemption, err := client.Redemptions.GetForInvoice("1108")
	if !invoked {
		t.Fatal("handler not invoked")
	} else if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if redemption != nil {
		t.Fatalf("expected redemption to be nil: %#v", redemption)
	}
}

func TestRedemptions_RedeemCoupon(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/coupons/special/redeem", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		var given bytes.Buffer
		given.ReadFrom(r.Body)
		expected := "<redemption><account_code>1</account_code><currency>USD</currency></redemption>"
		if expected != given.String() {
			t.Fatalf("unexpected input: %s", given.String())
		}

		w.WriteHeader(201)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><redemption></redemption>`)
	})

	r, _, err := client.Redemptions.Redeem("special", "1", "USD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected redeeming add on to return OK")
	}
}

func TestRedemptions_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/27/redemption", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(204)
	})

	r, err := client.Redemptions.Delete("27")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected delete add on to return OK")
	}
}
