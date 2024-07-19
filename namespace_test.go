package stringable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNamespace_Adapt(t *testing.T) {
	ns := NewNamespace()
	typ, adaptor := ToAnyStringableAdaptor(func(b *bool) (Stringable, error) {
		return (*YesNo)(b), nil
	})
	ns.Adapt(typ, adaptor)

	assert.Contains(t, ns.adaptors, typ)

	var yesno bool = true
	sb, err := ns.New(&yesno)
	assert.NoError(t, err)
	assert.NoError(t, sb.FromString("no"))
	assert.False(t, yesno)
	assert.ErrorContains(t, sb.FromString("false"), "invalid value")
}

func TestNamespace_NewWithHybridInstanceCreated(t *testing.T) {
	ns := NewNamespace()

	orange := &textMarshalerAndUnmarshalerOrange{Content: "orange"}
	sb, err := ns.New(orange)
	assert.NoError(t, err)

	text, err := sb.ToString()
	assert.NoError(t, err)
	assert.Equal(t, "orange", text)

	err = sb.FromString("red orange")
	assert.NoError(t, err)
	assert.Equal(t, "red orange", orange.Content)
}
