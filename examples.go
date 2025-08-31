package main

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
)

func examples(stats []*dirStats, tag tagPath) error {
	const wantCount = 5
	type exampleType struct {
		Element  *etree.Element
		MovieDir movieDir
	}
	examples := make([]exampleType, 0)
statloop:
	for _, stat := range stats {
		if len(examples) >= wantCount {
			break statloop
		}
		if _, ok := stat.Nfo.TagCounts[tag]; ok {
			// Tag is present in this dirStats
			path, err := tag.XmlPath()
			if err != nil {
				return err
			}
			for elem := range stat.Nfo.Doc.FindElementsPathSeq(path) {
				if len(examples) >= wantCount {
					break statloop
				}
				examples = append(examples, exampleType{
					Element:  elem,
					MovieDir: stat.Dir,
				})
			}
		}
	}

	if len(examples) == 0 {
		fmt.Printf("No examples found for tag %s\n\n", tag)
		return nil
	}

	fmt.Printf("Examples for tag %s:\n", tag)
	for _, example := range examples {
		example := exampleType{
			Element:  example.Element.Copy(),
			MovieDir: example.MovieDir,
		}
		indentSettings := etree.NewIndentSettings()
		indentSettings.UseTabs = false
		example.Element.IndentWithSettings(indentSettings)
		sb := strings.Builder{}
		example.Element.WriteTo(&sb, &etree.WriteSettings{})
		output := sb.String()
		output = strings.ReplaceAll(output, "\n", "\n      ")
		fmt.Printf(" * %s:\n      %s\n", example.MovieDir.Name(), output)
	}
	fmt.Println()

	return nil
}
