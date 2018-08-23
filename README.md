# Recurly
Recurly is a Go (golang) API Client for the [Recurly](https://recurly.com/) API.

 [![Build Status](https://travis-ci.org/launchpadcentral/recurly.svg?branch=master)](https://travis-ci.org/launchpadcentral/recurly)  [![GoDoc](https://godoc.org/github.com/launchpadcentral/recurly?status.svg)](https://godoc.org/github.com/launchpadcentral/recurly)

## Announcement: Credit Invoice Release 
As of v2.10, Recurly has released data structure and API changes to support the [credit invoice release](https://docs.recurly.com/docs/credit-invoices-release). All sites created after May 8, 2018, have this turned on automatically. Existing sites must turn it on no later than Nov 1, 2018. New features will be applicable once the feature is turned on; legacy (existing) invoices will continue to use some deprecated features until they are closed. Most new code can coexist with existing, allowing you to write code to support your transition smoothly. Note that the `new_dunning_event` webhook data structure will change and both may be sent while you have both legacy and new invoices with dunning events. Please review the `Parse` and `ParseDeprecated` methods if you listen for this webhook.  Deprecated code will be removed with any library updates no earlier than December 1, 2018. 

## References
 * [API Reference](http://godoc.org/github.com/launchpadcentral/recurly)
 * [Recurly API Documentation](https://dev.recurly.com/docs/)
 * [recurly.js Documentation](https://docs.recurly.com/js/)
 * Documentation and examples below. Unit tests also provide thorough examples.

## Installation
Install using the "go get" command:
```
go get github.com/launchpadcentral/recurly
```

### Example

```go
import "github.com/launchpadcentral/recurly"
```

Construct a new Recurly Client and then work off of that. For example, to list
accounts:
```go
client := recurly.NewClient("subdomain", "apiKey", nil)
resp, accounts, err := client.Accounts.List(recurly.Params{"per_page": 20})
```

recurly.Response embeds http.Response and provides some convenience methods:
```go
if resp.IsOK() {
    fmt.Println("Response was a 200-299 status code")
} else if resp.IsError() {
    fmt.Println("Response was NOT a 200-299 status code")

    // Loop through errors (422 status code only)
    for _, e := range resp.Errors {
        fmt.Printf("Message: %s; Field: %s; Symbol: %s\n", e.Message, e.Field, e.Symbol)
    }
}

if resp.IsClientError() {
    fmt.Println("You messed up. Response was a 400-499 status code")
} else if resp.IsServerError() {
    fmt.Println("Try again later. Response was a 500-599 status code")
}

// Get status code from http.response
if resp.StatusCode == 422 {
    // ...
}
```

## Usage
The basic usage format is to create a client, and then operate directly off of each
of the services.

The services are (each link to the GoDoc documentation):
 * [Accounts](https://godoc.org/github.com/launchpadcentral/recurly#AccountsService)
 * [AddOns](https://godoc.org/github.com/launchpadcentral/recurly#AddOnsService)
 * [Adjustments](https://godoc.org/github.com/launchpadcentral/recurly#AdjustmentsService)
 * [Billing](https://godoc.org/github.com/launchpadcentral/recurly#BillingService)
 * [Coupons](https://godoc.org/github.com/launchpadcentral/recurly#CouponsService)
 * [CreditPayments](https://godoc.org/github.com/launchpadcentral/recurly#CreditPaymentsService)
 * [Invoices](https://godoc.org/github.com/launchpadcentral/recurly#InvoicesService)
 * [Plans](https://godoc.org/github.com/launchpadcentral/recurly#PlansService)
 * [Redemptions](https://godoc.org/github.com/launchpadcentral/recurly#RedemptionsService)
 * [Subscriptions](https://godoc.org/github.com/launchpadcentral/recurly#SubscriptionsService)
 * [Transactions](https://godoc.org/github.com/launchpadcentral/recurly#TransactionsService)
 * [Purchases](https://godoc.org/github.com/launchpadcentral/recurly#PurchasesService)

Each of the services correspond to their respective sections in the
[Recurly API Documentation](https://dev.recurly.com/docs/).

Here are a few examples:

### Create Account
```go
resp, a, err := client.Accounts.Create(recurly.Account{
    Code: "1",
    FirstName: "Verena",
    LastName: "Example",
    Email: "verena@example.com",
})

if resp.IsOK() {
    log.Printf("Account successfully created. Hosted Login Token: %s", a.HostedLoginToken)
}
```

### Get Account
```go
resp, a, err := client.Accounts.Get("1")
if resp.IsOK() {
    log.Printf("Account Found: %+v", a)
}
```

### Get Accounts (pagination example)
All paginated methods (usually named List or List*) support a ```per_page``` and ```cursor``` parameter. Example usage:

```go
resp, accounts, err := client.Accounts.List(recurly.Params{"per_page": 10})

if resp.IsError() {
    // Error occurred
}

for i, a := range accounts {
    // Loop through accounts
}

// Check for next page
next := resp.Next()
if next == "" {
    // No next page
}

// Retrieve next page
resp, accounts, err := client.Accounts.List(recurly.Params{
    "per_page": 10,
    "cursor": next,
})

// Check for prev page
prev := resp.Prev()
if prev == "" {
    // No prev page
}

// Retrieve prev page
resp, accounts, err := client.Accounts.List(recurly.Params{
    "per_page": 10,
    "cursor": prev,
})
```

### Close account
```go
resp, err := client.Accounts.Close("1")
```

### Reopen account
```go
resp, err := client.Accounts.Reopen("1")
```

### Create Billing Info Using recurly.js Token
```go
// 1 is the account code
resp, b, err := client.Billing.CreateWithToken("1", token)
```

### Update Billing Info Using recurly.js Token
```go
// 1 is the account code
resp, b, err := client.Billing.UpdateWithToken("1", token)
```

### Create Billing with Credit Card
```go
resp, b, err := client.Billing.Create("1", recurly.Billing{
    FirstName: "Verena",
    LastName:  "Example",
    Address:   "123 Main St.",
    City:      "San Francisco",
    State:     "CA",
    Zip:       "94105",
    Country:   "US",
    Number:    4111111111111111,
    Month:     10,
    Year:      2020,
})
```

### Create Billing With Bank account
```go
resp, b, err := client.Billing.Create("134", recurly.Billing{
    FirstName:     "Verena",
    LastName:      "Example",
    Address:       "123 Main St.",
    City:          "San Francisco",
    State:         "CA",
    Zip:           "94105",
    Country:       "US",
    NameOnAccount: "Acme, Inc",
    RoutingNumber: "123456780",
    AccountNumber: "111111111",
    AccountType:   "checking",
})
```

### Creating Subscriptions
Subscriptions have different formats for creating and reading.
Due to that, they have a special use case when creating -- a ```NewSubscription```
struct respectively. `NewSubscription` structs are only used for creating.

When updating a subscription, you should use the ```UpdateSubscription``` struct.
All other creates/updates throughout use the same struct to create/update as to read.

```go
// s will return a Subscription struct after creating using the
// NewSubscription struct.
resp, s, err := client.Subscriptions.Create(recurly.NewSubscription{
    PlanCode: "gold",
    Currency: "EUR",
    Account: recurly.Account{
        Code:      "b6f5783",
        Email:     "verena@example.com",
        FirstName: "Verena",
        LastName:  "Example",
        BillingInfo: &recurly.Billing{
            Number:            4111111111111111,
            Month:             12,
            Year:              2017,
            VerificationValue: 123,
            Address:           "400 Alabama St",
            City:              "San Francisco",
            State:             "CA",
            Zip:               "94110",
        },
    },
})
```

## Working with Null* Types
This package has a few null types that ensure that zero values will marshal
or unmarshal properly.

For example, booleans have a zero value of ```false``` in Go. If you need to
explicitly send a false value, go will see that as a zero value and the omitempty
option will ensure it doesn't get sent.

Likewise if you attempt to unmarshal empty/nil values into a struct, you will
also get errors. The Null types help ensure things work as expected.

### NullBool
NullBool is a basic struct that looks like this:

```go
NullBool struct {
    Bool  bool
    Valid bool
}
```
The Valid field determines if the boolean value stored in Bool was intentionally
set there, or if it should be discarded since the default will be false.

Here's how to work with NullBool:
```go
// Create a new NullBool:
t := recurly.NewBool(true)

// Check if the value held in the bool is what you expected
fmt.Printf("%v", t.Is(true)) // true
fmt.Printf("%v", t.Is(false)) // false
```

If, however, NullBool looked like this:
```go
recurly.NullBool{
    Bool: false,
    Valid: false,
}
```

Then those checks will always return false:
```go
fmt.Printf("%v", t.Is(true)) // false
fmt.Printf("%v", t.Is(false)) // false
```

### NullInt
NullInt works the same way as NullBool, but for integers.

```go
i := recurly.NewInt(0)
i = recurly.NewInt(1)
i = recurly.NewInt(50)
```

### NullTime
NullTime won't breakdown if an empty string / nil value is returned from the Recurly
API. It also ensures times are always in UTC.

```go
t := time.Now()
nt := recurly.NewTime(t) // time is now in UTC
fmt.Println(t.String()) // 2015-08-03T19:11:33Z
```

You can then use s.Account.Code to retrieve account info, or s.Invoice.Code to
retrieve invoice info.

## Transaction errors
In addition to the Errors property in the recurly.Response, response also
contains a TransactionError field for Transaction Errors.

Be sure to check `resp.TransactionError` for any API calls that may return a transaction
error for additional info. The `TransactionError` struct is defined like this:
```go
TransactionError struct {
	XMLName          xml.Name `xml:"transaction_error"`
	ErrorCode        string   `xml:"error_code,omitempty"`
	ErrorCategory    string   `xml:"error_category,omitempty"`
	MerchantMessage  string   `xml:"merchant_message,omitempty"`
	CustomerMessage  string   `xml:"customer_message,omitempty"`
	GatewayErrorCode string   `xml:"gateway_error_code,omitempty"`
}
```

[Link to transaction error documentation](https://recurly.readme.io/v2.0/page/transaction-errors).

## Using webhooks
Initial webhook support is in place. The following webhooks are supported:

Account Notifications
 - `NewAccountNotification`
 - `UpdatedAccountNotification`
 - `CanceledAccountNotification`
 - `BillingInfoUpdatedNotification`
 - `BillingInfoUpdateFailedNotification`

Subscription Notifications
 - `NewSubscriptionNotification`
 - `UpdatedSubscriptionNotification`
 - `RenewedSubscriptionNotification`
 - `ExpiredSubscriptionNotification`
 - `CanceledSubscriptionNotification`
 - `PausedSubscriptionNotification`
 - `ResumedSubscriptionNotification`
 - `ScheduledSubscriptionPauseNotification`
 - `SubscriptionPauseModifiedNotification`
 - `PausedSubscriptionRenewalNotification`
 - `SubscriptionPauseCanceledNotification`
 - `ReactivatedAccountNotification`

 Invoice Notifications
 - `NewInvoiceNotification`
 - `PastDueInvoiceNotification`
 - `ProcessingInvoiceNotification`
 - `ClosedInvoiceNotification`

Payment Notifications
 - `SuccessfulPaymentNotification`
 - `FailedPaymentNotification`
 - `VoidPaymentNotification`
 - `SuccessfulRefundNotification`
 - `ScheduledPaymentNotification`
 - `ProcessingPaymentNotification`
 
 Dunning Event Notifications
 - `NewDunningEventNotification`
     
Webhooks can be used by passing an `io.Reader` to `webhooks.Parse`, then using a switch statement with type assertions to determine the webhook returned.

PRs are welcome for additional webhooks.

## License
recurly is available under the [MIT License](http://opensource.org/licenses/MIT).
