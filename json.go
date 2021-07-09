package jsonparser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

const (
	optRequired = "required"
	optNotEmpty = "notEmpty"
	optMin      = "min:"
	optMax      = "max:"
)

type fieldOpt struct {
	required bool
	notEmpty bool
	min      *int
	max      *int
}

type validError struct {
	path   string
	reason string
}

// Error returns error text
func (e validError) Error() string {
	return fmt.Sprintf("value [%s] %s", e.path, e.reason)
}

func newError(reason string, path ...string) error {
	return validError{
		path:   strings.Join(path, ""),
		reason: reason,
	}
}

func parseTag(data string) (string, fieldOpt) {
	values := strings.Split(data, ",")

	opt := fieldOpt{}

	for _, o := range values[1:] {
		switch o {
		case optRequired:
			opt.required = true
		case optNotEmpty:
			opt.notEmpty = true
		}

		if strings.HasPrefix(o, optMin) {
			if v, err := strconv.Atoi(o[4:]); err == nil {
				opt.min = &v
			}

			continue
		}

		if strings.HasPrefix(o, optMax) {
			if v, err := strconv.Atoi(o[4:]); err == nil {
				opt.max = &v
			}

			continue
		}
	}

	return values[0], opt
}

func fieldIs(t reflect.Type, kind reflect.Kind) bool {
	for {
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
			continue
		}

		if t.Kind() == kind {
			return true
		}

		return false
	}
}

func parseJsonObject(r io.Reader, prefix string, val reflect.Value) error {
	data := map[string]json.RawMessage{}

	err := json.NewDecoder(r).Decode(&data)
	if err != nil {
		return err
	}

	for {
		if val.Kind() == reflect.Ptr {
			if val.Elem().Kind() == reflect.Invalid {
				newVal := reflect.New(val.Type().Elem())
				val.Set(newVal)
				val = newVal
			}

			val = val.Elem()
		} else {
			break
		}
	}

	for i := 0; i < val.NumField(); i++ {
		tag := val.Type().Field(i).Tag.Get("json")

		var (
			name string
			opts fieldOpt
		)

		if len(tag) == 0 {
			name = val.Type().Field(i).Name
		} else {
			name, opts = parseTag(tag)
		}

		jsonValue, ok := data[name]
		if !ok {
			if opts.required {
				return newError("is required", prefix, name)
			}

			continue
		}

		refField := val.Field(i)

		err := parseJson(bytes.NewReader(jsonValue), fmt.Sprintf("%s.", name), opts, refField.Addr())
		if err != nil {
			return err
		}

		if opts.notEmpty && refField.IsZero() {
			return newError("must be not empty", prefix, name)
		}
	}

	return nil
}

func parseJsonSlice(r io.Reader, prefix string, opt fieldOpt, val reflect.Value) error {
	var data []json.RawMessage

	err := json.NewDecoder(r).Decode(&data)
	if err != nil {
		return err
	}

	if opt.min != nil && len(data) < *opt.min {
		return newError("count of items less than expected", prefix)
	}

	if opt.max != nil && len(data) > *opt.max {
		return newError("count of items more than expected", prefix)
	}

	for {
		if val.Kind() != reflect.Ptr {
			break
		}

		val = val.Elem()
	}

	for i, d := range data {
		newVal := reflect.New(val.Type().Elem())
		val.Set(reflect.Append(val, newVal.Elem()))

		err := parseJson(bytes.NewReader(d), fmt.Sprintf("%s[%d].", prefix, i), opt, val.Index(i).Addr())
		if err != nil {
			return err
		}
	}

	return nil
}

func parseJson(r io.Reader, prefix string, opt fieldOpt, val reflect.Value) error {
	if fieldIs(val.Type(), reflect.Struct) {
		return parseJsonObject(r, prefix, val)
	}

	if fieldIs(val.Type(), reflect.Slice) {
		return parseJsonSlice(r, prefix, opt, val)

	}

	return json.NewDecoder(r).Decode(val.Interface())
}

type Decoder struct {
	r io.Reader
}

// NewDecoder created and return new decoder
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// Decode run parsing and validation JSON from reader
func (dec *Decoder) Decode(v interface{}) error {
	return parseJson(dec.r, "", fieldOpt{}, reflect.ValueOf(v))
}

// Unmarshal run parsing and validation JSON
func Unmarshal(data []byte, v interface{}) error {
	return NewDecoder(bytes.NewBuffer(data)).Decode(v)
}
