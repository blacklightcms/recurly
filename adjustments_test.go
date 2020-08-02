package recurly_test

import (
	"bytes"
	"context"
	"encoding/xml"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/blacklightcms/recurly"
	"github.com/google/go-cmp/cmp"
)

// Ensure structs are encoded to XML properly.
func TestAdjustments_Encoding(t *testing.T) {
	now := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		v        recurly.Adjustment
		expected string
	}{
		{
			expected: MustCompactString(`
				<adjustment></adjustment>
			`),
		},
		{
			v: recurly.Adjustment{UnitAmountInCents: recurly.NewInt(2000), Currency: "USD"},
			expected: MustCompactString(`
				<adjustment>
					<unit_amount_in_cents>2000</unit_amount_in_cents>
					<currency>USD</currency>
				</adjustment>
			`),
		},
		{
			v: recurly.Adjustment{Origin: "external_gift_card", UnitAmountInCents: recurly.NewInt(-2000), Currency: "USD"},
			expected: MustCompactString(`
				<adjustment>
					<origin>external_gift_card</origin>
					<unit_amount_in_cents>-2000</unit_amount_in_cents>
					<currency>USD</currency>
				</adjustment>
			`),
		},

		{
			v: recurly.Adjustment{Description: "Charge for extra bandwidth", ProductCode: "bandwidth", UnitAmountInCents: recurly.NewInt(2000), Currency: "USD"},
			expected: MustCompactString(`
				<adjustment>
					<description>Charge for extra bandwidth</description>
					<product_code>bandwidth</product_code>
					<unit_amount_in_cents>2000</unit_amount_in_cents>
					<currency>USD</currency>
				</adjustment>
			`),
		},
		{
			v: recurly.Adjustment{Quantity: 1, UnitAmountInCents: recurly.NewInt(2000), Currency: "CAD"},
			expected: MustCompactString(`
				<adjustment>
					<unit_amount_in_cents>2000</unit_amount_in_cents>
					<quantity>1</quantity>
					<currency>CAD</currency>
				</adjustment>
			`),
		},
		{
			v: recurly.Adjustment{AccountingCode: "bandwidth", UnitAmountInCents: recurly.NewInt(2000), Currency: "CAD"},
			expected: MustCompactString(`
				<adjustment>
					<accounting_code>bandwidth</accounting_code>
					<unit_amount_in_cents>2000</unit_amount_in_cents>
					<currency>CAD</currency>
				</adjustment>
			`),
		},
		{
			v: recurly.Adjustment{TaxExempt: recurly.NewBool(false), UnitAmountInCents: recurly.NewInt(2000), Currency: "USD"},
			expected: MustCompactString(`
				<adjustment>
					<unit_amount_in_cents>2000</unit_amount_in_cents>
					<currency>USD</currency>
					<tax_exempt>false</tax_exempt>
				</adjustment>
			`),
		},
		{
			v: recurly.Adjustment{TaxCode: "digital", UnitAmountInCents: recurly.NewInt(2000), Currency: "USD"},
			expected: MustCompactString(`
				<adjustment>
					<unit_amount_in_cents>2000</unit_amount_in_cents>
					<currency>USD</currency>
					<tax_code>digital</tax_code>
				</adjustment>
			`),
		},
		{
			v: recurly.Adjustment{StartDate: recurly.NewTime(now), UnitAmountInCents: recurly.NewInt(2000), Currency: "USD"},
			expected: MustCompactString(`
				<adjustment>
					<unit_amount_in_cents>2000</unit_amount_in_cents>
					<currency>USD</currency>
					<start_date>2000-01-01T00:00:00Z</start_date>
				</adjustment>
			`),
		},
		{
			v: recurly.Adjustment{StartDate: recurly.NewTime(now), EndDate: recurly.NewTime(now), UnitAmountInCents: recurly.NewInt(2000), Currency: "USD"},
			expected: MustCompactString(`
				<adjustment>
					<unit_amount_in_cents>2000</unit_amount_in_cents>
					<currency>USD</currency>
					<start_date>2000-01-01T00:00:00Z</start_date>
					<end_date>2000-01-01T00:00:00Z</end_date>
				</adjustment>
			`),
		},
		{
			v: recurly.Adjustment{UnitAmountInCents: recurly.NewInt(2000), Currency: "USD", RevenueScheduleType: recurly.RevenueScheduleTypeAtInvoice},
			expected: MustCompactString(`
				<adjustment>
					<revenue_schedule_type>at_invoice</revenue_schedule_type>
					<unit_amount_in_cents>2000</unit_amount_in_cents>
					<currency>USD</currency>
				</adjustment>
			`),
		},
		{
			v: recurly.Adjustment{UnitAmountInCents: recurly.NewInt(2000), Currency: "USD", AvalaraServiceType: 600, AvalaraTransactionType: 3},
			expected: MustCompactString(`
				<adjustment>
					<unit_amount_in_cents>2000</unit_amount_in_cents>
					<currency>USD</currency>
					<avalara_transaction_type>3</avalara_transaction_type>
					<avalara_service_type>600</avalara_service_type>
				</adjustment>
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

func TestAdjustments_ListAccount(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	var invocations int
	s.HandleFunc("GET", "/v2/accounts/1/adjustments", func(w http.ResponseWriter, r *http.Request) {
		invocations++
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("adjustments.xml"))
	}, t)

	pager := client.Adjustments.ListAccount("1", nil)
	for pager.Next() {
		var a []recurly.Adjustment
		if err := pager.Fetch(context.Background(), &a); err != nil {
			t.Fatal(err)
		} else if !s.Invoked {
			t.Fatal("expected s to be invoked")
		} else if diff := cmp.Diff(a, []recurly.Adjustment{*NewTestAdjustment()}); diff != "" {
			t.Fatal(diff)
		}
	}
	if invocations != 1 {
		t.Fatalf("unexpected number of invocations: %d", invocations)
	}
}

func TestAdjustments_Get(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/adjustments/626db120a84102b1809909071c701c60", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write(MustOpenFile("adjustment.xml"))
		}, t)

		if a, err := client.Adjustments.Get(context.Background(), "626db120-a841-02b1-8099-09071c701c60"); err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(a, NewTestAdjustment()); diff != "" {
			t.Fatal(diff)
		} else if !s.Invoked {
			t.Fatal("expected fn invocation")
		}
	})

	// Ensure a 404 returns nil values.
	t.Run("ErrNotFound", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/adjustments/626db120a84102b1809909071c701c60", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}, t)

		if a, err := client.Adjustments.Get(context.Background(), "626db120-a841-02b1-8099-09071c701c60"); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if a != nil {
			t.Fatalf("expected nil: %#v", a)
		}
	})
}

func TestAdjustments_Create(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("POST", "/v2/accounts/1/adjustments", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			w.Write(MustOpenFile("adjustment.xml"))
		}, t)

		if a, err := client.Adjustments.Create(context.Background(), "1", recurly.Adjustment{}); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(a, NewTestAdjustment()); diff != "" {
			t.Fatal(diff)
		}
	})

	t.Run("Credit", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("POST", "/v2/accounts/1/adjustments", func(w http.ResponseWriter, r *http.Request) {
			if str := MustReadAllString(r.Body); str != MustCompactString(`
				<adjustment>
					<description>Description</description>
					<unit_amount_in_cents>-100</unit_amount_in_cents>
					<currency>USD</currency>
				</adjustment>
			`) {
				t.Fatal(str)
			}
			w.WriteHeader(http.StatusCreated)
			w.Write(MustOpenFile("adjustment.xml"))
		}, t)

		if a, err := client.Adjustments.Create(context.Background(), "1", recurly.Adjustment{
			UnitAmountInCents: recurly.NewInt(-100),
			Description:       "Description",
			Currency:          "USD",
		}); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(a, NewTestAdjustment()); diff != "" {
			t.Fatal(diff)
		}
	})
}

func TestAdjustments_Delete(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("DELETE", "/v2/adjustments/945a4cb9afd64300b97b138407a51aef", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}, t)

	if err := client.Adjustments.Delete(context.Background(), "945a4cb9-afd6-4300-b97b-138407a51aef"); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	}
}

// Returns an Adjustment corresponding to testdata/adjustment.xml.
func NewTestAdjustment() *recurly.Adjustment {
	return &recurly.Adjustment{
		XMLName:                xml.Name{Local: "adjustment"},
		AccountCode:            "100",
		InvoiceNumber:          1108,
		SubscriptionUUID:       "453f6aa0995e2d52c0d3e6453e9341da",
		UUID:                   "626db120a84102b1809909071c701c60",
		State:                  "invoiced",
		Description:            "One-time Charged Fee",
		ProductCode:            "basic",
		Origin:                 "debit",
		UnitAmountInCents:      recurly.NewInt(2000),
		Quantity:               1,
		OriginalAdjustmentUUID: "2cc95aa62517e56d5bec3a48afa1b3b9",
		TaxInCents:             175,
		TotalInCents:           2175,
		Currency:               "USD",
		TaxType:                "usst",
		TaxRegion:              "CA",
		TaxRate:                0.0875,
		TaxExempt:              recurly.NewBool(false),
		TaxDetails: []recurly.TaxDetail{
			{
				XMLName:    xml.Name{Local: "tax_detail"},
				Name:       "california",
				Type:       "state",
				TaxRate:    0.065,
				TaxInCents: 130,
				Billable:   recurly.NewBool(true),
				Level:      "state",
			},
			{
				XMLName:    xml.Name{Local: "tax_detail"},
				Name:       "san mateo county",
				Type:       "county",
				TaxRate:    0.01,
				TaxInCents: 20,
			},
			{
				XMLName:    xml.Name{Local: "tax_detail"},
				Name:       "sf municipal tax",
				Type:       "city",
				TaxRate:    0.0,
				TaxInCents: 0,
			},
			{
				XMLName:    xml.Name{Local: "tax_detail"},
				Type:       "special",
				TaxRate:    0.0125,
				TaxInCents: 25,
			},
		},
		StartDate: recurly.NewTime(MustParseTime("2015-02-04T23:13:07Z")),
		CreatedAt: recurly.NewTime(MustParseTime("2015-02-04T23:13:07Z")),
	}
}
