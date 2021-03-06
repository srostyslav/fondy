# Fondy package for Go language

This package provides Go implementation of Fondy api.

## Installation

Use the `go` command:

	$ go get github.com/srostyslav/fondy-go

## Requirements

Fondy package tested against Go 1.13.

Support only 2.0 protocol & json format

## Example

```go
package main

import (
    "github.com/srostyslav/fondy"
    "github.com/satori/go.uuid"
    "fmt"
)

func main() {
    api := fondy.NewApi(&fondy.ApiOptions{MerchantID: 1396424, SecretKey: "test"})
    
    data := &fondy.Checkout{
		Amount: 100,
		Currency: "USD",
                OrderDesc: "Pay for order",
		OrderID: uuid.NewV4().String(),
	}
	if url, err := api.CheckoutUrl(data); err != nil {
		panic(err)
	} else {
		fmt.Println("go to pay: " + url)
	}
}

```
