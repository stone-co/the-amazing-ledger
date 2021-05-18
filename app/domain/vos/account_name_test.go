package vos

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app"
)

func TestNewAccountName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Successfully creates an account name with minimum inputs",
			args: args{
				name: "assets.bacen.conta_liquidacao.tesouraria",
			},
			err: nil,
		},
		{
			name: "Successfully creates an account name with 4 levels",
			args: args{
				name: "liability.clients.available." + uuid.New().String(),
			},
			err: nil,
		},
		{
			name: "Successfully creates an account name with 4 levels and one detail",
			args: args{
				name: "liability.clients.available." + uuid.New().String() + "/mydetail",
			},
			err: nil,
		},
		{
			name: "Successfully creates an account name with 4 levels and a lot of details",
			args: args{
				name: "liability.clients.available." + uuid.New().String() + "/mydetail1/mydetail2/mydetail3",
			},
			err: nil,
		},
		{
			name: "Error when account has only 1 level",
			args: args{
				name: "assets",
			},
			err: app.ErrInvalidAccountStructure,
		},
		{
			name: "Error when account has only 2 levels",
			args: args{
				name: "assets.bacen",
			},
			err: app.ErrInvalidAccountStructure,
		},
		{
			name: "Error when account omits level 1",
			args: args{
				name: ".conta_liquidacao.tesouraria",
			},
			err: app.ErrInvalidAccountStructure,
		},
		{
			name: "Error when account omits level 2",
			args: args{
				name: "assets..tesouraria",
			},
			err: app.ErrInvalidAccountStructure,
		},
		{
			name: "Error when account omits level 3",
			args: args{
				name: "assets.conta_liquidacao.",
			},
			err: app.ErrInvalidAccountStructure,
		},
		{
			name: "Error when level 1 is not one of the predefined values (assets, liability, ...)",
			args: args{
				name: "xpto.conta_liquidacao.tesouraria",
			},
			err: app.ErrInvalidAccountStructure,
		},
		{
			name: "Error when sending empty suffix",
			args: args{
				name: "xpto.conta_liquidacao.tesouraria.",
			},
			err: app.ErrInvalidAccountStructure,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAccountName(tt.args.name)
			assert.Equal(t, tt.err, err)
			if err == nil {
				assert.Equal(t, tt.args.name, got.Name())
			}
		})
	}
}

func TestNewAccountNameIsSplitted(t *testing.T) {
	newUUID := uuid.New().String()

	tests := []struct {
		test        string
		account     string
		expClass    string
		expSubgroup string
		expID       string
		expSuffix   string
	}{
		{
			test:        "Successfully get data from a valid account",
			account:     "assets.conta_liquidacao.tesouraria",
			expClass:    "assets",
			expSubgroup: "conta_liquidacao",
			expID:       "tesouraria",
			expSuffix:   "",
		},
		{
			test:        "Successfully get data from a valid account",
			account:     "liability.available." + newUUID,
			expClass:    "liability",
			expSubgroup: "available",
			expID:       newUUID,
			expSuffix:   "",
		},
		{
			test:        "Successfully get data from a valid account",
			account:     "liability.available." + newUUID + ".mydetail",
			expClass:    "liability",
			expSubgroup: "available",
			expID:       newUUID,
			expSuffix:   "mydetail",
		},
		{
			test:        "Successfully get data from a valid account",
			account:     "liability.available." + newUUID + ".mydetail1.mydetail2.mydetail3",
			expClass:    "liability",
			expSubgroup: "available",
			expID:       newUUID,
			expSuffix:   "mydetail1.mydetail2.mydetail3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			got, err := NewAccountName(tt.account)
			assert.Equal(t, nil, err)
			assert.Equal(t, tt.account, got.Name())
			assert.Equal(t, tt.expClass, got.Class.String())
			assert.Equal(t, tt.expSubgroup, got.Subgroup)
			assert.Equal(t, tt.expID, got.Account)
			assert.Equal(t, tt.expSuffix, got.Suffix)
		})
	}
}
