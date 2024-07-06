package stringable

import (
	"encoding"
	"reflect"
)

type hybrid struct {
	StringMarshaler
	StringUnmarshaler
}

func (h *hybrid) ToString() (string, error) {
	if h.StringMarshaler != nil {
		return h.StringMarshaler.ToString()
	}
	return "", ErrNotStringMarshaler
}

func (h *hybrid) FromString(s string) error {
	if h.StringUnmarshaler != nil {
		return h.StringUnmarshaler.FromString(s)
	}
	return ErrNotStringUnmarshaler
}

func (h *hybrid) IsValid() bool {
	return h.StringMarshaler != nil || h.StringUnmarshaler != nil
}

// createHybridStringable tries to create a hybrid Stringable from a
// reflect.Value. It will make the most of the interfaces rv has implemented,
// including stringable.StringMarshaler, stringable.StringUnmarshaler,
// encoding.TextMarshaler, and encoding.TextUnmarshaler. Returns nil if the
// reflect.Value does not implement any of the above.
func createHybridStringable(rv reflect.Value) Stringable {
	h := &hybrid{}

	// Check stringable.StringMarshaler and encoding.TextMarshaler.
	if rv.Type().Implements(stringMarshalerType) {
		h.StringMarshaler = rv.Interface().(StringMarshaler)
	} else if rv.Type().Implements(textMarshalerType) {
		h.StringMarshaler = &textMarshaler{
			rv.Interface().(encoding.TextMarshaler),
			nil,
		}
	}

	// Check stringable.StringUnmarshaler and encoding.TextUnmarshaler.
	if rv.Type().Implements(stringUnmarshalerType) {
		h.StringUnmarshaler = rv.Interface().(StringUnmarshaler)
	} else if rv.Type().Implements(textUnmarshalerType) {
		h.StringUnmarshaler = &textMarshaler{
			nil,
			rv.Interface().(encoding.TextUnmarshaler),
		}
	}

	if h.IsValid() {
		return h
	}

	return nil
}

type textMarshaler struct {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
}

func (w textMarshaler) ToString() (string, error) {
	b, err := w.TextMarshaler.MarshalText()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (w textMarshaler) FromString(s string) error {
	return w.TextUnmarshaler.UnmarshalText([]byte(s))
}

var (
	stringMarshalerType   = typeOf[StringMarshaler]()
	stringUnmarshalerType = typeOf[StringUnmarshaler]()
	textMarshalerType     = typeOf[encoding.TextMarshaler]()
	textUnmarshalerType   = typeOf[encoding.TextUnmarshaler]()
)
