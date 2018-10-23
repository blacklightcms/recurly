package recurly

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNullInt(t *testing.T) {
	if diff := cmp.Diff(NewInt(1), NullInt{Int: 1, Valid: true}); diff != "" {
		t.Fatal(diff)
	} else if diff := cmp.Diff(NewInt(0), NullInt{Int: 0, Valid: true}); diff != "" {
		t.Fatal(diff)
	}

	jsontests := []struct {
		v        NullInt
		expected string
	}{
		{v: NewInt(5), expected: "5"},
		{v: NewInt(0), expected: "0"},
		{v: NullInt{}, expected: ""},
	}
	for _, tt := range jsontests {
		bytes, _ := json.Marshal(tt.v)
		if diff := cmp.Diff(string(bytes), tt.expected); diff != "" {
			t.Fatal(diff)
		}
	}

	type s struct {
		XMLName xml.Name `xml:"s"`
		Name    string   `xml:"name"`
		Amount  NullInt  `xml:"amount,omitempty"`
	}

	tests := []struct {
		s        s
		expected string
	}{
		{s: s{XMLName: xml.Name{Local: "s"}, Name: "Bob", Amount: NewInt(1)}, expected: "<s><name>Bob</name><amount>1</amount></s>"},
		{s: s{XMLName: xml.Name{Local: "s"}, Name: "Bob", Amount: NewInt(0)}, expected: "<s><name>Bob</name><amount>0</amount></s>"},
		{s: s{XMLName: xml.Name{Local: "s"}, Name: "Bob"}, expected: "<s><name>Bob</name></s>"},
	}

	for i, tt := range tests {
		var given bytes.Buffer
		if err := xml.NewEncoder(&given).Encode(tt.s); err != nil {
			t.Errorf("(%d): unexpected error: %v", i, err)
		} else if tt.expected != given.String() {
			t.Errorf("(%d): unexpected value: %s", i, given.String())
		}

		var dst s
		if err := xml.NewDecoder(bytes.NewBufferString(tt.expected)).Decode(&dst); err != nil {
			t.Errorf("(%d) unexpected error: %s", i, err)
		} else if diff := cmp.Diff(tt.s, dst); diff != "" {
			t.Errorf("(%d): %s", i, diff)
		}
	}
}
