package types

import (
	"bytes"
	"encoding/xml"
	"reflect"
	"testing"
)

func TestUnitAmount(t *testing.T) {
	type s struct {
		Amount UnitAmount `xml:"amount,omitempty"`
	}

	tests := []struct {
		v        s
		expected string
	}{
		{v: s{Amount: UnitAmount{USD: 1000}}, expected: "<s><amount><USD>1000</USD></amount></s>"},
		{v: s{Amount: UnitAmount{USD: 800, EUR: 650}}, expected: "<s><amount><USD>800</USD><EUR>650</EUR></amount></s>"},
		{v: s{Amount: UnitAmount{EUR: 650}}, expected: "<s><amount><EUR>650</EUR></amount></s>"},
		{v: s{}, expected: "<s></s>"},
		{v: s{Amount: UnitAmount{USD: 1}}, expected: "<s><amount><USD>1</USD></amount></s>"},
	}

	for _, tt := range tests {
		var given bytes.Buffer
		if err := xml.NewEncoder(&given).Encode(tt.expected); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		buf := bytes.NewBufferString(tt.expected)
		var dst s
		if err := xml.NewDecoder(buf).Decode(&dst); err != nil {
			t.Fatalf("unexpected error: %v", err)
		} else if !reflect.DeepEqual(tt.v, dst) {
			t.Fatalf("unexpected value: %T %v", dst, dst)
		}
	}
}
