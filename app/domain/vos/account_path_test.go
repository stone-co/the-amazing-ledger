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

	type wants struct {
		name string
		err  error
	}
	tests := []struct {
		name  string
		wants wants
	}{
		{
			name: "Successfully creates an account name with minimum inputs",
			wants: wants{
				name: "assets.bacen.conta_liquidacao",
				err:  nil,
			},
		},
		{
			name: "Successfully creates an account name with depth 4",
			wants: wants{
				name: "liability.clients.available." + id,
				err:  nil,
			},
		},
		{
			name: "Successfully creates an account name with depth 5",
			wants: wants{
				name: "liability.clients.available." + id + ".mydetail",
				err:  nil,
			},
		},
		{
			name: "Successfully creates an account name with depth 6",
			wants: wants{
				name: "liability.clients.available." + id + ".mydetail1.mydetail2_mydetail3",
				err:  nil,
			},
		},
		{
			name: "Error when account has depth 1",
			wants: wants{
				name: "assets",
				err:  app.ErrInvalidAccountStructure,
			},
		},
		{
			name: "Error when account has a depth 2",
			wants: wants{
				name: "assets.bacen",
				err:  app.ErrInvalidAccountStructure,
			},
		},
		{
			name: "Error when account has an empty depth",
			wants: wants{
				name: "liability..treasury",
				err:  app.ErrInvalidAccountComponentSize,
			},
		},
		{
			name: "Error when account omits level 3",
			wants: wants{
				name: "assets.conta_liquidacao.",
				err:  app.ErrInvalidAccountComponentSize,
			},
		},
		{
			name: "Error when depth 1 value is not one of the available",
			wants: wants{
				name: "xpto.conta_liquidacao.tesouraria",
				err:  app.ErrAccountPathViolation,
			},
		},
		{
			name: "Error when account has invalid characters",
			wants: wants{
				name: "assets.conta_liquidacao." + uuid.New().String(),
				err:  app.ErrInvalidAccountComponentCharacters,
			},
		},
		{
			name: "Error when label has more thant 255 characters characters",
			wants: wants{
				name: "assets.bacen.conta_liquidacao." + strings.Repeat("a", maxLabelLength+1),
				err:  app.ErrInvalidAccountComponentSize,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAccountPath(tt.wants.name)
			assert.ErrorIs(t, err, tt.wants.err)

			if err == nil {
				assert.Equal(t, tt.wants.name, got.Name())
			}
		})
	}
}
