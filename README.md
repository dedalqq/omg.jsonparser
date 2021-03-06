# omg.jsonParser

[![Go](https://github.com/dedalqq/omg.jsonparser/actions/workflows/go.yml/badge.svg)](https://github.com/dedalqq/omg.jsonparser/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/dedalqq/omg.jsonparser.svg)](https://pkg.go.dev/github.com/dedalqq/omg.jsonparser)
[![Coverage Status](https://coveralls.io/repos/github/dedalqq/omg.jsonparser/badge.svg?branch=master)](https://coveralls.io/github/dedalqq/omg.jsonparser?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/dedalqq/omg.jsonparser)](https://goreportcard.com/report/github.com/dedalqq/omg.jsonparser)

omg.jsonParser is a simple JSON parser with a simple condition validation. It's a wrapper on standard go JSON lib. With it help you can add the validation condition via golang structure fields tags.

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

## Tag struct

```
json:"[name][,option]..."
```

### Tag examples

Here are some the tag examples:

```go
package main

type data struct {
    F1 string   `json:"f1,notEmpty"`          // the field must be not empty but can be "null" or may not exist
    F2 string   `json:"f2,notEmpty,required"` // the field is required and must be not empty but may be the "null" value
    F3 string   `json:"f3,notEmpty,notNull"`  // the field must be not empty and not "null" but may not exist
    F4 []string `json:"f4,notNull,min:3"`     // the field must be not "null" and contains 3 or more items but may not exist
}
```

## Available tags options

| Name        | Types    | Description                                                                            |
| ----------- | -------- | -------------------------------------------------------------------------------------- |
| `required`  | any      | The field must be exist with any value or `null`                                       |
| `notEmpty`  | any      | The field can be not exist but if exist value must be not zero value but can be `null` |
| `notEmpty`  | slice    | The field must have one element or more but may be `null` or not exist                 |
| `notNull`   | any      | The field should not be null, but may not exist                                        |
| `uniq`      | []string | The strings slice must contains only unique strings                                    |
| `min:n`     | slice    | The slice must have `n` items or more                                                  |
| `max:n`     | slice    | The slice must have `n` items or less                                                  |
| `min:n`     | string   | The string must have `n` runes or more                                                 |
| `max:n`     | string   | The string must have `n` runes or less                                                 |
| `min:n`     | int/uint | The value must be equals `n` or more                                                   |
| `max:n`     | int/uint | The value must be equals `n` or less                                                   |
