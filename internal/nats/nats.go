package nats

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

type Store interface {
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
	Delete(key string) error
}

type Consumer struct {
	sub  *nats.Subscription
	conn *nats.Conn
}

func NewConsumer(url string) (*Consumer, error) {
	connect, err := nats.Connect(url)
	if err != nil {
		return &Consumer{}, fmt.Errorf("connect: %v", err.Error())
	}

	stream, err := connect.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return &Consumer{}, fmt.Errorf("get JetStream: %v", err.Error())
	}

	sub, err := stream.SubscribeSync("ORDERS.*")
	if err != nil {
		return &Consumer{}, fmt.Errorf("cound not subscribe to a subject: %v", err)
	}

	return &Consumer{
		sub:  sub,
		conn: connect,
	}, nil
}

func (c *Consumer) Run(ctx context.Context) error {
	var err error
	defer func(sub *nats.Subscription) {
		if unSubErr := sub.Unsubscribe(); unSubErr != nil {
			err = fmt.Errorf("could not unsubscribe properly: %v", unSubErr)
		}
	}(c.sub)
	defer func(sub *nats.Subscription) {
		if drainErr := sub.Drain(); drainErr != nil {
			err = fmt.Errorf("could not unsubscribe properly: %v", drainErr)
		}
	}(c.sub)

streamLoop:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("done")
			break streamLoop
		default:
			msg, err := c.sub.NextMsg(time.Second)
			if err == nats.ErrTimeout {
				log.Println("could not read next message: timeout")
				continue
			}
			if err != nil {
				// migrate to zap logger
				_, _ = fmt.Fprintf(os.Stderr, "next msg: %s", err.Error())
			}
			_ = msg.Ack()

			_, err = UnmarshalAndValidate(msg.Data)
			if err != nil {
				log.Printf("could not unmarshal or validate message: %v", err)
				continue
			}
			//
		}
	}
	if err != nil {
		return err
	}
	return nil
}
