package n3grpc

import (
	"fmt"
	"testing"

	//"github.com/nsip/n3-transport/messages"
	"../messages"
)

func TestNewPublisher(t *testing.T) {
	defer func() { PH(recover(), "./log.txt", true) }()

	pub, e := NewPublisher("localhost", 5778)
	PE(e)

	tuple, e := messages.NewTuple("subject1", "predicate1longlonglonglonglonglonglonglong", "obj1")
	PE(e)

	for i := 0; i < 50; i++ {
		tuple.Version = int64(i)
		pub.Publish(tuple, "namespace", "contextName")
	}
	pub.Close() 

	tuples := pub.Query(tuple, "namespace", "contextName")
	for _, t := range tuples {
		fmt.Println(*t)
	}
}
