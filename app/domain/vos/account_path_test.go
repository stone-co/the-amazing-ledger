package vos

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app"
)

func TestNewAccount(t *testing.T) {
	id := strings.ReplaceAll(uuid.New().String(), "-", "_")

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
				name: "assets.bacen.conta_liquidacao",
			},
			err: nil,
		},
		{
			name: "Successfully creates an account name with depth 4",
			args: args{
				name: "liability.clients.available." + id,
			},
			err: nil,
		},
		{
			name: "Successfully creates an account name with depth 5",
			args: args{
				name: "liability.clients.available." + id + ".mydetail",
			},
			err: nil,
		},
		{
			name: "Successfully creates an account name with depth 6",
			args: args{
				name: "liability.clients.available." + id + ".mydetail1.mydetail2_mydetail3",
			},
			err: nil,
		},
		{
			name: "Error when account has depth 1",
			args: args{
				name: "assets",
			},
			err: app.ErrInvalidAccountStructure,
		},
		{
			name: "Error when account has a depth 2",
			args: args{
				name: "assets.bacen",
			},
			err: app.ErrInvalidAccountStructure,
		},
		{
			name: "Error when account has an empty depth",
			args: args{
				name: "liability..treasury",
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
			name: "Error when depth 1 value is not one of the available",
			args: args{
				name: "xpto.conta_liquidacao.tesouraria",
			},
			err: app.ErrInvalidAccountStructure,
		},
		{
			name: "Error when account has invalid characters",
			args: args{
				name: "assets.conta_liquidacao." + uuid.New().String(),
			},
			err: app.ErrInvalidAccountStructure,
		},
		{
			name: "Error when label has more thant 255 characters characters",
			args: args{
				name: "assets.bacen.conta_liquidacao." + strings.Repeat("a", maxLabelLength+1),
			},
			err: app.ErrInvalidAccountStructure,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAccountPath(tt.args.name)
			assert.ErrorIs(t, err, tt.err)

			if err == nil {
				assert.Equal(t, tt.args.name, got.Name())
			}
		})
	}
}
