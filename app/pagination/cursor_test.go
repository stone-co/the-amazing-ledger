package pagination

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app"
)

func TestNewCursor(t *testing.T) {
	tests := []struct {
		name        string
		arg         interface{}
		want        func(interface{}) (Cursor, error)
		expectedErr error
	}{
		{
			name: "valid case",
			arg:  map[string]interface{}{"abc": 123},
			want: func(v interface{}) (Cursor, error) {
				return json.Marshal(v)
			},
			expectedErr: nil,
		},
		{
			name: "invalid case",
			arg:  make(chan int),
			want: func(v interface{}) (Cursor, error) {
				return nil, nil
			},
			expectedErr: app.ErrInvalidPageCursor,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCursor(tt.arg)
			assert.ErrorIs(t, err, tt.expectedErr)

			w, err := tt.want(tt.arg)
			assert.NoError(t, err)
			assert.Equal(t, w, got)
		})
	}
}
