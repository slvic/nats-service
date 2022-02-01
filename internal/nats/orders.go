package nats

import (
	"encoding/json"
	"fmt"

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
		return fmt.Errorf("data is empty")
	}

	var newOrder = types.Order{}
	err := json.Unmarshal(message, &newOrder)
	if err != nil {
		return fmt.Errorf("could not unmarshal data: %v", err)
	}

	if errors := newOrder.Validate(); errors != nil {
		return fmt.Errorf("could not validate data: %v", errors)
	}

	err = o.ordersService.SaveOrUpdate(newOrder)
	return err
}
