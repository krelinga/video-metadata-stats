package main

import "fmt"

func dirsWithNfo(dirs []movieDir) ([]movieDir, error) {
	dirsWithNfo := make([]movieDir, 0)
	dirsWithoutNfo := make([]movieDir, 0)
	for _, dir := range dirs {
		hasNfo, err := dir.HasNfoFile()
		if err != nil {
			return nil, err
		}
		if hasNfo {
			dirsWithNfo = append(dirsWithNfo, dir)
		} else {
			dirsWithoutNfo = append(dirsWithoutNfo, dir)
		}
	}
	if len(dirsWithoutNfo) > 0 {
		fmt.Println("Directories without NFO files:")
		for _, dir := range dirsWithoutNfo {
			fmt.Printf(" * %s\n", dir.Name())
		}
		fmt.Println()
	}
	return dirsWithNfo, nil
}
