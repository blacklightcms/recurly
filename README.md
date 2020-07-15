# Recurly Client for Go

 [![Build Status](https://travis-ci.org/blacklightcms/recurly.svg?branch=master)](https://travis-ci.org/blacklightcms/recurly)  [![GoDoc](https://godoc.org/github.com/blacklightcms/recurly?status.svg)](https://godoc.org/github.com/blacklightcms/recurly/)

 Recurly is a Go (golang) API Client for the [Recurly](https://recurly.com/) API. It is actively maintained, unit tested, and uses no external dependencies. The vast majority of the API is implemented.

 Supports:
  - Recurly API `v2.27`
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
  - Shipping Methods
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
 - [API Docs](https://godoc.org/github.com/blacklightcms/recurly/)
 - [Examples](https://godoc.org/github.com/blacklightcms/recurly/#pkg-examples)

## Note on v1 and breaking changes
If migrating from a previous version of the library, there was a large refactor with breaking changes released to address some design issues with the library. See the [migration guide](https://github.com/blacklightcms/recurly/wiki/v1-Migration-Guide) for steps on how to migrate to the latest version.

This is recommended for all users.

## Quickstart

Construct a new Recurly client, then use the various services on the client to access different parts of the Recurly API. For example:

```go
client := recurly.NewClient("your-subdomain", "APIKEY")

// Retrieve an account
a, err := client.Accounts.Get(context.Background(), "1")
```

## Examples and How To
Please go through [examples](https://godoc.org/github.com/blacklightcms/recurly/#pkg-examples) for detailed examples of using this package.

The examples explain important cases like:

- Null Types
- Error Handling
- Get Methods
- Pagination

Here are a few snippets to demonstrate library usage.

### Create Account
```go
account, err := client.Accounts.Create(ctx, recurly.Account{
    Code: "1",
    FirstName: "Verena",
    LastName: "Example",
    Email: "verena@example.com",
})
```

> **NOTE**: An account can also be created along a subscription by embedding the 
> account in the subscription during creation. The purchases API also supports 
> this, and likely other endpoints. See Recurly's documentation for details.

### Get Account
```go
account, err := client.Accounts.Get(ctx, "1")
if err != nil {
    return err
} else if account == nil {
    // account not found
    // Note: this nil, nil response on 404s is unique to Get() methods
    // See GoDoc for details.
}
```

### Create Billing Info
```go
// Using token obtained with recurly.js
// If you want to set billing info directly, omit the token and set the
// corresponding fields on the recurly.Billing struct.
billing, err := client.Billing.Create("1", recurly.Billing{
    Token: token,
})
```
> **NOTE**: See the error handling section in GoDoc for how to handle transaction errors

### Creating Purchases

```go
purchase, err := c.Client.Purchases.Create(ctx, recurly.Purchase{
    Account: recurly.Account{
	    Code: "1",
    },
    Adjustments: []recurly.Adjustment{{
	    UnitAmountInCents: recurly.NewInt(100),
	    Description:       "Purchase Description",
	    ProductCode:       "product_code",
    }},
    CollectionMethod: recurly.CollectionMethodAutomatic,
    Currency:         "USD",
})
if err != nil {
    // NOTE: See GoDoc for how to handle failed transaction errors
}
```

> **NOTE**: The purchases API supports subscriptions, adjustments, shipping addresses,
> shipping fees, and more. This is one of many possible examples. See the underlying
> structs and [Recurly's documentation](https://dev.recurly.com/docs/create-purchase) for more info.

### Creating Subscriptions
```go
subscription, err := client.Subscriptions.Create(ctx, recurly.NewSubscription{
    PlanCode: "gold",
    Currency: "USD",
    Account: recurly.Account{
        // Note: Set the Code for an existing account
        // To create a new account, omit Code but provide other fields
    },
})
if err != nil {
    // NOTE: See GoDoc for how to handle failed transaction errors
    return err
}
```
> **NOTE**: Recurly offers several other ways to create subscriptions, often embedded 
> within other requests (such as the `Purchases.Create()` call). See Recurly's 
> documentation for more details.

## Webhooks
This library supports webhooks via the `webhooks` sub package. 

The usage is to parse the webhook from a reader, then use a switch statement 
to determine the type of webhook received.

```go
// import "github.com/blacklightcms/recurly/webhooks"

hook, err := webhooks.Parse(r)
if e, ok := err.(*webhooks.ErrUnknownNotification); ok {
    // e.Name() holds the name of the notification
} else if err != nil {
    // all other errors
}

// Use a switch statement to determine the type of webhook received.
switch h := hook.(type) {
case *webhooks.AccountNotification:
    // h.Account
case *webhooks.PaymentNotification:
    // h.Account
    // h.Transaction
case *webhooks.SubscriptionNotification:
    // h.Account
    // h.Subscription
default:
    // webhook not listed above
}
```

## Testing
Once you've imported this library into your application, you will want to add tests.

Internally this library sets up a test HTTPs server and validates methods, paths, 
query strings, request body, and returns XML. You will not need to worry about those internals
when testing your own code that uses this library.

Instead we recommend using the `mock` package. The `mock` package provides mocks 
for all of the different services in this library.

For examples of how to test your code using mocks, visit the [GoDoc examples](https://godoc.org/github.com/blacklightcms/recurly/mock/).

> **NOTE**: If you need to go beyond mocks and test requests/responses, `testing.go` exports `TestServer`. This is how the library tests itself. See the GoDoc or the `*_test.go` files for usage examples.

## Contributing

We use [`dep`](https://github.com/golang/dep) for dependency management. If you 
do not have it installed, see the [installation instructions](https://github.com/golang/dep#installation).

To contribute: fork and clone the repository, `cd` into the directory, and run:

```shell
dep ensure
```

That will ensure you have [`google/go-cmp`](https://github.com/google/go-cmp) which is used to run tests.

If you plan on submitting a patch, please write tests for it.