package testdata

import (
	"strings"

	"github.com/google/uuid"
)

func GenerateAccountQuery() string {
	return "liability.clients.available." + strings.ReplaceAll(uuid.New().String(), "-", "_") + ".*"
}

func GenerateInvalidAccountQuery() string {
	return "liability.clients." + strings.ReplaceAll(uuid.New().String(), "-", "_") + ".*"
}
