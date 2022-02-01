package memory

import (
	"errors"
	"sync"
)

type Store struct {
	sync.RWMutex
	items map[string]Item
}

type Item struct {
	Value interface{}
}

func NewStore() *Store {

	items := make(map[string]Item)

	store := Store{
		items: items,
	}

	return &store
}

func (c *Store) Set(key string, value interface{}) {
	c.Lock()

	defer c.Unlock()

	c.items[key] = Item{
		Value: value,
	}

}

func (c *Store) Get(key string) (interface{}, bool) {

	c.RLock()

	defer c.RUnlock()

	item, found := c.items[key]

	if !found {
		return nil, false
	}

	return item.Value, true
}

func (c *Store) Delete(key string) error {

	c.Lock()

	defer c.Unlock()

	if _, found := c.items[key]; !found {
		return errors.New("Key not found")
	}

	delete(c.items, key)

	return nil
}
