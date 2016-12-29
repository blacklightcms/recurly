package types

import "encoding/xml"

// NullMarshal can be embedded in structs that are read-only or just should
// never be marshaled
type NullMarshal struct{}

// MarshalXML ensures that nullMarshal doesn't marshal any xml.
func (nm NullMarshal) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return nil
}
