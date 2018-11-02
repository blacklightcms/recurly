// +build go1.9

package recurly

import (
	"encoding/json"
	"encoding/xml"
	"strings"
)

// NullInt is used for properly handling int types that could be null.
type NullInt struct {
	Int   int
	Valid bool
}

// NewInt builds a new NullInt struct.
func NewInt(i int) NullInt {
	return NullInt{Int: i, Valid: true}
}

// MarshalJSON marshals an int based on whether valid is true
func (n NullInt) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Int)
	}
	return []byte("null"), nil
}

// UnmarshalXML unmarshals an int properly, as well as marshaling an empty string to nil.
func (n *NullInt) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v struct {
		Int int    `xml:",chardata"`
		Nil string `xml:"nil,attr"`
	}
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	} else if strings.EqualFold(v.Nil, "nil") || strings.EqualFold(v.Nil, "true") {
		return nil
	}
	*n = NullInt{Int: v.Int, Valid: true}
	return nil
}

// MarshalXML marshals NullInts greater than zero to XML. Otherwise nothing is
// marshaled.
func (n NullInt) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if n.Valid {
		e.EncodeElement(n.Int, start)
	}
	return nil
}
