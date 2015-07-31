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

// TestTransactionEncoding ensures structs are encoded to XML properly.
// Because Recurly supports partial updates, it's important that only defined
// fields are handled properly -- including types like booleans and integers which
// have zero values that we want to send.
func TestTransactionsEncoding(t *testing.T) {
	suite := []map[string]interface{}{
		map[string]interface{}{"struct": Transaction{}, "xml": "<transaction><amount_in_cents>0</amount_in_cents><currency></currency><details><account></account></details></transaction>"},
	}

	for _, s := range suite {
		buf := new(bytes.Buffer)
		err := xml.NewEncoder(buf).Encode(s["struct"])
		if err != nil {
			t.Errorf("TestTransactionEncoding Error: %s", err)
		}

		if buf.String() != s["xml"] {
			t.Errorf("TestTransactionEncoding Error: Expected %s, given %s", s["xml"], buf.String())
		}
	}
}

func TestTransactionsList(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/transactions", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("TestTransactionsList Error: Expected %s request, given %s", "GET", r.Method)
		}
		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?>
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

	r, transactions, err := client.Transactions.List(Params{"per_page": 1})
	if err != nil {
		t.Errorf("TestTransactionsList Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestTransactionsList Error: Expected list transactions to return OK")
	}

	if len(transactions) != 1 {
		t.Fatalf("TestTransactionsList Error: Expected 1 transaction returned, given %d", len(transactions))
	}

	if r.Request.URL.Query().Get("per_page") != "1" {
		t.Errorf("TestTransactionsList Error: Expected per_page parameter of 1, given %s", r.Request.URL.Query().Get("per_page"))
	}

	ts, _ := time.Parse(datetimeFormat, "2015-06-10T15:25:06Z")
	for _, given := range transactions {
		expected := Transaction{
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
			CVVResult: &TransactionResult{
				Code:    "M",
				Message: "Match",
			},
			AVSResult: &TransactionResult{
				Code:    "D",
				Message: "Street address and postal code match.",
			},
			CreatedAt: NewTime(ts),
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
		}

		if !reflect.DeepEqual(expected, given) {
			t.Errorf("TestTransactionsList Error: expected transaction to equal %#v, given %#v", expected, given)
		}
	}
}

func TestTransactionsListAccount(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/1/transactions", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("TestTransactionsListAccount Error: Expected %s request, given %s", "GET", r.Method)
		}
		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?>
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

	r, transactions, err := client.Transactions.ListAccount("1", Params{"per_page": 1})
	if err != nil {
		t.Errorf("TestTransactionsListAccount Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestTransactionsListAccount Error: Expected list for account transactions to return OK")
	}

	if len(transactions) != 1 {
		t.Fatalf("TestTransactionsListAccount Error: Expected 1 transaction returned, given %d", len(transactions))
	}

	if r.Request.URL.Query().Get("per_page") != "1" {
		t.Errorf("TestTransactionsListAccount Error: Expected per_page parameter of 1, given %s", r.Request.URL.Query().Get("per_page"))
	}

	ts, _ := time.Parse(datetimeFormat, "2015-06-10T15:25:06Z")
	for _, given := range transactions {
		expected := Transaction{
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
			CVVResult: &TransactionResult{
				Code:    "M",
				Message: "Match",
			},
			AVSResult: &TransactionResult{
				Code:    "D",
				Message: "Street address and postal code match.",
			},
			CreatedAt: NewTime(ts),
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
		}

		if !reflect.DeepEqual(expected, given) {
			t.Errorf("TestTransactionsListAccount Error: expected transaction to equal %#v, given %#v", expected, given)
		}
	}
}

