package recurly

import (
	"encoding/xml"
	"regexp"
)

type (
	// href takes href links and extracts out the code/ID into a struct
	href struct {
		nullMarshal
		HREF string
		Code string
	}
)

var rxHREF = regexp.MustCompile(`([^/]+)$`)

// UnmarshalXML unmarshals an int properly, as well as marshaling an empty string to nil.
func (h *href) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v struct {
		HREF string `xml:"href,attr"`
	}
	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}

	str := rxHREF.FindString(v.HREF)

	*h = href{
		HREF: v.HREF,
		Code: str,
	}

	return nil
}
