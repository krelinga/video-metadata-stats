package main

import (
	"fmt"
	"log"
)

func main() {
	dirs, err := listMovieDirs()
	if err != nil {
		log.Fatal(err)
	}
	dirsWithNfo := make([]movieDir, 0)
	for _, dir := range dirs {
		hasNfo, err := dir.HasNfoFile()
		if err != nil {
			log.Fatal(err)
		}
		if hasNfo {
			dirsWithNfo = append(dirsWithNfo, dir)
		}
	}
	stats, err := computeDirStats(dirsWithNfo)
	if err != nil {
		log.Fatal(err)
	}
	for _, stat := range stats {
		fmt.Printf("Directory: %s, Has Images: %v\n", stat.Dir.Name(), stat.HasImages)
	}

	tagsByImages(stats)
}
