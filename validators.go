package validator

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type validatorI interface {
	validate(tagValue string, v reflect.Value) error
}

func ptr[T any](v T) *T {
	return &v
}

type minV struct{}

func (minV) validate(tagValue string, v reflect.Value) error {
	k := v.Kind()
	if k == reflect.Pointer {
		k = v.Elem().Kind()
		v = v.Elem()
	}
	switch {
	case isIntKind(k):
		req, _ := strconv.ParseInt(tagValue, 0, 64)
		val := v.Int()
		if req > val {
			return fmt.Errorf("invalid value: %d, minimum value is %s", val, tagValue)
		}

	case isUintKind(k):
		req, _ := strconv.ParseUint(tagValue, 0, 64)
		val := v.Uint()
		if req > val {
			return fmt.Errorf("invalid value: %d, minimum value is %s", val, tagValue)
		}

	case isFloatKind(k):
		req, _ := strconv.ParseFloat(tagValue, 64)
		val := v.Float()
		if req > val {
			return fmt.Errorf("invalid value: %f, minimum value is %s", val, tagValue)
		}

	case isLenKind(k):
		req, _ := strconv.Atoi(tagValue)
		val := v.Len()
		if int(req) > val {
			return fmt.Errorf("invalid length: %d, minimum value is %s", val, tagValue)
		}

	case !v.IsValid():
		return nil

	default:
		return fmt.Errorf("invalid type for min validator")
	}

	return nil
}

type maxV struct{}

func (maxV) validate(tagValue string, v reflect.Value) error {
	k := v.Kind()
	if k == reflect.Pointer {
		k = v.Elem().Kind()
		v = v.Elem()
	}
	switch {
	case isIntKind(k):
		req, _ := strconv.ParseInt(tagValue, 0, 64)
		val := v.Int()
		if val > req {
			return fmt.Errorf("invalid value: %d, maximum value is %s", val, tagValue)
		}

	case isUintKind(k):
		req, _ := strconv.ParseUint(tagValue, 0, 64)
		val := v.Uint()
		if val > req {
			return fmt.Errorf("invalid value: %d, maximum value is %s", val, tagValue)
		}

	case isFloatKind(k):
		req, _ := strconv.ParseFloat(tagValue, 64)
		val := v.Float()
		if val > req {
			return fmt.Errorf("invalid value: %f, maximum value is %s", val, tagValue)
		}

	case isLenKind(k):
		req, _ := strconv.Atoi(tagValue)
		val := v.Len()
		if val > int(req) {
			return fmt.Errorf("invalid length: %d, maximum value is %s", val, tagValue)
		}

	case !v.IsValid():
		return nil

	default:
		return fmt.Errorf("invalid type for max validator")
	}

	return nil
}

type regexV struct{}

func (regexV) validate(tagValue string, v reflect.Value) error {
	if !v.IsValid() {
		return nil
	}

	if !strings.HasPrefix(tagValue, "^") && !strings.HasSuffix(tagValue, "$") {
		tagValue = "^" + tagValue + "$"
	}

	re, err := regexp.Compile(tagValue)
	if err != nil {
		return fmt.Errorf("invalid regex: %s", tagValue)
	}

	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			// return fmt.Errorf("invalid value: nil")
			return nil
		}
		v = v.Elem()
	}
	val := v.String()
	// if len(val) == 0 {
	// 	return fmt.Errorf("invalid value: empty string")
	// }
	if !re.MatchString(val) {
		return fmt.Errorf("invalid value: %s does not match regex %s", val, tagValue)
	}

	return nil
}

type flagsV struct{}

func (flagsV) validate(tagValue string, v reflect.Value) error {
	for _, flag := range strings.Split(tagValue, ",") {
		switch strings.TrimSpace(flag) {
		case "required":

			if v.Kind() == reflect.Pointer {
				if v.IsNil() {
					return fmt.Errorf("required value not filled")
				}
				v = v.Elem()
			}
			if !v.IsValid() || v.IsZero() {
				return fmt.Errorf("required value not filled")
			}
		}
	}

	return nil
}

