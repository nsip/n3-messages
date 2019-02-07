package n3grpc

import (
	"testing"

	"github.com/nsip/n3-messages/messages"
	"github.com/nsip/n3-messages/messages/pb"
)

func TestPublish(t *testing.T) {
	defer func() { PH(recover(), "./log.txt") }()
	clt := NewClient("localhost", 5778)
	tuple := Must(messages.NewTuple("subject1", "predicate1longlonglonglonglonglonglonglong", "obj1")).(*pb.SPOTuple)

	for i := 0; i < 50; i++ {
		tuple.Version = int64(i)
		clt.Publish(tuple, "namespace", "contextName")
	}
	clt.Close()
}

func TestQuery(t *testing.T) {
	defer func() { PH(recover(), "./log.txt") }()
	clt := NewClient("localhost", 5778)
	tuple := Must(messages.NewTuple("subject1", "predicate1longlonglonglonglonglonglonglong", "obj1")).(*pb.SPOTuple)
	for _, t := range clt.Query(tuple, "namespace", "contextName") {
		fPln(*t)
	}
}

func TestDelete(t *testing.T) {
	defer func() { PH(recover(), "./log.txt") }()
	clt := NewClient("localhost", 5778)
	tuple := Must(messages.NewTuple("subject1", "predicate1longlonglonglonglonglonglonglong", "obj1")).(*pb.SPOTuple)
	n := clt.Delete(tuple, "namespace", "contextName")
	fPln(n)
}
