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

var params = service.PostgresDBParams{
	Port:     5432,
	DbName:   "app",
	Host:     "localhost",
	User:     "app_user",
	Password: "123",
}

func main() {
	var err error
	logger, err = service.NewPostgresTransactionLogger(params)
	if err != nil {
		fmt.Printf("failed to create event logger: %s", err)
		return
	}
	err = logger.Initialize(storage)
	defer logger.Close() // Закрытие файла в конце main
	if err != nil {
		fmt.Printf("failed to initialize transaction logger: %s", err)
		return
	}
	r := mux.NewRouter()
	handler := service.NewHttpHandler(storage, logger)
	r.HandleFunc("/v1/{key}", handler.Put).Methods("PUT")
	r.HandleFunc("/v1/{key}", handler.Get).Methods("GET")
	r.HandleFunc("/v1/delete/{key}", handler.Delete).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}
