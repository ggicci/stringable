package stringable

import (
	"fmt"
	"reflect"
)

func isNil(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.Interface, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}

func panicOnError(err error) {
	if err != nil {
		panic(fmt.Errorf("httpin: %w", err))
	}
}

// typeOf returns the reflect.Type of a given type.
// e.g. typeOf[int]() returns reflect.typeOf(0)
func typeOf[T any]() reflect.Type {
	var zero [0]T
	return reflect.TypeOf(zero).Elem()
}

func pointerize[T any](v T) *T {
	return &v
}

// dereferencedType returns the underlying type of a pointer.
func dereferencedType(v any) reflect.Type {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	return rv.Type()
}
