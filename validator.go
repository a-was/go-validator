package validator

import (
	"fmt"
	"reflect"
)

type validatorElem struct {
	tag       string
	validator validatorI
}

var validators = []validatorElem{
	{tag: "env", validator: envV{}},
	{tag: "default", validator: defaultV{}},
	{tag: "flags", validator: flagsV{}},
	{tag: "min", validator: minV{}},
	{tag: "max", validator: maxV{}},
	{tag: "regex", validator: regexV{}},
}

type validator struct {
	errors Errors
}

func Validate[T any](strct *T) error {
	v := validator{}
	return v.validateStruct("", strct)
}

// strct needs to be pointer to struct
func (v *validator) validateStruct(baseFieldPath string, strct any) error {
	rv := reflect.ValueOf(strct).Elem()
	// needs to resolve interface to call FieldByName
	if rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}
	// needs to resolve pointer to call FieldByName
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	rt := rv.Type()
	// rt has to be kind struct
	if rt.Kind() != reflect.Struct {
		panic(fmt.Sprintf("validator: cannot validate type %s", rt.Kind()))
	}

	defaultField, ok := rt.FieldByName("_")
	var defaultTags map[string]string // tag name -> tag value
	if ok {
		defaultTags = map[string]string{}
		for _, elem := range validators {
			if val, ok := defaultField.Tag.Lookup(elem.tag); ok {
				defaultTags[elem.tag] = val
			}
		}
	}

	for i := 0; i < rt.NumField(); i++ {
		fieldT := rt.Field(i)
		if fieldT.Name == "_" {
			continue
		}

		fieldV := rv.Field(i)
		resolved := fieldV
		if resolved.Kind() == reflect.Pointer {
			resolved = resolved.Elem()
		}
		if resolved.Kind() == reflect.Struct {
			v.validateStruct(
				fmt.Sprintf("%s%s.", baseFieldPath, fieldT.Name),
				resolved.Addr().Interface(),
			)
		}

		// validate field
		for _, elem := range validators {
			if _, ok := defaultTags[elem.tag]; ok {
				// will be validated by default tags, no need to do it now
				continue
			}

			val, ok := fieldT.Tag.Lookup(elem.tag)
			if !ok {
				continue
			}
			if err := elem.validator.validate(val, fieldV); err != nil {
				v.errors = append(v.errors, fmt.Errorf("%s%s: %s", baseFieldPath, fieldT.Name, err.Error()))
			}
		}
		// validate default tags
		for _, elem := range validators {
			val, ok := defaultTags[elem.tag]
			if !ok {
				continue
			}
			if err := elem.validator.validate(val, fieldV); err != nil {
				v.errors = append(v.errors, fmt.Errorf("%s%s: %s", baseFieldPath, fieldT.Name, err.Error()))
			}
		}
	}

	if len(v.errors) > 0 {
		return v.errors
	}
	return nil
}
