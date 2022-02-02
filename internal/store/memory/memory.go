package memory

import (
	"sync"
)

type Store struct {
	sync.RWMutex
	items map[string][]byte
}

func New() *Store {

	items := make(map[string][]byte)

	store := Store{
		items: items,
	}

	return &store
}

func (c *Store) Set(key string, value []byte) {
	c.Lock()

	defer c.Unlock()

	c.items[key] = value
}

func (c *Store) Get(key string) ([]byte, bool) {

	c.RLock()

	defer c.RUnlock()

	item, found := c.items[key]

	if !found {
		return nil, false
	}

	return item, true
}
