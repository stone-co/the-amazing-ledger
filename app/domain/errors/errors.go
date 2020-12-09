package errors

import (
	"errors"
	"fmt"
	"strings"
)

const (
	ErrInvalidData             = DomainError("invalid data")
	ErrInvalidEntriesNumber    = DomainError("invalid entries number")
	ErrInvalidBalance          = DomainError("invalid balance")
	ErrIdempotencyKey          = DomainError("idempotency key violation")
	ErrInvalidVersion          = DomainError("invalid version")
	ErrInvalidAccountStructure = DomainError("invalid account structure")
	ErrInvalidClassName        = DomainError("invalid class name")
	ErrAccountNotFound         = DomainError("account not found")
)

type DomainError string

func (err DomainError) Error() string {
	return string(err)
}

func (err DomainError) Is(target error) bool {
	ts := target.Error()
	es := string(err)
	return ts == es || strings.HasPrefix(ts, es+": ")
}

func (err DomainError) Cause(inner string) error {
	return wrapCause{msg: string(err), err: errors.New(inner)}
}

type wrapCause struct {
	err error
	msg string
}

func (err wrapCause) Error() string {
	if err.err != nil {
		return fmt.Sprintf("%s: %v", err.msg, err.err)
	}
	return err.msg
}

func (err wrapCause) Unwrap() error {
	return err.err
}

func (err wrapCause) Is(target error) bool {
	return DomainError(err.msg).Is(target)
}
