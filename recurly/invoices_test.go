package recurly

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"testing"
	"time"
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
        		<subscription href="https://your-subdomain.recurly.com/v2/subscriptions/17caaca1716f33572edc8146e0aaefde"/>
        		<original_invoice href="https://your-subdomain.recurly.com/v2/invoices/938571" />
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

	resp, invoices, err := client.Invoices.List(Params{"per_page": 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected list invoices to return OK")
	} else if resp.Request.URL.Query().Get("per_page") != "1" {
		t.Fatalf("expected per_page parameter of 1, given %s", resp.Request.URL.Query().Get("per_page"))
	} else if !reflect.DeepEqual(invoices, []Invoice{{
		XMLName: xml.Name{Local: "invoice"},
		Account: href{
			HREF: "https://your-subdomain.recurly.com/v2/accounts/1",
			Code: "1",
		},
		Address: Address{
			Address: "400 Alabama St.",
			City:    "San Francisco",
			State:   "CA",
			Zip:     "94110",
			Country: "US",
		},
		Subscription: href{
			HREF: "https://your-subdomain.recurly.com/v2/subscriptions/17caaca1716f33572edc8146e0aaefde",
			Code: "17caaca1716f33572edc8146e0aaefde",
		},
		OriginalInvoice: href{
			HREF: "https://your-subdomain.recurly.com/v2/invoices/938571",
			Code: "938571",
		},
		UUID:             "421f7b7d414e4c6792938e7c49d552e9",
		State:            InvoiceStateOpen,
		InvoiceNumber:    1005,
		SubtotalInCents:  1200,
		TaxInCents:       0,
		TotalInCents:     1200,
		Currency:         "USD",
		CreatedAt:        newTimeFromString("2011-08-25T12:00:00Z"),
		TaxType:          "usst",
		TaxRegion:        "CA",
		TaxRate:          float64(0),
		NetTerms:         NewInt(0),
		CollectionMethod: "automatic",
		LineItems: []Adjustment{
			Adjustment{
				XMLName: xml.Name{Local: "adjustment"},
				Account: href{
					HREF: "https://your-subdomain.recurly.com/v2/accounts/100",
					Code: "100",
				},
				Invoice: href{
					HREF: "https://your-subdomain.recurly.com/v2/invoices/1108",
					Code: "1108",
				},
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
				Taxable:                NewBool(false),
				TaxExempt:              NewBool(false),
				StartDate:              newTimeFromString("2011-08-31T03:30:00Z"),
				CreatedAt:              newTimeFromString("2011-08-31T03:30:00Z"),
			},
		},
	}}) {
		t.Fatalf("unexpected invoices: %v", invoices)
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
        		<subscription href="https://your-subdomain.recurly.com/v2/subscriptions/17caaca1716f33572edc8146e0aaefde"/>
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

	resp, invoices, err := client.Invoices.ListAccount("1", Params{"per_page": 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected list invoices to return OK")
	} else if pp := resp.Request.URL.Query().Get("per_page"); pp != "1" {
		t.Fatalf("unexpected per_page: %s", pp)
	} else if !reflect.DeepEqual(invoices, []Invoice{Invoice{
		XMLName: xml.Name{Local: "invoice"},
		Account: href{
			HREF: "https://your-subdomain.recurly.com/v2/accounts/1",
			Code: "1",
		},
		Address: Address{
			Address: "400 Alabama St.",
			City:    "San Francisco",
			State:   "CA",
			Zip:     "94110",
			Country: "US",
		},
		Subscription: href{
			HREF: "https://your-subdomain.recurly.com/v2/subscriptions/17caaca1716f33572edc8146e0aaefde",
			Code: "17caaca1716f33572edc8146e0aaefde",
		},
		UUID:             "421f7b7d414e4c6792938e7c49d552e9",
		State:            InvoiceStateOpen,
		InvoiceNumber:    1005,
		SubtotalInCents:  1200,
		TaxInCents:       0,
		TotalInCents:     1200,
		Currency:         "USD",
		CreatedAt:        newTimeFromString("2011-08-25T12:00:00Z"),
		TaxType:          "usst",
		TaxRegion:        "CA",
		TaxRate:          float64(0),
		NetTerms:         NewInt(0),
		CollectionMethod: "automatic",
		LineItems: []Adjustment{
			Adjustment{
				XMLName: xml.Name{Local: "adjustment"},
				Account: href{
					HREF: "https://your-subdomain.recurly.com/v2/accounts/100",
					Code: "100",
				},
				Invoice: href{
					HREF: "https://your-subdomain.recurly.com/v2/invoices/1108",
					Code: "1108",
				},
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
				Taxable:                NewBool(false),
				TaxExempt:              NewBool(false),
				StartDate:              newTimeFromString("2011-08-31T03:30:00Z"),
				CreatedAt:              newTimeFromString("2011-08-31T03:30:00Z"),
			},
		},
	}}) {
		t.Fatalf("unexpected invoices: %v", invoices)
	}
}

func TestInvoices_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/invoices/1402", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
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
    		<subscription href="https://your-subdomain.recurly.com/v2/subscriptions/17caaca1716f33572edc8146e0aaefde"/>
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
                <transaction href="https://your-subdomain.recurly.com/v2/transactions/a13acd8fe4294916b79aec87b7ea441f" type="credit_card">
                    <account href="https://your-subdomain.recurly.com/v2/accounts/1"/>
                    <invoice href="https://your-subdomain.recurly.com/v2/invoices/1108"/>
                    <subscription href="https://your-subdomain.recurly.com/v2/subscriptions/17caaca1716f33572edc8146e0aaefde"/>
                    <uuid>a13acd8fe4294916b79aec87b7ea441f</uuid>
                    <action>purchase</action>
                    <amount_in_cents type="integer">1000</amount_in_cents>
                    <tax_in_cents type="integer">0</tax_in_cents>
                    <currency>USD</currency>
                    <status>success</status>
                    <payment_method>credit_card</payment_method>
                    <reference>5416477</reference>
                    <source>subscription</source>
                    <recurring type="boolean">true</recurring>
                    <test type="boolean">true</test>
                    <voidable type="boolean">true</voidable>
                    <refundable type="boolean">true</refundable>
                    <ip_address>127.0.0.1</ip_address>
                    <cvv_result code="M">Match</cvv_result>
                    <avs_result code="D">Street address and postal code match.</avs_result>
                    <avs_result_street nil="nil"/>
                    <avs_result_postal nil="nil"/>
                    <created_at type="datetime">2015-06-10T15:25:06Z</created_at>
                    <details>
                        <account>
                            <account_code>1</account_code>
                            <first_name>Verena</first_name>
                            <last_name>Example</last_name>
                            <company nil="nil"/>
                            <email>verena@test.com</email>
                            <billing_info type="credit_card">
                                <first_name>Verena</first_name>
                                <last_name>Example</last_name>
                                <address1>123 Main St.</address1>
                                <address2 nil="nil"/>
                                <city>San Francisco</city>
                                <state>CA</state>
                                <zip>94105</zip>
                                <country>US</country>
                                <phone nil="nil"/>
                                <vat_number nil="nil"/>
                                <card_type>Visa</card_type>
                                <year type="integer">2017</year>
                                <month type="integer">11</month>
                                <first_six>411111</first_six>
                                <last_four>1111</last_four>
                            </billing_info>
                        </account>
                    </details>
                    <a name="refund" href="https://your-subdomain.recurly.com/v2/transactions/a13acd8fe4294916b79aec87b7ea441f" method="delete"/>
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

	ts, _ := time.Parse(datetimeFormat, "2011-08-25T12:00:00Z")
	if !reflect.DeepEqual(invoice, &Invoice{
		XMLName: xml.Name{Local: "invoice"},
		Account: href{
			HREF: "https://your-subdomain.recurly.com/v2/accounts/1",
			Code: "1",
		},
		Address: Address{
			Address: "400 Alabama St.",
			City:    "San Francisco",
			State:   "CA",
			Zip:     "94110",
			Country: "US",
		},
		Subscription: href{
			HREF: "https://your-subdomain.recurly.com/v2/subscriptions/17caaca1716f33572edc8146e0aaefde",
			Code: "17caaca1716f33572edc8146e0aaefde",
		},
		UUID:             "421f7b7d414e4c6792938e7c49d552e9",
		State:            InvoiceStateOpen,
		InvoiceNumber:    1005,
		SubtotalInCents:  1200,
		TaxInCents:       0,
		TotalInCents:     1200,
		Currency:         "USD",
		CreatedAt:        NewTime(ts),
		TaxType:          "usst",
		TaxRegion:        "CA",
		TaxRate:          float64(0),
		NetTerms:         NewInt(0),
		CollectionMethod: "automatic",
		LineItems: []Adjustment{
			Adjustment{
				XMLName: xml.Name{Local: "adjustment"},
				Account: href{
					HREF: "https://your-subdomain.recurly.com/v2/accounts/100",
					Code: "100",
				},
				Invoice: href{
					HREF: "https://your-subdomain.recurly.com/v2/invoices/1108",
					Code: "1108",
				},
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
				Taxable:                NewBool(false),
				TaxExempt:              NewBool(false),
				StartDate:              newTimeFromString("2011-08-31T03:30:00Z"),
				CreatedAt:              newTimeFromString("2011-08-31T03:30:00Z"),
			},
		},
		Transactions: []Transaction{
			Transaction{
				XMLName: xml.Name{Local: "transaction"},
				Invoice: href{
					HREF: "https://your-subdomain.recurly.com/v2/invoices/1108",
					Code: "1108",
				},
				Subscription: href{
					HREF: "https://your-subdomain.recurly.com/v2/subscriptions/17caaca1716f33572edc8146e0aaefde",
					Code: "17caaca1716f33572edc8146e0aaefde",
				},
				UUID:          "a13acd8fe4294916b79aec87b7ea441f",
				Action:        "purchase",
				AmountInCents: 1000,
				TaxInCents:    0,
				Currency:      "USD",
				Status:        "success",
				PaymentMethod: "credit_card",
				Reference:     "5416477",
				Source:        "subscription",
				Recurring:     NewBool(true),
				Test:          true,
				Voidable:      NewBool(true),
				Refundable:    NewBool(true),
				IPAddress:     net.ParseIP("127.0.0.1"),
				CVVResult: CVVResult{
					transactionResult{
						Code:    "M",
						Message: "Match",
					},
				},
				AVSResult: AVSResult{
					transactionResult{
						Code:    "D",
						Message: "Street address and postal code match.",
					},
				},
				CreatedAt: newTimeFromString("2015-06-10T15:25:06Z"),
				Account: Account{
					XMLName:   xml.Name{Local: "account"},
					Code:      "1",
					FirstName: "Verena",
					LastName:  "Example",
					Email:     "verena@test.com",
					BillingInfo: &Billing{
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
						LastFour:  1111,
					},
				},
			},
		},
	}) {
		t.Fatalf("unexpected invoice: %v", invoice)
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
	if !reflect.DeepEqual(expected, pdf) {
		t.Fatalf("unexpected pdf invoice: %s", pdf)
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
	if !reflect.DeepEqual(expected, pdf) {
		t.Fatalf("unexpected pdf: %v", pdf)
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
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><invoice></invoice>`)
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
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><invoice></invoice>`)
	})

	resp, _, err := client.Invoices.Create("10", Invoice{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected create invoice to return OK")
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

	resp, _, err := client.Invoices.MarkAsPaid(1402)
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
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><invoice></invoice>`)
	})

	resp, _, err := client.Invoices.MarkAsFailed(1402)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if resp.IsError() {
		t.Fatal("expected create invoice to return OK")
	}
}
