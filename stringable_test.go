package stringable

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/ggicci/stringable/internal"
	"github.com/stretchr/testify/assert"
)

func TestNew_string(t *testing.T) {
	var s string = "hello"
	sv, err := New(&s)
	assert.NoError(t, err)
	got, err := sv.ToString()
	assert.NoError(t, err)
	assert.Equal(t, "hello", got)
	sv.FromString("world")
	assert.Equal(t, "world", s)
}

func TestNew_bool(t *testing.T) {
	var b bool = true
	sv, err := New(&b)
	assert.NoError(t, err)
	got, err := sv.ToString()
	assert.NoError(t, err)
	assert.Equal(t, "true", got)
	sv.FromString("false")
	assert.Equal(t, false, b)

	assert.Error(t, sv.FromString("hello"))
}

func TestNew_int(t *testing.T) {
	testInteger[int](t, 2045, "hello")
}

func TestNew_int8(t *testing.T) {
	testInteger[int8](t, int8(127), "128")
}

func TestNew_int16(t *testing.T) {
	testInteger[int16](t, int16(32767), "32768")
}

func TestNew_int32(t *testing.T) {
	testInteger[int32](t, int32(2147483647), "2147483648")
}

func TestNew_int64(t *testing.T) {
	testInteger[int64](t, int64(9223372036854775807), "9223372036854775808")
}

func TestNew_uint(t *testing.T) {
	testInteger[uint](t, uint(2045), "-1")
}

func TestNew_uint8(t *testing.T) {
	testInteger[uint8](t, uint8(255), "256")
}

func TestNew_uint16(t *testing.T) {
	testInteger[uint16](t, uint16(65535), "65536")
}

func TestNew_uint32(t *testing.T) {
	testInteger[uint32](t, uint32(4294967295), "4294967296")
}

func TestNew_uint64(t *testing.T) {
	testInteger[uint64](t, uint64(18446744073709551615), "18446744073709551616")
}

func TestNew_float32(t *testing.T) {
	testInteger[float32](t, float32(3.1415926), "hello")
}

func TestNew_float64(t *testing.T) {
	testInteger[float64](t, float64(3.14159265358979323846264338327950288419716939937510582097494459), "hello")
}

func TestNew_complex64(t *testing.T) {
	testInteger[complex64](t, complex64(3.1415926+2.71828i), "hello")
}

func TestNew_complex128(t *testing.T) {
	testInteger[complex128](t, complex128(3.14159265358979323846264338327950288419716939937510582097494459+2.71828182845904523536028747135266249775724709369995957496696763i), "hello")
}

func TestNew_Time(t *testing.T) {
	var now = time.Now()
	rvTime := reflect.ValueOf(now)
	sb, err := New(rvTime)
	assert.Error(t, err)
	assert.Nil(t, sb)

	rvTimePointer := reflect.ValueOf(&now)
	sv, err := New(rvTimePointer)
	assert.NoError(t, err)

	// RFC3339Nano
	testTime(t, sv, "1991-11-10T08:00:00+08:00", time.Date(1991, 11, 10, 8, 0, 0, 0, time.FixedZone("Asia/Shanghai", +8*3600)), "1991-11-10T00:00:00Z")
	// Date string
	testTime(t, sv, "1991-11-10", time.Date(1991, 11, 10, 0, 0, 0, 0, time.UTC), "1991-11-10T00:00:00Z")

	// Unix timestamp
	testTime(t, sv, "678088800", time.Date(1991, 6, 28, 6, 0, 0, 0, time.UTC), "1991-06-28T06:00:00Z")

	// Unix timestamp fraction
	testTime(t, sv, "678088800.123456789", time.Date(1991, 6, 28, 6, 0, 0, 123456789, time.UTC), "1991-06-28T06:00:00.123456789Z")

	// Unsupported format
	assert.Error(t, sv.FromString("hello"))
}

func TestNew_ByteSlice(t *testing.T) {
	var b []byte = []byte("hello")
	rvByteSlice := reflect.ValueOf(b)
	sb, err := New(rvByteSlice)
	assert.Error(t, err)
	assert.Nil(t, sb)

	rvByteSlicePointer := reflect.ValueOf(&b)
	sv, err := New(rvByteSlicePointer)
	assert.NoError(t, err)
	got, err := sv.ToString()
	assert.NoError(t, err)
	assert.Equal(t, "aGVsbG8=", got)

	sv.FromString("d29ybGQ=")
	assert.Equal(t, []byte("world"), b)

	assert.Error(t, sv.FromString("hello"))
}

type StructNotStringable struct {
	Name string
}

func TestNew_ErrNotPointer(t *testing.T) {
	for _, v := range getBuiltinInstances() {
		rv := reflect.ValueOf(v)
		sb, err := New(rv)
		assert.ErrorIs(t, err, ErrNotPointer)
		assert.Nil(t, sb)
	}

	var s StructNotStringable
	sb, err := New(s)
	assert.ErrorIs(t, err, ErrNotPointer)
	assert.Nil(t, sb)
}

func TestNew_ErrNilPointer(t *testing.T) {
	var b *bool
	sb, err := New(b)
	assert.ErrorIs(t, err, ErrNilPointer)
	assert.Nil(t, sb)
}

func TestNew_ErrUnsupportedType(t *testing.T) {
	var s StructNotStringable
	sv, err := New(&s)
	assert.ErrorIs(t, err, ErrUnsupportedType)
	assert.Nil(t, sv)
}

type Numeric interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64 | complex64 | complex128
}

func testInteger[T Numeric](t *testing.T, vSuccess T, invalidStr string) {
	rv := reflect.ValueOf(vSuccess)
	sb, err := New(rv)
	assert.Error(t, err)
	assert.Nil(t, sb)

	rvPointer := reflect.ValueOf(&vSuccess)
	sv, err := New(rvPointer)
	assert.NoError(t, err)
	got, err := sv.ToString()
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("%v", vSuccess), got)
	sv.FromString("2")
	assert.Equal(t, T(2), vSuccess)

	assert.Error(t, sv.FromString(invalidStr))
}

func testTime(t *testing.T, sv Stringable, fromStr string, expected time.Time, expectedToStr string) {
	assert.NoError(t, sv.FromString(fromStr))
	assert.True(t, equalTime(expected, time.Time(*sv.(*internal.Time))))
	ts, err := sv.ToString()
	assert.NoError(t, err)
	assert.Equal(t, expectedToStr, ts)
}

func equalTime(expected, actual time.Time) bool {
	return expected.UTC() == actual.UTC()
}

func getBuiltinInstances() []any {
	return []any{
		string("hello"),
		bool(true),
		int(1),
		int8(1),
		int16(1),
		int32(1),
		int64(1),
		uint(1),
		uint8(1),
		uint16(1),
		uint32(1),
		uint64(1),
		float32(1.0),
		float64(1.0),
		complex64(1.0),
		complex128(1.0),
		time.Now(),
		[]byte("hello"),
	}
}
