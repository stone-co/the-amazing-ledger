package vos

import (
	"strings"

	"github.com/stone-co/the-amazing-ledger/app"
)

const (
	AccountStructureSep = ":"
	AccountSuffixSep    = "/"
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
// AccountName can have a fifth level in her structure: "class:group:subgroup:id/suffix", where:
//   - suffix is "free text";
//   - it is considered a suffix everything after the id
//   - "/" must be used as a separator
//
// Some examples:
//   - "assets:bacen:conta_liquidacao:tesouraria"
//   - "liability:clients:available:96a131a8-c4ac-495e-8971-fcecdbdd003a"
//   - "liability:clients:available:96a131a8-c4ac-495e-8971-fcecdbdd003a/somedetail"
//   - "liability:clients:available:96a131a8-c4ac-495e-8971-fcecdbdd003a/detail1/detail2"
type AccountName struct {
	Class    *AccountClass
	Group    string
	Subgroup string
	ID       string
	Suffix   string
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

	var suffix string
	identifiers := strings.SplitN(levels[idLevel], AccountSuffixSep, 2)
	if len(identifiers) == 1 {
		suffix = ""
	} else {
		suffix = identifiers[1]
	}

	return &AccountName{
		Class:    accountClass,
		Group:    levels[groupLevel],
		Subgroup: levels[subgroupLevel],
		ID:       identifiers[0],
		Suffix:   suffix,
	}, nil
}

func (a AccountName) Name() string {
	return FormatAccount(a.Class.String(), a.Group, a.Subgroup, a.ID, a.Suffix)
}

func FormatAccount(class, group, subgroup, id, suffix string) string {
	name := class + AccountStructureSep + group + AccountStructureSep + subgroup + AccountStructureSep + id
	if suffix != "" {
		name += AccountSuffixSep + suffix
	}
	return name
}
