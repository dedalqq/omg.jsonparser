# omg.jsonParser

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
		Name string `json:"name,notEmpty"`
	}{}

	err := jsonparser.Unmarshal([]byte(jsonData), &st)
	if err != nil {
		fmt.Println(err.Error()) // value [name] must be not empty
	}
}

```