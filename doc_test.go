package recurly_test

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/blacklightcms/recurly"
)

var client *recurly.Client

func init() {
	client = recurly.NewClient("your-subdomain", "APIKEY")
}

func ExampleNewClient() {
	// Create a client pointing to https://your-subdomain.recurly.com
	client := recurly.NewClient("your-subdomain", "APIKEY")

	// Optionally overwrite the underlying HTTP client with your own
	client.Client = &http.Client{
		Timeout:   5 * time.Second,
		Transport: &http.Transport{},
	}
}

func Example_AccountsService_Create() {
	_, err := client.Accounts.Create(context.Background(), recurly.Account{
		Code:      "1",
		FirstName: "Verena",
		LastName:  "Example",
		Email:     "verena@example.com",
	})
	if err != nil {
		panic(err)
	}
}

func ExampleNullBool() {
	b := recurly.NewBool(true)
	fmt.Println(b.Bool())

	value, ok := b.Value()
	fmt.Println(value, ok)

	// Output:
	// true
	// true true
}

func ExampleNullBool_invalid() {
	var b recurly.NullBool
	fmt.Println(b.Bool())

	value, ok := b.Value()
	fmt.Println(value, ok)

	// Output:
	// false
	// false false
}

func ExampleNullInt() {
	i := recurly.NewInt(100)
	fmt.Println(i.Int())

	value, ok := i.Value()
	fmt.Println(value, ok)

	// Output:
	// 100
	// 100 true
}

func ExampleNullInt_invalid() {
	var i recurly.NullInt
	fmt.Println(i.Int())

	value, ok := i.Value()
	fmt.Println(value, ok)

	// Output:
	// 0
	// 0 false
}

func ExampleNullTime() {
	t := recurly.NewTime(time.Date(2018, 5, 13, 0, 0, 0, 0, time.UTC))
	fmt.Println(t.Time())

	value, ok := t.Value()
	fmt.Println(value, ok)

	// Output:
	// 2018-05-13 00:00:00 +0000 UTC
	// 2018-05-13 00:00:00 +0000 UTC true
}

func ExampleNullTime_invalid() {
	var t recurly.NullTime
	fmt.Println(t.Time())

	value, ok := t.Value()
	fmt.Println(value, ok)

	// Output:
	// 0001-01-01 00:00:00 +0000 UTC
	// 0001-01-01 00:00:00 +0000 UTC false
}
