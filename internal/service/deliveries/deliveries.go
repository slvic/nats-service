package deliveries

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"os"
	"time"
)

type Stream struct {
	subscriptions map[string]*nats.Subscription
	jetstream     nats.JetStream
}

func NewStream(url string) (*Stream, error) {
	connect, err := nats.Connect(url)
	if err != nil {
		return &Stream{}, fmt.Errorf("connect: %v", err.Error())
	}

	stream, err := connect.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return &Stream{}, fmt.Errorf("get JetStream: %v", err.Error())
	}
	return &Stream{jetstream: stream}, nil
}

func (s *Stream) Subscribe(subject string) error {
	if _, ok := s.subscriptions[subject]; ok {
		return fmt.Errorf("subscription already exists")
	}

	sub, err := s.jetstream.SubscribeSync(subject)
	if err != nil {
		return fmt.Errorf("cound not subscribe to a subject: %v", err)
	}
	s.subscriptions[subject] = sub
	return nil
}

func (s *Stream) Publish(subject string, message string) error {
	_, err := s.jetstream.PublishAsync(subject, []byte(message))
	if err != nil {
		return fmt.Errorf("could not publish a message: %v", err)
	}
	return nil
}

func (s *Stream) Run(ctx context.Context) error {
	var err error
	defer func(subscriptions map[string]*nats.Subscription) {
		for _, sub := range subscriptions {
			if unSubErr := sub.Unsubscribe(); unSubErr != nil {
				err = fmt.Errorf("could not unsubscribe properly: %v", unSubErr)
			}
		}
	}(s.subscriptions)
	defer func(subscriptions map[string]*nats.Subscription) {
		for _, sub := range subscriptions {
			if drainErr := sub.Drain(); drainErr != nil {
				err = fmt.Errorf("could not unsubscribe properly: %v", drainErr)
			}
		}
	}(s.subscriptions)

streamLoop:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("done")
			break streamLoop
		default:
			for _, sub := range s.subscriptions {
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
	}
	if err != nil {
		return err
	}
	return nil
}
