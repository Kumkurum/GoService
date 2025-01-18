package tests

import (
	"gRPCServer/internal/service"
	"testing"
)

func TestWritePut(t *testing.T) {
	var fileLoger, err = service.NewFileTransactionLogger("testLogger")
	if err != nil {
		t.Error(err)
	}
	fileLoger.WritePut("testKey", "testValue")

}
