package recurly

import (
	"encoding/xml"
	"time"
)

type (
	// NullTime is used for properly handling time.Time types that could be null.
	NullTime struct {
		*time.Time
		Raw string `xml:",innerxml"`
	}
)

// NewTime generates a new NullTime.
func NewTime(t time.Time) NullTime {
	return NullTime{Time: &t}
}

// UnmarshalXML unmarshals an int properly, as well as marshaling an empty string to nil.
func (t *NullTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	err := d.DecodeElement(&v, &start)
	if err == nil && v != "" {
		parsed, err := time.Parse("2006-01-02T15:04:05Z07:00", v)
		if err != nil {
			return err
		}

		*t = NullTime{Time: &parsed}
	}

	return nil
}

// MarshalXML marshals times into their proper format. Otherwise nothing is
// marshaled. All times are sent in UTC.
func (t NullTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if t.Time != nil {
		t.Time.UTC()
		e.EncodeElement(t.Time.Format("2006-01-02T15:04:05Z07:00"), start)
	}

	return nil
}

func (t NullTime) String() string {
	if t.Time != nil {
		return t.Time.String()
	}

	return ""
}
