package entities

import "strings"

type AccountClass struct {
	string
}

var accountClasses = map[string]struct{}{
	"liability": empty,
	"assets":    empty,
	"income":    empty,
	"expense":   empty,
	"equity":    empty,
}

var empty struct{}

func NewAccountClassFromString(class string) (*AccountClass, error) {
	class = strings.ToLower(class)

	_, ok := accountClasses[class]
	if !ok {
		return nil, ErrInvalidClassName
	}

	return &AccountClass{class}, nil
}

func (ac AccountClass) String() string {
	return ac.string
}
