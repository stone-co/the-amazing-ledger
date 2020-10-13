package entities

import (
	"sync"
)

type CachedAccountInfo struct {
	sync.Mutex
	Version Version
}

type CachedAccounts struct {
	objects *sync.Map
}

func NewCachedAccounts() *CachedAccounts {
	return &CachedAccounts{
		objects: new(sync.Map),
	}
}

func (c *CachedAccounts) LoadOrStore(accountID string) *CachedAccountInfo {
	object := &CachedAccountInfo{
		Version: NewAccountVersion,
	}

	objectInMap, _ := c.objects.LoadOrStore(accountID, object)
	return objectInMap.(*CachedAccountInfo)
}

func (c *CachedAccounts) Store(accountID string, version Version) {
	object := &CachedAccountInfo{
		Version: version,
	}

	c.objects.Store(accountID, object)
}
