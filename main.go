package main

import (
	"log"
)

func main() {
	dirs, err := listMovieDirs()
	if err != nil {
		log.Fatal(err)
	}

	dirsWithNfo, err := dirsWithNfo(dirs)
	if err != nil {
		log.Fatal(err)
	}

	stats, err := computeDirStats(dirsWithNfo)
	if err != nil {
		log.Fatal(err)
	}

	tags := tagsAbsentInDirsWithoutImages(stats)

	for _, tag := range tags {
		if err := examples(stats, tag); err != nil {
			log.Fatal(err)
		}
	}
}
