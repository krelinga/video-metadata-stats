package main

func tagStats() error {
	dirs, err := listMovieDirs()
	if err != nil {
		return err
	}

	dirsWithNfo, err := dirsWithNfo(dirs)
	if err != nil {
		return err
	}

	stats, err := computeDirStats(dirsWithNfo)
	if err != nil {
		return err
	}

	tags := tagsAbsentInDirsWithoutImages(stats)

	for _, tag := range tags {
		if err := examples(stats, tag); err != nil {
			return err
		}
	}

	return nil
}
