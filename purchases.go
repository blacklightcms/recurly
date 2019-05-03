package recurly

import "encoding/xml"

// Purchase represents an individual checkout holding at least one
// subscription OR one adjustment
type Purchase struct {
	XMLName               xml.Name               `xml:"purchase"`
	Account               Account                `xml:"account,omitempty"`
	Adjustments           []Adjustment           `xml:"adjustments>adjustment,omitempty"`
	CollectionMethod      string                 `xml:"collection_method,omitempty"`
	Currency              string                 `xml:"currency"`
	PONumber              string                 `xml:"po_number,omitempty"`
	NetTerms              NullInt                `xml:"net_terms,omitempty"`
	GiftCard              string                 `xml:"gift_card>redemption_code,omitempty"`
	CouponCodes           []string               `xml:"coupon_codes>coupon_code,omitempty"`
	Subscriptions         []PurchaseSubscription `xml:"subscriptions>subscription,omitempty"`
	CustomerNotes         string                 `xml:"customer_notes,omitempty"`
	TermsAndConditions    string                 `xml:"terms_and_conditions,omitempty"`
	VATReverseChargeNotes string                 `xml:"vat_reverse_charge_notes,omitempty"`
	ShippingAddressID     int64                  `xml:"shipping_address_id,omitempty"`
	GatewayCode           string                 `xml:"gateway_code,omitempty"`
}

// PurchaseSubscription represents a subscription to purchase
// some new subscription fields are moved to purchase object level
// recurly does not accept these fields at subscription level
type PurchaseSubscription struct {
	XMLName              xml.Name             `xml:"subscription"`
	PlanCode             string               `xml:"plan_code"`
	SubscriptionAddOns   *[]SubscriptionAddOn `xml:"subscription_add_ons>subscription_add_on,omitempty"`
	UnitAmountInCents    int                  `xml:"unit_amount_in_cents,omitempty"`
	Quantity             int                  `xml:"quantity,omitempty"`
	TrialEndsAt          NullTime             `xml:"trial_ends_at,omitempty"`
	StartsAt             NullTime             `xml:"starts_at,omitempty"`
	TotalBillingCycles   int                  `xml:"total_billing_cycles,omitempty"`
	RenewalBillingCycles NullInt              `xml:"renewal_billing_cycles,omitempty"`
	NextBillDate         NullTime             `xml:"next_bill_date,omitempty"`
	AutoRenew            bool                 `xml:"auto_renew,omitempty"`
	CustomFields         *CustomFields        `xml:"custom_fields,omitempty"`
}
