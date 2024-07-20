// stringable is a tiny package that helps converting values from/to a string.
package stringable

// Stringable defines a type to be able to convert from/to a string.
type Stringable interface {
	StringMarshaler
	StringUnmarshaler
}

// StringMarshaler defines a type to be able to convert to a string.
type StringMarshaler interface {
	ToString() (string, error)
}

// StringUnmarshaler defines a type to be able to convert from a string.
type StringUnmarshaler interface {
	FromString(string) error
}

// New creates a Stringable instance from the given value. Note that
// this method is a wrapper around the default namespace's New method.
// Which means it doesn't support override/adapt existing types. Please
// read Namespace.New to learn more.
func New(v any) (Stringable, error) {
	return defaultNS.New(v)
}
