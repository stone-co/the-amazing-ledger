package vos

import (
	"strings"

	"github.com/stone-co/the-amazing-ledger/app"
)

const (
	AccountStructureSep    = "."
	minimumStructureLevels = 3
	maximumStructureLevels = 4
)

const (
	classLevel uint8 = iota
	subgroupLevel
	accountLevel
	suffixLevel
)

// AccountName must have at least 3 levels in its structure: "class.subgroup.account", where:
//   - class could be liability, assets, income, expense or equity;
//   - group, subgroup and id are "free text";
//   - "." must be used as a separator.
//
// AccountName can have a fourth level in its structure: "class.subgroup.account.suffix", where:
//   - suffix is "free text";
//   - it is considered a suffix everything after the id
//
// Some examples:
//   - "assets.conta_liquidacao.tesouraria"
//   - "liability.available.96a131a8-c4ac-495e-8971-fcecdbdd003a"
//   - "liability.available.96a131a8-c4ac-495e-8971-fcecdbdd003a.some_detail"
//   - "liability.clients.available.96a131a8-c4ac-495e-8971-fcecdbdd003a.detail1.detail2"
type AccountName struct {
	Class    AccountClass
	Subgroup string
	Account  string
	Suffix   string
}

func NewAccountName(name string) (AccountName, error) {
	name = strings.ToLower(name)

	levels := strings.SplitN(name, AccountStructureSep, maximumStructureLevels)
	if len(levels) < minimumStructureLevels {
		return AccountName{}, app.ErrInvalidAccountStructure
	}

	for _, v := range levels {
		if len(v) == 0 {
			return AccountName{}, app.ErrInvalidAccountStructure
		}
	}

	accountClass, err := NewAccountClassFromString(levels[classLevel])
	if err != nil {
		return AccountName{}, app.ErrInvalidAccountStructure
	}

	suffix := ""
	if len(levels) > minimumStructureLevels {
		suffix = levels[suffixLevel]
	}

	return AccountName{
		Class:    accountClass,
		Subgroup: levels[subgroupLevel],
		Account:  levels[accountLevel],
		Suffix:   suffix,
	}, nil
}

func (a AccountName) Name() string {
	return FormatAccount(a.Class.String(), a.Subgroup, a.Account, a.Suffix)
}

func FormatAccount(class, subgroup, account, suffix string) string {
	name := class + AccountStructureSep + subgroup + AccountStructureSep + account
	if suffix != "" {
		name += AccountStructureSep + suffix
	}
	return name
}
