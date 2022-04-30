package models

import (
	"testing"
)

func Test_StorageToJSON(t *testing.T) {
	storage := new(Storage)
	b, err := storage.ToJSON()
	if err != nil {
		t.Fatal(err)
	}
	t.Log((string(b)))
}
