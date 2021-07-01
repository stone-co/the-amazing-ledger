package vos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSyntheticReport(t *testing.T) {
	accountPathLiability := "liability"
	accountPathAssets := "assets"

	accountPath1, err := NewAccountPath(accountPathLiability)
	assert.NotNil(t, err)
	accountPath2, err := NewAccountPath(accountPathAssets)
	assert.NotNil(t, err)

	paths := []Path{
		{
			accountPath1,
			200,
			300,
		},
		{
			accountPath2,
			400,
			500,
		},
	}

	syntheticReport, err := NewSyntheticReport(600, 800, paths)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(syntheticReport.Paths))
}
