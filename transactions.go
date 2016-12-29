package recurly

import (
	"encoding/xml"
	"net"
)

type (
	// Transaction represents an individual transaction.
	Transaction struct {
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
		CVVResult        CVVResult // Read only
		AVSResult        AVSResult // Read only
		AVSResultStreet  string    // Read only
		AVSResultPostal  string    // Read only
		CreatedAt        NullTime  // Read only
		Account          Account
	}

	TransactionResult struct {
		NullMarshal
		Code    string `xml:"code,attr"`
		Message string `xml:",innerxml"`
	}

	// CVVResult holds transaction results for CVV fields.
	// https://www.chasepaymentech.com/card_verification_codes.html
	CVVResult struct {
		TransactionResult
	}

	// AVSResult holds transaction results for address verification.
	// http://developer.authorize.net/tools/errorgenerationguide/
	AVSResult struct {
		TransactionResult
	}
)

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
		XMLName          xml.Name   `xml:"transaction"`
		InvoiceNumber    HrefInt    `xml:"invoice"`      // use hrefInt for parsing
		SubscriptionUUID HrefString `xml:"subscription"` // use hrefString for parsing
		UUID             string     `xml:"uuid,omitempty"`
		Action           string     `xml:"action,omitempty"`
		AmountInCents    int        `xml:"amount_in_cents"`
		TaxInCents       int        `xml:"tax_in_cents,omitempty"`
		Currency         string     `xml:"currency"`
		Status           string     `xml:"status,omitempty"`
		PaymentMethod    string     `xml:"payment_method,omitempty"`
		Reference        string     `xml:"reference,omitempty"`
		Source           string     `xml:"source,omitempty"`
		Recurring        NullBool   `xml:"recurring,omitempty"`
		Test             bool       `xml:"test,omitempty"`
		Voidable         NullBool   `xml:"voidable,omitempty"`
		Refundable       NullBool   `xml:"refundable,omitempty"`
		IPAddress        net.IP     `xml:"ip_address,omitempty"`
		CVVResult        CVVResult  `xml:"cvv_result"`
		AVSResult        AVSResult  `xml:"avs_result"`
		AVSResultStreet  string     `xml:"avs_result_street,omitempty"`
		AVSResultPostal  string     `xml:"avs_result_postal,omitempty"`
		CreatedAt        NullTime   `xml:"created_at,omitempty"`
		Account          Account    `xml:"details>account"`
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
		CVVResult:        v.CVVResult,
		AVSResult:        v.AVSResult,
		AVSResultStreet:  v.AVSResultStreet,
		AVSResultPostal:  v.AVSResultPostal,
		CreatedAt:        v.CreatedAt,
		Account:          v.Account,
	}

	return nil
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

const (
	// TransactionStatusSuccess is the status for a successful transaction.
	TransactionStatusSuccess = "success"

	// TransactionStatusFailed is the status for a failed transaction.
	TransactionStatusFailed = "failed"

	// TransactionStatusVoid is the status for a voided transaction.
	TransactionStatusVoid = "void"
)
