package app

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/slvic/nats-service/configs"
	"github.com/slvic/nats-service/internal/http"
	worker "github.com/slvic/nats-service/internal/nats"
	"github.com/slvic/nats-service/internal/service/deliveries"
	"github.com/slvic/nats-service/internal/store/memory"
	"github.com/slvic/nats-service/internal/store/persistent"
	"go.uber.org/zap"
)

func Initialize(ctx context.Context, logger *zap.Logger, store *memory.Store, dbConfig *configs.DatabaseConfig) error {
	dbConnectionString := fmt.Sprintf("user=%s dbname=%s password=%s port=%s sslmode=%s",
		dbConfig.DBUser,
		dbConfig.DBName,
		dbConfig.DBPassword,
		dbConfig.DBPort,
		dbConfig.SSLMode)
	database, err := persistent.New("postgres", dbConnectionString)
	fmt.Println(dbConnectionString)
	if err != nil {
		return fmt.Errorf("could not create new database: %s", err.Error())
	}

	err = LoadOrdersFromDB(database, store, logger)
	if err != nil {
		return fmt.Errorf("could not load orders from db: %s", err.Error())
	}

	storeService := deliveries.New(store, database, logger)

	newWorker, err := worker.New(nats.DefaultURL, logger)
	if err != nil {
		return fmt.Errorf("could not create new worker: %s", err.Error())
	}
	ordersHandler := worker.NewOrdersHandler(storeService, logger)
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

func LoadOrdersFromDB(db *persistent.Database, store *memory.Store, logger *zap.Logger) error {
	allMessages, err := db.GetAll()
	if err != nil {
		logger.Error("database", zap.Error(err))
		return fmt.Errorf("could not get all stream messages: %s", err.Error())
	}

	for _, message := range allMessages {
		store.Set(message.Id, message.Data)
	}
	logger.Info("all stream messages loaded")
	return nil
}
