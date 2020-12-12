package vos

import (
	"strings"
	"testing"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stretchr/testify/assert"
)

func TestNewAccountClassFromString(t *testing.T) {
	t.Run("Valid account class must return true", func(t *testing.T) {
		validClasses := []string{"liability", "assets", "income", "expense", "equity", "InCOme"}
		for _, class := range validClasses {
			got, err := NewAccountClassFromString(class)
			assert.Nil(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, strings.ToLower(class), got.String())
		}
	})

	t.Run("Invalid account class must return false", func(t *testing.T) {
		invalidClasses := []string{"ASSE1S", "XPTO"}
		for _, class := range invalidClasses {
			got, err := NewAccountClassFromString(class)
			assert.True(t, app.ErrInvalidClassName.Is(err))
			assert.Nil(t, got)
		}
	})
}
