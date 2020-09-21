package entities

import "sync/atomic"

type Version uint64

const (
	AnyAccountVersion Version = 0
	NewAccountVersion Version = 1
)

func (v Version) Current() Version {
	return v
}

func (v *Version) Next() Version {
	return Version(atomic.AddUint64((*uint64)(v), 1))
}
