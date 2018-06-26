package recurly

import "encoding/xml"

// Purchase represents an individual checkout holding at least one
// subscription OR one adjustment
type Purchase struct {
	XMLName               xml.Name          `xml:"purchase"`
	Account               Account           `xml:"account,omitempty"`
	Adjustments           []Adjustment      `xml:"adjustments>adjustment,omitempty"`
	CollectionMethod      string            `xml:"collection_method,omitempty"`
	Currency              string            `xml:"currency"`
	PONumber              string            `xml:"po_number,omitempty"`
	NetTerms              NullInt           `xml:"net_terms,omitempty"`
	GiftCard              string            `xml:"gift_card>redemption_code,omitempty"`
	CouponCodes           []string          `xml:"coupon_codes>coupon_code,omitempty"`
	Subscriptions         []NewSubscription `xml:"subscriptions>subscription,omitempty"`
	CustomerNotes         string            `xml:"customer_notes,omitempty"`
	TermsAndConditions    string            `xml:"terms_and_conditions,omitempty"`
	VATReverseChargeNotes string            `xml:"vat_reverse_charge_notes,omitempty"`
	ShippingAddressID     int64             `xml:"shipping_address_id,omitempty"`
}
