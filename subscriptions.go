package recurly

import (
	"encoding/xml"
	"sort"
	"strings"
)

// SanitizeUUID returns the uuid without dashes.
func SanitizeUUID(id string) string {
	return strings.TrimSpace(strings.Replace(id, "-", "", -1))
}

const (
	// SubscriptionStateActive represents subscriptions that are valid for the
	// current time. This includes subscriptions in a trial period
	SubscriptionStateActive = "active"

	// SubscriptionStateCanceled are subscriptions that are valid for
	// the current time but will not renew because a cancelation was requested
	SubscriptionStateCanceled = "canceled"

	// SubscriptionStateExpired are subscriptions that have expired and are no longer valid
	SubscriptionStateExpired = "expired"

	// SubscriptionStateFuture are subscriptions that will start in the
	// future, they are not active yet
	SubscriptionStateFuture = "future"

	// SubscriptionStateInTrial are subscriptions that are active or canceled
	// and are in a trial period
	SubscriptionStateInTrial = "in_trial"

	// SubscriptionStateLive are all subscriptions that are not expired
	SubscriptionStateLive = "live"

	// SubscriptionStatePastDue are subscriptions that are active or canceled
	// and have a past-due invoice
	SubscriptionStatePastDue = "past_due"

	// SubscriptionStatePaused are subscriptions that are in a paused state
	// and will not be billed for the set RemainingPauseCycles
	SubscriptionStatePaused = "paused"
)

// Subscription represents an individual subscription.
type Subscription struct {
	XMLName                xml.Name             `xml:"subscription"`
	Plan                   NestedPlan           `xml:"plan,omitempty"`
	AccountCode            string               `xml:"-"`
	InvoiceNumber          int                  `xml:"-"`
	UUID                   string               `xml:"uuid,omitempty"`
	State                  string               `xml:"state,omitempty"`
	UnitAmountInCents      int                  `xml:"unit_amount_in_cents,omitempty"`
	Currency               string               `xml:"currency,omitempty"`
	Quantity               int                  `xml:"quantity,omitempty"`
	TotalAmountInCents     int                  `xml:"total_amount_in_cents,omitempty"`
	ActivatedAt            NullTime             `xml:"activated_at,omitempty"`
	CanceledAt             NullTime             `xml:"canceled_at,omitempty"`
	ExpiresAt              NullTime             `xml:"expires_at,omitempty"`
	CurrentPeriodStartedAt NullTime             `xml:"current_period_started_at,omitempty"`
	CurrentPeriodEndsAt    NullTime             `xml:"current_period_ends_at,omitempty"`
	TrialStartedAt         NullTime             `xml:"trial_started_at,omitempty"`
	TrialEndsAt            NullTime             `xml:"trial_ends_at,omitempty"`
	TaxInCents             int                  `xml:"tax_in_cents,omitempty"`
	TaxType                string               `xml:"tax_type,omitempty"`
	TaxRegion              string               `xml:"tax_region,omitempty"`
	TaxRate                float64              `xml:"tax_rate,omitempty"`
	PONumber               string               `xml:"po_number,omitempty"`
	NetTerms               NullInt              `xml:"net_terms,omitempty"`
	SubscriptionAddOns     []SubscriptionAddOn  `xml:"subscription_add_ons>subscription_add_on,omitempty"`
	CurrentTermStartedAt   NullTime             `xml:"current_term_started_at,omitempty"`
	CurrentTermEndsAt      NullTime             `xml:"current_term_ends_at,omitempty"`
	PendingSubscription    *PendingSubscription `xml:"pending_subscription,omitempty"`
	Invoice                *Invoice             `xml:"-"`
	RemainingPauseCycles   int                  `xml:"remaining_pause_cycles,omitempty"`
	CollectionMethod       string               `xml:"collection_method"`
	AutoRenew              bool                 `xml:"auto_renew,omitempty"`
	RenewalBillingCycles   NullInt              `xml:"renewal_billing_cycles"`
	CustomFields           *CustomFields        `xml:"custom_fields,omitempty"`
}

// UnmarshalXML unmarshals transactions and handles intermediary state during unmarshaling
// for types like href.
func (s *Subscription) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type subscriptionAlias Subscription
	var v struct {
		subscriptionAlias
		XMLName           xml.Name           `xml:"subscription"`
		AccountCode       hrefString         `xml:"account"`
		InvoiceNumber     hrefInt            `xml:"invoice"`
		InvoiceCollection *InvoiceCollection `xml:"invoice_collection"`
	}
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	*s = Subscription(v.subscriptionAlias)
	s.XMLName = v.XMLName
	s.AccountCode = string(v.AccountCode)
	s.InvoiceNumber = int(v.InvoiceNumber)

	if v.InvoiceCollection != nil {
		s.Invoice = v.InvoiceCollection.ChargeInvoice
	}

	return nil
}

// MakeUpdate creates an UpdateSubscription with values that need to be passed
// on update to be retained (meaning nil/zero values will delete that value).
// After calling MakeUpdate you should modify the struct with your updates.
// Once you're ready you can call client.Subscriptions.Update
func (s Subscription) MakeUpdate() UpdateSubscription {
	return UpdateSubscription{
		// NetTerms need to be copied over because on update they default to 0.
		// This ensures the NetTerms don't get overridden.
		NetTerms:           s.NetTerms,
		SubscriptionAddOns: &s.SubscriptionAddOns,
	}
}

type NestedPlan struct {
	Code string `xml:"plan_code,omitempty"`
	Name string `xml:"name,omitempty"`
}

