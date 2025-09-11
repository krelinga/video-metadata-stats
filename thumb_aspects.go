package main

import (
	"cmp"
	"fmt"
	"net/url"
	"regexp"
	"slices"
	"strings"

	"github.com/beevik/etree"
)

func thumbAspects() error {
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
	aspectsMap := make(map[string]int)
	thumbPath, err := etree.CompilePath("/movie/thumb")
	if err != nil {
		return err
	}
	urlPathRegex, err := regexp.Compile(`movies/\d+/([^/]+)/.*$`)
	if err != nil {
		return err
	}
	for _, stat := range stats {
		for _, thumb := range stat.Nfo.Doc.FindElementsPath(thumbPath) {
			aspect := thumb.SelectAttr("aspect")
			if aspect == nil {
				continue
			}
			value := aspect.Value
			if len(value) == 0 {
				continue
			}
			rawThumbUrl := thumb.Text()
			parsed, err := url.Parse(rawThumbUrl)
			if err != nil {
				return err
			}
			var urlPart string
			if strings.Contains(parsed.Hostname(), "image.tmdb.org") {
				urlPart = "tmdb"
			} else {
				matches := urlPathRegex.FindStringSubmatch(parsed.Path)
				if len(matches) != 2 {
					return fmt.Errorf("unexpected thumb URL format: %q in %q", rawThumbUrl, stat.Dir.Name())
				}
				urlPart = matches[1]
			}
			aspectKey := fmt.Sprintf("%s (%s)", value, urlPart)
			aspectsMap[aspectKey]++
		}
	}

	type kv struct {
		k string
		v int
	}
	var sorted []kv
	for k, v := range aspectsMap {
		sorted = append(sorted, kv{k, v})
	}
	slices.SortFunc(sorted, func(a, b kv) int {
		return cmp.Compare(a.k, b.k)
	})

	fmt.Println("Thumb Aspects and Directory counts:")
	for _, kv := range sorted {
		fmt.Printf(" * Aspect: %s, Directories: %d\n", kv.k, kv.v)
	}

	return nil
}
