package recurly_test

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/launchpadcentral/recurly"
	"github.com/google/go-cmp/cmp"
)

func TestInvoices_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/invoices", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
        <invoices type="array">
        	<invoice href="https://your-subdomain.recurly.com/v2/invoices/1005">
        		<account href="https://your-subdomain.recurly.com/v2/accounts/1"/>
        		<address>
        			<address1>400 Alabama St.</address1>
        			<address2></address2>
        			<city>San Francisco</city>
        			<state>CA</state>
        			<zip>94110</zip>
        			<country>US</country>
        			<phone></phone>
        		</address>
        		<subscription href="https://your-subdomain.recurly.com/v2/accounts/1/subscriptions"/>
        		<original_invoice href="https://your-subdomain.recurly.com/v2/invoices/938571" />
        		<uuid>421f7b7d414e4c6792938e7c49d552e9</uuid>
        		<state>open</state>
        		<invoice_number_prefix></invoice_number_prefix>
        		<invoice_number type="integer">1005</invoice_number>
        		<po_number nil="nil"></po_number>
        		<vat_number nil="nil"></vat_number>
        		<discount_in_cents type="integer">17</discount_in_cents>
        		<subtotal_in_cents type="integer">1200</subtotal_in_cents>
        		<tax_in_cents type="integer">0</tax_in_cents>
        		<total_in_cents type="integer">1200</total_in_cents>
        		<currency>USD</currency>
        		<created_at type="datetime">2011-08-25T12:00:00Z</created_at>
        		<closed_at nil="nil"></closed_at>
        		<tax_type>usst</tax_type>
        		<tax_region>CA</tax_region>
        		<tax_rate type="float">0</tax_rate>
        		<net_terms type="integer">0</net_terms>
        		<collection_method>automatic</collection_method>
        		<redemption href="https://your-subdomain.recurly.com/v2/invoices/e3f0a9e084a2468480d00ee61b090d4d/redemption"/>
                <credit_payments type="array">
                    <credit_payment href="https://your-subdomain.recurly.com/v2/credit_payments/451b7869bbb20d766b1604492d97a740">
                        <account href="https://your-subdomain.recurly.com/v2/accounts/1"/>
                        <uuid>451b7869bbb20d766b1604492d97a740</uuid>
                        <action>payment</action>
                        <currency>USD</currency>
                        <amount_in_cents type="integer">10000</amount_in_cents>
                        <original_invoice href="https://your-subdomain.recurly.com/v2/invoices/5397"/>
                        <applied_to_invoice href="https://your-subdomain.recurly.com/v2/invoices/5404"/>
                        <created_at type="datetime">2018-05-29T16:13:39Z</created_at>
                        <updated_at type="datetime">2018-05-29T16:13:39Z</updated_at>
                        <voided_at nil="nil"></voided_at>
                    </credit_payment>
                </credit_payments>
                <line_items type="array">
                    <adjustment href="https://your-subdomain.recurly.com/v2/adjustments/626db120a84102b1809909071c701c60" type="charge">
                        <account href="https://your-subdomain.recurly.com/v2/accounts/100"/>
                        <invoice href="https://your-subdomain.recurly.com/v2/invoices/1108"/>
                        <subscription href="https://your-subdomain.recurly.com/v2/subscriptions/17caaca1716f33572edc8146e0aaefde"/>
                        <uuid>626db120a84102b1809909071c701c60</uuid>
                        <state>invoiced</state>
                        <description>One-time Charged Fee</description>
                        <accounting_code nil="nil"/>
                        <product_code>basic</product_code>
                        <origin>debit</origin>
                        <unit_amount_in_cents type="integer">2000</unit_amount_in_cents>
                        <quantity type="integer">1</quantity>
                        <original_adjustment_uuid>2cc95aa62517e56d5bec3a48afa1b3b9</original_adjustment_uuid> <!-- Only shows if adjustment is a credit created from another credit. -->
                        <discount_in_cents type="integer">0</discount_in_cents>
                        <tax_in_cents type="integer">180</tax_in_cents>
                        <total_in_cents type="integer">2180</total_in_cents>
                        <currency>USD</currency>
                        <taxable type="boolean">false</taxable>
                        <tax_exempt type="boolean">false</tax_exempt>
                        <tax_code nil="nil"/>
                        <start_date type="datetime">2011-08-31T03:30:00Z</start_date>
                        <end_date nil="nil"/>
                        <created_at type="datetime">2011-08-31T03:30:00Z</created_at>
                    </adjustment>
        		</line_items>
        		<transactions type="array">
        		</transactions>
        	</invoice>
        </invoices>`)
	})

	resp, invoices, err := client.Invoices.List(recurly.Params{"per_page": 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected list invoices to return OK")
	} else if resp.Request.URL.Query().Get("per_page") != "1" {
		t.Fatalf("expected per_page parameter of 1, given %s", resp.Request.URL.Query().Get("per_page"))
	} else if diff := cmp.Diff(invoices, []recurly.Invoice{{
		XMLName:     xml.Name{Local: "invoice"},
		AccountCode: "1",
		Address: recurly.Address{
			Address: "400 Alabama St.",
			City:    "San Francisco",
			State:   "CA",
			Zip:     "94110",
			Country: "US",
		},
		OriginalInvoiceNumber: 938571,
		UUID:             "421f7b7d414e4c6792938e7c49d552e9",
		State:            recurly.InvoiceStateOpenDeprecated,
		InvoiceNumber:    1005,
		DiscountInCents:  17,
		SubtotalInCents:  1200,
		TaxInCents:       0,
		TotalInCents:     1200,
		Currency:         "USD",
		CreatedAt:        recurly.NewTimeFromString("2011-08-25T12:00:00Z"),
		TaxType:          "usst",
		TaxRegion:        "CA",
		TaxRate:          float64(0),
		NetTerms:         recurly.NewInt(0),
		CollectionMethod: "automatic",
		CreditPayments: []recurly.CreditPayment{{
			XMLName:               xml.Name{Local: "credit_payment"},
			AccountCode:           "1",
			UUID:                  "451b7869bbb20d766b1604492d97a740",
			Action:                "payment",
			AmountInCents:         10000,
			Currency:              "USD",
			OriginalInvoiceNumber: 5397,
			AppliedToInvoice:      5404,
			CreatedAt:             recurly.NewTimeFromString("2018-05-29T16:13:39Z"),
			UpdatedAt:             recurly.NewTimeFromString("2018-05-29T16:13:39Z"),
		}},
		LineItems: []recurly.Adjustment{
			{
				AccountCode:            "100",
				InvoiceNumber:          1108,
				SubscriptionUUID:       "17caaca1716f33572edc8146e0aaefde",
				UUID:                   "626db120a84102b1809909071c701c60",
				State:                  "invoiced",
				Description:            "One-time Charged Fee",
				ProductCode:            "basic",
				Origin:                 "debit",
				UnitAmountInCents:      2000,
				Quantity:               1,
				OriginalAdjustmentUUID: "2cc95aa62517e56d5bec3a48afa1b3b9",
				TaxInCents:             180,
				TotalInCents:           2180,
				Currency:               "USD",
				Taxable:                recurly.NewBool(false),
				TaxExempt:              recurly.NewBool(false),
				StartDate:              recurly.NewTimeFromString("2011-08-31T03:30:00Z"),
				CreatedAt:              recurly.NewTimeFromString("2011-08-31T03:30:00Z"),
			},
		},
	}}); diff != "" {
		t.Fatal(diff)
	}
}

func TestInvoices_ListAccount(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/1/invoices", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
        <invoices type="array">
        	<invoice href="https://your-subdomain.recurly.com/v2/invoices/1005">
        		<account href="https://your-subdomain.recurly.com/v2/accounts/1"/>
        		<address>
        			<address1>400 Alabama St.</address1>
        			<address2></address2>
        			<city>San Francisco</city>
        			<state>CA</state>
        			<zip>94110</zip>
        			<country>US</country>
        			<phone></phone>
        		</address>
        		<subscription href="https://your-subdomain.recurly.com/v2/accounts/1/subscriptions"/>
        		<uuid>421f7b7d414e4c6792938e7c49d552e9</uuid>
        		<state>open</state>
        		<invoice_number_prefix></invoice_number_prefix>
        		<invoice_number type="integer">1005</invoice_number>
        		<po_number nil="nil"></po_number>
        		<vat_number nil="nil"></vat_number>
        		<subtotal_in_cents type="integer">1200</subtotal_in_cents>
        		<tax_in_cents type="integer">0</tax_in_cents>
        		<total_in_cents type="integer">1200</total_in_cents>
        		<currency>USD</currency>
        		<created_at type="datetime">2011-08-25T12:00:00Z</created_at>
        		<closed_at nil="nil"></closed_at>
        		<tax_type>usst</tax_type>
        		<tax_region>CA</tax_region>
        		<tax_rate type="float">0</tax_rate>
        		<net_terms type="integer">0</net_terms>
        		<collection_method>automatic</collection_method>
        		<redemption href="https://your-subdomain.recurly.com/v2/invoices/e3f0a9e084a2468480d00ee61b090d4d/redemption"/>
        		<line_items type="array">
                    <adjustment href="https://your-subdomain.recurly.com/v2/adjustments/626db120a84102b1809909071c701c60" type="charge">
                        <account href="https://your-subdomain.recurly.com/v2/accounts/100"/>
                        <invoice href="https://your-subdomain.recurly.com/v2/invoices/1108"/>
                        <subscription href="https://your-subdomain.recurly.com/v2/subscriptions/17caaca1716f33572edc8146e0aaefde"/>
                        <uuid>626db120a84102b1809909071c701c60</uuid>
                        <state>invoiced</state>
                        <description>One-time Charged Fee</description>
                        <accounting_code nil="nil"/>
                        <product_code>basic</product_code>
                        <origin>debit</origin>
                        <unit_amount_in_cents type="integer">2000</unit_amount_in_cents>
                        <quantity type="integer">1</quantity>
                        <original_adjustment_uuid>2cc95aa62517e56d5bec3a48afa1b3b9</original_adjustment_uuid> <!-- Only shows if adjustment is a credit created from another credit. -->
                        <discount_in_cents type="integer">0</discount_in_cents>
                        <tax_in_cents type="integer">180</tax_in_cents>
                        <total_in_cents type="integer">2180</total_in_cents>
                        <currency>USD</currency>
                        <taxable type="boolean">false</taxable>
                        <tax_exempt type="boolean">false</tax_exempt>
                        <tax_code nil="nil"/>
                        <start_date type="datetime">2011-08-31T03:30:00Z</start_date>
                        <end_date nil="nil"/>
                        <created_at type="datetime">2011-08-31T03:30:00Z</created_at>
                    </adjustment>
        		</line_items>
        		<transactions type="array">
        		</transactions>
        	</invoice>
        </invoices>`)
	})

	resp, invoices, err := client.Invoices.ListAccount("1", recurly.Params{"per_page": 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected list invoices to return OK")
	} else if pp := resp.Request.URL.Query().Get("per_page"); pp != "1" {
		t.Fatalf("unexpected per_page: %s", pp)
	} else if diff := cmp.Diff(invoices, []recurly.Invoice{
		{
			XMLName:     xml.Name{Local: "invoice"},
			AccountCode: "1",
			Address: recurly.Address{
				Address: "400 Alabama St.",
				City:    "San Francisco",
				State:   "CA",
				Zip:     "94110",
				Country: "US",
			},
			UUID:             "421f7b7d414e4c6792938e7c49d552e9",
			State:            recurly.InvoiceStateOpenDeprecated,
			InvoiceNumber:    1005,
			SubtotalInCents:  1200,
			TaxInCents:       0,
			TotalInCents:     1200,
			Currency:         "USD",
			CreatedAt:        recurly.NewTimeFromString("2011-08-25T12:00:00Z"),
			TaxType:          "usst",
			TaxRegion:        "CA",
			TaxRate:          float64(0),
			NetTerms:         recurly.NewInt(0),
			CollectionMethod: "automatic",
			LineItems: []recurly.Adjustment{
				{
					AccountCode:            "100",
					InvoiceNumber:          1108,
					SubscriptionUUID:       "17caaca1716f33572edc8146e0aaefde",
					UUID:                   "626db120a84102b1809909071c701c60",
					State:                  "invoiced",
					Description:            "One-time Charged Fee",
					ProductCode:            "basic",
					Origin:                 "debit",
					UnitAmountInCents:      2000,
					Quantity:               1,
					OriginalAdjustmentUUID: "2cc95aa62517e56d5bec3a48afa1b3b9",
					TaxInCents:             180,
					TotalInCents:           2180,
					Currency:               "USD",
					Taxable:                recurly.NewBool(false),
					TaxExempt:              recurly.NewBool(false),
					StartDate:              recurly.NewTimeFromString("2011-08-31T03:30:00Z"),
					CreatedAt:              recurly.NewTimeFromString("2011-08-31T03:30:00Z"),
				},
			},
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestInvoices_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/invoices/5558", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
			<invoice href="https://blacklighttest.recurly.com/v2/invoices/5558">
				<account href="https://blacklighttest.recurly.com/v2/accounts/1"/>
				<address>
					<address1>400 Alabama St.</address1>
					<city>San Francisco</city>
					<state>CA</state>
					<zip>94110</zip>
					<country>US</country>
					<phone nil="nil"></phone>
				</address>
				<subscriptions href="https://blacklighttest.recurly.com/v2/invoices/5558/subscriptions"/>
				<uuid>421f7b7d414e4c6792938e7c49d552e9</uuid>
				<state>paid</state>
				<invoice_number_prefix></invoice_number_prefix>
				<invoice_number type="integer">5558</invoice_number>
				<vat_number nil="nil"></vat_number>
				<tax_in_cents type="integer">0</tax_in_cents>
				<total_in_cents type="integer">153000</total_in_cents>
				<currency>USD</currency>
				<created_at type="datetime">2018-06-05T15:44:57Z</created_at>
				<updated_at type="datetime">2018-06-05T15:44:57Z</updated_at>
				<attempt_next_collection_at nil="nil"></attempt_next_collection_at>
				<closed_at type="datetime">2018-06-05T15:44:57Z</closed_at>
				<customer_notes></customer_notes>
				<recovery_reason nil="nil"></recovery_reason>
				<subtotal_before_discount_in_cents type="integer">153000</subtotal_before_discount_in_cents>
				<subtotal_in_cents type="integer">153000</subtotal_in_cents>
				<discount_in_cents type="integer">0</discount_in_cents>
				<due_on type="datetime">2018-06-05T15:44:57Z</due_on>
				<balance_in_cents type="integer">0</balance_in_cents>
				<type>charge</type>
				<origin>purchase</origin>
				<credit_invoices href="https://blacklighttest.recurly.com/v2/invoices/5558/credit_invoices"/>
				<refundable_total_in_cents type="integer">153000</refundable_total_in_cents>
				<credit_payments type="array">
				</credit_payments>
				<net_terms type="integer">0</net_terms>
				<collection_method>automatic</collection_method>
				<po_number nil="nil"></po_number>
				<terms_and_conditions></terms_and_conditions>
				<line_items type="array">
					<adjustment href="https://blacklighttest.recurly.com/v2/adjustments/626db120a84102b1809909071c701c60" type="charge">
						<account href="https://blacklighttest.recurly.com/v2/accounts/1"/>
						<invoice href="https://blacklighttest.recurly.com/v2/invoices/5558"/>
						<subscription href="https://blacklighttest.recurly.com/v2/subscriptions/453f6aa0995e2d52c0d3e6453e9341da"/>
						<credit_adjustments href="https://blacklighttest.recurly.com/v2/adjustments/626db120a84102b1809909071c701c60/credit_adjustments"/>
						<refundable_total_in_cents type="integer">150000</refundable_total_in_cents>
						<uuid>626db120a84102b1809909071c701c60</uuid>
						<state>invoiced</state>
						<description>License</description>
						<accounting_code nil="nil"></accounting_code>
						<product_code>license</product_code>
						<origin>add_on</origin>
						<unit_amount_in_cents type="integer">150000</unit_amount_in_cents>
						<quantity type="integer">1</quantity>
						<discount_in_cents type="integer">0</discount_in_cents>
						<tax_in_cents type="integer">0</tax_in_cents>
						<total_in_cents type="integer">150000</total_in_cents>
						<currency>USD</currency>
						<proration_rate nil="nil"></proration_rate>
						<taxable type="boolean">false</taxable>
						<start_date type="datetime">2018-06-05T15:44:56Z</start_date>
						<end_date type="datetime">2018-07-05T15:44:56Z</end_date>
						<created_at type="datetime">2018-06-05T15:44:57Z</created_at>
						<updated_at type="datetime">2018-06-05T15:44:57Z</updated_at>
						<revenue_schedule_type>evenly</revenue_schedule_type>
					</adjustment>
					<adjustment href="https://blacklighttest.recurly.com/v2/adjustments/453f6aa1473a0620e4411a4fc88122cf" type="charge">
						<account href="https://blacklighttest.recurly.com/v2/accounts/1"/>
						<invoice href="https://blacklighttest.recurly.com/v2/invoices/5558"/>
						<subscription href="https://blacklighttest.recurly.com/v2/subscriptions/453f6aa0995e2d52c0d3e6453e9341da"/>
						<credit_adjustments href="https://blacklighttest.recurly.com/v2/adjustments/453f6aa1473a0620e4411a4fc88122cf/credit_adjustments"/>
						<refundable_total_in_cents type="integer">3000</refundable_total_in_cents>
						<uuid>453f6aa1473a0620e4411a4fc88122cf</uuid>
						<state>invoiced</state>
						<description>Domains</description>
						<accounting_code nil="nil"></accounting_code>
						<product_code>domains</product_code>
						<origin>add_on</origin>
						<unit_amount_in_cents type="integer">1500</unit_amount_in_cents>
						<quantity type="integer">2</quantity>
						<discount_in_cents type="integer">0</discount_in_cents>
						<tax_in_cents type="integer">0</tax_in_cents>
						<total_in_cents type="integer">3000</total_in_cents>
						<currency>USD</currency>
						<proration_rate nil="nil"></proration_rate>
						<taxable type="boolean">false</taxable>
						<start_date type="datetime">2018-06-05T15:44:56Z</start_date>
						<end_date type="datetime">2018-07-05T15:44:56Z</end_date>
						<created_at type="datetime">2018-06-05T15:44:57Z</created_at>
						<updated_at type="datetime">2018-06-05T15:44:57Z</updated_at>
						<revenue_schedule_type>evenly</revenue_schedule_type>
					</adjustment>
				</line_items>
				<transactions type="array">
					<transaction href="https://blacklighttest.recurly.com/v2/transactions/453f6aa17014a4624b636e455c8876ca" type="credit_card">
						<account href="https://blacklighttest.recurly.com/v2/accounts/1"/>
						<invoice href="https://blacklighttest.recurly.com/v2/invoices/5558"/>
						<subscriptions href="https://blacklighttest.recurly.com/v2/transactions/453f6aa17014a4624b636e455c8876ca/subscriptions"/>
						<uuid>453f6aa17014a4624b636e455c8876ca</uuid>
						<action>purchase</action>
						<amount_in_cents type="integer">153000</amount_in_cents>
						<tax_in_cents type="integer">0</tax_in_cents>
						<currency>USD</currency>
						<status>success</status>
						<payment_method>credit_card</payment_method>
						<reference>1909673</reference>
						<source>subscription</source>
						<recurring type="boolean">false</recurring>
						<test type="boolean">true</test>
						<voidable type="boolean">true</voidable>
						<refundable type="boolean">true</refundable>
						<ip_address>127.0.0.1</ip_address>
						<gateway_type>test</gateway_type>
						<origin>token_api</origin>
						<description nil="nil"></description>
						<message>Successful test transaction</message>
						<approval_code nil="nil"></approval_code>
						<failure_type nil="nil"></failure_type>
						<gateway_error_codes nil="nil"></gateway_error_codes>
						<cvv_result code="" nil="nil"></cvv_result>
						<avs_result code="D">Street address and postal code match.</avs_result>
						<avs_result_street nil="nil"></avs_result_street>
						<avs_result_postal nil="nil"></avs_result_postal>
						<created_at type="datetime">2018-06-05T15:44:56Z</created_at>
						<collected_at type="datetime">2018-06-05T15:44:56Z</collected_at>
						<updated_at type="datetime">2018-06-05T15:44:57Z</updated_at>
						<details>
							<account>
								<account_code>1</account_code>
								<first_name>Verena</first_name>
								<last_name>Example</last_name>
								<email>verena@test.com</email>
								<billing_info type="credit_card">
									<first_name>Verena</first_name>
									<last_name>Example</last_name>
									<address1>123 Main St.</address1>
									<city>San Francisco</city>
									<state>CA</state>
									<zip>94105</zip>
									<country>US</country>
									<phone nil="nil"></phone>
									<vat_number nil="nil"></vat_number>
									<card_type>Visa</card_type>
									<year type="integer">2017</year>
									<month type="integer">11</month>
									<first_six>411111</first_six>
									<last_four>1111</last_four>
								</billing_info>
							</account>
						</details>
					</transaction>
				</transactions>
				<a name="refund" href="https://blacklighttest.recurly.com/v2/invoices/5558/refund" method="post"/>
			</invoice>`)
	})

	resp, invoice, err := client.Invoices.Get(5558)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected get invoice to return OK")
	}

	ts, _ := time.Parse(recurly.DateTimeFormat, "2018-06-05T15:44:57Z")
	if diff := cmp.Diff(invoice, &recurly.Invoice{
		XMLName:     xml.Name{Local: "invoice"},
		AccountCode: "1",
		Address: recurly.Address{
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
		CreatedAt:        recurly.NewTime(ts),
		UpdatedAt:        recurly.NewTime(ts),
		ClosedAt:         recurly.NewTimeFromString("2018-06-05T15:44:57Z"),
		DueOn:            recurly.NewTimeFromString("2018-06-05T15:44:57Z"),
		Type:             "charge",
		Origin:           "purchase",
		TaxRate:          float64(0),
		NetTerms:         recurly.NewInt(0),
		CollectionMethod: "automatic",
		LineItems: []recurly.Adjustment{
			{
				AccountCode:       "1",
				InvoiceNumber:     5558,
				SubscriptionUUID:  "453f6aa0995e2d52c0d3e6453e9341da",
				UUID:              "626db120a84102b1809909071c701c60",
				State:             "invoiced",
				Description:       "License",
				ProductCode:       "license",
				Origin:            "add_on",
				UnitAmountInCents: 150000,
				Quantity:          1,
				TaxInCents:        0,
				TotalInCents:      150000,
				Currency:          "USD",
				Taxable:           recurly.NewBool(false),
				StartDate:         recurly.NewTimeFromString("2018-06-05T15:44:56Z"),
				EndDate:           recurly.NewTimeFromString("2018-07-05T15:44:56Z"),
				CreatedAt:         recurly.NewTimeFromString("2018-06-05T15:44:57Z"),
				UpdatedAt:         recurly.NewTimeFromString("2018-06-05T15:44:57Z"),
			},
			{
				AccountCode:       "1",
				InvoiceNumber:     5558,
				SubscriptionUUID:  "453f6aa0995e2d52c0d3e6453e9341da",
				UUID:              "453f6aa1473a0620e4411a4fc88122cf",
				State:             "invoiced",
				Description:       "Domains",
				ProductCode:       "domains",
				Origin:            "add_on",
				UnitAmountInCents: 1500,
				Quantity:          2,
				TaxInCents:        0,
				TotalInCents:      3000,
				Currency:          "USD",
				Taxable:           recurly.NewBool(false),
				StartDate:         recurly.NewTimeFromString("2018-06-05T15:44:56Z"),
				EndDate:           recurly.NewTimeFromString("2018-07-05T15:44:56Z"),
				CreatedAt:         recurly.NewTimeFromString("2018-06-05T15:44:57Z"),
				UpdatedAt:         recurly.NewTimeFromString("2018-06-05T15:44:57Z"),
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
				AVSResult: recurly.AVSResult{
					recurly.TransactionResult{
						Code:    "D",
						Message: "Street address and postal code match.",
					},
				},
				CreatedAt: recurly.NewTimeFromString("2018-06-05T15:44:56Z"),
				Account: recurly.Account{
					XMLName:   xml.Name{Local: "account"},
					Code:      "1",
					FirstName: "Verena",
					LastName:  "Example",
					Email:     "verena@test.com",
					BillingInfo: &recurly.Billing{
						XMLName:   xml.Name{Local: "billing_info"},
						FirstName: "Verena",
						LastName:  "Example",
						Address:   "123 Main St.",
						City:      "San Francisco",
						State:     "CA",
						Zip:       "94105",
						Country:   "US",
						CardType:  "Visa",
						Year:      2017,
						Month:     11,
						FirstSix:  411111,
						LastFour:  "1111",
					},
				},
			},
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestInvoices_Get_ErrNotFound(t *testing.T) {
	setup()
	defer teardown()

	var invoked bool
	mux.HandleFunc("/v2/invoices/1402", func(w http.ResponseWriter, r *http.Request) {
		invoked = true
		w.WriteHeader(http.StatusNotFound)
	})

	_, invoice, err := client.Invoices.Get(1402)
	if !invoked {
		t.Fatal("handler not invoked")
	} else if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if invoice != nil {
		t.Fatalf("expected invoice to be nil: %#v", invoice)
	}
}

// Ensures transactions are ordered by created at date.
func TestInvoices_Get_TransactionsOrder(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/invoices/1402", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
		<invoice href="https://your-subdomain.recurly.com/v2/invoices/1402">
			<transactions type="array">
				<transaction href="https://your-subdomain.recurly.com/v2/transactions/20150611" type="credit_card">
					<uuid>20150611</uuid>
					<created_at type="datetime">2015-06-11T15:25:06Z</created_at>
				</transaction>
				<transaction href="https://your-subdomain.recurly.com/v2/transactions/20160101" type="credit_card">
					<uuid>20160101</uuid>
					<created_at type="datetime">2016-01-01T15:25:06Z</created_at>
				</transaction>
				<transaction href="https://your-subdomain.recurly.com/v2/transactions/20150609" type="credit_card">
					<uuid>20150609</uuid>
					<created_at type="datetime">2015-06-09T15:25:06Z</created_at>
				</transaction>
				<transaction href="https://your-subdomain.recurly.com/v2/transactions/20150610" type="credit_card">
					<uuid>20150610</uuid>
					<created_at type="datetime">2015-06-10T15:25:06Z</created_at>
				</transaction>
			</transactions>
		</invoice>`)
	})

	resp, invoice, err := client.Invoices.Get(1402)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected get invoice to return OK")
	}

	// Verify transactions are in the correct order.
	if invoice.Transactions[0].UUID != "20150609" { // June 09 2015
		t.Fatalf("unexpected uuid(0): %s", invoice.Transactions[0].UUID)
	} else if invoice.Transactions[1].UUID != "20150610" { // June 10 2015
		t.Fatalf("unexpected uuid(1): %s", invoice.Transactions[1].UUID)
	} else if invoice.Transactions[2].UUID != "20150611" { // June 11 2015
		t.Fatalf("unexpected uuid(2): %s", invoice.Transactions[2].UUID)
	} else if invoice.Transactions[3].UUID != "20160101" { // Jan 01 2016
		t.Fatalf("unexpected uuid(3): %s", invoice.Transactions[3].UUID)
	}
}

