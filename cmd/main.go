package main

import (
	"fmt"
	"gRPCServer/internal/service"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var (
	storage = service.NewStorage()
	logger  service.TransactionLogger
)

func main() {
	var err error
	logger, err = service.NewFileTransactionLogger("transaction.log")
	if err != nil {
		fmt.Errorf("failed to create event logger: %w", err)
		return
	}
	err = logger.Initialize(storage)
	if err != nil {
		fmt.Errorf("failed to initialize transaction logger: %v\n", err)
		return
	}
	r := mux.NewRouter()

	handler := service.NewHttpHandler(storage, logger)
	r.HandleFunc("/v1/{key}", handler.Put).Methods("PUT")
	r.HandleFunc("/v1/{key}", handler.Get).Methods("GET")
	r.HandleFunc("/v1/{key}", handler.Delete).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}
