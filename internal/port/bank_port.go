package port

import (
	"context"

	"github.com/aditya3232/my-grpc-proto/protogen/go/bank"
	"google.golang.org/grpc"
)

type BankClientPort interface {
	GetCurrentBalance(ctx context.Context, in *bank.CurrentBalanceRequest, opts ...grpc.CallOption) (*bank.CurrentBalanceResponse, error)
	FetchExchangeRates(ctx context.Context, in *bank.ExchangeRateRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[bank.ExchangeRateResponse], error)
	SummarizeTransactions(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[bank.Transaction, bank.TransactionSummary], error)
	TransferMultiple(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[bank.TransferRequest, bank.TransferResponse], error)
	CreateAccount(ctx context.Context, in *bank.CreateAccountRequest, opts ...grpc.CallOption) (*bank.CreateAccountResponse, error)
}