func TestInvoices_GetPDF(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/invoices/1402", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		} else if r.Header.Get("Accept") != "application/pdf" {
			t.Fatalf("unexpected Accept heading: %s", r.Header.Get("Accept"))
		} else if r.Header.Get("Accept-Language") != "English" {
			t.Fatalf("unexpected Accept-Language header: %s", r.Header.Get("Accept-Language"))
		}

		w.WriteHeader(200)
		fmt.Fprint(w, "binary pdf text")
	})

	resp, pdf, err := client.Invoices.GetPDF(1402, "")
	if err != nil {
		t.Fatalf("error occurred making API call. Err: %s", err)
	} else if resp.IsError() {
		t.Fatal("expected get invoice to return OK")
	}

	expected := bytes.NewBufferString("binary pdf text")
	if !bytes.Equal(expected.Bytes(), pdf.Bytes()) {
		t.Fatalf("unexpected bytes: have=%v want %v", expected, pdf)
	}
}

func TestInvoices_GetPDFLanguage(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/invoices/1402", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		} else if r.Header.Get("Accept") != "application/pdf" {
			t.Fatalf("unexpected Accept heading: %s", r.Header.Get("Accept"))
		} else if r.Header.Get("Accept-Language") != "French" {
			t.Fatalf("unexpected Accept-Language header: %s", r.Header.Get("Accept-Language"))
		}

		w.WriteHeader(200)
		fmt.Fprint(w, "binary pdf text")
	})

	resp, pdf, err := client.Invoices.GetPDF(1402, "French")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected get invoice to return OK")
	}

	expected := bytes.NewBufferString("binary pdf text")
	if !bytes.Equal(expected.Bytes(), pdf.Bytes()) {
		t.Fatalf("unexpected bytes: have=%v want %v", expected, pdf)
	}
}

