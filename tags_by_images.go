package main

import (
	"fmt"
	"slices"
)

func tagsAbsentInDirsWithoutImages(stats []*dirStats) []tagPath {
	allTags := make(map[tagPath]struct{})
	for _, stat := range stats {
		for tag := range stat.Nfo.TagCounts {
			allTags[tag] = struct{}{}
		}
	}
	statsWithoutImages := make([]*dirStats, 0)
	for _, stat := range stats {
		if !stat.HasImages {
			statsWithoutImages = append(statsWithoutImages, stat)
		}
	}
	noImageTags := make(map[tagPath]struct{})
	for _, stat := range statsWithoutImages {
		for tag := range stat.Nfo.TagCounts {
			noImageTags[tag] = struct{}{}
		}
	}

	tagsAbsentInDirsWithoutImages := make([]tagPath, 0)
	for tag := range allTags {
		if _, ok := noImageTags[tag]; !ok {
			tagsAbsentInDirsWithoutImages = append(tagsAbsentInDirsWithoutImages, tag)
		}
	}

	if len(tagsAbsentInDirsWithoutImages) > 0 {
		slices.Sort(tagsAbsentInDirsWithoutImages)
		fmt.Println("Tags absent in dirs without images:")
		for _, tag := range tagsAbsentInDirsWithoutImages {
			fmt.Printf(" * %s\n", tag)
		}
	}
	return tagsAbsentInDirsWithoutImages
}
