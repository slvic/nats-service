package app

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/slvic/nats-service/internal/configs"
	"github.com/slvic/nats-service/internal/http"
	worker "github.com/slvic/nats-service/internal/nats"
	"github.com/slvic/nats-service/internal/service"
	"github.com/slvic/nats-service/internal/service/deliveries"
	"github.com/slvic/nats-service/internal/store/memory"
	"github.com/slvic/nats-service/internal/store/persistent"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type App struct {
	nats *worker.NATS
	http *http.Server
}

func Initialize(ctx context.Context) (*App, error) {
	dbConfig, err := configs.NewDbConfig()
	if err != nil {
		return nil, err
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, fmt.Errorf("could not create new logger: %s", err.Error())
	}
	store := memory.New()
	dbConnectionString := fmt.Sprintf("user=%s dbname=%s password=%s port=%s sslmode=%s",
		dbConfig.DBUser,
		dbConfig.DBName,
		dbConfig.DBPassword,
		dbConfig.DBPort,
		dbConfig.SSLMode)
	database, err := persistent.New("postgres", dbConnectionString)
	if err != nil {
		return nil, fmt.Errorf("could not create new database: %s", err.Error())
	}

	err = loadOrdersFromDB(ctx, database, store, logger)
	if err != nil {
		return nil, fmt.Errorf("could not load orders from db: %s", err.Error())
	}

	connection, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, fmt.Errorf("connect: %s", err.Error())
	}

	elector, err := service.NewGraftElector(connection)
	if err != nil {
		return nil, err
	}

	storeService := deliveries.New(store, database, elector, logger)
	newWorker, err := worker.New(connection, logger)
	if err != nil {
		return nil, fmt.Errorf("could not create new worker: %s", err.Error())
	}
	ordersHandler := worker.NewOrdersHandler(storeService, logger)
	err = newWorker.AddWorker("ORDERS.*", ordersHandler)
	if err != nil {
		return nil, err
	}

	httpServer := http.New(storeService, logger)

	return &App{
		nats: newWorker,
		http: httpServer,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	errGroup, ctx := errgroup.WithContext(ctx)
	errGroup.Go(func() error {
		return a.nats.Start(ctx)
	})
	errGroup.Go(func() error {
		return a.http.Start(ctx)
	})
	errGroup.Go(func() error {
		<-ctx.Done()
		a.nats.Stop()
		return a.http.Stop(ctx)
	})

	return errGroup.Wait()
}

func loadOrdersFromDB(ctx context.Context, db *persistent.Database, store *memory.Store, logger *zap.Logger) error {
	allMessages, err := db.GetAll(ctx)
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
