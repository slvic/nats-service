package nats

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
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
	stream  nats.JetStream
	logger  *zap.Logger
}

func New(url string, logger *zap.Logger) (*NATS, error) {
	connect, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("connect: %s", err.Error())
	}

	stream, err := connect.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return nil, fmt.Errorf("get JetStream: %s", err.Error())
	}

	return &NATS{
		subs:    map[*nats.Subscription]Handler{},
		connect: connect,
		stream:  stream,
		logger:  logger,
	}, nil
}

func (n *NATS) AddWorker(subject string, handler Handler) error {
	sub, err := n.stream.SubscribeSync(subject)
	if err != nil {
		return fmt.Errorf("cound not subscribe to a subject: %s", err.Error())
	}
	n.subs[sub] = handler
	return nil
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
					n.logger.Info("stream loop stopping...")
					break streamLoop
				default:
					msg, err := sub.NextMsg(msgWaitTimeout)
					if err == nats.ErrTimeout {
						continue
					}
					if err != nil {
						n.logger.Error("next msg", zap.Error(err))
						time.Sleep(time.Second * 5)
					}

					if err := handler.Handle(msg.Data); err != nil {
						n.logger.Error("cant handle message", zap.Error(err))
						err = msg.Nak()
						if err != nil {
							n.logger.Error("nak", zap.Error(err))
							continue
						}
					}
					err = msg.Ack()
					if err != nil {
						n.logger.Error("ack", zap.Error(err))
						time.Sleep(time.Second * 5)
					}
				}
			}
		}()
	}

	wg.Wait()

	if !n.connect.IsClosed() {
		n.connect.Close()
	}
	n.logger.Info("stream loop stopped")
	return nil
}
