# stringable

A tiny go package that helps converting values from/to a string.

[![Go](https://github.com/ggicci/stringable/actions/workflows/go.yaml/badge.svg)](https://github.com/ggicci/stringable/actions/workflows/go.yaml)
[![codecov](https://codecov.io/gh/ggicci/stringable/graph/badge.svg?token=YU7FGGOY60)](https://codecov.io/gh/ggicci/stringable)
[![Go Report Card](https://goreportcard.com/badge/github.com/ggicci/stringable)](https://goreportcard.com/report/github.com/ggicci/stringable)
[![Go Reference](https://pkg.go.dev/badge/github.com/ggicci/stringable.svg)](https://pkg.go.dev/github.com/ggicci/stringable)

## Basic API

```go
var yesno bool
sb, err := stringable.New(&yesno)

sb.FromString("true")
sb.ToString()
```

## Supported Builtin Types

- string, bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, complex64, complex128
- `time.Time`
- `[]byte`

## The Hybrid Stringable Instance

When calling `stringable.New(x)` with an instance `x` that is not a Stringable itself, nor any of the above builtin types, it will try to create a _"hybrid" Stringable instance_ from `x` for you.

Here is how the "hybrid" Stringable instance will be created:

1. Create a hybrid instance `h` from the given instance `x`;
2. If `x` has implemented one of [`stringable.StringMarshaler`](https://pkg.go.dev/github.com/ggicci/stringable#StringMarshaler) and [`encoding.TextMarshaler`](https://pkg.go.dev/encoding#TextMarshaler), `h` will use it as the implementation of `stringable.StringMarshaler`, i.e. the `ToString()` method;
3. If `x` has implemented one of [`stringable.StringUnmarshaler`](https://pkg.go.dev/github.com/ggicci/stringable#StringUnmarshaler) and [`encoding.TextUnmarshaler`](https://pkg.go.dev/encoding#TextUnmarshaler), `h` will use it as the implementation of `stringable.StringUnmarshaler`, i.e. the `FromString()` method;
4. As long as `h` has an implementation of either `stringable.StringMarshaler` or `stringable.StringUnmarshaler`, we consider `h` is a valid `Stringable` instance. You can require both by passing in a [`CompleteHybrid()` option](#hybrid-options) to `New` method. For a valid `h`, `stringable.New(x)` will return `h`. Otherwise, an `ErrUnsupportedType` occurs.

Example:

```go
type Location struct {
	X int
	Y int
}

func (l *Location) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("L(%d,%d)", l.X, l.Y)), nil
}

loc := &Location{3, 4}
sb, err := stringable.New(loc) // err is nil

sb.ToString() // L(3,4)
sb.FromString("L(5,6)") // ErrNotStringUnmarshaler, "not a StringUnmarshaler"
```

### Hybrid Options

1. `New(v, NoHybrid())`: prevent `New` from trying to create a hybrid instance from `v` at all. Instead, returns `ErrUnsupportedType`.
2. `New(v, CompleteHybrid())`: still allow `New` trying to create a hybrid instance from `v` if necessary, but with the present of `CompleteHybrid()` option, the returned hybrid instance must have a valid implementation of both `FromString` and `ToString`.

## Adapt/Override Existing Types

The [`Namespace.Adapt()`](https://pkg.go.dev/github.com/ggicci/stringable#Namespace.Adapt) API is used to customize the behaviour of `stringable.Stringable` of a specific type. The principal is to create a **type alias** to the target type you want to override, and implement the `Stringable` interface on the new type.

When should you use this API?

1. change the conversion logic of the builtin types.
2. change the conversion logic of existing types that are "hybridizable", but you don't want to change their implementations.

For example, the default support of `bool` type in this package uses `strconv.ParseBool` method to convert strings like "true", "TRUE", "f", "0", etc. to a bool value. If you want to support also converting "YES", "NO", "はい" to a bool value, you can implement a custom bool type and register it to a `Namespace` instance:

```go
type YesNo bool

func (yn YesNo) ToString() (string, error) {
	if yn {
		return "yes", nil
	} else {
		return "no", nil
	}
}

func (yn *YesNo) FromString(s string) error {
	switch strings.ToLower(s) {
	case "yes":
		*yn = true
	case "no":
		*yn = false
	default:
		return errors.New("invalid value")
	}
	return nil
}

func main() {
	ns := stringable.NewNamespace()
	typ, adaptor := ToAnyStringableAdaptor(func(b *bool) (Stringable, error) {
		return (*YesNo)(b), nil
	})
	ns.Adapt(typ, adaptor)

	var yesno bool = true
	sb, err := ns.New(&yesno)
}
```
