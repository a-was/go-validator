package validator

import "reflect"

func isIntKind(k reflect.Kind) bool {
	switch k {
	case
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return true
	}
	return false
}

func isUintKind(k reflect.Kind) bool {
	switch k {
	case
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr:
		return true
	}
	return false
}

func isFloatKind(k reflect.Kind) bool {
	switch k {
	case
		reflect.Float32,
		reflect.Float64:
		return true
	}
	return false
}

func isLenKind(k reflect.Kind) bool {
	switch k {
	case
		reflect.Array,
		reflect.Slice,
		reflect.Map,
		reflect.String:
		return true
	}
	return false
}
