package main

import (
	"github.com/renanvicente/grpc_sample/hexarch/core"
	"github.com/renanvicente/grpc_sample/hexarch/frontend"
	"github.com/renanvicente/grpc_sample/hexarch/transact"
	"log"
	"os"
)

func main() {
	// Create our TransactionLogger. This is an adapter that will plug
	// into the core application's TransactionLogger port.
	tl, err := transact.NewTransactionLogger(os.Getenv("TLOG_TYPE"))

	if err != nil {
		log.Fatal(err)
	}

	// Create Core and tell it which TransactionLogger to use.
	// This is an example of a "driven agent"
	store := core.NewKeyValueStore(tl)
	log.Println("here")

	store.Restore()

	// Create the frontend.
	// This is an example of a "driving agent."
	fe, err := frontend.NewFrontEnd(os.Getenv("FRONTEND_TYPE"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("here")
	log.Fatal(fe.Start(store))
}


