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
	moviesMissingImages := []string{
		"A Few Good Men (1992)",
		"Deep Water (2006)",
		"Ice People (2009)",
		"Over the Top (1987)",
		"Primer (2004)",
		"The Last Boy Scout (1991)",
	}
	dirs, err := listMovieDirs()
	if err != nil {
		t.Fatalf("Failed to list movie directories: %v", err)
	}
	for _, d := range dirs {
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

		if hasImages, err := d.HasImages(); err != nil {
			t.Errorf("Failed to check for images in directory: %s, error: %v", d.Path(), err)
			continue
		} else if !hasImages && !slices.Contains(moviesMissingImages, d.Name()) {
			t.Errorf("No images found in directory: %s", d.Path())
		} else if hasImages && slices.Contains(moviesMissingImages, d.Name()) {
			t.Errorf("Unexpected images found in directory: %s", d.Path())
		}
	}
}
