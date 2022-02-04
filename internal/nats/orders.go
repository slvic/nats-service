package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/graft"
	"github.com/slvic/nats-service/internal/types"
	"go.uber.org/zap"
)

type OrdersService interface {
	SaveOrUpdate(ctx context.Context, node *graft.Node, order types.Order, rawOrder []byte) error
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

func (o *OrdersHandler) Handle(ctx context.Context, node *graft.Node, message []byte) error {
	if len(message) == 0 {
		return nil
	}
	fmt.Println("I Handle hehe: ", node.Id(), node.State().String())
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

	err = o.ordersService.SaveOrUpdate(ctx, node, newOrder, message)
	return err
}
