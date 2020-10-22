package main

import (
	"log"
	"reflect"

	"github.com/stone-co/the-amazing-ledger/clients/grpc/ledger"
)

func main() {
	log.Println("Ledger client example starting...")
	defer log.Println("Ledger client example finishing...")

	// Connect to the Ledger gRPC server
	host := "localhost"
	port := 50051
	conn, err := ledger.Connect(host, port)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	defineTransactionWithThreeEntries(conn)
	invalidTransactionsTests(conn)
}

func AssertEqual(expected, actual interface{}) {
	if reflect.DeepEqual(expected, actual) == false {
		log.Fatalf("Expected: %v Actual: %v\n", expected, actual)
	}
}
