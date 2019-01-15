// server.go

package n3grpc

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/nsip/n3-messages/messages/pb"
	"google.golang.org/grpc"
)

//
// implementation of the grpc server used by the node
// to receive 'raw' messages from outside the n3
// environment
// API is defined in messages/pb/n3msg.proto
//

type MessageHandler func(*pb.N3Message)
type QueryHandler func(*pb.N3Message) []*pb.SPOTuple

type APIServer struct {
	msgHandler   MessageHandler
	queryHandler QueryHandler
	grpcServer   *grpc.Server
}

//
// creates the API server which implements the Publish method
// and wraps in a generic grpc server which can be launched
// with normal tcp dial params elsewhere - e.g. in the
// context of the n3 node, from within the client
//
// the server will simply consume (nullify) messages unless a
// handler is provided to supply some business logic.
//
func NewAPIServer() *APIServer {

	apiServer := &APIServer{
		grpcServer: grpc.NewServer(),
		msgHandler: func(msg *pb.N3Message) { msg = nil },
		queryHandler: func(msg *pb.N3Message) []*pb.SPOTuple {
			msg = nil
			return nil
		},
	}
	// bind this api server logic into the generic grpc server
	pb.RegisterAPIServer(apiServer.grpcServer, apiServer)

	return apiServer
}

//
// creates a running api server on the given port
//
func (s *APIServer) Start(port int) error {

	// get underlying tcp connection
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	PE1(err, "cannot start grpc server: ")

	// launch the server
	go func() {
		log.Println("grpc server running on port", port)
		s.grpcServer.Serve(lis)
	}()

	return nil
}

//
// the server simply receives n3messages in a stream from clients.
// the handler allows business logic to be added by the creator of
// server to process messages
//
func (s *APIServer) SetMessageHandler(mh MessageHandler, qh QueryHandler) {
	s.msgHandler = mh
	s.queryHandler = qh
}

// Publish : Recv, SendAndClose
func (s *APIServer) Publish(stream pb.API_PublishServer) error {
	var msgCount int64
	for {
		n3msg, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.TxSummary{
				MsgCount: msgCount,
			})
		}
		if err != nil {
			return err
		}
		// do something with the msg
		s.msgHandler(n3msg)
		msgCount++
	}
}

// Query : Send
func (s *APIServer) Query(msg *pb.N3Message, stream pb.API_QueryServer) error {
	tuples := s.queryHandler(msg)
	// tuples = []*pb.SPOTuple{}
	for _, t := range tuples {
		if err := stream.Send(t); err != nil {
			return err
		}
	}
	return nil
}
