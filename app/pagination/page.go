package pagination

import (
	"encoding/base64"
	"encoding/json"

	"github.com/stone-co/the-amazing-ledger/app"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

const (
	_defaultPageSize = 10
	_maxPageSize     = 50
)

type Page struct {
	Size   int
	Cursor Cursor
}

func NewPage(p *proto.RequestPagination) (Page, error) {
	if p == nil {
		return Page{Size: _defaultPageSize}, nil
	}

	if p.PageSize == 0 {
		return Page{}, app.ErrInvalidPageSize
	}

	size := int(p.PageSize)
	if size > _maxPageSize {
		size = _maxPageSize
	}

	var (
		data []byte
		err  error
	)

	if p.PageToken != "" {
		data, err = base64.StdEncoding.DecodeString(p.PageToken)
		if err != nil {
			return Page{}, app.ErrInvalidPageCursor
		}
	}

	return Page{
		Size:   size,
		Cursor: data,
	}, nil
}

func (p Page) Extract(data interface{}) error {
	err := json.Unmarshal(p.Cursor, data)
	if err != nil {
		return app.ErrInvalidPageCursor
	}

	return nil
}
