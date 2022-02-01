package persistent

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/slvic/nats-service/internal/service"
	"github.com/slvic/nats-service/internal/types"
)

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

func (d *Database) SaveOrUpdateMany(data []string) error {
	orders, err := service.UnmarshalAndValidate(data)
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
	order, err := service.UnmarshalAndValidate([]string{data})
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

func (d *Database) GetAll() ([]types.Order, []error) {
	var errors []error
	var orders []types.Order
	var dbOrder *types.Order
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
