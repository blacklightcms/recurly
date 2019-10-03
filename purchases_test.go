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
func TestPurchases_Purchase_Encoding(t *testing.T) {
	tests := []struct {
		v        recurly.Purchase
		expected string
	}{
		{
			expected: MustCompactString(`
				<purchase>
					<account></account>
					<adjustments></adjustments>
					<currency></currency>
					<gift_card></gift_card>
					<coupon_codes></coupon_codes>
					<subscriptions></subscriptions>
				</purchase>
			`),
		},
		{
			v:        *NewTestPurchase(),
			expected: string(MustOpenCompactXMLFile("purchase.xml")),
		},
		{
			v: recurly.Purchase{ShippingAddressID: 2438622711411416831},
			expected: MustCompactString(`
				<purchase>
					<account></account>
					<adjustments></adjustments>
					<currency></currency>
					<gift_card></gift_card>
					<coupon_codes></coupon_codes>
					<subscriptions></subscriptions>
					<shipping_address_id>2438622711411416831</shipping_address_id>
				</purchase>
			`),
		},
		{
			v: recurly.Purchase{
				ShippingFees: &[]recurly.ShippingFee{
					{
						ShippingMethodCode:    "foo",
						ShippingAmountInCents: recurly.NewInt(0),
					},
					{
						ShippingMethodCode:    "bar",
						ShippingAmountInCents: recurly.NewInt(10),
					},
				},
				ShippingAddressID: 1,
				TransactionType:   "moto",
			},
			expected: MustCompactString(`
				<purchase>
					<account></account>
					<adjustments></adjustments>
					<currency></currency>
					<gift_card></gift_card>
					<coupon_codes></coupon_codes>
					<subscriptions></subscriptions>
					<shipping_address_id>1</shipping_address_id>
					<shipping_fees>
						<shipping_fee>
							<shipping_method_code>foo</shipping_method_code>
							<shipping_amount_in_cents>0</shipping_amount_in_cents>
						</shipping_fee>
						<shipping_fee>
							<shipping_method_code>bar</shipping_method_code>
							<shipping_amount_in_cents>10</shipping_amount_in_cents>
						</shipping_fee>
					</shipping_fees>
					<transaction_type>moto</transaction_type>
				</purchase>
			`),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			buf := new(bytes.Buffer)
			if err := xml.NewEncoder(buf).Encode(tt.v); err != nil {
				t.Fatal(err)
			} else if tt.expected != buf.String() {
				t.Fatal(buf.String()+"\n\n", tt.expected)
			}
		})
	}
}

func TestPurchases_Create(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("POST", "/v2/purchases", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write(MustOpenFile("invoice_collection.xml"))
	}, t)

	if collection, err := client.Purchases.Create(context.Background(), recurly.Purchase{}); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(collection, NewTestInvoiceCollection()); diff != "" {
		t.Fatal(diff)
	}
}

func TestPurchases_Preview(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("POST", "/v2/purchases/preview", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("invoice_collection.xml"))
	}, t)

	if collection, err := client.Purchases.Preview(context.Background(), recurly.Purchase{}); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(collection, NewTestInvoiceCollection()); diff != "" {
		t.Fatal(diff)
	}
}

func TestPurchases_Authorize(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("POST", "/v2/purchases/authorize", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("purchase.xml"))
	}, t)

	if purchase, err := client.Purchases.Authorize(context.Background(), recurly.Purchase{}); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(purchase, NewTestPurchase()); diff != "" {
		t.Fatal(diff)
	}
}

func TestPurchases_Pending(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("POST", "/v2/purchases/pending", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("purchase.xml"))
	}, t)

	if purchase, err := client.Purchases.Pending(context.Background(), recurly.Purchase{}); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(purchase, NewTestPurchase()); diff != "" {
		t.Fatal(diff)
	}
}

func TestPurchases_Capture(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("POST", "/v2/purchases/a13acd8fe4294916b79aec87b7ea441f/capture", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("invoice_collection.xml"))
	}, t)

	if collection, err := client.Purchases.Capture(context.Background(), "a13acd8f-e429-4916-b79a-ec87b7ea441f"); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(collection, NewTestInvoiceCollection()); diff != "" {
		t.Fatal(diff)
	}
}

func TestPurchases_Cancel(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	s.HandleFunc("POST", "/v2/purchases/a13acd8fe4294916b79aec87b7ea441f/cancel", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("invoice_collection.xml"))
	}, t)

	if collection, err := client.Purchases.Cancel(context.Background(), "a13acd8f-e429-4916-b79a-ec87b7ea441f"); !s.Invoked {
		t.Fatal("expected fn invocation")
	} else if err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(collection, NewTestInvoiceCollection()); diff != "" {
		t.Fatal(diff)
	}
}

// Returns a Purchase corresponding to testdata/purchase.xml.
func NewTestPurchase() *recurly.Purchase {
	return &recurly.Purchase{
		XMLName:               xml.Name{Local: "purchase"},
		CollectionMethod:      "automatic",
		Currency:              "USD",
		CustomerNotes:         "Some notes for the customer.",
		TermsAndConditions:    "Our company terms and conditions.",
		VATReverseChargeNotes: "Vat reverse charge notes.",
		GatewayCode:           "test-gateway-code",
		Account: recurly.Account{
			XMLName: xml.Name{Local: "account"},
			Code:    "c442b36c-c64f-41d7-b8e1-9c04e7a6ff82",
			ShippingAddresses: &[]recurly.ShippingAddress{{
				XMLName:   xml.Name{Local: "shipping_address"},
				FirstName: "Lon",
				LastName:  "Doner",
				Address:   "221B Baker St.",
				City:      "London",
				Zip:       "W1K 6AH",
				Country:   "GB",
				Nickname:  "Home",
			}},
			BillingInfo: &recurly.Billing{
				XMLName:   xml.Name{Local: "billing_info"},
				Address:   "400 Alabama St",
				City:      "San Francisco",
				Country:   "US",
				FirstName: "Benjamin",
				LastName:  "Du Monde",
				Month:     12,
				Number:    4111111111111111,
				State:     "CA",
				Year:      2019,
				Zip:       "94110",
			},
		},
		Adjustments: []recurly.Adjustment{{
			XMLName:           xml.Name{Local: "adjustment"},
			ProductCode:       "4549449c-5870-4845-b672-1d07f15e87dd",
			Quantity:          1,
			UnitAmountInCents: recurly.NewInt(1000),
			Description:       "Description of this adjustment",
		}},
		Subscriptions: []recurly.PurchaseSubscription{{
			XMLName:  xml.Name{Local: "subscription"},
			PlanCode: "plan1",
		}},
		CouponCodes: []string{"coupon1", "coupon2"},
		GiftCard:    "ABC1234",
	}
}
