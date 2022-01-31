package persistent

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
)

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

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

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

type Database struct {
	pg *sql.DB
}

func NewDb(driver string, connection string) (*Database, error) {
	db, err := sql.Open(driver, connection)
	if err != nil {
		return nil, fmt.Errorf("could not open db: %v", err)
	}

	database := &Database{pg: db}
	return database, nil
}

func (o *Order) Validate() []error {
	var errors []error
	if o.Uid == "" {
		errors = append(errors, fmt.Errorf("order uid required, got: %v", o.Uid))
	}
	if o.TrackNumber == "" {
		errors = append(errors, fmt.Errorf("order track number required, got: %v", o.TrackNumber))
	}
	if o.Entry == "" {
		errors = append(errors, fmt.Errorf("order entry required, got: %v", o.Entry))
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
		errors = append(errors, fmt.Errorf("order locale required, got: %v", o.Locale))
	}
	if o.InternalSignature == "" {
		errors = append(errors, fmt.Errorf("order internal signature required, got: %v", o.InternalSignature))
	}
	if o.CustomerID == "" {
		errors = append(errors, fmt.Errorf("order customer ID required, got: %v", o.CustomerID))
	}
	if o.DeliveryService == "" {
		errors = append(errors, fmt.Errorf("order delivery service required, got: %v", o.DeliveryService))
	}
	if o.Shardkey == "" {
		errors = append(errors, fmt.Errorf("order shard key required, got: %v", o.Shardkey))
	}
	if o.SmID == 0 {
		errors = append(errors, fmt.Errorf("order smID required, got: %v", o.SmID))
	}
	if o.DateCreated == "" {
		errors = append(errors, fmt.Errorf("order date created required, got: %v", o.DateCreated))
	}
	if o.OofShard == "" {
		errors = append(errors, fmt.Errorf("order oof shard created required, got: %v", o.OofShard))
	}
	return errors
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
func (i *Item) Validate() []error {
	var errors []error
	if i.ChrtID == 0 {
		errors = append(errors, fmt.Errorf("item chrt ID required, got: %v", i.ChrtID))
	}
	if i.TrackNumber == "" {
		errors = append(errors, fmt.Errorf("item track number required, got: %v", i.TrackNumber))
	}
	if i.Price == 0 {
		errors = append(errors, fmt.Errorf("item price required, got: %v", i.Price))
	}
	if i.Rid == "" {
		errors = append(errors, fmt.Errorf("item rid required, got: %v", i.Rid))
	}
	if i.Name == "" {
		errors = append(errors, fmt.Errorf("item name required, got: %v", i.Name))
	}
	if i.Sale == 0 {
		errors = append(errors, fmt.Errorf("item sale required, got: %v", i.Sale))
	}
	if i.Size == "" {
		errors = append(errors, fmt.Errorf("item size required, got: %v", i.Size))
	}
	if i.TotalPrice == 0 {
		errors = append(errors, fmt.Errorf("item total price required, got: %v", i.TotalPrice))
	}
	if i.NmID == 0 {
		errors = append(errors, fmt.Errorf("item nm ID required, got: %v", i.NmID))
	}
	if i.Brand == "" {
		errors = append(errors, fmt.Errorf("item brand required, got: %v", i.Brand))
	}
	if i.Status == 0 {
		errors = append(errors, fmt.Errorf("item status required, got: %v", i.Status))
	}
	return errors
}
func (p *Payment) Validate() []error {
	var errors []error
	if p.Transaction == "" {
		errors = append(errors, fmt.Errorf("payment transaction required, got: %v", p.Transaction))
	}
	if p.RequestID == "" {
		errors = append(errors, fmt.Errorf("payment request ID required, got: %v", p.RequestID))
	}
	if p.Currency == "" {
		errors = append(errors, fmt.Errorf("payment currency required, got: %v", p.Currency))
	}
	if p.Provider == "" {
		errors = append(errors, fmt.Errorf("payment provider required, got: %v", p.Provider))
	}
	if p.Amount == 0 {
		errors = append(errors, fmt.Errorf("payment amount required, got: %v", p.Amount))
	}
	if p.PaymentDt == 0 {
		errors = append(errors, fmt.Errorf("payment payment Dt required, got: %v", p.PaymentDt))
	}
	if p.Bank == "" {
		errors = append(errors, fmt.Errorf("payment bank required, got: %v", p.Bank))
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

func unmarshalAndValidate(data []string) ([]Order, error) {
	if data == nil {
		return nil, fmt.Errorf("data is empty")
	}
	var orders []Order
	var newOrder = &Order{}
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

func (d *Database) SaveOrUpdateMany(data []string) error {
	orders, err := unmarshalAndValidate(data)
	if err != nil {
		return fmt.Errorf("could not unmarshal or validate while save or update many: %v", err)
	}

	queryString := `INSERT INTO order(id, data) VALUES `
	values := []interface{}{}
	for _, row := range orders {
		queryString += `(?, ?),`
		values = append(values, row.Uid, row)
	}
	queryString = queryString[:len(queryString)-1]
	queryString += `ON CONFLICT (id) DO UPDATE SET data = EXCLUDED.data`

	_, err = d.pg.Query(queryString, values)
	if err != nil {
		return fmt.Errorf("could not execute save or update many query: %v", err)
	}
	return nil
}

func (d *Database) SaveOrUpdate(data string) error {
	order, err := unmarshalAndValidate([]string{data})
	if err != nil {
		return fmt.Errorf("could not unmarshal or validate while save or update: %v", err)
	}

	queryString := `INSERT INTO order(id, data) VALUES (?, ?) ON CONFLICT (id) DO UPDATE SET data = EXCLUDED.data`
	values := []interface{}{order[0].Uid, order[0]}

	_, err = d.pg.Query(queryString, values)
	if err != nil {
		return fmt.Errorf("could not execute save or update query: %v", err)
	}
	return nil
}

func (d *Database) GetAll() ([]Order, []error) {
	var errors []error
	var orders []Order
	var dbOrder *Order
	queryString := `SELECT data FROM orders`

	rows, err := d.pg.Query(queryString)
	if err != nil {
		return nil, []error{fmt.Errorf("could not execute get all query: %v", err)}
	}

	for rows.Next() {
		err := rows.Scan(dbOrder)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		orders = append(orders, *dbOrder)
	}
	return orders, nil
}
