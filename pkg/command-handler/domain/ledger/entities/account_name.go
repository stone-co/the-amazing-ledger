package entities

import "strings"

var (
	classTypes = "|liability|assets|income|expense|equity|"
)

// AccountName contains exactly 4 levels in her structure: "class:group:subgroup:id", where:
//   - class could be liability, assets, income, expense or equity;
//   - group, subgroup and id are "free text".
//
// ":" must be used as the separator.
type AccountName struct { // TODO: could be just "Account", but already exists the type "Account"
	name string
}

func NewAccountName(name string) (*AccountName, error) {
	name = strings.ToLower(name)

	levels := strings.Split(name, ":")
	if len(levels) < 4 {
		return nil, ErrInvalidAccountStructure
	}

	for _, v := range levels {
		if len(v) == 0 {
			return nil, ErrInvalidAccountStructure
		}
	}

	if !strings.Contains(classTypes, levels[0]) {
		return nil, ErrInvalidAccountStructure
	}

	return &AccountName{name}, nil
}

func (a AccountName) Name() string {
	return a.name
}