package entities

import (
	"encoding/json"
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
		metadata  json.RawMessage
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
				account:   "asset.bacen.conta_liquidacao.tesouraria",
				version:   vos.NextAccountVersion,
				amount:    123,
				metadata:  json.RawMessage(`{}`),
			},
			expectedErr: nil,
		},
		{
			name: "Invalid when entry id is invalid",
			args: args{
				id:        uuid.Nil,
				operation: vos.CreditOperation,
				account:   "asset.bacen.conta_liquidacao.tesouraria",
				version:   vos.NextAccountVersion,
				amount:    123,
				metadata:  json.RawMessage(`{}`),
			},
			expectedErr: app.ErrInvalidEntryID,
		},
		{
			name: "Invalid when operation is invalid",
			args: args{
				id:        uuid.New(),
				operation: vos.InvalidOperation,
				account:   "asset.bacen.conta_liquidacao.tesouraria",
				version:   vos.NextAccountVersion,
				amount:    123,
				metadata:  json.RawMessage(`{}`),
			},
			expectedErr: app.ErrInvalidOperation,
		},
		{
			name: "Invalid when amount is zero",
			args: args{
				id:        uuid.New(),
				operation: vos.CreditOperation,
				account:   "asset.bacen.conta_liquidacao.tesouraria",
				version:   vos.NextAccountVersion,
				amount:    0,
				metadata:  json.RawMessage(`{}`),
			},
			expectedErr: app.ErrInvalidAmount,
		},
		{
			name: "Invalid when amount < zero",
			args: args{
				id:        uuid.New(),
				operation: vos.CreditOperation,
				account:   "asset.bacen.conta_liquidacao.tesouraria",
				version:   vos.NextAccountVersion,
				amount:    -1,
				metadata:  json.RawMessage(`{}`),
			},
			expectedErr: app.ErrInvalidAmount,
		},
		{
			name: "Invalid when account structure has less than 3 levels",
			args: args{
				id:        uuid.New(),
				operation: vos.CreditOperation,
				account:   "asset.bacen",
				version:   vos.NextAccountVersion,
				amount:    123,
				metadata:  json.RawMessage(`{}`),
			},
			expectedErr: app.ErrInvalidAccountStructure,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			entry, err := NewEntry(tt.args.id, tt.args.operation, tt.args.account, tt.args.version, tt.args.amount, tt.args.metadata)
			assert.ErrorIs(t, err, tt.expectedErr)

			if err != nil {
				assert.Empty(t, entry)
			} else {
				assert.Equal(t, tt.args.id, entry.ID)
				assert.Equal(t, tt.args.operation, entry.Operation)
				assert.Equal(t, tt.args.account, entry.Account.Value())
				assert.Equal(t, tt.args.version, entry.Version)
				assert.Equal(t, tt.args.amount, entry.Amount)
				assert.Equal(t, string(tt.args.metadata), string(entry.Metadata))
			}
		})
	}
}
