package nats

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	msgWaitTimeout = time.Second
)

type Handler interface {
	Handle(message []byte) error
}

type NATS struct {
	subs    map[*nats.Subscription]Handler
	connect *nats.Conn
}

func New(url string) (*NATS, error) {
	connect, err := nats.Connect(url)
	if err != nil {
		return &NATS{}, fmt.Errorf("connect: %s", err.Error())
	}

	stream, err := connect.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return &NATS{}, fmt.Errorf("get JetStream: %s", err.Error())
	}

	orders, err := stream.SubscribeSync("ORDERS.*")
	if err != nil {
		return &NATS{}, fmt.Errorf("cound not subscribe to a subject: %v", err)
	}

	return &NATS{
		subs: map[*nats.Subscription]Handler{
			orders: NewOrdersHandler(nil), // pass implementation
		},
		connect: connect,
	}, nil
}

func (n *NATS) Run(ctx context.Context) error {
	wg := sync.WaitGroup{}

	for sub, handler := range n.subs {
		wg.Add(1)
		sub, handler := sub, handler
		go func() {
		streamLoop:
			for {
				select {
				case <-ctx.Done():
					wg.Done()
					break streamLoop
				default:
					msg, err := sub.NextMsg(msgWaitTimeout)
					if err == nats.ErrTimeout {
						log.Println("could not read next message: timeout")
						continue
					}
					if err != nil {
						_, _ = fmt.Fprintf(os.Stderr, "next msg: %s", err.Error())
						time.Sleep(time.Second * 5) // waiting for him to get better
					}

					if err := handler.Handle(msg.Data); err != nil {
						_ = msg.Nak()
					}
				}
			}
		}()
	}

	wg.Wait()

	if !n.connect.IsClosed() {
		n.connect.Close()
	}

	return nil
}
