package pagination

import (
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

func TestNewPage(t *testing.T) {
	tests := []struct {
		name        string
		p           *proto.RequestPagination
		want        Page
		expectedErr error
	}{
		{
			name: "correct case - empty token",
			p: &proto.RequestPagination{
				PageSize:  15,
				PageToken: "",
			},
			want: Page{
				Size:   15,
				Cursor: nil,
			},
			expectedErr: nil,
		},
		{
			name: "correct case - filled token",
			p: &proto.RequestPagination{
				PageSize:  10,
				PageToken: base64.StdEncoding.EncodeToString([]byte("abc")),
			},
			want: Page{
				Size:   10,
				Cursor: []byte("abc"),
			},
			expectedErr: nil,
		},
		{
			name: "page size greater than maximum",
			p: &proto.RequestPagination{
				PageSize:  100,
				PageToken: "",
			},
			want: Page{
				Size:   50,
				Cursor: nil,
			},
			expectedErr: nil,
		},
		{
			name: "empty case",
			p:    nil,
			want: Page{
				Size:   10,
				Cursor: nil,
			},
			expectedErr: nil,
		},
		{
			name: "invalid page size",
			p: &proto.RequestPagination{
				PageSize:  0,
				PageToken: "",
			},
			want:        Page{},
			expectedErr: app.ErrInvalidPageSize,
		},
		{
			name: "invalid page token",
			p: &proto.RequestPagination{
				PageSize:  10,
				PageToken: "!@#$%^&*()",
			},
			want:        Page{},
			expectedErr: app.ErrInvalidPageCursor,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPage(tt.p)

			assert.ErrorIs(t, err, tt.expectedErr)
			if tt.expectedErr != nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestPage_Extract(t *testing.T) {
	type data struct {
		Abc time.Time `json:"abc"`
	}

	dt := data{Abc: time.Now().UTC().Round(time.Nanosecond)}
	jsonStr, err := json.Marshal(dt)
	assert.NoError(t, err)

	tests := []struct {
		name        string
		page        Page
		expectedErr error
	}{
		{
			name: "correct case",
			page: Page{
				Size:   10,
				Cursor: jsonStr,
			},
			expectedErr: nil,
		},
		{
			name: "invalid format",
			page: Page{
				Size:   10,
				Cursor: []byte("invalid"),
			},
			expectedErr: app.ErrInvalidPageCursor,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d data
			err := tt.page.Extract(&d)
			assert.ErrorIs(t, err, tt.expectedErr)
			if err == nil {
				assert.Equal(t, d.Abc, dt.Abc)
			}
		})
	}
}
