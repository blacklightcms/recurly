package recurly

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
)

// AutomatedExportsService manages the interactions for automated exports.
type AutomatedExportsService interface {
	// Get retrieves export file.
	//
	// https://dev.recurly.com/docs/download-export-file
	Get(ctx context.Context, date time.Time, fileName string) (*AutomatedExport, error)
}

// AutomatedExport holds export file info.
type AutomatedExport struct {
	XMLName     xml.Name `xml:"export_file"`
	ExpiresAt   NullTime `xml:"expires_at,omitempty"`
	DownloadURL string   `xml:"download_url,omitempty"`
}

var _ AutomatedExportsService = &automatedExportsImpl{}

// automatedExportsImpl implements AutomatedExportsService.
type automatedExportsImpl serviceImpl

func (s *automatedExportsImpl) Get(ctx context.Context, date time.Time, fileName string) (*AutomatedExport, error) {
	d := fmt.Sprintf("%02d-%02d-%02d", date.Year(), date.Month(), date.Day())
	path := fmt.Sprintf("/export_dates/%s/export_files/%s", d, fileName)
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var dst AutomatedExport
	if _, err := s.client.do(ctx, req, &dst); err != nil {
		if e, ok := err.(*ClientError); ok && e.Response.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &dst, nil
}
