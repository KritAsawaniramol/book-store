package grpccon

import (
	"errors"
	"log"
	"net"

	"github.com/kritAsawaniramol/book-store/module/user/userPb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type (
	//Creator class
	GrpcClientFactoryHandler interface {
		User() userPb.UserGrpcServiceClient
	}

	grpcClientFactory struct {
		client *grpc.ClientConn
	}
)

// User implements GrpcClientFactoryHandler.
func (g *grpcClientFactory) User() userPb.UserGrpcServiceClient {
	return userPb.NewUserGrpcServiceClient(g.client)
}

func NewGrpcServer(host string) (*grpc.Server, net.Listener) {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("Error: Failed to listen: %v", err)
	}
	return grpcServer, listener
}

func NewGrpcClient(grpcUrl string) (GrpcClientFactoryHandler, error) {
	creds := insecure.NewCredentials()

	clientConn, err := grpc.NewClient(grpcUrl, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Printf("error: NewGrpcClient: %s\n", err.Error())
		return nil, errors.New("error: grpc client connection failed")
	}

	return &grpcClientFactory{
		client: clientConn,
	}, nil
}
