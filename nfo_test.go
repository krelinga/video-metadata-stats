package main

import (
	"testing"
)

func TestReadNfo(t *testing.T) {
	movieDirs, err := listMovieDirs()
	if err != nil {
		t.Fatalf("Failed to list movie directories: %v", err)
	}
	doc, err := readNfo(movieDirs[0].NfoPath())
	if err != nil {
		t.Fatalf("Failed to read NFO file: %v", err)
	}
	if doc == nil {
		t.Fatal("Expected non-nil document")
	}
}