package port

import (
	"context"

	"github.com/aditya3232/my-grpc-proto/protogen/go/hello"
	"google.golang.org/grpc"
)

// hello service client proto
type HelloClientPort interface {
	SayHello(ctx context.Context, in *hello.HelloRequest, opts ...grpc.CallOption) (*hello.HelloResponse, error)
	SayManyHellos(ctx context.Context, in *hello.HelloRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[hello.HelloResponse], error)
}
