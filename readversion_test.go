package main

import (
	"testing"
)

func TestReadVersion(t *testing.T) {
	version, err := readVersion()
	if err != nil {
		t.Fatalf("Failed to get version")
	}

	if version != Version {
		t.Errorf("Version not matched")
	}
}
