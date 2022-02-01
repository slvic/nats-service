package types

import "fmt"

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

func (d *Delivery) Validate() []error {
	var errors []error
	if d.Name == "" {
		errors = append(errors, fmt.Errorf("delivery name required, got: %v", d.Name))
	}
	if d.Phone == "" {
		errors = append(errors, fmt.Errorf("delivery phone required, got: %v", d.Phone))
	}
	if d.Zip == "" {
		errors = append(errors, fmt.Errorf("delivery zip required, got: %v", d.Zip))
	}
	if d.City == "" {
		errors = append(errors, fmt.Errorf("delivery city required, got: %v", d.City))
	}
	if d.Address == "" {
		errors = append(errors, fmt.Errorf("delivery address required, got: %v", d.Address))
	}
	if d.Region == "" {
		errors = append(errors, fmt.Errorf("delivery region required, got: %v", d.Region))
	}
	if d.Email == "" {
		errors = append(errors, fmt.Errorf("delivery email required, got: %v", d.Email))
	}
	return errors
}
