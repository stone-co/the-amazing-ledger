package entities

import "strings"

var (
	classTypes      = "|liability|assets|income|expense|equity|"
	structureSep    = ":"
	structureLevels = 4
	classTypeLevel  = 0
)

// AccountName must have 4 levels in her structure: "class:group:subgroup:id", where:
//   - class could be liability, assets, income, expense or equity;
//   - group, subgroup and id are "free text";
//   - ":" must be used as a separator.
//
// Some examples:
//   - "assets:bacen:conta_liquidacao:tesouraria"
//   - "liability:clients:available:96a131a8-c4ac-495e-8971-fcecdbdd003a"
//   - "liability:clients:available:96a131a8-c4ac-495e-8971-fcecdbdd003a/somedetail"
//   - "liability:clients:available:96a131a8-c4ac-495e-8971-fcecdbdd003a/detail1/detail2"
type AccountName struct { // TODO: could be just "Account", but already exists the type "Account"
	name string
}

func NewAccountName(name string) (*AccountName, error) {
	name = strings.ToLower(name)

	levels := strings.Split(name, structureSep)
	if len(levels) != structureLevels {
		return nil, ErrInvalidAccountStructure
	}

	for _, v := range levels {
		if len(v) == 0 {
			return nil, ErrInvalidAccountStructure
		}
	}

	if !strings.Contains(classTypes, levels[classTypeLevel]) {
		return nil, ErrInvalidAccountStructure
	}

	return &AccountName{name}, nil
}

func (a AccountName) Name() string {
	return a.name
}
