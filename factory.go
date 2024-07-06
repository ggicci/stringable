package stringable

import (
	"fmt"
	"reflect"
)

var defaultFactory = NewFactory()

type Factory struct {
	adaptors map[reflect.Type]AnyStringableAdaptor
}

// NewFactory creates a Factory where you can register adaptors to
// override/adapt the converting behaviours of existing types.
func NewFactory() *Factory {
	return &Factory{
		adaptors: make(map[reflect.Type]AnyStringableAdaptor),
	}
}

// New creates a Stringable instance from the given value.
func (c *Factory) New(v any) (Stringable, error) {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}
	if rv.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("%w: value must be a non-nil pointer", ErrNotPointer)
	}
	if rv.IsNil() {
		return nil, fmt.Errorf("%w: value must be a non-nil pointer", ErrNilPointer)
	}

	baseType := rv.Type().Elem()

	// Check if there is a custom adaptor for the base type.
	if adapt, ok := c.adaptors[baseType]; ok {
		return adapt(rv.Interface())
	}

	// Check if there is a built-in adaptor for the base type.
	if adapt, ok := builtinStringableAdaptors[baseType]; ok {
		return adapt(rv.Interface())
	}

	// Try to create a hybrid Stringable from the reflect.Value.
	if hybrid := createHybridStringable(rv); hybrid != nil {
		return hybrid, nil
	}

	return nil, unsupportedType(baseType)
}

// Adapt registers a custom adaptor for the given type. You can call
// ToAnyStringableAdaptor to create an adaptor of a specific type.
//
// Example:
//
//	core := stringable.NewCore()
//	typ, adaptor := stringable.ToAnyStringableAdaptor[bool](func(b *bool) (stringable.Stringable, error) {})
//	core.Adapt(typ, adaptor)
func (c *Factory) Adapt(typ reflect.Type, adaptor AnyStringableAdaptor) {
	c.adaptors[typ] = adaptor
}

func unsupportedType(rt reflect.Type) error {
	return fmt.Errorf("%w: %v", ErrUnsupportedType, rt)
}
