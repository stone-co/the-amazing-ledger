package entities

import "errors"

var (
	NoError                 = errors.New("")
	ErrInvalidData          = errors.New("invalid data")
	ErrInvalidEntriesNumber = errors.New("invalid entries number")
	ErrInvalidBalance       = errors.New("invalid balance")
	ErrIdempotencyKey       = errors.New("idempotency key violation")
	ErrInvalidVersion       = errors.New("invalid version")
)
