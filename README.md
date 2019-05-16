# Recurly Client for Go

 [![Build Status](https://travis-ci.org/blacklightcms/recurly.svg?branch=master)](https://travis-ci.org/blacklightcms/recurly)  [![GoDoc](https://godoc.org/github.com/blacklightcms/recurly?status.svg)](https://godoc.org/github.com/blacklightcms/recurly)

 Recurly is a Go (golang) API Client for the [Recurly](https://recurly.com/) API. It is actively maintained, unit tested, and uses no external dependencies. The vast majority of the API is implemented.

 Supports:
  - Recurly API `v2.19`
  - Accounts
  - Add Ons
  - Adjustments
  - Billing
  - Coupons
  - Credit Payments
  - Invoices
  - Plans
  - Purchases
  - Redemptions
  - Shipping Addresses
  - Subscriptions
  - Transactions

## Installation
Install:

```shell
go get github.com/blacklightcms/recurly
```

Import:
```go
import "github.com/blacklightcms/recurly"
```

Resources:
 - [API Docs](https://godoc.org/github.com/blacklightcms/recurly)
 - [Examples](https://godoc.org/github.com/blacklightcms/recurly#pkg-examples)

## Quickstart

Construct a new Recurly client, then use the various services on the client to access different parts of the Recurly API. For example:

```go
client := recurly.NewClient("your-subdomain", "APIKEY")

// Retrieve an account
a, err := client.Accounts.Get(context.Background(), "1")
```

## Examples and How To
Please go through [examples](https://godoc.org/github.com/blacklightcms/recurly#pkg-examples) for detailed examples of using this package.

There are a few high-level notes below that will be helpful in getting the most out of this library.

- [Null Types](#null-types)
- [Error Handling](#error-handling)
- [Get Methods](#get-methods)
- [Pagination](#pagination)

## Null Types
In some cases, Recurly requires values to be sent if valid, but the XML tag to be omitted if not valid.

If you want to create a subscription with a price of `$0`, using an `int` type with an `omitempty` XML tag the value would never actually be sent. `0` is the zero value for an `int`.

You want this:

```xml
<subscription>
   <unit_amount_in_cents>0</unit_amount_in_cents>
</subscription>
```

But instead send this:

```xml
<subscription></subscription>
```

For `int`, `bool`, and `time.Time` types, null types provide a way to differentiate between zero values 
and real values. Let's look at integers:

Zero as a valid value:
```go
v := recurly.NewInt(0)
i := v.Int() // returns 0
i, ok := v.Value() // returns 0, true
```

Zero as an invalid value:
```go
// Zero value, int value is 0 but flagged as invalid
var nullInt recurly.NullInt
i := nullInt.Int() // returns 0
i, ok := nullInt.Value() // returns 0, false
```


## Error Handling
There are four important error types to know about. While you can generally just check for a non-nil 
error (`if err != nil`), this section will cover more advanced error usage you can implement where it makes sense for your application.

### ClientError
`ClientError` is returned for all 400-level responses with the exception of
1) `429 Too Many Requests`. See `RateLimitError`.
2) `422 Unprocessable Entity` with a failed transaction. See `TransactionFailedError`.

Inspecting this error can be useful when you are looking for specific status codes 
and/or you want to look at any validation errors from Recurly.

```go
sub, err := client.Invoices.Create(ctx, "1", recurly.Invoice{...})
if err != nil {
    if e, ok := err.(*recurly.ClientError); ok {
        // e.ValidationErrors contains any validation errors from Recurly

        // Use e.Is() to see if any of the errors contain <symbol>will_not_invoice</symbol> (for example)
        if e.Is("will_not_invoice") {
            // ...
        }

        // e.Response contains the *http.Response. You can check the status code as well.
        if e.Response.Code == http.StatusBadRequest {
            // ...
        }
    }
}
```

### TransactionFailedError
`TransactionFailedError` is returned for any endpoint where a transaction was attempted and failed. 
It is recommended that you check for this error when using any endpoint that 
creates a transaction.

```go
_, err := client.Purchases.Create(ctx, recurly.Purchase{...})
if e, ok := err.(*recurly.TransactionFailedError); ok {
    // e.Transaction holds the failed transaction (if available)
    // e.TransactionError holds the specific error. See godoc for specific fields.
} else if err != nil {
    // Handle all other errors
}
```

### ServerError
`ServerError` operates the same way as `ClientError`, except it's for 500-level responses and only 
contains the `*http.Response`. This allows you to differentiate retriable errors from bad requests. 

> NOTE: You will generally just capture this as a generic error (`if err != nil`) unless you need to look for something specific.

### RateLimitError
`RateLimitError` is returned when your request has exceeded your rate limit. It contains 
information on the rate limit and when the limit will reset.

> NOTE: You will generally just capture this as a generic error (`if err != nil`) unless you need to look for something specific.

```go
a, err := client.Accounts.Get(ctx, "1")
if err != nil {
    if e, ok := err.(*recurly.RateLimitError); ok {
        // e.Rate.Limit holds the total request limit during the 5 minute window
        // e.Reset holds the time when the current window will completely reset
    }
}
```

You can easily combine these if checking multiple errors using a type switch:

```go
_, err := client.Purchases.Create(ctx, recurly.Purchase{...})
if err != nil {
    switch err := err.(type) {
    case *recurly.TransactionFailedErr:
        // Determine why the transaction failed
    case *recurly.ClientError:
        // Inspect error for details of what went wrong
    case *recurly.ServerError:
        // Retryable
    default:
        return err
    }
}
```

### Get Methods
When retrieving an individual item (such as an account, invoice, subscription): if the item is not found, a nil item and nil error will be returned. This is standard for all functions named `Get()`. All other functions would return a `ClientError` when encountering a `404 Not Found` status code.

```go
a, err := client.Accounts.Get(ctx, "1")
if err != nil {
    return err
} else if a == nil {
    // Account not found
}
```

### Pagination
Any function named `List()` or `List*()` is a pagination function. They all return a pager which standardizes how pagination works.

In those cases, pagination always works the same way. Here is an example of how to paginate accounts:

```go
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
    a, err := pager.Fetch(ctx)
    if err != nil {
        return err
    }
    // Do something with a
}
```

You can also let the library paginate for you and return all of the results in a slice.

```go
pager := client.Accounts.List(nil)
a, err := pager.FetchAll(ctx)
if err != nil {
    return err
}
```

## Migration
If migrating from a previous version of the library, there was a large refactor with breaking changes released to address some design issues with the library. See the migration guide for steps on how to migrate to the latest version.

This is recommended for all users.