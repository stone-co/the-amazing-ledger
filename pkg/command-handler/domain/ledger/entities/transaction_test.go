package entities

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewTransaction(t *testing.T) {
	type args struct {
		id      uuid.UUID
		entries []Entry
	}

	id := uuid.New()
	validTwoEntries := []Entry{
		*NewEntry(uuid.New(), DebitOperation, "account/111", AnyAccountVersion, 123),
		*NewEntry(uuid.New(), CreditOperation, "account/222", AnyAccountVersion, 123),
	}
	validThreeEntries := []Entry{
		*NewEntry(uuid.New(), DebitOperation, "account/333", AnyAccountVersion, 400),
		*NewEntry(uuid.New(), CreditOperation, "account/444", AnyAccountVersion, 300),
		*NewEntry(uuid.New(), CreditOperation, "account/555", AnyAccountVersion, 100),
	}

	tests := []struct {
		name        string
		args        args
		want        *Transaction
		expectedErr error
	}{
		{
			name: "Invalid entries number when the transaction has no entries",
			args: args{
				id: id,
			},
			want:        nil,
			expectedErr: ErrInvalidEntriesNumber,
		},
		{
			name: "Invalid entries number when the transaction has 1 entry",
			args: args{
				id: id,
			},
			want:        nil,
			expectedErr: ErrInvalidEntriesNumber,
		},
		{
			name: "Valid transaction with 2 entries",
			args: args{
				id:      id,
				entries: validTwoEntries,
			},
			want: &Transaction{
				ID:      id,
				Entries: validTwoEntries,
			},
			expectedErr: nil,
		},
		{
			name: "Valid transaction with 3 entries",
			args: args{
				id:      id,
				entries: validThreeEntries,
			},
			want: &Transaction{
				ID:      id,
				Entries: validThreeEntries,
			},
			expectedErr: nil,
		},
		{
			name: "Invalid transaction with 2 entries and balance != 0",
			args: args{
				id: id,
				entries: []Entry{
					*NewEntry(uuid.New(), DebitOperation, "account/111", AnyAccountVersion, 123),
					*NewEntry(uuid.New(), CreditOperation, "account/222", AnyAccountVersion, 234),
				},
			},
			want:        nil,
			expectedErr: ErrInvalidBalance,
		},
		{
			name: "Invalid transaction with 3 entries and balance != 0",
			args: args{
				id: id,
				entries: []Entry{
					*NewEntry(uuid.New(), DebitOperation, "account/111", AnyAccountVersion, 400),
					*NewEntry(uuid.New(), CreditOperation, "account/222", AnyAccountVersion, 200),
					*NewEntry(uuid.New(), CreditOperation, "account/333", AnyAccountVersion, 100),
				},
			},
			want:        nil,
			expectedErr: ErrInvalidBalance,
		},
		{
			name: "Invalid transaction with empty ID",
			args: args{
				entries: validTwoEntries,
			},
			want:        nil,
			expectedErr: ErrInvalidData,
		},
		{
			name: "Amount must be > 1",
			args: args{
				id: id,
				entries: []Entry{
					*NewEntry(uuid.New(), DebitOperation, "account/111", AnyAccountVersion, 0),
					*NewEntry(uuid.New(), CreditOperation, "account/222", AnyAccountVersion, 0),
				},
			},
			want:        nil,
			expectedErr: ErrInvalidData,
		},
		{
			name: "Creating a transaction with an invalid operation must fail",
			args: args{
				id: id,
				entries: []Entry{
					*NewEntry(uuid.New(), InvalidOperation, "account/111", AnyAccountVersion, 123),
					*NewEntry(uuid.New(), CreditOperation, "account/222", AnyAccountVersion, 123),
				},
			},
			want:        nil,
			expectedErr: ErrInvalidData,
		},
		{
			name: "Creating a transaction with an invalid operation must fail",
			args: args{
				id: id,
				entries: []Entry{
					*NewEntry(uuid.New(), InvalidOperation, "account/111", AnyAccountVersion, 123),
					*NewEntry(uuid.New(), CreditOperation, "account/222", AnyAccountVersion, 123),
				},
			},
			want:        nil,
			expectedErr: ErrInvalidData,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTransaction(tt.args.id, tt.args.entries...)
			assert.True(t, errors.Is(err, tt.expectedErr))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}