// SubscriptionAddOn are add ons to subscriptions.
// https://docs.com/api/subscriptions/subscription-add-ons
type SubscriptionAddOn struct {
	XMLName           xml.Name `xml:"subscription_add_on"`
	Type              string   `xml:"add_on_type,omitempty"`
	Code              string   `xml:"add_on_code"`
	UnitAmountInCents int      `xml:"unit_amount_in_cents"`
	Quantity          int      `xml:"quantity,omitempty"`
}

// PendingSubscription are updates to the subscription or subscription add ons that
// will be made on the next renewal.
type PendingSubscription struct {
	XMLName            xml.Name            `xml:"pending_subscription"`
	Plan               NestedPlan          `xml:"plan,omitempty"`
	UnitAmountInCents  int                 `xml:"unit_amount_in_cents,omitempty"`
	Quantity           int                 `xml:"quantity,omitempty"` // Quantity of subscriptions
	SubscriptionAddOns []SubscriptionAddOn `xml:"subscription_add_ons>subscription_add_on,omitempty"`
}

// NewSubscription is used to create new subscriptions.
type NewSubscription struct {
	XMLName                 xml.Name             `xml:"subscription"`
	PlanCode                string               `xml:"plan_code"`
	Account                 Account              `xml:"account"`
	SubscriptionAddOns      *[]SubscriptionAddOn `xml:"subscription_add_ons>subscription_add_on,omitempty"`
	CouponCode              string               `xml:"coupon_code,omitempty"`
	UnitAmountInCents       int                  `xml:"unit_amount_in_cents,omitempty"`
	Currency                string               `xml:"currency"`
	Quantity                int                  `xml:"quantity,omitempty"`
	TrialEndsAt             NullTime             `xml:"trial_ends_at,omitempty"`
	StartsAt                NullTime             `xml:"starts_at,omitempty"`
	TotalBillingCycles      int                  `xml:"total_billing_cycles,omitempty"`
	RenewalBillingCycles    NullInt              `xml:"renewal_billing_cycles"`
	FirstRenewalDate        NullTime             `xml:"first_renewal_date,omitempty"`
	CollectionMethod        string               `xml:"collection_method,omitempty"`
	AutoRenew               bool                 `xml:"auto_renew,omitempty"`
	NetTerms                NullInt              `xml:"net_terms,omitempty"`
	PONumber                string               `xml:"po_number,omitempty"`
	Bulk                    bool                 `xml:"bulk,omitempty"`
	TermsAndConditions      string               `xml:"terms_and_conditions,omitempty"`
	CustomerNotes           string               `xml:"customer_notes,omitempty"`
	VATReverseChargeNotes   string               `xml:"vat_reverse_charge_notes,omitempty"`
	BankAccountAuthorizedAt NullTime             `xml:"bank_account_authorized_at,omitempty"`
	CustomFields            *CustomFields        `xml:"custom_fields,omitempty"`
}

// NewSubscriptionResponse is used to unmarshal either the subscription or the transaction.
type NewSubscriptionResponse struct {
	Subscription *Subscription
	Transaction  *Transaction // UnprocessableEntity errors return only the transaction
}

// UpdateSubscription is used to update subscriptions
type UpdateSubscription struct {
	XMLName              xml.Name             `xml:"subscription"`
	Timeframe            string               `xml:"timeframe,omitempty"`
	PlanCode             string               `xml:"plan_code,omitempty"`
	Quantity             int                  `xml:"quantity,omitempty"`
	UnitAmountInCents    int                  `xml:"unit_amount_in_cents,omitempty"`
	RenewalBillingCycles NullInt              `xml:"renewal_billing_cycles"`
	CollectionMethod     string               `xml:"collection_method,omitempty"`
	AutoRenew            bool                 `xml:"auto_renew,omitempty"`
	NetTerms             NullInt              `xml:"net_terms,omitempty"`
	PONumber             string               `xml:"po_number,omitempty"`
	SubscriptionAddOns   *[]SubscriptionAddOn `xml:"subscription_add_ons>subscription_add_on,omitempty"`
}

// SubscriptionNotes is used to update a subscription's notes.
type SubscriptionNotes struct {
	XMLName               xml.Name `xml:"subscription"`
	TermsAndConditions    string   `xml:"terms_and_conditions,omitempty"`
	CustomerNotes         string   `xml:"customer_notes,omitempty"`
	VATReverseChargeNotes string   `xml:"vat_reverse_charge_notes,omitempty"`
}

// CustomFields represents custom key value pairs.
// Note that custom fields must be enabled on your Recurly site and must be added in
// the dashboard before they can be used.
type CustomFields map[string]string

// UnmarshalXML unmarshals custom_fields.
func (c *CustomFields) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v struct {
		XMLName xml.Name `xml:"custom_fields"`
		Fields  []struct {
			Name  string `xml:"name"`
			Value string `xml:"value"`
		} `xml:"custom_field"`
	}
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	m := make(map[string]string, len(v.Fields))
	for _, f := range v.Fields {
		m[f.Name] = f.Value
	}
	*c = m

	return nil
}

//MarshalXML marshals custom_fields.
func (c CustomFields) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(c) == 0 {
		return nil
	}

	type xmlMapEntry struct {
		XMLName struct{} `xml:"custom_field"`
		Name    string   `xml:"name"`
		Value   string   `xml:"value"`
	}

	e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "custom_fields"}})

	// Ensure key field order, otherwise causes rendered xml can have order difference
	var keys []string
	for k := range c {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		n := &xmlMapEntry{
			Name:  k,
			Value: c[k],
		}
		e.Encode(n)
	}
	e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "custom_fields"}})
	return nil
}
