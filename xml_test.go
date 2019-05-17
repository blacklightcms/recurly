package recurly_test

import (
	"encoding/json"
	"encoding/xml"
	"testing"

	"github.com/blacklightcms/recurly"
	"github.com/google/go-cmp/cmp"
)

func TestXML_NullBool(t *testing.T) {
	t.Run("ZeroValue", func(t *testing.T) {
		var b recurly.NullBool
		if value, ok := b.Value(); ok {
			t.Fatal("expected ok to be false")
		} else if value != false {
			t.Fatal("expected false")
		}
	})

	t.Run("True", func(t *testing.T) {
		b := recurly.NewBool(true)
		if value, ok := b.Value(); !ok {
			t.Fatal("expected ok to be true")
		} else if value != true {
			t.Fatal("expected true")
		}
	})

	t.Run("False", func(t *testing.T) {
		b := recurly.NewBool(false)
		if value, ok := b.Value(); !ok {
			t.Fatal("expected ok to be true")
		} else if value != false {
			t.Fatal("expected false")
		}
	})

	type testStruct struct {
		XMLName xml.Name         `xml:"test"`
		Value   recurly.NullBool `xml:"b"`
	}

	t.Run("Encode", func(t *testing.T) {
		for i, tt := range []struct {
			value  recurly.NullBool
			expect string
		}{
			{value: recurly.NewBool(true), expect: `<test><b>true</b></test>`},
			{value: recurly.NewBool(false), expect: `<test><b>false</b></test>`},
			{expect: `<test></test>`}, // zero value
		} {
			if xml, err := xml.Marshal(testStruct{Value: tt.value}); err != nil {
				t.Fatalf("%d %#v", i, err)
			} else if string(xml) != tt.expect {
				t.Fatalf("%d %s", i, string(xml))
			}
		}
	})

	t.Run("Decode", func(t *testing.T) {
		for i, tt := range []struct {
			expect recurly.NullBool
			input  string
		}{
			{expect: recurly.NewBool(true), input: `<test><b>true</b></test>`},
			{expect: recurly.NewBool(false), input: `<test><b>false</b></test>`},
			{input: `<test></test>`},
		} {
			var dst testStruct
			if err := xml.Unmarshal([]byte(tt.input), &dst); err != nil {
				t.Fatalf("%d %#v", i, err)
			} else if diff := cmp.Diff(testStruct{XMLName: xml.Name{Local: "test"}, Value: tt.expect}, dst); diff != "" {
				t.Fatalf("%d %s", i, diff)
			}
		}
	})

	t.Run("JSON", func(t *testing.T) {
		for i, tt := range []struct {
			b      recurly.NullBool
			expect string
		}{
			{b: recurly.NewBool(true), expect: "true"},
			{b: recurly.NewBool(false), expect: "false"},
			{expect: "null"}, // zero value
		} {
			if b, err := json.Marshal(tt.b); err != nil {
				t.Fatalf("%d %#v", i, err)
			} else if string(b) != tt.expect {
				t.Fatalf("%d %s", i, string(b))
			}
		}
	})
}

func TestXML_NullBoolPtr(t *testing.T) {
	boolVal := true

	b := recurly.NewBoolPtr(&boolVal)
	if value, ok := b.Value(); !ok {
		t.Fatal("expected ok to be true")
	} else if value != true {
		t.Fatal("expected true")
	}
}

func TestXML_NullInt(t *testing.T) {
	t.Run("ZeroValue", func(t *testing.T) {
		var i recurly.NullInt
		if value, ok := i.Value(); ok {
			t.Fatal("expected ok to be false")
		} else if value != 0 {
			t.Fatalf("unexpected value: %d", value)
		}
	})

	i := recurly.NewInt(1)
	if value, ok := i.Value(); !ok {
		t.Fatal("expected ok to be true")
	} else if value != 1 {
		t.Fatalf("unexpected value: %d", value)
	}

	i = recurly.NewInt(0)
	if value, ok := i.Value(); !ok {
		t.Fatal("expected ok to be true")
	} else if value != 0 {
		t.Fatalf("unexpected value: %d", value)
	}

	type testStruct struct {
		XMLName xml.Name        `xml:"test"`
		Value   recurly.NullInt `xml:"i"`
	}

	t.Run("Encode", func(t *testing.T) {
		for i, tt := range []struct {
			value  recurly.NullInt
			expect string
		}{
			{value: recurly.NewInt(1), expect: `<test><i>1</i></test>`},
			{value: recurly.NewInt(0), expect: `<test><i>0</i></test>`},
			{value: recurly.NewInt(-1), expect: `<test><i>-1</i></test>`},
			{expect: `<test></test>`}, // zero value
		} {
			if xml, err := xml.Marshal(testStruct{Value: tt.value}); err != nil {
				t.Fatalf("%d %#v", i, err)
			} else if string(xml) != tt.expect {
				t.Fatalf("%d %s", i, string(xml))
			}
		}
	})

	t.Run("Decode", func(t *testing.T) {
		for i, tt := range []struct {
			expect recurly.NullInt
			input  string
		}{
			{expect: recurly.NewInt(1), input: `<test><i>1</i></test>`},
			{expect: recurly.NewInt(0), input: `<test><i>0</i></test>`},
			{input: `<test></test>`}, // zero value
		} {
			var dst testStruct
			if err := xml.Unmarshal([]byte(tt.input), &dst); err != nil {
				t.Fatalf("%d %#v", i, err)
			} else if diff := cmp.Diff(testStruct{XMLName: xml.Name{Local: "test"}, Value: tt.expect}, dst); diff != "" {
				t.Fatalf("%d %s", i, diff)
			}
		}
	})
}

