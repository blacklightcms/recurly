package recurly

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
)

// AdjustmentsService manages the interactions for adjustments.
type AdjustmentsService interface {
	// List returns a pager to paginate adjustments for an account. PagerOptions are
	// used to optionally filter the results.
	//
	// https://dev.recurly.com/docs/list-an-accounts-adjustments
	ListAccount(accountCode string, opts *PagerOptions) Pager

	// Get retrieves an adjustment. If the add on does not exist,
	// a nil adjustment and nil error are returned.
	//
	// NOTE: Link below is from v2.9, the last documented version showing
	// this endpoint. The endpoint appears to still be valid in v2.19, waiting
	// to hear back from Recurly support.
	// https://dev.recurly.com/v2.9/docs/get-an-adjustment
	Get(ctx context.Context, uuid string) (*Adjustment, error)

	// Create creates a one-time charge on an account. Charges are not invoiced
	// or collected immediately. Non-invoiced charges will automatically be
	// invoiced when the account's subscription renews, or you trigger a
	// collection by posting an invoice. Charges may be removed from an account
	// if they have not been invoiced.
	//
	// For a charge, set a.UnitAmountInCents to a positive number.
	// For a credit, set a.UnitAmountInCents to a negative amount.
	//
	// https://dev.recurly.com/docs/create-a-charge
	// https://dev.recurly.com/docs/create-a-credit
	Create(ctx context.Context, accountCode string, a Adjustment) (*Adjustment, error)

	// Delete deletes an adjustment from an account. Only non-invoiced adjustments
	// can be deleted.
	//
	// https://dev.recurly.com/docs/delete-an-adjustment
	Delete(ctx context.Context, uuid string) error
}

// Adjustment state constants.
const (
	AdjustmentStatePending = "pending"
	AdjustmentStateInvoied = "invoiced"
)

// Revenue schedule type constants.
const (
	RevenueScheduleTypeNever        = "never"
	RevenueScheduleTypeAtRangeStart = "at_range_start"
	RevenueScheduleTypeAtInvoice    = "at_invoice"
	RevenueScheduleTypeEvenly       = "evenly"       // if 'end_date' is set
	RevenueScheduleTypeAtRangeEnd   = "at_range_end" // if 'end_date' is set
)

// Adjustment works with charges and credits on a given account.
//
// https://dev.recurly.com/docs/adjustment-object
type Adjustment struct {
	XMLName                xml.Name    `xml:"adjustment"`
	AccountCode            string      `xml:"-"` // Read only
	InvoiceNumber          int         `xml:"-"` // Read only
	SubscriptionUUID       string      `xml:"-"` // Read only
	UUID                   string      `xml:"uuid,omitempty"`
	State                  string      `xml:"state,omitempty"`
	Description            string      `xml:"description,omitempty"`
	AccountingCode         string      `xml:"accounting_code,omitempty"`
	RevenueScheduleType    string      `xml:"revenue_schedule_type,omitempty"`
	ProductCode            string      `xml:"product_code,omitempty"`
	Origin                 string      `xml:"origin,omitempty"`
	UnitAmountInCents      NullInt     `xml:"unit_amount_in_cents,omitempty"`
	Quantity               int         `xml:"quantity,omitempty"`
	OriginalAdjustmentUUID string      `xml:"original_adjustment_uuid,omitempty"`
	DiscountInCents        int         `xml:"discount_in_cents,omitempty"`
	TaxInCents             int         `xml:"tax_in_cents,omitempty"`
	TotalInCents           int         `xml:"total_in_cents,omitempty"`
	Currency               string      `xml:"currency"`
	TaxCode                string      `xml:"tax_code,omitempty"`
	TaxType                string      `xml:"tax_type,omitempty"`
	TaxRegion              string      `xml:"tax_region,omitempty"`
	TaxRate                float64     `xml:"tax_rate,omitempty"`
	TaxExempt              NullBool    `xml:"tax_exempt,omitempty"`
	TaxDetails             []TaxDetail `xml:"tax_details>tax_detail,omitempty"`

	// The following are only valid with an `Avalara for Communications` integration
	AvalaraTransactionType int `xml:"avalara_transaction_type,omitempty"`
	AvalaraServiceType     int `xml:"avalara_service_type,omitempty"`

	StartDate NullTime `xml:"start_date,omitempty"`
	EndDate   NullTime `xml:"end_date,omitempty"`
	CreatedAt NullTime `xml:"created_at,omitempty"`
	UpdatedAt NullTime `xml:"updated_at,omitempty"`
}

