package vos

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stretchr/testify/assert"
)

func TestNewPartialAccountName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Successfully creates with an empty account name",
			args: args{
				name: "",
			},
			err: nil,
		},
		{
			name: "Successfully creates an account name with 1 level",
			args: args{
				name: "liability",
			},
			err: nil,
		},
		{
			name: "Successfully creates an account name with 2 levels",
			args: args{
				name: "liability:aaa",
			},
			err: nil,
		},
		{
			name: "Successfully creates an account name with 3 levels",
			args: args{
				name: "liability:aaa:ccc",
			},
			err: nil,
		},
		{
			name: "Error when account has 4 leves",
			args: args{
				name: "liability:aaa:ccc:ddd",
			},
			err: app.ErrInvalidAccountStructure,
		},
		{
			name: "Error when account has more than 4 levels",
			args: args{
				name: "liability:clients:available:" + uuid.New().String() + ":invalid",
			},
			err: app.ErrInvalidAccountStructure,
		},
		{
			name: "Error when account omits level 1",
			args: args{
				name: ":bacen:conta_liquidacao",
			},
			err: app.ErrInvalidAccountStructure,
		},
		{
			name: "Error when account omits level 2",
			args: args{
				name: "assets::conta_liquidacao",
			},
			err: app.ErrInvalidAccountStructure,
		},
		{
			name: "Error when account omits level 3",
			args: args{
				name: "assets:bacen:",
			},
			err: app.ErrInvalidAccountStructure,
		},
		{
			name: "Error when level 1 is not one of the predefined values (assets, liability, ...)",
			args: args{
				name: "xpto:bacen:conta_liquidacao",
			},
			err: app.ErrInvalidAccountStructure,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAccountPath(tt.args.name)
			assert.Equal(t, tt.err, err)
			if err == nil {
				assert.Equal(t, tt.args.name, got.Name())
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
			account:     "assets:bacen:conta_liquidacao",
			expClass:    "assets",
			expGroup:    "bacen",
			expSubgroup: "conta_liquidacao",
		},
		{
			test:        "Successfully get data from a valid account",
			account:     "liability:clients:available",
			expClass:    "liability",
			expGroup:    "clients",
			expSubgroup: "available",
		},
	}

	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			got, err := NewAccountPath(tt.account)
			assert.Equal(t, nil, err)
			assert.Equal(t, tt.account, got.Name())
			assert.Equal(t, tt.expClass, got.Class.String())
			assert.Equal(t, tt.expGroup, got.Group)
			assert.Equal(t, tt.expSubgroup, got.Subgroup)
		})
	}
}
