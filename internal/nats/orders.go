package nats

import (
	"encoding/json"

	"github.com/slvic/nats-service/internal/types"
)

type OrdersService interface {
	SaveOrUpdate(order types.Order) error
}

type OrdersHandler struct {
	ordersService OrdersService
}

func NewOrdersHandler(ordersService OrdersService) *OrdersHandler {
	return &OrdersHandler{
		ordersService: ordersService,
	}
}

func (o *OrdersHandler) Handle(message []byte) error {
	if len(message) == 0 {
		return nil
	}

	var newOrder = types.Order{}
	err := json.Unmarshal(message, &newOrder)
	if err != nil {
		// zap
		return nil
	}

	if errors := newOrder.Validate(); errors != nil {
		// zap
		return nil
	}

	err = o.ordersService.SaveOrUpdate(newOrder)
	return err
}
