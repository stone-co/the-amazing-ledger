package vos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSyntheticReport(t *testing.T) {
	version := Version(1)
	accountPathLiability := "liability"
	accountPathAssets := "assets"

	paths := []Path{
		{
			*&accountPathLiability,
			200,
			300,
		},
		{
			*&accountPathAssets,
			400,
			500,
		},
	}

	syntheticReport, err := NewSyntheticReport(600, 800, paths, version)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(syntheticReport.Paths))
}
