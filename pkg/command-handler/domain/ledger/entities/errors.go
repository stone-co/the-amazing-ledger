package entities

import (
	"fmt"
	"strings"
)

const (
	ErrInvalidData             = EntityError("invalid data")
	ErrInvalidEntriesNumber    = EntityError("invalid entries number")
	ErrInvalidBalance          = EntityError("invalid balance")
	ErrIdempotencyKey          = EntityError("idempotency key violation")
	ErrInvalidVersion          = EntityError("invalid version")
	ErrInvalidAccountStructure = EntityError("invalid account structure")
	ErrInvalidClassName        = EntityError("invalid class name")
	ErrAccountNotFound         = EntityError("account not found")
)

type EntityError string

func (err EntityError) Error() string {
	return string(err)
}

func (err EntityError) Is(target error) bool {
	ts := target.Error()
	es := string(err)
	return ts == es || strings.HasPrefix(ts, es+": ")
}

func (err EntityError) cause(inner error) error {
	return wrapCause{msg: string(err), err: inner}
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
	return EntityError(err.msg).Is(target)
}
