package service

import (
	"encoding/json"
	"fmt"

	"github.com/slvic/nats-service/internal/types"
)

func UnmarshalAndValidate(data []byte) (types.Order, error) {
	if data == nil {
		return types.Order{}, fmt.Errorf("data is empty")
	}
	var newOrder = &types.Order{}
	err := json.Unmarshal(data, newOrder)
	if err != nil {
		return types.Order{}, fmt.Errorf("could not unmarshal data: %v", err)
	}

	if errors := newOrder.Validate(); errors != nil {
		return types.Order{}, fmt.Errorf("could not validate data: %v", errors)
	}
	return *newOrder, nil
}
