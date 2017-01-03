package recurly

import "encoding/xml"

type (
	// NullInt is used for properly handling int types that could be null.
	NullInt struct {
		Int   int
		Valid bool
	}
)

// NewInt builds a new NullInt struct.
func NewInt(i int) NullInt {
	return NullInt{Int: i, Valid: true}
}

// UnmarshalXML unmarshals an int properly, as well as marshaling an empty string to nil.
func (n *NullInt) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v int
	err := d.DecodeElement(&v, &start)
	if err == nil {
		*n = NullInt{Int: v, Valid: true}
	}

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
