package vos

import (
	"strings"

	"github.com/stone-co/the-amazing-ledger/app"
)

// AccountPath could have between 0 (empty) and 3 levels in her structure: "class:group:subgroup".
// And, as AccountName, AccountPath has pre-defined classes.
//
// Some examples:
//   - ""
//   - "liability"
//   - "liability:clients"
//   - "liability:clients:available"
//   - "assets:bacen:conta_liquidacao"
type AccountPath struct {
	Class       *AccountClass
	Group       string
	Subgroup    string
	TotalLevels int
}

func NewAccountPath(name string) (*AccountPath, error) {
	name = strings.ToLower(name)

	var levels []string

	if len(name) > 0 {
		levels = strings.Split(name, AccountStructureSep)
		if len(levels) > structureLevels-1 {
			return nil, app.ErrInvalidAccountStructure
		}
	}

	account := &AccountPath{
		TotalLevels: len(levels),
	}

	if account.TotalLevels > 0 {
		for _, v := range levels {
			if len(v) == 0 {
				return nil, app.ErrInvalidAccountStructure
			}
		}
	}

	if account.TotalLevels >= 1 {
		var err error
		account.Class, err = NewAccountClassFromString(levels[classLevel])
		if err != nil {
			return nil, app.ErrInvalidAccountStructure
		}
	}

	if account.TotalLevels >= 2 {
		account.Group = levels[groupLevel]
	}

	if account.TotalLevels >= 3 {
		account.Subgroup = levels[subgroupLevel]
	}

	return account, nil
}

func (a AccountPath) Name() string {
	if a.TotalLevels == 0 {
		return ""
	}

	str := a.Class.String()
	if a.TotalLevels == 1 {
		return str
	}

	str += AccountStructureSep + a.Group
	if a.TotalLevels == 2 {
		return str
	}

	str += AccountStructureSep + a.Subgroup
	return str
}
