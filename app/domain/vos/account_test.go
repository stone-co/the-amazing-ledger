package vos

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app"
)

func TestNewAccount(t *testing.T) {
	tests := []struct {
		name       string
		account    string
		singleOnly bool
		want       Account
		wantErr    error
	}{
		{
			name:    "Single simple",
			account: "asset.account.example",
			want: Account{
				account:    "asset.account.example",
				represents: Single,
			},
			wantErr: nil,
		},
		{
			name:    "Single complete",
			account: "asset.account.abc_123",
			want: Account{
				account:    "asset.account.abc_123",
				represents: Single,
			},
			wantErr: nil,
		},
		{
			name:    "Single with upper",
			account: "asset.Account.example",
			want: Account{
				account:    "asset.account.example",
				represents: Single,
			},
			wantErr: nil,
		},
		{
			name:    "suffix simple",
			account: "*.asset.account",
			want: Account{
				account:    "*.asset.account",
				represents: Group,
			},
			wantErr: nil,
		},
		{
			name:    "prefix simple",
			account: "asset.account.*",
			want: Account{
				account:    "asset.account.*",
				represents: Group,
			},
			wantErr: nil,
		},
		{
			name:    "group simple",
			account: "asset.*.account",
			want: Account{
				account:    "asset.*.account",
				represents: Group,
			},
			wantErr: nil,
		},
		{
			name:    "group prefix composed",
			account: "asset.*.account.*",
			want: Account{
				account:    "asset.*.account.*",
				represents: Group,
			},
			wantErr: nil,
		},
		{
			name:    "group prefix simple",
			account: "asset.account*",
			want: Account{
				account:    "asset.account*",
				represents: Group,
			},
			wantErr: nil,
		},
		{
			name:    "group suffix composed",
			account: "*.asset.*.account",
			want: Account{
				account:    "*.asset.*.account",
				represents: Group,
			},
			wantErr: nil,
		},
		{
			name:    "group suffix simple",
			account: "*asset.account",
			want: Account{
				account:    "*asset.account",
				represents: Group,
			},
			wantErr: nil,
		},
		{
			name:    "empty account",
			account: "",
			want:    Account{},
			wantErr: app.ErrInvalidAccountStructure,
		},
		{
			name:    "empty component beginning",
			account: ".account",
			want:    Account{},
			wantErr: app.ErrInvalidAccountComponentSize,
		},
		{
			name:    "empty component middle",
			account: "asset..account",
			want:    Account{},
			wantErr: app.ErrInvalidAccountComponentSize,
		},
		{
			name:    "empty component end",
			account: "asset.",
			want:    Account{},
			wantErr: app.ErrInvalidAccountComponentSize,
		},
		{
			name:    "non ascii characters",
			account: "asset.áccount",
			want:    Account{},
			wantErr: app.ErrInvalidAccountComponentCharacters,
		},
		{
			name:       "Single only should fail if any '*' is present",
			account:    "*.account",
			singleOnly: true,
			want:       Account{},
			wantErr:    app.ErrInvalidSingleAccountComponentCharacters,
		},
		{
			name:    "should fail with invalid class",
			account: "assets.account",
			want:    Account{},
			wantErr: app.ErrAccountPathViolation,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				got Account
				err error
			)

			if tt.singleOnly {
				got, err = NewSingleAccount(tt.account)
			} else {
				got, err = NewAccount(tt.account)
			}

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
