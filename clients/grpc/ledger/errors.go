package ledger

const (
	ErrInvalidTransactionID    = ClientError("invalid transaction id")
	ErrInvalidEntryID          = ClientError("invalid entry id")
	ErrInvalidOperation        = ClientError("invalid operation")
	ErrInvalidAmount           = ClientError("invalid amount")
	ErrInvalidEntriesNumber    = ClientError("invalid entries number")
	ErrInvalidBalance          = ClientError("invalid balance")
	ErrIdempotencyKeyViolation = ClientError("idempotency key violation")
	ErrInvalidVersion          = ClientError("invalid version")
	ErrAccountNotFound         = ClientError("account not found")
	ErrInvalidAccountStructure = ClientError("invalid account structure")
	ErrConnectionFailed        = ClientError("connection failed")
	ErrUndefined               = ClientError("undefined")
)

type ClientError string

func (err ClientError) Error() string {
	return string(err)
}

func (err ClientError) Is(target error) bool {
	if target == nil {
		return false
	}
	ts := target.Error()
	es := string(err)
	return ts == es
}
