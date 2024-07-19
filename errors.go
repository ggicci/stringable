package stringable

import "errors"

var (
	ErrUnsupportedType      = errors.New("unsupported type")
	ErrTypeMismatch         = errors.New("type mismatch")
	ErrNotStringMarshaler   = errors.New("not a StringMarshaler")
	ErrNotStringUnmarshaler = errors.New("not a StringUnmarshaler")
	ErrNotPointer           = errors.New("not a pointer")
	ErrNilPointer           = errors.New("nil pointer")
	ErrMissingMarshaler     = errors.New("missing marshaler")
	ErrMissingUnmarshaler   = errors.New("missing unmarshaler")
)