// MarshalXML marshals only the fields needed for creating/updating adjustments
// with the recurly API.
func (a Adjustment) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.Encode(struct {
		XMLName                xml.Name `xml:"adjustment"`
		Description            string   `xml:"description,omitempty"`
		AccountingCode         string   `xml:"accounting_code,omitempty"`
		RevenueScheduleType    string   `xml:"revenue_schedule_type,omitempty"`
		ProductCode            string   `xml:"product_code,omitempty"`
		Origin                 string   `xml:"origin,omitempty"`
		UnitAmountInCents      NullInt  `xml:"unit_amount_in_cents,omitempty"`
		Quantity               int      `xml:"quantity,omitempty"`
		Currency               string   `xml:"currency,omitempty"`
		TaxCode                string   `xml:"tax_code,omitempty"`
		TaxExempt              NullBool `xml:"tax_exempt,omitempty"`
		AvalaraTransactionType int      `xml:"avalara_transaction_type,omitempty"`
		AvalaraServiceType     int      `xml:"avalara_service_type,omitempty"`
		StartDate              NullTime `xml:"start_date,omitempty"`
		EndDate                NullTime `xml:"end_date,omitempty"`
	}{
		Description:            a.Description,
		AccountingCode:         a.AccountingCode,
		RevenueScheduleType:    a.RevenueScheduleType,
		ProductCode:            a.ProductCode,
		Origin:                 a.Origin,
		UnitAmountInCents:      a.UnitAmountInCents,
		Quantity:               a.Quantity,
		Currency:               a.Currency,
		TaxCode:                a.TaxCode,
		TaxExempt:              a.TaxExempt,
		AvalaraServiceType:     a.AvalaraServiceType,
		AvalaraTransactionType: a.AvalaraTransactionType,
		StartDate:              a.StartDate,
		EndDate:                a.EndDate,
	})
}

// UnmarshalXML unmarshal an adjustment.
func (a *Adjustment) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type adjustmentAlias Adjustment
	var v struct {
		XMLName xml.Name `xml:"adjustment"`
		adjustmentAlias
		AccountCode      href    `xml:"account,omitempty"`
		InvoiceNumber    hrefInt `xml:"invoice,omitempty"`
		SubscriptionUUID href    `xml:"subscription,omitempty"`
	}
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}

	*a = Adjustment(v.adjustmentAlias)
	a.XMLName = v.XMLName
	a.AccountCode = v.AccountCode.LastPartOfPath()
	a.InvoiceNumber = v.InvoiceNumber.LastPartOfPath()
	a.SubscriptionUUID = v.SubscriptionUUID.LastPartOfPath()
	return nil
}

// TaxDetail holds tax information and is embedded in an Adjustment.
// TaxDetails are a read only field, so they shouldn't marshal.
type TaxDetail struct {
	XMLName    xml.Name `xml:"tax_detail"`
	Name       string   `xml:"name,omitempty"`
	Type       string   `xml:"type,omitempty"`
	TaxRate    float64  `xml:"tax_rate,omitempty"`
	TaxInCents int      `xml:"tax_in_cents,omitempty"`
	Level      string   `xml:"level,omitempty"`
	Billable   NullBool `xml:"billable,omitempty"`
}

var _ AdjustmentsService = &adjustmentsImpl{}

// adjustmentsImpl implements AdjustmentsService.
type adjustmentsImpl serviceImpl

func (s *adjustmentsImpl) ListAccount(accountCode string, opts *PagerOptions) Pager {
	path := fmt.Sprintf("/accounts/%s/adjustments", accountCode)
	return s.client.newPager("GET", path, opts)
}

func (s *adjustmentsImpl) Get(ctx context.Context, uuid string) (*Adjustment, error) {
	path := fmt.Sprintf("/adjustments/%s", sanitizeUUID(uuid))
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var dst Adjustment
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		if e, ok := err.(*ClientError); ok && e.Response.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &dst, nil
}

func (s *adjustmentsImpl) Create(ctx context.Context, accountCode string, a Adjustment) (*Adjustment, error) {
	path := fmt.Sprintf("/accounts/%s/adjustments", accountCode)
	req, err := s.client.newRequest("POST", path, a)
	if err != nil {
		return nil, err
	}

	var dst Adjustment
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *adjustmentsImpl) Delete(ctx context.Context, uuid string) error {
	path := fmt.Sprintf("/adjustments/%s", sanitizeUUID(uuid))
	req, err := s.client.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(ctx, req, nil)
	return err
}
