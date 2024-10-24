package main

import (
	"os"
	"testing"
)

func TestDeleteKV(t *testing.T) {
	testFile := "test.db"
	ds, err := InitDiskStore(testFile)
	if err != nil {
		t.Fatalf("failed to create test file %v", err)
	}
	defer ds.Close()
	defer func() {
		err := os.Remove(testFile)
		if err != nil {
			t.Fatalf("failed to delete test file %v", err)
		}
	}()

	k := "test key"
	v := "test val"
	err = ds.Set(k, v)
	if err != nil {
		t.Fatalf("expected value: 'test val, got '%s'", v)
	}
	err = ds.Delete(k)
	if err != nil {
		t.Fatalf("failed to delete key: %v", err)
	}

	val, _ := ds.Get(k)
	if val != "" {
		t.Fatalf("value should be an empty string")
	}

	ds.List(testFile)
}
