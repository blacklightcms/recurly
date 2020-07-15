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
func TestItems_Encoding(t *testing.T) {
	tests := []struct {
		v        recurly.Item
		expected string
	}{
		{
			expected: MustCompactString(`
				<item></item>
			`),
		},
		{
			v: recurly.Item{Code: "abc"},
			expected: MustCompactString(`
				<item>
					<item_code>abc</item_code>
				</item>
			`),
		},
		{
			v: recurly.Item{Name: "Item"},
			expected: MustCompactString(`
				<item>
					<name>Item</name>
				</item>
			`),
		},
		{
			v: recurly.Item{State: "active"},
			expected: MustCompactString(`
				<item>
					<state>active</state>
				</item>
			`),
		},
		{
			v: recurly.Item{Description: "An Item"},
			expected: MustCompactString(`
				<item>
					<description>An Item</description>
				</item>
			`),
		},
		{
			v: recurly.Item{ExternalSKU: "BCN-ZZ-ABC-07"},
			expected: MustCompactString(`
				<item>
					<external_sku>BCN-ZZ-ABC-07</external_sku>
				</item>
			`),
		},
		{
			v: recurly.Item{TaxExempt: recurly.NewBool(true)},
			expected: MustCompactString(`
				<item>
					<tax_exempt>true</tax_exempt>
				</item>
			`),
		},
		{
			v: recurly.Item{TaxExempt: recurly.NewBool(false)},
			expected: MustCompactString(`
				<item>
					<tax_exempt>false</tax_exempt>
				</item>
			`),
		},
		{
			v: recurly.Item{AvalaraServiceType: 6},
			expected: MustCompactString(`
				<item>
					<avalara_service_type>6</avalara_service_type>
				</item>
			`),
		},
		{
			v: recurly.Item{AvalaraTransactionType: 300},
			expected: MustCompactString(`
				<item>
					<avalara_transaction_type>300</avalara_transaction_type>
				</item>
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

func TestItems_List(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	var invocations int
	s.HandleFunc("GET", "/v2/items", func(w http.ResponseWriter, r *http.Request) {
		invocations++
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("items.xml"))
	}, t)

	pager := client.Items.List(nil)
	for pager.Next() {
		var a []recurly.Item
		if err := pager.Fetch(context.Background(), &a); err != nil {
			t.Fatal(err)
		} else if !s.Invoked {
			t.Fatal("expected s to be invoked")
		} else if diff := cmp.Diff(a, []recurly.Item{*NewTestItem()}); diff != "" {
			t.Fatal(diff)
		}
	}
	if invocations != 1 {
		t.Fatalf("unexpected number of invocations: %d", invocations)
	}
}

func TestItems_Get(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/items/pink_sweaters", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write(MustOpenFile("item.xml"))
		}, t)

		if a, err := client.Items.Get(context.Background(), "pink_sweaters"); err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(a, NewTestItem()); diff != "" {
			t.Fatal(diff)
		} else if !s.Invoked {
			t.Fatal("expected fn invocation")
		}
	})

	// Ensure a 404 returns nil values.
	t.Run("ErrNotFound", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/items/pink_sweaters", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}, t)

		if a, err := client.Items.Get(context.Background(), "pink_sweaters"); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if a != nil {
			t.Fatalf("expected nil item: %#v", a)
		}
	})
}

func TestItems_Create(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("POST", "/v2/items", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write(MustOpenFile("item.xml"))
	}, t)

	if a, err := client.Items.Create(context.Background(), recurly.Item{}); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(a, NewTestItem()); diff != "" {
		t.Fatal(diff)
	}
}

func TestItems_Update(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("PUT", "/v2/items/pink_sweaters", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("item.xml"))
	}, t)

	if a, err := client.Items.Update(context.Background(), "pink_sweaters", recurly.Item{}); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(a, NewTestItem()); diff != "" {
		t.Fatal(diff)
	}
}

func TestItems_Deactivate(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("DELETE", "/v2/items/pink_sweaters", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}, t)

	if err := client.Items.Deactivate(context.Background(), "pink_sweaters"); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	}
}

// Returns an item corresponding to testdata/item.xml
func NewTestItem() *recurly.Item {
	ts := MustParseTime("2019-11-21T15:55:19Z")
	t := recurly.NewTime(ts)
	return &recurly.Item{
		XMLName:        xml.Name{Local: "item"},
		Code:           "pink_sweaters",
		Name:           "Pink Sweaters",
		Description:    "Favorite Pink Sweaters",
		ExternalSKU:    "PS1234",
		AccountingCode: "ps0000193",
		State:          "active",
		TaxExempt:      recurly.NewBool(false),
		CustomFields:   &recurly.CustomFields{"size": "large"},
		CreatedAt:      t,
		UpdatedAt:      t,
	}
}
