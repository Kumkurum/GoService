package sevice

import (
	"errors"
	"sync"
)

type Storage struct {
	sync.RWMutex
	store map[string]string
}

func NewStorage() Storage {
	return Storage{store: make(map[string]string)}
}

func (storage Storage) Put(key string, value string) error {
	storage.Lock()
	storage.store[key] = value
	storage.Unlock()
	return nil
}

var ErrorNoSuchKey = errors.New("no such key")

func (storage Storage) Get(key string) (string, error) {
	storage.RLock()
	value, ok := storage.store[key]
	storage.RUnlock()
	if !ok {
		return "", ErrorNoSuchKey
	}
	return value, nil
}
func (storage Storage) Delete(key string) error {
	delete(storage.store, key)
	return nil
}
