package recurly

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// SubscriptionsService manages the interactions for subscriptions.
type SubscriptionsService interface {
	// List returns a pager to paginate subscription. PagerOptions are used to
	// optionally filter the results.
	//
	// https://dev.recurly.com/docs/list-subscriptions
	List(opts *PagerOptions) Pager

	// ListAccount returns a pager to paginate subscriptions for an account.
	// PagerOptions are used to optionally filter the results.
	//
	// https://dev.recurly.com/docs/list-accounts-subscriptions
	ListAccount(accountCode string, opts *PagerOptions) Pager

	// Get retrieves a subscription. If the subscription does not exist,
	// a nil subscription and nil error are returned.
	//
	// https://dev.recurly.com/docs/lookup-subscription-details
	Get(ctx context.Context, uuid string) (*Subscription, error)

	// Create creates a subscription. You can optionally include subscription
	// add ons. See Recurly's documentation for specfics.
	//
	// https://dev.recurly.com/docs/create-subscription
	// https://dev.recurly.com/docs/subscription-add-ons
	Create(ctx context.Context, sub NewSubscription) (*Subscription, error)

	// Preview returns a preview for a new subscription applied to an account.
	//
	// https://dev.recurly.com/docs/preview-subscription
	Preview(ctx context.Context, sub NewSubscription) (*Subscription, error)

	// Update updates a subscription that takes place immediately or at renewal
	// based on sub.Timeframe. You can optionally send subscription add ons.
	// If the subscription has add ons and you omit, the subscription will
	// be updated and all add ons will be removed.
	//
	// https://dev.recurly.com/docs/update-subscription
	// https://dev.recurly.com/docs/update-subscription-with-add-ons
	Update(ctx context.Context, uuid string, sub UpdateSubscription) (*Subscription, error)

	// UpdateNotes updates a subscription's invoice notes before the next renewal.
	// Updating notes will not trigger the renewal.
	//
	// https://dev.recurly.com/docs/update-subscription-notes
	UpdateNotes(ctx context.Context, uuid string, n SubscriptionNotes) (*Subscription, error)

	// PreviewChange previews a subscription change applied to an account without
	// committing a subscription change or posting an invoice.
	//
	// https://dev.recurly.com/docs/preview-subscription-change
	PreviewChange(ctx context.Context, uuid string, sub UpdateSubscription) (*Subscription, error)

	// Cancel cancels a subscription so it remains active and then expires at
	// the end of the current bill cycle.
	//
	// https://dev.recurly.com/docs/cancel-subscription
	Cancel(ctx context.Context, uuid string) (*Subscription, error)

	// Reactive reactivates a canceled subscription so it renews at the end
	// of the current bill cycle.
	//
	// https://dev.recurly.com/docs/reactivate-canceled-subscription
	Reactivate(ctx context.Context, uuid string) (*Subscription, error)

	// Terminate terminates a subscription and refunds according to refundType.
	// Valid values for refundType: 'partial', 'full', or 'none'.
	// See Recurly's documentation for more details.
	//
	// https://dev.recurly.com/docs/terminate-subscription
	Terminate(ctx context.Context, uuid string, refundType string) (*Subscription, error)

	// Pause schedules a pause or updates remaining pause cycles for a subscription.
	//
	// https://dev.recurly.com/docs/pause-subscription
	Pause(ctx context.Context, uuid string, cycles int) (*Subscription, error)

	// Postpone changes the next bill date (for an active subscription) or
	// changes when the trial expires (for subscriptions in trial period).
	// See Recurly's documentation for details.
	//
	// https://dev.recurly.com/docs/postpone-subscription
	Postpone(ctx context.Context, uuid string, dt time.Time, bulk bool) (*Subscription, error)

	// Resume reactivates a paused subscription, starting a new billing cycle.
	//
	// https://dev.recurly.com/docs/resume-subscription
	Resume(ctx context.Context, uuid string) (*Subscription, error)

	// Immediately converts a trial subscription to paid
	//
	// https://dev.recurly.com/docs/convert-trial
	ConvertTrial(ctx context.Context, uuid string) (*Subscription, error)
}

