package jsonparser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	optRequired = "required"
	optNotEmpty = "notEmpty"
	optNotNull  = "notNull"
	optMin      = "min:"
	optMax      = "max:"
)

type fieldOpt struct {
	required bool
	notEmpty bool
	notNull  bool
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
		case optNotNull:
			opt.notNull = true
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

func valueIsZero(v reflect.Value) bool {
	for {
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return false
			}

			v = v.Elem()
			continue
		}

		return v.IsZero()
	}
}

func validateMinMax(opt fieldOpt, prefix string, v int, errorPrefix string) error {
	if opt.min != nil && v < *opt.min {
		return newError(fmt.Sprintf("%s less than expected", errorPrefix), prefix)
	}

	if opt.max != nil && v > *opt.max {
		return newError(fmt.Sprintf("%s more than expected", errorPrefix), prefix)
	}

	return nil
}

func getString(v reflect.Value) (string, bool) {
	for {
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
			continue
		}

		if v.Kind() != reflect.String {
			return "", false
		}

		return v.String(), true
	}
}

func getInt(v reflect.Value) (int64, bool) {
	for {
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
			continue
		}

		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return v.Int(), true
		default:
			return 0, false
		}
	}
}

func getUint(v reflect.Value) (uint64, bool) {
	for {
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
			continue
		}

		switch v.Kind() {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return v.Uint(), true
		default:
			return 0, false
		}
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

		jsonValue, valueExist := data[name]
		if !valueExist {
			if opts.required {
				return newError("is required", prefix, name)
			}

			continue
		}

		if opts.notNull && string(jsonValue) == "null" {
			return newError("must be not null", prefix, name)
		}

		refField := val.Field(i)

		err := parseJson(bytes.NewReader(jsonValue), fmt.Sprintf("%s%s.", prefix, name), opts, refField.Addr())
		if err != nil {
			return err
		}

		if opts.notEmpty && valueIsZero(refField) && string(jsonValue) != "null" {
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

	err = validateMinMax(opt, prefix, len(data), "count of items")
	if err != nil {
		return err
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

func parseJsonValue(r io.Reader, prefix string, opt fieldOpt, val reflect.Value) error {
	tempVal := reflect.New(val.Type().Elem())

	err := json.NewDecoder(r).Decode(tempVal.Interface())
	if err != nil {
		return newError(err.Error(), prefix)
	}

	if str, ok := getString(tempVal); ok {
		err := validateMinMax(opt, prefix, utf8.RuneCountInString(str), "count of runes in a string")
		if err != nil {
			return err
		}
	} else if i, ok := getInt(tempVal); ok {
		err := validateMinMax(opt, prefix, (int)(i), "value")
		if err != nil {
			return err
		}
	} else if ui, ok := getUint(tempVal); ok {
		err := validateMinMax(opt, prefix, (int)(ui), "value")
		if err != nil {
			return err
		}
	}

	val.Elem().Set(tempVal.Elem())

	return nil
}

func parseJson(r io.Reader, prefix string, opt fieldOpt, val reflect.Value) error {
	if fieldIs(val.Type(), reflect.Struct) {
		return parseJsonObject(r, prefix, val)
	}

	if fieldIs(val.Type(), reflect.Slice) {
		return parseJsonSlice(r, prefix, opt, val)
	}

	return parseJsonValue(r, prefix, opt, val)
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
