package vos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion_CheckConstants(t *testing.T) {
	assert.Equal(t, IgnoreAccountVersion, Version(-1))
	assert.Equal(t, NextAccountVersion, Version(0))
}
