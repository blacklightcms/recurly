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
func TestAddOns_Encoding(t *testing.T) {
	tests := []struct {
		v        recurly.AddOn
		expected string
	}{
		{
			expected: MustCompactString(`
				<add_on>
				</add_on>
			`),
		},
		{
			v: recurly.AddOn{Code: "xyz"},
			expected: MustCompactString(`
				<add_on>
					<add_on_code>xyz</add_on_code>
				</add_on>
			`),
		},
		{
			v: recurly.AddOn{Name: "IP Addresses"},
			expected: MustCompactString(`
				<add_on>
					<name>IP Addresses</name>
				</add_on>
			`),
		},
		{
			v: recurly.AddOn{DefaultQuantity: recurly.NewInt(0)},
			expected: MustCompactString(`
				<add_on>
					<default_quantity>0</default_quantity>
				</add_on>
			`),
		},
		{
			v: recurly.AddOn{DefaultQuantity: recurly.NewInt(1)},
			expected: MustCompactString(`
				<add_on>
					<default_quantity>1</default_quantity>
				</add_on>
			`),
		},
		{
			v: recurly.AddOn{DisplayQuantityOnHostedPage: recurly.NewBool(true)},
			expected: MustCompactString(`
				<add_on>
					<display_quantity_on_hosted_page>true</display_quantity_on_hosted_page>
				</add_on>
			`),
		},
		{
			v: recurly.AddOn{DisplayQuantityOnHostedPage: recurly.NewBool(false)},
			expected: MustCompactString(`
				<add_on>
					<display_quantity_on_hosted_page>false</display_quantity_on_hosted_page>
				</add_on>
			`),
		},
		{
			v: recurly.AddOn{TaxCode: "digital"},
			expected: MustCompactString(`
				<add_on>
					<tax_code>digital</tax_code>
				</add_on>
			`),
		},
		{
			v: recurly.AddOn{UnitAmountInCents: recurly.UnitAmount{USD: 200}},
			expected: MustCompactString(`
				<add_on>
					<unit_amount_in_cents>
						<USD>200</USD>
					</unit_amount_in_cents>
				</add_on>
			`),
		},
		{
			v: recurly.AddOn{AccountingCode: "abc123"},
			expected: MustCompactString(`
				<add_on>
					<accounting_code>abc123</accounting_code>
				</add_on>
			`),
		},
		{
			v: recurly.AddOn{ItemCode: "pink_sweaters"},
			expected: MustCompactString(`
				<add_on>
					<item_code>pink_sweaters</item_code>
				</add_on>
			`),
		},
		{
			v: recurly.AddOn{ExternalSKU: "BC-123-ABC"},
			expected: MustCompactString(`
				<add_on>
					<external_sku>BC-123-ABC</external_sku>
				</add_on>
			`),
		},
		{
			v: recurly.AddOn{ItemState: "active"},
			expected: MustCompactString(`
				<add_on>
					<item_state>active</item_state>
				</add_on>
			`),
		},
		{
			v: recurly.AddOn{AvalaraServiceType: 300, AvalaraTransactionType: 6},
			expected: MustCompactString(`
				<add_on>
					<avalara_transaction_type>6</avalara_transaction_type>
					<avalara_service_type>300</avalara_service_type>
				</add_on>
			`),
		},
		{
			v: recurly.AddOn{TierType: "flat"},
			expected: MustCompactString(`
				<add_on>
					<tier_type>flat</tier_type>
				</add_on>
			`),
		},
		{
			v: recurly.AddOn{Tiers: &[]recurly.Tier{*NewTestTier()}},
			expected: MustCompactString(`
				<add_on>
					<tiers>
						<tier>
							<unit_amount_in_cents>
								<USD>100</USD>
							</unit_amount_in_cents>
							<ending_quantity>500</ending_quantity>
						</tier>
					</tiers>
				</add_on>
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

func TestUnitAmount(t *testing.T) {
	type s struct {
		Amount recurly.UnitAmount `xml:"amount,omitempty"`
	}

	tests := []struct {
		v        s
		expected string
	}{
		{
			expected: "<s></s>",
		},
		{v: s{Amount: recurly.UnitAmount{USD: 1000}},
			expected: MustCompactString(`
				<s>
					<amount>
						<USD>1000</USD>
					</amount>
				</s>
			`),
		},
		{v: s{Amount: recurly.UnitAmount{USD: 800, EUR: 650}},
			expected: MustCompactString(`
				<s>
					<amount>
						<USD>800</USD>
						<EUR>650</EUR>
					</amount>
				</s>
			`),
		},
		{v: s{Amount: recurly.UnitAmount{EUR: 650}},
			expected: MustCompactString(`
				<s>
					<amount>
						<EUR>650</EUR>
					</amount>
				</s>
			`),
		},
		{v: s{Amount: recurly.UnitAmount{GBP: 3000}},
			expected: MustCompactString(`
				<s>
					<amount>
						<GBP>3000</GBP>
					</amount>
				</s>
			`),
		},
		{v: s{Amount: recurly.UnitAmount{CAD: 300}},
			expected: MustCompactString(`
				<s>
					<amount>
						<CAD>300</CAD>
					</amount>
				</s>
			`),
		},
		{v: s{Amount: recurly.UnitAmount{AUD: 400}},
			expected: MustCompactString(`
				<s>
					<amount>
						<AUD>400</AUD>
					</amount>
				</s>
			`),
		},
		{v: s{Amount: recurly.UnitAmount{USD: 1}},
			expected: MustCompactString(`
				<s>
					<amount>
						<USD>1</USD>
					</amount>
				</s>
			`),
		},
	}

	buf := new(bytes.Buffer)
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if err := xml.NewEncoder(buf).Encode(tt.v); err != nil {
				t.Fatal(err)
			} else if buf.String() != tt.expected {
				t.Fatal(buf.String())
			}
			buf.Reset()
		})
	}
}

func TestAddOns_List(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	var invocations int
	s.HandleFunc("GET", "/v2/plans/gold/add_ons", func(w http.ResponseWriter, r *http.Request) {
		invocations++
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("add_ons.xml"))
	}, t)

	pager := client.AddOns.List("gold", nil)
	for pager.Next() {
		var a []recurly.AddOn
		if err := pager.Fetch(context.Background(), &a); err != nil {
			t.Fatal(err)
		} else if !s.Invoked {
			t.Fatal("expected s to be invoked")
		} else if diff := cmp.Diff(a, []recurly.AddOn{*NewTestAddOn()}); diff != "" {
			t.Fatal(diff)
		}
	}
	if invocations != 1 {
		t.Fatalf("unexpected number of invocations: %d", invocations)
	}
}

func TestAddOns_Get(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/plans/gold/add_ons/ipaddresses", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write(MustOpenFile("add_on.xml"))
		}, t)

		if a, err := client.AddOns.Get(context.Background(), "gold", "ipaddresses"); err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(a, NewTestAddOn()); diff != "" {
			t.Fatal(diff)
		} else if !s.Invoked {
			t.Fatal("expected fn invocation")
		}
	})

	// Ensure a 404 returns nil values.
	t.Run("ErrNotFound", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/plans/gold/add_ons/ipaddresses", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}, t)

		if a, err := client.AddOns.Get(context.Background(), "gold", "ipaddresses"); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if a != nil {
			t.Fatalf("expected nil: %#v", a)
		}
	})
}

func TestAddOns_Create(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("POST", "/v2/plans/gold/add_ons", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write(MustOpenFile("add_on.xml"))
	}, t)

	if a, err := client.AddOns.Create(context.Background(), "gold", recurly.AddOn{}); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(a, NewTestAddOn()); diff != "" {
		t.Fatal(diff)
	}
}

func TestAddOns_Update(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("PUT", "/v2/plans/gold/add_ons/ipaddresses", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("add_on.xml"))
	}, t)

	if a, err := client.AddOns.Update(context.Background(), "gold", "ipaddresses", recurly.AddOn{}); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(a, NewTestAddOn()); diff != "" {
		t.Fatal(diff)
	}
}

func TestAddOns_Delete(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("DELETE", "/v2/plans/gold/add_ons/ipaddresses", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}, t)

	if err := client.AddOns.Delete(context.Background(), "gold", "ipaddresses"); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	}
}

// Returns add on corresponding to testdata/add_on.xml
func NewTestAddOn() *recurly.AddOn {
	return &recurly.AddOn{
		XMLName:                     xml.Name{Local: "add_on"},
		Code:                        "ipaddresses",
		Name:                        "IP Addresses",
		DefaultQuantity:             recurly.NewInt(1),
		DisplayQuantityOnHostedPage: recurly.NewBool(false),
		TaxCode:                     "digital",
		UnitAmountInCents: recurly.UnitAmount{
			USD: 200,
		},
		TierType:       "volume",
		Tiers:          &[]recurly.Tier{*NewTestTier()},
		AccountingCode: "abc123",
		CreatedAt:      recurly.NewTime(MustParseTime("2011-06-28T12:34:56Z")),
	}
}

func NewTestTier() *recurly.Tier {
	return &recurly.Tier{
		XMLName:        xml.Name{Local: "tier"},
		EndingQuantity: 500,
		UnitAmountInCents: recurly.UnitAmount{
			USD: 100,
		},
	}
}
