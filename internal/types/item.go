package types

import "fmt"

type Item struct {
	ChrtID      int64  `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int64  `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int64  `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int64  `json:"total_price"`
	NmID        int64  `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int64  `json:"status"`
}

func (i *Item) Validate() []error {
	var errors []error
	if i.ChrtID == 0 {
		errors = append(errors, fmt.Errorf("item chrt ID required, got: %v", i.ChrtID))
	}
	if i.TrackNumber == "" {
		errors = append(errors, fmt.Errorf("item track number required, got: %s", i.TrackNumber))
	}
	if i.Price == 0 {
		errors = append(errors, fmt.Errorf("item price required, got: %v", i.Price))
	}
	if i.Rid == "" {
		errors = append(errors, fmt.Errorf("item rid required, got: %s", i.Rid))
	}
	if i.Name == "" {
		errors = append(errors, fmt.Errorf("item name required, got: %s", i.Name))
	}
	if i.Sale == 0 {
		errors = append(errors, fmt.Errorf("item sale required, got: %v", i.Sale))
	}
	if i.Size == "" {
		errors = append(errors, fmt.Errorf("item size required, got: %s", i.Size))
	}
	if i.TotalPrice == 0 {
		errors = append(errors, fmt.Errorf("item total price required, got: %v", i.TotalPrice))
	}
	if i.NmID == 0 {
		errors = append(errors, fmt.Errorf("item nm ID required, got: %v", i.NmID))
	}
	if i.Brand == "" {
		errors = append(errors, fmt.Errorf("item brand required, got: %s", i.Brand))
	}
	if i.Status == 0 {
		errors = append(errors, fmt.Errorf("item status required, got: %v", i.Status))
	}
	return errors
}
