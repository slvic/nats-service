package deliveries

import (
	"fmt"

	"github.com/slvic/nats-service/internal/store/memory"
	"github.com/slvic/nats-service/internal/store/persistent"
	"github.com/slvic/nats-service/internal/types"
)

type Memory interface {
	Get(key string) ([]byte, bool)
	Set(key string, value []byte)
}

type Database interface {
	SaveOrUpdate(order types.Order, rawOrder []byte) error
	GetAll() ([]types.Message, error)
}

type Deliverer struct {
	store Memory
	db    Database
}

func New(store *memory.Store, database *persistent.Database) *Deliverer {
	return &Deliverer{
		store: store,
		db:    database,
	}
}
func (d *Deliverer) SaveOrUpdate(order types.Order, rawOrder []byte) error {
	d.store.Set(order.Uid, rawOrder)
	err := d.db.SaveOrUpdate(order, rawOrder)
	if err != nil {
		return fmt.Errorf("could not save or update order: %s", err.Error())
	}
	return nil
}

func (d *Deliverer) GetMessageById(id string) ([]byte, error) {
	order, found := d.store.Get(id)
	if !found {
		return nil, fmt.Errorf("message not found")
	}
	return order, nil
}
