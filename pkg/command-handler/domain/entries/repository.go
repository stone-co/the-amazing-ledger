package entries

import "github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/entries/entities"

type Repository interface {
	Create(*[]entities.Entry) error
}
