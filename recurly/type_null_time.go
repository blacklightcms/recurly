package recurly

import (
	"encoding/xml"
	"time"
)

const (
	datetimeFormat = "2006-01-02T15:04:05Z07:00"
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
	t = t.UTC()
	return NullTime{Time: &t}
}

// newTimeFromString generates a new NullTime based on a
// time string in the datetimeFormat format.
// This is primarily used in unit testing.
func newTimeFromString(str string) NullTime {
	t, _ := time.Parse(datetimeFormat, str)
	return NullTime{Time: &t}
}

// UnmarshalXML unmarshals an int properly, as well as marshaling an empty string to nil.
func (t *NullTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	err := d.DecodeElement(&v, &start)
	if err == nil && v != "" {
		parsed, err := time.Parse(datetimeFormat, v)
		if err != nil {
			return err
		}

		*t = NewTime(parsed)
	}

	return nil
}

// MarshalXML marshals times into their proper format. Otherwise nothing is
// marshaled. All times are sent in UTC.
func (t NullTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if t.Time != nil {
		e.EncodeElement(t.String(), start)
	}

	return nil
}

// String returns a string representation of the time in UTC using the
// datetimeFormat constant as the format.
func (t NullTime) String() string {
	if t.Time != nil {
		return t.Time.UTC().Format(datetimeFormat)
	}

	return ""
}
