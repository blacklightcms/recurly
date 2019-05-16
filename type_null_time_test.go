package recurly

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestNulTime(t *testing.T) {
	t1, _ := time.Parse(DateTimeFormat, "2011-10-25T12:00:00-07:00")
	given1 := NewTime(t1)
	utc1 := t1.UTC()
	expected1 := NullTime{Time: utc1, Valid: true}

	t2 := t1.AddDate(0, 3, 7)
	given2 := NewTime(t2)
	utc2 := t2.UTC()
	expected2 := NullTime{Time: utc2, Valid: true}

	given3 := NullTime{}

	if given3.String() != "" {
		t.Fatalf("expected nil time to print empty string, given %s", given3.String())
	} else if diff := cmp.Diff(expected1, given1); diff != "" {
		t.Fatal(diff)
	} else if diff := cmp.Diff(expected2, given2); diff != "" {
		t.Fatal(diff)
	}

	// check marshal interface
	b, err := json.Marshal(given1)
	if err != nil {
		t.Fatalf("json marshaling error %s", err.Error())
	} else {
		tb, _ := json.Marshal(given1.Time)
		if diff := cmp.Diff(tb, b); diff != "" {
			t.Fatal(diff)
		}
	}

	b, err = json.Marshal(given3)
	if err != nil {
		t.Fatalf("json marshaling error %s", err.Error())
	} else if !bytes.Equal(b, []byte("null")) {
		t.Fatal("null time not marshaled to null")
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

		if diff := cmp.Diff(str, dest); diff != "" {
			t.Fatalf("(%d): %s", i, diff)
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
