package nats

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Handler interface {
	Handle(ctx context.Context, message []byte) error
}

type NATS struct {
	subs    map[*nats.Subscription]Handler
	connect *nats.Conn
	stream  nats.JetStream
	logger  *zap.Logger
}

func New(connect *nats.Conn, logger *zap.Logger) (*NATS, error) {
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

func (n *NATS) Start(ctx context.Context) error {
	wg := sync.WaitGroup{}

	for sub, handler := range n.subs {
		wg.Add(1)
		sub, handler := sub, handler
		go func() {
			for {
				msg, err := sub.NextMsgWithContext(ctx)
				if err == context.Canceled {
					wg.Done()
					return
				}
				if err != nil {
					n.logger.Error("next msg", zap.Error(err))
					time.Sleep(time.Second * 5)
				}

				if msg == nil {
					continue
				}

				if err := handler.Handle(ctx, msg.Data); err != nil {
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
		}()
	}

	wg.Wait()

	return nil
}

func (n *NATS) Stop() {
	if !n.connect.IsClosed() {
		n.connect.Close()
	}
	n.logger.Info("stream loop stopped")
}
