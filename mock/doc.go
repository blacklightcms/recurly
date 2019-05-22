/*
Package mock provides mocks attached to the Recurly client for testing
code using the Recurly API client.

This package makes it easy to test code using the Recurly API by focusing on
the arguments passed into each function and returning the expected result as
structs from the Recurly package. There is no need to deal with XML or make any HTTP requests.

Simple Setup

The basic setup for a test involves creating a new mock client, attaching mocks,
and returning expected results.

	func TestFoo(t *testing.T) {
		client := mock.NewClient("subdomain", "key")
		client.Billing.OnGet = func(ctx context.Context, accountCode string) (*recurly.Billing, error) {
			if accountCode != "10" {
				t.Fatalf("unexpected account code: %s", accountCode)
			}
			return &recurly.Billing{
				FirstName: "Foo",
				LastName: "Bar",
			}
		}

		if b, err := client.Billing.Get(context.Background(), "10"); err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(b, &recurly.Billing{
			FirstName: "Foo",
			LastName: "Bar",
		}); diff != "" {
			t.Fatal(diff)
		} else if !client.Billing.GetInvoked {
			t.Fatal("expected Get() to be invoked")
		}
	}

More Common Setup

If you created your own wrapper type to the library, let's call it PaymentsProvider:

	// MyBillingInfo is a custom billing holder. *recurly.Billing is converted to
	// *MyBillingInfo by combining FirstName and LastName into a single Name field.
	type MyBillingInfo{
		Name string
	}

	type PaymentsProvider interface {
		GetBilling(ctx context.Context, id int) (*MyBillingInfo, error)
	}

	// Provider implements PaymentsProvider.
	type Provider struct {
		Client *recurly.Client
	}

	func New(subdomain, key string) *Provider {
		return &Provider{Client: recurly.NewClient(subdomain, key)}
	}

	// GetBilling calls Recurly and converts *recurly.BillingInfo to *MyBillingInfo.
	func (p *Provider) GetBilling(ctx context.Context, id int) (*MyBillingInfo, error) {
		// Retrieve billing info from Recurly.
		b, err := p.Client.Billing.Get(ctx, strconv.Atoi(id))
		if err != nil {
			return nil, err
		} else if b == nil {
			return nil, nil
		}

		// Convert to MyBillingInfo.
		return &MyBillingInfo{
			Name: fmt.Sprintf("%s %s", b.FirstName, b.LastName),
		}
	}

Then in your test suite you might configure your own wrapper similar to mock.Client

	package foo_test

	import (
		"context"

		"github.com/your-project/foo"
		"github.com/blacklightcms/recurly/mock"
	)

	// Provider is a test wrapper for foo.Provider.
	type Provider struct {
		*foo.Provider

		Billing mock.Billing
	}

	// NewProvider returns a new test provider.
	func NewProvider() *Provider {
		// Init Provider.
		p := &Provider{Provider: foo.New("foo", "bar")}

		// Point Recurly Client's Billing Service to your mock
		// attached to p.
		p.Provider.Client.Billing = &p.Billing
		return p
	}

	func TestBilling(t *testing.T) {
		// Init test Provider.
		p := NewProvider()

		// Mock Recurly's response
		p.Billing.OnGet = func(ctx context.Context, id int) (*recurly.Billing, error) {
			return &recurly.Billing{
				FirstName: "Verena",
				LastName: "Example",
			}
		}

		// Call your provider and assert *MyBillingInfo
		if b, err := p.GetBilling(context.Background(), 10); err != nil {
			t.Fatal(err)
		} else if diff := cmp.Diff(b, &foo.MyBillingInfo{
			Name: "Verena Example",
		}); diff != "" {
			t.Fatal(diff)
		} else if !p.Billing.GetInvoked {
			t.Fatal("expected Get() to be invoked")
		}
	}

See examples for more.

*/
package mock
