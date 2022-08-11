package recurly_test

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/autopilot3/recurly"
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

func ExampleNewTestServer() {
	client, s := recurly.NewTestServer()
	defer s.Close()

	// NOTE: This example doesn't have access to *testing.T so it passes nil
	// as the last argument to s.HandlFunc(). When setting up your tests,
	// be sure to pass *testing.T instead of nil.
	s.HandleFunc("GET", "/v2/accounts/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
		<?xml version="1.0" encoding="UTF-8"?>
		<account href="https://your-subdomain.recurly.com/v2/accounts/1">
		   <email>verena@example.com</email>
		   <first_name>Verena</first_name>
		   <last_name>Example</last_name>
		</account>
		`))
	}, nil)

	a, _ := client.Accounts.Get(context.Background(), "1")
	fmt.Printf("%t\n", s.Invoked)
	fmt.Printf("%t\n", a.FirstName == "Verena")
	fmt.Printf("%t\n", a.LastName == "Example")
	fmt.Printf("%t\n", a.Email == "verena@example.com")

	// Output:
	// true
	// true
	// true
	// true
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
