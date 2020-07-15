package recurly_test

import (
	"bytes"
	"context"
	"encoding/xml"
	"net/http"
	"strconv"
	"testing"

	"github.com/blacklightcms/recurly"
	"github.com/google/go-cmp/cmp"
)

// Ensure structs are encoded to XML properly.
func TestPlans_Encoding(t *testing.T) {
	tests := []struct {
		v        recurly.Plan
		expected string
	}{
		// name is a required field. It should always be present.
		{
			expected: MustCompactString(`
				<plan>
					<name></name>
				</plan>
			`),
		},
		{
			v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, Description: "abc"},
			expected: MustCompactString(`
				<plan>
					<name>Gold plan</name>
					<description>abc</description>
					<unit_amount_in_cents>
					<USD>1500</USD>
					</unit_amount_in_cents>
				</plan>
			`),
		},
		{
			v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, AccountingCode: "gold"},
			expected: MustCompactString(`
				<plan>
					<name>Gold plan</name>
					<accounting_code>gold</accounting_code>
					<unit_amount_in_cents>
					<USD>1500</USD>
					</unit_amount_in_cents>
				</plan>
			`),
		},
		{
			v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, IntervalUnit: "months"},
			expected: MustCompactString(`
				<plan>
					<name>Gold plan</name>
					<plan_interval_unit>months</plan_interval_unit>
					<unit_amount_in_cents>
					<USD>1500</USD>
					</unit_amount_in_cents>
				</plan>
			`),
		},
		{
			v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, IntervalLength: 1},
			expected: MustCompactString(`
				<plan>
					<name>Gold plan</name>
					<plan_interval_length>1</plan_interval_length>
					<unit_amount_in_cents>
					<USD>1500</USD>
					</unit_amount_in_cents>
				</plan>
			`),
		},
		{
			v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, TrialIntervalUnit: "days"},
			expected: MustCompactString(`
				<plan>
					<name>Gold plan</name>
					<trial_interval_unit>days</trial_interval_unit>
					<unit_amount_in_cents>
					<USD>1500</USD>
					</unit_amount_in_cents>
				</plan>
			`),
		},
		{
			v: recurly.Plan{Name: "Gold plan", AutoRenew: true, UnitAmountInCents: recurly.UnitAmount{USD: 1500}, TrialIntervalLength: 10},
			expected: MustCompactString(`
				<plan>
					<name>Gold plan</name>
					<trial_interval_length>10</trial_interval_length>
					<auto_renew>true</auto_renew>
					<unit_amount_in_cents>
					<USD>1500</USD>
					</unit_amount_in_cents>
				</plan>
			`),
		},
		{
			v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, IntervalUnit: "months"},
			expected: MustCompactString(`
				<plan>
					<name>Gold plan</name>
					<plan_interval_unit>months</plan_interval_unit>
					<unit_amount_in_cents>
					<USD>1500</USD>
					</unit_amount_in_cents>
				</plan>
			`),
		},
		{
			v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, SetupFeeInCents: recurly.UnitAmount{USD: 1000, EUR: 800}},
			expected: MustCompactString(`
				<plan>
					<name>Gold plan</name>
					<unit_amount_in_cents>
					<USD>1500</USD>
					</unit_amount_in_cents>
					<setup_fee_in_cents>
					<USD>1000</USD>
					<EUR>800</EUR>
					</setup_fee_in_cents>
				</plan>
			`),
		},
		{
			v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, TotalBillingCycles: recurly.NewInt(24)},
			expected: MustCompactString(`
				<plan>
					<name>Gold plan</name>
					<total_billing_cycles>24</total_billing_cycles>
					<unit_amount_in_cents>
					<USD>1500</USD>
					</unit_amount_in_cents>
				</plan>
			`),
		},
		{
			v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, UnitName: "unit"},
			expected: MustCompactString(`
				<plan>
					<name>Gold plan</name>
					<unit_name>unit</unit_name>
					<unit_amount_in_cents>
					<USD>1500</USD>
					</unit_amount_in_cents>
				</plan>
			`),
		},
		{
			v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, DisplayQuantity: recurly.NewBool(true)},
			expected: MustCompactString(`
				<plan>
					<name>Gold plan</name>
					<display_quantity>true</display_quantity>
					<unit_amount_in_cents>
					<USD>1500</USD>
					</unit_amount_in_cents>
				</plan>
			`),
		},
		{
			v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, DisplayQuantity: recurly.NewBool(false)},
			expected: MustCompactString(`
				<plan>
					<name>Gold plan</name>
					<display_quantity>false</display_quantity>
					<unit_amount_in_cents>
					<USD>1500</USD>
					</unit_amount_in_cents>
				</plan>
			`),
		},
		{
			v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, SuccessURL: "https://example.com/success"},
			expected: MustCompactString(`
				<plan>
					<name>Gold plan</name>
					<success_url>https://example.com/success</success_url>
					<unit_amount_in_cents>
					<USD>1500</USD>
					</unit_amount_in_cents>
				</plan>
			`),
		},
		{
			v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, CancelURL: "https://example.com/cancel"},
			expected: MustCompactString(`
				<plan>
					<name>Gold plan</name>
					<cancel_url>https://example.com/cancel</cancel_url>
					<unit_amount_in_cents>
					<USD>1500</USD>
					</unit_amount_in_cents>
				</plan>
			`),
		},
		{
			v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, TaxExempt: recurly.NewBool(true)},
			expected: MustCompactString(`
				<plan>
					<name>Gold plan</name>
					<tax_exempt>true</tax_exempt>
					<unit_amount_in_cents>
					<USD>1500</USD>
					</unit_amount_in_cents>
				</plan>
			`),
		},
		{
			v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, TaxExempt: recurly.NewBool(false)},
			expected: MustCompactString(`
				<plan>
					<name>Gold plan</name>
					<tax_exempt>false</tax_exempt>
					<unit_amount_in_cents>
					<USD>1500</USD>
					</unit_amount_in_cents>
				</plan>
			`),
		},
		{
			v: recurly.Plan{Name: "Gold plan", UnitAmountInCents: recurly.UnitAmount{USD: 1500}, TaxCode: "physical"},
			expected: MustCompactString(`
				<plan>
					<name>Gold plan</name>
					<tax_code>physical</tax_code>
					<unit_amount_in_cents>
					<USD>1500</USD>
					</unit_amount_in_cents>
				</plan>
			`),
		},
		{
			v: recurly.Plan{Name: "Gold", AllowAnyItemOnSubscription: recurly.NewBool(true)},
			expected: MustCompactString(`
				<plan>
					<name>Gold</name>
					<allow_any_item_on_subscription>true</allow_any_item_on_subscription>
				</plan>
			`),
		},
		{
			v: recurly.Plan{Name: "Gold", AvalaraServiceType: 600, AvalaraTransactionType: 3},
			expected: MustCompactString(`
				<plan>
					<name>Gold</name>
					<avalara_transaction_type>3</avalara_transaction_type>
					<avalara_service_type>600</avalara_service_type>
				</plan>
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

func TestPlans_List(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	var invocations int
	s.HandleFunc("GET", "/v2/plans", func(w http.ResponseWriter, r *http.Request) {
		invocations++
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("plans.xml"))
	}, t)

	pager := client.Plans.List(nil)
	for pager.Next() {
		var plans []recurly.Plan
		if err := pager.Fetch(context.Background(), &plans); err != nil {
			t.Fatal(err)
		} else if !s.Invoked {
			t.Fatal("expected s to be invoked")
		} else if diff := cmp.Diff(plans, []recurly.Plan{*NewTestPlan()}); diff != "" {
			t.Fatal(diff)
		}
	}
	if invocations != 1 {
		t.Fatalf("unexpected number of invocations: %d", invocations)
	}
}

func TestPlans_Get(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/plans/gold", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write(MustOpenFile("plan.xml"))
		}, t)

		if plan, err := client.Plans.Get(context.Background(), "gold"); err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(plan, NewTestPlan()); diff != "" {
			t.Fatal(diff)
		} else if !s.Invoked {
			t.Fatal("expected fn invocation")
		}
	})

	// Ensure a 404 returns nil values.
	t.Run("ErrNotFound", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/plans/gold", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}, t)

		if plan, err := client.Plans.Get(context.Background(), "gold"); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if plan != nil {
			t.Fatalf("expected nil: %#v", plan)
		}
	})
}

func TestPlans_Create(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("POST", "/v2/plans", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write(MustOpenFile("plan.xml"))
	}, t)

	if plan, err := client.Plans.Create(context.Background(), recurly.Plan{}); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(plan, NewTestPlan()); diff != "" {
		t.Fatal(diff)
	}
}

func TestPlans_Update(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("PUT", "/v2/plans/gold", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("plan.xml"))
	}, t)

	if plan, err := client.Plans.Update(context.Background(), "gold", recurly.Plan{}); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(plan, NewTestPlan()); diff != "" {
		t.Fatal(diff)
	}
}

func TestPlans_Delete(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("DELETE", "/v2/plans/gold", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}, t)

	if err := client.Plans.Delete(context.Background(), "gold"); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	}
}

func NewTestPlan() *recurly.Plan {
	return &recurly.Plan{
		XMLName:                  xml.Name{Local: "plan"},
		Code:                     "gold",
		Name:                     "Gold plan",
		DisplayDonationAmounts:   recurly.NewBool(false),
		DisplayQuantity:          recurly.NewBool(false),
		DisplayPhoneNumber:       recurly.NewBool(false),
		BypassHostedConfirmation: recurly.NewBool(false),
		UnitName:                 "unit",
		IntervalUnit:             "months",
		IntervalLength:           1,
		TrialIntervalUnit:        "days",
		TaxExempt:                recurly.NewBool(false),
		UnitAmountInCents: recurly.UnitAmount{
			USD: 6000,
			EUR: 4500,
		},
		SetupFeeInCents: recurly.UnitAmount{
			USD: 1000,
			EUR: 800,
		},
		CreatedAt: recurly.NewTime(MustParseTime("2015-05-29T17:38:15Z")),
	}
}
