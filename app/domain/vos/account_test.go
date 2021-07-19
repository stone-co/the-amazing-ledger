package vos

import (
	"reflect"
	"testing"
)

func TestNewAccount(t *testing.T) {
	tests := []struct {
		name       string
		account    string
		singleOnly bool
		want       Account
		wantErr    bool
	}{
		{
			name:    "Single simple",
			account: "asset.account.example",
			want: Account{
				account:    "asset.account.example",
				represents: Single,
			},
			wantErr: false,
		},
		{
			name:    "Single complete",
			account: "asset.account.abc_123",
			want: Account{
				account:    "asset.account.abc_123",
				represents: Single,
			},
			wantErr: false,
		},
		{
			name:    "Single with upper",
			account: "asset.Account.example",
			want: Account{
				account:    "asset.account.example",
				represents: Single,
			},
			wantErr: false,
		},
		{
			name:    "suffix simple",
			account: "*.asset.account",
			want: Account{
				account:    "*.asset.account",
				represents: Group,
			},
			wantErr: false,
		},
		{
			name:    "prefix simple",
			account: "asset.account.*",
			want: Account{
				account:    "asset.account.*",
				represents: Group,
			},
			wantErr: false,
		},
		{
			name:    "group simple",
			account: "asset.*.account",
			want: Account{
				account:    "asset.*.account",
				represents: Group,
			},
			wantErr: false,
		},
		{
			name:    "group prefix composed",
			account: "asset.*.account.*",
			want: Account{
				account:    "asset.*.account.*",
				represents: Group,
			},
			wantErr: false,
		},
		{
			name:    "group prefix simple",
			account: "asset.account*",
			want: Account{
				account:    "asset.account*",
				represents: Group,
			},
			wantErr: false,
		},
		{
			name:    "group suffix composed",
			account: "*.asset.*.account",
			want: Account{
				account:    "*.asset.*.account",
				represents: Group,
			},
			wantErr: false,
		},
		{
			name:    "group suffix simple",
			account: "*asset.account",
			want: Account{
				account:    "*asset.account",
				represents: Group,
			},
			wantErr: false,
		},
		{
			name:    "empty account",
			account: "",
			want:    Account{},
			wantErr: true,
		},
		{
			name:    "empty component beginning",
			account: ".account",
			want:    Account{},
			wantErr: true,
		},
		{
			name:    "empty component middle",
			account: "asset..account",
			want:    Account{},
			wantErr: true,
		},
		{
			name:    "empty component end",
			account: "asset.",
			want:    Account{},
			wantErr: true,
		},
		{
			name:    "non ascii characters",
			account: "asset.áccount",
			want:    Account{},
			wantErr: true,
		},
		{
			name:       "Single only should fail if any '*' is present",
			account:    "*.account",
			singleOnly: true,
			want:       Account{},
			wantErr:    true,
		},
		{
			name:    "should fail with invalid class",
			account: "assets.account",
			want:    Account{},
			wantErr: true,
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

			if (err != nil) != tt.wantErr {
				t.Errorf("NewAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAccount() got = %v, want %v", got, tt.want)
			}
		})
	}
}
