package app

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/slvic/nats-service/internal/configs"
	"github.com/slvic/nats-service/internal/http"
	worker "github.com/slvic/nats-service/internal/nats"
	"github.com/slvic/nats-service/internal/service/deliveries"
	"github.com/slvic/nats-service/internal/store/memory"
	"github.com/slvic/nats-service/internal/store/persistent"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type App struct {
	Nats *worker.NATS
	Http *http.Server
}

func Initialize(logger *zap.Logger, store *memory.Store, dbConfig *configs.DatabaseConfig) (*App, error) {
	dbConnectionString := fmt.Sprintf("user=%s dbname=%s password=%s port=%s sslmode=%s",
		dbConfig.DBUser,
		dbConfig.DBName,
		dbConfig.DBPassword,
		dbConfig.DBPort,
		dbConfig.SSLMode)
	database, err := persistent.New("postgres", dbConnectionString)
	if err != nil {
		return &App{}, fmt.Errorf("could not create new database: %s", err.Error())
	}

	err = LoadOrdersFromDB(database, store, logger)
	if err != nil {
		return &App{}, fmt.Errorf("could not load orders from db: %s", err.Error())
	}

	storeService := deliveries.New(store, database, logger)

	newWorker, err := worker.New(nats.DefaultURL, logger)
	if err != nil {
		return &App{}, fmt.Errorf("could not create new worker: %s", err.Error())
	}
	ordersHandler := worker.NewOrdersHandler(storeService, logger)
	err = newWorker.AddWorker("ORDERS.*", ordersHandler)
	if err != nil {
		return &App{}, err
	}

	httpServer := http.New(storeService, logger)

	return &App{
		Nats: newWorker,
		Http: httpServer,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	errGroup := new(errgroup.Group)
	errGroup.Go(func() error {
		err := a.Nats.Run(ctx)
		if err != nil {
			return err
		}
		return nil
	})
	err := a.Http.Start(ctx)
	if err != nil {
		return err
	}
	if err := errGroup.Wait(); err != nil {
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