func TestInvoices_Preview(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/1/invoices/preview", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(201)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><invoice_collection>
			<charge_invoice href="">
			</charge_invoice>
			<credit_invoices type="array">
			</credit_invoices>
		</invoice_collection>`)
	})

	resp, _, err := client.Invoices.Preview("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected create invoice to return OK")
	}
}

func TestInvoices_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/10/invoices", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(201)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><invoice_collection>
			<charge_invoice href="https://your-subdomain.recurly.com/v2/invoices/1016">
				<account href="https://your-subdomain.recurly.com/v2/accounts/1"/>
				<address>
				<address1>123 Main St.</address1>
				<address2 nil="nil"/>
				<city>San Francisco</city>
				<state>CA</state>
				<zip>94105</zip>
				<country>US</country>
				<phone nil="nil"/>
				</address>
				<uuid>43adb97640cc05dee0b10042e596307f</uuid>
				<state>pending</state>
				<invoice_number_prefix/>
				<invoice_number type="integer">1016</invoice_number>
				<vat_number nil="nil"/>
				<tax_in_cents type="integer">425</tax_in_cents>
				<total_in_cents type="integer">3425</total_in_cents>
				<currency>USD</currency>
				<created_at type="datetime">2018-03-19T15:43:41Z</created_at>
				<updated_at type="datetime">2018-03-19T15:43:41Z</updated_at>
				<attempt_next_collection_at type="datetime">2018-03-20T15:43:41Z</attempt_next_collection_at>
				<closed_at nil="nil"/>
				<customer_notes nil="nil"/>
				<recovery_reason nil="nil"/>
				<subtotal_before_discount_in_cents type="integer">5000</subtotal_before_discount_in_cents>
				<subtotal_in_cents type="integer">3000</subtotal_in_cents>
				<discount_in_cents type="integer">17</discount_in_cents>
				<due_on type="datetime">2018-03-20T15:43:41Z</due_on>
				<net_terms type="integer">0</net_terms>
				<collection_method>manual</collection_method>
				<po_number nil="nil"/>
				<terms_and_conditions nil="nil"/>
				<tax_type>usst</tax_type>
				<tax_region>CA</tax_region>
				<line_items type="array">
				<adjustment href="https://your-subdomain.recurly.com/v2/adjustments/43adb5a639dc950ff620de42e6be4141" type="charge">
					<!-- Detail. -->
				</adjustment>
				</line_items>
				<transactions type="array">
				</transactions>
				<a name="mark_successful" href="https://your-subdomain.recurly.com/v2/invoices/1016/mark_successful" method="put"/>
				<a name="mark_failed" href="https://your-subdomain.recurly.com/v2/invoices/1016/mark_failed" method="put"/>
			</charge_invoice>
			<credit_invoices type="array">
			</credit_invoices>
			</invoice_collection>`)
	})

	resp, invoice, err := client.Invoices.Create("10", recurly.Invoice{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected create invoice to return OK")
	} else if diff := cmp.Diff(invoice, &recurly.Invoice{
		XMLName:     xml.Name{Local: "invoice"},
		AccountCode: "1",
		Address: recurly.Address{
			Address: "123 Main St.",
			City:    "San Francisco",
			State:   "CA",
			Zip:     "94105",
			Country: "US",
		},
		UUID:                    "43adb97640cc05dee0b10042e596307f",
		State:                   recurly.ChargeInvoiceStatePending,
		InvoiceNumber:           1016,
		DiscountInCents:         17,
		SubtotalInCents:         3000,
		TaxInCents:              425,
		TotalInCents:            3425,
		Currency:                "USD",
		DueOn:                   recurly.NewTimeFromString("2018-03-20T15:43:41Z"),
		CreatedAt:               recurly.NewTimeFromString("2018-03-19T15:43:41Z"),
		UpdatedAt:               recurly.NewTimeFromString("2018-03-19T15:43:41Z"),
		AttemptNextCollectionAt: recurly.NewTimeFromString("2018-03-20T15:43:41Z"),
		TaxType:                 "usst",
		TaxRegion:               "CA",
		NetTerms:                recurly.NewInt(0),
		CollectionMethod:        recurly.CollectionMethodManual,
		LineItems:               []recurly.Adjustment{{}},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestInvoices_Create_Params(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/10/invoices", func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		defer r.Body.Close()
		if !bytes.Equal(b, []byte("<invoice><po_number>ABC</po_number><net_terms>30</net_terms><collection_method>COLLECTION_METHOD</collection_method><terms_and_conditions>TERMS</terms_and_conditions><customer_notes>CUSTOMER_NOTES</customer_notes><vat_reverse_charge_notes>VAT_REVERSE_CHARGE_NOTES</vat_reverse_charge_notes></invoice>")) {
			t.Fatalf("unexpected input: %s", string(b))
		}
		w.WriteHeader(201)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><invoice_collection>
			<charge_invoice href="">
			<account href="https://your-subdomain.recurly.com/v2/accounts/1"/>
			<uuid>43adfe52c21cbb221557a24940bcd7e5</uuid>
			<state>pending</state>
			</charge_invoice>
			<credit_invoices type="array">
			</credit_invoices>
		</invoice_collection>`)
	})

	// Fields ordered in same order as struct xml tags, XML above in same order
	// for equality check.
	resp, _, err := client.Invoices.Create("10", recurly.Invoice{
		PONumber:              "ABC",
		NetTerms:              recurly.NewInt(30),
		CollectionMethod:      "COLLECTION_METHOD",
		TermsAndConditions:    "TERMS",
		CustomerNotes:         "CUSTOMER_NOTES",
		VatReverseChargeNotes: "VAT_REVERSE_CHARGE_NOTES",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected create invoice to return OK")
	}
}

func TestInvoices_Collect(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/invoices/1010/collect", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(201)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><invoice></invoice>`)
	})

	resp, _, err := client.Invoices.Collect(1010)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected create invoice to return OK")
	}
}

