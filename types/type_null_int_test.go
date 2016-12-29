package types

import (
	"bytes"
	"encoding/xml"
	"reflect"
	"testing"
)

func TestNullInt(t *testing.T) {
	if !reflect.DeepEqual(NewInt(1), NullInt{Int: 1, Valid: true}) {
		t.Fatalf("unexpected value: %v", NewInt(1))
	} else if !reflect.DeepEqual(NewInt(0), NullInt{Int: 0, Valid: true}) {
		t.Fatalf("unexpected value: %v", NewInt(0))
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
		} else if !reflect.DeepEqual(tt.s, dst) {
			t.Errorf("(%d): unexpected value: %v", i, dst)
		}
	}
}
