package mock

import (
	"context"

	"github.com/blacklightcms/recurly"
)

var _ recurly.Pager = &Pager{}

type Pager struct {
	OnCount      func(ctx context.Context) (int, error)
	CountInvoked bool

	OnNext      func() bool
	NextInvoked bool

	OnCursor      func() string
	CursorInvoked bool

	OnFetch      func(ctx context.Context, dst interface{}) error
	FetchInvoked bool

	OnFetchAll      func(ctx context.Context, dst interface{}) error
	FetchAllInvoked bool
}

func (m *Pager) Count(ctx context.Context) (int, error) {
	m.CountInvoked = true
	return m.OnCount(ctx)
}

func (m *Pager) Next() bool {
	m.NextInvoked = true
	return m.OnNext()
}

func (m *Pager) Cursor() string {
	m.CursorInvoked = true
	return m.OnCursor()
}

func (m *Pager) Fetch(ctx context.Context, dst interface{}) error {
	m.FetchInvoked = true
	return m.OnFetch(ctx, dst)
}

func (m *Pager) FetchAll(ctx context.Context, dst interface{}) error {
	m.FetchAllInvoked = true
	return m.OnFetchAll(ctx, dst)
}
