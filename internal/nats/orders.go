package nats

import (
	"context"
	"encoding/json"

	"github.com/slvic/nats-service/internal/types"
	"go.uber.org/zap"
)

type OrdersService interface {
	SaveOrUpdate(ctx context.Context, order types.Order, rawOrder []byte) error
}

type OrdersHandler struct {
	ordersService OrdersService
	logger        *zap.Logger
}

func NewOrdersHandler(ordersService OrdersService, logger *zap.Logger) *OrdersHandler {
	return &OrdersHandler{
		ordersService: ordersService,
		logger:        logger,
	}
}

func (o *OrdersHandler) Handle(ctx context.Context, message []byte) error {
	if len(message) == 0 {
		return nil
	}
	var newOrder = types.Order{}
	err := json.Unmarshal(message, &newOrder)
	if err != nil {
		o.logger.Error("could not unmarshal a message", zap.Error(err))
		return nil
	}

	if errors := newOrder.Validate(); errors != nil {
		o.logger.Error("could not validate a message", zap.Error(err))
		return nil
	}

	err = o.ordersService.SaveOrUpdate(ctx, newOrder, message)
	return err
}
