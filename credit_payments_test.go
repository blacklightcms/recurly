package recurly_test

import (
	"context"
	"encoding/xml"
	"net/http"
	"testing"

	"github.com/blacklightcms/recurly"
	"github.com/google/go-cmp/cmp"
)

func TestCreditPayments_List(t *testing.T) {
	client, s := NewServer()
	defer s.Close()

	var invocations int
	s.HandleFunc("GET", "/v2/credit_payments", func(w http.ResponseWriter, r *http.Request) {
		invocations++
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("credit_payments.xml"))
	}, t)

	pager := client.CreditPayments.List(nil)
	for pager.Next() {
		if pmts, err := pager.Fetch(context.Background()); err != nil {
			t.Fatal(err)
		} else if !s.Invoked {
			t.Fatal("expected s to be invoked")
		} else if diff := cmp.Diff(pmts, []recurly.CreditPayment{*NewTestCreditPayment()}); diff != "" {
			t.Fatal(diff)
		}
	}
	if invocations != 1 {
		t.Fatalf("unexpected number of invocations: %d", invocations)
	}
}

func TestCreditPayments_ListAccount(t *testing.T) {
	client, s := NewServer()
	defer s.Close()

	var invocations int
	s.HandleFunc("GET", "/v2/accounts/1/credit_payments", func(w http.ResponseWriter, r *http.Request) {
		invocations++
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("credit_payments.xml"))
	}, t)

	pager := client.CreditPayments.ListAccount("1", nil)
	for pager.Next() {
		if pmts, err := pager.Fetch(context.Background()); err != nil {
			t.Fatal(err)
		} else if !s.Invoked {
			t.Fatal("expected s to be invoked")
		} else if diff := cmp.Diff(pmts, []recurly.CreditPayment{*NewTestCreditPayment()}); diff != "" {
			t.Fatal(diff)
		}
	}
	if invocations != 1 {
		t.Fatalf("unexpected number of invocations: %d", invocations)
	}
}

func TestCreditPayments_Get(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		client, s := NewServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/credit_payments/2cc95aa62517e56d5bec3a48afa1b3b9", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write(MustOpenFile("credit_payment.xml"))
		}, t)

		if pmt, err := client.CreditPayments.Get(context.Background(), "2cc95aa62517e56d5bec3a48afa1b3b9"); err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(pmt, NewTestCreditPayment()); diff != "" {
			t.Fatal(diff)
		} else if !s.Invoked {
			t.Fatal("expected fn invocation")
		}
	})

	// Ensure a 404 returns nil values.
	t.Run("ErrNotFound", func(t *testing.T) {
		client, s := NewServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/credit_payments/2cc95aa62517e56d5bec3a48afa1b3b9", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}, t)

		if pmt, err := client.CreditPayments.Get(context.Background(), "2cc95aa62517e56d5bec3a48afa1b3b9"); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if pmt != nil {
			t.Fatalf("expected nil: %#v", pmt)
		}
	})
}

// Returns a CreditPayment corresponding to testdata/credit_payment.xml.
func NewTestCreditPayment() *recurly.CreditPayment {
	return &recurly.CreditPayment{
		XMLName:                   xml.Name{Local: "credit_payment"},
		UUID:                      "3d3f6754c6df41b9d2a32e43029adc55",
		AccountCode:               "3465345645345",
		Action:                    recurly.CreditPaymentActionRefund,
		Currency:                  "USD",
		AmountInCents:             1000,
		OriginalInvoiceNumber:     1000,
		AppliedToInvoice:          1000,
		OriginalCreditPaymentUUID: "3d3f6754c6df41b9d2a32e43029adc55",
		RefundTransactionUUID:     "3e823e405e7f752988536947c08349ae",
		CreatedAt:                 recurly.NewTime(MustParseTime("2017-07-06T15:51:38Z")),
		UpdatedAt:                 recurly.NewTime(MustParseTime("2017-07-06T15:51:38Z")),
	}
}
