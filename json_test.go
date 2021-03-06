package jsonparser

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

type testData struct {
	json   string
	st     interface{}
	err    error
	result string
}

func test(t *testing.T, td ...testData) {
	for i, td := range td {
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

func TestParseJson(t *testing.T) {
	test(
		t,
		testData{
			`{"a": "b"}`,
			&struct {
				A string `json:"a"`
			}{},
			nil,
			`{"a":"b"}`,
		},
		testData{
			`{"A": "b"}`,
			&struct {
				A string
			}{},
			nil,
			`{"A":"b"}`,
		},
		testData{
			`{"a": "b"}`,
			&struct {
				A string `json:"a"`
				B string `json:"b,required"`
			}{},
			fmt.Errorf("value [b] is required"),
			`{"a":"b","b":""}`,
		},
		testData{
			`{"a": "b", "b": null}`,
			&struct {
				A string `json:"a"`
				B string `json:"b,required"`
			}{},
			nil,
			`{"a":"b","b":""}`,
		},
		testData{
			`{"a": "b", "b": null}`,
			&struct {
				A string `json:"a"`
				B string `json:"b,required,notEmpty"`
			}{},
			nil,
			`{"a":"b","b":""}`,
		},
		testData{
			`{"a": "b", "b": ""}`,
			&struct {
				A string `json:"a"`
				B string `json:"b,required,notEmpty"`
			}{},
			fmt.Errorf("value [b] must be not empty"),
			`{"a":"b","b":""}`,
		},
		testData{
			`{"a": "b", "b": "c"}`,
			&struct {
				A string `json:"a"`
				B string `json:"b,required,notEmpty"`
			}{},
			nil,
			`{"a":"b","b":"c"}`,
		},
		testData{
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
		testData{
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
		testData{
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
	)
}

func TestNotEmpty(t *testing.T) {
	test(
		t,
		testData{
			`{"a": "b"}`,
			&struct {
				A string `json:"a,notEmpty"`
			}{},
			nil,
			`{"a":"b"}`,
		},
		testData{
			`{"a": ""}`,
			&struct {
				A string `json:"a,notEmpty"`
			}{},
			fmt.Errorf("value [a] must be not empty"),
			`{"a":""}`,
		},
		testData{
			`{"a": null}`,
			&struct {
				A string `json:"a,notEmpty"`
			}{},
			nil,
			`{"a":""}`,
		},
		testData{
			`{}`,
			&struct {
				A string `json:"a,notEmpty"`
			}{},
			nil,
			`{"a":""}`,
		},
		testData{
			`{"a": "b"}`,
			&struct {
				A *string `json:"a,notEmpty"`
			}{},
			nil,
			`{"a":"b"}`,
		},
		testData{
			`{"a": ""}`,
			&struct {
				A *string `json:"a,notEmpty"`
			}{},
			fmt.Errorf("value [a] must be not empty"),
			`{"a":""}`,
		},
		testData{
			`{"a": null}`,
			&struct {
				A *string `json:"a,notEmpty"`
			}{},
			nil,
			`{"a":null}`,
		},
		testData{
			`{}`,
			&struct {
				A *string `json:"a,notEmpty"`
			}{},
			nil,
			`{"a":null}`,
		},
	)
}

func TestNotNull(t *testing.T) {
	test(
		t,
		testData{
			`{"a": "b"}`,
			&struct {
				A string `json:"a,notNull"`
			}{},
			nil,
			`{"a":"b"}`,
		},
		testData{
			`{"a": ""}`,
			&struct {
				A string `json:"a,notNull"`
			}{},
			nil,
			`{"a":""}`,
		},
		testData{
			`{"a": null}`,
			&struct {
				A string `json:"a,notNull"`
			}{},
			fmt.Errorf("value [a] must be not null"),
			`{"a":""}`,
		},
		testData{
			`{}`,
			&struct {
				A string `json:"a,notNull"`
			}{},
			nil,
			`{"a":""}`,
		},
		testData{
			`{"a": "b"}`,
			&struct {
				A *string `json:"a,notNull"`
			}{},
			nil,
			`{"a":"b"}`,
		},
		testData{
			`{"a": ""}`,
			&struct {
				A *string `json:"a,notNull"`
			}{},
			nil,
			`{"a":""}`,
		},
		testData{
			`{"a": null}`,
			&struct {
				A *string `json:"a,notNull"`
			}{},
			fmt.Errorf("value [a] must be not null"),
			`{"a":null}`,
		},
		testData{
			`{}`,
			&struct {
				A *string `json:"a,notNull"`
			}{},
			nil,
			`{"a":null}`,
		},
	)
}

func TestParsPointer(t *testing.T) {
	test(
		t,
		testData{
			`{"a": "b"}`,
			&struct {
				A *string `json:"a"`
			}{},
			nil,
			`{"a":"b"}`,
		},
		testData{
			`{"a": 123}`,
			&struct {
				A *int `json:"a"`
			}{},
			nil,
			`{"a":123}`,
		},
		testData{
			`{"a": 123}`,
			&struct {
				A *uint `json:"a"`
			}{},
			nil,
			`{"a":123}`,
		},
	)
}

func TestCustomTypes(t *testing.T) {
	type (
		customString string
		customInt    int
	)

	test(
		t,
		testData{
			`{"a": "b"}`,
			&struct {
				A customString `json:"a"`
			}{},
			nil,
			`{"a":"b"}`,
		},
		testData{
			`{"a": 123}`,
			&struct {
				A customInt `json:"a"`
			}{},
			nil,
			`{"a":123}`,
		},
		testData{
			`{"a": "b"}`,
			&struct {
				A *customString `json:"a"`
			}{},
			nil,
			`{"a":"b"}`,
		},
		testData{
			`{"a": 123}`,
			&struct {
				A *customInt `json:"a"`
			}{},
			nil,
			`{"a":123}`,
		},
	)
}

func TestParseSlice(t *testing.T) {
	test(
		t,
		testData{
			`{"a": ["b", "c"]}`,
			&struct {
				A []string `json:"a"`
			}{},
			nil,
			`{"a":["b","c"]}`,
		},
		testData{
			`{"a": ["b", "c"]}`,
			&struct {
				A []*string `json:"a"`
			}{},
			nil,
			`{"a":["b","c"]}`,
		},
		testData{
			`{"a": [{"b": 1}, {"b": 2}]}`,
			&struct {
				A []struct {
					B int `json:"b"`
				} `json:"a"`
			}{},
			nil,
			`{"a":[{"b":1},{"b":2}]}`,
		},
		testData{
			`{"a": [{"b": 1}, {}]}`,
			&struct {
				A []struct {
					B int `json:"b,required"`
				} `json:"a"`
			}{},
			fmt.Errorf("value [a.[1].b] is required"),
			`{"a":[{"b":1},{"b":0}]}`,
		},
		testData{
			`{"a": []}`,
			&struct {
				A []int `json:"a,notEmpty"`
			}{},
			fmt.Errorf("value [a] must be not empty"),
			`{"a":null}`,
		},
		testData{
			`{"a": null}`,
			&struct {
				A []int `json:"a,notEmpty"`
			}{},
			nil,
			`{"a":null}`,
		},
	)
}

func TestMinMax(t *testing.T) {
	test(
		t,
		testData{
			`{"a": [{"b": 1}, {}]}`,
			&struct {
				A []struct {
					B int `json:"b,required"`
				} `json:"a,min:3"`
			}{},
			fmt.Errorf("value [a.] count of items less than expected"),
			`{"a":null}`,
		}, testData{
			`{"a": [{"b": 1}, {}]}`,
			&struct {
				A []struct {
					B int `json:"b,required"`
				} `json:"a,max:1"`
			}{},
			fmt.Errorf("value [a.] count of items more than expected"),
			`{"a":null}`,
		},
		testData{
			`{"a": "test"}`,
			&struct {
				A string `json:"a,max:5"`
			}{},
			nil,
			`{"a":"test"}`,
		},
		testData{
			`{"a": "test"}`,
			&struct {
				A string `json:"a,max:3"`
			}{},
			fmt.Errorf("value [a.] count of runes in a string more than expected"),
			`{"a":""}`,
		},
		testData{
			`{"a": 4}`,
			&struct {
				A int `json:"a,max:5"`
			}{},
			nil,
			`{"a":4}`,
		},
		testData{
			`{"a": 4}`,
			&struct {
				A int `json:"a,max:3"`
			}{},
			fmt.Errorf("value [a.] value more than expected"),
			`{"a":0}`,
		},
		testData{
			`{"a": null}`,
			&struct {
				A int `json:"a,max:5"`
			}{},
			nil,
			`{"a":0}`,
		},
		testData{
			`{"a": null}`,
			&struct {
				A int `json:"a,max:3"`
			}{},
			nil,
			`{"a":0}`,
		},
		testData{
			`{}`,
			&struct {
				A int `json:"a,max:5"`
			}{},
			nil,
			`{"a":0}`,
		},
		testData{
			`{}`,
			&struct {
				A int `json:"a,max:3"`
			}{},
			nil,
			`{"a":0}`,
		},
	)
}

func TestIncorrectType(t *testing.T) {
	test(
		t,
		testData{
			`{"a": {"b": 1}}`,
			&struct {
				A struct {
					B string `json:"b"`
				} `json:"a"`
			}{},
			fmt.Errorf("value [a.b.] json: cannot unmarshal number into Go value of type string"),
			`{"a":{"b":""}}`,
		},
	)
}

func TestUniq(t *testing.T) {
	test(
		t,
		testData{
			`{"a": []}`,
			&struct {
				A []string `json:"a,uniq"`
			}{},
			nil,
			`{"a":null}`,
		},
		testData{
			`{"a": ["b"]}`,
			&struct {
				A []string `json:"a,uniq"`
			}{},
			nil,
			`{"a":["b"]}`,
		},
		testData{
			`{"a": ["b", "c"]}`,
			&struct {
				A []string `json:"a,uniq"`
			}{},
			nil,
			`{"a":["b","c"]}`,
		},
		testData{
			`{"a": ["b", "c", "d"]}`,
			&struct {
				A []string `json:"a,uniq"`
			}{},
			nil,
			`{"a":["b","c","d"]}`,
		},
		testData{
			`{"a": ["a", "a"]}`,
			&struct {
				A []string `json:"a,uniq"`
			}{},
			fmt.Errorf("value [a.] contains repeated values"),
			`{"a":null}`,
		},
		testData{
			`{"a": ["b", "c", "d", "c"]}`,
			&struct {
				A []string `json:"a,uniq"`
			}{},
			fmt.Errorf("value [a.] contains repeated values"),
			`{"a":null}`,
		},
	)
}

func TestUnmarshal(t *testing.T) {
	data := []byte(`{"a": "b"}`)

	st := struct {
		A string `json:"a"`
	}{}

	err := Unmarshal(data, &st)
	if err != nil {
		t.Fail()
	}
}
