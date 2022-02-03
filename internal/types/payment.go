package types

import "fmt"

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int64  `json:"amount"`
	PaymentDt    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int64  `json:"delivery_cost"`
	GoodsTotal   int64  `json:"goods_total"`
	CustomFee    int64  `json:"custom_fee"`
}

func (p *Payment) Validate() []error {
	var errors []error
	if p.Transaction == "" {
		errors = append(errors, fmt.Errorf("payment transaction required, got: %s", p.Transaction))
	}
	if p.RequestID == "" {
		errors = append(errors, fmt.Errorf("payment request ID required, got: %s", p.RequestID))
	}
	if p.Currency == "" {
		errors = append(errors, fmt.Errorf("payment currency required, got: %s", p.Currency))
	}
	if p.Provider == "" {
		errors = append(errors, fmt.Errorf("payment provider required, got: %s", p.Provider))
	}
	if p.Amount == 0 {
		errors = append(errors, fmt.Errorf("payment amount required, got: %v", p.Amount))
	}
	if p.PaymentDt == 0 {
		errors = append(errors, fmt.Errorf("payment payment Dt required, got: %v", p.PaymentDt))
	}
	if p.Bank == "" {
		errors = append(errors, fmt.Errorf("payment bank required, got: %s", p.Bank))
	}
	if p.DeliveryCost == 0 {
		errors = append(errors, fmt.Errorf("payment delivery cost required, got: %v", p.DeliveryCost))
	}
	if p.GoodsTotal == 0 {
		errors = append(errors, fmt.Errorf("payment goods total required, got: %v", p.GoodsTotal))
	}
	if p.CustomFee == 0 {
		errors = append(errors, fmt.Errorf("payment custom fee required, got: %v", p.CustomFee))
	}
	return errors
}