type envV struct{}

func (envV) validate(tagValue string, v reflect.Value) error {
	if !v.CanSet() {
		panic("validator: Using `env` validator with unsettable value")
	}

	if !v.IsValid() {
		return nil
	}

	fromEnv := os.Getenv(tagValue)
	if fromEnv == "" {
		return nil
	}

	var pointer bool
	k := v.Kind()
	if k == reflect.Pointer {
		k = v.Type().Elem().Kind()
		pointer = true
	}
	switch {
	case isIntKind(k):
		req, err := strconv.ParseInt(fromEnv, 0, 64)
		if err != nil {
			return fmt.Errorf("invalid `env` value %s", fromEnv)
		}
		if pointer {
			v.Set(reflect.ValueOf(ptr(req)))
		} else {
			v.SetInt(req)
		}

	case isUintKind(k):
		req, err := strconv.ParseUint(fromEnv, 0, 64)
		if err != nil {
			return fmt.Errorf("invalid `env` value %s", fromEnv)
		}
		if pointer {
			v.Set(reflect.ValueOf(ptr(req)))
		} else {
			v.SetUint(req)
		}

	case isFloatKind(k):
		req, err := strconv.ParseFloat(fromEnv, 64)
		if err != nil {
			return fmt.Errorf("invalid `env` value %s", fromEnv)
		}
		if pointer {
			v.Set(reflect.ValueOf(ptr(req)))
		} else {
			v.SetFloat(req)
		}

	case k == reflect.String:
		if pointer {
			v.Set(reflect.ValueOf(ptr(fromEnv)))
		} else {
			v.SetString(fromEnv)
		}

	case k == reflect.Bool:
		req, err := strconv.ParseBool(fromEnv)
		if err != nil {
			return fmt.Errorf("invalid `env` value %s", fromEnv)
		}
		if pointer {
			v.Set(reflect.ValueOf(ptr(req)))
		} else {
			v.SetBool(req)
		}

	default:
		return fmt.Errorf("invalid type for `env` validator")
	}

	return nil
}

type defaultV struct{}

func (defaultV) validate(tagValue string, v reflect.Value) error {
	if !v.CanSet() {
		panic("validator: Using `default` validator with unsettable value")
	}

	if !v.IsValid() {
		return nil
	}

	if tagValue == "" {
		return nil
	}

	var pointer bool
	k := v.Kind()
	if k == reflect.Pointer {
		k = v.Type().Elem().Kind()
		pointer = true
	}
	switch {
	case isIntKind(k):
		req, err := strconv.ParseInt(tagValue, 0, 64)
		if err != nil {
			return fmt.Errorf("invalid `default` value %s", tagValue)
		}
		if pointer {
			v.Set(reflect.ValueOf(ptr(req)))
		} else {
			v.SetInt(req)
		}

	case isUintKind(k):
		req, err := strconv.ParseUint(tagValue, 0, 64)
		if err != nil {
			return fmt.Errorf("invalid `default` value %s", tagValue)
		}
		if pointer {
			v.Set(reflect.ValueOf(ptr(req)))
		} else {
			v.SetUint(req)
		}

	case isFloatKind(k):
		req, err := strconv.ParseFloat(tagValue, 64)
		if err != nil {
			return fmt.Errorf("invalid `default` value %s", tagValue)
		}
		if pointer {
			v.Set(reflect.ValueOf(ptr(req)))
		} else {
			v.SetFloat(req)
		}

	case k == reflect.String:
		if pointer {
			v.Set(reflect.ValueOf(ptr(tagValue)))
		} else {
			v.SetString(tagValue)
		}

	case k == reflect.Bool:
		req, err := strconv.ParseBool(tagValue)
		if err != nil {
			return fmt.Errorf("invalid `default` value %s", tagValue)
		}
		if pointer {
			v.Set(reflect.ValueOf(ptr(req)))
		} else {
			v.SetBool(req)
		}

	default:
		return fmt.Errorf("invalid type for `default` validator")
	}

	return nil
}
