package recurly

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// pager paginates API calls.
type pager struct {
	client *Client

	method string
	path   string
	opts   *PagerOptions

	count  *int
	cursor string

	expectResults bool
}

// returns a new pager and initializes params if nil. It ensures no cursor
// is set.
func (c *Client) newPager(method, path string, opts *PagerOptions) *pager {
	if opts == nil {
		opts = &PagerOptions{}
	}
	return &pager{
		client: c,
		method: method,
		path:   path,
		opts:   opts,

		expectResults: true,
	}
}

// Count returns a total count of the request. Calling this function
// more than once will return the value from the first call.
func (p *pager) Count(ctx context.Context) (int, error) {
	if p.count != nil {
		return *p.count, nil
	}

	req, err := p.client.newPagerRequest("HEAD", p.path, p.opts, nil)
	if err != nil {
		return 0, err
	}

	resp, err := p.client.do(ctx, req, nil)
	if err != nil {
		return 0, err
	}

	if count := resp.Header.Get("X-Records"); count == "" {
		return 0, nil
	} else if i, err := strconv.Atoi(count); err != nil {
		return 0, err
	} else {
		p.count = &i
		return i, nil
	}
}

// Next prepares returns true if there is a next result expected, or false if
// there is no next result.
func (p *pager) Next() bool {
	return p.expectResults
}

// fetch retrieves the results and populates dst, setting the next
// cursor.
func (p *pager) fetch(ctx context.Context, dst interface{}) error {
	select {
	default:
	case <-ctx.Done():
		return ctx.Err()
	}

	if !p.expectResults {
		return errors.New("no more results")
	}
	p.opts.cursor = p.cursor

	req, err := p.client.newPagerRequest(p.method, p.path, p.opts, nil)
	if err != nil {
		return err
	}

	resp, err := p.client.do(ctx, req, dst)
	if err != nil {
		p.cursor = ""
		p.expectResults = false
		return err
	} else if p.cursor = resp.NextCursor; p.cursor == "" {
		p.expectResults = false
	}
	return nil
}

// setMaxPerPage sets the request to return the maximum number of results per page
// recurly will provide. This is useful for FetchAll() requests where it's best
// to limit the total number of HTTP requests made to retrieve all of the
// results.
func (p *pager) setMaxPerPage() {
	p.opts.PerPage = 200
}

// PagerOptions are used to send pagination parameters with paginated requests.
type PagerOptions struct {
	// Results per page. If not provides, Recurly defaults to 50.
	PerPage int

	// The field to sort by (e.g. created_at). See Recurly's documentation.
	Sort string

	// asc or desc
	Order string

	// Returns records greater than or equal to BeginTime.
	BeginTime NullTime

	// Returns records less than or equal to EndTime.
	EndTime NullTime

	State string // supported by some endpoints. Check Recurly's documenation.
	Type  string // supported by some endpoints. Check Recurly's documentation.

	// query is for any one-off URL params used by a specific endpoint.
	// Values sent as time.Time or recurly.NullTime will be automatically
	// converted to a valid datetime format for Recurly.
	query query

	cursor string // managed internally by this library
}

type query map[string]interface{}

func (q query) append(u *url.URL) {
	if len(q) == 0 {
		return
	}

	vals := u.Query()
	for key, val := range q {
		switch v := val.(type) {
		case string:
			if v != "" {
				vals.Add(key, v)
			}
		case time.Time:
			if !v.IsZero() {
				vals.Add(key, v.UTC().Format(DateTimeFormat))
			}
		case NullTime:
			if v.Time != nil && !v.IsZero() {
				vals.Add(key, v.String())
			}
		case bool:
			vals.Add(key, fmt.Sprintf("%t", v))
		case int, int64, uint, uint64:
			vals.Add(key, fmt.Sprintf("%d", v))
		default:
			vals.Add(key, fmt.Sprintf("%v", val))
		}
	}
	u.RawQuery = vals.Encode()
}

// append appends params to a URL.
func (p PagerOptions) append(u *url.URL) {
	if p.query == nil {
		p.query = map[string]interface{}{}
	}
	if p.PerPage > 0 {
		p.query["per_page"] = p.PerPage
	}
	if p.BeginTime.Time != nil && !p.BeginTime.IsZero() {
		p.query["begin_time"] = p.BeginTime
	}
	if p.EndTime.Time != nil && !p.EndTime.IsZero() {
		p.query["end_time"] = p.BeginTime
	}

	p.query["sort"] = p.Sort
	p.query["order"] = p.Order
	p.query["state"] = p.State
	p.query["type"] = p.Type
	p.query["cursor"] = p.cursor
	p.query.append(u)
}
