package app

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/slvic/nats-service/internal/http"
	worker "github.com/slvic/nats-service/internal/nats"
	"github.com/slvic/nats-service/internal/service/deliveries"
	"github.com/slvic/nats-service/internal/store/memory"
	"github.com/slvic/nats-service/internal/store/persistent"
	"go.uber.org/zap"
)

func Initialize(ctx context.Context, logger *zap.Logger, store *memory.Store) error {
	// add config
	database, err := persistent.New("postgres", "user=postgres dbname=nats-service password=postgres port=5433 sslmode=disable")
	if err != nil {
		return fmt.Errorf("could not create new database: %s", err.Error())
	}

	err = LoadOrdersFromDB(database, store)
	if err != nil {
		return fmt.Errorf("could not load orders from db: %s", err.Error())
	}

	storeService := deliveries.New(store, database)

	newWorker, err := worker.New(nats.DefaultURL, logger)
	if err != nil {
		return fmt.Errorf("could not create new worker: ")
	}
	ordersHandler := worker.NewOrdersHandler(storeService)
	err = newWorker.AddWorker("ORDERS.*", ordersHandler)
	if err != nil {
		return err
	}
	go func() {
		natsErr := newWorker.Run(ctx)
		if natsErr != nil {
			err = natsErr
		}
	}()
	if err != nil {
		return err
	}

	server := http.New(storeService, logger)
	err = server.Start()
	if err != nil {
		return err
	}
	return nil
}

func LoadOrdersFromDB(db *persistent.Database, store *memory.Store) error {
	allMessages, err := db.GetAll()
	if err != nil {
		return fmt.Errorf("could not get all stream messages: %s", err.Error())
	}

	for _, message := range allMessages {
		store.Set(message.Id, message.Data)
	}
	return nil
}
