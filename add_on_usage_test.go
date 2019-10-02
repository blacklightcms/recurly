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
func TestAddOnUsage_Encoding(t *testing.T) {

	now := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		v        recurly.AddOnUsage
		expected string
	}{
		{
			expected: MustCompactString(`
				<usage>
				</usage>
			`),
		},
		{
			v: recurly.AddOnUsage{Id: 123456},
			expected: MustCompactString(`
				<usage>
					<id>123456</id>
				</usage>
			`),
		},
		{
			v: recurly.AddOnUsage{Amount: 100},
			expected: MustCompactString(`
				<usage>
					<amount>100</amount>
				</usage>
			`),
		},
		{
			v: recurly.AddOnUsage{MerchantTag: "some_merchant"},
			expected: MustCompactString(`
				<usage>
					<merchant_tag>some_merchant</merchant_tag>
				</usage>
			`),
		},
		{
			v: recurly.AddOnUsage{RecordingTimestamp: recurly.NewTime(now)},
			expected: MustCompactString(`
				<usage>
					<recording_timestamp>2000-01-01T00:00:00Z</recording_timestamp>
				</usage>
			`),
		},
		{
			v: recurly.AddOnUsage{UsageTimestamp: recurly.NewTime(now)},
			expected: MustCompactString(`
				<usage>
					<usage_timestamp>2000-01-01T00:00:00Z</usage_timestamp>
				</usage>
			`),
		},
		{
			v: recurly.AddOnUsage{CreatedAt: recurly.NewTime(now)},
			expected: MustCompactString(`
				<usage>
					<created_at>2000-01-01T00:00:00Z</created_at>
				</usage>
			`),
		},
		{
			v: recurly.AddOnUsage{UpdatedAt: recurly.NewTime(now)},
			expected: MustCompactString(`
				<usage>
					<updated_at>2000-01-01T00:00:00Z</updated_at>
				</usage>
			`),
		},
		{
			v: recurly.AddOnUsage{BilledAt: recurly.NewTime(now)},
			expected: MustCompactString(`
				<usage>
					<billed_at>2000-01-01T00:00:00Z</billed_at>
				</usage>
			`),
		},
		{
			v: recurly.AddOnUsage{UsageType: "price"},
			expected: MustCompactString(`
				<usage>
					<usage_type>price</usage_type>
				</usage>
			`),
		},

		{
			v: recurly.AddOnUsage{UnitAmountInCents: 313},
			expected: MustCompactString(`
				<usage>
					<unit_amount_in_cents>313</unit_amount_in_cents>
				</usage>
			`),
		},
		{
			v: recurly.AddOnUsage{UsagePercentage: recurly.NewFloat(0.50)},
			expected: MustCompactString(`
				<usage>
					<usage_percentage>0.5</usage_percentage>
				</usage>
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

func TestAddOnUsage_List(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	var invocations int
	s.HandleFunc("GET", "/v2/subscriptions/1122334455/add_ons/addOnCode/usage", func(w http.ResponseWriter, r *http.Request) {
		invocations++
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("add_on_usages.xml"))
	}, t)

	pager := client.AddOnUsages.List("1122334455", "addOnCode", nil)
	for pager.Next() {
		var a []recurly.AddOnUsage
		if err := pager.Fetch(context.Background(), &a); err != nil {
			t.Fatal(err)
		} else if !s.Invoked {
			t.Fatal("expected s to be invoked")
		} else if diff := cmp.Diff(a, []recurly.AddOnUsage{*NewTestAddOnUsage()}); diff != "" {
			t.Fatal(diff)
		}
	}
	if invocations != 1 {
		t.Fatalf("unexpected number of invocations: %d", invocations)
	}
}

func TestAddOnUsage_Get(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/subscriptions/1122334455/add_ons/addOnCode/usage/1234", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write(MustOpenFile("add_on_usage.xml"))
		}, t)

		if a, err := client.AddOnUsages.Get(context.Background(), "1122334455", "addOnCode", "1234"); err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(a, NewTestAddOnUsage()); diff != "" {
			t.Fatal(diff)
		} else if !s.Invoked {
			t.Fatal("expected fn invocation")
		}
	})

	// Ensure a 404 returns nil values.
	t.Run("ErrNotFound", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/subscriptions/1122334455/add_ons/addOnCode/usage/8888", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}, t)

		if a, err := client.AddOnUsages.Get(context.Background(), "1122334455", "addOnCode", "8888"); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if a != nil {
			t.Fatalf("expected nil: %#v", a)
		}
	})
}

func TestAddOnUsage_Create(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("POST", "/v2/subscriptions/1122334455/add_ons/addOnCode/usage", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write(MustOpenFile("add_on_usage.xml"))
	}, t)

	if a, err := client.AddOnUsages.Create(context.Background(), "1122334455", "addOnCode", recurly.AddOnUsage{}); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(a, NewTestAddOnUsage()); diff != "" {
		t.Fatal(diff)
	}
}

func TestAddOnUsage_Update(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("PUT", "/v2/subscriptions/1122334455/add_ons/addOnCode/usage/1234", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("add_on_usage.xml"))
	}, t)

	if a, err := client.AddOnUsages.Update(context.Background(), "1122334455", "addOnCode", "1234", recurly.AddOnUsage{}); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(a, NewTestAddOnUsage()); diff != "" {
		t.Fatal(diff)
	}
}

func TestAddOnUsage_Delete(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("DELETE", "/v2/subscriptions/1122334455/add_ons/addOnCode/usage/1234", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}, t)

	if err := client.AddOnUsages.Delete(context.Background(), "1122334455", "addOnCode", "1234"); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	}
}


// Returns add on corresponding to testdata/add_on.xml
func NewTestAddOnUsage() *recurly.AddOnUsage {

	now := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	return &recurly.AddOnUsage{
		XMLName:            xml.Name{Local: "usage"},
		Amount:             1,
		MerchantTag:        "Order ID: 4939853977878713",
		RecordingTimestamp: recurly.NewTime(now),
		UsageTimestamp:     recurly.NewTime(now),
		CreatedAt:          recurly.NewTime(now),
		UpdatedAt:          recurly.NullTime{},
		BilledAt:           recurly.NullTime{},
		UsageType:          "price",
		UnitAmountInCents:  45,
		UsagePercentage:    recurly.NullFloat{},
	}
}