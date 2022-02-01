package deliveries

import (
	"fmt"

	"github.com/slvic/nats-service/internal/store/memory"
)

type Deliverer struct {
	store *memory.Store
}

func NewDeliverer(store *memory.Store) *Deliverer {
	return &Deliverer{
		store: store,
	}
}

func (d *Deliverer) GetMessagesBySubject(subject string) (string, error) {
	_, found := d.store.Get(subject)
	if !found {
		return "", fmt.Errorf("subject not found")
	}

	//
	return "", nil
}
