package recurly_test

import (
	"bytes"
	"context"
	"encoding/xml"
	"net"
	"net/http"
	"strconv"
	"testing"

	"github.com/blacklightcms/recurly"
	"github.com/google/go-cmp/cmp"
)

func TestInvoices_Encoding(t *testing.T) {
	tests := []struct {
		v        recurly.CollectInvoice
		expected string
	}{
		{
			v: recurly.CollectInvoice{
				TransactionType: "moto",
			},
			expected: MustCompactString(`
				<invoice>
					<transaction_type>moto</transaction_type>
				</invoice>
		`),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			buf := new(bytes.Buffer)
			if err := xml.NewEncoder(buf).Encode(tt.v); err != nil {
				t.Fatal(err)
			} else if buf.String() != tt.expected {
				t.Fatal(buf.String())
			}
		})
	}
}

func TestInvoices_List(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	var invocations int
	s.HandleFunc("GET", "/v2/invoices", func(w http.ResponseWriter, r *http.Request) {
		invocations++
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("invoices.xml"))
	}, t)

	pager := client.Invoices.List(nil)
	for pager.Next() {
		var invoices []recurly.Invoice
		if err := pager.Fetch(context.Background(), &invoices); err != nil {
			t.Fatal(err)
		} else if !s.Invoked {
			t.Fatal("expected s to be invoked")
		} else if diff := cmp.Diff(invoices, []recurly.Invoice{*NewTestInvoice()}); diff != "" {
			t.Fatal(diff)
		}
	}
	if invocations != 1 {
		t.Fatalf("unexpected number of invocations: %d", invocations)
	}
}

func TestInvoices_ListAccount(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	var invocations int
	s.HandleFunc("GET", "/v2/accounts/1/invoices", func(w http.ResponseWriter, r *http.Request) {
		invocations++
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("invoices.xml"))
	}, t)

	pager := client.Invoices.ListAccount("1", nil)
	for pager.Next() {
		var invoices []recurly.Invoice
		if err := pager.Fetch(context.Background(), &invoices); err != nil {
			t.Fatal(err)
		} else if !s.Invoked {
			t.Fatal("expected s to be invoked")
		} else if diff := cmp.Diff(invoices, []recurly.Invoice{*NewTestInvoice()}); diff != "" {
			t.Fatal(diff)
		}
	}
	if invocations != 1 {
		t.Fatalf("unexpected number of invocations: %d", invocations)
	}
}

func TestInvoices_Get(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/invoices/5558", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write(MustOpenFile("invoice.xml"))
		}, t)

		if invoice, err := client.Invoices.Get(context.Background(), 5558); err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(invoice, NewTestInvoice()); diff != "" {
			t.Fatal(diff)
		} else if !s.Invoked {
			t.Fatal("expected fn invocation")
		}
	})

	// Ensure a 404 returns nil values.
	t.Run("ErrNotFound", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/invoices/5558", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}, t)

		if invoice, err := client.Invoices.Get(context.Background(), 5558); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if invoice != nil {
			t.Fatalf("expected nil: %#v", invoice)
		}
	})
}

func TestInvoices_GetPDF(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/invoices/5558", func(w http.ResponseWriter, r *http.Request) {
			if h := r.Header.Get("Accept"); h != "application/pdf" {
				t.Fatalf("unexpected 'Accept' header: %q", h)
			} else if h := r.Header.Get("Accept-Language"); h != "English" {
				t.Fatalf("unexpected 'Accept-Language' header: %q", h)
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("binary pdf text"))
		}, t)

		if b, err := client.Invoices.GetPDF(context.Background(), 5558, "English"); err != nil {
			t.Fatal(err)
		} else if b.String() != "binary pdf text" {
			t.Fatal(b.String())
		} else if !s.Invoked {
			t.Fatal("expected fn invocation")
		}
	})

	// Ensure a 404 returns nil values.
	t.Run("ErrNotFound", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/invoices/5558", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}, t)

		if b, err := client.Invoices.GetPDF(context.Background(), 5558, "English"); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if b != nil {
			t.Fatalf("expected nil: %#v", b)
		}
	})
}

