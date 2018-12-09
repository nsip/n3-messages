package n3grpc

import (
	"fmt"
	"testing"
	"time"

	"../messages"
	"../messages/pb" //"github.com/nsip/n3-transport/messages/pb"
)

func TestServer(t *testing.T) {
	defer func() { PH(recover(), "./log.txt", true) }()

	svr := NewAPIServer()

	pHandler := func(n3msg *pb.N3Message) {
		tuple, e := messages.DecodeTuple(n3msg.Payload)
		PE(e)
		fmt.Println(*tuple)
	}

	qHandler := func(n3msg *pb.N3Message) []*pb.SPOTuple {
		tuple, e := messages.DecodeTuple(n3msg.Payload)
		PE(e)
		return []*pb.SPOTuple{tuple, tuple, tuple}
	}

	svr.SetMessageHandler(pHandler, qHandler)
	svr.Start(5778)
	time.Sleep(500 * time.Second)
}