func TestInvoices_Collect_ErrBadRequest(t *testing.T) {
	setup()
	defer teardown()

	var invoked bool
	mux.HandleFunc("/v2/invoices/1010/collect", func(w http.ResponseWriter, r *http.Request) {
		invoked = true
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><error></error>`)
	})

	_, invoice, err := client.Invoices.Collect(1010)
	if !invoked {
		t.Fatal("handler not invoked")
	} else if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if invoice != nil {
		t.Fatal("unexpected invoice")
	}
}

func TestInvoices_MarkPaid(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/invoices/1402/mark_successful", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><invoice></invoice>`)
	})

	resp, _, err := client.Invoices.MarkPaid(1402)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected create invoice to return OK")
	}
}

func TestInvoices_MarkFailed(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/invoices/1402/mark_failed", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><invoice_collection>
			<charge_invoice href="">
			</charge_invoice>
			<credit_invoices type="array">
			</credit_invoices>
		</invoice_collection>`)
	})

	resp, _, err := client.Invoices.MarkFailed(1402)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected create invoice to return OK")
	}
}

func TestInvoices_RefundVoidOpenAmount(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/invoices/1010/refund", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(201)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><invoice></invoice>`)
	})

	resp, _, err := client.Invoices.RefundVoidOpenAmount(1010, 100, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected create open amount refund to return OK")
	}
}

