package vo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperationType_DefiningValuesByEnumType(t *testing.T) {
	op := CreditOperation
	assert.Equal(t, "credit", op.String())

	op = DebitOperation
	assert.Equal(t, "debit", op.String())

	op = InvalidOperation
	assert.Equal(t, "invalid_op_type", op.String())
}

func TestOperationType_DefiningValuesByString(t *testing.T) {
	var op OperationType
	op = OperationTypeFromString("credit")
	assert.Equal(t, "credit", op.String())

	op = OperationTypeFromString("debit")
	assert.Equal(t, "debit", op.String())

	op = OperationTypeFromString("CreDiT")
	assert.Equal(t, "credit", op.String())

	op = OperationTypeFromString("CreDeb")
	assert.Equal(t, "invalid_op_type", op.String())
}
