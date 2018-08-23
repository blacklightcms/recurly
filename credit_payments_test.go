package recurly_test

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"testing"

	"github.com/launchpadcentral/recurly"
	"github.com/google/go-cmp/cmp"
)

func TestCreditPayments_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/credit_payments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
		<credit_payments type="array">
			<credit_payment href="https://api.recurly.com/v2/credit_payments/3d3f6754c6df41b9d2a32e43029adc55" type="payment">
				<account href="https://api.recurly.com/v2/accounts/3465345645345"/>
				<uuid>3d3f6754c6df41b9d2a32e43029adc55</uuid>
				<action>payment</action>
				<currency>USD</currency>
				<amount_in_cents type="integer">1000</amount_in_cents>
				<original_invoice href="https://api.recurly.com/v2/invoices/1000"/>
				<applied_to_invoice href="https://api.recurly.com/v2/invoices/1001"/>
				<created_at type="datetime">2017-07-06T15:51:38Z</created_at>
				<updated_at type="datetime">2017-07-06T15:51:38Z</updated_at>
				<voided_at type="datetime">2017-07-22T15:51:38Z</voided_at>
			</credit_payment>
			<credit_payment href="https://api.recurly.com/v2/credit_payments/3e8764beb04add789fb3b54778838e17" type="charge">
				<account href="https://api.recurly.com/v2/accounts/3465345645345"/>
				<uuid>3d3f6754c6df41b9d2a32e43029adc55</uuid>
				<action>refund</action>
				<currency>USD</currency>
				<amount_in_cents>1000</amount_in_cents>
				<original_invoice href="https://api.recurly.com/v2/invoices/1001"/>
				<applied_to_invoice href="https://api.recurly.com/v2/invoices/1002"/>
				<original_credit_payment href="https://api.recurly.com/v2/credit_payments/3d3f6754c6df41b9d2a32e43029adc55"/>
				<refund_transaction href="https://api.recurly.com/v2/transactions/3e823e405e7f752988536947c08349ae"/>
				<created_at type="datetime">2017-07-06T15:51:38Z</created_at>
				<updated_at type="datetime">2017-07-06T15:51:38Z</updated_at>
				<voided_at nil="nil"></voided_at>
			</credit_payment>
		</credit_payments>`)
	})

	resp, creditPayments, err := client.CreditPayments.List(recurly.Params{"per_page": 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected list credit payments to return OK")
	} else if resp.Request.URL.Query().Get("per_page") != "1" {
		t.Fatalf("expected per_page parameter of 1, given %s", resp.Request.URL.Query().Get("per_page"))
	} else if diff := cmp.Diff(creditPayments, []recurly.CreditPayment{
		{
			XMLName:               xml.Name{Local: "credit_payment"},
			UUID:                  "3d3f6754c6df41b9d2a32e43029adc55",
			AccountCode:           "3465345645345",
			Action:                recurly.CreditPaymentActionPayment,
			Currency:              "USD",
			AmountInCents:         1000,
			OriginalInvoiceNumber: 1000,
			AppliedToInvoice:      1001,
			CreatedAt:             recurly.NewTimeFromString("2017-07-06T15:51:38Z"),
			UpdatedAt:             recurly.NewTimeFromString("2017-07-06T15:51:38Z"),
			VoidedAt:              recurly.NewTimeFromString("2017-07-22T15:51:38Z"),
		},
		{
			XMLName:                   xml.Name{Local: "credit_payment"},
			UUID:                      "3d3f6754c6df41b9d2a32e43029adc55",
			AccountCode:               "3465345645345",
			Action:                    recurly.CreditPaymentActionRefund,
			Currency:                  "USD",
			AmountInCents:             1000,
			OriginalInvoiceNumber:     1001,
			AppliedToInvoice:          1002,
			OriginalCreditPaymentUUID: "3d3f6754c6df41b9d2a32e43029adc55",
			RefundTransactionUUID:     "3e823e405e7f752988536947c08349ae",
			CreatedAt:                 recurly.NewTimeFromString("2017-07-06T15:51:38Z"),
			UpdatedAt:                 recurly.NewTimeFromString("2017-07-06T15:51:38Z"),
		},
	}); diff != "" {
		t.Fatalf(diff)
	}
}

func TestCreditPayments_ListAccount(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/1/credit_payments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
		<credit_payments type="array">
			<credit_payment href="https://api.recurly.com/v2/credit_payments/3d3f6754c6df41b9d2a32e43029adc55" type="payment">
				<account href="https://api.recurly.com/v2/accounts/3465345645345"/>
				<uuid>3d3f6754c6df41b9d2a32e43029adc55</uuid>
				<action>payment</action>
				<currency>USD</currency>
				<amount_in_cents type="integer">1000</amount_in_cents>
				<original_invoice href="https://api.recurly.com/v2/invoices/1000"/>
				<applied_to_invoice href="https://api.recurly.com/v2/invoices/1001"/>
				<created_at type="datetime">2017-07-06T15:51:38Z</created_at>
				<updated_at type="datetime">2017-07-06T15:51:38Z</updated_at>
				<voided_at type="datetime">2017-07-22T15:51:38Z</voided_at>
			</credit_payment>
			<credit_payment href="https://api.recurly.com/v2/credit_payments/3e8764beb04add789fb3b54778838e17" type="charge">
				<account href="https://api.recurly.com/v2/accounts/3465345645345"/>
				<uuid>3d3f6754c6df41b9d2a32e43029adc55</uuid>
				<action>refund</action>
				<currency>USD</currency>
				<amount_in_cents>1000</amount_in_cents>
				<original_invoice href="https://api.recurly.com/v2/invoices/1000"/>
				<applied_to_invoice href="https://api.recurly.com/v2/invoices/1000"/>
				<original_credit_payment href="https://api.recurly.com/v2/credit_payments/3d3f6754c6df41b9d2a32e43029adc55"/>
				<refund_transaction href="https://api.recurly.com/v2/transactions/3e823e405e7f752988536947c08349ae"/>
				<created_at type="datetime">2017-07-06T15:51:38Z</created_at>
				<updated_at type="datetime">2017-07-06T15:51:38Z</updated_at>
				<voided_at nil="nil"></voided_at>
			</credit_payment>
		</credit_payments>`)
	})

	resp, creditPayments, err := client.CreditPayments.ListAccount("1", recurly.Params{"per_page": 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected list credit payments to return OK")
	} else if pp := resp.Request.URL.Query().Get("per_page"); pp != "1" {
		t.Fatalf("unexpected per_page: %s", pp)
	} else if diff := cmp.Diff(creditPayments, []recurly.CreditPayment{
		{
			XMLName:               xml.Name{Local: "credit_payment"},
			UUID:                  "3d3f6754c6df41b9d2a32e43029adc55",
			AccountCode:           "3465345645345",
			Action:                recurly.CreditPaymentActionPayment,
			Currency:              "USD",
			AmountInCents:         1000,
			OriginalInvoiceNumber: 1000,
			AppliedToInvoice:      1001,
			CreatedAt:             recurly.NewTimeFromString("2017-07-06T15:51:38Z"),
			UpdatedAt:             recurly.NewTimeFromString("2017-07-06T15:51:38Z"),
			VoidedAt:              recurly.NewTimeFromString("2017-07-22T15:51:38Z"),
		},
		{
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
			CreatedAt:                 recurly.NewTimeFromString("2017-07-06T15:51:38Z"),
			UpdatedAt:                 recurly.NewTimeFromString("2017-07-06T15:51:38Z"),
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestCreditPayments_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/credit_payments/2cc95aa62517e56d5bec3a48afa1b3b9", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
		<credit_payment href="https://api.recurly.com/v2/credit_payments/3e8764beb04add789fb3b54778838e17" type="charge">
			<account href="https://api.recurly.com/v2/accounts/3465345645345"/>
			<uuid>3d3f6754c6df41b9d2a32e43029adc55</uuid>
			<action>refund</action>
			<currency>USD</currency>
			<amount_in_cents>1000</amount_in_cents>
			<original_invoice href="https://api.recurly.com/v2/invoices/1000"/>
			<applied_to_invoice href="https://api.recurly.com/v2/invoices/1000"/>
			<original_credit_payment href="https://api.recurly.com/v2/credit_payments/3d3f6754c6df41b9d2a32e43029adc55"/>
			<refund_transaction href="https://api.recurly.com/v2/transactions/3e823e405e7f752988536947c08349ae"/>
			<created_at type="datetime">2017-07-06T15:51:38Z</created_at>
			<updated_at type="datetime">2017-07-06T15:51:38Z</updated_at>
			<voided_at nil="nil"></voided_at>
		</credit_payment>`)
	})

	resp, creditPayment, err := client.CreditPayments.Get("2cc95aa62517e56d5bec3a48afa1b3b9")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected get invoice to return OK")
	}

	if diff := cmp.Diff(creditPayment, &recurly.CreditPayment{
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
		CreatedAt:                 recurly.NewTimeFromString("2017-07-06T15:51:38Z"),
		UpdatedAt:                 recurly.NewTimeFromString("2017-07-06T15:51:38Z"),
	}); diff != "" {
		t.Fatal(diff)
	}
}
