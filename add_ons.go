package recurly

import "encoding/xml"

// AddOn represents an individual add on linked to a plan.
type AddOn struct {
	XMLName                     xml.Name   `xml:"add_on"`
	Code                        string     `xml:"add_on_code,omitempty"`
	Name                        string     `xml:"name,omitempty"`
	DefaultQuantity             NullInt    `xml:"default_quantity,omitempty"`
	DisplayQuantityOnHostedPage NullBool   `xml:"display_quantity_on_hosted_page,omitempty"`
	TaxCode                     string     `xml:"tax_code,omitempty"`
	UnitAmountInCents           UnitAmount `xml:"unit_amount_in_cents,omitempty"`
	AccountingCode              string     `xml:"accounting_code,omitempty"`
	CreatedAt                   NullTime   `xml:"created_at,omitempty"`
}
