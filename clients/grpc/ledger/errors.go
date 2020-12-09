package ledger

import (
	"fmt"
	"strings"
)

const (
	ErrInvalidData             = ClientError("invalid data")
	ErrInvalidEntriesNumber    = ClientError("invalid entries number")
	ErrInvalidVersion          = ClientError("invalid version")
	ErrInvalidAccountStructure = ClientError("invalid account structure")
	ErrAccountNotFound         = ClientError("account not found")
	ErrConnectionFailed        = ClientError("connection failed")
	ErrUndefined               = ClientError("undefined")
	ErrUnmaped                 = ClientError("unmaped")
)

type ClientError string

func (err ClientError) Error() string {
	return string(err)
}

func (err ClientError) Is(target error) bool {
	ts := target.Error()
	es := string(err)
	return ts == es || strings.HasPrefix(ts, es+": ")
}

func (err ClientError) cause(inner error) error {
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
	return ClientError(err.msg).Is(target)
}