func TestInvoices_Preview(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("POST", "/v2/accounts/1/invoices/preview", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write(MustOpenFile("invoice_collection.xml"))
	}, t)

	if invoice, err := client.Invoices.Preview(context.Background(), "1"); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(invoice, NewTestInvoice()); diff != "" {
		t.Fatal(diff)
	}
}

func TestInvoices_Create(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("POST", "/v2/accounts/1/invoices", func(w http.ResponseWriter, r *http.Request) {
		if str := MustReadAllString(r.Body); str != MustCompactString(`
			<invoice>
				<po_number>ABC</po_number>
				<net_terms>30</net_terms>
				<collection_method>COLLECTION_METHOD</collection_method>
				<terms_and_conditions>TERMS</terms_and_conditions>
				<customer_notes>CUSTOMER_NOTES</customer_notes>
				<vat_reverse_charge_notes>VAT_REVERSE_CHARGE_NOTES</vat_reverse_charge_notes>
			</invoice>
		`) {
			t.Fatal(str)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(MustOpenFile("invoice_collection.xml"))
	}, t)

	if invoice, err := client.Invoices.Create(context.Background(), "1", recurly.Invoice{
		PONumber:              "ABC",
		NetTerms:              recurly.NewInt(30),
		CollectionMethod:      "COLLECTION_METHOD",
		TermsAndConditions:    "TERMS",
		CustomerNotes:         "CUSTOMER_NOTES",
		VatReverseChargeNotes: "VAT_REVERSE_CHARGE_NOTES",
	}); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(invoice, NewTestInvoice()); diff != "" {
		t.Fatal(diff)
	}
}

func TestInvoices_Collect(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("PUT", "/v2/invoices/1010/collect", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("invoice.xml"))
	}, t)

	if invoice, err := client.Invoices.Collect(context.Background(), 1010, recurly.CollectInvoice{}); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(invoice, NewTestInvoice()); diff != "" {
		t.Fatal(diff)
	}
}

func TestInvoices_MarkPaid(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("PUT", "/v2/invoices/1010/mark_successful", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("invoice.xml"))
	}, t)

	if invoice, err := client.Invoices.MarkPaid(context.Background(), 1010); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(invoice, NewTestInvoice()); diff != "" {
		t.Fatal(diff)
	}
}

func TestInvoices_MarkFailed(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("PUT", "/v2/invoices/1010/mark_failed", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("invoice_collection.xml"))
	}, t)

	if invoice, err := client.Invoices.MarkFailed(context.Background(), 1010); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(invoice, NewTestInvoice()); diff != "" {
		t.Fatal(diff)
	}
}

