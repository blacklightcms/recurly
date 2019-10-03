package recurly

import (
	"context"
	"encoding/xml"
	"fmt"
)

// PurchasesService manages the interactions for a purchase
// involving at least one adjustment or one subscription.
type PurchasesService interface {
	// Create a purchase. See Recurly's documentation for more details.
	//
	// https://dev.recurly.com/docs/create-purchase
	Create(ctx context.Context, p Purchase) (*InvoiceCollection, error)

	// Preview a purchase. See Recurly's documentation for more details.
	//
	// https://dev.recurly.com/docs/preview-purchase
	Preview(ctx context.Context, p Purchase) (*InvoiceCollection, error)

	// Authorize creates a pending purchase that can be activated at a later
	// time once payment has been completed on an external source (e.g. Adyen's
	// Hosted Payment Pages).
	//
	// p.Account.Email and p.Account.Billing.ExternalHPPType appear to be required.
	//
	// https://dev.recurly.com/docs/authorize-purchase
	Authorize(ctx context.Context, p Purchase) (*Purchase, error)

	// Pending is used for Adyen HPP transaction requests. This runs the validations
	// but not the transactions. See Recurly's documentation for more info.
	//
	// https://dev.recurly.com/docs/pending-purchase
	Pending(ctx context.Context, p Purchase) (*Purchase, error)

	// Capture an open Authorization request. See Recurly's documentation
	// for details.
	//
	// https://dev.recurly.com/docs/capture-purchase
	Capture(ctx context.Context, transactionUUID string) (*InvoiceCollection, error)

	// Cancel an open Authorization request. See Recurly's documentation
	// for details.
	//
	// https://dev.recurly.com/docs/cancel-purchase
	Cancel(ctx context.Context, transactionUUID string) (*InvoiceCollection, error)
}

// Purchase represents an individual checkout holding at least one
// subscription OR one adjustment.
// NOTE: Adjustments cannot contain a Currency field. Use Purchase.Currency instead.
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
	ShippingFees          *[]ShippingFee         `xml:"shipping_fees>shipping_fee,omitempty"`
	TransactionType       string                 `xml:"transaction_type,omitempty"` // Create only
}

// PurchaseSubscription represents a subscription to purchase some new subscription.
// This is different from the Subscription struct in that only fields allowed to
// be used with the purchases API are available.
type PurchaseSubscription struct {
	XMLName               xml.Name             `xml:"subscription"`
	PlanCode              string               `xml:"plan_code"`
	SubscriptionAddOns    *[]SubscriptionAddOn `xml:"subscription_add_ons>subscription_add_on,omitempty"`
	UnitAmountInCents     NullInt              `xml:"unit_amount_in_cents,omitempty"`
	Quantity              int                  `xml:"quantity,omitempty"`
	TrialEndsAt           NullTime             `xml:"trial_ends_at,omitempty"`
	StartsAt              NullTime             `xml:"starts_at,omitempty"`
	TotalBillingCycles    int                  `xml:"total_billing_cycles,omitempty"`
	RenewalBillingCycles  NullInt              `xml:"renewal_billing_cycles,omitempty"`
	NextBillDate          NullTime             `xml:"next_bill_date,omitempty"`
	AutoRenew             bool                 `xml:"auto_renew,omitempty"`
	CustomFields          *CustomFields        `xml:"custom_fields,omitempty"`
	ShippingAddress       *ShippingAddress     `xml:"shipping_address,omitempty"`
	ShippingAddressID     int64                `xml:"shipping_address_id,omitempty"`
	ShippingMethodCode    string               `xml:"shipping_method_code,omitempty"`
	ShippingAmountInCents NullInt              `xml:"shipping_amount_in_cents,omitempty"`
}

// ShippingFee holds shipping fees for a purchase.
type ShippingFee struct {
	XMLName               xml.Name `xml:"shipping_fee"`
	ShippingMethodCode    string   `xml:"shipping_method_code,omitempty"`
	ShippingAmountInCents NullInt  `xml:"shipping_amount_in_cents,omitempty"`
}

var _ PurchasesService = &purchasesImpl{}

// purchasesImpl implements PurchasesService.
type purchasesImpl serviceImpl

func (s *purchasesImpl) Create(ctx context.Context, p Purchase) (*InvoiceCollection, error) {
	req, err := s.client.newRequest("POST", "/purchases", p)
	if err != nil {
		return nil, err
	}

	var dst InvoiceCollection
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *purchasesImpl) Preview(ctx context.Context, p Purchase) (*InvoiceCollection, error) {
	req, err := s.client.newRequest("POST", "/purchases/preview", p)
	if err != nil {
		return nil, err
	}

	var dst InvoiceCollection
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *purchasesImpl) Authorize(ctx context.Context, p Purchase) (*Purchase, error) {
	req, err := s.client.newRequest("POST", "/purchases/authorize", p)
	if err != nil {
		return nil, err
	}

	var dst Purchase
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *purchasesImpl) Pending(ctx context.Context, p Purchase) (*Purchase, error) {
	req, err := s.client.newRequest("POST", "/purchases/pending", p)
	if err != nil {
		return nil, err
	}

	var dst Purchase
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *purchasesImpl) Capture(ctx context.Context, transactionUUID string) (*InvoiceCollection, error) {
	path := fmt.Sprintf("/purchases/%s/capture", sanitizeUUID(transactionUUID))
	req, err := s.client.newRequest("POST", path, nil)
	if err != nil {
		return nil, err
	}

	var dst InvoiceCollection
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *purchasesImpl) Cancel(ctx context.Context, transactionUUID string) (*InvoiceCollection, error) {
	path := fmt.Sprintf("/purchases/%s/cancel", sanitizeUUID(transactionUUID))
	req, err := s.client.newRequest("POST", path, nil)
	if err != nil {
		return nil, err
	}

	var dst InvoiceCollection
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, nil
}
