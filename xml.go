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
type NullBool struct {
	value bool
	valid bool
}

// NewBool returns NullBool with a valid value of b.
func NewBool(b bool) NullBool {
	return NullBool{
		value: b,
		valid: true,
	}
}

// NewBoolPtr returns a new bool from a pointer.
func NewBoolPtr(b *bool) NullBool {
	if b == nil {
		return NullBool{}
	}
	return NewBool(*b)
}

// Bool returns the bool value, regardless of validity. Use Value() if
// you need to know whether the value is valid.
func (n NullBool) Bool() bool {
	return n.value
}

// BoolPtr returns a pointer to the bool value, or nil if the value is not valid.
func (n NullBool) BoolPtr() *bool {
	if n.valid {
		return &n.value
	}
	return nil
}

// Value returns the value of NullBool. The value should only be considered
// valid if ok returns true.
func (n NullBool) Value() (value bool, ok bool) {
	return n.value, n.valid
}

// Equal compares the equality of two NullBools.
func (n NullBool) Equal(v NullBool) bool {
	return n.value == v.value && n.valid == v.valid
}

// MarshalJSON marshals a bool based on whether valid is true.
func (n NullBool) MarshalJSON() ([]byte, error) {
	if n.valid {
		return json.Marshal(n.value)
	}
	return []byte("null"), nil
}

// UnmarshalXML unmarshals an bool properly, as well as marshaling an empty string to nil.
func (n *NullBool) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	if err := d.DecodeElement(&v, &start); err == nil {
		if val, err := strconv.ParseBool(v); err == nil {
			*n = NewBool(val)
		} else {
			*n = NullBool{valid: false}
		}
	}
	return nil
}

// MarshalXML marshals NullBools to XML. Otherwise nothing is
// marshaled.
func (n NullBool) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if n.valid {
		return e.EncodeElement(n.value, start)
	}
	return nil
}

// NullInt is used for properly handling int types that could be null.
type NullInt struct {
	value int
	valid bool
}

// NewInt returns NullInt with a valid value of i.
func NewInt(i int) NullInt {
	return NullInt{value: i, valid: true}
}

// NewIntPtr returns a new bool from a pointer.
func NewIntPtr(i *int) NullInt {
	if i == nil {
		return NullInt{}
	}
	return NewInt(*i)
}

// Int returns the int value, regardless of validity. Use Value() if
// you need to know whether the value is valid.
func (n NullInt) Int() int {
	return n.value
}

// IntPtr returns a pointer to the int value, or nil if the value is not valid.
func (n NullInt) IntPtr() *int {
	if n.valid {
		return &n.value
	}
	return nil
}

// Value returns the value of NullInt. The value should only be considered
// valid if ok returns true.
func (n NullInt) Value() (value int, ok bool) {
	return n.value, n.valid
}

// Equal compares the equality of two NullInt.
func (n NullInt) Equal(v NullInt) bool {
	return n.value == v.value && n.valid == v.valid
}

// MarshalJSON marshals an int based on whether valid is true.
func (n NullInt) MarshalJSON() ([]byte, error) {
	if n.valid {
		return json.Marshal(n.value)
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
	*n = NewInt(v.Int)
	return nil
}

// MarshalXML marshals NullInts greater than zero to XML. Otherwise nothing is
// marshaled.
func (n NullInt) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if n.valid {
		return e.EncodeElement(n.value, start)
	}
	return nil
}

// NullFloat is used for properly handling float types that could be null. (float64 is returned)
type NullFloat struct {
	value float64
	valid bool
}

// NullFloat returns NullFloat with a valid value of f.
func NewFloat(f float64) NullFloat {
	return NullFloat{value: f, valid: true}
}

// NewFloatPtr returns a new NullFloat from a pointer to float64.
func NewFloatPtr(f *float64) NullFloat {
	if f == nil {
		return NullFloat{}
	}
	return NewFloat(*f)
}

// Float64 returns the float64 value, regardless of validity. Use Value() if
// you need to know whether the value is valid.
func (n NullFloat) Float64() float64 {
	return n.value
}

