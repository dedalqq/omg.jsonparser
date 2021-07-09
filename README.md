# omg.jsonParser

[![Build Status](https://travis-ci.com/dedalqq/omg.jsonparser.svg?branch=master)](https://travis-ci.com/dedalqq/omg.jsonparser)
[![Go Reference](https://pkg.go.dev/badge/github.com/dedalqq/omg.jsonparser.svg)](https://pkg.go.dev/github.com/dedalqq/omg.jsonparser)

omg.jsonParser is a simple JSON parser with a simple condition validation. It's a wrapper on standard go JSON lib.

## example

```go
package main

import (
	"fmt"

	"github.com/dedalqq/omg.jsonparser"
)

func main() {
	jsonData := `{"name": ""}`

	st := struct {
		Name string `json:"name,notEmpty"` // added notEmpty for enable validation for it field
	}{}

	err := jsonparser.Unmarshal([]byte(jsonData), &st)
	if err != nil {
		fmt.Println(err.Error()) // print: value [name] must be not empty
	}
}

```