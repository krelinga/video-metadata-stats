package main

import (
	"os"
	"path/filepath"
)

type movieDir string

func (d movieDir) Path() string {
	return string(d)
}

func (d movieDir) Name() string {
	return filepath.Base(d.Path())
}

func (d movieDir) NfoPath() string {
	return filepath.Join(d.Path(), d.Name()+".nfo")
}

func (d movieDir) HasImages() (bool, error) {
	files, err := os.ReadDir(d.Path())
	if err != nil {
		return false, err
	}
	isImageFile := func(name string) bool {
		ext := filepath.Ext(name)
		switch ext {
		case ".jpg", ".jpeg", ".png", ".gif":
			return true
		default:
			return false
		}
	}
	for _, f := range files {
		if isImageFile(f.Name()) {
			return true, nil
		}
	}
	return false, nil
}

func (d movieDir) HasNfoFile() (bool, error) {
	stat, err := os.Stat(d.NfoPath())
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return !stat.IsDir(), nil
}
