package mock_test

import (
	"context"
	"testing"

	"github.com/autopilot3/recurly"
	"github.com/autopilot3/recurly/mock"
	"github.com/google/go-cmp/cmp"
)

func ExampleNewClient(t *testing.T) {
	// Test data.
	const AccountID = "10"

	// Initialize client.
	client := mock.NewClient("your-subdomain", "key")

	// Setup mock.
	client.Accounts.OnGet = func(ctx context.Context, id string) (*recurly.Account, error) {
		if id != AccountID {
			t.Fatalf("unexpected account id: %s", id)
		}
		return &recurly.Account{
			Code:      AccountID,
			State:     "active",
			Email:     "verena@example.com",
			FirstName: "Verena",
			LastName:  "Example",
		}, nil
	}

	// Retrieve account.
	// Verify results and ensure client.Accounts.Get() was invoked.
	if a, err := client.Accounts.Get(context.Background(), "10"); err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(a, &recurly.Account{
		Code:      AccountID,
		State:     "active",
		Email:     "verena@example.com",
		FirstName: "Verena",
		LastName:  "Example",
	}); diff != "" {
		t.Fatal(diff)
	} else if !client.Accounts.GetInvoked {
		t.Fatal("expected Accounts.Get() to be invoked")
	}
}

func ExampleNewClient_pager(t *testing.T) {
	// Test data.
	results := []recurly.Account{
		{
			Code:      "1",
			State:     "active",
			Email:     "verena@foo.com",
			FirstName: "Verena",
			LastName:  "Example",
		},
		{
			Code:      "2",
			State:     "active",
			Email:     "verena@bar.com",
			FirstName: "Foo",
			LastName:  "Bar",
		},
	}

	// Initialize client.
	client := mock.NewClient("your-subdomain", "key")

	// Setup pager mock.
	var invocations int
	accountsPager := &mock.Pager{
		OnNext: func() bool {
			return invocations < len(results)
		},
		OnFetch: func(ctx context.Context, dst interface{}) error {
			// Fetch should not be called more times than Next() returns true.
			if invocations == len(results) {
				t.Fatalf("unexpected invocation: %d", invocations)
			}

			// Assert that dst is what we expect.
			accounts, ok := dst.(*[]recurly.Account)
			if !ok {
				t.Fatalf("unexpected type for dst: %T", dst)
			}

			// Set a slice with a length of 1 using the invocation number
			// as the index
			*accounts = []recurly.Account{results[invocations]}
			invocations++ // increment
			return nil
		},
	}

	// Setup client.Accounts.List() mock.
	// Ensure List() is called correctly, then return mock pager.
	client.Accounts.OnList = func(opts *recurly.PagerOptions) recurly.Pager {
		if diff := cmp.Diff(opts, &recurly.PagerOptions{
			PerPage: 1,
			State:   "active",
		}); diff != "" {
			t.Fatal(diff)
		}
		return accountsPager
	}

	// Paginate.
	pager := client.Accounts.List(&recurly.PagerOptions{
		PerPage: 1,
		State:   "active",
	})

	var accounts []recurly.Account
	for pager.Next() {
		var a []recurly.Account
		if err := pager.Fetch(context.Background(), &a); err != nil {
			t.Fatal(err)
		}
		accounts = append(accounts, a...)
	}

	// Check results
	if diff := cmp.Diff(accounts, results); diff != "" {
		t.Fatal(diff)
	} else if !client.Accounts.ListInvoked {
		t.Fatal("expected Accounts.List() to be invoked")
	} else if !accountsPager.NextInvoked {
		t.Fatal("expected Next() to be invoked on pager")
	} else if !accountsPager.FetchInvoked {
		t.Fatal("expected Fetch() to be invoked on pager")
	} else if invocations != len(results) {
		t.Fatalf("unexpected number of invocations: %d", invocations)
	}
}
