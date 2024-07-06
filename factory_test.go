package stringable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFactory_Adapt(t *testing.T) {
	factory := NewFactory()
	typ, adaptor := ToAnyStringableAdaptor(func(b *bool) (Stringable, error) {
		return (*YesNo)(b), nil
	})
	factory.Adapt(typ, adaptor)

	assert.Contains(t, factory.adaptors, typ)

	var yesno bool = true
	sb, err := factory.New(&yesno)
	assert.NoError(t, err)
	assert.NoError(t, sb.FromString("no"))
	assert.False(t, yesno)
	assert.ErrorContains(t, sb.FromString("false"), "invalid value")
}

func TestFactory_NewWithHybridInstanceCreated(t *testing.T) {
	factory := NewFactory()

	orange := &textMarshalerAndUnmarshalerOrange{Content: "orange"}
	sb, err := factory.New(orange)
	assert.NoError(t, err)

	text, err := sb.ToString()
	assert.NoError(t, err)
	assert.Equal(t, "orange", text)

	err = sb.FromString("red orange")
	assert.NoError(t, err)
	assert.Equal(t, "red orange", orange.Content)
}
