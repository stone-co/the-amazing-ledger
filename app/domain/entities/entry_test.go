package entities

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func TestNewEntry(t *testing.T) {
	type args struct {
		id        uuid.UUID
		operation vos.OperationType
		account   string
		version   vos.Version
		amount    int
	}

	testCases := []struct {
		name        string
		args        args
		expectedErr error
	}{
		{
			name: "Successfully creates an entry with minimum inputs",
			args: args{
				id:        uuid.New(),
				operation: vos.CreditOperation,
				account:   "assets.bacen.conta_liquidacao.tesouraria",
				version:   vos.NextAccountVersion,
				amount:    123,
			},
			expectedErr: nil,
		},
		{
			name: "Invalid when entry id is invalid",
			args: args{
				id:        uuid.Nil,
				operation: vos.CreditOperation,
				account:   "assets.bacen.conta_liquidacao.tesouraria",
				version:   vos.NextAccountVersion,
				amount:    123,
			},
			expectedErr: app.ErrInvalidEntryID,
		},
		{
			name: "Invalid when operation is invalid",
			args: args{
				id:        uuid.New(),
				operation: vos.InvalidOperation,
				account:   "assets.bacen.conta_liquidacao.tesouraria",
				version:   vos.NextAccountVersion,
				amount:    123,
			},
			expectedErr: app.ErrInvalidOperation,
		},
		{
			name: "Invalid when amount is zero",
			args: args{
				id:        uuid.New(),
				operation: vos.CreditOperation,
				account:   "assets.bacen.conta_liquidacao.tesouraria",
				version:   vos.NextAccountVersion,
				amount:    0,
			},
			expectedErr: app.ErrInvalidAmount,
		},
		{
			name: "Invalid when amount < zero",
			args: args{
				id:        uuid.New(),
				operation: vos.CreditOperation,
				account:   "assets.bacen.conta_liquidacao.tesouraria",
				version:   vos.NextAccountVersion,
				amount:    -1,
			},
			expectedErr: app.ErrInvalidAmount,
		},
		{
			name: "Invalid when account structure has less than 3 levels",
			args: args{
				id:        uuid.New(),
				operation: vos.CreditOperation,
				account:   "assets.bacen",
				version:   vos.NextAccountVersion,
				amount:    123,
			},
			expectedErr: app.ErrInvalidAccountStructure,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			entry, err := NewEntry(tt.args.id, tt.args.operation, tt.args.account, tt.args.version, tt.args.amount)
			assert.ErrorIs(t, err, tt.expectedErr)

			if err != nil {
				assert.Empty(t, entry)
			} else {
				assert.Equal(t, tt.args.id, entry.ID)
				assert.Equal(t, tt.args.operation, entry.Operation)
				assert.Equal(t, tt.args.account, entry.Account.Name())
				assert.Equal(t, tt.args.version, entry.Version)
				assert.Equal(t, tt.args.amount, entry.Amount)
			}
		})
	}
}
