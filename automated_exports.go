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

	// ListDates returns a list of dates with export files.
	//
	// https://dev.recurly.com/v2.8/docs/list-export-dates
	ListDates(opts *PagerOptions) Pager

	// ListFiles returns a list of files available for the date specified.
	//
	// https://dev.recurly.com/v2.8/docs/list-export-files
	ListFiles(date time.Time, opts *PagerOptions) Pager
}

// AutomatedExport holds export file info.
type AutomatedExport struct {
	XMLName     xml.Name `xml:"export_file"`
	ExpiresAt   NullTime `xml:"expires_at,omitempty"`
	DownloadURL string   `xml:"download_url,omitempty"`
}

// ExportDate holds export date info.
type ExportDate struct {
	XMLName xml.Name `xml:"export_date"`
	Date    string   `xml:"date,omitempty"`
}

// ExportFile holds export file info.
type ExportFile struct {
	XMLName xml.Name `xml:"export_file"`
	Name    string   `xml:"name,omitempty"`
}

var _ AutomatedExportsService = &automatedExportsImpl{}

// automatedExportsImpl implements AutomatedExportsService.
type automatedExportsImpl serviceImpl

func (s *automatedExportsImpl) Get(ctx context.Context, date time.Time, fileName string) (*AutomatedExport, error) {
	d := date.Format("2006-01-02")
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

func (s *automatedExportsImpl) ListDates(opts *PagerOptions) Pager {
	return s.client.newPager("GET", "/export_dates", opts)
}

func (s *automatedExportsImpl) ListFiles(date time.Time, opts *PagerOptions) Pager {
	d := date.Format("2006-01-02")
	path := fmt.Sprintf("/export_dates/%s/export_files", d)
	return s.client.newPager("GET", path, opts)
}
