// publisher.go

package n3grpc

import (
	"context"
	"io"
	"log"

	"../messages"
	"../messages/pb"
)

//
// publisher is the business-level wrapper around the
// grpc client, that exposes domain-friendly methods
// for sending new tuples/messages to an n3 node
//
type Client struct {
	grpcClient pb.APIClient
	stream     pb.API_PublishClient /* Send, CloseAndRecv */
	streamQ    pb.API_QueryClient
}

func NewClient(host string, port int) *Client {

	clt := Must(newAPIClient(host, port)).(pb.APIClient)
	stream, err := clt.Publish(context.Background())
	PE1(err, "unable to create stream connection to server")
	return &Client{grpcClient: clt, stream: stream}
}

func (clt *Client) Delete(tuple *pb.SPOTuple, namespace, ctxName string) int {

	// set the message envelope parameters & payload
	n3msg := &pb.N3Message{
		Payload:   Must(messages.EncodeTuple(tuple)).([]byte),
		NameSpace: namespace,
		CtxName:   ctxName,
	}

	txsum, err := clt.grpcClient.Delete(context.Background(), n3msg)
	PE(err)

	return int(txsum.MsgCount)
}

func (clt *Client) Query(tuple *pb.SPOTuple, namespace, ctxName string) (rt []*pb.SPOTuple) {

	// set the message envelope parameters & payload
	n3msg := &pb.N3Message{
		Payload:   Must(messages.EncodeTuple(tuple)).([]byte),
		NameSpace: namespace,
		CtxName:   ctxName,
	}

	stream, err := clt.grpcClient.Query(context.Background(), n3msg)
	PE(err)

	for {
		t, err := stream.Recv()
		if err == io.EOF {
			break
		}
		PE(err)
		//fPf("client got: %v", *t)
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
func (clt *Client) Publish(tuple *pb.SPOTuple, namespace, ctxName string) error {

	// set the message envelope parameters & payload
	n3msg := &pb.N3Message{
		Payload:   Must(messages.EncodeTuple(tuple)).([]byte),
		NameSpace: namespace,
		CtxName:   ctxName,
	}

	// send the message
	PE(clt.stream.Send(n3msg))

	return nil
}

//
// closes the network connections to the grpc server, and retrieves
// a tx summary report of how many messages have been sent from this
// publisher.
//
func (clt *Client) Close() {

	// get the tx report
	reply, err := clt.stream.CloseAndRecv()
	PE1(err, fSf("%v.CloseAndRecv() got error %v, want %v", clt.stream, err, nil))
	log.Printf("Tx Summary: %v", reply)
}
