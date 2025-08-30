package main

import (
	"testing"
)

func TestReadNfo(t *testing.T) {
	movieDirs, err := listMovieDirs()
	if err != nil {
		t.Fatalf("Failed to list movie directories: %v", err)
	}
	movieDir := movieDirs[0] // Just test the first one for simplicity
	t.Logf("Testing NFO read for movie directory: %s", movieDir.Path())
	doc, err := readNfo(movieDir.NfoPath())
	if err != nil {
		t.Fatalf("Failed to read NFO file: %v", err)
	}
	if doc == nil {
		t.Fatal("Expected non-nil document")
	}
	if len(doc.TagCounts) == 0 {
		t.Fatal("Expected non-empty tag counts")
	}
	for tag, count := range doc.TagCounts {
		t.Logf("Found %d occurrences of tag: %s", count, tag)
	}
}