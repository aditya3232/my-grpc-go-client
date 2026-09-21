package hello

import (
	"context"
	"log"

	"github.com/aditya3232/my-grpc-go-client/internal/port"
	"github.com/aditya3232/my-grpc-proto/protogen/go/hello"
	"google.golang.org/grpc"
)

type HelloAdapter struct {
	helloClient port.HelloClientPort
}

func NewHelloAdapter(conn *grpc.ClientConn) (*HelloAdapter, error) {
	client := hello.NewHelloServiceClient(conn)

	return &HelloAdapter{
		helloClient: client,
	}, nil
}

func (a *HelloAdapter) SayHello(ctx context.Context, name string) (*hello.HelloResponse, error) {
	helloRequest := &hello.HelloRequest{
		Name: name,
	}

	greet, err := a.helloClient.SayHello(ctx, helloRequest)
	if err != nil {
		log.Fatalln("Error on SayHello : ", err)
	}

	return greet, nil
}
