package entities

import "errors"

var (
	ErrInvalidData             = errors.New("invalid data")
	ErrInvalidEntriesNumber    = errors.New("invalid entries number")
	ErrInvalidBalance          = errors.New("invalid balance")
	ErrIdempotencyKey          = errors.New("idempotency key violation")
	ErrInvalidVersion          = errors.New("invalid version")
	ErrInvalidAccountStructure = errors.New("invalid account structure")
)
