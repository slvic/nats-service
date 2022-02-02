package deliveries

import (
	"fmt"

	"github.com/slvic/nats-service/internal/store/memory"
	"github.com/slvic/nats-service/internal/types"
)

type Deliverer struct {
	store *memory.Store
}

func New(store *memory.Store) *Deliverer {
	return &Deliverer{
		store: store,
	}
}
func (d *Deliverer) SaveOrUpdate(order types.Order) error {
	d.store.Set(order.Uid, order)
	// possible db failure
	return nil
}

func (d *Deliverer) GetMessageById(id string) (types.Order, error) {
	order, found := d.store.Get(id)
	if !found {
		return types.Order{}, fmt.Errorf("message not found")
	}
	return order, nil
}
