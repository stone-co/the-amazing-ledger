package ledger

const (
	ErrInvalidTransactionID    = ClientError("invalid transaction id")
	ErrInvalidEntryID          = ClientError("invalid entry id")
	ErrInvalidOperation        = ClientError("invalid operation")
	ErrInvalidAmount           = ClientError("invalid amount")
	ErrInvalidEntriesNumber    = ClientError("invalid entries number")
	ErrInvalidVersion          = ClientError("invalid version")
	ErrInvalidBalance          = ClientError("invalid balance")
	ErrInvalidAccountStructure = ClientError("invalid account structure")
	ErrAccountNotFound         = ClientError("account not found")
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
