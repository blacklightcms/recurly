package recurly_test

import (
	"bytes"
	"context"
	"encoding/xml"
	"log"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/blacklightcms/recurly"
	"github.com/google/go-cmp/cmp"
)

// Ensure structs are encoded to XML properly.
func TestAutomatedExports_Encoding(t *testing.T) {
	now := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		v        interface{}
		expected string
	}{
		{
			v: recurly.AutomatedExport{ExpiresAt: recurly.NewTime(now), DownloadURL: "https://recurly.com/sub.csv.gz"},
			expected: MustCompactString(`
				<export_file>
					<expires_at>2000-01-01T00:00:00Z</expires_at>
					<download_url>https://recurly.com/sub.csv.gz</download_url>
				</export_file>
			`),
		},
		{
			v: recurly.ExportDate{Date: "2019-10-10"},
			expected: MustCompactString(`
				<export_date>
					<date>2019-10-10</date>
				</export_date>
			`),
		},
		{
			v: recurly.ExportFile{Name: "account_notes_created.csv.gz"},
			expected: MustCompactString(`
				<export_file>
					<name>account_notes_created.csv.gz</name>
				</export_file>
			`),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			buf := new(bytes.Buffer)
			if err := xml.NewEncoder(buf).Encode(tt.v); err != nil {
				t.Fatal(err)
			} else if buf.String() != tt.expected {
				log.Print(tt.expected)
				t.Fatal(buf.String())
			}
		})
	}
}

func TestAutomatedExports_Get(t *testing.T) {
	now := time.Date(2015, 2, 4, 23, 13, 7, 0, time.UTC)

	t.Run("OK", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/export_dates/2015-02-04/export_files/sub.csv.gz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write(MustOpenFile("automated_export.xml"))
		}, t)

		if a, err := client.AutomatedExports.Get(context.Background(), now, "sub.csv.gz"); err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(a, NewTestAutomatedExport()); diff != "" {
			t.Fatal(diff)
		} else if !s.Invoked {
			t.Fatal("expected fn invocation")
		}
	})

	// Ensure a 404 returns nil values.
	t.Run("ErrNotFound", func(t *testing.T) {
		client, s := recurly.NewTestServer()
		defer s.Close()

		s.HandleFunc("GET", "/v2/export_dates/2015-02-04/export_files/sub.csv.gz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}, t)

		if a, err := client.AutomatedExports.Get(context.Background(), now, "sub.csv.gz"); !s.Invoked {
			t.Fatal("expected fn invocation")
		} else if err != nil {
			t.Fatal(err)
		} else if a != nil {
			t.Fatalf("expected nil: %#v", a)
		}
	})
}

func TestAutomatedExports_ListDates(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	var invocations int
	s.HandleFunc("GET", "/v2/export_dates", func(w http.ResponseWriter, r *http.Request) {
		invocations++
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("export_dates.xml"))
	}, t)

	pager := client.AutomatedExports.ListDates(nil)
	for pager.Next() {
		var dates []recurly.ExportDate
		if err := pager.Fetch(context.Background(), &dates); err != nil {
			t.Fatal(err)
		} else if !s.Invoked {
			t.Fatal("expected s to be invoked")
		} else if diff := cmp.Diff(dates, []recurly.ExportDate{*NewTestExportDate()}); diff != "" {
			t.Fatal(diff)
		}
	}
	if invocations != 1 {
		t.Fatalf("unexpected number of invocations: %d", invocations)
	}
}

func TestAutomatedExports_ListFiles(t *testing.T) {
	now := time.Date(2019, 10, 10, 23, 13, 7, 0, time.UTC)
	client, s := recurly.NewTestServer()
	defer s.Close()

	var invocations int
	s.HandleFunc("GET", "/v2/export_dates/2019-10-10/export_files", func(w http.ResponseWriter, r *http.Request) {
		invocations++
		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("export_files.xml"))
	}, t)

	pager := client.AutomatedExports.ListFiles(now, nil)
	for pager.Next() {
		var files []recurly.ExportFile
		if err := pager.Fetch(context.Background(), &files); err != nil {
			t.Fatal(err)
		} else if !s.Invoked {
			t.Fatal("expected s to be invoked")
		} else if diff := cmp.Diff(files, []recurly.ExportFile{*NewTestExportFile()}); diff != "" {
			t.Fatal(diff)
		}
	}
	if invocations != 1 {
		t.Fatalf("unexpected number of invocations: %d", invocations)
	}
}

// Returns an AutomatedExport corresponding to testdata/adjustment.xml.
func NewTestAutomatedExport() *recurly.AutomatedExport {
	return &recurly.AutomatedExport{
		XMLName:     xml.Name{Local: "export_file"},
		ExpiresAt:   recurly.NewTime(MustParseTime("2015-02-04T23:13:07Z")),
		DownloadURL: "https://recurly.s3.amazonaws.com/file",
	}
}

// Returns an ExportDate corresponding to testdata/export_dates.xml.
func NewTestExportDate() *recurly.ExportDate {
	return &recurly.ExportDate{
		XMLName: xml.Name{Local: "export_date"},
		Date:    "2019-10-10",
	}
}

// Returns an ExportFile corresponding to testdata/export_files.xml.
func NewTestExportFile() *recurly.ExportFile {
	return &recurly.ExportFile{
		XMLName: xml.Name{Local: "export_file"},
		Name:    "account_notes_created.csv.gz",
	}
}
