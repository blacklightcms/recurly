# Go Recurly
Recurly is a Go (golang) API Client for the [Recurly](https://recurly.com/) API.

 [![Build Status](https://travis-ci.org/blacklightcms/go-recurly.svg?branch=master)](https://travis-ci.org/blacklightcms/go-recurly)  [![GoDoc](https://godoc.org/github.com/blacklightcms/go-recurly/recurly?status.svg)](https://godoc.org/github.com/blacklightcms/go-recurly/recurly)  

## References
 * [API Reference](http://godoc.org/github.com/blacklightcms/go-recurly/recurly)
 * [Recurly API Documentation](https://dev.recurly.com/docs/)
 * [recurly.js Documentation](https://docs.recurly.com/js/)
 * Documentation and examples for the library are coming soon. In the meantime,
 checkout the unit tests for thorough examples and usage.

## Installation
Install using the "go get" command:
```
go get github.com/blacklightcms/go-recurly/recurly
```

### Example

```go
import "github.com/blacklightcms/go-recurly/recurly"
```

Construct a new Recurly Client and then work off of that. For example, to list
accounts:
```go
client, err := recurly.NewClient("subdomain", "apiKey", nil)
resp, accounts, err := client.Accounts.List({"per_page": 20})
```

The recurly.Response class provides some convenience methods:
```go

if resp.IsOK() {
    fmt.Println("Response was a 200-299 status code")
}

if resp.IsError() {
    fmt.Println("Response was NOT a 200-299 status code")

    // Loop through errors (422 status code only)
    for _, e := range resp.Errors() {
        fmt.Printf("Message: %s; Field: %s; Symbol: %s\n", e.Message, e.Field, e.Symbol)
    }
}

if resp.IsClientError() {
    fmt.Println("You messed up. Response was a 400-499 status code")
}

if resp.IsServerError() {
    fmt.Println("Try again later. Response was a 500-599 status code")
}
```

## Roadmap
The API is still being finalized and may change over the coming weeks. Here is
what's coming before things stabilize:
 * Support for paginating beyond the first page with cursors needs to be completed
 * ~~Coupons, coupon redemptions, invoices, and transactions. All other
 portions of the API are complete.~~
 * ~~Documentation~~ and more usage examples.
 * There is currently no support for updating billing info with a credit card or
 bank account directly. Using [recurly.js](https://docs.recurly.com/js/) token is the only supported method currently.
 Because the the token method using [recurly.js](https://docs.recurly.com/js/) is the recommended method, this
 is currently a low priority. The placeholder functions are already in place so
 this will not affect API stability of the library.
 * [Webhook](https://dev.recurly.com/page/webhooks) support. This will come last after API stability.

Once that notice is removed things will be stable. Contributions are welcome.

## License
go-recurly is available under the [MIT License](http://opensource.org/licenses/MIT).
