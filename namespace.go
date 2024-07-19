package stringable

import (
	"fmt"
	"reflect"
	"time"

	"github.com/ggicci/stringable/internal"
)

var defaultNS = NewNamespace()

type Namespace struct {
	adaptors map[reflect.Type]AnyStringableAdaptor
}

// NewNamespace creates a namespace where you can register adaptors to
// override/adapt the converting behaviours of existing types.
func NewNamespace() *Namespace {
	return &Namespace{
		adaptors: make(map[reflect.Type]AnyStringableAdaptor),
	}
}

// New creates a Stringable instance from the given value.
func (c *Namespace) New(v any, opts ...Option) (Stringable, error) {
	if vs, ok := v.(Stringable); ok {
		return vs, nil
	}

	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return c.createStringable(v, options)
}

func (c *Namespace) createStringable(v any, opts *options) (Stringable, error) {
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
	if !opts.Has(optionNoHybrid) {
		h := createHybridStringable(rv)
		if h != nil {
			if opts.Has(optionCompleteHybrid) {
				if err := h.(*hybrid).validateAsComplete(); err != nil {
					return nil, err
				}
			}
			return h, nil
		}
	}

	return nil, unsupportedType(baseType)
}

// Adapt registers a custom adaptor for the given type. You can call
// ToAnyStringableAdaptor to create an adaptor of a specific type.
//
// Example:
//
//	ns := stringable.NewNamespace()
//	typ, adaptor := stringable.ToAnyStringableAdaptor[bool](func(b *bool) (stringable.Stringable, error) {})
//	ns.Adapt(typ, adaptor)
func (c *Namespace) Adapt(typ reflect.Type, adaptor AnyStringableAdaptor) {
	c.adaptors[typ] = adaptor
}

func unsupportedType(rt reflect.Type) error {
	return fmt.Errorf("%w: %v", ErrUnsupportedType, rt)
}

func init() {
	builtinStringable[string](func(v *string) (Stringable, error) { return (*internal.String)(v), nil })
	builtinStringable[bool](func(v *bool) (Stringable, error) { return (*internal.Bool)(v), nil })
	builtinStringable[int](func(v *int) (Stringable, error) { return (*internal.Int)(v), nil })
	builtinStringable[int8](func(v *int8) (Stringable, error) { return (*internal.Int8)(v), nil })
	builtinStringable[int16](func(v *int16) (Stringable, error) { return (*internal.Int16)(v), nil })
	builtinStringable[int32](func(v *int32) (Stringable, error) { return (*internal.Int32)(v), nil })
	builtinStringable[int64](func(v *int64) (Stringable, error) { return (*internal.Int64)(v), nil })
	builtinStringable[uint](func(v *uint) (Stringable, error) { return (*internal.Uint)(v), nil })
	builtinStringable[uint8](func(v *uint8) (Stringable, error) { return (*internal.Uint8)(v), nil })
	builtinStringable[uint16](func(v *uint16) (Stringable, error) { return (*internal.Uint16)(v), nil })
	builtinStringable[uint32](func(v *uint32) (Stringable, error) { return (*internal.Uint32)(v), nil })
	builtinStringable[uint64](func(v *uint64) (Stringable, error) { return (*internal.Uint64)(v), nil })
	builtinStringable[float32](func(v *float32) (Stringable, error) { return (*internal.Float32)(v), nil })
	builtinStringable[float64](func(v *float64) (Stringable, error) { return (*internal.Float64)(v), nil })
	builtinStringable[complex64](func(v *complex64) (Stringable, error) { return (*internal.Complex64)(v), nil })
	builtinStringable[complex128](func(v *complex128) (Stringable, error) { return (*internal.Complex128)(v), nil })
	builtinStringable[time.Time](func(v *time.Time) (Stringable, error) { return (*internal.Time)(v), nil })
	builtinStringable[[]byte](func(b *[]byte) (Stringable, error) { return (*internal.ByteSlice)(b), nil })
}
