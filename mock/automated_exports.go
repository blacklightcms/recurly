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
}

func (m *AutomatedExportsService) Get(ctx context.Context, date time.Time, fileName string) (*recurly.AutomatedExport, error) {
	m.GetInvoked = true
	return m.OnGet(ctx, date, fileName)
}
