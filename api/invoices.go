package api

import (
	"bytes"
	"encoding/xml"
	"fmt"

	recurly "github.com/blacklightcms/go-recurly"
)

var _ recurly.InvoicesService = &InvoicesService{}

// InvoicesService handles communication with theinvoices related methods
// of the recurly API.
type InvoicesService struct {
	client *Client
}

// List returns a list of all invoices.
// https://dev.recurly.com/docs/list-invoices
func (s *InvoicesService) List(params recurly.Params) (*recurly.Response, []recurly.Invoice, error) {
	req, err := s.client.newRequest("GET", "invoices", params, nil)
	if err != nil {
		return nil, nil, err
	}

	var p struct {
		XMLName  xml.Name          `xml:"invoices"`
		Invoices []recurly.Invoice `xml:"invoice"`
	}
	resp, err := s.client.do(req, &p)

	return resp, p.Invoices, err
}

// ListAccount returns a list of all invoices for an account.
// https://dev.recurly.com/docs/list-an-accounts-invoices
func (s *InvoicesService) ListAccount(accountCode string, params recurly.Params) (*recurly.Response, []recurly.Invoice, error) {
	action := fmt.Sprintf("accounts/%s/invoices", accountCode)
	req, err := s.client.newRequest("GET", action, params, nil)
	if err != nil {
		return nil, nil, err
	}

	var p struct {
		XMLName  xml.Name          `xml:"invoices"`
		Invoices []recurly.Invoice `xml:"invoice"`
	}
	resp, err := s.client.do(req, &p)

	return resp, p.Invoices, err
}

// Get returns detailed information about an invoice including line items and
// payments.
// https://dev.recurly.com/docs/lookup-invoice-details
func (s *InvoicesService) Get(invoiceNumber int) (*recurly.Response, *recurly.Invoice, error) {
	action := fmt.Sprintf("invoices/%d", invoiceNumber)
	req, err := s.client.newRequest("GET", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst recurly.Invoice
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// GetPDF retrieves the invoice as a PDF.
// The language parameters allows you to specify a language to translate the
// invoice into. If empty, English will be used. Options: Danish, German,
// Spanish, French, Hindi, Japanese, Dutch, Portuguese, Russian, Turkish, Chinese.
// https://dev.recurly.com/docs/retrieve-a-pdf-invoice
func (s *InvoicesService) GetPDF(invoiceNumber int, language string) (*recurly.Response, *bytes.Buffer, error) {
	action := fmt.Sprintf("invoices/%d", invoiceNumber)
	req, err := s.client.newRequest("GET", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	if language == "" {
		language = "English"
	}

	req.Header.Set("Accept", "application/pdf")
	req.Header.Set("Accept-Language", language)

	var pdf bytes.Buffer
	resp, err := s.client.do(req, &pdf)

	return resp, &pdf, err
}

// Preview allows you to display the invoice details, including estimated tax,
// before you post it.
// https://dev.recurly.com/docs/post-an-invoice-invoice-pending-charges-on-an-acco
func (s *InvoicesService) Preview(accountCode string) (*recurly.Response, *recurly.Invoice, error) {
	action := fmt.Sprintf("accounts/%s/invoices/preview", accountCode)
	req, err := s.client.newRequest("POST", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst recurly.Invoice
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// Create posts an accounts pending charges to a new invoice on that account.
// When you post one-time charges to an account, these will remain pending
// until they are invoiced. An account is automatically invoiced when the
// subscription renews. However, there are times when it is appropriate to
// invoice an account before the renewal. If the subscriber has a yearly
// subscription, you might want to collect the one-time charges well before the renewal.
// https://dev.recurly.com/docs/post-an-invoice-invoice-pending-charges-on-an-acco
func (s *InvoicesService) Create(accountCode string, invoice recurly.Invoice) (*recurly.Response, *recurly.Invoice, error) {
	action := fmt.Sprintf("accounts/%s/invoices", accountCode)
	req, err := s.client.newRequest("POST", action, nil, invoice)
	if err != nil {
		return nil, nil, err
	}

	var dst recurly.Invoice
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// MarkPaid marks an invoice as paid successfully.
// https://dev.recurly.com/docs/mark-an-invoice-as-paid-successfully
func (s *InvoicesService) MarkPaid(invoiceNumber int) (*recurly.Response, *recurly.Invoice, error) {
	action := fmt.Sprintf("invoices/%d/mark_successful", invoiceNumber)
	req, err := s.client.newRequest("PUT", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst recurly.Invoice
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// MarkFailed marks an invoice as failed.
// https://dev.recurly.com/docs/mark-an-invoice-as-failed-collection
func (s *InvoicesService) MarkFailed(invoiceNumber int) (*recurly.Response, *recurly.Invoice, error) {
	action := fmt.Sprintf("invoices/%d/mark_failed", invoiceNumber)
	req, err := s.client.newRequest("PUT", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst recurly.Invoice
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}
