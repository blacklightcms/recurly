package recurly

import (
	"encoding/json"
	"encoding/xml"
	"time"
)

// DateTimeFormat is the format Recurly uses to represent datetimes.
const DateTimeFormat = "2006-01-02T15:04:05Z07:00"

// NullTime is used for properly handling time.Time types that could be null.
type NullTime struct {
	*time.Time `json:"time,omitempty"`
	Raw        string `xml:",innerxml"`
}

// NewTime generates a new NullTime.
func NewTime(t time.Time) NullTime {
	t = t.UTC()
	return NullTime{Time: &t}
}

// NewTimeFromString generates a new NullTime based on a
// time string in the DateTimeFormat format.
// This is primarily used in unit testing.
func NewTimeFromString(str string) NullTime {
	t, _ := time.Parse(DateTimeFormat, str)
	return NullTime{Time: &t}
}

// UnmarshalXML unmarshals an int properly, as well as marshaling an empty string to nil.
func (t *NullTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	err := d.DecodeElement(&v, &start)
	if err == nil && v != "" {
		parsed, err := time.Parse(DateTimeFormat, v)
		if err != nil {
			return err
		}

		*t = NewTime(parsed)
	}

	return nil
}

func (t NullTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time)
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
// DateTimeFormat constant as the format.
func (t NullTime) String() string {
	if t.Time != nil {
		return t.Time.UTC().Format(DateTimeFormat)
	}

	return ""
}
