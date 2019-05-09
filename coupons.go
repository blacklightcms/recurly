package recurly

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// CouponsService manages the interactions for coupons.
type CouponsService interface {
	// List returns a pager to paginate coupons. PagerOptions are used to optionally
	// filter the results.
	// https://dev.recurly.com/docs/list-active-coupons
	List(opts *PagerOptions) *CouponsPager

	// Get retrieves a coupon. If the coupon does not exist,
	// a nil coupon and nil error is returned.
	// https://dev.recurly.com/docs/lookup-a-coupon
	Get(ctx context.Context, code string) (*Coupon, error)

	// Create creates a new coupon. Please note: coupons cannot be updated
	// after being created.
	// https://dev.recurly.com/docs/create-coupon
	Create(ctx context.Context, c Coupon) (*Coupon, error)

	// Update edits a redeemable coupon to extend redemption rules.
	// Only redeemable coupons can be edited an only certain fields can
	// be sent. See documentation for the list of valid fields.
	// Any fields provided outside of valid fields will be automatically
	// excluded before sending to Recurly.
	//
	// https://dev.recurly.com/docs/edit-coupon
	Update(ctx context.Context, code string, c Coupon) (*Coupon, error)

	// Restore restores an expired coupon so it can be redeemed again.
	// You can change editable fields in this call. See documentation for the
	// list of valid fields. Any fields provided outside of valid fields
	// will be automatically excluded before sending to Recurly.
	//
	// https://dev.recurly.com/docs/restore-coupon
	Restore(ctx context.Context, code string, c Coupon) (*Coupon, error)

	// Delete expires a coupon so customers can no longer redeem it.
	// https://dev.recurly.com/docs/deactivate-coupon
	Delete(ctx context.Context, code string) error

	// Generate creates unique codes for a bulk coupon. A bulk coupon can
	// have up to 100,000 unique codes total. The generate endpoint allows
	// up to 200 unique codes at a time. The endpoint can be called
	// multiple times to create the number of coupon codes you need.
	//
	// The response will return a CouponsPager to view the unique codes.
	// https://dev.recurly.com/docs/generate-unique-codes
	Generate(ctx context.Context, code string, n int) (*CouponsPager, error)
}

// Coupon represents an individual coupon on your site.
type Coupon struct {
	XMLName                  xml.Name    `xml:"coupon"`
	ID                       int         `xml:"id,omitempty"`
	Code                     string      `xml:"coupon_code"`
	Type                     string      `xml:"coupon_type,omitempty"`
	Name                     string      `xml:"name"`
	RedemptionResource       string      `xml:"redemption_resource,omitempty"`
	State                    string      `xml:"state,omitempty"`
	AppliesToAllPlans        bool        `xml:"applies_to_all_plans,omitempty"`
	Duration                 string      `xml:"duration,omitempty"`
	DiscountType             string      `xml:"discount_type"`
	AppliesToNonPlanCharges  bool        `xml:"applies_to_non_plan_charges,omitempty"`
	Description              string      `xml:"description,omitempty"`
	InvoiceDescription       string      `xml:"invoice_description,omitempty"`
	DiscountPercent          NullInt     `xml:"discount_percent,omitempty"`
	DiscountInCents          *UnitAmount `xml:"discount_in_cents,omitempty"`
	RedeemByDate             NullTime    `xml:"redeem_by_date,omitempty"`
	MaxRedemptions           NullInt     `xml:"max_redemptions,omitempty"`
	CreatedAt                NullTime    `xml:"created_at,omitempty"`
	UpdatedAt                NullTime    `xml:"updated_at,omitempty"`
	DeletedAt                NullTime    `xml:"deleted_at,omitempty"`
	TemporalUnit             string      `xml:"temporal_unit,omitempty"`
	TemporalAmount           NullInt     `xml:"temporal_amount,omitempty"`
	MaxRedemptionsPerAccount NullInt     `xml:"max_redemptions_per_account,omitempty"`
	UniqueCodeTemplate       string      `xml:"unique_code_template,omitempty"`
	UniqueCouponCodeCount    NullInt     `xml:"unique_coupon_codes_count,omitempty"`
	PlanCodes                []string    `xml:"plan_codes>plan_code,omitempty"`
	// SingleUse deprecated. Please use duration instead.
}

// returns the editable fields used for the Redeem() and Restore() endpoints.
func (c Coupon) editableFields() couponEditableFields {
	return couponEditableFields{
		Name:                     c.Name,
		Description:              c.Description,
		InvoiceDescription:       c.InvoiceDescription,
		RedeemByDate:             c.RedeemByDate,
		MaxRedemptions:           c.MaxRedemptions,
		MaxRedemptionsPerAccount: c.MaxRedemptionsPerAccount,
	}
}

