package recurly

import (
	"bytes"
	"encoding/xml"
	"reflect"
	"testing"
)

func TestNullBool(t *testing.T) {
	given1 := NewBool(true)
	expected1 := NullBool{Bool: true, Valid: true}

	given2 := NewBool(false)
	expected2 := NullBool{Bool: false, Valid: true}

	given3 := NullBool{Bool: false, Valid: false}

	if !reflect.DeepEqual(expected1, given1) {
		t.Errorf("TestNullBool Error (1): Expected %#v, given %#v", expected1, given1)
	}

	if !given1.Is(true) {
		t.Errorf("TestNullBool Error (1): Expected Is() to return %v, given %v", true, given1.Is(true))
	}

	if !reflect.DeepEqual(expected2, given2) {
		t.Errorf("TestNullBool Error (2): Expected %#v, given %#v", expected2, given2)
	}

	if !given2.Is(false) {
		t.Errorf("TestNullBool Error (2): Expected Is() to return %t, given %t", false, given2.Is(false))
	}

	if given3.Is(false) {
		t.Errorf("TestNullBool Error (3): Expected Is() on invalid bool to return %t, given %t", false, given3.Is(false))
	}

	type s struct {
		XMLName xml.Name `xml:"s"`
		Name    string   `xml:"name"`
		Exempt  NullBool `xml:"exempt,omitempty"`
	}

	suite := []map[string]interface{}{
		map[string]interface{}{"struct": s{XMLName: xml.Name{Local: "s"}, Name: "Bob", Exempt: given1}, "expected": "<s><name>Bob</name><exempt>true</exempt></s>"},
		map[string]interface{}{"struct": s{XMLName: xml.Name{Local: "s"}, Name: "Bob", Exempt: given2}, "expected": "<s><name>Bob</name><exempt>false</exempt></s>"},
		map[string]interface{}{"struct": s{XMLName: xml.Name{Local: "s"}, Name: "Bob", Exempt: given3}, "expected": "<s><name>Bob</name></s>"},
	}

	for i, test := range suite {
		str := test["struct"].(s)
		expected := test["expected"].(string)
		given := new(bytes.Buffer)
		if err := xml.NewEncoder(given).Encode(str); err != nil {
			t.Errorf("TestNullBool Error Suite (%d): Error encoding. Error: %s", i, err)
		}

		if expected != given.String() {
			t.Errorf("TestNullBool Error Suite (%d): Expected %s, given %s", i, expected, given.String())
		}

		given = bytes.NewBufferString(expected)
		var dest s
		if err := xml.NewDecoder(given).Decode(&dest); err != nil {
			t.Errorf("TestNullBool Error Suite (%d): Error decoding. Error: %s", i, err)
		}

		if !reflect.DeepEqual(str, dest) {
			t.Errorf("TestNullBool Error Suite (%d): Expected unmarshal to be %#v, given %#v", i, str, dest)
		}
	}
}
