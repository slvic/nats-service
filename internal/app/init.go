package app

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/slvic/nats-service/internal/http"
	worker "github.com/slvic/nats-service/internal/nats"
	"github.com/slvic/nats-service/internal/service/deliveries"
	"github.com/slvic/nats-service/internal/store/memory"
	"go.uber.org/zap"
)

func Initialize(ctx context.Context, logger *zap.Logger, store *memory.Store) error {
	storeService := deliveries.New(store)

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
