package pagination

import (
	"encoding/base64"
	"encoding/json"

	"github.com/stone-co/the-amazing-ledger/app"
)

type Cursor []byte

func NewCursor(v interface{}) (Cursor, error) {
	jsonStr, err := json.Marshal(v)
	if err != nil {
		return nil, app.ErrInvalidPageCursor
	}

	return jsonStr, nil
}

func (c Cursor) Tokenize() string {
	return base64.StdEncoding.EncodeToString(c)
}
