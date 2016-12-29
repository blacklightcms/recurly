package recurly

import "encoding/xml"

// Plan represents an individual plan on your site.
type Plan struct {
	XMLName                  xml.Name   `xml:"plan"`
	Code                     string     `xml:"plan_code,omitempty"`
	Name                     string     `xml:"name"`
	Description              string     `xml:"description,omitempty"`
	SuccessURL               string     `xml:"success_url,omitempty"`
	CancelURL                string     `xml:"cancel_url,omitempty"`
	DisplayDonationAmounts   NullBool   `xml:"display_donation_amounts,omitempty"`
	DisplayQuantity          NullBool   `xml:"display_quantity,omitempty"`
	DisplayPhoneNumber       NullBool   `xml:"display_phone_number,omitempty"`
	BypassHostedConfirmation NullBool   `xml:"bypass_hosted_confirmation,omitempty"`
	UnitName                 string     `xml:"unit_name,omitempty"`
	PaymentPageTOSLink       string     `xml:"payment_page_tos_link,omitempty"`
	IntervalUnit             string     `xml:"plan_interval_unit,omitempty"`
	IntervalLength           int        `xml:"plan_interval_length,omitempty"`
	TrialIntervalUnit        string     `xml:"trial_interval_unit,omitempty"`
	TrialIntervalLength      int        `xml:"trial_interval_length,omitempty"`
	TotalBillingCycles       NullInt    `xml:"total_billing_cycles,omitempty"`
	AccountingCode           string     `xml:"accounting_code,omitempty"`
	CreatedAt                NullTime   `xml:"created_at,omitempty"`
	TaxExempt                NullBool   `xml:"tax_exempt,omitempty"`
	TaxCode                  string     `xml:"tax_code,omitempty"`
	UnitAmountInCents        UnitAmount `xml:"unit_amount_in_cents"`
	SetupFeeInCents          UnitAmount `xml:"setup_fee_in_cents,omitempty"`
}
