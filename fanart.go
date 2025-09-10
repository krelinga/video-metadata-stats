package main

import (
	"cmp"
	"context"
	"fmt"
	"net/url"
	"os"
	"slices"
	"strings"

	"github.com/beevik/etree"
	"github.com/krelinga/go-tmdb"
)

func countDirsWithFanart(stats []*dirStats) {
	dirsWithFanartCount := make(map[int]int)

	fanartTagPath := newTagPath("movie", "fanart", "thumb")
	for _, stat := range stats {
		count := stat.Nfo.TagCounts[fanartTagPath]
		dirsWithFanartCount[count]++
	}

	type kv struct {
		k, v int
	}
	var sorted []kv
	for k, v := range dirsWithFanartCount {
		sorted = append(sorted, kv{k, v})
	}
	slices.SortFunc(sorted, func(a, b kv) int {
		return cmp.Compare(a.k, b.k)
	})
	fmt.Println("Fanart counts:")
	for _, kv := range sorted {
		fmt.Printf(" * Files with %d fanart tags: %d\n", kv.k, kv.v)
	}
}

func extractDomain(rawUrl string) (string, error) {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}
	return parsedUrl.Hostname(), nil
}

func countDirsByFanartDomain(stats []*dirStats) error {
	domainToCount := make(map[string]int)
	fanartTagPath, err := etree.CompilePath("/movie/fanart/thumb")
	if err != nil {
		return err
	}
	for _, stat := range stats {
		for _, elem := range stat.Nfo.Doc.FindElementsPath(fanartTagPath) {
			domain, err := extractDomain(elem.Text())
			if err != nil {
				return err
			}
			domainToCount[domain]++
		}
	}

	type kv struct {
		k string
		v int
	}
	var sorted []kv
	for k, v := range domainToCount {
		sorted = append(sorted, kv{k, v})
	}
	slices.SortFunc(sorted, func(a, b kv) int {
		return cmp.Compare(b.k, a.k) // Descending order
	})
	fmt.Println("Fanart domain counts:")
	for _, kv := range sorted {
		fmt.Printf(" * %s : %d\n", kv.k, kv.v)
	}
	return nil
}

func countDirsByFanartSource(stats []*dirStats) error {
	sourceToCount := make(map[string]int)
	fanartTagPath, err := etree.CompilePath("/movie/fanart/thumb")
	if err != nil {
		return err
	}
	tmdbClient := tmdb.ClientOptions{
		APIReadAccessToken: os.Getenv("TMDB_READ_ACCESS_TOKEN"),
	}.NewClient()
	for _, stat := range stats {
		tmdbId, err := stat.Nfo.TmdbId()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not get TMDB ID for %s: %v\n", stat.Dir.Name(), err)
			continue
		}
		movie, err := tmdb.GetMovie(context.Background(), tmdbClient, tmdbId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not fetch movie details for %s: %v\n", stat.Dir.Name(), err)
			continue
		}
		for _, elem := range stat.Nfo.Doc.FindElementsPath(fanartTagPath) {
			fanartUrl := elem.Text()
			if backdropPath, err := movie.BackdropPath(); err == nil && strings.HasSuffix(fanartUrl, backdropPath) {
				sourceToCount["TMDB Backdrop"]++
				continue
			}
			if posterPath, err := movie.PosterPath(); err == nil && strings.HasSuffix(fanartUrl, posterPath) {
				sourceToCount["TMDB Poster"]++
				continue
			}
		}
	}

	type kv struct {
		k string
		v int
	}
	var sorted []kv
	for k, v := range sourceToCount {
		sorted = append(sorted, kv{k, v})
	}
	slices.SortFunc(sorted, func(a, b kv) int {
		return cmp.Compare(b.k, a.k) // Descending order
	})
	fmt.Println("Fanart source counts:")
	for _, kv := range sorted {
		fmt.Printf(" * %s : %d\n", kv.k, kv.v)
	}

	return nil
}

func fanart() error {
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

	countDirsWithFanart(stats)
	fmt.Println()
	if err := countDirsByFanartDomain(stats); err != nil {
		return err
	}
	fmt.Println()
	if err := countDirsByFanartSource(stats); err != nil {
		return err
	}

	return nil
}
