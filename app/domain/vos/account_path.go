package vos

import (
	"regexp"
	"strings"

	"github.com/stone-co/the-amazing-ledger/app"
)

type DepthConfig struct {
	Restrictions map[string]struct{}
	Name         string
}

type AccountConfig struct {
	MinimumDepth   int
	MaximumDepth   int
	DepthConfigs   map[int]DepthConfig
	DepthSeparator string
}

var _empty = struct{}{}

var _defaultConfig = AccountConfig{
	MinimumDepth: 3,
	MaximumDepth: 0,
	DepthConfigs: map[int]DepthConfig{
		0: {
			Restrictions: map[string]struct{}{
				"liability": _empty,
				"assets":    _empty,
				"income":    _empty,
				"expense":   _empty,
				"equity":    _empty,
			},
			Name: "class",
		},
	},
	DepthSeparator: ".",
}

var regexOnlyAlphanumericAndUnderscore = regexp.MustCompile(`^[a-zA-Z0-9_]*$`)

const maxLabelLength = 255

// TODO: better docs

// AccountPath can be as deep as needed, limited by AccountConfig MinimumDepth and MaximumDepth.
// None of the values can be '' (empty string), more than 255 characters or characteres other than
// alphanumeric and underscore (_).
// Depth restrictions can be applied by using DepthConfig. The default configuration for example:
//	- the first depth is called class
// 	- it can only be one of the following:
//		- liability
//		- assets
//		- income
//		- expense
//		- equity
//	- account need to have a depth of at least 3
//	- there are no upper limits
// 	- '.' must be used as a separator.
//
// Some examples:
//   - assets.account.treasury
//   - liability.available.96a131a8_c4ac_495e_8971_fcecdbdd003a
//   - liability.available.96a131a8_c4ac_495e_8971_fcecdbdd003a.some_detail
//   - liability.clients.available.96a131a8_c4ac_495e_8971_fcecdbdd003a.detail1.detail2
type AccountPath struct {
	path string
}

func NewAccountPath(path string) (AccountPath, error) {
	path = strings.ToLower(path)

	components := strings.Split(path, _defaultConfig.DepthSeparator)
	size := len(components)

	if size == 0 {
		return AccountPath{}, app.ErrInvalidAccountStructure
	}

	if _defaultConfig.MaximumDepth != 0 && size > _defaultConfig.MaximumDepth {
		return AccountPath{}, app.ErrInvalidAccountStructure
	}

	if size < _defaultConfig.MinimumDepth {
		return AccountPath{}, app.ErrInvalidAccountStructure
	}

	for i, component := range components {
		if component == "" || len(component) > maxLabelLength {
			return AccountPath{}, app.ErrInvalidAccountStructure
		}

		if !regexOnlyAlphanumericAndUnderscore.MatchString(component) {
			return AccountPath{}, app.ErrInvalidAccountStructure
		}

		config, ok := _defaultConfig.DepthConfigs[i]
		if !ok {
			continue
		}

		if _, ok = config.Restrictions[component]; !ok {
			return AccountPath{}, app.ErrInvalidAccountStructure
		}
	}

	return AccountPath{path}, nil
}

func (a AccountPath) Name() string {
	return a.path
}
