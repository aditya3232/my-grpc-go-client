package bank

import (
	"context"
	"errors"
	"io"
	"log"
	"sync"

	dbank "github.com/aditya3232/my-grpc-go-client/internal/application/domain/bank"
	"github.com/aditya3232/my-grpc-go-client/internal/port"
	"github.com/aditya3232/my-grpc-proto/protogen/go/bank"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BankAdapter struct {
	bankClient port.BankClientPort
}

func NewBankAdapter(conn *grpc.ClientConn) (*BankAdapter, error) {
	client := bank.NewBankServiceClient(conn)

	return &BankAdapter{
		bankClient: client,
	}, nil
}

func (a *BankAdapter) GetCurrentBalance(ctx context.Context, acct string) (*bank.CurrentBalanceResponse, error) {
	bankRequest := &bank.CurrentBalanceRequest{
		AccoutNumber: acct,
	}

	bal, err := a.bankClient.GetCurrentBalance(ctx, bankRequest)
	if err != nil {
		st, _ := status.FromError(err)
		log.Fatalln("[FATAL] Error on GetCurrentBalance : ", st)
	}

	return bal, nil
}

func (a *BankAdapter) FetchExchangeRates(ctx context.Context, fromCur string, toCur string) {
	bankRequest := &bank.ExchangeRateRequest{
		FromCurrency: fromCur,
		ToCurrency:   toCur,
	}

	exchangeRateStream, err := a.bankClient.FetchExchangeRates(ctx, bankRequest)
	if err != nil {
		log.Fatalln("[FATAL] Error on fetchExchangeRates : ", err)
	}

	for {
		rate, err := exchangeRateStream.Recv()
		switch {
		case errors.Is(err, io.EOF):
			return
		case err != nil:
			st, _ := status.FromError(err)
			if st.Code() == codes.InvalidArgument {
				log.Fatalln("[FATAL] Error on fetchExchangeRates : ", st.Message())
			}
		}

		log.Printf("Rate at %v from %v to %v is %v \n", rate.Timestamp, rate.FromCurrency, rate.ToCurrency, rate.Rate)
	}
}

func (a *BankAdapter) SummarizeTransactions(ctx context.Context, acct string, tx []dbank.Transaction) {
	txStream, err := a.bankClient.SummarizeTransactions(ctx)
	if err != nil {
		log.Fatalln("[FATAL] Error on SummarizeTransactions : ", err)
	}

	for _, t := range tx {
		ttype := bank.TransactionType_TRANSACTION_TYPE_UNSPECIFIED

		switch t.TransactionType {
		case dbank.TransactionTypeIn:
			ttype = bank.TransactionType_TRANSACTION_TYPE_IN
		case dbank.TransactionTypeOut:
			ttype = bank.TransactionType_TRANSACTION_TYPE_OUT
		}

		bankRequest := &bank.Transaction{
			AccountNumber: acct,
			Type:          ttype,
			Amount:        t.Amount,
			Notes:         t.Notes,
		}

		txStream.Send(bankRequest)
	}

	summary, err := txStream.CloseAndRecv()
	if err != nil {
		st, _ := status.FromError(err)
		log.Fatalln("[FATAL] Error on SummarizeTransactions : ", st)
	}

	log.Println(summary)
}

func (a *BankAdapter) TransferMultiple(ctx context.Context, trf []dbank.TransferTransaction) {
	trfStream, err := a.bankClient.TransferMultiple(ctx)
	if err != nil {
		log.Fatalln("[FATAL] Error on TransferMultiple : ", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for _, tt := range trf {
			req := &bank.TransferRequest{
				FromAccountNumber: tt.FromAccountNumber,
				ToAccountNumber:   tt.ToAccountNumber,
				Currency:          tt.Currency,
				Amount:            tt.Amount,
			}
			trfStream.Send(req)
		}

		trfStream.CloseSend()
	}()

	go func() {
		defer wg.Done()
		for {
			res, err := trfStream.Recv()

			switch {
			case errors.Is(err, io.EOF):
				return
			case err != nil:
				handleTransferErrorGrpc(err)
			}

			log.Printf("Transfer status %v on %v \n", res.Status, res.Timestamp)
		}
	}()

	wg.Wait()
}

func handleTransferErrorGrpc(err error) {
	st := status.Convert(err)

	log.Printf("Error %v on TransferMultiple : %v \n", st.Code(), st.Message())

	for _, detail := range st.Details() {
		switch t := detail.(type) {
		case *errdetails.PreconditionFailure:
			for _, violation := range t.GetViolations() {
				log.Println("[VIOLATION] ", violation)
			}
		case *errdetails.ErrorInfo:
			log.Printf("Error on : %v, with reason %v \n", t.Domain, t.Reason)
			for i, v := range t.GetMetadata() {
				log.Printf("%v : %v \n", i, v)
			}
		}
	}
}
