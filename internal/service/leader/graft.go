package leader

import (
	"fmt"

	"github.com/nats-io/graft"
	"github.com/nats-io/nats.go"
)

type GraftElector struct {
	node *graft.Node
}

func NewGraftElector(connect *nats.Conn) (*GraftElector, error) {
	ci := graft.ClusterInfo{Name: "order_service", Size: 3}
	rpc, err := graft.NewNatsRpc(&connect.Opts)
	if err != nil {
		return nil, fmt.Errorf("could not create instance of nats rpc driver: %s", err.Error())
	}
	errChan := make(chan error)
	stateChangeChan := make(chan graft.StateChange)
	handler := graft.NewChanHandler(stateChangeChan, errChan)
	node, err := graft.New(ci, handler, rpc, "/tmp/leader.log")
	if err != nil {
		return nil, fmt.Errorf("could not create leader node: %s", err.Error())
	}
	return &GraftElector{node: node}, nil
}

func (s *GraftElector) AmILeader() bool {
	return s.node.State() == graft.LEADER
}
