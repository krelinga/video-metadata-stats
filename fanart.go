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

	statWithFanartMatchingAnyImageBackdrop := 0
	statWithoutFanartMatchingAnyImageBackdrop := 0

statloop:
	for _, stat := range stats {
		tmdbId, err := stat.Nfo.TmdbId()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not get TMDB ID for %s: %v\n", stat.Dir.Name(), err)
			continue
		}
		movie, err := tmdb.GetMovie(context.Background(), tmdbClient, tmdbId, tmdb.WithAppendToResponse("images"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not fetch movie details for %s: %v\n", stat.Dir.Name(), err)
			continue
		}
		for _, elem := range stat.Nfo.Doc.FindElementsPath(fanartTagPath) {
			fanartUrl := elem.Text()
			movieHasBackdrop := func() bool {
				_, err := movie.BackdropPath()
				return err == nil
			}()
			movieImagesHaveAnyBackdrop := func() bool {
				if images, err := movie.Images(); err != nil {
					return false
				} else if backdrops, err := images.Backdrops(); err != nil {
					return false
				} else {
					return len(backdrops) > 0
				}
			}()
			fanartMatchesAnyImageBackdrop := func() bool {
				if images, err := movie.Images(); err != nil {
					return false
				} else if backdrops, err := images.Backdrops(); err != nil {
					return false
				} else {
					for _, backdrop := range backdrops {
						if filepath, err := backdrop.FilePath(); err == nil && strings.HasSuffix(fanartUrl, filepath) {
							return true
						}
					}
				}
				return false
			}()
			if fanartMatchesAnyImageBackdrop {
				statWithFanartMatchingAnyImageBackdrop++
			} else {
				statWithoutFanartMatchingAnyImageBackdrop++
			}

			if backdropPath, err := movie.BackdropPath(); err == nil && strings.HasSuffix(fanartUrl, backdropPath) {
				key := fmt.Sprintf("TMDB Backdrop (Movie Level), movieHasBackdrop=%t, movieImagesHaveAnyBackdrop=%t", movieHasBackdrop, movieImagesHaveAnyBackdrop)
				sourceToCount[key]++
				continue
			}
			if posterPath, err := movie.PosterPath(); err == nil && strings.HasSuffix(fanartUrl, posterPath) {
				sourceToCount["TMDB Poster"]++
				continue
			}
			if images, err := movie.Images(); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Could not fetch images for %s: %v\n", stat.Dir.Name(), err)
			} else {
				if backdrops, err := images.Backdrops(); err == nil {
					for _, backdrop := range backdrops {
						if filepath, err := backdrop.FilePath(); err == nil && strings.HasSuffix(fanartUrl, filepath) {
							key := fmt.Sprintf("TMDB Backdrop (Additional Images), movieHasBackdrop=%t, movieImagesHaveAnyBackdrop=%t", movieHasBackdrop, movieImagesHaveAnyBackdrop)
							sourceToCount[key]++
							continue statloop
						}
					}
				}
				if posters, err := images.Posters(); err == nil {
					for _, poster := range posters {
						if filepath, err := poster.FilePath(); err == nil && strings.HasSuffix(fanartUrl, filepath) {
							sourceToCount["TMDB Poster (Additional Images)"]++
							continue statloop
						}
					}
				}
				if logos, err := images.Logos(); err == nil {
					for _, logo := range logos {
						if filepath, err := logo.FilePath(); err == nil && strings.HasSuffix(fanartUrl, filepath) {
							sourceToCount["TMDB Logo (Additional Images)"]++
							continue statloop
						}
					}
				}
			}
			sourceToCount["Unknown Source"]++
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
	fmt.Println()
	fmt.Printf("Stats: Directories with fanart matching any image backdrop: %d\n", statWithFanartMatchingAnyImageBackdrop)
	fmt.Printf("Stats: Directories without fanart matching any image backdrop: %d\n", statWithoutFanartMatchingAnyImageBackdrop)

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
