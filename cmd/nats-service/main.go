package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/slvic/nats-service/configs"
	"github.com/slvic/nats-service/internal/app"
	"github.com/slvic/nats-service/internal/store/memory"
	"go.uber.org/zap"
)

func run(ctx context.Context) error {
	dbConfig, err := configs.NewDbConfig()
	logger, err := zap.NewDevelopment()
	if err != nil {
		return fmt.Errorf("could not create new logger: %s", err.Error())
	}
	store := memory.New()

	err = app.Initialize(ctx, logger, store, dbConfig)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL, syscall.SIGINT)
	defer cancel()
	if err := run(ctx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "app run: %v\n", err.Error())
	}
}
