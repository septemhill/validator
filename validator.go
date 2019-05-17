package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const tagName = "validator"

func stringValidator(val string, kv []string) bool {
	for n := range kv {
		pair := strings.Split(kv[n], ":")
		switch pair[0] {
		case "min":
			if min, err := strconv.ParseInt(pair[1], 10, 64); err != nil || int(min) > len(val) {
				return false
			}
		case "max":
			if max, err := strconv.ParseInt(pair[1], 10, 64); err != nil || int(max) < len(val) {
				return false
			}
		case "regex":
			re := regexp.MustCompile(pair[1])
			if !re.Match([]byte(val)) {
				return false
			}
		}
	}

	return true
}

func integerValidator(val int64, kv []string) bool {
	for n := range kv {
		pair := strings.Split(kv[n], ":")
		switch pair[0] {
		case "min":
			if min, err := strconv.ParseInt(pair[1], 10, 64); err != nil || min > val {
				return false
			}
		case "max":
			if max, err := strconv.ParseInt(pair[1], 10, 64); err != nil || max < val {
				return false
			}
		}
	}

	return true
}

func floatValidator(val float64, kv []string) bool {
	for n := range kv {
		pair := strings.Split(kv[n], ":")
		switch pair[0] {
		case "min":
			if min, err := strconv.ParseFloat(pair[1], 64); err != nil || min > val {
				return false
			}
		case "max":
			if max, err := strconv.ParseFloat(pair[1], 64); err != nil || max < val {
				return false
			}
		}
	}

	return true
}

func validate(value reflect.Value, kv []string) bool {
	switch value.Interface().(type) {
	case int, int8, int16, int32, int64:
		return integerValidator(value.Int(), kv)
	case float32, float64:
		return floatValidator(value.Float(), kv)
	case string:
		return stringValidator(value.String(), kv)
	default:
		return false
	}
}

func fieldsValidate(v interface{}) bool {
	dt := reflect.TypeOf(v)
	dv := reflect.ValueOf(v)

	for i := 0; i < dt.NumField(); i++ {
		switch dt.Field(i).Type.Kind() {
		case reflect.Struct:
			return fieldsValidate(dv.Field(i).Interface())
		case reflect.Array:
		case reflect.Slice:
			data := reflect.ValueOf(dv.Field(i).Interface())
			for j := 0; j < data.Len(); j++ {
				if !fieldsValidate(data.Index(j).Interface()) {
					return false
				}
			}
		case reflect.Ptr:
			return fieldsValidate(dv.Field(i).Elem().Interface())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Float32, reflect.Float64, reflect.String:
			ti := dt.Field(i).Tag.Get(tagName)
			if !validate(dv.Field(i), strings.Split(ti, ",")) {
				return false
			}
		default:
			fmt.Println("Unsupport data type", dt.Field(i).Name)
		}
	}

	return true
}

func Validate(v interface{}) bool {
	dt := reflect.TypeOf(v)
	dv := reflect.ValueOf(v)

	switch dt.Kind() {
	case reflect.Struct:
		return fieldsValidate(v)
	case reflect.Ptr:
		return Validate(dv.Elem().Interface())
	default:
		return false
	}
}
