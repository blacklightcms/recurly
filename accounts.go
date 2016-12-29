package recurly

import (
	"encoding/xml"
	"time"
)

type (
	// Account represents an individual account on your site
	Account struct {
		XMLName          xml.Name `xml:"account"`
		Code             string   `xml:"account_code,omitempty"`
		State            string   `xml:"state,omitempty"`
		Username         string   `xml:"username,omitempty"`
		Email            string   `xml:"email,omitempty"`
		FirstName        string   `xml:"first_name,omitempty"`
		LastName         string   `xml:"last_name,omitempty"`
		CompanyName      string   `xml:"company_name,omitempty"`
		VATNumber        string   `xml:"vat_number,omitempty"`
		TaxExempt        NullBool `xml:"tax_exempt,omitempty"`
		BillingInfo      *Billing `xml:"billing_info,omitempty"`
		Address          Address  `xml:"address,omitempty"`
		AcceptLanguage   string   `xml:"accept_language,omitempty"`
		HostedLoginToken string   `xml:"hosted_login_token,omitempty"`
		CreatedAt        NullTime `xml:"created_at,omitempty"`
	}

	// Address is used for embedded addresses within other structs.
	Address struct {
		Address  string `xml:"address1,omitempty"`
		Address2 string `xml:"address2,omitempty"`
		City     string `xml:"city,omitempty"`
		State    string `xml:"state,omitempty"`
		Zip      string `xml:"zip,omitempty"`
		Country  string `xml:"country,omitempty"`
		Phone    string `xml:"phone,omitempty"`
	}

	// Note holds account notes.
	Note struct {
		XMLName   xml.Name  `xml:"note"`
		Message   string    `xml:"message,omitempty"`
		CreatedAt time.Time `xml:"created_at,omitempty"`
	}
)

// MarshalXML ensures addresses marshal to nil if empty without the need
// to use pointers.
func (a Address) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if a.Address == "" && a.Address2 == "" && a.City == "" && a.State == "" && a.Zip == "" && a.Country == "" && a.Phone == "" {
		return nil
	}

	e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "address"}})
	if a.Address != "" {
		s := xml.StartElement{Name: xml.Name{Local: "address1"}}
		e.EncodeToken(s)
		e.EncodeToken(xml.CharData([]byte(a.Address)))
		e.EncodeToken(xml.EndElement{Name: s.Name})
	}

	if a.Address2 != "" {
		s := xml.StartElement{Name: xml.Name{Local: "address2"}}
		e.EncodeToken(s)
		e.EncodeToken(xml.CharData([]byte(a.Address2)))
		e.EncodeToken(xml.EndElement{Name: s.Name})
	}

	if a.City != "" {
		s := xml.StartElement{Name: xml.Name{Local: "city"}}
		e.EncodeToken(s)
		e.EncodeToken(xml.CharData([]byte(a.City)))
		e.EncodeToken(xml.EndElement{Name: s.Name})
	}

	if a.State != "" {
		s := xml.StartElement{Name: xml.Name{Local: "state"}}
		e.EncodeToken(s)
		e.EncodeToken(xml.CharData([]byte(a.State)))
		e.EncodeToken(xml.EndElement{Name: s.Name})
	}

	if a.Zip != "" {
		s := xml.StartElement{Name: xml.Name{Local: "zip"}}
		e.EncodeToken(s)
		e.EncodeToken(xml.CharData([]byte(a.Zip)))
		e.EncodeToken(xml.EndElement{Name: s.Name})
	}

	if a.Country != "" {
		s := xml.StartElement{Name: xml.Name{Local: "country"}}
		e.EncodeToken(s)
		e.EncodeToken(xml.CharData([]byte(a.Country)))
		e.EncodeToken(xml.EndElement{Name: s.Name})
	}

	if a.Phone != "" {
		s := xml.StartElement{Name: xml.Name{Local: "phone"}}
		e.EncodeToken(s)
		e.EncodeToken(xml.CharData([]byte(a.Phone)))
		e.EncodeToken(xml.EndElement{Name: s.Name})
	}

	e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "address"}})

	return nil
}
