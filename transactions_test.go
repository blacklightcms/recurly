package recurly_test

import (
	"context"
	"encoding/xml"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/blacklightcms/recurly"
	"github.com/google/go-cmp/cmp"
)

func TestTransactions_List(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	var invocations int
	s.HandleFunc("GET", "/v2/transactions", func(w http.ResponseWriter, r *http.Request) {
		invocations++
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("transactions.xml"))
	}, t)

	pager := client.Transactions.List(nil)
	for pager.Next() {
		var transactions []recurly.Transaction
		if err := pager.Fetch(context.Background(), &transactions); err != nil {
			t.Fatal(err)
		} else if !s.Invoked {
			t.Fatal("expected s to be invoked")
		} else if diff := cmp.Diff(transactions, []recurly.Transaction{*NewTestTransaction()}); diff != "" {
			t.Fatal(diff)
		}
	}
	if invocations != 1 {
		t.Fatalf("unexpected number of invocations: %d", invocations)
	}
}

func TestTransactions_ListAccount(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	var invocations int
	s.HandleFunc("GET", "/v2/accounts/1/transactions", func(w http.ResponseWriter, r *http.Request) {
		invocations++
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("transactions.xml"))
	}, t)

	pager := client.Transactions.ListAccount("1", nil)
	for pager.Next() {
		var transactions []recurly.Transaction
		if err := pager.Fetch(context.Background(), &transactions); err != nil {
			t.Fatal(err)
		} else if !s.Invoked {
			t.Fatal("expected s to be invoked")
		} else if diff := cmp.Diff(transactions, []recurly.Transaction{*NewTestTransaction()}); diff != "" {
			t.Fatal(diff)
		}
	}
	if invocations != 1 {
		t.Fatalf("unexpected number of invocations: %d", invocations)
	}
}

func TestTransactions_Get(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/transactions/a13acd8fe4294916b79aec87b7ea441f", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write(MustOpenFile("transaction.xml"))
		}, t)

		if transaction, err := client.Transactions.Get(context.Background(), "a13acd8f-e429-4916-b79a-ec87b7ea441f"); err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(transaction, NewTestTransaction()); diff != "" {
			t.Fatal(diff)
		} else if !s.Invoked {
			t.Fatal("expected fn invocation")
		}
	})

	// Retrieving a refunded transaction should succeed
	t.Run("Refund OK", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/transactions/5346552e216a445f82b524bb9d1d27aa", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write(MustOpenFile("transaction_refunded.xml"))
		}, t)

		if transaction, err := client.Transactions.Get(context.Background(), "5346552e-216a-445f-82b5-24bb9d1d27aa"); err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(transaction, NewTestTransactionRefunded()); diff != "" {
			t.Fatal(diff)
		} else if !s.Invoked {
			t.Fatal("expected fn invocation")
		}
	})

	// Retrieving a failed transaction should hold the transaction errors.
	t.Run("TransactionFailed", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/transactions/a13acd8fe4294916b79aec87b7ea441f", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write(MustOpenFile("transaction_failed.xml"))
		}, t)

		if transaction, err := client.Transactions.Get(context.Background(), "a13acd8f-e429-4916-b79a-ec87b7ea441f"); err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(transaction, NewTestTransactionFailed()); diff != "" {
			t.Fatal(diff)
		} else if !s.Invoked {
			t.Fatal("expected fn invocation")
		}
	})

	// Ensure a 404 returns nil values.
	t.Run("ErrNotFound", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/transactions/a13acd8fe4294916b79aec87b7ea441f", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}, t)

		if transaction, err := client.Transactions.Get(context.Background(), "a13acd8f-e429-4916-b79a-ec87b7ea441f"); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if transaction != nil {
			t.Fatalf("expected nil: %#v", transaction)
		}
	})
}

