package memory

import (
	"errors"
	"sync"

	"github.com/slvic/nats-service/internal/types"
)

type Store struct {
	sync.RWMutex
	items map[string]types.Order
}

func NewStore() *Store {

	items := make(map[string]types.Order)

	store := Store{
		items: items,
	}

	return &store
}

func (c *Store) Set(key string, value types.Order) {
	c.Lock()

	defer c.Unlock()

	c.items[key] = value
}

func (c *Store) Get(key string) (types.Order, bool) {

	c.RLock()

	defer c.RUnlock()

	item, found := c.items[key]

	if !found {
		return types.Order{}, false
	}

	return item, true
}

func (c *Store) Delete(key string) error {

	c.Lock()

	defer c.Unlock()

	if _, found := c.items[key]; !found {
		return errors.New("key not found")
	}

	delete(c.items, key)

	return nil
}
