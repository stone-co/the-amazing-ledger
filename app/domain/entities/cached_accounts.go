package entities

import (
	"sync"

	"github.com/stone-co/the-amazing-ledger/app/domain/vo"
)

type CachedAccountInfo struct {
	sync.Mutex
	CurrentVersion vo.Version
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
		CurrentVersion: vo.NewAccountVersion,
	}

	objectInMap, _ := c.objects.LoadOrStore(accountID, object)
	return objectInMap.(*CachedAccountInfo)
}

func (c *CachedAccounts) Store(accountID string, version vo.Version) {
	object := &CachedAccountInfo{
		CurrentVersion: version,
	}

	c.objects.Store(accountID, object)
}
