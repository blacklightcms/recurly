package recurly

import (
	"encoding/xml"

	"github.com/blacklightcms/go-recurly/types"
)

// AddOn represents an individual add on linked to a plan.
type AddOn struct {
	XMLName                     xml.Name         `xml:"add_on"`
	Code                        string           `xml:"add_on_code,omitempty"`
	Name                        string           `xml:"name,omitempty"`
	DefaultQuantity             types.NullInt    `xml:"default_quantity,omitempty"`
	DisplayQuantityOnHostedPage types.NullBool   `xml:"display_quantity_on_hosted_page,omitempty"`
	TaxCode                     string           `xml:"tax_code,omitempty"`
	UnitAmountInCents           types.UnitAmount `xml:"unit_amount_in_cents,omitempty"`
	AccountingCode              string           `xml:"accounting_code,omitempty"`
	CreatedAt                   types.NullTime   `xml:"created_at,omitempty"`
}
