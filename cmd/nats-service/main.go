package main

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func run(ctx context.Context) error {
	connect, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return fmt.Errorf("connect: %v", err.Error())
	}

	stream, err := connect.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return fmt.Errorf("get JetStream: %v", err.Error())
	}

	sub, err := stream.SubscribeSync("ORDERS.*")
	if err != nil {
		return fmt.Errorf("subscribe: %v", err.Error())
	}
	defer sub.Unsubscribe()
	defer sub.Drain()

streamLoop:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("done")
			break streamLoop
		default:
			msg, err := sub.NextMsg(time.Second)
			if err == nats.ErrTimeout {
				continue
			}
			if err != nil {
				// migrate to zap logger
				_, _ = fmt.Fprintf(os.Stderr, "next msg: %s", err.Error())
			}
			_ = msg.Ack()
			fmt.Println("new message:", string(msg.Data))
		}
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
