package entities

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
				name: "assets:bacen:conta_liquidacao:tesouraria",
			},
			err: nil,
		},
		{
			name: "Successfully creates an account name with more than 4 leves",
			args: args{
				name: "liability:clients:available:" + uuid.New().String() + ":detail",
			},
			err: nil,
		},
		{
			name: "Error when account has only 1 level",
			args: args{
				name: "assets",
			},
			err: ErrInvalidAccountStructure,
		},
		{
			name: "Error when account has only 2 leves",
			args: args{
				name: "assets:bacen",
			},
			err: ErrInvalidAccountStructure,
		},
		{
			name: "Error when account has only 3 levels",
			args: args{
				name: "assets:bacen:conta_liquidacao",
			},
			err: ErrInvalidAccountStructure,
		},
		{
			name: "Error when account omits level 1",
			args: args{
				name: ":bacen:conta_liquidacao:tesouraria",
			},
			err: ErrInvalidAccountStructure,
		},
		{
			name: "Error when account omits level 2",
			args: args{
				name: "assets::conta_liquidacao:tesouraria",
			},
			err: ErrInvalidAccountStructure,
		},
		{
			name: "Error when account omits level 3",
			args: args{
				name: "assets:bacen::tesouraria",
			},
			err: ErrInvalidAccountStructure,
		},
		{
			name: "Error when account omits level 4",
			args: args{
				name: "assets:bacen:conta_liquidacao:",
			},
			err: ErrInvalidAccountStructure,
		},
		{
			name: "Error when level 1 is not one of the predefined values (assets, liability, ...)",
			args: args{
				name: "xpto:bacen:conta_liquidacao:tesouraria",
			},
			err: ErrInvalidAccountStructure,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAccountName(tt.args.name)
			assert.Equal(t, tt.err, err)
			if err == nil {
				assert.Equal(t, &AccountName{tt.args.name}, got)
			}
		})
	}
}
