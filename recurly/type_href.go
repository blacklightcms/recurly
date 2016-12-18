package recurly

import (
	"encoding/xml"
	"regexp"
	"strconv"
)

var rxHREF = regexp.MustCompile(`([^/]+)$`)

type hrefString string

// UnmarshalXML unmarshals an int properly, as well as marshaling an empty string to nil.
func (h *hrefString) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v struct {
		HREF string `xml:"href,attr"`
	}
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	} else if v.HREF == "" {
		return nil
	}

	*h = hrefString(rxHREF.FindString(v.HREF))
	return nil
}

type hrefInt int

// UnmarshalXML unmarshals an int properly, as well as marshaling an empty string to nil.
func (h *hrefInt) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v struct {
		HREF string `xml:"href,attr"`
	}
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	} else if v.HREF == "" {
		return nil
	}

	i, err := strconv.Atoi(rxHREF.FindString(v.HREF))
	if err != nil {
		return err
	}

	*h = hrefInt(i)
	return nil
}