func TestXML_NullIntPtr(t *testing.T) {
	i := recurly.NewIntPtr(nil)
	if value, ok := i.Value(); ok {
		t.Fatal("expected ok to be false")
	} else if value != 0 {
		t.Fatalf("unexpected value: %d", value)
	}

	intVal := 1
	i = recurly.NewIntPtr(&intVal)
	if value, ok := i.Value(); !ok {
		t.Fatal("expected ok to be true")
	} else if value != 1 {
		t.Fatalf("unexpected value: %d", value)
	}
}

func TestXML_NullTime(t *testing.T) {
	t.Run("ZeroValue", func(t *testing.T) {
		var rt recurly.NullTime
		if value, ok := rt.Value(); ok {
			t.Fatal("expected ok to be false")
		} else if !value.IsZero() {
			t.Fatalf("expected zero time: %s", value.String())
		}
	})

	v := MustParseTime("2011-10-25T12:00:00Z")

	rt := recurly.NewTime(v)
	if value, ok := rt.Value(); !ok {
		t.Fatal("expected ok to be true")
	} else if !value.Equal(v) {
		t.Fatalf("unexpected value: %v", value)
	}

	t.Run("Encode", func(t *testing.T) {
		// Value
		if b, err := xml.Marshal(struct {
			XMLName xml.Name         `xml:"test"`
			Time    recurly.NullTime `xml:"time"`
		}{
			Time: recurly.NewTime(v),
		}); err != nil {
			t.Fatal(err)
		} else if string(b) != `<test><time>2011-10-25T12:00:00Z</time></test>` {
			t.Fatal(string(b))
		}

		// No value.
		if b, err := xml.Marshal(struct {
			XMLName xml.Name         `xml:"test"`
			Time    recurly.NullTime `xml:"time"`
		}{
			Time: recurly.NullTime{},
		}); err != nil {
			t.Fatal(err)
		} else if string(b) != `<test></test>` {
			t.Fatal(string(b))
		}
	})

	t.Run("Decode", func(t *testing.T) {
		// Value
		var s struct {
			XMLName xml.Name         `xml:"test"`
			Time    recurly.NullTime `xml:"time"`
		}
		if err := xml.Unmarshal([]byte(`<test><time>2011-10-25T12:00:00Z</time></test>`), &s); err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(s.Time, recurly.NewTime(v)); diff != "" {
			t.Fatal(diff)
		}

		// Reset and try no value.
		s.XMLName = xml.Name{}
		s.Time = recurly.NullTime{}
		if err := xml.Unmarshal([]byte(`<test></test>`), &s); err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(s.Time, recurly.NullTime{}); diff != "" {
			t.Fatal(diff)
		}

		// Reset and try with tag but no value. Recurly often uses this format.
		s.XMLName = xml.Name{}
		s.Time = recurly.NullTime{}
		if err := xml.Unmarshal([]byte(`<test><time nil="nil"/></test>`), &s); err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(s.Time, recurly.NullTime{}); diff != "" {
			t.Fatal(diff)
		}
	})
}

func TestXML_NullTimePtr(t *testing.T) {
	rt := recurly.NewTimePtr(nil)
	if value, ok := rt.Value(); ok {
		t.Fatal("expected ok to be false")
	} else if !value.IsZero() {
		t.Fatalf("expected zero time: %s", value.String())
	}

	v := MustParseTime("2011-10-25T12:00:00Z")
	rt = recurly.NewTimePtr(&v)
	if value, ok := rt.Value(); !ok {
		t.Fatal("expected ok to be true")
	} else if !value.Equal(v) {
		t.Fatalf("unexpected value: %v", value)
	}
}
