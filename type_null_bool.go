package recurly

import (
	"encoding/xml"
	"strconv"
)

// NullBool is used for properly handling bool types with the Recurly API.
// Without it, setting a false boolean value will be ignored when encoding
// xml requests to the Recurly API.
type NullBool struct {
	Bool  bool
	Valid bool
}

// NewBool creates a new NullBool.
func NewBool(b bool) NullBool {
	return NullBool{
		Bool:  b,
		Valid: true,
	}
}

// Is checks to see if the boolean is valid and equivalent
func (n NullBool) Is(b bool) bool {
	return n.Valid && n.Bool == b
}

// UnmarshalXML unmarshals an bool properly, as well as marshaling an empty string to nil.
func (n *NullBool) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	err := d.DecodeElement(&v, &start)
	if err == nil {
		val, _ := strconv.ParseBool(v)
		*n = NullBool{Bool: val, Valid: true}
	}

	return nil
}

// MarshalXML marshals NullBools to XML. Otherwise nothing is
// marshaled.
func (n NullBool) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if n.Valid {
		e.EncodeElement(n.Bool, start)
	}

	return nil
}
