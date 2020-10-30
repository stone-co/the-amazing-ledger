package entities

import (
	"testing"

	"gotest.tools/assert"
)

func TestAccountInfoBalance(t *testing.T) {
	accountInfo := NewAccountInfo("stone/banking", Version(3), 100, 50)

	assert.Equal(t, accountInfo.Balance(), 50)
}

func TestAccountInfoInvertedBalance(t *testing.T) {
	accountInfo := NewAccountInfo("stone/banking", Version(3), 100, 50)

	assert.Equal(t, accountInfo.InvertedBalance(), -50)
}