// Float64Ptr returns a pointer to the float64 value, or nil if the value is not valid.
func (n NullFloat) Float64Ptr() *float64 {
	if n.valid {
		return &n.value
	}
	return nil
}

// Value returns the value of NullFloat. The value should only be considered
// valid if ok returns true.
func (n NullFloat) Value() (value float64, ok bool) {
	return n.value, n.valid
}

// Equal compares the equality of two NullFloat.
func (n NullFloat) Equal(v NullFloat) bool {
	return n.value == v.value && n.valid == v.valid
}

// MarshalJSON marshals an float64 based on whether valid is true.
func (n NullFloat) MarshalJSON() ([]byte, error) {
	if n.valid {
		return json.Marshal(n.value)
	}
	return []byte("null"), nil
}

// UnmarshalXML unmarshals an float64 properly, as well as marshaling an empty string to nil.
func (n *NullFloat) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v struct {
		Float float64 `xml:",chardata"`
		Nil   string  `xml:"nil,attr"`
	}
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	} else if strings.EqualFold(v.Nil, "nil") || strings.EqualFold(v.Nil, "true") {
		return nil
	}
	*n = NewFloat(v.Float)
	return nil
}

// MarshalXML marshals NullFloat greater than zero to XML. Otherwise nothing is
// marshaled.
func (n NullFloat) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if n.valid {
		return e.EncodeElement(n.value, start)
	}
	return nil
}

// DateTimeFormat is the format Recurly uses to represent datetimes.
const DateTimeFormat = "2006-01-02T15:04:05Z07:00"

// NullTime is used for properly handling time.Time types that could be null.
type NullTime struct {
	value time.Time
	valid bool
}

// NewTime returns NullTime with a valid value of t. The NullTime value
// will be considered invalid if t.IsZero() returns true.
func NewTime(t time.Time) NullTime {
	if t.IsZero() {
		return NullTime{}
	}
	t = t.UTC()
	return NullTime{value: t, valid: true}
}

// NewTimePtr returns NullTime from a pointer.
func NewTimePtr(t *time.Time) NullTime {
	if t == nil || t.IsZero() {
		return NullTime{}
	}
	return NewTime(*t)
}

// Time returns the time value, regardless of validity. Use Value() if
// you need to know whether the value is valid.
func (n NullTime) Time() time.Time {
	return n.value
}

// TimePtr returns a pointer to the time.Time value, or nil if the value is not valid.
func (n NullTime) TimePtr() *time.Time {
	if n.valid {
		return &n.value
	}
	return nil
}

// Value returns the value of NullTime. The value should only be considered
// valid if ok returns true.
func (n NullTime) Value() (value time.Time, ok bool) {
	return n.value, n.valid
}

// Equal compares the equality of two NullTime.
func (n NullTime) Equal(v NullTime) bool {
	return n.value.Equal(v.value) && n.valid == v.valid
}

// UnmarshalXML unmarshals an int properly, as well as marshaling an empty string to nil.
func (n *NullTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	if err := d.DecodeElement(&v, &start); err == nil && v != "" {
		parsed, err := time.Parse(DateTimeFormat, v)
		if err != nil {
			return err
		}
		*n = NewTime(parsed)
	}
	return nil
}

// MarshalJSON method has to be added here due to embeded interface json marshal issue in Go
// with panic on nil time field
func (n NullTime) MarshalJSON() ([]byte, error) {
	if n.valid {
		return json.Marshal(n.value)
	}
	return []byte("null"), nil
}

// MarshalXML marshals times into their proper format. Otherwise nothing is
// marshaled. All times are sent in UTC.
func (n NullTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if n.valid {
		return e.EncodeElement(n.String(), start)
	}
	return nil
}

// String returns a string representation of the time in UTC using the
// DateTimeFormat constant as the format.
func (n NullTime) String() string {
	if n.valid {
		return n.value.UTC().Format(DateTimeFormat)
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
