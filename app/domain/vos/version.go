package vos

type Version int64

const (
	IgnoreAccountVersion Version = -1
	NextAccountVersion   Version = 0
)

func (v Version) AsInt64() int64 {
	return int64(v)
}
