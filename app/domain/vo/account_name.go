package vo

import (
	"strings"

	"github.com/stone-co/the-amazing-ledger/app"
)

const (
	AccountStructureSep = ":"
	structureLevels     = 4
)

const (
	classLevel int = iota
	groupLevel
	subgroupLevel
	idLevel
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
	Class    *AccountClass
	Group    string
	Subgroup string
	ID       string
}

func NewAccountName(name string) (*AccountName, error) {
	name = strings.ToLower(name)

	levels := strings.Split(name, AccountStructureSep)
	if len(levels) != structureLevels {
		return nil, app.ErrInvalidAccountStructure
	}

	for _, v := range levels {
		if len(v) == 0 {
			return nil, app.ErrInvalidAccountStructure
		}
	}

	accountClass, err := NewAccountClassFromString(levels[classLevel])
	if err != nil {
		return nil, app.ErrInvalidAccountStructure
	}

	return &AccountName{
		Class:    accountClass,
		Group:    levels[groupLevel],
		Subgroup: levels[subgroupLevel],
		ID:       levels[idLevel],
	}, nil
}

func (a AccountName) Name() string {
	return FormatAccount(a.Class.String(), a.Group, a.Subgroup, a.ID)
}

func FormatAccount(class, group, subgroup, id string) string {
	return class + AccountStructureSep + group + AccountStructureSep + subgroup + AccountStructureSep + id
}
