package persistent

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/slvic/nats-service/internal/types"
)

type Database struct {
	pg *sql.DB
}

func New(driver string, connection string) (*Database, error) {
	db, err := sql.Open(driver, connection)
	if err != nil {
		return nil, fmt.Errorf("could not open db: %s", err.Error())
	}

	database := &Database{pg: db}
	return database, nil
}

func (d *Database) SaveOrUpdate(ctx context.Context, order types.Order, rawOrder []byte) error {
	queryString := `INSERT INTO orders (id, data) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET data = EXCLUDED.data`
	values := []interface{}{order.Uid, rawOrder}

	_, err := d.pg.QueryContext(ctx, queryString, values...)
	if err != nil {
		return fmt.Errorf("could not execute save or update query: %s", err.Error())
	}
	return nil
}

func (d *Database) GetAll(ctx context.Context) ([]types.Message, error) {
	var orders []types.Message
	dbOrder := new(types.Message)
	queryString := `SELECT * FROM orders`

	rows, err := d.pg.QueryContext(ctx, queryString)
	if err != nil {
		return nil, fmt.Errorf("could not execute get all query: %s", err.Error())
	}
	defer func(rows *sql.Rows) {
		rowsCloseErr := rows.Close()
		if rowsCloseErr != nil {
			err = fmt.Errorf("could not close rows: %s", rowsCloseErr.Error())
		}
	}(rows)

	for rows.Next() {
		err := rows.Scan(&dbOrder.Id, &dbOrder.Data)
		if err != nil {
			return nil, err
		}
		orders = append(orders, *dbOrder)
	}
	return orders, nil
}
