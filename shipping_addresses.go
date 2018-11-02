package recurly

import (
	"encoding/xml"
)

// ShippingAddress represents a shipping address
type ShippingAddress struct {
	XMLName     xml.Name `xml:"shipping_address"`
	AccountCode string   `xml:"account,omitempty"`
	ID          int64    `xml:"id,omitempty"`
	FirstName   string   `xml:"first_name"`
	LastName    string   `xml:"last_name"`
	Nickname    string   `xml:"nickname,omitempty"`
	Address     string   `xml:"address1"`
	Address2    string   `xml:"address2"`
	Company     string   `xml:"company,omitempty"`
	City        string   `xml:"city"`
	State       string   `xml:"state"`
	Zip         string   `xml:"zip"`
	Country     string   `xml:"country"`
	Phone       string   `xml:"phone,omitempty"`
	Email       string   `xml:"email,omitempty"`
	VATNumber   string   `xml:"vat_number,omitempty"`
	CreatedAt   NullTime `xml:"created_at,omitempty"`
	UpdatedAt   NullTime `xml:"updated_at,omitempty"`
}

// UnmarshalXML unmarshals shipping addresses and handles intermediary state during unmarshaling
// for types like href.
func (s *ShippingAddress) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type shippingAddressAlias ShippingAddress
	var a struct {
		shippingAddressAlias
		XMLName xml.Name   `xml:"shipping_address"`
		Account hrefString `xml:"account,omitempty"`
	}
	if err := d.DecodeElement(&a, &start); err != nil {
		return err
	}
	*s = ShippingAddress(a.shippingAddressAlias)
	s.AccountCode = string(a.Account)
	return nil
}
