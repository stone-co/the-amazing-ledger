package vos

import (
	"strings"

	"github.com/stone-co/the-amazing-ledger/app"
)

// AccountPath could have between 0 (empty) and n levels in her structure: "class.subgroup.account...".
// And, as AccountName, AccountPath has pre-defined classes.
//
// Some examples:
//   - ""
//   - "liability"
//   - "liability.clients"
//   - "liability.clients.available"

type AccountPath struct {
	Class       AccountClass
	Subgroup    string
	TotalLevels int
}

func NewAccountPath(name string) (*AccountPath, error) {
	name = strings.ToLower(name)

	var levels []string

	if len(name) > 0 {
		levels = strings.SplitN(name, AccountStructureSep, maximumStructureLevels)
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

	if len(levels) >= 1 {
		var err error
		account.Class, err = NewAccountClassFromString(levels[classLevel])
		if err != nil {
			return nil, app.ErrInvalidAccountStructure
		}
	}

	if len(levels) >= 2 {
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

	str += AccountStructureSep + a.Subgroup
	return str
}
