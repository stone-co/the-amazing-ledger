package vo

import "sync/atomic"

type Version uint64

const (
	AnyAccountVersion Version = 0
	NewAccountVersion Version = 1
)

func (v Version) Current() Version {
	return v
}

func (v Version) ToUInt64() uint64 {
	return uint64(v)
}

func (v *Version) Next() Version {
	return Version(atomic.AddUint64((*uint64)(v), 1))
}
