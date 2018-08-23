package recurly_test

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/launchpadcentral/recurly"
	"github.com/google/go-cmp/cmp"
)

// TestTransactionEncoding ensures structs are encoded to XML properly.
// Because Recurly supports partial updates, it's important that only defined
// fields are handled properly -- including types like booleans and integers which
// have zero values that we want to send.
func TestTransactions_Encoding(t *testing.T) {
	var transaction recurly.Transaction
	buf, err := xml.Marshal(transaction)
	if err != nil {
		t.Fatal(err)
	}

	if string(buf) != "<transaction><amount_in_cents>0</amount_in_cents><currency></currency><account></account></transaction>" {
		t.Fatalf("unexpected encoding: %s", string(buf))
	}
}

func TestTransactions_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/transactions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
        <transactions type="array">
        	<transaction href="https://your-subdomain.recurly.com/v2/transactions/a13acd8fe4294916b79aec87b7ea441f" type="credit_card">
        		<account href="https://your-subdomain.recurly.com/v2/accounts/1"/>
        		<invoice href="https://your-subdomain.recurly.com/v2/invoices/1108"/>
        		<subscriptions href="https://your-subdomain.recurly.com/v2/transactions/a13acd8fe4294916b79aec87b7ea441f/subscriptions"/>
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
        </transactions>`)
	})

	r, transactions, err := client.Transactions.List(recurly.Params{"per_page": 1})
	if err != nil {
		t.Fatalf("TestTransactionsList Error: Error occurred making API call. Err: %s", err)
	} else if r.IsError() {
		t.Fatal("TestTransactionsList Error: Expected list transactions to return OK")
	} else if pp := r.Request.URL.Query().Get("per_page"); pp != "1" {
		t.Fatalf("unexpected per_page: %s", pp)
	}

	if diff := cmp.Diff(transactions, []recurly.Transaction{
		{
			InvoiceNumber: 1108,
			UUID:          "a13acd8fe4294916b79aec87b7ea441f",
			Action:        "purchase",
			AmountInCents: 1000,
			TaxInCents:    0,
			Currency:      "USD",
			Status:        "success",
			PaymentMethod: "credit_card",
			Reference:     "5416477",
			Source:        "subscription",
			Recurring:     recurly.NewBool(true),
			Test:          true,
			Voidable:      recurly.NewBool(true),
			Refundable:    recurly.NewBool(true),
			IPAddress:     net.ParseIP("127.0.0.1"),
			CVVResult: recurly.CVVResult{
				recurly.TransactionResult{
					Code:    "M",
					Message: "Match",
				},
			},
			AVSResult: recurly.AVSResult{
				recurly.TransactionResult{
					Code:    "D",
					Message: "Street address and postal code match.",
				},
			},
			CreatedAt: recurly.NewTime(time.Date(2015, time.June, 10, 15, 25, 6, 0, time.UTC)),
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
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestTransactions_ListAccount(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/1/transactions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
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
        </transactions>`)
	})

	r, transactions, err := client.Transactions.ListAccount("1", recurly.Params{"per_page": 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected list for account transactions to return OK")
	} else if pp := r.Request.URL.Query().Get("per_page"); pp != "1" {
		t.Fatalf("unexpected per_page: %s", pp)
	}

	if diff := cmp.Diff(transactions, []recurly.Transaction{
		{
			InvoiceNumber: 1108,
			UUID:          "a13acd8fe4294916b79aec87b7ea441f",
			Action:        "purchase",
			AmountInCents: 1000,
			TaxInCents:    0,
			Currency:      "USD",
			Status:        "success",
			PaymentMethod: "credit_card",
			Reference:     "5416477",
			Source:        "subscription",
			Recurring:     recurly.NewBool(true),
			Test:          true,
			Voidable:      recurly.NewBool(true),
			Refundable:    recurly.NewBool(true),
			IPAddress:     net.ParseIP("127.0.0.1"),
			CVVResult: recurly.CVVResult{
				recurly.TransactionResult{
					Code:    "M",
					Message: "Match",
				},
			},
			AVSResult: recurly.AVSResult{
				recurly.TransactionResult{
					Code:    "D",
					Message: "Street address and postal code match.",
				},
			},
			CreatedAt: recurly.NewTime(time.Date(2015, time.June, 10, 15, 25, 6, 0, time.UTC)),
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
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestTransactions_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/transactions/a13acd8fe4294916b79aec87b7ea441f", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
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
    		<description>Order #717</description>
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
    	</transaction>`)
	})

	r, transaction, err := client.Transactions.Get("a13acd8f-e4294916b79aec87-b7ea441f") // UUID contains dashes
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected get transaction to return OK")
	}

	if diff := cmp.Diff(transaction, &recurly.Transaction{
		InvoiceNumber: 1108,
		UUID:          "a13acd8fe4294916b79aec87b7ea441f", // UUID has been sanitized
		Action:        "purchase",
		AmountInCents: 1000,
		TaxInCents:    0,
		Currency:      "USD",
		Status:        "success",
		Description:   "Order #717",
		PaymentMethod: "credit_card",
		Reference:     "5416477",
		Source:        "subscription",
		Recurring:     recurly.NewBool(true),
		Test:          true,
		Voidable:      recurly.NewBool(true),
		Refundable:    recurly.NewBool(true),
		IPAddress:     net.ParseIP("127.0.0.1"),
		CVVResult: recurly.CVVResult{
			recurly.TransactionResult{
				Code:    "M",
				Message: "Match",
			},
		},
		AVSResult: recurly.AVSResult{
			recurly.TransactionResult{
				Code:    "D",
				Message: "Street address and postal code match.",
			},
		},
		CreatedAt: recurly.NewTime(time.Date(2015, time.June, 10, 15, 25, 6, 0, time.UTC)),
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
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestTransactions_Get_ErrNotFound(t *testing.T) {
	setup()
	defer teardown()

	var invoked bool
	mux.HandleFunc("/v2/transactions/a13acd8fe4294916b79aec87b7ea441f", func(w http.ResponseWriter, r *http.Request) {
		invoked = true
		w.WriteHeader(http.StatusNotFound)
	})

	_, transaction, err := client.Transactions.Get("a13acd8fe4294916b79aec87b7ea441f")
	if !invoked {
		t.Fatal("handler not invoked")
	} else if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if transaction != nil {
		t.Fatalf("expected transaction to be nil: %#v", transaction)
	}
}

func TestTransactions_New(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/transactions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		defer r.Body.Close()
		expected := `<transaction><amount_in_cents>100</amount_in_cents><currency>USD</currency><description>description</description><product_code>code</product_code><account><account_code>25</account_code></account></transaction>`
		var given bytes.Buffer
		given.ReadFrom(r.Body)
		if expected != given.String() {
			t.Fatalf("unexpected input: %s", given.String())
		}

		w.WriteHeader(200)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><transaction></transaction>`)
	})

	r, _, err := client.Transactions.Create(
		recurly.Transaction{
			AmountInCents: 100,
			Currency:      "USD",
			Account: recurly.Account{
				Code: "25",
			},
			Description: "description",
			ProductCode: "code",
		})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if r.IsError() {
		t.Fatal("expected create transaction to return OK")
	}
}

