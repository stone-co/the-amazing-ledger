package vos

import "strings"

type OperationType int8

const (
	InvalidOperation OperationType = iota
	DebitOperation
	CreditOperation
)

var _operationTypes = []string{"invalid_op_type", "debit", "credit"}

func (ot OperationType) String() string {
	return _operationTypes[ot]
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
