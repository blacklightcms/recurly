package recurly_test

import (
	"encoding/json"
	"encoding/xml"
	"testing"

	"github.com/blacklightcms/recurly"
	"github.com/google/go-cmp/cmp"
)

func TestXML_NullBool(t *testing.T) {
	if b := recurly.NewBool(true); !b.Is(true) {
		t.Fatal("expected true")
	} else if b.Is(false) {
		t.Fatal("expected false")
	} else if diff := cmp.Diff(b, recurly.NullBool{Bool: true, Valid: true}); diff != "" {
		t.Fatal(diff)
	}

	if b := recurly.NewBool(false); !b.Is(false) {
		t.Fatal("expected true")
	} else if b.Is(true) {
		t.Fatal("expected false")
	} else if diff := cmp.Diff(b, recurly.NullBool{Bool: false, Valid: true}); diff != "" {
		t.Fatal(diff)
	}

	type testStruct struct {
		XMLName xml.Name         `xml:"test"`
		Value   recurly.NullBool `xml:"b"`
	}

	t.Run("Encode", func(t *testing.T) {
		for i, tt := range []struct {
			b      bool
			valid  bool
			expect string
		}{
			{b: true, valid: true, expect: `<test><b>true</b></test>`},
			{b: false, valid: true, expect: `<test><b>false</b></test>`},
			{b: true, valid: false, expect: `<test></test>`},
			{b: false, valid: false, expect: `<test></test>`},
		} {
			value := recurly.NewBool(tt.b)
			value.Valid = tt.valid
			if xml, err := xml.Marshal(testStruct{Value: value}); err != nil {
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
			{expect: recurly.NullBool{Bool: false, Valid: false}, input: `<test></test>`},
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
			{b: recurly.NullBool{Bool: true, Valid: false}, expect: "null"},
			{b: recurly.NullBool{Bool: false, Valid: false}, expect: "null"},
		} {
			if b, err := json.Marshal(tt.b); err != nil {
				t.Fatalf("%d %#v", i, err)
			} else if string(b) != tt.expect {
				t.Fatalf("%d %s", i, string(b))
			}
		}
	})
}

func TestXML_NullInt(t *testing.T) {
	b := recurly.NewInt(1)
	if diff := cmp.Diff(b, recurly.NullInt{Int: 1, Valid: true}); diff != "" {
		t.Fatal(diff)
	}

	b = recurly.NewInt(0)
	if diff := cmp.Diff(b, recurly.NullInt{Int: 0, Valid: true}); diff != "" {
		t.Fatal(diff)
	}

	type testStruct struct {
		XMLName xml.Name        `xml:"test"`
		Value   recurly.NullInt `xml:"i"`
	}

	t.Run("Encode", func(t *testing.T) {
		for i, tt := range []struct {
			i      int
			valid  bool
			expect string
		}{
			{i: 1, valid: true, expect: `<test><i>1</i></test>`},
			{i: 0, valid: true, expect: `<test><i>0</i></test>`},
			{i: 1, valid: false, expect: `<test></test>`},
			{i: 0, valid: false, expect: `<test></test>`},
		} {
			value := recurly.NewInt(tt.i)
			value.Valid = tt.valid
			if xml, err := xml.Marshal(testStruct{Value: value}); err != nil {
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
			{expect: recurly.NullInt{Int: 0, Valid: false}, input: `<test></test>`},
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

func TestXML_NullTime(t *testing.T) {
	v := MustParseTime("2011-10-25T12:00:00Z")

	rt := recurly.NewTime(v)
	if diff := cmp.Diff(rt, recurly.NullTime{Time: &v}); diff != "" {
		t.Fatal(diff)
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
