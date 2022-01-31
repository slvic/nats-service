package persistent

import (
	"encoding/json"
	"fmt"
	"github.com/slvic/nats-service/internal/store/types"
)

func UnmarshalAndValidate(data []string) ([]types.Order, error) {
	if data == nil {
		return nil, fmt.Errorf("data is empty")
	}
	var orders []types.Order
	var newOrder = &types.Order{}
	for _, dataRow := range data {
		err := json.Unmarshal([]byte(dataRow), newOrder)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal data row error: %v", err)
		}
		if errors := newOrder.Validate(); errors != nil {
			return nil, fmt.Errorf("could not validate data: %v", errors)
		}
		orders = append(orders, *newOrder)
	}
	return orders, nil
}
