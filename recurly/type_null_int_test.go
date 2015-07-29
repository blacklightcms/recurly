package recurly

import (
	"bytes"
	"encoding/xml"
	"reflect"
	"testing"
)

func TestNullInt(t *testing.T) {
	given1 := NewInt(1)
	expected1 := NullInt{Int: 1, Valid: true}

	given2 := NewInt(0)
	expected2 := NullInt{Int: 0, Valid: true}

	given3 := NullInt{Int: 0, Valid: false}

	if !reflect.DeepEqual(expected1, given1) {
		t.Errorf("TestNullInt Error (1): Expected %#v, given %#v", expected1, given1)
	}

	if !reflect.DeepEqual(expected2, given2) {
		t.Errorf("TestNullInt Error (2): Expected %#v, given %#v", expected2, given2)
	}

	type s struct {
		XMLName xml.Name `xml:"s"`
		Name    string   `xml:"name"`
		Amount  NullInt  `xml:"amount,omitempty"`
	}

	suite := []map[string]interface{}{
		map[string]interface{}{"struct": s{XMLName: xml.Name{Local: "s"}, Name: "Bob", Amount: given1}, "expected": "<s><name>Bob</name><amount>1</amount></s>"},
		map[string]interface{}{"struct": s{XMLName: xml.Name{Local: "s"}, Name: "Bob", Amount: given2}, "expected": "<s><name>Bob</name><amount>0</amount></s>"},
		map[string]interface{}{"struct": s{XMLName: xml.Name{Local: "s"}, Name: "Bob", Amount: given3}, "expected": "<s><name>Bob</name></s>"},
	}

	for i, test := range suite {
		str := test["struct"].(s)
		expected := test["expected"].(string)
		given := new(bytes.Buffer)
		if err := xml.NewEncoder(given).Encode(str); err != nil {
			t.Errorf("TestNullInt Error Suite (%d): Error encoding. Error: %s", i, err)
		}

		if expected != given.String() {
			t.Errorf("TestNullInt Error Suite (%d): Expected %s, given %s", i, expected, given.String())
		}

		given = bytes.NewBufferString(expected)
		var dest s
		if err := xml.NewDecoder(given).Decode(&dest); err != nil {
			t.Errorf("TestNullInt Error Suite (%d): Error decoding. Error: %s", i, err)
		}

		if !reflect.DeepEqual(str, dest) {
			t.Errorf("TestNullInt Error Suite (%d): Expected unmarshal to be %#v, given %#v", i, str, dest)
		}
	}
}
