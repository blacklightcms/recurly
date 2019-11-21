package recurly

import (
	"context"
	"encoding/xml"
	"fmt"
	"net"
	"net/http"
	"sort"
	"strconv"
)

// TransactionsService manages the interactions for transactions.
type TransactionsService interface {
	// List returns a pager to paginate transactions. PagerOptions are used to
	// optionally filter the results.
	//
	// https://dev.recurly.com/docs/list-transactions
	List(opts *PagerOptions) Pager

	// ListAccount returns a pager to paginate transactions for an account.
	// PagerOptions are used to optionally filter the results.
	//
	// https://dev.recurly.com/docs/list-accounts-transactions
	ListAccount(accountCode string, opts *PagerOptions) Pager

	// Get retrieves a transaction. If the transaction does not exist,
	// a nil transaction and nil error are returned.
	//
	// https://dev.recurly.com/docs/lookup-transaction
	Get(ctx context.Context, uuid string) (*Transaction, error)
}

// Transaction constants.
// https://docs.recurly.com/docs/transactions
const (
	TransactionStatusSuccess = "success"
	TransactionStatusFailed  = "failed"
	TransactionStatusVoid    = "void"
)

// Transaction is an individual transaction.
type Transaction struct {
	InvoiceNumber           int               // Read only
	OriginalTransactionUUID string            // Read only
	UUID                    string            `xml:"uuid,omitempty"` // Read only
	Action                  string            `xml:"action,omitempty"`
	AmountInCents           int               `xml:"amount_in_cents"`
	TaxInCents              int               `xml:"tax_in_cents,omitempty"`
	Currency                string            `xml:"currency"`
	Status                  string            `xml:"status,omitempty"`
	Description             string            `xml:"description,omitempty"`
	ProductCode             string            `xml:"-"` // Write only field, saved on the invoice line item but not the transaction
	PaymentMethod           string            `xml:"payment_method,omitempty"`
	Reference               string            `xml:"reference,omitempty"`
	Source                  string            `xml:"source,omitempty"`
	Recurring               NullBool          `xml:"recurring,omitempty"`
	Test                    bool              `xml:"test,omitempty"`
	Voidable                NullBool          `xml:"voidable,omitempty"`
	Refundable              NullBool          `xml:"refundable,omitempty"`
	IPAddress               net.IP            `xml:"ip_address,omitempty"`
	TransactionError        *TransactionError `xml:"transaction_error,omitempty"` // Read only
	CVVResult               CVVResult         `xml:"cvv_result,omitempty"`        // Read only
	AVSResult               AVSResult         `xml:"avs_result,omitempty"`        // Read only
	AVSResultStreet         string            `xml:"avs_result_street,omitempty"` // Read only
	AVSResultPostal         string            `xml:"avs_result_postal,omitempty"` // Read only
	CreatedAt               NullTime          `xml:"created_at,omitempty"`        // Read only
	Account                 Account           `xml:"details>account"`             // Read only
	GatewayType             string            `xml:"gateway_type,omitempty"`      // Read only
	Origin                  string            `xml:"origin,omitempty"`            // Read only
	Message                 string            `xml:"message,omitempty"`           // Read only
	ApprovalCode            string            `xml:"approval_code,omitempty"`     // Read only

}

// TransactionError is an error encounted from your payment gateway that
// recurly has standardized.
//
// https://dev.recurly.com/page/transaction-errors
type TransactionError struct {
	XMLName                   xml.Name `xml:"transaction_error"`
	ErrorCode                 string   `xml:"error_code,omitempty"`
	ErrorCategory             string   `xml:"error_category,omitempty"`
	MerchantMessage           string   `xml:"merchant_message,omitempty"`
	CustomerMessage           string   `xml:"customer_message,omitempty"`
	GatewayErrorCode          string   `xml:"gateway_error_code,omitempty"`
	ThreeDSecureActionTokenID string   `xml:"three_d_secure_action_token_id,omitempty"`
}

// UnmarshalXML unmarshals transactions and handles intermediary state during unmarshaling
// for types like href.
func (t *Transaction) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type transactionAlias Transaction
	var v struct {
		transactionAlias
		XMLName                 xml.Name `xml:"transaction"`
		InvoiceNumber           href     `xml:"invoice"`
		OriginalTransactionUUID href     `xml:"original_transaction"`
	}
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}

	*t = Transaction(v.transactionAlias)

	t.InvoiceNumber, _ = strconv.Atoi(v.InvoiceNumber.LastPartOfPath())
	t.OriginalTransactionUUID = v.OriginalTransactionUUID.LastPartOfPath()
	return nil
}

// CVVResult holds transaction results for CVV fields.
// https://www.chasepaymentech.com/card_verification_codes.html
type CVVResult struct {
	NullMarshal
	Code    string `xml:"code,attr"`
	Message string `xml:",innerxml"`
}

// AVSResult holds transaction results for address verification.
// http://developer.authorize.net/tools/errorgenerationguide/
type AVSResult struct {
	NullMarshal
	Code    string `xml:"code,attr"`
	Message string `xml:",innerxml"`
}

// Transactions is a sortable slice of Transaction.
type Transactions []Transaction

// Sort sorts transactions in ascending order.
func (t Transactions) Sort() {
	sort.Slice(t, func(i, j int) bool {
		return t[i].CreatedAt.Time().Before(t[j].CreatedAt.Time())
	})
}

var _ TransactionsService = &transactionsImpl{}

// transactionsImpl implements TransactionsService.
type transactionsImpl serviceImpl

func (s *transactionsImpl) List(opts *PagerOptions) Pager {
	return s.client.newPager("GET", "/transactions", opts)
}

func (s *transactionsImpl) ListAccount(accountCode string, opts *PagerOptions) Pager {
	path := fmt.Sprintf("/accounts/%s/transactions", accountCode)
	return s.client.newPager("GET", path, opts)
}

func (s *transactionsImpl) Get(ctx context.Context, uuid string) (*Transaction, error) {
	path := fmt.Sprintf("/transactions/%s", sanitizeUUID(uuid))
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var dst Transaction
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		if e, ok := err.(*ClientError); ok && e.Response.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &dst, nil
}
