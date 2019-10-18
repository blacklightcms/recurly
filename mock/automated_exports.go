package mock

import (
	"context"
	"github.com/blacklightcms/recurly"
)

var _ recurly.AutomatedExportsService = &AutomatedExportsService{}

// AdjustmentsService manages the interactions for adjustments.
type AutomatedExportsService struct {
	OnGet      func(ctx context.Context, date string, fileName string) (*recurly.AutomatedExport, error)
	GetInvoked bool
}

func (m *AutomatedExportsService) Get(ctx context.Context, date string, fileName string) (*recurly.AutomatedExport, error) {
	m.GetInvoked = true
	return m.OnGet(ctx, date, fileName)
}
