package recurly

import (
	"encoding/xml"

	"github.com/blacklightcms/go-recurly/types"
)

type (
	// Adjustment works with charges and credits on a given account.
	Adjustment struct {
		AccountCode            string
		InvoiceNumber          int
		UUID                   string
		State                  string
		Description            string
		AccountingCode         string
		ProductCode            string
		Origin                 string
		UnitAmountInCents      int
		Quantity               int
		OriginalAdjustmentUUID string
		DiscountInCents        int
		TaxInCents             int
		TotalInCents           int
		Currency               string
		Taxable                types.NullBool
		TaxCode                string
		TaxType                string
		TaxRegion              string
		TaxRate                float64
		TaxExempt              types.NullBool
		TaxDetails             []TaxDetail
		StartDate              types.NullTime
		EndDate                types.NullTime
		CreatedAt              types.NullTime
	}

	// TaxDetail holds tax information and is embedded in an Adjustment.
	// TaxDetails are a read only field, so theys houldn't marshall
	TaxDetail struct {
		XMLName    xml.Name `xml:"tax_detail"`
		Name       string   `xml:"name,omitempty"`
		Type       string   `xml:"type,omitempty"`
		TaxRate    float64  `xml:"tax_rate,omitempty"`
		TaxInCents int      `xml:"tax_in_cents,omitempty"`
	}
)

// MarshalXML marshals only the fields needed for creating/updating adjustments
// with the recurly API.
func (a Adjustment) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	v := struct {
		XMLName           xml.Name       `xml:"adjustment"`
		Description       string         `xml:"description,omitempty"`
		AccountingCode    string         `xml:"accounting_code,omitempty"`
		UnitAmountInCents int            `xml:"unit_amount_in_cents"`
		Quantity          int            `xml:"quantity,omitempty"`
		Currency          string         `xml:"currency"`
		TaxCode           string         `xml:"tax_code,omitempty"`
		TaxExempt         types.NullBool `xml:"tax_exempt,omitempty"`
	}{
		Description:       a.Description,
		AccountingCode:    a.AccountingCode,
		UnitAmountInCents: a.UnitAmountInCents,
		Quantity:          a.Quantity,
		Currency:          a.Currency,
		TaxCode:           a.TaxCode,
		TaxExempt:         a.TaxExempt,
	}

	return e.Encode(v)
}

// UnmarshalXML unmarshal a coupon redemption object. Minaly converts href links
// for coupons and accounts to CouponCode and AccountCodes.
func (a *Adjustment) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v struct {
		XMLName                xml.Name         `xml:"adjustment"`
		AccountCode            types.HrefString `xml:"account,omitempty"` // Read only
		InvoiceNumber          types.HrefInt    `xml:"invoice,omitempty"` // Read only
		UUID                   string           `xml:"uuid,omitempty"`
		State                  string           `xml:"state,omitempty"`
		Description            string           `xml:"description,omitempty"`
		AccountingCode         string           `xml:"accounting_code,omitempty"`
		ProductCode            string           `xml:"product_code,omitempty"`
		Origin                 string           `xml:"origin,omitempty"`
		UnitAmountInCents      int              `xml:"unit_amount_in_cents"`
		Quantity               int              `xml:"quantity,omitempty"`
		OriginalAdjustmentUUID string           `xml:"original_adjustment_uuid,omitempty"`
		DiscountInCents        int              `xml:"discount_in_cents,omitempty"`
		TaxInCents             int              `xml:"tax_in_cents,omitempty"`
		TotalInCents           int              `xml:"total_in_cents,omitempty"`
		Currency               string           `xml:"currency"`
		Taxable                types.NullBool   `xml:"taxable,omitempty"`
		TaxCode                string           `xml:"tax_code,omitempty"`
		TaxType                string           `xml:"tax_type,omitempty"`
		TaxRegion              string           `xml:"tax_region,omitempty"`
		TaxRate                float64          `xml:"tax_rate,omitempty"`
		TaxExempt              types.NullBool   `xml:"tax_exempt,omitempty"`
		TaxDetails             []TaxDetail      `xml:"tax_details>tax_detail,omitempty"`
		StartDate              types.NullTime   `xml:"start_date,omitempty"`
		EndDate                types.NullTime   `xml:"end_date,omitempty"`
		CreatedAt              types.NullTime   `xml:"created_at,omitempty"`
	}
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	*a = Adjustment{
		AccountCode:            string(v.AccountCode),
		InvoiceNumber:          int(v.InvoiceNumber),
		UUID:                   v.UUID,
		State:                  v.State,
		Description:            v.Description,
		AccountingCode:         v.AccountingCode,
		ProductCode:            v.ProductCode,
		Origin:                 v.Origin,
		UnitAmountInCents:      v.UnitAmountInCents,
		Quantity:               v.Quantity,
		OriginalAdjustmentUUID: v.OriginalAdjustmentUUID,
		DiscountInCents:        v.DiscountInCents,
		TaxInCents:             v.TaxInCents,
		TotalInCents:           v.TotalInCents,
		Currency:               v.Currency,
		Taxable:                v.Taxable,
		TaxCode:                v.TaxCode,
		TaxType:                v.TaxType,
		TaxRegion:              v.TaxRegion,
		TaxRate:                v.TaxRate,
		TaxExempt:              v.TaxExempt,
		TaxDetails:             v.TaxDetails,
		StartDate:              v.StartDate,
		EndDate:                v.EndDate,
		CreatedAt:              v.CreatedAt,
	}

	return nil
}
