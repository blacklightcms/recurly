package mock

import "io"

// Webhooks represents the interactions available for webhooks.
type Webhooks struct {
	OnParse      func(r io.Reader) (interface{}, error)
	ParseInvoked bool
}

func (m *Webhooks) Parse(r io.Reader) (interface{}, error) {
	m.ParseInvoked = true
	return m.OnParse(r)
}
