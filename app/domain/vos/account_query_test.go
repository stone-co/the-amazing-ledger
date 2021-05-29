package vos

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app"
)

func TestNewPartialAccountName(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		expectedName string
		err          error
	}{
		{
			name:         "Successfully creates a query with depth 1",
			path:         "liability",
			expectedName: "liability",
			err:          nil,
		},
		{
			name:         "Successfully creates a query with depth 2",
			path:         "liability.aaa",
			expectedName: "liability.aaa",
			err:          nil,
		},
		{
			name:         "Successfully creates a query with depth 3",
			path:         "liability.aaa.ccc",
			expectedName: "liability.aaa.ccc",
			err:          nil,
		},
		{
			name:         "Successfully creates a query with depth 4",
			path:         "liability.clients.available.sub1",
			expectedName: "liability.clients.available.sub1",
			err:          nil,
		},
		{
			name:         "Error when creating an empty query",
			path:         "",
			expectedName: "",
			err:          app.ErrInvalidAccountStructure,
		},
		{
			name:         "Error when account omits level 1",
			path:         ".bacen.conta_liquidacao",
			expectedName: "",
			err:          app.ErrInvalidAccountStructure,
		},
		{
			name:         "Error when account omits level 2",
			path:         "assets..conta_liquidacao",
			expectedName: "",
			err:          app.ErrInvalidAccountStructure,
		},
		{
			name:         "Error when account omits level 3",
			path:         "assets.bacen.",
			expectedName: "",
			err:          app.ErrInvalidAccountStructure,
		},
		{
			name:         "Error when level 1 is not one of the predefined values (assets, liability, ...)",
			path:         "xpto.bacen.conta_liquidacao",
			expectedName: "",
			err:          app.ErrInvalidAccountStructure,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAccountQuery(tt.path)
			assert.Equal(t, tt.err, err)
			if err == nil {
				assert.Equal(t, tt.expectedName, got.Value())
			}
		})
	}
}
