package main

import (
	"github.com/beevik/etree"
)

// readNfo reads the XML file at the given path and returns the etree.Document representation of the DOM.
func readNfo(path string) (*etree.Document, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(path); err != nil {
		return nil, err
	}
	return doc, nil
}
