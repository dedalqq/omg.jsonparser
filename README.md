# omg.jsonParser

[![Build Status](https://travis-ci.com/dedalqq/omg.jsonparser.svg?branch=master)](https://travis-ci.com/dedalqq/omg.jsonparser)
[![Go Reference](https://pkg.go.dev/badge/github.com/dedalqq/omg.jsonparser.svg)](https://pkg.go.dev/github.com/dedalqq/omg.jsonparser)
[![Coverage Status](https://coveralls.io/repos/github/dedalqq/omg.jsonparser/badge.svg?branch=master)](https://coveralls.io/github/dedalqq/omg.jsonparser?branch=master)

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
    F1 string   `json:"f1,notEmpty"`          // the field must be not empty
    F2 string   `json:"f2,notEmpty,required"` // the field is required and must be not empty but may be the null value
    F3 string   `json:"f3,notEmpty,notNull"`  // the field must be not empty and not null but may not exist
    F4 []string `json:"f4,notNull,min:3"`     // the field must be not null and contains more 3 items but may not exist
}
```

## Available tags options

| Name        | Types  | Description                                                                            |
| ----------- | ------ | -------------------------------------------------------------------------------------- |
| `required`  | any    | The field must be exist with any value or `null`                                       |
| `notEmpty`  | any    | The field can be not exist but if exist value must be not zero value but can be `null` |
| `notNull`   | any    | The field should not be null, but may not exist                                        |
| `min:n`     | slice  | The slice must have n items or more                                                    |
| `max:n`     | slice  | The slice must have n items or less                                                    |
| `min:n`     | string | The string must have n runes or more                                                   |
| `max:n`     | string | The string must have n runes or less                                                   |