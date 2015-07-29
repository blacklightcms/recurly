package recurly

import (
	"encoding/xml"
	"fmt"
)

type (
	// AdjustmentsService handles communication with the adjustments related methods
	// of the recurly API.
	AdjustmentsService struct {
		client *Client
	}

	// Adjustment works with charges and credits on a given account.
	Adjustment struct {
		XMLName                xml.Name      `xml:"adjustment"`
		Account                href          `xml:"account,omitempty"`
		Invoice                href          `xml:"invoice,omitempty"`
		UUID                   string        `xml:"uuid,omitempty"`
		State                  string        `xml:"state,omitempty"`
		Description            string        `xml:"description,omitempty"`
		AccountingCode         string        `xml:"accounting_code,omitempty"`
		ProductCode            string        `xml:"product_code,omitempty"`
		Origin                 string        `xml:"origin,omitempty"`
		UnitAmountInCents      int           `xml:"unit_amount_in_cents"`
		Quantity               int           `xml:"quantity,omitempty"`
		OriginalAdjustmentUUID string        `xml:"original_adjustment_uuid,omitempty"`
		DiscountInCents        int           `xml:"discount_in_cents,omitempty"`
		TaxInCents             int           `xml:"tax_in_cents,omitempty"`
		TotalInCents           int           `xml:"total_in_cents,omitempty"`
		Currency               string        `xml:"currency"`
		Taxable                NullBool      `xml:"taxable,omitempty"`
		TaxCode                string        `xml:"tax_code,omitempty"`
		TaxType                string        `xml:"tax_type,omitempty"`
		TaxRegion              string        `xml:"tax_region,omitempty"`
		TaxRate                float64       `xml:"tax_rate,omitempty"`
		TaxExempt              NullBool      `xml:"tax_exempt,omitempty"`
		TaxDetails             *[]TaxDetails `xml:"tax_details>tax_detail,omitempty"`
		StartDate              NullTime      `xml:"start_date,omitempty"`
		EndDate                NullTime      `xml:"end_date,omitempty"`
		CreatedAt              NullTime      `xml:"created_at,omitempty"`
	}

	// TaxDetails holds tax information and is embedded in an Adjustment.
	TaxDetails struct {
		XMLName    xml.Name `xml:"tax_detail"`
		Name       string   `xml:"name,omitempty"`
		Type       string   `xml:"type,omitempty"`
		TaxRate    float64  `xml:"tax_rate,omitempty"`
		TaxInCents int      `xml:"tax_in_cents,omitempty"`
	}
)

// List retrieves all charges and credits issued for an account
// https://docs.recurly.com/api/adjustments#list-adjustments
func (service AdjustmentsService) List(accountCode string, params Params) (*Response, []Adjustment, error) {
	action := fmt.Sprintf("accounts/%s/adjustments", accountCode)
	req, err := service.client.newRequest("GET", action, params, nil)
	if err != nil {
		return nil, nil, err
	}

	var a struct {
		XMLName     xml.Name     `xml:"adjustments"`
		Adjustments []Adjustment `xml:"adjustment"`
	}
	res, err := service.client.do(req, &a)

	return res, a.Adjustments, err
}

// Get returns information about a single adjustment.
// https://docs.recurly.com/api/adjustments#get-adjustments
func (service AdjustmentsService) Get(uuid string) (*Response, Adjustment, error) {
	action := fmt.Sprintf("adjustments/%s", uuid)
	req, err := service.client.newRequest("GET", action, nil, nil)
	if err != nil {
		return nil, Adjustment{}, err
	}

	var a Adjustment
	res, err := service.client.do(req, &a)

	return res, a, err
}

// Create creates a one-time charge on an account. Charges are not invoiced or
// collected immediately. Non-invoiced charges will automatically be invoices
// when the account's subscription renews, or you trigger a collection by
// posting an invoice. Charges may be removed from an account if they have
// not been invoiced.
// https://docs.recurly.com/api/adjustments#create-adjustment
func (service AdjustmentsService) Create(accountCode string, a Adjustment) (*Response, Adjustment, error) {
	action := fmt.Sprintf("accounts/%s/adjustments", accountCode)
	req, err := service.client.newRequest("POST", action, nil, a)
	if err != nil {
		return nil, Adjustment{}, err
	}

	var dest Adjustment
	res, err := service.client.do(req, &dest)

	return res, a, err
}

// Delete removes a non-invoiced adjustment from an account.
// https://docs.recurly.com/api/adjustments#delete-adjustment
func (service AdjustmentsService) Delete(uuid string) (*Response, error) {
	action := fmt.Sprintf("adjustments/%s", uuid)
	req, err := service.client.newRequest("DELETE", action, nil, nil)
	if err != nil {
		return nil, err
	}

	return service.client.do(req, nil)
}
