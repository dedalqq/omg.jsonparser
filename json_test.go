package jsonparser

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestParseJsonBody(t *testing.T) {
	testData := []struct {
		json   string
		st     interface{}
		err    error
		result string
	}{
		{
			`{"a": "b"}`,
			&struct {
				A string `json:"a"`
			}{},
			nil,
			`{"a":"b"}`,
		},
		{
			`{"a": "b"}`,
			&struct {
				A string `json:"a"`
				B string `json:"b,required"`
			}{},
			fmt.Errorf("value [b] is required"),
			`{"a":"b","b":""}`,
		},
		{
			`{"a": "b", "b": null}`,
			&struct {
				A string `json:"a"`
				B string `json:"b,required"`
			}{},
			nil,
			`{"a":"b","b":""}`,
		},
		{
			`{"a": "b", "b": null}`,
			&struct {
				A string `json:"a"`
				B string `json:"b,required,notEmpty"`
			}{},
			fmt.Errorf("value [b] must be not empty"),
			`{"a":"b","b":""}`,
		},
		{
			`{"a": "b", "b": ""}`,
			&struct {
				A string `json:"a"`
				B string `json:"b,required,notEmpty"`
			}{},
			fmt.Errorf("value [b] must be not empty"),
			`{"a":"b","b":""}`,
		},
		{
			`{"a": "b", "b": "c"}`,
			&struct {
				A string `json:"a"`
				B string `json:"b,required,notEmpty"`
			}{},
			nil,
			`{"a":"b","b":"c"}`,
		},
		{
			`{"a": "b", "b": {}}`,
			&struct {
				A string `json:"a"`
				B struct {
					C string `json:"c,required"`
				} `json:"b"`
			}{},
			fmt.Errorf("value [b.c] is required"),
			`{"a":"b","b":{"c":""}}`,
		},
		{
			`{"a": "b", "b": {}}`,
			&struct {
				A string `json:"a"`
				B *struct {
					C string `json:"c,required"`
				} `json:"b"`
			}{},
			fmt.Errorf("value [b.c] is required"),
			`{"a":"b","b":{"c":""}}`,
		},
		{
			`{"a": "b", "b": {"c": "d"}}`,
			&struct {
				A string `json:"a"`
				B *struct {
					C string `json:"c,required"`
				} `json:"b"`
			}{},
			nil,
			`{"a":"b","b":{"c":"d"}}`,
		},
		{
			`{"a": ["b", "c"]}`,
			&struct {
				A []string `json:"a"`
			}{},
			nil,
			`{"a":["b","c"]}`,
		},
		{
			`{"a": ["b", "c"]}`,
			&struct {
				A []*string `json:"a"`
			}{},
			nil,
			`{"a":["b","c"]}`,
		},
		{
			`{"a": [{"b": 1}, {"b": 2}]}`,
			&struct {
				A []struct {
					B int `json:"b"`
				} `json:"a"`
			}{},
			nil,
			`{"a":[{"b":1},{"b":2}]}`,
		},
		{
			`{"a": [{"b": 1}, {}]}`,
			&struct {
				A []struct {
					B int `json:"b,required"`
				} `json:"a"`
			}{},
			fmt.Errorf("value [a.[1].b] is required"),
			`{"a":[{"b":1},{"b":0}]}`,
		},
		{
			`{"a": [{"b": 1}, {}]}`,
			&struct {
				A []struct {
					B int `json:"b,required"`
				} `json:"a,min:3"`
			}{},
			fmt.Errorf("value [a.] count of items less than expected"),
			`{"a":null}`,
		},
		{
			`{"a": [{"b": 1}, {}]}`,
			&struct {
				A []struct {
					B int `json:"b,required"`
				} `json:"a,max:1"`
			}{},
			fmt.Errorf("value [a.] count of items more than expected"),
			`{"a":null}`,
		},
	}

	for i, td := range testData {
		t.Run(fmt.Sprintf("Case %d", i+1), func(t *testing.T) {
			err := NewDecoder(strings.NewReader(td.json)).Decode(td.st)

			if err == nil {
				if td.err != nil {
					t.Fail()
				}
			} else {
				if td.err == nil || err.Error() != td.err.Error() {
					fmt.Println(err.Error())
					t.Fail()
				}
			}

			d, err := json.Marshal(td.st)
			if err != nil {
				t.Fail()
			}

			if string(d) != td.result {
				fmt.Println(fmt.Sprintf("%v", string(d)))
				t.Fail()
			}
		})
	}
}
