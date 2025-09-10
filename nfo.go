package main

import (
	"fmt"
	"strconv"

	"github.com/beevik/etree"
)

type nfo struct {
	TagCounts map[tagPath]int
	Doc       *etree.Document
}

func (n *nfo) TmdbId() (int32, error) {
	elem := n.Doc.FindElement("//movie/tmdbid")
	if elem == nil {
		return 0, fmt.Errorf("tmdbid not found")
	}
	id, err := strconv.ParseInt(elem.Text(), 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(id), nil
}

func makeTagCounts(doc *etree.Document) map[tagPath]int {
	counts := make(map[tagPath]int)
	recursiveTagCounts(&doc.Element, newTagPath(), counts)
	return counts
}

func recursiveTagCounts(elem *etree.Element, current tagPath, counts map[tagPath]int) {
	for _, child := range elem.ChildElements() {
		childPath := current.Append(child.Tag)
		counts[childPath]++
		recursiveTagCounts(child, childPath, counts)
	}
}

// readNfo reads the XML file at the given path and returns the etree.Document representation of the DOM.
func readNfo(path string) (*nfo, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(path); err != nil {
		return nil, err
	}

	info := &nfo{
		Doc:       doc,
		TagCounts: makeTagCounts(doc),
	}
	return info, nil
}
