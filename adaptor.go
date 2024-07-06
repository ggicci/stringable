package stringable

import (
	"fmt"
	"reflect"
)

type StringableAdaptor[T any] func(*T) (Stringable, error)
type AnyStringableAdaptor func(any) (Stringable, error)

func ToAnyStringableAdaptor[T any](adapt StringableAdaptor[T]) (reflect.Type, AnyStringableAdaptor) {
	return typeOf[T](), func(v any) (Stringable, error) {
		if cv, ok := v.(*T); ok {
			return adapt(cv)
		} else {
			return nil, fmt.Errorf("%w: cannot convert %T to %s", ErrTypeMismatch, v, typeOf[*T]())
		}
	}
}

var builtinStringableAdaptors = make(map[reflect.Type]AnyStringableAdaptor)

func builtinStringable[T any](adaptor StringableAdaptor[T]) {
	typ, anyAdaptor := ToAnyStringableAdaptor[T](adaptor)
	builtinStringableAdaptors[typ] = anyAdaptor
}