func TestInvoices_RefundVoidOpenAmount(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("POST", "/v2/invoices/1010/refund", func(w http.ResponseWriter, r *http.Request) {
		if str := MustReadAllString(r.Body); str != MustCompactString(`
			<invoice>
				<amount_in_cents>1000</amount_in_cents>
				<refund_method>credit_first</refund_method>
				<external_refund>true</external_refund>
				<credit_customer_notes>notes</credit_customer_notes>
				<payment_method>METHOD</payment_method>
				<description>description</description>
				<refunded_at>2011-04-10T07:00:00Z</refunded_at>
			</invoice>
		`) {
			t.Fatal(str)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(MustOpenFile("invoice.xml"))
	}, t)

	if invoice, err := client.Invoices.RefundVoidOpenAmount(context.Background(), 1010, recurly.InvoiceRefund{
		AmountInCents:       recurly.NewInt(1000),
		RefundMethod:        "credit_first",
		ExternalRefund:      recurly.NewBool(true),
		CreditCustomerNotes: "notes",
		PaymentMethod:       "METHOD",
		Description:         "description",
		RefundedAt:          recurly.NewTime(MustParseTime("2011-04-10T07:00:00Z")),
	}); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(invoice, NewTestInvoice()); diff != "" {
		t.Fatal(diff)
	}
}

func TestInvoices_RefundVoidLineItems(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("POST", "/v2/invoices/1010/refund", func(w http.ResponseWriter, r *http.Request) {
		if str := MustReadAllString(r.Body); str != MustCompactString(`
			<invoice>
				<line_items>
					<adjustment>
						<uuid>2bc33a7469dc1458f455634212acdcd6</uuid>
						<quantity>1</quantity>
						<prorate>false</prorate>
					</adjustment>
					<adjustment>
						<uuid>2bc33a746a89d867df47024fd6b261b6</uuid>
						<quantity>1</quantity>
						<prorate>true</prorate>
					</adjustment>
				</line_items>
				<amount_in_cents>1000</amount_in_cents>
				<refund_method>credit_first</refund_method>
				<external_refund>true</external_refund>
				<credit_customer_notes>notes</credit_customer_notes>
				<payment_method>METHOD</payment_method>
				<description>description</description>
				<refunded_at>2011-04-10T07:00:00Z</refunded_at>
			</invoice>
		`) {
			t.Fatal(str)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(MustOpenFile("invoice.xml"))
	}, t)

	if invoice, err := client.Invoices.RefundVoidLineItems(context.Background(), 1010, recurly.InvoiceLineItemsRefund{
		LineItems: []recurly.VoidLineItem{
			{
				UUID:     "2bc33a7469dc1458f455634212acdcd6",
				Quantity: 1,
				Prorate:  recurly.NewBool(false),
			},
			{
				UUID:     "2bc33a746a89d867df47024fd6b261b6",
				Quantity: 1,
				Prorate:  recurly.NewBool(true),
			},
		},
		InvoiceRefund: recurly.InvoiceRefund{
			AmountInCents:       recurly.NewInt(1000),
			RefundMethod:        "credit_first",
			ExternalRefund:      recurly.NewBool(true),
			CreditCustomerNotes: "notes",
			PaymentMethod:       "METHOD",
			Description:         "description",
			RefundedAt:          recurly.NewTime(MustParseTime("2011-04-10T07:00:00Z")),
		},
	}); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(invoice, NewTestInvoice()); diff != "" {
		t.Fatal(diff)
	}
}

func TestInvoices_VoidCreditInvoice(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("PUT", "/v2/invoices/1010/void", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("invoice.xml"))
	}, t)

	if invoice, err := client.Invoices.VoidCreditInvoice(context.Background(), 1010); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(invoice, NewTestInvoice()); diff != "" {
		t.Fatal(diff)
	}
}

func TestInvoices_RecordPayment(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("POST", "/v2/invoices/1010/transactions", func(w http.ResponseWriter, r *http.Request) {
		if str := MustReadAllString(r.Body); str != MustCompactString(`
			<transaction>
				<payment_method>check</payment_method>
				<collected_at>2017-01-03T00:00:00Z</collected_at>
				<amount_in_cents>1000</amount_in_cents>
				<description>Paid with a check</description>
			</transaction>
		`) {
			t.Fatal(str)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("transaction.xml"))
	}, t)

	if transaction, err := client.Invoices.RecordPayment(context.Background(), recurly.OfflinePayment{
		InvoiceNumber: 1010,
		PaymentMethod: recurly.PaymentMethodCheck,
		Amount:        1000,
		CollectedAt:   recurly.NewTime(MustParseTime("2017-01-03T00:00:00Z")),
		Description:   "Paid with a check",
	}); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(transaction, NewTestTransaction()); diff != "" {
		t.Fatal(diff)
	}
}

