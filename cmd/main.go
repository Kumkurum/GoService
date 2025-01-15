package main

import (
	"errors"
	"fmt"
	"gRPCServer/internal/service"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

func keyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Получить ключ из запроса
	key := vars["key"]
	value, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	// Тело запроса хранит значение
	if err != nil {
		// Если возникла ошибка, сообщить о ней
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}
	err = storage.Put(key, string(value))
	// Сохранить значение как строку
	if err != nil {
		// Если возникла ошибка, сообщить о ней
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated) // Все хорошо! Вернуть StatusCreated
}

func keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Извлечь ключ из запроса
	key := vars["key"]
	value, err := storage.Get(key)
	// Получить значение для данного ключа
	if errors.Is(err, service.ErrorNoSuchKey) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(value))
	// Записать значение в ответ
}

var (
	storage service.Storage = service.NewStorage()
	logger  service.TransactionLogger
)

func initializeTransactionLog() error {
	var err error
	logger, err = service.NewFileTransactionLogger("transaction.log")
	if err != nil {
		return fmt.Errorf("failed to create event logger: %w", err)
	}
	events, errorsLog := logger.ReadEvents()
	e, ok := service.Event{}, true

	for ok && err == nil {
		select {
		case err, ok = <-errorsLog:
		case e, ok = <-events:
			switch e.EventType {
			case service.EventDelete:
				err = storage.Delete(e.Key)
			case service.EventPut:
				err = storage.Put(e.Key, string(e.Value))
			}
		}
	}
	logger.Run()
	return err
}
func main() {
	/*port := ":8080"
	fmt.Println("Hello World")
	server := grpc.NewServer()
	l, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
	}
	fmt.Printf("Listening on %s", port)
	if err := server.Serve(l); err != nil {
		fmt.Printf("failed to serve: %v", err)
	}*/
	err := initializeTransactionLog()
	if err != nil {
		return
	}
	r := mux.NewRouter()
	// Зарегистрировать keyValuePutHandler как обработчик HTTP-запросов PUT,
	// в которых указан путь "/v1/{key}"
	r.HandleFunc("/v1/{key}", keyValuePutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", keyValueGetHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}
