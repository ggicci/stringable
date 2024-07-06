package stringable

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNil(t *testing.T) {
	assert.False(t, isNil(reflect.ValueOf("hello")))
	assert.True(t, isNil(reflect.ValueOf((*string)(nil))))
}

func TestPanicOnError(t *testing.T) {
	panicOnError(nil)

	assert.PanicsWithError(t, "httpin: "+assert.AnError.Error(), func() {
		panicOnError(assert.AnError)
	})
}

func TestTypeOf(t *testing.T) {
	assert.Equal(t, reflect.TypeOf(0), typeOf[int]())
}

func TestPointerize(t *testing.T) {
	assert.Equal(t, 102, *pointerize[int](102))
}

func TestDereferencedType(t *testing.T) {
	type Object struct{}

	var o = new(Object)
	var po = &o
	var ppo = &po
	assert.Equal(t, reflect.TypeOf(Object{}), dereferencedType(Object{}))
	assert.Equal(t, reflect.TypeOf(Object{}), dereferencedType(o))
	assert.Equal(t, reflect.TypeOf(Object{}), dereferencedType(po))
	assert.Equal(t, reflect.TypeOf(Object{}), dereferencedType(ppo))
}
