package recurly_test

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/launchpadcentral/recurly"
)

func TestPurchases_Purchase_Encoding(t *testing.T) {
	tests := []struct {
		v        recurly.Purchase
		expected string
	}{
		{
			expected: "<purchase><account></account><adjustments></adjustments><currency></currency><gift_card></gift_card><coupon_codes></coupon_codes><subscriptions></subscriptions></purchase>",
		},
		{
			// From Recurly examples: https://dev.recurly.com/docs/create-purchase
			v: recurly.Purchase{
				CollectionMethod:      "automatic",
				Currency:              "USD",
				CustomerNotes:         "Some notes for the customer.",
				TermsAndConditions:    "Our company terms and conditions.",
				VATReverseChargeNotes: "Vat reverse charge notes.",
				Account: recurly.Account{
					Code: "c442b36c-c64f-41d7-b8e1-9c04e7a6ff82",
					ShippingAddresses: &[]recurly.ShippingAddress{
						recurly.ShippingAddress{
							FirstName: "Lon",
							LastName:  "Doner",
							Address:   "221B Baker St.",
							City:      "London",
							Zip:       "W1K 6AH",
							Country:   "GB",
							Nickname:  "Home",
						},
					},
					BillingInfo: &recurly.Billing{
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
				Adjustments: []recurly.Adjustment{
					recurly.Adjustment{
						ProductCode: "4549449c-5870-4845-b672-1d07f15e87dd",
						Quantity:    1,
						// RevenueScheduleType: "at_invoice",
						UnitAmountInCents: 1000,
						Description:       "Description of this adjustment",
					},
				},
				Subscriptions: []recurly.NewSubscription{
					recurly.NewSubscription{PlanCode: "plan1"},
				},
				CouponCodes: []string{"coupon1", "coupon2"},
				GiftCard:    "ABC1234",
			},
			expected: "<purchase>" +
				"<account><account_code>c442b36c-c64f-41d7-b8e1-9c04e7a6ff82</account_code>" +
				"<billing_info><first_name>Benjamin</first_name><last_name>Du Monde</last_name><address1>400 Alabama St</address1><city>San Francisco</city><state>CA</state><zip>94110</zip><country>US</country><number>4111111111111111</number><month>12</month><year>2019</year></billing_info>" +
				"<shipping_addresses><shipping_address><first_name>Lon</first_name><last_name>Doner</last_name><nickname>Home</nickname><address1>221B Baker St.</address1><address2></address2><city>London</city><state></state><zip>W1K 6AH</zip><country>GB</country></shipping_address></shipping_addresses>" +
				"</account>" +
				"<adjustments><adjustment>" +
				"<description>Description of this adjustment</description><product_code>4549449c-5870-4845-b672-1d07f15e87dd</product_code><unit_amount_in_cents>1000</unit_amount_in_cents><quantity>1</quantity>" +
				// TODO: RevenueScheduleType not yet modeled
				// "<revenue_schedule_type>at_invoice</revenue_schedule_type>"+
				"</adjustment></adjustments>" +
				"<collection_method>automatic</collection_method>" +
				"<currency>USD</currency>" +
				"<gift_card><redemption_code>ABC1234</redemption_code></gift_card>" +
				"<coupon_codes><coupon_code>coupon1</coupon_code><coupon_code>coupon2</coupon_code></coupon_codes>" +
				"<subscriptions><subscription><plan_code>plan1</plan_code><account></account><currency></currency></subscription></subscriptions>" +
				"<customer_notes>Some notes for the customer.</customer_notes>" +
				"<terms_and_conditions>Our company terms and conditions.</terms_and_conditions>" +
				"<vat_reverse_charge_notes>Vat reverse charge notes.</vat_reverse_charge_notes>" +
				"</purchase>",
		},
		{
			v:        recurly.Purchase{ShippingAddressID: 2438622711411416831},
			expected: "<purchase><account></account><adjustments></adjustments><currency></currency><gift_card></gift_card><coupon_codes></coupon_codes><subscriptions></subscriptions><shipping_address_id>2438622711411416831</shipping_address_id></purchase>",
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var given bytes.Buffer
			if err := xml.NewEncoder(&given).Encode(tt.v); err != nil {
				t.Fatalf("(%d) unexpected encode error: %v", i, err)
			} else if tt.expected != given.String() {
				t.Fatalf("(%d) unexpected value: \n%s\n%s", i, given.String(), tt.expected)
			}
		})
	}
}

func TestPurchases_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/purchases", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(201)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
		<invoice_collection>
			<charge_invoice href="https://your-subdomain.recurly.com/v2/invoices/2531">
				<account href="https://your-subdomain.recurly.com/v2/accounts/1"/>
				<address>
					<address1>123 Main St.</address1>
					<address2></address2>
					<city>SF</city>
					<state>CA</state>
					<zip>94105</zip>
					<country>US</country>
					<phone></phone>
				</address>
				<uuid>458465c8bd2035928b441940d1a17179</uuid>
				<state>paid</state>
				<invoice_number_prefix></invoice_number_prefix>
				<invoice_number type="integer">2531</invoice_number>
				<vat_number></vat_number>
				<tax_in_cents type="integer">111</tax_in_cents>
				<total_in_cents type="integer">1411</total_in_cents>
				<currency>USD</currency>
				<created_at type="datetime">2018-06-19T01:13:26Z</created_at>
				<updated_at type="datetime">2018-06-19T01:13:26Z</updated_at>
				<attempt_next_collection_at nil="nil"></attempt_next_collection_at>
				<closed_at type="datetime">2018-06-19T01:13:26Z</closed_at>
				<customer_notes nil="nil"></customer_notes>
				<recovery_reason nil="nil"></recovery_reason>
				<subtotal_before_discount_in_cents type="integer">1300</subtotal_before_discount_in_cents>
				<subtotal_in_cents type="integer">1300</subtotal_in_cents>
				<discount_in_cents type="integer">0</discount_in_cents>
				<due_on type="datetime">2018-06-19T01:13:26Z</due_on>
				<net_terms type="integer">0</net_terms>
				<collection_method>automatic</collection_method>
				<po_number nil="nil"></po_number>
				<terms_and_conditions nil="nil"></terms_and_conditions>
				<tax_type>usst</tax_type>
				<tax_region>CA</tax_region>
				<tax_rate type="float">0.085</tax_rate>
				<line_items type="array">
				<adjustment href="https://your-subdomain.recurly.com/v2/adjustments/458465c8a977efb18125714c5abc4119" type="charge">
					<account href="https://your-subdomain.recurly.com/v2/accounts/1"/>
					<invoice href="https://your-subdomain.recurly.com/v2/invoices/2531"/>
					<uuid>458465c8a977efb18125714c5abc4119</uuid>
					<state>invoiced</state>
					<description>Description for Adjustment</description>
					<accounting_code nil="nil"></accounting_code>
					<product_code>freeform_98765</product_code>
					<origin>debit</origin>
					<unit_amount_in_cents type="integer">100</unit_amount_in_cents>
					<quantity type="integer">13</quantity>
					<discount_in_cents type="integer">0</discount_in_cents>
					<tax_in_cents type="integer">111</tax_in_cents>
					<total_in_cents type="integer">1411</total_in_cents>
					<currency>USD</currency>
					<proration_rate nil="nil"></proration_rate>
					<taxable type="boolean">false</taxable>
					<tax_type>usst</tax_type>
					<tax_region>CA</tax_region>
					<tax_rate type="float">0.085</tax_rate>
					<tax_exempt type="boolean">false</tax_exempt>
					<start_date type="datetime">2018-06-19T01:13:26Z</start_date>
					<end_date nil="nil"></end_date>
					<created_at type="datetime">2018-06-19T01:13:26Z</created_at>
					<updated_at type="datetime">2018-06-19T01:13:26Z</updated_at>
					<revenue_schedule_type></revenue_schedule_type>
				</adjustment>
				</line_items>
				<transactions type="array">
				<transaction href="https://your-subdomain.recurly.com/v2/transactions/458465c9016060dbde60734789869a2d" type="credit_card">
					<account href="https://your-subdomain.recurly.com/v2/accounts/1"/>
					<invoice href="https://your-subdomain.recurly.com/v2/invoices/2531"/>
					<uuid>458465c9016060dbde60734789869a2d</uuid>
					<action>purchase</action>
					<amount_in_cents type="integer">1411</amount_in_cents>
					<tax_in_cents type="integer">111</tax_in_cents>
					<currency>USD</currency>
					<status>success</status>
					<payment_method>credit_card</payment_method>
					<reference>8001561</reference>
					<source>transaction</source>
					<recurring type="boolean">false</recurring>
					<test type="boolean">true</test>
					<voidable type="boolean">true</voidable>
					<refundable type="boolean">true</refundable>
					<ip_address nil="nil"></ip_address>
					<gateway_type>test</gateway_type>
					<origin>api</origin>
					<description>Transaction Description</description>
					<message>Successful test transaction</message>
					<approval_code nil="nil"></approval_code>
					<failure_type nil="nil"></failure_type>
					<gateway_error_codes nil="nil"></gateway_error_codes>
					<cvv_result code="" nil="nil"></cvv_result>
					<avs_result code="D">Street address and postal code match.</avs_result>
					<avs_result_street nil="nil"></avs_result_street>
					<avs_result_postal nil="nil"></avs_result_postal>
					<created_at type="datetime">2018-06-19T01:13:26Z</created_at>
					<collected_at type="datetime">2018-06-19T01:13:26Z</collected_at>
					<updated_at type="datetime">2018-06-19T01:13:26Z</updated_at>
					<details>
					<account>
						<account_code>1</account_code>
						<first_name>Verena</first_name>
						<last_name>Example</last_name>
						<company nil="nil"></company>
						<email>verena@example.com</email>
						<billing_info type="credit_card">
							<first_name>Verena</first_name>
							<last_name>Example</last_name>
							<address1>123 Main St.</address1>
							<address2></address2>
							<city>SF</city>
							<state>CA</state>
							<zip>94105</zip>
							<country>US</country>
							<phone></phone>
							<vat_number></vat_number>
							<card_type>Visa</card_type>
							<year type="integer">2020</year>
							<month type="integer">6</month>
							<first_six>411111</first_six>
							<last_four>1111</last_four>
						</billing_info>
					</account>
					</details>
					</transaction>
				</transactions>
				<a name="refund" href="https://your-subdomain.recurly.com/v2/invoices/2531/refund" method="post"/>
			</charge_invoice>
			<credit_invoices type="array">
			</credit_invoices>
		</invoice_collection>`)
	})

	expectedInvoiceNumber := 2531
	expectedTax := 111
	expectedTotal := 1411

	r, invs, err := client.Purchases.Create(recurly.Purchase{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected create purchase to return OK")
	} else if invs.ChargeInvoice.InvoiceNumber != expectedInvoiceNumber {
		t.Fatalf("expected invoice number %d; got %d\n", expectedInvoiceNumber, invs.ChargeInvoice.InvoiceNumber)
	} else if invs.ChargeInvoice.TaxInCents != expectedTax {
		t.Fatalf("expected tax %d; got %d\n", expectedTax, invs.ChargeInvoice.TaxInCents)
	} else if invs.ChargeInvoice.TotalInCents != expectedTotal {
		t.Fatalf("expected total %d; got %d\n", expectedTotal, invs.ChargeInvoice.TotalInCents)
	}
}

func TestPurchases_Preview(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/purchases/preview", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(201)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
		<invoice_collection>
			<charge_invoice href="https://your-subdomain.recurly.com/v2/invoices/2531">
				<account href="https://your-subdomain.recurly.com/v2/accounts/1"/>
				<address>
					<address1>123 Main St.</address1>
					<address2></address2>
					<city>SF</city>
					<state>CA</state>
					<zip>94105</zip>
					<country>US</country>
					<phone></phone>
				</address>
				<uuid>458465c8bd2035928b441940d1a17179</uuid>
				<state>paid</state>
				<invoice_number_prefix></invoice_number_prefix>
				<invoice_number type="integer">2531</invoice_number>
				<vat_number></vat_number>
				<tax_in_cents type="integer">111</tax_in_cents>
				<total_in_cents type="integer">1411</total_in_cents>
				<currency>USD</currency>
				<created_at type="datetime">2018-06-19T01:13:26Z</created_at>
				<updated_at type="datetime">2018-06-19T01:13:26Z</updated_at>
				<attempt_next_collection_at nil="nil"></attempt_next_collection_at>
				<closed_at type="datetime">2018-06-19T01:13:26Z</closed_at>
				<customer_notes nil="nil"></customer_notes>
				<recovery_reason nil="nil"></recovery_reason>
				<subtotal_before_discount_in_cents type="integer">1300</subtotal_before_discount_in_cents>
				<subtotal_in_cents type="integer">1300</subtotal_in_cents>
				<discount_in_cents type="integer">0</discount_in_cents>
				<due_on type="datetime">2018-06-19T01:13:26Z</due_on>
				<net_terms type="integer">0</net_terms>
				<collection_method>automatic</collection_method>
				<po_number nil="nil"></po_number>
				<terms_and_conditions nil="nil"></terms_and_conditions>
				<tax_type>usst</tax_type>
				<tax_region>CA</tax_region>
				<tax_rate type="float">0.085</tax_rate>
				<line_items type="array">
				<adjustment href="https://your-subdomain.recurly.com/v2/adjustments/458465c8a977efb18125714c5abc4119" type="charge">
					<account href="https://your-subdomain.recurly.com/v2/accounts/1"/>
					<invoice href="https://your-subdomain.recurly.com/v2/invoices/2531"/>
					<uuid>458465c8a977efb18125714c5abc4119</uuid>
					<state>invoiced</state>
					<description>Description for Adjustment</description>
					<accounting_code nil="nil"></accounting_code>
					<product_code>freeform_98765</product_code>
					<origin>debit</origin>
					<unit_amount_in_cents type="integer">100</unit_amount_in_cents>
					<quantity type="integer">13</quantity>
					<discount_in_cents type="integer">0</discount_in_cents>
					<tax_in_cents type="integer">111</tax_in_cents>
					<total_in_cents type="integer">1411</total_in_cents>
					<currency>USD</currency>
					<proration_rate nil="nil"></proration_rate>
					<taxable type="boolean">false</taxable>
					<tax_type>usst</tax_type>
					<tax_region>CA</tax_region>
					<tax_rate type="float">0.085</tax_rate>
					<tax_exempt type="boolean">false</tax_exempt>
					<start_date type="datetime">2018-06-19T01:13:26Z</start_date>
					<end_date nil="nil"></end_date>
					<created_at type="datetime">2018-06-19T01:13:26Z</created_at>
					<updated_at type="datetime">2018-06-19T01:13:26Z</updated_at>
					<revenue_schedule_type></revenue_schedule_type>
				</adjustment>
				</line_items>
				<transactions type="array">
				<transaction href="https://your-subdomain.recurly.com/v2/transactions/458465c9016060dbde60734789869a2d" type="credit_card">
					<account href="https://your-subdomain.recurly.com/v2/accounts/1"/>
					<invoice href="https://your-subdomain.recurly.com/v2/invoices/2531"/>
					<uuid>458465c9016060dbde60734789869a2d</uuid>
					<action>purchase</action>
					<amount_in_cents type="integer">1411</amount_in_cents>
					<tax_in_cents type="integer">111</tax_in_cents>
					<currency>USD</currency>
					<status>success</status>
					<payment_method>credit_card</payment_method>
					<reference>8001561</reference>
					<source>transaction</source>
					<recurring type="boolean">false</recurring>
					<test type="boolean">true</test>
					<voidable type="boolean">true</voidable>
					<refundable type="boolean">true</refundable>
					<ip_address nil="nil"></ip_address>
					<gateway_type>test</gateway_type>
					<origin>api</origin>
					<description>Transaction Description</description>
					<message>Successful test transaction</message>
					<approval_code nil="nil"></approval_code>
					<failure_type nil="nil"></failure_type>
					<gateway_error_codes nil="nil"></gateway_error_codes>
					<cvv_result code="" nil="nil"></cvv_result>
					<avs_result code="D">Street address and postal code match.</avs_result>
					<avs_result_street nil="nil"></avs_result_street>
					<avs_result_postal nil="nil"></avs_result_postal>
					<created_at type="datetime">2018-06-19T01:13:26Z</created_at>
					<collected_at type="datetime">2018-06-19T01:13:26Z</collected_at>
					<updated_at type="datetime">2018-06-19T01:13:26Z</updated_at>
					<details>
					<account>
						<account_code>1</account_code>
						<first_name>Verena</first_name>
						<last_name>Example</last_name>
						<company nil="nil"></company>
						<email>verena@example.com</email>
						<billing_info type="credit_card">
							<first_name>Verena</first_name>
							<last_name>Example</last_name>
							<address1>123 Main St.</address1>
							<address2></address2>
							<city>SF</city>
							<state>CA</state>
							<zip>94105</zip>
							<country>US</country>
							<phone></phone>
							<vat_number></vat_number>
							<card_type>Visa</card_type>
							<year type="integer">2020</year>
							<month type="integer">6</month>
							<first_six>411111</first_six>
							<last_four>1111</last_four>
						</billing_info>
					</account>
					</details>
					</transaction>
				</transactions>
				<a name="refund" href="https://your-subdomain.recurly.com/v2/invoices/2531/refund" method="post"/>
			</charge_invoice>
			<credit_invoices type="array">
			</credit_invoices>
		</invoice_collection>`)
	})

	expectedInvoiceNumber := 2531
	expectedTax := 111
	expectedTotal := 1411

	r, invs, err := client.Purchases.Preview(recurly.Purchase{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected preview purchase to return OK")
	} else if invs.ChargeInvoice.InvoiceNumber != expectedInvoiceNumber {
		t.Fatalf("expected invoice number %d; got %d\n", expectedInvoiceNumber, invs.ChargeInvoice.InvoiceNumber)
	} else if invs.ChargeInvoice.TaxInCents != expectedTax {
		t.Fatalf("expected tax %d; got %d\n", expectedTax, invs.ChargeInvoice.TaxInCents)
	} else if invs.ChargeInvoice.TotalInCents != expectedTotal {
		t.Fatalf("expected total %d; got %d\n", expectedTotal, invs.ChargeInvoice.TotalInCents)
	}
}
