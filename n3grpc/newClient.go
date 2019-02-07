// client.go

package n3grpc

import (
	"github.com/nsip/n3-messages/messages/pb"
	"google.golang.org/grpc"
)

//
// creates a grpc client that connects to the specified server hostname:port
// for practical use, is embedded in the publisher component.
//
func newAPIClient(host string, port int) (pb.APIClient, error) {
	serverAddr := fSf("%s:%d", host, port)
	conn := Must(grpc.Dial(serverAddr, grpc.WithInsecure())).(*grpc.ClientConn) // TODO: upgrade to non-insecure if required ie. jwt/tls etc.
	return pb.NewAPIClient(conn), nil
}