func TestInvoices_VoidCreditInvoice(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/invoices/1010/void", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(201)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><invoice></invoice>`)
	})

	resp, _, err := client.Invoices.VoidCreditInvoice(1010)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected void credit invoice to return OK")
	}
}

func TestInvoices_RefundVoidOpenAmount_Params(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/invoices/1010/refund", func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		defer r.Body.Close()
		if !bytes.Equal(b, []byte("<invoice><amount_in_cents>100</amount_in_cents><refund_method>credit_first</refund_method></invoice>")) {
			t.Fatalf("unexpected input: %s", string(b))
		}
		w.WriteHeader(201)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><invoice></invoice>`)
	})

	// Fields ordered in same order as struct xml tags, XML above in same order
	// for equality check.
	resp, _, err := client.Invoices.RefundVoidOpenAmount(1010, 100, recurly.VoidRefundMethodCreditFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected create open amount refund to return OK")
	}
}

func TestInvoices_RecordPayment(t *testing.T) {
	setup()
	defer teardown()

	var invoked bool
	mux.HandleFunc("/v2/invoices/1402/transactions", func(w http.ResponseWriter, r *http.Request) {
		invoked = true
		if r.Method != "POST" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		defer r.Body.Close()
		if !bytes.Equal(b, []byte("<transaction><payment_method>check</payment_method><collected_at>2017-01-03T00:00:00Z</collected_at><amount_in_cents>1000</amount_in_cents><description>Paid with a check</description></transaction>")) {
			t.Fatalf("unexpected input: %s", string(b))
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><transaction></transaction>`)
	})

	date := time.Date(2017, 1, 3, 0, 0, 0, 0, time.UTC)
	resp, _, err := client.Invoices.RecordPayment(recurly.OfflinePayment{
		InvoiceNumber: 1402,
		PaymentMethod: recurly.PaymentMethodCheck,
		Amount:        1000,
		CollectedAt:   &date,
		Description:   "Paid with a check",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected create invoice to return OK")
	} else if !invoked {
		t.Fatal("handler not invoked")
	}
}
