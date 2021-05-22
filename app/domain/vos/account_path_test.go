package vos

import (
	"testing"

	"github.com/google/uuid"
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
			name:         "Successfully creates with an empty account name",
			path:         "",
			expectedName: "",
			err:          nil,
		},
		{
			name:         "Successfully creates an account name with 1 level",
			path:         "liability",
			expectedName: "liability",
			err:          nil,
		},
		{
			name:         "Successfully creates an account name with 2 levels",
			path:         "liability.aaa",
			expectedName: "liability.aaa",
			err:          nil,
		},
		{
			name:         "Successfully creates an account name with 3 levels",
			path:         "liability.aaa.ccc",
			expectedName: "liability.aaa",
			err:          nil,
		},
		{
			name:         "Successfully creates an account name with 4 levels",
			path:         "liability.clients.available." + uuid.New().String() + ".sub1",
			expectedName: "liability.clients",
			err:          nil,
		},
		{
			name:         "Successfully creates an account name with 5 levels",
			path:         "liability.clients.available." + uuid.New().String() + ".sub1.sub2",
			expectedName: "liability.clients",
			err:          nil,
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
			got, err := NewAccountPath(tt.path)
			assert.Equal(t, tt.err, err)
			if err == nil {
				assert.Equal(t, tt.expectedName, got.Name())
			}
		})
	}
}

func TestNewPartialAccountNameIsSplitted(t *testing.T) {
	tests := []struct {
		test        string
		account     string
		expClass    string
		expGroup    string
		expSubgroup string
	}{
		{
			test:        "Successfully get data from a valid account",
			account:     "assets.conta_liquidacao",
			expClass:    "assets",
			expSubgroup: "conta_liquidacao",
		},
		{
			test:        "Successfully get data from a valid account",
			account:     "liability.available",
			expClass:    "liability",
			expSubgroup: "available",
		},
	}

	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			got, err := NewAccountPath(tt.account)
			assert.Equal(t, nil, err)
			assert.Equal(t, tt.account, got.Name())
			assert.Equal(t, tt.expClass, got.Class.String())
			assert.Equal(t, tt.expSubgroup, got.Subgroup)
		})
	}
}
