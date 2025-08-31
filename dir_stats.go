package main

type dirStats struct {
	Dir movieDir
	Nfo *nfo
	HasImages bool
}

func computeDirStats(dirs []movieDir) ([]*dirStats, error) {
	stats := make([]*dirStats, 0, len(dirs))
	for _, dir := range dirs {
		nfo, err := readNfo(dir.NfoPath())
		if err != nil {
			return nil, err
		}
		hasImages, err := dir.HasImages()
		if err != nil {
			return nil, err
		}
		stats = append(stats, &dirStats{
			Dir:     dir,
			Nfo:    nfo,
			HasImages: hasImages,
		})
	}
	return stats, nil
}