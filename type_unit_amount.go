package recurly

import "encoding/xml"

// UnitAmount is used in plans where unit amounts are represented in cents
// in both EUR and USD.
type UnitAmount struct {
	USD int `xml:"USD,omitempty" json:"USD,omitempty"`
	EUR int `xml:"EUR,omitempty" json:"EUR,omitempty"`
	GBP int `xml:"GBP,omitempty" json:"GBP,omitempty"`
	CAD int `xml:"CAD,omitempty" json:"CAD,omitempty"`
	AUD int `xml:"AUD,omitempty" json:"AUD,omitempty"`
}

type uaAlias struct {
	USD int `xml:"USD,omitempty" json:"USD,omitempty"`
	EUR int `xml:"EUR,omitempty" json:"EUR,omitempty"`
	GBP int `xml:"GBP,omitempty" json:"GBP,omitempty"`
	CAD int `xml:"CAD,omitempty" json:"CAD,omitempty"`
	AUD int `xml:"AUD,omitempty" json:"AUD,omitempty"`
}

// UnmarshalXML unmarshals an int properly, as well as marshaling an empty string to nil.
func (u *UnitAmount) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v uaAlias
	err := d.DecodeElement(&v, &start)
	if err == nil && (v.USD > 0 || v.EUR > 0 || v.CAD > 0 || v.GBP > 0 || v.AUD > 0) {
		*u = UnitAmount{USD: v.USD, EUR: v.EUR, CAD: v.CAD, GBP: v.GBP, AUD: v.AUD}
	}

	return nil
}

// MarshalXML marshals NullBools to XML. Otherwise nothing is
// marshaled.
func (u UnitAmount) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if u.USD > 0 || u.EUR > 0 || u.CAD > 0 || u.GBP > 0 || u.AUD > 0 {
		v := (uaAlias)(u)
		e.EncodeElement(v, start)
	}

	return nil
}
