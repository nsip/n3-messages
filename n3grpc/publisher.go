// publisher.go

package n3grpc

import (
	"context"
	"io"
	"log"

	"../messages"    //"github.com/nsip/n3-transport/messages"
	"../messages/pb" //"github.com/nsip/n3-transport/messages/pb"
)

//
// publisher is the business-level wrapper around the
// grpc client, that exposes domain-friendly methods
// for sending new tuples/messages to an n3 node
//
type Publisher struct { /* Publisher should be changed to another name */
	grpcClient pb.APIClient
	stream     pb.API_PublishClient /* Send, CloseAndRecv */
	streamQ    pb.API_QueryClient
}

func NewPublisher(host string, port int) (*Publisher, error) {

	client, err := newAPIClient(host, port)
	PE1(err, "unable to create grpc client")

	stream, err := client.Publish(context.Background())
	PE1(err, "unable to create stream connection to server")

	pub := &Publisher{grpcClient: client, stream: stream}
	return pub, nil
}

func (pub *Publisher) Query(tuple *pb.SPOTuple, namespace, contextName string) (rt []*pb.SPOTuple) {

	// encode the tuple
	payload, err := messages.EncodeTuple(tuple)
	PE(err)

	// set the message envelope parameters & payload
	n3msg := &pb.N3Message{
		Payload:   payload,
		NameSpace: namespace,
		CtxName:   contextName,
	}

	stream, err := pub.grpcClient.Query(context.Background(), n3msg)
	PE(err)

	//rt := []*pb.SPOTuple{}

	for {
		t, err := stream.Recv()
		if err == io.EOF {
			break
		}
		PE(err)
		//fmt.Printf("client got: %v", *t)
		rt = append(rt, t)
	}

	return
}

//
// constructs a trasport-level message from the tup,e and the delivery params (namespace, contextName)
// and sends to the grpc server on the n3 node
//
// tuple - an SPO Tuple (no version required, will be assinged by node)
// namespace - owner of the delivery context - a base58 pub key id
// contextName - the delivery context for this data
//
func (pub *Publisher) Publish(tuple *pb.SPOTuple, namespace, contextName string) error {

	// encode the tuple
	payload, err := messages.EncodeTuple(tuple)
	if err != nil {
		return err
	}

	// set the message envelope parameters & payload
	n3msg := &pb.N3Message{
		Payload:   payload,
		NameSpace: namespace,
		CtxName:   contextName,
	}

	// send the message
	err = pub.stream.Send(n3msg)
	if err != nil {
		return err
	}

	return nil
}

//
// closes the network connections to the grpc server, and retrieves
// a tx summary report of how many messages have been sent from this
// publisher.
//
func (pub *Publisher) Close() {

	// get the tx report
	reply, err := pub.stream.CloseAndRecv()
	if err != nil {
		log.Printf("%v.CloseAndRecv() got error %v, want %v", pub.stream, err, nil)
	}
	log.Printf("Tx Summary: %v", reply)

}