// Subscription state constants.
// https://docs.recurly.com/docs/subscriptions
const (
	SubscriptionStateActive   = "active"
	SubscriptionStateCanceled = "canceled"
	SubscriptionStateExpired  = "expired"
	SubscriptionStateFuture   = "future"
	SubscriptionStateInTrial  = "in_trial"
	SubscriptionStateLive     = "live"
	SubscriptionStatePastDue  = "past_due"
	SubscriptionStatePaused   = "paused"
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
	PausedAt               NullTime             `xml:"paused_at,omitempty"`
	ResumeAt               NullTime             `xml:"resume_at,omitempty"`
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
	InvoiceCollection      *InvoiceCollection   `xml:"invoice_collection,omitempty"`
	RemainingPauseCycles   int                  `xml:"remaining_pause_cycles,omitempty"`
	CollectionMethod       string               `xml:"collection_method"`
	CustomerNotes          string               `xml:"customer_notes,omitempty"`
	AutoRenew              bool                 `xml:"auto_renew,omitempty"`
	RenewalBillingCycles   NullInt              `xml:"renewal_billing_cycles,omitempty"`
	RemainingBillingCycles NullInt              `xml:"remaining_billing_cycles,omitempty"`
	GatewayCode            string               `xml:"gateway_code,omitempty"`
	CustomFields           *CustomFields        `xml:"custom_fields,omitempty"`
}

// NewSubscription is used to create new subscriptions.
type NewSubscription struct {
	XMLName                 xml.Name             `xml:"subscription"`
	PlanCode                string               `xml:"plan_code"`
	Account                 Account              `xml:"account"`
	SubscriptionAddOns      *[]SubscriptionAddOn `xml:"subscription_add_ons>subscription_add_on,omitempty"`
	CouponCode              string               `xml:"coupon_code,omitempty"`
	UnitAmountInCents       NullInt              `xml:"unit_amount_in_cents,omitempty"`
	Currency                string               `xml:"currency"`
	Quantity                int                  `xml:"quantity,omitempty"`
	TrialEndsAt             NullTime             `xml:"trial_ends_at,omitempty"`
	StartsAt                NullTime             `xml:"starts_at,omitempty"`
	TotalBillingCycles      int                  `xml:"total_billing_cycles,omitempty"`
	RenewalBillingCycles    NullInt              `xml:"renewal_billing_cycles"`
	NextBillDate            NullTime             `xml:"next_bill_date,omitempty"`
	CollectionMethod        string               `xml:"collection_method,omitempty"`
	AutoRenew               bool                 `xml:"auto_renew,omitempty"`
	NetTerms                NullInt              `xml:"net_terms,omitempty"`
	PONumber                string               `xml:"po_number,omitempty"`
	Bulk                    bool                 `xml:"bulk,omitempty"`
	TermsAndConditions      string               `xml:"terms_and_conditions,omitempty"`
	CustomerNotes           string               `xml:"customer_notes,omitempty"`
	VATReverseChargeNotes   string               `xml:"vat_reverse_charge_notes,omitempty"`
	BankAccountAuthorizedAt NullTime             `xml:"bank_account_authorized_at,omitempty"`
	RevenueScheduleType     string               `xml:"revenue_schedule_type,omitempty"`
	ShippingAddress         *ShippingAddress     `xml:"shipping_address,omitempty"`
	ShippingAddressID       int64                `xml:"shipping_address_id,omitempty"`
	ImportedTrial           NullBool             `xml:"imported_trial,omitempty"`
	GatewayCode             string               `xml:"gateway_code,omitempty"`
	ShippingMethodCode      string               `xml:"shipping_method_code,omitempty"`
	ShippingAmountInCents   NullInt              `xml:"shipping_amount_in_cents,omitempty"`
	CustomFields            *CustomFields        `xml:"custom_fields,omitempty"`
	TransactionType         string               `xml:"transaction_type,omitempty"` // Optional transaction type. Currently accepts "moto"
}

// UnmarshalXML unmarshals transactions and handles intermediary state during unmarshaling
// for types like href.
func (s *Subscription) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type subscriptionAlias Subscription
	var v struct {
		subscriptionAlias
		XMLName       xml.Name `xml:"subscription"`
		AccountCode   href     `xml:"account"`
		InvoiceNumber href     `xml:"invoice"`
	}
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}

	*s = Subscription(v.subscriptionAlias)
	s.XMLName = v.XMLName
	s.AccountCode = v.AccountCode.LastPartOfPath()
	s.InvoiceNumber, _ = strconv.Atoi(v.InvoiceNumber.LastPartOfPath())
	return nil
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
	UnitAmountInCents NullInt  `xml:"unit_amount_in_cents,omitempty"`
	Quantity          int      `xml:"quantity,omitempty"`
	AddOnSource       string   `xml:"add_on_source,omitempty"`
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