type couponEditableFields struct {
	XMLName                  xml.Name `xml:"coupon"`
	Name                     string   `xml:"name,omitempty"`
	Description              string   `xml:"description,omitempty"`
	InvoiceDescription       string   `xml:"invoice_description,omitempty"`
	RedeemByDate             NullTime `xml:"redeem_by_date,omitempty"`
	MaxRedemptions           NullInt  `xml:"max_redemptions"`
	MaxRedemptionsPerAccount NullInt  `xml:"max_redemptions_per_account"`
}

// CouponsPager paginates coupons.
type CouponsPager struct {
	*pager
}

// Fetch fetches the next set of results.
func (p *CouponsPager) Fetch(ctx context.Context) ([]Coupon, error) {
	var dst struct {
		// This pager needs to process both 'coupons' and 'unique_coupon_codes'
		// as the top-level XML tags. We intentionally don't set the xml
		// tag name so it works with both.
		// See the Generate() method for specifics.
		XMLName xml.Name
		Coupons []Coupon `xml:"coupon"`
	}
	if err := p.fetch(ctx, &dst); err != nil {
		return nil, err
	}
	return dst.Coupons, nil
}

// FetchAll paginates all records and returns a cumulative list.
func (p *CouponsPager) FetchAll(ctx context.Context) ([]Coupon, error) {
	p.setMaxPerPage()

	var all []Coupon
	for p.Next() {
		v, err := p.Fetch(ctx)
		if err != nil {
			return nil, err
		}
		all = append(all, v...)
	}
	return all, nil
}

var _ CouponsService = &couponsImpl{}

// couponsImpl implements CouponsService.
type couponsImpl serviceImpl

func (s *couponsImpl) List(opts *PagerOptions) *CouponsPager {
	return &CouponsPager{
		pager: s.client.newPager("GET", "/coupons", opts),
	}
}

func (s *couponsImpl) Get(ctx context.Context, code string) (*Coupon, error) {
	path := fmt.Sprintf("/coupons/%s", code)
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var dst Coupon
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		if e, ok := err.(*ClientError); ok && e.Response.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &dst, nil
}

func (s *couponsImpl) Create(ctx context.Context, c Coupon) (*Coupon, error) {
	req, err := s.client.newRequest("POST", "/coupons", c)
	if err != nil {
		return nil, err
	}

	var dst Coupon
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *couponsImpl) Update(ctx context.Context, code string, c Coupon) (*Coupon, error) {
	path := fmt.Sprintf("/coupons/%s", code)
	req, err := s.client.newRequest("PUT", path, c.editableFields())
	if err != nil {
		return nil, err
	}

	var dst Coupon
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *couponsImpl) Restore(ctx context.Context, code string, c Coupon) (*Coupon, error) {
	path := fmt.Sprintf("/coupons/%s/restore", code)
	req, err := s.client.newRequest("PUT", path, c.editableFields())
	if err != nil {
		return nil, err
	}

	var dst Coupon
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *couponsImpl) Delete(ctx context.Context, code string) error {
	path := fmt.Sprintf("/coupons/%s", code)
	req, err := s.client.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(ctx, req, nil)
	return err
}

func (s *couponsImpl) Generate(ctx context.Context, code string, n int) (*CouponsPager, error) {
	path := fmt.Sprintf("/coupons/%s/generate", code)
	req, err := s.client.newRequest("POST", path, struct {
		XMLName             xml.Name `xml:"coupon"`
		NumberOfUniqueCodes int      `xml:"number_of_unique_codes"`
	}{
		NumberOfUniqueCodes: n,
	})
	if err != nil {
		return nil, err
	}

	resp, err := s.client.do(ctx, req, nil)
	if err != nil {
		return nil, err
	}

	// The location header holds the path to the generated coupons, but
	// the /v2 prefix needs to be stripped.
	u, err := url.Parse(resp.Header.Get("Location"))
	if err != nil {
		return nil, err
	}
	u.Path = strings.TrimPrefix(u.Path, "/v2")

	// Setup pager and attach params and cursor.
	pager := s.client.newPager("GET", u.Path, nil)
	pager.opts.PerPage, _ = strconv.Atoi(u.Query().Get("per_page"))
	pager.cursor = u.Query().Get("cursor")

	return &CouponsPager{pager: pager}, nil
}
