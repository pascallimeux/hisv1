package main

import (
	"os"
	"testing"
)

func TestHelloNominal(t *testing.T) {
	VERSION := "version1"
	version := VERSION
	if version != VERSION {
		t.Log("getversion", version, "was not", VERSION, "as expected")
		t.FailNow()
	}
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
