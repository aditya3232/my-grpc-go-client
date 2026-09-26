package port

import (
	"context"

	resl "github.com/aditya3232/my-grpc-proto/protogen/go/resiliency"
	"google.golang.org/grpc"
)

type ResiliencyClientPort interface {
	UnaryResiliency(ctx context.Context, in *resl.ResiliencyRequest, opts ...grpc.CallOption) (*resl.ResiliencyResponse, error)
	ServerStreamingResiliency(ctx context.Context, in *resl.ResiliencyRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[resl.ResiliencyResponse], error)
	ClientStreamingResiliency(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[resl.ResiliencyRequest, resl.ResiliencyResponse], error)
	BiDirectionalResiliency(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[resl.ResiliencyRequest, resl.ResiliencyResponse], error)
}

type ResiliencyWithMetadataClientPort interface {
	UnaryResiliencyWithMetadata(ctx context.Context, in *resl.ResiliencyRequest, opts ...grpc.CallOption) (*resl.ResiliencyResponse, error)
	ServerStreamingResiliencyWithMetadata(ctx context.Context, in *resl.ResiliencyRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[resl.ResiliencyResponse], error)
	ClientStreamingResiliencyWithMetadata(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[resl.ResiliencyRequest, resl.ResiliencyResponse], error)
	BiDirectionalResiliencyWithMetadata(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[resl.ResiliencyRequest, resl.ResiliencyResponse], error)
}
