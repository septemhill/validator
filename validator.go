package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type dataType int

const validateTag = "validator"

const (
	primitiveType dataType = iota
	pointerType
	structType
	unsupportType
)

var kvStack = make([][]string, 0)
var cnt = int64(0)

func top() []string {
	return kvStack[len(kvStack)-1]
}

func pop() []string {
	l := len(kvStack)
	v := kvStack[l-1]
	kvStack = kvStack[0 : l-1]
	return v
}

func push(kv []string) {
	kvStack = append(kvStack, kv)
}

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

func primitiveValidate(value reflect.Value) bool {
	switch value.Interface().(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return integerValidator(value.Int(), top())
	case float32, float64:
		return floatValidator(value.Float(), top())
	case string:
		return stringValidator(value.String(), top())
	default:
		return false
	}
}

func primitiveTypeCheck(value reflect.Value) dataType {
	switch value.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.String:
		return primitiveType
	case reflect.Ptr:
		return pointerType
	case reflect.Struct:
		return structType
	default:
		return unsupportType
	}
}

func structValidate(v interface{}) bool {
	dt := reflect.TypeOf(v)
	dv := reflect.ValueOf(v)
	flag := true

	for i := 0; i < dt.NumField(); i++ {
		tags := dt.Field(i).Tag.Get(validateTag)
		push(strings.Split(tags, ","))

		switch dt.Field(i).Type.Kind() {
		case reflect.Array, reflect.Slice:
			if dv.Field(i).Len() > 0 {
				t := primitiveTypeCheck(dv.Field(i).Index(0))
				if t == primitiveType {
					for j := 0; j < dv.Field(i).Len(); j++ {
						if !primitiveValidate(dv.Field(i).Index(j)) {
							return false
						}
					}
				} else if t == pointerType {
					for j := 0; j < dv.Field(i).Len(); j++ {
						if !primitiveValidate(dv.Field(i).Index(j).Elem()) {
							return false
						}
					}
				} else if t == structType {
					for j := 0; j < dv.Field(i).Len(); j++ {
						if !structValidate(dv.Field(i).Index(j).Interface()) {
							return false
						}
					}
				} else {
					fmt.Println("KerKer")
				}
			}
		case reflect.Struct:
			flag = structValidate(dv.Field(i).Interface())
		case reflect.Ptr:
			flag = structValidate(dv.Field(i).Elem().Interface())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64, reflect.String:
			flag = primitiveValidate(dv.Field(i))
		default:
		}

		pop()

		if !flag {
			return false
		}
	}

	return true
}

func Validate(v interface{}) bool {
	dt := reflect.TypeOf(v)
	dv := reflect.ValueOf(v)

	switch dt.Kind() {
	case reflect.Struct:
		return structValidate(v)
	case reflect.Ptr:
		return Validate(dv.Elem().Interface())
	default:
		fmt.Println("unexpected data type")
		return false
	}
}
