# omg.jsonParser

[![Build Status](https://travis-ci.com/dedalqq/omg.jsonparser.svg?branch=master)](https://travis-ci.com/dedalqq/omg.jsonparser)
[![Go Reference](https://pkg.go.dev/badge/github.com/dedalqq/omg.jsonparser.svg)](https://pkg.go.dev/github.com/dedalqq/omg.jsonparser)
[![Coverage Status](https://coveralls.io/repos/github/dedalqq/omg.jsonparser/badge.svg?branch=master)](https://coveralls.io/github/dedalqq/omg.jsonparser?branch=master)

omg.jsonParser is a simple JSON parser with a simple condition validation. It's a wrapper on standard go JSON lib.

## Example

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

## Available tags

| Name        | Description                                                                        | Types |
| ----------- | ---------------------------------------------------------------------------------- | ----- |
| `required`  | field must be exist with any value or `null`                                       | any   |
| `notEmpty`  | field can be not exist but if exist value must be not zero value but can be `null` | any   |
| `min:n`     | slice must have n items or more                                                    | slice |
| `max:n`     | slice must have n items or less                                                    | slice |