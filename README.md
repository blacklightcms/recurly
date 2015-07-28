# Recurly API Client for Go

## Description
An implementation of the recurly API in golang.

References
 * https://dev.recurly.com/docs/

## Installation
```
go get github.com/blacklightcms/go-recurly/recurly
```

## Usage

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

## Notes
The API is still being finalized and may change over the coming weeks.
 * Support for paginating beyond the first page with cursors needs to be completed
 * Unit tests are nearly complete and will be coming over the coming days.
 * Documentation and more usage examples will be coming as well

Once that notice is removed things will be stable. Contributions are welcome.

## License
Licensed under the MIT.
