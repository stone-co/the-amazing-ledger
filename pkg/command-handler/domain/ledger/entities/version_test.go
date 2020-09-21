package entities

import (
	"testing"

	"gotest.tools/assert"
)

func TestVersion_Next(t *testing.T) {
	var v Version
	assert.Equal(t, Version(0), v.Current())
	assert.Equal(t, Version(1), v.Next())
	assert.Equal(t, Version(1), v.Current())
	assert.Equal(t, Version(2), v.Next())
	assert.Equal(t, Version(2), v.Current())
}
