package main

import "path/filepath"

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