package recurly

import (
	"bytes"
	"encoding/xml"
	"reflect"
	"testing"
)

func TestTypeHREFUnmarshal(t *testing.T) {
	type h struct {
		XMLName xml.Name `xml:"foo"`
		Account href     `xml:"account"`
		Invoice href     `xml:"invoice"`
	}

	expected := h{
		XMLName: xml.Name{Local: "foo"},
		Account: href{
			HREF: "https://your-subdomain.recurly.com/v2/accounts/100",
			Code: "100",
		},
		Invoice: href{
			HREF: "https://your-subdomain.recurly.com/v2/invoices/1108",
			Code: "1108",
		},
	}

	str := bytes.NewBufferString(`<foo><account href="https://your-subdomain.recurly.com/v2/accounts/100"/>
    <invoice href="https://your-subdomain.recurly.com/v2/invoices/1108"/></foo>`)

	var given h
	if err := xml.NewDecoder(str).Decode(&given); err != nil {
		t.Errorf("TestTypeHREFUnmarshal Error: error decoding xml. Err: %s", err)
	}

	if !reflect.DeepEqual(expected, given) {
		t.Errorf("TestTypeHREFUnmarshal Error: Expected unmarshal to be %#v, given %#v", expected, given)
	}
}

func TestTypeHREFMarshal(t *testing.T) {
	type h struct {
		XMLName xml.Name `xml:"foo"`
		Name    string   `xml:"name"`
		Account href     `xml:"account"`
		Invoice href     `xml:"invoice"`
	}

	v := h{
		Name: "Bob",
		Account: href{
			HREF: "https://your-subdomain.recurly.com/v2/accounts/100",
			Code: "100",
		},
		Invoice: href{
			HREF: "https://your-subdomain.recurly.com/v2/invoices/1108",
			Code: "1108",
		},
	}

	expected := `<foo><name>Bob</name></foo>`

	given := new(bytes.Buffer)
	if err := xml.NewEncoder(given).Encode(v); err != nil {
		t.Errorf("TestTypeHREFMarshal Error: error encoding xml. Err: %s", err)
	}

	if expected != given.String() {
		t.Errorf("TestTypeHREFMarshal Error: Expected marshal to be %s, given %s", expected, given.String())
	}
}
