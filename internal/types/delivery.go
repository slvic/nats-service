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
		errors = append(errors, fmt.Errorf("delivery name required, got: %s", d.Name))
	}
	if d.Phone == "" {
		errors = append(errors, fmt.Errorf("delivery phone required, got: %s", d.Phone))
	}
	if d.Zip == "" {
		errors = append(errors, fmt.Errorf("delivery zip required, got: %s", d.Zip))
	}
	if d.City == "" {
		errors = append(errors, fmt.Errorf("delivery city required, got: %s", d.City))
	}
	if d.Address == "" {
		errors = append(errors, fmt.Errorf("delivery address required, got: %s", d.Address))
	}
	if d.Region == "" {
		errors = append(errors, fmt.Errorf("delivery region required, got: %s", d.Region))
	}
	if d.Email == "" {
		errors = append(errors, fmt.Errorf("delivery email required, got: %s", d.Email))
	}
	return errors
}
