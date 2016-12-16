package recurly

import "encoding/xml"

// NullMarshal can be embedded in structs that are read-only or just should
// never be marshaled
type nullMarshal struct{}

// MarshalXML ensures that nullMarshal doesn't marshal any xml.
func (nm nullMarshal) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return nil
}
