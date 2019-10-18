package mock

import (
	"context"
	"time"

	"github.com/blacklightcms/recurly"
)

var _ recurly.AutomatedExportsService = &AutomatedExportsService{}

// AutomatedExportsService manages the interactions for automated exports.
type AutomatedExportsService struct {
	OnGet      func(ctx context.Context, date time.Time, fileName string) (*recurly.AutomatedExport, error)
	GetInvoked bool

	OnListDates      func(opts *recurly.PagerOptions) recurly.Pager
	ListDatesInvoked bool

	OnListFiles      func(date time.Time, opts *recurly.PagerOptions) recurly.Pager
	ListFilesInvoked bool
}

func (m *AutomatedExportsService) Get(ctx context.Context, date time.Time, fileName string) (*recurly.AutomatedExport, error) {
	m.GetInvoked = true
	return m.OnGet(ctx, date, fileName)
}

func (m *AutomatedExportsService) ListDates(opts *recurly.PagerOptions) recurly.Pager {
	m.ListDatesInvoked = true
	return m.OnListDates(opts)
}

func (m *AutomatedExportsService) ListFiles(date time.Time, opts *recurly.PagerOptions) recurly.Pager {
	m.ListFilesInvoked = true
	return m.OnListFiles(date, opts)
}
