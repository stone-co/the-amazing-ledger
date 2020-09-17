package entities

import "strings"

type OperationType int

const (
	InvalidOperation OperationType = iota
	DebitOperation
	CreditOperation
)

func (ot OperationType) String() string {
	return [...]string{"invalid_op_type", "debit", "credit"}[ot]
}

func OperationTypeFromString(opType string) OperationType {
	opType = strings.ToLower(opType)
	if opType == "debit" {
		return DebitOperation
	} else if opType == "credit" {
		return CreditOperation
	}

	return InvalidOperation
}