// Returns a Transaction corresponding to testdata/transaction.xml.
func NewTestTransaction() *recurly.Transaction {
	return &recurly.Transaction{
		InvoiceNumber: 1108,
		UUID:          "a13acd8fe4294916b79aec87b7ea441f", // UUID has been sanitized
		Action:        "purchase",
		AmountInCents: 1000,
		TaxInCents:    0,
		Currency:      "USD",
		Status:        "success",
		Description:   "Order #717",
		PaymentMethod: "credit_card",
		Reference:     "5416477",
		Source:        "subscription",
		Recurring:     recurly.NewBool(true),
		Test:          true,
		Voidable:      recurly.NewBool(true),
		Refundable:    recurly.NewBool(true),
		IPAddress:     net.ParseIP("127.0.0.1"),
		CVVResult: recurly.CVVResult{
			Code:    "M",
			Message: "Match",
		},
		AVSResult: recurly.AVSResult{
			Code:    "D",
			Message: "Street address and postal code match.",
		},
		CreatedAt: recurly.NewTime(time.Date(2015, time.June, 10, 15, 25, 6, 0, time.UTC)),
		Account: recurly.Account{
			XMLName:   xml.Name{Local: "account"},
			Code:      "1",
			FirstName: "Verena",
			LastName:  "Example",
			Email:     "verena@test.com",
			BillingInfo: &recurly.Billing{
				XMLName:     xml.Name{Local: "billing_info"},
				FirstName:   "Verena",
				LastName:    "Example",
				Address:     "123 Main St.",
				City:        "San Francisco",
				State:       "CA",
				Zip:         "94105",
				Country:     "US",
				CardType:    "Visa",
				Year:        2017,
				Month:       11,
				FirstSix:    "411111",
				LastFour:    "1111",
				PaymentType: "credit_card",
			},
		},
	}
}

// Returns a Transaction corresponding to testdata/transaction_failed.xml
// as well as the transaction portion of testdata/errors_transaction_failed.xml.
func NewTestTransactionFailed() *recurly.Transaction {
	return &recurly.Transaction{
		UUID:          "3d1c6aa86e3d447eb0f3b4a6e3e074d9",
		Action:        "purchase",
		AmountInCents: 4900,
		TaxInCents:    0,
		Currency:      "USD",
		Status:        "declined",
		Test:          true,
		Voidable:      recurly.NewBool(false),
		Refundable:    recurly.NewBool(false),
		TransactionError: &recurly.TransactionError{
			XMLName:                   xml.Name{Local: "transaction_error"},
			ErrorCode:                 "fraud_security_code",
			ErrorCategory:             "fraud",
			MerchantMessage:           "The payment gateway declined the transaction because the security code (CVV) did not match.",
			CustomerMessage:           "The security code you entered does not match. Please update the CVV and try again.",
			GatewayErrorCode:          "301",
			ThreeDSecureActionTokenID: "ABCDEFGHIJKL012345",
		},
		CVVResult: recurly.CVVResult{
			Code:    "N",
			Message: "No Match",
		},
		AVSResult: recurly.AVSResult{
			Code:    "D",
			Message: "Street address and postal code match.",
		},
		AVSResultStreet: "Y",
		AVSResultPostal: "Y",
		CreatedAt:       recurly.NewTime(MustParseTime("2011-10-17T17:24:53Z")),
		Account: recurly.Account{
			XMLName: xml.Name{Local: "account"},
			Code:    "1",
			Email:   "verena@example.com",
			BillingInfo: &recurly.Billing{
				XMLName:     xml.Name{Local: "billing_info"},
				CardType:    "Visa",
				Year:        2015,
				Month:       11,
				FirstSix:    "400000",
				LastFour:    "0101",
				PaymentType: "credit_card",
			},
		},
	}
}

// Returns a Transaction corresponding to testdata/transaction_refunded.xml
func NewTestTransactionRefunded() *recurly.Transaction {
	return &recurly.Transaction{
		InvoiceNumber:           1108,
		UUID:                    "5346552e216a445f82b524bb9d1d27aa", // UUID has been sanitized
		OriginalTransactionUUID: "b9d02bfaa8bf401abf2b18db76863ac4",
		Action:                  "refund",
		AmountInCents:           1000,
		TaxInCents:              0,
		Currency:                "USD",
		Status:                  "success",
		Description:             "Order #717",
		PaymentMethod:           "credit_card",
		Reference:               "5416477",
		Source:                  "transaction",
		Recurring:               recurly.NewBool(false),
		Test:                    true,
		Voidable:                recurly.NewBool(false),
		Refundable:              recurly.NewBool(false),
		IPAddress:               net.ParseIP("127.0.0.1"),
		CreatedAt:               recurly.NewTime(time.Date(2015, time.June, 10, 15, 25, 6, 0, time.UTC)),
		Account: recurly.Account{
			XMLName:   xml.Name{Local: "account"},
			Code:      "1",
			FirstName: "Verena",
			LastName:  "Example",
			Email:     "verena@test.com",
			BillingInfo: &recurly.Billing{
				XMLName:     xml.Name{Local: "billing_info"},
				FirstName:   "Verena",
				LastName:    "Example",
				Address:     "123 Main St.",
				City:        "San Francisco",
				State:       "CA",
				Zip:         "94105",
				Country:     "US",
				CardType:    "Visa",
				Year:        2017,
				Month:       11,
				FirstSix:    "411111",
				LastFour:    "1111",
				PaymentType: "credit_card",
			},
		},
	}
}
