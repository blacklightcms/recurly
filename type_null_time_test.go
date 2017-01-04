package recurly

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestNulTime(t *testing.T) {
	t1, _ := time.Parse(DateTimeFormat, "2011-10-25T12:00:00-07:00")
	given1 := NewTime(t1)
	utc1 := t1.UTC()
	expected1 := NullTime{Time: &utc1}

	t2 := t1.AddDate(0, 3, 7)
	given2 := NewTime(t2)
	utc2 := t2.UTC()
	expected2 := NullTime{Time: &utc2}

	given3 := NullTime{Time: nil}

	if given3.String() != "" {
		t.Fatalf("expected nil time to print empty string, given %s", given3.String())
	} else if !reflect.DeepEqual(expected1, given1) {
		t.Fatalf("(1): Expected %#v, given %#v", expected1, given1)
	} else if !reflect.DeepEqual(expected2, given2) {
		t.Fatalf("(2): Expected %#v, given %#v", expected2, given2)
	}

	type s struct {
		XMLName xml.Name `xml:"s"`
		Name    string   `xml:"name"`
		Stamp   NullTime `xml:"stamp,omitempty"`
	}

	suite := []map[string]interface{}{
		{"struct": s{XMLName: xml.Name{Local: "s"}, Name: "A", Stamp: given1}, "expected": fmt.Sprintf("<s><name>A</name><stamp>%s</stamp></s>", utc1.Format(DateTimeFormat))},
		{"struct": s{XMLName: xml.Name{Local: "s"}, Name: "B", Stamp: given2}, "expected": fmt.Sprintf("<s><name>B</name><stamp>%s</stamp></s>", utc2.Format(DateTimeFormat))},
		{"struct": s{XMLName: xml.Name{Local: "s"}, Name: "C", Stamp: given3}, "expected": "<s><name>C</name></s>"},
	}

	for i, test := range suite {
		str := test["struct"].(s)
		expected := test["expected"].(string)
		given := new(bytes.Buffer)
		if err := xml.NewEncoder(given).Encode(str); err != nil {
			t.Fatalf("(%d): Error encoding. Error: %s", i, err)
		}

		if expected != given.String() {
			t.Fatalf("(%d): Expected %s, given %s", i, expected, given.String())
		}

		given = bytes.NewBufferString(expected)
		var dest s
		if err := xml.NewDecoder(given).Decode(&dest); err != nil {
			t.Fatalf("(%d): Error decoding. Error: %s", i, err)
		}

		if !reflect.DeepEqual(str, dest) {
			t.Fatalf("(%d): Expected unmarshal to be %#v, given %#v", i, str, dest)
		}
	}

	// Decode Error
	var dest s
	errBuf := bytes.NewBufferString("<s><name>B</name><stamp>ABC</stamp></s>")
	if err := xml.NewDecoder(errBuf).Decode(&dest); err == nil {
		t.Fatal("expected time.Parse error. None given.", err)
	}

	if dest.Stamp.String() != "" {
		t.Fatalf("expected time.Parse error to result in empty String(), given %s", dest.Stamp.String())
	}
}
