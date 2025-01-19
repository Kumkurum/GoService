package service

import (
	"errors"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type HttpHandler struct {
	storage *Storage
	logger  TransactionLogger
}

func NewHttpHandler(storage *Storage, logger TransactionLogger) *HttpHandler {
	return &HttpHandler{storage: storage, logger: logger}
}

func (h *HttpHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	err := h.storage.Delete(key)
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK) // Все хорошо! Вернуть StatusOK
	//h.logger.WriteDelete(key)
}

func (h *HttpHandler) Put(w http.ResponseWriter, r *http.Request) {
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
	err = h.storage.Put(key, string(value))
	// Сохранить значение как строку
	if err != nil {
		// Если возникла ошибка, сообщить о ней
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated) // Все хорошо! Вернуть StatusCreated
	//h.logger.WritePut(key, string(value))
}

func (h *HttpHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Извлечь ключ из запроса
	key := vars["key"]
	value, err := h.storage.Get(key)
	// Получить значение для данного ключа
	if errors.Is(err, ErrorNoSuchKey) {
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
