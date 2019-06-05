package recurly_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/blacklightcms/recurly"
	"github.com/google/go-cmp/cmp"
)

func TestPager(t *testing.T) {
	const CursorA = "1972702718353176814:A1465932489"

	client, s := recurly.NewTestServer()
	defer s.Close()

	var invocations int
	s.HandleFunc("GET", "/v2/accounts", func(w http.ResponseWriter, r *http.Request) {
		cursor := r.URL.Query().Get("cursor")
		switch invocations {
		case 0:
			if cursor != "" {
				t.Fatalf("unexpected cursor: %s", cursor)
			}
			w.Header().Set("Link", `<https://test.recurly.com/v2/accounts?cursor=`+CursorA+`>; rel="next"`)
		case 1:
			if cursor != CursorA {
				t.Fatalf("unexpected cursor: %s", cursor)
			}
		default:
			t.Fatalf("unexpected number of invocations")
		}

		query := r.URL.Query()
		query.Del("cursor") // conditionally checked above
		if diff := cmp.Diff(query, url.Values{
			"per_page":   []string{"50"},
			"sort":       []string{"created_at"},
			"order":      []string{"asc"},
			"state":      []string{"active"},
			"begin_time": []string{"2011-10-17T17:24:53Z"},
			"end_time":   []string{"2011-10-18T17:24:53Z"},
		}); diff != "" {
			t.Fatal(diff)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("accounts.xml"))
		invocations++
	}, t)

	pager := client.Accounts.List(&recurly.PagerOptions{
		PerPage:   50,
		Sort:      "created_at",
		Order:     "asc",
		State:     "active",
		BeginTime: recurly.NewTime(MustParseTime("2011-10-17T17:24:53Z")),
		EndTime:   recurly.NewTime(MustParseTime("2011-10-18T17:24:53Z")),
	})

	for pager.Next() {
		var a []recurly.Account
		if err := pager.Fetch(context.Background(), &a); err != nil {
			t.Fatal(err)
		} else if !s.Invoked {
			t.Fatal("expected s to be invoked")
		} else if diff := cmp.Diff(a, []recurly.Account{*NewTestAccount()}); diff != "" {
			t.Fatal(diff)
		}

		// Check cursor.
		switch invocations {
		case 1:
			if pager.Cursor() == CursorA {
				break
			}
			fallthrough
		case 2:
			if pager.Cursor() == "" {
				break
			}
			fallthrough
		default:
			t.Fatalf("unexpected cursors on invocation %d: cursor=%s", invocations, pager.Cursor())
		}

		s.Invoked = false
	}
}

// Verify pager can accept a cursor and send it through to Recurly as expected.
func TestPager_CursorProvided(t *testing.T) {
	const CursorA = "CURSOR_A"
	const CursorB = "CURSOR_B"

	client, s := recurly.NewTestServer()
	defer s.Close()

	var invocations int
	s.HandleFunc("GET", "/v2/accounts", func(w http.ResponseWriter, r *http.Request) {
		cursor := r.URL.Query().Get("cursor")
		switch invocations {
		case 0:
			if cursor != CursorA {
				t.Fatalf("unexpected cursor: %s", cursor)
			}
			w.Header().Set("Link", `<https://test.recurly.com/v2/accounts?cursor=`+CursorB+`>; rel="next"`)
		case 1:
			if cursor != CursorB {
				t.Fatalf("unexpected cursor: %s", cursor)
			}
		default:
			t.Fatalf("unexpected number of invocations")
		}

		w.WriteHeader(http.StatusOK)
		w.Write(MustOpenFile("accounts.xml"))
		invocations++
	}, t)

	pager := client.Accounts.List(&recurly.PagerOptions{
		Cursor: CursorA, // Provide starting cursor
	})

	// Verify starting cursor is set.
	if pager.Cursor() != CursorA {
		t.Fatalf("unexpected next cursor: %s", pager.Cursor())
	}

	for pager.Next() {
		var a []recurly.Account
		if err := pager.Fetch(context.Background(), &a); err != nil {
			t.Fatal(err)
		} else if !s.Invoked {
			t.Fatal("expected s to be invoked")
		} else if diff := cmp.Diff(a, []recurly.Account{*NewTestAccount()}); diff != "" {
			t.Fatal(diff)
		}

		// Check cursor.
		switch invocations {
		case 1:
			if pager.Cursor() == CursorB {
				break
			}
			fallthrough
		case 2:
			if pager.Cursor() == "" {
				break
			}
			fallthrough
		default:
			t.Fatalf("unexpected cursors on invocation %d: cursor=%s", invocations, pager.Cursor())
		}

		s.Invoked = false
	}
}
