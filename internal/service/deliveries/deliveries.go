package deliveries

import (
	"context"
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
	SaveOrUpdate(ctx context.Context, order types.Order, rawOrder []byte) error
	GetAll(ctx context.Context) ([]types.Message, error)
}

type LeaderElector interface {
	AmILeader() bool
}

type Deliverer struct {
	store   Memory
	db      Database
	elector LeaderElector
	logger  *zap.Logger
}

func New(store *memory.Store, database *persistent.Database, elector LeaderElector, logger *zap.Logger) *Deliverer {
	return &Deliverer{
		store:   store,
		db:      database,
		elector: elector,
		logger:  logger,
	}
}
func (d *Deliverer) SaveOrUpdate(ctx context.Context, order types.Order, rawOrder []byte) error {
	d.store.Set(order.Uid, rawOrder)
	if d.elector.AmILeader() {
		err := d.db.SaveOrUpdate(ctx, order, rawOrder)
		if err != nil {
			d.logger.Error("database", zap.Error(err))
			return fmt.Errorf("could not save or update order: %s", err.Error())
		}
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