// UpdateSubscription is used to update subscriptions
type UpdateSubscription struct {
	XMLName                xml.Name             `xml:"subscription"`
	Timeframe              string               `xml:"timeframe,omitempty"`
	PlanCode               string               `xml:"plan_code,omitempty"`
	Quantity               int                  `xml:"quantity,omitempty"`
	UnitAmountInCents      NullInt              `xml:"unit_amount_in_cents,omitempty"`
	CollectionMethod       string               `xml:"collection_method,omitempty"`
	NetTerms               NullInt              `xml:"net_terms,omitempty"`
	PONumber               string               `xml:"po_number,omitempty"`
	SubscriptionAddOns     *[]SubscriptionAddOn `xml:"subscription_add_ons>subscription_add_on,omitempty"`
	CouponCode             string               `xml:"coupon_code,omitempty"`
	RevenueScheduleType    string               `xml:"revenue_schedule_type,omitempty"`
	RemainingBillingCycles NullInt              `xml:"remaining_billing_cycles,omitempty"`
	ImportedTrial          NullBool             `xml:"imported_trial,omitempty"`
	RenewalBillingCycles   NullInt              `xml:"renewal_billing_cycles,omitempty"`
	AutoRenew              NullBool             `xml:"auto_renew,omitempty"`
	CustomFields           *CustomFields        `xml:"custom_fields,omitempty"`
	BillingInfo            *Billing             `xml:"billing_info,omitempty"`
	TransactionType        string               `xml:"transaction_type,omitempty"` // Optional transaction type. Currently accepts "moto"
}

// SubscriptionNotes is used to update a subscription's notes.
type SubscriptionNotes struct {
	XMLName               xml.Name      `xml:"subscription"`
	TermsAndConditions    string        `xml:"terms_and_conditions,omitempty"`
	CustomerNotes         string        `xml:"customer_notes,omitempty"`
	VATReverseChargeNotes string        `xml:"vat_reverse_charge_notes,omitempty"`
	GatewayCode           string        `xml:"gateway_code"`
	CustomFields          *CustomFields `xml:"custom_fields,omitempty"`
}

// CustomFields represents custom key value pairs.
// Note that custom fields must be enabled on your Recurly site and must be added in
// the dashboard before they can be used.
type CustomFields map[string]string

// UnmarshalXML unmarshals custom_fields.
func (c *CustomFields) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v struct {
		XMLName xml.Name      `xml:"custom_fields"`
		Fields  []customField `xml:"custom_field"`
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

	e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "custom_fields"}})
	e.Encode(c.xmlFields())
	e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "custom_fields"}})
	return nil
}

type customField struct {
	XMLName struct{} `xml:"custom_field"`
	Name    string   `xml:"name"`
	Value   string   `xml:"value"`
}

// xmlFields returns []customField from CustomFields. Results sorted for testing.
func (c CustomFields) xmlFields() []customField {
	var i int
	fields := make([]customField, len(c))
	for k := range c {
		fields[i] = customField{
			Name:  k,
			Value: c[k],
		}
		i++
	}
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})
	return fields
}

var _ SubscriptionsService = &subscriptionsImpl{}

// subscriptionsImpl implements SubscriptionsService.
type subscriptionsImpl serviceImpl

func (s *subscriptionsImpl) List(opts *PagerOptions) Pager {
	return s.client.newPager("GET", "/subscriptions", opts)
}

func (s *subscriptionsImpl) ListAccount(accountCode string, opts *PagerOptions) Pager {
	path := fmt.Sprintf("/accounts/%s/subscriptions", accountCode)
	return s.client.newPager("GET", path, opts)
}

func (s *subscriptionsImpl) Get(ctx context.Context, uuid string) (*Subscription, error) {
	path := fmt.Sprintf("/subscriptions/%s", sanitizeUUID(uuid))
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var dst Subscription
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		if e, ok := err.(*ClientError); ok && e.Response.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &dst, nil
}

func (s *subscriptionsImpl) Create(ctx context.Context, sub NewSubscription) (*Subscription, error) {
	req, err := s.client.newRequest("POST", "/subscriptions", sub)
	if err != nil {
		return nil, err
	}

	var dst Subscription
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, err
}

