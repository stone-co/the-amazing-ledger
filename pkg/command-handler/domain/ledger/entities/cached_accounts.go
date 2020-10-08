package entities

import (
	"sync"

	"github.com/google/uuid"
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

func (c *CachedAccounts) LoadOrStore(accountID uuid.UUID) *CachedAccountInfo {
	object := &CachedAccountInfo{
		Version: NewAccountVersion,
	}

	objectInMap, _ := c.objects.LoadOrStore(accountID, object)
	return objectInMap.(*CachedAccountInfo)
}

func (c *CachedAccounts) Store(accountID uuid.UUID, version Version) {
	object := &CachedAccountInfo{
		Version: version,
	}

	c.objects.Store(accountID, object)
}
