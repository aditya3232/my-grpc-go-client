package main

import (
	"context"
	"log"
	"time"

	// "github.com/aditya3232/my-grpc-go-client/internal/adapter/bank"
	// dbank "github.com/aditya3232/my-grpc-go-client/internal/application/domain/bank"
	// "github.com/aditya3232/my-grpc-go-client/internal/adapter/hello"

	"github.com/aditya3232/my-grpc-go-client/internal/adapter/resiliency"
	dresl "github.com/aditya3232/my-grpc-go-client/internal/application/domain/resiliency"
	resl "github.com/aditya3232/my-grpc-proto/protogen/go/resiliency"

	"github.com/sony/gobreaker/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Ganti [T any] dengan tipe return dari fungsi yang akan
// dibungkus dengan cbreaker.Execute(...)
var cbreaker *gobreaker.CircuitBreaker[*resl.ResiliencyResponse]

func init() {
	mybreaker := gobreaker.Settings{
		Name: "my-circuit-breaker",
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)

			log.Printf("Circuit breaker failure is %v, requests is %v, means failure ratio : %v \n",
				counts.TotalFailures, counts.Requests, failureRatio)

			return counts.Requests >= 3 && failureRatio >= 0.6
		},
		Timeout:     4 * time.Second,
		MaxRequests: 3,
		OnStateChange: func(name string, from, to gobreaker.State) {
			log.Printf("Circuit breaker %v changed state, from %v to %v \n", name, from, to)
		},
	}

	cbreaker = gobreaker.NewCircuitBreaker[*resl.ResiliencyResponse](mybreaker)
}

func main() {
	log.SetFlags(0)
	log.SetOutput(logWriter{})

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient("localhost:9090", opts...)
	if err != nil {
		log.Fatalln("Can not connect to gRPC server : ", err)
	}

	defer conn.Close()

	// helloAdapter, err := hello.NewHelloAdapter(conn)
	// if err != nil {
	// 	log.Fatalln("Can not create HelloAdapter : ", err)
	// }

	// runSayHello(helloAdapter, "iashiddiqi")
	// runSayManyHellos(helloAdapter, "aditya3232")
	// runSayHelloToEveryone(helloAdapter, []string{"invoker", "spectre", "juggernut", "muerta", "sven"})
	// runSayHelloContinous(helloAdapter, []string{"anna", "bella", "carol", "diana", "emma"})

	// bankAdapter, err := bank.NewBankAdapter(conn)
	// if err != nil {
	// 	log.Fatalln("Can not create BankAdapter : ", err)
	// }

	// runGetCurrentBalance(bankAdapter, "7835697001")
	// runFetchExchangeRates(bankAdapter, "USD", "IDR")
	// runSummarizeTransactions(bankAdapter, "7835697001", 10)
	// runTransferMultiple(bankAdapter, "7835697001", "7835697003", 10)

	resiliencyAdapter, err := resiliency.NewResiliencyAdapter(conn)

	if err != nil {
		log.Fatalln("Can not create ResiliencyAdapter :", err)
	}

	// runUnaryResiliencyWithTimeout(resiliencyAdapter, 2, 8, []uint32{dresl.OK}, 5*time.Second)
	// runServerStreamingResiliencyWithTimeout(resiliencyAdapter, 0, 3, []uint32{dresl.OK}, 15*time.Second)
	// runClientStreamingResiliencyWithTimeout(resiliencyAdapter, 0, 3, []uint32{dresl.OK}, 10, 10*time.Second)
	// runBiDirectionalResiliencyWithTimeout(resiliencyAdapter, 0, 3, []uint32{dresl.OK}, 10, 10*time.Second)

	for i := 0; i < 300; i++ {
		runUnaryResiliencyWithCircuitBreaker(resiliencyAdapter, 0, 0, []uint32{dresl.UNKNOWN, dresl.OK})
		time.Sleep(time.Second)
	}
}

// func runSayHello(adapter *hello.HelloAdapter, name string) {
// 	greet, err := adapter.SayHello(context.Background(), name)
// 	if err != nil {
// 		log.Fatalln("Can not call SayHello : ", err)
// 	}

// 	log.Println(greet.Greet)
// }

// func runSayManyHellos(adapter *hello.HelloAdapter, name string) {
// 	adapter.SayManyHellos(context.Background(), name)
// }

// func runSayHelloToEveryone(adapter *hello.HelloAdapter, names []string) {
// 	adapter.SayHelloToEveryone(context.Background(), names)
// }

