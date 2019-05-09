package recurly

import (
	"encoding/json"
	"encoding/xml"
	"strconv"
	"strings"
	"time"

	"net/url"
)

// NullBool is able to marshal or unmarshal bools to XML in order to differentiate
// between false as a value and false as a zero-value.
//
// Do not initialize as a struct, use NewBool().
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

// Is checks to see if the boolean is valid and equivalent.
func (n NullBool) Is(b bool) bool { return n.Valid && n.Bool == b }

// MarshalJSON marshals an bool based on whether valid is true.
func (n NullBool) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Bool)
	}
	return []byte("null"), nil
}

// UnmarshalXML unmarshals an bool properly, as well as marshaling an empty string to nil.
func (n *NullBool) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	if err := d.DecodeElement(&v, &start); err == nil {
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

// NullInt is used for properly handling int types that could be null.
//
// Do not initialize as a struct, use NewInt().
type NullInt struct {
	Int   int
	Valid bool
}

// NewInt builds a new NullInt struct.
func NewInt(i int) NullInt {
	return NullInt{Int: i, Valid: true}
}

// MarshalJSON marshals an int based on whether valid is true.
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

// DateTimeFormat is the format Recurly uses to represent datetimes.
const DateTimeFormat = "2006-01-02T15:04:05Z07:00"

// NullTime is used for properly handling time.Time types that could be null.
//
// Do not initialize as a struct, use NewTime().
type NullTime struct {
	*time.Time
	Raw string `xml:",innerxml"`
}

// NewTime generates a new NullTime.
func NewTime(t time.Time) NullTime {
	t = t.UTC()
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

// MarshalJSON method has to be added here due to embeded interface json marshal issue in Go
// with panic on nil time field
func (t NullTime) MarshalJSON() ([]byte, error) {
	if t.Time != nil {
		return json.Marshal(t.Time)
	}
	return []byte("null"), nil
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

// NullMarshal can be embedded in structs that are read-only or just should
// never be marshaled
type NullMarshal struct{}

// MarshalXML ensures that nullMarshal doesn't marshal any xml.
func (nm NullMarshal) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return nil
}

// unmarshals href types into strings.
type href url.URL

// UnmarshalXML unmarshals an int properly, as well as marshaling an empty string to nil.
func (h *href) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v struct {
		HREF string `xml:"href,attr"`
	}
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	} else if v.HREF == "" {
		return nil
	}
	u, err := url.Parse(v.HREF)
	if err != nil {
		return err
	}
	*h = href(*u)
	return nil
}

func (h *href) LastPartOfPath() string {
	if h.Path == "" {
		return ""
	}
	h.Path = strings.TrimSuffix(h.Path, "/")
	idx := strings.LastIndex(h.Path, "/")
	if idx == -1 || len(h.Path)-1 == idx {
		return ""
	}
	return h.Path[idx+1:]
}

type hrefInt url.URL

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
	u, err := url.Parse(v.HREF)
	if err != nil {
		return err
	}
	*h = hrefInt(*u)
	return nil
}

func (h *hrefInt) LastPartOfPath() int {
	if h.Path == "" {
		return 0
	}
	h.Path = strings.TrimSuffix(h.Path, "/")
	idx := strings.LastIndex(h.Path, "/")
	if idx == -1 || len(h.Path)-1 == idx {
		return 0
	}
	i, _ := strconv.Atoi(h.Path[idx+1:])
	return i
}
