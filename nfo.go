package main

import (
	"github.com/beevik/etree"
)

type nfo struct {
	TagCounts map[string]int
}

func makeTagCounts(doc *etree.Document) map[string]int {
	counts := make(map[string]int)
	movie := doc.FindElement("movie")
	if movie == nil {
		return counts
	}
	for _, elem := range movie.ChildElements() {
		counts[elem.Tag]++
	}
	return counts
}

// readNfo reads the XML file at the given path and returns the etree.Document representation of the DOM.
func readNfo(path string) (*nfo, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(path); err != nil {
		return nil, err
	}

	info := &nfo{
		TagCounts: makeTagCounts(doc),
	}
	return info, nil
}
