package main

import (
	"os"
	"path/filepath"
)

func listMovieDirs() ([]movieDir, error) {
	const baseDir = "/nas/media/Movies"
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, err
	}
	dirs := []movieDir{}
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, movieDir(filepath.Join(baseDir, entry.Name())))
		}
	}
	return dirs, nil
}