// Returns a Invoice corresponding to testdata/invoice.xml.
func NewTestInvoice() *recurly.Invoice {
	return &recurly.Invoice{
		XMLName:     xml.Name{Local: "invoice"},
		AccountCode: "1",
		Address: recurly.Address{
			XMLName: xml.Name{Local: "address"},
			Address: "400 Alabama St.",
			City:    "San Francisco",
			State:   "CA",
			Zip:     "94110",
			Country: "US",
		},
		UUID:             "421f7b7d414e4c6792938e7c49d552e9",
		State:            recurly.ChargeInvoiceStatePaid,
		InvoiceNumber:    5558,
		SubtotalInCents:  153000,
		TaxInCents:       0,
		TotalInCents:     153000,
		Currency:         "USD",
		CreatedAt:        recurly.NewTime(MustParseTime("2018-06-05T15:44:57Z")),
		UpdatedAt:        recurly.NewTime(MustParseTime("2018-06-05T15:44:57Z")),
		ClosedAt:         recurly.NewTime(MustParseTime("2018-06-05T15:44:57Z")),
		DueOn:            recurly.NewTime(MustParseTime("2018-06-05T15:44:57Z")),
		Type:             "charge",
		Origin:           "purchase",
		TaxRate:          float64(0),
		NetTerms:         recurly.NewInt(0),
		CollectionMethod: "automatic",
		TaxDetails: &[]recurly.TaxDetail{
			{
				XMLName:    xml.Name{Local: "tax_detail"},
				Name:       "california",
				Type:       "state",
				TaxRate:    0.065,
				TaxInCents: 130,
				Billable:   recurly.NewBool(true),
				Level:      "state",
			},
		},
		LineItems: []recurly.Adjustment{
			{
				XMLName:             xml.Name{Local: "adjustment"},
				AccountCode:         "1",
				InvoiceNumber:       5558,
				SubscriptionUUID:    "453f6aa0995e2d52c0d3e6453e9341da",
				UUID:                "626db120a84102b1809909071c701c60",
				State:               "invoiced",
				Description:         "License",
				RevenueScheduleType: recurly.RevenueScheduleTypeEvenly,
				ProductCode:         "license",
				Origin:              "add_on",
				UnitAmountInCents:   recurly.NewInt(150000),
				Quantity:            1,
				TaxInCents:          0,
				TotalInCents:        150000,
				Currency:            "USD",
				StartDate:           recurly.NewTime(MustParseTime("2018-06-05T15:44:56Z")),
				EndDate:             recurly.NewTime(MustParseTime("2018-07-05T15:44:56Z")),
				CreatedAt:           recurly.NewTime(MustParseTime("2018-06-05T15:44:57Z")),
				UpdatedAt:           recurly.NewTime(MustParseTime("2018-06-05T15:44:57Z")),
			},
			{
				XMLName:             xml.Name{Local: "adjustment"},
				AccountCode:         "1",
				InvoiceNumber:       5558,
				SubscriptionUUID:    "453f6aa0995e2d52c0d3e6453e9341da",
				UUID:                "453f6aa1473a0620e4411a4fc88122cf",
				State:               "invoiced",
				Description:         "Domains",
				RevenueScheduleType: recurly.RevenueScheduleTypeEvenly,
				ProductCode:         "domains",
				Origin:              "add_on",
				UnitAmountInCents:   recurly.NewInt(1500),
				Quantity:            2,
				TaxInCents:          0,
				TotalInCents:        3000,
				Currency:            "USD",
				StartDate:           recurly.NewTime(MustParseTime("2018-06-05T15:44:56Z")),
				EndDate:             recurly.NewTime(MustParseTime("2018-07-05T15:44:56Z")),
				CreatedAt:           recurly.NewTime(MustParseTime("2018-06-05T15:44:57Z")),
				UpdatedAt:           recurly.NewTime(MustParseTime("2018-06-05T15:44:57Z")),
			},
		},
		Transactions: []recurly.Transaction{
			{
				InvoiceNumber: 5558,
				UUID:          "453f6aa17014a4624b636e455c8876ca",
				Action:        "purchase",
				AmountInCents: 153000,
				TaxInCents:    0,
				Currency:      "USD",
				Status:        "success",
				PaymentMethod: "credit_card",
				Reference:     "1909673",
				Source:        "subscription",
				Recurring:     recurly.NewBool(false),
				Test:          true,
				Voidable:      recurly.NewBool(true),
				Refundable:    recurly.NewBool(true),
				IPAddress:     net.ParseIP("127.0.0.1"),
				GatewayType:   "test",
				Origin:        "token_api",
				Message:       "Successful test transaction",
				AVSResult: recurly.AVSResult{
					Code:    "D",
					Message: "Street address and postal code match.",
				},
				CreatedAt: recurly.NewTime(MustParseTime("2018-06-05T15:44:56Z")),
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
			},
		},
	}
}

// Returns an InvoiceCollection corresponding to testdata/invoice_collection.xml.
func NewTestInvoiceCollection() *recurly.InvoiceCollection {
	return &recurly.InvoiceCollection{
		XMLName:        xml.Name{Local: "invoice_collection"},
		ChargeInvoice:  NewTestInvoice(),
		CreditInvoices: []recurly.Invoice{},
	}
}