func TestGetTransaction(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/transactions/a13acd8fe4294916b79aec87b7ea441f", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("TestGetTransaction Error: Expected %s request, given %s", "GET", r.Method)
		}
		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?>
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
    	</transaction>`)
	})

	r, a, err := client.Transactions.Get("a13acd8fe4294916b79aec87b7ea441f")
	if err != nil {
		t.Errorf("TestGetTransaction Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestGetTransaction Error: Expected get transaction to return OK")
	}

	ts, _ := time.Parse(datetimeFormat, "2015-06-10T15:25:06Z")
	expected := Transaction{
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
		CVVResult: &TransactionResult{
			Code:    "M",
			Message: "Match",
		},
		AVSResult: &TransactionResult{
			Code:    "D",
			Message: "Street address and postal code match.",
		},
		CreatedAt: NewTime(ts),
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
	}

	if !reflect.DeepEqual(expected, a) {
		t.Errorf("TestGetTransaction Error: expected account to equal %#v, given %#v", expected, a)
	}
}

func TestCreateTransaction(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/transactions", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("TestCreateTransaction Error: Expected %s request, given %s", "POST", r.Method)
		}
		expected := `<transaction><amount_in_cents>100</amount_in_cents><currency>USD</currency><account><account_code>25</account_code></account></transaction>`
		given := new(bytes.Buffer)
		given.ReadFrom(r.Body)
		if expected != given.String() {
			t.Errorf("TestCreateTransaction Error: Expected request body of %s, given %s", expected, given.String())
		}

		rw.WriteHeader(200)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?><transaction></transaction>`)
	})

	r, _, err := client.Transactions.Create(CreateTransaction{
		AmountInCents: 100,
		Currency:      "USD",
		Account: Account{
			Code: "25",
		},
	})
	if err != nil {
		t.Errorf("TestCreateTransaction Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsError() {
		t.Fatal("TestCreateTransaction Error: Expected create transaction to return OK")
	}
}

func TestCreateTransactionFraudCard(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/transactions", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("TestCreateTransactionFraudCard Error: Expected %s request, given %s", "POST", r.Method)
		}

		rw.WriteHeader(422)
		fmt.Fprint(rw, `<?xml version="1.0" encoding="UTF-8"?>
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
			    <reference>6223543</reference>
			    <source>transaction</source>
			    <recurring type="boolean">false</recurring>
			    <test type="boolean">true</test>
			    <voidable type="boolean">false</voidable>
			    <refundable type="boolean">false</refundable>
			    <ip_address>184.23.184.210</ip_address>
			    <transaction_error>
			      <error_code>fraud_gateway</error_code>
			      <error_category>fraud</error_category>
			      <merchant_message>The payment gateway declined the transaction due to fraud filters enabled in your gateway.</merchant_message>
			      <customer_message>The transaction was declined. Please use a different card, contact your bank, or contact support.</customer_message>
			      <gateway_error_code nil="nil"></gateway_error_code>
			    </transaction_error>
			    <cvv_result code="" nil="nil"></cvv_result>
			    <avs_result code="" nil="nil"></avs_result>
			    <avs_result_street nil="nil"></avs_result_street>
			    <avs_result_postal nil="nil"></avs_result_postal>
			    <created_at type="datetime">2015-07-31T20:45:01Z</created_at>
			    <details>
			      <account>
			        <account_code>1</account_code>
			        <first_name>Verena</first_name>
			        <last_name>Example</last_name>
			        <company></company>
			        <email></email>
			        <billing_info type="credit_card">
			          <first_name>Verena</first_name>
			          <last_name>Example</last_name>
			          <address1>123 Main St.</address1>
			          <address2></address2>
			          <city>San Francisco</city>
			          <state>CA</state>
			          <zip>94133</zip>
			          <country>US</country>
			          <phone nil="nil"></phone>
			          <vat_number nil="nil"></vat_number>
			          <card_type>Visa</card_type>
			          <year type="integer">2020</year>
			          <month type="integer">10</month>
			          <first_six>400000</first_six>
			          <last_four>0085</last_four>
			        </billing_info>
			      </account>
			    </details>
			  </transaction>
			</errors>`)
	})

	r, _, err := client.Transactions.Create(CreateTransaction{
		AmountInCents: 100,
		Currency:      "USD",
		Account: Account{
			Code: "25",
			BillingInfo: &Billing{
				FirstName: "Verena",
				LastName:  "Example",
				Number:    4000000000000085,
				Month:     10,
				Year:      2020,
			},
		},
	})
	if err != nil {
		t.Errorf("TestCreateTransactionFraudCard Error: Error occurred making API call. Err: %s", err)
	}

	if r.IsOK() {
		t.Fatal("TestCreateTransactionFraudCard Error: Expected create fraudulent transaction to return error")
	}

	expected := TransactionError{
		XMLName:         xml.Name{Local: "transaction_error"},
		ErrorCode:       "fraud_gateway",
		ErrorCategory:   "fraud",
		MerchantMessage: "The payment gateway declined the transaction due to fraud filters enabled in your gateway.",
		CustomerMessage: "The transaction was declined. Please use a different card, contact your bank, or contact support.",
	}

	if !reflect.DeepEqual(expected, r.TransactionError) {
		t.Errorf("TestCreateTransactionFraudCard Error: Expected transaction error of %+v, given %+v", expected, r.TransactionError)
	}
}
