package entities

// ClassTypes has the set of class types.
var ClassTypes = AccountClassTypes{
	"liability": empty,
	"assets":    empty,
	"income":    empty,
	"expense":   empty,
	"equity":    empty,
}

type AccountClassTypes map[string]struct{}

var empty struct{}

func (act AccountClassTypes) Has(classType string) bool {
	_, ok := act[classType]
	return ok
}
