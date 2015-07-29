package recurly

import (
	"bytes"
	"encoding/xml"
	"reflect"
	"testing"
)

func TestUnitAmount(t *testing.T) {
	given1 := UnitAmount{USD: 1000}
	given2 := UnitAmount{USD: 800, EUR: 650}
	given3 := UnitAmount{EUR: 650}
	given4 := UnitAmount{}
	given5 := UnitAmount{USD: 1}

	type s struct {
		XMLName xml.Name   `xml:"s"`
		Name    string     `xml:"name"`
		Amount  UnitAmount `xml:"amount,omitempty"`
	}

	suite := []map[string]interface{}{
		map[string]interface{}{"struct": s{XMLName: xml.Name{Local: "s"}, Name: "Bob", Amount: given1}, "expected": "<s><name>Bob</name><amount><USD>1000</USD></amount></s>"},
		map[string]interface{}{"struct": s{XMLName: xml.Name{Local: "s"}, Name: "Bob", Amount: given2}, "expected": "<s><name>Bob</name><amount><USD>800</USD><EUR>650</EUR></amount></s>"},
		map[string]interface{}{"struct": s{XMLName: xml.Name{Local: "s"}, Name: "Bob", Amount: given3}, "expected": "<s><name>Bob</name><amount><EUR>650</EUR></amount></s>"},
		map[string]interface{}{"struct": s{XMLName: xml.Name{Local: "s"}, Name: "Bob", Amount: given4}, "expected": "<s><name>Bob</name></s>"},
		map[string]interface{}{"struct": s{XMLName: xml.Name{Local: "s"}, Name: "Bob", Amount: given5}, "expected": "<s><name>Bob</name><amount><USD>1</USD></amount></s>"},
	}

	for i, test := range suite {
		str := test["struct"].(s)
		expected := test["expected"].(string)
		given := new(bytes.Buffer)
		if err := xml.NewEncoder(given).Encode(str); err != nil {
			t.Errorf("TestUnitAmount Error Suite (%d): Error encoding. Error: %s", i, err)
		}

		if expected != given.String() {
			t.Errorf("TestUnitAmount Error Suite (%d): Expected %s, given %s", i, expected, given.String())
		}

		given = bytes.NewBufferString(expected)
		var dest s
		if err := xml.NewDecoder(given).Decode(&dest); err != nil {
			t.Errorf("TestUnitAmount Error Suite (%d): Error decoding. Error: %s", i, err)
		}

		if !reflect.DeepEqual(str, dest) {
			t.Errorf("TestUnitAmount Error Suite (%d): Expected unmarshal to be %#v, given %#v", i, str, dest)
		}
	}
}
