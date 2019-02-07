package n3grpc

import (
	"testing"
	"time"

	"github.com/nsip/n3-messages/messages"
	"github.com/nsip/n3-messages/messages/pb"
)

func TestServer(t *testing.T) {
	defer func() { PH(recover(), "./log.txt") }()

	svr := NewAPIServer()

	pHandler := func(n3msg *pb.N3Message) {
		tuple := Must(messages.DecodeTuple(n3msg.Payload)).(*pb.SPOTuple)
		fPln(*tuple)
	}

	qHandler := func(n3msg *pb.N3Message) []*pb.SPOTuple {
		tuple := Must(messages.DecodeTuple(n3msg.Payload)).(*pb.SPOTuple)
		return []*pb.SPOTuple{tuple, tuple, tuple}
	}

	dHandler := func(n3msg *pb.N3Message) int {
		return 123456
	}

	svr.SetMessageHandler(pHandler, qHandler, dHandler)
	svr.Start(5778)
	time.Sleep(500 * time.Second)
}
