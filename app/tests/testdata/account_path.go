package testdata

import (
	"strings"

	"github.com/google/uuid"
)

func GenerateAccountPath() string {
	return "liability.clients.available." + strings.ReplaceAll(uuid.New().String(), "-", "_")
}
