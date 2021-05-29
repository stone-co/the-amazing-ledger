package vos

import (
	"strings"

	"github.com/stone-co/the-amazing-ledger/app"
)

// AccountQuery can have depth between 0 (empty) and n in her structure: "class.subgroup.account...".
// And, as AccountPath, AccountQuery respects the depth configurations.
//
// Some examples:
//   - ""
//   - "liability"
//   - "liability.clients"
//   - "liability.clients.available"
type AccountQuery struct {
	path string
}

func NewAccountQuery(query string) (AccountQuery, error) {
	query = strings.ToLower(query)

	components := strings.Split(query, _defaultConfig.DepthSeparator)
	size := len(components)

	if size == 0 {
		return AccountQuery{}, app.ErrInvalidAccountStructure
	}

	for i, component := range components {
		if len(component) == 0 {
			return AccountQuery{}, app.ErrInvalidAccountStructure
		}

		config, ok := _defaultConfig.DepthConfigs[i]
		if !ok {
			continue
		}

		if _, ok = config.Restrictions[component]; !ok {
			return AccountQuery{}, app.ErrInvalidAccountStructure
		}
	}

	return AccountQuery{query}, nil
}

func (q AccountQuery) Value() string {
	return q.path
}
