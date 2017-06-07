package recurly

import (
	"encoding/xml"
	"net"
)

const (
	// TransactionStatusSuccess is the status for a successful transaction.
	TransactionStatusSuccess = "success"

	// TransactionStatusFailed is the status for a failed transaction.
	TransactionStatusFailed = "failed"

	// TransactionStatusVoid is the status for a voided transaction.
	TransactionStatusVoid = "void"
)

// Transaction represents an individual transaction.
type Transaction struct {
	InvoiceNumber    int    // Read only
	SubscriptionUUID string // Read only
	UUID             string // Read only
	Action           string
	AmountInCents    int
	TaxInCents       int
	Currency         string
	Status           string
	PaymentMethod    string
	Reference        string
	Source           string
	Recurring        NullBool
	Test             bool
	Voidable         NullBool
	Refundable       NullBool
	IPAddress        net.IP
	TransactionError TransactionError // Read only
	CVVResult        CVVResult        // Read only
	AVSResult        AVSResult        // Read only
	AVSResultStreet  string           // Read only
	AVSResultPostal  string           // Read only
	CreatedAt        NullTime         // Read only
	Account          Account
}

// TransactionError is an error encounted from your payment gateway that
// recurly has standardized.
// https://recurly.readme.io/v2.0/page/transaction-errors
type TransactionError struct {
	XMLName          xml.Name `xml:"transaction_error"`
	ErrorCode        string   `xml:"error_code,omitempty"`
	ErrorCategory    string   `xml:"error_category,omitempty"`
	MerchantMessage  string   `xml:"merchant_message,omitempty"`
	CustomerMessage  string   `xml:"customer_message,omitempty"`
	GatewayErrorCode string   `xml:"gateway_error_code,omitempty"`
}

// MarshalXML marshals a transaction sending only the fields recurly allows for writes.
// Read only fields are not encoded, and account is written as <account></account>
// instead of as <details><account></account></details> (like it is in Transaction).
func (t Transaction) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	dst := struct {
		XMLName       xml.Name `xml:"transaction"`
		Action        string   `xml:"action,omitempty"`
		AmountInCents int      `xml:"amount_in_cents"`
		TaxInCents    int      `xml:"tax_in_cents,omitempty"`
		Currency      string   `xml:"currency"`
		Status        string   `xml:"status,omitempty"`
		PaymentMethod string   `xml:"payment_method,omitempty"`
		Reference     string   `xml:"reference,omitempty"`
		Source        string   `xml:"source,omitempty"`
		Recurring     NullBool `xml:"recurring,omitempty"`
		Test          bool     `xml:"test,omitempty"`
		Voidable      NullBool `xml:"voidable,omitempty"`
		Refundable    NullBool `xml:"refundable,omitempty"`
		IPAddress     net.IP   `xml:"ip_address,omitempty"`
		Account       Account  `xml:"account"`
	}{
		Action:        t.Action,
		AmountInCents: t.AmountInCents,
		TaxInCents:    t.TaxInCents,
		Currency:      t.Currency,
		Status:        t.Status,
		PaymentMethod: t.PaymentMethod,
		Reference:     t.Reference,
		Source:        t.Source,
		Recurring:     t.Recurring,
		Test:          t.Test,
		Voidable:      t.Voidable,
		Refundable:    t.Refundable,
		IPAddress:     t.IPAddress,
		Account:       t.Account,
	}
	e.Encode(dst)
	return nil
}

// UnmarshalXML unmarshals transactions and handles intermediary state during unmarshaling
// for types like href.
func (t *Transaction) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v struct {
		XMLName          xml.Name         `xml:"transaction"`
		InvoiceNumber    hrefInt          `xml:"invoice"`      // use hrefInt for parsing
		SubscriptionUUID hrefString       `xml:"subscription"` // use hrefString for parsing
		UUID             string           `xml:"uuid,omitempty"`
		Action           string           `xml:"action,omitempty"`
		AmountInCents    int              `xml:"amount_in_cents"`
		TaxInCents       int              `xml:"tax_in_cents,omitempty"`
		Currency         string           `xml:"currency"`
		Status           string           `xml:"status,omitempty"`
		PaymentMethod    string           `xml:"payment_method,omitempty"`
		Reference        string           `xml:"reference,omitempty"`
		Source           string           `xml:"source,omitempty"`
		Recurring        NullBool         `xml:"recurring,omitempty"`
		Test             bool             `xml:"test,omitempty"`
		Voidable         NullBool         `xml:"voidable,omitempty"`
		Refundable       NullBool         `xml:"refundable,omitempty"`
		IPAddress        net.IP           `xml:"ip_address,omitempty"`
		TransactionError TransactionError `xml:"transaction_error,omitempty"`
		CVVResult        CVVResult        `xml:"cvv_result"`
		AVSResult        AVSResult        `xml:"avs_result"`
		AVSResultStreet  string           `xml:"avs_result_street,omitempty"`
		AVSResultPostal  string           `xml:"avs_result_postal,omitempty"`
		CreatedAt        NullTime         `xml:"created_at,omitempty"`
		Account          Account          `xml:"details>account"`
	}
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	*t = Transaction{
		InvoiceNumber:    int(v.InvoiceNumber),
		SubscriptionUUID: string(v.SubscriptionUUID),
		UUID:             v.UUID,
		Action:           v.Action,
		AmountInCents:    v.AmountInCents,
		TaxInCents:       v.TaxInCents,
		Currency:         v.Currency,
		Status:           v.Status,
		PaymentMethod:    v.PaymentMethod,
		Reference:        v.Reference,
		Source:           v.Source,
		Recurring:        v.Recurring,
		Test:             v.Test,
		Voidable:         v.Voidable,
		Refundable:       v.Refundable,
		IPAddress:        v.IPAddress,
		TransactionError: v.TransactionError,
		CVVResult:        v.CVVResult,
		AVSResult:        v.AVSResult,
		AVSResultStreet:  v.AVSResultStreet,
		AVSResultPostal:  v.AVSResultPostal,
		CreatedAt:        v.CreatedAt,
		Account:          v.Account,
	}

	return nil
}

type TransactionResult struct {
	NullMarshal
	Code    string `xml:"code,attr"`
	Message string `xml:",innerxml"`
}

// CVVResult holds transaction results for CVV fields.
// https://www.chasepaymentech.com/card_verification_codes.html
type CVVResult struct {
	TransactionResult
}

// IsMatch returns true if the CVV code is a match.
func (c CVVResult) IsMatch() bool {
	return c.Code == "M" || c.Code == "Y"
}

// IsNoMatch returns true if the CVV code did not match.
func (c CVVResult) IsNoMatch() bool {
	return c.Code == "N"
}

// NotProcessed returns true if the CVV code was not processed.
func (c CVVResult) NotProcessed() bool {
	return c.Code == "P"
}

// ShouldHaveBeenPresent returns true if the CVV code should have been present
// on the card but was not indicated.
func (c CVVResult) ShouldHaveBeenPresent() bool {
	return c.Code == "S"
}

// UnableToProcess returns true when the issuer was unable to process the CVV.
func (c CVVResult) UnableToProcess() bool {
	return c.Code == "U"
}

// AVSResult holds transaction results for address verification.
// http://developer.authorize.net/tools/errorgenerationguide/
type AVSResult struct {
	TransactionResult
}

// Transactions is a sortable slice of Transaction.
// It implements sort.Interface.
type Transactions []Transaction

func (s Transactions) Len() int {
	return len(s)
}

func (s Transactions) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Transactions) Less(i, j int) bool {
	return s[i].CreatedAt.Time.Before(*s[j].CreatedAt.Time)
}
