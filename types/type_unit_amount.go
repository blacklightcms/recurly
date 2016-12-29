package types

import "encoding/xml"

type (
	// UnitAmount is used in plans where unit amounts are represented in cents
	// in both EUR and USD.
	UnitAmount struct {
		USD int `xml:"USD,omitempty"`
		EUR int `xml:"EUR,omitempty"`
	}

	uaAlias struct {
		USD int `xml:"USD,omitempty"`
		EUR int `xml:"EUR,omitempty"`
	}
)

// UnmarshalXML unmarshals an int properly, as well as marshaling an empty string to nil.
func (u *UnitAmount) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v uaAlias
	err := d.DecodeElement(&v, &start)
	if err == nil && (v.USD > 0 || v.EUR > 0) {
		*u = UnitAmount{USD: v.USD, EUR: v.EUR}
	}

	return nil
}

// MarshalXML marshals NullBools to XML. Otherwise nothing is
// marshaled.
func (u UnitAmount) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if u.USD > 0 || u.EUR > 0 {
		v := (uaAlias)(u)
		e.EncodeElement(v, start)
	}

	return nil
}
