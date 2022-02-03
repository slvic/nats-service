package deliveries

import (
	"fmt"

	"github.com/slvic/nats-service/internal/store/memory"
	"github.com/slvic/nats-service/internal/store/persistent"
	"github.com/slvic/nats-service/internal/types"
	"go.uber.org/zap"
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
	store  Memory
	db     Database
	logger *zap.Logger
}

func New(store *memory.Store, database *persistent.Database, logger *zap.Logger) *Deliverer {
	return &Deliverer{
		store:  store,
		db:     database,
		logger: logger,
	}
}
func (d *Deliverer) SaveOrUpdate(order types.Order, rawOrder []byte) error {
	d.store.Set(order.Uid, rawOrder)
	err := d.db.SaveOrUpdate(order, rawOrder)
	if err != nil {
		d.logger.Error("database", zap.Error(err))
		return fmt.Errorf("could not save or update order: %s", err.Error())
	}
	return nil
}

func (d *Deliverer) GetMessageById(id string) ([]byte, error) {
	order, found := d.store.Get(id)
	if !found {
		d.logger.Error("wrong message id", zap.String("id", id))
		return nil, fmt.Errorf("message not found")
	}
	return order, nil
}