// func runSayHelloContinous(adapter *hello.HelloAdapter, names []string) {
// 	adapter.SayHelloContinous(context.Background(), names)
// }

// func runGetCurrentBalance(adapter *bank.BankAdapter, acct string) {
// 	bal, err := adapter.GetCurrentBalance(context.Background(), acct)
// 	if err != nil {
// 		log.Fatalln("Failed to call GetCurrentBalance ", err)
// 	}

// 	log.Println(bal)
// }

// func runFetchExchangeRates(adapter *bank.BankAdapter, fromCur string, toCur string) {
// 	adapter.FetchExchangeRates(context.Background(), fromCur, toCur)
// }

// func runSummarizeTransactions(adapter *bank.BankAdapter, acct string, numDummyTransactions int) {
// 	var tx []dbank.Transaction

// 	for i := 1; i <= numDummyTransactions; i++ {
// 		ttype := dbank.TransactionTypeIn

// 		if i%3 == 0 {
// 			ttype = dbank.TransactionTypeOut
// 		}

// 		t := dbank.Transaction{
// 			Amount:          float64(rand.Intn(500) + 10),
// 			TransactionType: ttype,
// 			Notes:           fmt.Sprintf("Dummy transaction %v", i),
// 		}

// 		tx = append(tx, t)
// 	}

// 	adapter.SummarizeTransactions(context.Background(), acct, tx)
// }

// func runTransferMultiple(adapter *bank.BankAdapter, fromAcct string, toAcct string, numDummyTransactions int) {
// 	var trf []dbank.TransferTransaction

// 	for i := 1; i <= numDummyTransactions; i++ {
// 		tr := dbank.TransferTransaction{
// 			FromAccountNumber: fromAcct,
// 			ToAccountNumber:   toAcct,
// 			Currency:          "USD",
// 			Amount:            float64(rand.Intn(200) + 5),
// 		}

// 		trf = append(trf, tr)
// 	}

// 	adapter.TransferMultiple(context.Background(), trf)
// }

// func runUnaryResiliencyWithTimeout(adapter *resiliency.ResiliencyAdapter, minDelaySecond int32,
// 	maxDelaySecond int32, statusCodes []uint32, timeout time.Duration) {
// 	ctx, cancel := context.WithTimeout(context.Background(), timeout)

// 	defer cancel()

// 	res, err := adapter.UnaryResiliency(ctx, minDelaySecond, maxDelaySecond, statusCodes)

// 	if err != nil {
// 		log.Fatalln("Failed to call UnaryResiliency :", err)
// 	}

// 	log.Println(res.DummyString)
// }

// func runServerStreamingResiliencyWithTimeout(adapter *resiliency.ResiliencyAdapter,
// 	minDelaySecond int32, maxDelaySecond int32, statusCodes []uint32, timeout time.Duration) {
// 	ctx, cancel := context.WithTimeout(context.Background(), timeout)
// 	defer cancel()

// 	adapter.ServerStreamingResiliency(ctx, minDelaySecond, maxDelaySecond, statusCodes)
// }

// func runClientStreamingResiliencyWithTimeout(adapter *resiliency.ResiliencyAdapter,
// 	minDelaySecond int32, maxDelaySecond int32, statusCodes []uint32,
// 	count int, timeout time.Duration) {
// 	ctx, cancel := context.WithTimeout(context.Background(), timeout)
// 	defer cancel()

// 	adapter.ClientStreamingResiliency(ctx, minDelaySecond, maxDelaySecond, statusCodes, count)
// }

// func runBiDirectionalResiliencyWithTimeout(adapter *resiliency.ResiliencyAdapter,
// 	minDelaySecond int32, maxDelaySecond int32, statusCodes []uint32,
// 	count int, timeout time.Duration) {
// 	ctx, cancel := context.WithTimeout(context.Background(), timeout)
// 	defer cancel()

// 	adapter.BiDirectionalResiliency(ctx, minDelaySecond, maxDelaySecond, statusCodes, count)
// }

func runUnaryResiliencyWithCircuitBreaker(adapter *resiliency.ResiliencyAdapter, minDelaySecond int32,
	maxDelaySecond int32, statusCodes []uint32) {

	res, err := cbreaker.Execute(func() (*resl.ResiliencyResponse, error) {
		return adapter.UnaryResiliency(context.Background(), minDelaySecond, maxDelaySecond, statusCodes)
	})

	if err != nil {
		log.Println("Error on UnaryResiliency : ", err)
		return
	}

	log.Println(res.DummyString)
}
