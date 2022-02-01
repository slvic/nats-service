package deliveries

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/slvic/nats-service/internal/service"
	"github.com/slvic/nats-service/internal/store/memory"
)

type Consumer struct {
	jetstream     nats.JetStream
	subscriptions map[string]*nats.Subscription
	store         *memory.Store
}

func NewConsumer(url string, store *memory.Store) (*Consumer, error) {
	connect, err := nats.Connect(url)
	if err != nil {
		return &Consumer{}, fmt.Errorf("connect: %v", err.Error())
	}

	stream, err := connect.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return &Consumer{}, fmt.Errorf("get JetStream: %v", err.Error())
	}
	return &Consumer{
		jetstream: stream,
		store:     store,
	}, nil
}

func (c *Consumer) Subscribe(subject string) error {
	if _, ok := c.subscriptions[subject]; ok {
		return fmt.Errorf("subscription already exists")
	}

	sub, err := c.jetstream.SubscribeSync(subject)
	if err != nil {
		return fmt.Errorf("cound not subscribe to a subject: %v", err)
	}
	c.subscriptions[subject] = sub
	return nil
}

func (c *Consumer) Publish(subject string, message []byte) error {
	_, err := c.jetstream.PublishAsync(subject, message)
	if err != nil {
		return fmt.Errorf("could not publish a message: %v", err)
	}
	return nil
}

func (c *Consumer) GetMessagesBySubject(subject string) (string, error) {
	_, found := c.store.Get(subject)
	if !found {
		return "", fmt.Errorf("subject not found")
	}

	//
	return "", nil
}

func (c *Consumer) Run(ctx context.Context) error {
	var err error
	defer func(subscriptions map[string]*nats.Subscription) {
		for _, sub := range subscriptions {
			if unSubErr := sub.Unsubscribe(); unSubErr != nil {
				err = fmt.Errorf("could not unsubscribe properly: %v", unSubErr)
			}
		}
	}(c.subscriptions)
	defer func(subscriptions map[string]*nats.Subscription) {
		for _, sub := range subscriptions {
			if drainErr := sub.Drain(); drainErr != nil {
				err = fmt.Errorf("could not unsubscribe properly: %v", drainErr)
			}
		}
	}(c.subscriptions)

streamLoop:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("done")
			break streamLoop
		default:
			for _, sub := range c.subscriptions {
				msg, err := sub.NextMsg(time.Second)
				if err == nats.ErrTimeout {
					log.Println("could not read next message: timeout")
					continue
				}
				if err != nil {
					// migrate to zap logger
					_, _ = fmt.Fprintf(os.Stderr, "next msg: %s", err.Error())
				}
				_ = msg.Ack()

				order, err := service.UnmarshalAndValidate(msg.Data)
				if err != nil {
					log.Println("could not unmarshal or validate message: %v", err)
					continue
				}
				c.store.Set(sub.Subject, order)
			}
		}
	}
	if err != nil {
		return err
	}
	return nil
}
