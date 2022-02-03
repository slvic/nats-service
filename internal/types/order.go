package types

import "fmt"

type Order struct {
	Uid               string   `json:"order_uid"`
	TrackNumber       string   `json:"track_number"`
	Entry             string   `json:"entry"`
	Delivery          Delivery `json:"delivery"`
	Payment           Payment  `json:"payment"`
	Items             []Item   `json:"items"`
	Locale            string   `json:"locale"`
	InternalSignature string   `json:"internal_signature"`
	CustomerID        string   `json:"customer_id"`
	DeliveryService   string   `json:"delivery_service"`
	Shardkey          string   `json:"shardkey"`
	SmID              int64    `json:"sm_id"`
	DateCreated       string   `json:"date_created"`
	OofShard          string   `json:"oof_shard"`
}

func (o Order) Validate() []error {
	var errors []error
	if o.Uid == "" {
		errors = append(errors, fmt.Errorf("order uid required, got: %s", o.Uid))
	}
	if o.TrackNumber == "" {
		errors = append(errors, fmt.Errorf("order track number required, got: %s", o.TrackNumber))
	}
	if o.Entry == "" {
		errors = append(errors, fmt.Errorf("order entry required, got: %s", o.Entry))
	}
	if o.Delivery.Validate() != nil {
		errors = append(errors, fmt.Errorf("order delivery required, got: %v", o.Delivery))
	}
	if o.Payment.Validate() != nil {
		errors = append(errors, fmt.Errorf("order payment required, got: %v", o.Payment))
	}
	for _, item := range o.Items {
		if item.Validate() != nil {
			errors = append(errors, fmt.Errorf("order items required, got: %v", o.Items))
		}
	}
	if o.Locale == "" {
		errors = append(errors, fmt.Errorf("order locale required, got: %s", o.Locale))
	}
	if o.InternalSignature == "" {
		errors = append(errors, fmt.Errorf("order internal signature required, got: %s", o.InternalSignature))
	}
	if o.CustomerID == "" {
		errors = append(errors, fmt.Errorf("order customer ID required, got: %s", o.CustomerID))
	}
	if o.DeliveryService == "" {
		errors = append(errors, fmt.Errorf("order delivery service required, got: %s", o.DeliveryService))
	}
	if o.Shardkey == "" {
		errors = append(errors, fmt.Errorf("order shard key required, got: %s", o.Shardkey))
	}
	if o.SmID == 0 {
		errors = append(errors, fmt.Errorf("order smID required, got: %v", o.SmID))
	}
	if o.DateCreated == "" {
		errors = append(errors, fmt.Errorf("order date created required, got: %s", o.DateCreated))
	}
	if o.OofShard == "" {
		errors = append(errors, fmt.Errorf("order oof shard created required, got: %s", o.OofShard))
	}
	return errors
}
