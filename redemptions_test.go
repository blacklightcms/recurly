package recurly_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/blacklightcms/recurly"
	"github.com/google/go-cmp/cmp"
)

func TestRedemptions_ListAccount(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	var invocations int
	s.HandleFunc("GET", "/v2/accounts/1/redemptions", func(w http.ResponseWriter, r *http.Request) {
		invocations++
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("redemptions.xml"))
	}, t)

	pager := client.Redemptions.ListAccount("1", nil)
	for pager.Next() {
		var redemptions []recurly.Redemption
		if err := pager.Fetch(context.Background(), &redemptions); err != nil {
			t.Fatal(err)
		} else if !s.Invoked {
			t.Fatal("expected s to be invoked")
		} else if diff := cmp.Diff(redemptions, []recurly.Redemption{*NewTestRedemption()}); diff != "" {
			t.Fatal(diff)
		}
	}
	if invocations != 1 {
		t.Fatalf("unexpected number of invocations: %d", invocations)
	}
}

func TestRedemptions_ListInvoice(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	var invocations int
	s.HandleFunc("GET", "/v2/invoices/1010/redemptions", func(w http.ResponseWriter, r *http.Request) {
		invocations++
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("redemptions.xml"))
	}, t)

	pager := client.Redemptions.ListInvoice(1010, nil)
	for pager.Next() {
		var redemptions []recurly.Redemption
		if err := pager.Fetch(context.Background(), &redemptions); err != nil {
			t.Fatal(err)
		} else if !s.Invoked {
			t.Fatal("expected s to be invoked")
		} else if diff := cmp.Diff(redemptions, []recurly.Redemption{*NewTestRedemption()}); diff != "" {
			t.Fatal(diff)
		}
	}
	if invocations != 1 {
		t.Fatalf("unexpected number of invocations: %d", invocations)
	}
}

func TestRedemptions_ListSubscription(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	var invocations int
	s.HandleFunc("GET", "/v2/subscriptions/37bfef7a8e44cfc3817b7a43eba8a6e6/redemptions", func(w http.ResponseWriter, r *http.Request) {
		invocations++
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("redemptions.xml"))
	}, t)

	pager := client.Redemptions.ListSubscription("37bfef7a-8e44-cfc3-817b-7a43eba8a6e6", nil) // UUID should be sanitized
	for pager.Next() {
		var redemptions []recurly.Redemption
		if err := pager.Fetch(context.Background(), &redemptions); err != nil {
			t.Fatal(err)
		} else if !s.Invoked {
			t.Fatal("expected s to be invoked")
		} else if diff := cmp.Diff(redemptions, []recurly.Redemption{*NewTestRedemption()}); diff != "" {
			t.Fatal(diff)
		}
	}
	if invocations != 1 {
		t.Fatalf("unexpected number of invocations: %d", invocations)
	}
}

func TestRedemptions_Redeem(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("POST", "/v2/coupons/special/redeem", func(w http.ResponseWriter, r *http.Request) {
			if str := MustReadAllString(r.Body); str != MustCompactString(`
				<redemption>
					<account_code>1</account_code>
					<currency>USD</currency>
				</redemption>
			`) {
				t.Fatal(str)
			}
			w.WriteHeader(http.StatusCreated)
			w.Write(MustOpenFile("redemption.xml"))
		}, t)

		if redemption, err := client.Redemptions.Redeem(context.Background(), "special", recurly.CouponRedemption{
			AccountCode: "1",
			Currency:    "USD",
		}); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(redemption, NewTestRedemption()); diff != "" {
			t.Fatal(diff)
		}
	})

	t.Run("Subscription", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("POST", "/v2/coupons/special/redeem", func(w http.ResponseWriter, r *http.Request) {
			if str := MustReadAllString(r.Body); str != MustCompactString(`
				<redemption>
					<account_code>1</account_code>
					<currency>USD</currency>
					<subscription_uuid>37bfef7a8e44cfc3817b7a43eba8a6e6</subscription_uuid>
				</redemption>
			`) {
				t.Fatal(str)
			}
			w.WriteHeader(http.StatusCreated)
			w.Write(MustOpenFile("redemption.xml"))
		}, t)

		if redemption, err := client.Redemptions.Redeem(context.Background(), "special", recurly.CouponRedemption{
			AccountCode:      "1",
			Currency:         "USD",
			SubscriptionUUID: "37bfef7a-8e44-cfc3-817b-7a43eba8a6e6",
		}); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(redemption, NewTestRedemption()); diff != "" {
			t.Fatal(diff)
		}
	})
}

func TestRedemptions_Delete(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("DELETE", "/v2/accounts/1/redemptions/3223", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}, t)

	if err := client.Redemptions.Delete(context.Background(), "1", "3223"); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	}
}

// Returns a Redemption corresponding to testdata/redemption.xml.
func NewTestRedemption() *recurly.Redemption {
	return &recurly.Redemption{
		UUID:                   "374a1c75374bd81493a3f7425db0a2b8",
		SubscriptionUUID:       "37bfef7a8e44cfc3817b7a43eba8a6e6",
		AccountCode:            "1",
		CouponCode:             "special",
		SingleUse:              true,
		TotalDiscountedInCents: 0,
		Currency:               "USD",
		State:                  recurly.RedemptionStateActive,
		CreatedAt:              recurly.NewTime(MustParseTime("2016-07-11T18:56:20Z")),
		UpdatedAt:              recurly.NewTime(MustParseTime("2016-07-11T18:56:20Z")),
	}
}
