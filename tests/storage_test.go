package tests

import (
	"errors"
	"gRPCServer/internal/service"
	"testing"
)

func TestPut(t *testing.T) {
	storage := service.NewStorage()
	if err := storage.Put("test", "test"); err != nil {
		t.Error(err)
	}
	if err := storage.Put("test", "test"); err != nil {
		t.Error(err)
	}
	if err := storage.Put("1", "2"); err != nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {

	storage := service.NewStorage()
	_ = storage.Put("test", "test")
	_ = storage.Put("test", "test")
	_ = storage.Put("2", "1")

	if val, err := storage.Get("test"); err != nil {
		t.Error(err)
	} else if val != "test" {
		t.Error(val)
	}
	if val, err := storage.Get("2"); err != nil {
		t.Error(err)
	} else if val != "1" {
		t.Error(val)
	}
	if val, err := storage.Get("test"); err != nil {
		t.Error(err)
	} else if val != "test" {
		t.Error(val)
	}
	_ = storage.Delete("test")
	if _, err := storage.Get("test"); !errors.Is(err, service.ErrorNoSuchKey) {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	storage := service.NewStorage()
	_ = storage.Put("test", "test")
	_ = storage.Put("test", "test")
	_ = storage.Put("2", "1")

	if err := storage.Delete("test"); err != nil {
		t.Error(err)
	}
	if err := storage.Delete("test"); err != nil {
		t.Error(err)
	}

}
