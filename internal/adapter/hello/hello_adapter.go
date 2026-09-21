package hello

import (
	"context"
	"errors"
	"io"
	"log"
	"time"

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

func (a *HelloAdapter) SayManyHellos(ctx context.Context, name string) {
	helloRequest := &hello.HelloRequest{
		Name: name,
	}

	greetStream, err := a.helloClient.SayManyHellos(ctx, helloRequest)
	if err != nil {
		log.Fatalln("Error on SayManyHellos : ", err)
	}

	for {
		greet, err := greetStream.Recv()
		switch {
		case errors.Is(err, io.EOF):
			return
		case err != nil:
			log.Fatalln("Error on SayManyHellos : ", err)
		}

		log.Println(greet.Greet)
	}
}

func (a *HelloAdapter) SayHelloToEveryone(ctx context.Context, names []string) {
	greetStream, err := a.helloClient.SayHelloToEveryone(ctx)
	if err != nil {
		log.Fatalln("Error on SayHelloToEveryone : ", err)
	}

	for _, name := range names {
		req := &hello.HelloRequest{
			Name: name,
		}

		greetStream.Send(req)
		time.Sleep(500 * time.Millisecond)
	}

	res, err := greetStream.CloseAndRecv()
	if err != nil {
		log.Fatalln("Error on SayHelloToEveryone : ", err)
	}

	log.Println(res.Greet)
}
