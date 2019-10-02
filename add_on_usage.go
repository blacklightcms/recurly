package recurly

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
)

// AddOnUsageService manages the interactions for add-on usages.
type AddOnUsageService interface {
	// List returns a pager to paginate usages for add-on in subscription. PagerOptions are used to
	// optionally filter the results.
	//
	// https://dev.recurly.com/docs/list-add-ons-usage
	List(subUUID, addOnCode string, opts *PagerOptions) Pager

	// Create creates a new usage for add-on in subscription.
	//
	// https://dev.recurly.com/docs/log-usage
	Create(ctx context.Context, subUUID, addOnCode string, usage AddOnUsage) (*AddOnUsage, error)

	// Get retrieves an usage. If the usage does not exist,
	// a nil usage and nil error are returned.
	//
	// https://dev.recurly.com/docs/lookup-usage-record
	Get(ctx context.Context, subUUID, addOnCode, usageId string) (*AddOnUsage, error)

	// Update updates the usage information. Once usage is billed, only MerchantTag can be updated.
	//
	// https://dev.recurly.com/docs/update-usage
	Update(ctx context.Context, subUUID, addOnCode, usageId string, usage AddOnUsage) (*AddOnUsage, error)

	// Delete removes an usage from subscription add-on. If usage is billed, it can't be removed
	//
	// https://dev.recurly.com/docs/delete-a-usage-record
	Delete(ctx context.Context, subUUID, addOnCode, usageId string) error
}

// Usage is a billable event or group of events recorded on a purchased usage-based add-on and billed in arrears each billing cycle.
//
// https://dev.recurly.com/docs/usage-record-object
type AddOnUsage struct {
	XMLName            xml.Name `xml:"usage"`
	Id                 int      `xml:"id,omitempty"`
	Amount             int      `xml:"amount,omitempty"`
	MerchantTag        string   `xml:"merchant_tag,omitempty"`
	RecordingTimestamp NullTime `xml:"recording_timestamp,omitempty"`
	UsageTimestamp     NullTime `xml:"usage_timestamp,omitempty"`
	CreatedAt          NullTime `xml:"created_at,omitempty"`
	UpdatedAt          NullTime `xml:"updated_at,omitempty"`
	BilledAt           NullTime `xml:"billed_at,omitempty"`
	UsageType          string   `xml:"usage_type,omitempty"`
	UnitAmountInCents  int      `xml:"unit_amount_in_cents,omitempty"`
	UsagePercentage    NullFloat  `xml:"usage_percentage,omitempty"`
}

var _ AddOnUsageService = &addOnUsageServiceImpl{}

type addOnUsageServiceImpl serviceImpl

func (s *addOnUsageServiceImpl) List(subUUID, addOnCode string, opts *PagerOptions) Pager {

	path := fmt.Sprintf("/subscriptions/%s/add_ons/%s/usage", subUUID, addOnCode)
	fmt.Printf("\n huh, path=[%s]\n", path)
	return s.client.newPager("GET", path, opts)
}

func (s *addOnUsageServiceImpl) Get(ctx context.Context, subUUID, addOnCode, usageId string) (*AddOnUsage, error) {
	path := fmt.Sprintf("/subscriptions/%s/add_ons/%s/usage/%s", subUUID, addOnCode, usageId)
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var dst AddOnUsage
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		if e, ok := err.(*ClientError); ok && e.Response.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &dst, nil
}

func (s *addOnUsageServiceImpl) Create(ctx context.Context, subUUID, addOnCode string, usage AddOnUsage) (*AddOnUsage, error) {

	path := fmt.Sprintf("/subscriptions/%s/add_ons/%s/usage", subUUID, addOnCode)
	req, err := s.client.newRequest("POST", path, usage)
	if err != nil {
		return nil, err
	}

	var dst AddOnUsage
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *addOnUsageServiceImpl) Update(ctx context.Context, subUUID, addOnCode, usageId string, usage AddOnUsage) (*AddOnUsage, error) {

	path := fmt.Sprintf("/subscriptions/%s/add_ons/%s/usage/%s", subUUID, addOnCode, usageId)
	req, err := s.client.newRequest("PUT", path, usage)
	if err != nil {
		return nil, err
	}

	var dst AddOnUsage
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *addOnUsageServiceImpl) Delete(ctx context.Context, subUUID, addOnCode, usageId string) error {

	path := fmt.Sprintf("/subscriptions/%s/add_ons/%s/usage/%s", subUUID, addOnCode, usageId)
	req, err := s.client.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(ctx, req, nil)
	return err
}
