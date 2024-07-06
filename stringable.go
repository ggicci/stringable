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

// New creates a Stringable instance from the given value. Remember that this
// function is a wrapper around defaultFactory.New, which means you're not able
// override/adapt existing types. If you want to override, you should create a
// new Factory and register the custom types using the Adapt method.
func New(v any) (Stringable, error) {
	return defaultFactory.New(v)
}