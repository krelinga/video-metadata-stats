package main

import (
	"os"
	"slices"
	"testing"
)

func TestListMovieDirs(t *testing.T) {
	moviesMissingNfo := []string{
		"All Is Lost (2013)",
		"Arctic (2018)",
		"Predator (1987)",
	}
	dirs, err := listMovieDirs()
	if err != nil {
		t.Fatalf("Failed to list movie directories: %v", err)
	}
	for _, d := range dirs {
		// t.Logf("Found movie directory: %s, name: %s, nfo: %s", d.Path(), d.Name(), d.NfoPath())

		info, err := os.Stat(d.Path())
		if err != nil {
			t.Errorf("Path does not exist: %s, error: %v", d.Path(), err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("Path is not a directory: %s", d.Path())
		}

		if !slices.Contains(moviesMissingNfo, d.Name()) {
			nfoInfo, err := os.Stat(d.NfoPath())
			if err != nil {
				t.Errorf("NFO file does not exist: %s, error: %v", d.NfoPath(), err)
				continue
			}
			if !nfoInfo.Mode().IsRegular() {
				t.Errorf("NFO path is not a regular file: %s", d.NfoPath())
			}
		}

	}
}