func (s *subscriptionsImpl) Preview(ctx context.Context, sub NewSubscription) (*Subscription, error) {
	req, err := s.client.newRequest("POST", "/subscriptions/preview", sub)
	if err != nil {
		return nil, err
	}

	var dst Subscription
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, err
}

func (s *subscriptionsImpl) Update(ctx context.Context, uuid string, sub UpdateSubscription) (*Subscription, error) {
	path := fmt.Sprintf("/subscriptions/%s", sanitizeUUID(uuid))
	req, err := s.client.newRequest("PUT", path, sub)
	if err != nil {
		return nil, err
	}

	var dst Subscription
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, err
}

func (s *subscriptionsImpl) UpdateNotes(ctx context.Context, uuid string, n SubscriptionNotes) (*Subscription, error) {
	path := fmt.Sprintf("/subscriptions/%s/notes", sanitizeUUID(uuid))
	req, err := s.client.newRequest("PUT", path, n)
	if err != nil {
		return nil, err
	}

	var dst Subscription
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, err
}

func (s *subscriptionsImpl) PreviewChange(ctx context.Context, uuid string, sub UpdateSubscription) (*Subscription, error) {
	path := fmt.Sprintf("/subscriptions/%s/preview", sanitizeUUID(uuid))
	req, err := s.client.newRequest("POST", path, sub)
	if err != nil {
		return nil, err
	}

	var dst Subscription
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, err
}

func (s *subscriptionsImpl) Cancel(ctx context.Context, uuid string) (*Subscription, error) {
	path := fmt.Sprintf("/subscriptions/%s/cancel", sanitizeUUID(uuid))
	req, err := s.client.newRequest("PUT", path, nil)
	if err != nil {
		return nil, err
	}

	var dst Subscription
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, err
}

func (s *subscriptionsImpl) Reactivate(ctx context.Context, uuid string) (*Subscription, error) {
	path := fmt.Sprintf("/subscriptions/%s/reactivate", sanitizeUUID(uuid))
	req, err := s.client.newRequest("PUT", path, nil)
	if err != nil {
		return nil, err
	}

	var dst Subscription
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, err
}

func (s *subscriptionsImpl) Terminate(ctx context.Context, uuid string, refundType string) (*Subscription, error) {
	path := fmt.Sprintf("/subscriptions/%s/terminate", sanitizeUUID(uuid))
	req, err := s.client.newQueryRequest("PUT", path, query{
		"refund": refundType,
	}, nil)
	if err != nil {
		return nil, err
	}

	var dst Subscription
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *subscriptionsImpl) Pause(ctx context.Context, uuid string, cycles int) (*Subscription, error) {
	path := fmt.Sprintf("/subscriptions/%s/pause", sanitizeUUID(uuid))
	req, err := s.client.newRequest("PUT", path, struct {
		XMLName              xml.Name `xml:"subscription"`
		RemainingPauseCycles int      `xml:"remaining_pause_cycles"`
	}{
		RemainingPauseCycles: cycles,
	})
	if err != nil {
		return nil, err
	}

	var dst Subscription
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, err
}

func (s *subscriptionsImpl) Postpone(ctx context.Context, uuid string, dt time.Time, bulk bool) (*Subscription, error) {
	path := fmt.Sprintf("/subscriptions/%s/postpone", sanitizeUUID(uuid))
	req, err := s.client.newQueryRequest("PUT", path, query{
		"bulk":              bulk,
		"next_renewal_date": dt,
	}, nil)
	if err != nil {
		return nil, err
	}

	var dst Subscription
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, err
}

func (s *subscriptionsImpl) Resume(ctx context.Context, uuid string) (*Subscription, error) {
	path := fmt.Sprintf("/subscriptions/%s/resume", sanitizeUUID(uuid))
	req, err := s.client.newRequest("PUT", path, nil)
	if err != nil {
		return nil, err
	}

	var dst Subscription
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, err
}

func (s *subscriptionsImpl) ConvertTrial(ctx context.Context, uuid string) (*Subscription, error) {
	path := fmt.Sprintf("/subscriptions/%s/convert_trial", sanitizeUUID(uuid))
	req, err := s.client.newRequest("PUT", path, nil)
	if err != nil {
		return nil, err
	}

	var dst Subscription
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, err
}

// sanitizeUUID returns the uuid without dashes.
func sanitizeUUID(uuid string) string {
	if !strings.Contains(uuid, "-") {
		return uuid
	}
	return strings.TrimSpace(strings.Replace(uuid, "-", "", -1))
}
