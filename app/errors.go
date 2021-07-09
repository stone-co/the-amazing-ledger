package app

const (
	ErrInvalidTransactionID              = DomainError("invalid transaction id")
	ErrInvalidEntryID                    = DomainError("invalid entry id")
	ErrInvalidOperation                  = DomainError("invalid operation")
	ErrInvalidAmount                     = DomainError("invalid amount")
	ErrInvalidEntriesNumber              = DomainError("invalid entries number")
	ErrInvalidBalance                    = DomainError("invalid balance")
	ErrIdempotencyKeyViolation           = DomainError("idempotency key violation")
	ErrInvalidVersion                    = DomainError("invalid version")
	ErrAccountNotFound                   = DomainError("account not found")
	ErrInvalidAccountStructure           = DomainError("account does not meet minimum or maximum supported sizes")
	ErrInvalidAccountComponentSize       = DomainError("account component cannot be empty and must be less than 256 characters")
	ErrInvalidAccountComponentCharacters = DomainError("only alphanumeric and underscore characters are supported")
	ErrAccountPathViolation              = DomainError("invalid depth value")
	ErrInvalidSyntheticReportStructure   = DomainError("invalid synthetic report structure")
	ErrEmptyAccountPath                  = DomainError("empty account path")
)

type DomainError string

func (err DomainError) Error() string {
	return string(err)
}
