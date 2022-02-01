package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	worker "github.com/slvic/nats-service/internal/nats"
	"go.uber.org/zap"
)

func run(ctx context.Context) error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return fmt.Errorf("could not create new logger: %s", err.Error())
	}

	newWorker, err := worker.New(nats.DefaultURL, logger)
	if err != nil {
		return fmt.Errorf("could not create new worker: ")
	}

	err = newWorker.Run(ctx)
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
