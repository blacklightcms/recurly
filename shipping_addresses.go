package recurly

// ShippingAddress represents a shipping address
type ShippingAddress struct {
	Nickname  string `xml:"nickname,omitempty"`
	Address   string `xml:"address1,omitempty"`
	Address2  string `xml:"address2,omitempty"`
	Company   string `xml:"company,omitempty"`
	City      string `xml:"city,omitempty"`
	State     string `xml:"state,omitempty"`
	Zip       string `xml:"zip,omitempty"`
	Country   string `xml:"country,omitempty"`
	Phone     string `xml:"phone,omitempty"`
	Email     string `xml:"email,omitempty"`
	VATNumber string `xml:"vat_number,omitempty"`
}
