/*
Package recurly provides a client for using the Recurly API.

Usage

Construct a new Recurly client, then use the various services on the client to
access different parts of the Recurly API. For example:

	import "github.com/autopilot3/recurly"

	func main() {
		client := recurly.NewClient("your-subdomain", "APIKEY")

		// Retrieve an account
		a, err := client.Accounts.Get(context.Background(), "1")
	}

See the examples section for more usage examples.

Null Types

Null types provide a way to differentiate between zero values and real values.
For example, 0 is the zero value for ints; false for bools. Because those are
valid values, null types allow us to differentiate between bool false as a
valid value (we want to send to Recurly) and bool false as the zero value (we
do not want to send to recurly)

There are three null types: recurly.NullInt, recurly.NullBool, and recurly.NullTime

	a := recurly.Adjustment{
		UnitAmountInCents: recurly.NewInt(0),
		Taxable: recurly.NewBool(false),
		StartDate: recurly.NewTime(time.Now()),
	}

	// Null Int
	i := a.UnitAmountInCents.Int() // 0
	i, ok := a.UnitAmountInCents.Value() // 0, true
	iPtr := a.UnitAmountInCents.IntPtr() // Return *int, will be nil if value is invalid

	// NullBool
	b := a.Taxable.Bool() // false
	b, ok := a.Taxable.Value() // false, true
	bPtr := a.Taxable.BoolPtr() // Return *bool, will be nil if value is invalid

	// NullTime
	t := a.StartDate.Time() // time.Time
	t, ok := a.StartDate.Value() // time.Time, true
	tPtr := a.StartDate.TimePtr() // Return *time.Time, will be nil if value is invalid

If you have a pointer value, you can use New*Ptr() to return a new null type, where
the value is valid only if the pointer is non-nil:

	i := recurly.NewIntPtr(v) // where v is *int
	b := recurly.NewBoolPtr(b) // where b is *bool
	t := recurly.NewTimePtr(t) // where t is *time.Time. If non-nil, t.IsZero() must be false to be considered valid

Error Handling

Generally, checking that err != nil is sufficient to catch errors. However there
are some circumstances where you may need to know more about the specific error.

ClientError is returned for all 400-level responses with the exception of rate limit
errors and failed transactions. Here is an example of working with client errors:

	sub, err := client.Invoices.Create(ctx, "1", recurly.Invoice{})
	if e, ok := err.(*recurly.ClientError); ok {
		// Check status code (e.g. 404)
		if e.Response.Code == http.StatusNotFound {
			return err
		}

		// Check for a specific validation symbol in one of the error messages
		if e.Is("will_not_invoice") {
			return err
		}

		// Or loop through all of the validation errors
		if len(e.ValidationErrors) > 0 {
			for _, ve := range e.ValidationErrors {
				// Access each validation error here
			}
		}
		return err
	} else if err != nil {
		return err
	}

TransactionFailedError is returned for any endpoint where a transaction was
attempted and failed. It is highly recommended that you check for this error when
using any endpoint that creates a transaction.

	_, err := client.Purchases.Create(ctx, recurly.Purchase{})
	if e, ok := err.(*recurly.TransactionFailedError); ok {
		// e.Transaction holds the failed transaction (if available)
		// e.TransactionError holds the specific error. See the TransactionError
		// struct for specifics on the fields.
	} else if err != nil {
		// Handle all other errors
	}

ServerError operates the same way as ClientError, except it's returned for 500-level
responses. It only contains the *http.Response. This allows you to differentiate
retriable errors (e.g. 503 Service Unavailable) from bad requests (e.g.
400 Bad Request).

RateLimitError is returned when the rate limit is exceeded. The Rate field contains
information on the amount of requests and when the rate limit will reset.

For more on errors, see the examples section below.

Get Methods

When retrieving an individual item (e.g. account, invoice, subscription): if the
item is not found, Recurly will return a 404 Not Found. Typically this will return
a *recurly.ClientError. The only exception is for any function named 'Get':
a nil item and nil error will be returned if the item is not found.

	a, err := client.Accounts.Get(ctx, "1")
	if err != nil {
		return err
	} else if a == nil {
		// account not found
	}


Pagination

All requests for resource collections support pagination. Pagination options are
described in the recurly.PagerOptions struct and passed to the list methods directly.

	client := recurly.NewClient("your-subdomain", "APIKEY")

	// Initialize a pager with any pagination options needed.
	pager := client.Accounts.List(&recurly.PagerOptions{
		State: recurly.AccountStateActive,
	})

	// Count the records (if desired)
	count, err := pager.Count(ctx)
	if err != nil {
		return err
	}

	// Or iterate through each of the pages
	for pager.Next() {
		var accounts []recurly.Account
		if err := pager.Fetch(ctx, &accounts); err != nil {
			return err
		}
	}

You can also let the library paginate for you and return all of the results at once:

	pager := client.Accounts.List(nil)
	var accounts []recurly.Account
	if err := pager.FetchAll(ctx, &accounts); err != nil {
		return err
	}

In some cases, you may want to paginate non-consecutively. For example, if you have
paginated results being sent to a frontend, and the frontend is providing your
app the next cursor.

In that case you can obtain the next cursor like this (although it may be empty):

	cursor := pager.Cursor()

If you have a cursor, you can provide *PagerOptions with it to start paginating
the next result set:

	pager := client.Accounts.List(&recurly.PagerOptions{
		Cursor: cursor,
	})

*/
package recurly