func TestTransactions_Err_FraudCard(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/transactions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("unexpected method: %s", r.Method)
		}

		w.WriteHeader(422)
		fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
			<errors>
			  <transaction_error>
			    <error_code>fraud_gateway</error_code>
			    <error_category>fraud</error_category>
			    <merchant_message>The payment gateway declined the transaction due to fraud filters enabled in your gateway.</merchant_message>
			    <customer_message>The transaction was declined. Please use a different card, contact your bank, or contact support.</customer_message>
			    <gateway_error_code nil="nil"></gateway_error_code>
			  </transaction_error>
			  <error field="transaction.account.base" symbol="fraud_gateway">The transaction was declined. Please use a different card, contact your bank, or contact support.</error>
			  <transaction href="https://your-subdomain.recurly.com/v2/transactions/3054a79e4c3ab4699f95be455f8653bb" type="credit_card">
			    <account href="https://your-subdomain.recurly.com/v2/accounts/Seantwo@recurly.com"/>
			    <uuid>3054a79e4c3ab4699f95be455f8653bb</uuid>
			    <action>purchase</action>
			    <amount_in_cents type="integer">100</amount_in_cents>
			    <tax_in_cents type="integer">0</tax_in_cents>
			    <currency>USD</currency>
			    <status>declined</status>
			    <payment_method>credit_card</payment_method>
			    <transaction_error>
			      <error_code>fraud_gateway</error_code>
			      <error_category>fraud</error_category>
			      <merchant_message>The payment gateway declined the transaction due to fraud filters enabled in your gateway.</merchant_message>
			      <customer_message>The transaction was declined. Please use a different card, contact your bank, or contact support.</customer_message>
			      <gateway_error_code nil="nil"></gateway_error_code>
			    </transaction_error>
			    <details>
			    </details>
			  </transaction>
			</errors>`)
	})

	r, transaction, err := client.Transactions.Create(recurly.Transaction{
		AmountInCents: 100,
		Currency:      "USD",
		Account: recurly.Account{
			Code: "25",
			BillingInfo: &recurly.Billing{
				FirstName: "Verena",
				LastName:  "Example",
				Number:    4000000000000085,
				Month:     10,
				Year:      2020,
			},
		},
	})
	if err != nil {
		t.Fatalf("error occurred making API call. Err: %s", err)
	} else if r.IsOK() {
		t.Fatal("expected create fraudulent transaction to return error")
	} else if diff := cmp.Diff(transaction, &recurly.Transaction{
		UUID:          "3054a79e4c3ab4699f95be455f8653bb",
		Action:        "purchase",
		AmountInCents: 100,
		Currency:      "USD",
		Status:        "declined",
		PaymentMethod: "credit_card",
		TransactionError: &recurly.TransactionError{
			XMLName:         xml.Name{Local: "transaction_error"},
			ErrorCode:       "fraud_gateway",
			ErrorCategory:   "fraud",
			MerchantMessage: "The payment gateway declined the transaction due to fraud filters enabled in your gateway.",
			CustomerMessage: "The transaction was declined. Please use a different card, contact your bank, or contact support.",
		},
	}); diff != "" {
		t.Fatal(diff)
	}
}

func TestCVV(t *testing.T) {
	c := recurly.CVVResult{recurly.TransactionResult{Code: "M"}}
	if !c.IsMatch() {
		t.Fatalf("expected %q code to be match", "M")
	} else if c.IsNoMatch() || c.NotProcessed() || c.ShouldHaveBeenPresent() || c.UnableToProcess() {
		t.Fatalf("expected %q code to ONLY be match", "M")
	}

	c = recurly.CVVResult{recurly.TransactionResult{Code: "N"}}
	if !c.IsNoMatch() {
		t.Fatalf("expected %q code to not be a match", "N")
	} else if c.IsMatch() || c.NotProcessed() || c.ShouldHaveBeenPresent() || c.UnableToProcess() {
		t.Fatalf("expected %q code to ONLY be match", "N")
	}

	c = recurly.CVVResult{recurly.TransactionResult{Code: "P"}}
	if !c.NotProcessed() {
		t.Fatalf("expected %q code to not be a match", "P")
	} else if c.IsMatch() || c.IsNoMatch() || c.ShouldHaveBeenPresent() || c.UnableToProcess() {
		t.Fatalf("expected %q code to ONLY be match", "P")
	}

	c = recurly.CVVResult{recurly.TransactionResult{Code: "S"}}
	if !c.ShouldHaveBeenPresent() {
		t.Fatalf("expected %q code to not be a match", "S")
	} else if c.IsMatch() || c.IsNoMatch() || c.NotProcessed() || c.UnableToProcess() {
		t.Fatalf("expected %q code to ONLY be match", "S")
	}

	c = recurly.CVVResult{recurly.TransactionResult{Code: "U"}}
	if !c.UnableToProcess() {
		t.Fatalf("expected %q code to not be a match", "U")
	} else if c.IsMatch() || c.IsNoMatch() || c.NotProcessed() || c.ShouldHaveBeenPresent() {
		t.Fatalf("expected %q code to ONLY be match", "U")
	}
}
