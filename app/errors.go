package app

import (
	"errors"
)

const (
	ErrInvalidTransactionID    = DomainError("invalid transaction id")
	ErrInvalidEntryID          = DomainError("invalid entry id")
	ErrInvalidOperation        = DomainError("invalid operation")
	ErrInvalidAmount           = DomainError("invalid amount")
	ErrInvalidEntriesNumber    = DomainError("invalid entries number")
	ErrInvalidBalance          = DomainError("invalid balance")
	ErrIdempotencyKeyViolation = DomainError("idempotency key violation")
	ErrInvalidVersion          = DomainError("invalid version")
	ErrAccountNotFound         = DomainError("account not found")
	ErrInvalidAccountStructure = DomainError("invalid account structure")
	ErrInvalidClassName        = DomainError("invalid class name")
)

type DomainError string

func (err DomainError) Error() string {
	return string(err)
}

func (err DomainError) Is(target error) bool {
	return errors.Is(err, target)
}
