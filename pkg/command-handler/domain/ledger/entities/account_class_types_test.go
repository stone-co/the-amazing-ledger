package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountClassTypes_Has(t *testing.T) {
	t.Run("Valid account class types must return true", func(t *testing.T) {
		assert.True(t, ClassTypes.Has("liability"))
		assert.True(t, ClassTypes.Has("assets"))
		assert.True(t, ClassTypes.Has("income"))
		assert.True(t, ClassTypes.Has("expense"))
		assert.True(t, ClassTypes.Has("equity"))
	})

	t.Run("Invalid account class types must return false", func(t *testing.T) {
		assert.False(t, ClassTypes.Has("LIABILITY"))
		assert.False(t, ClassTypes.Has("XPTO"))
		assert.False(t, ClassTypes.Has("Assets"))
	})
}
