package main

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/beevik/etree"
)

func mpaa() error {
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

	certPath, err := etree.CompilePath("/movie/certification")
	if err != nil {
		return err
	}
	mpaaPath, err := etree.CompilePath("/movie/mpaa")
	if err != nil {
		return err
	}
	nc17 := []*movieDir{}

	counts := make(map[string]int)
	for _, stat := range stats {
		var certValue, mpaaValue string
		if certs := stat.Nfo.Doc.FindElementsPath(certPath); len(certs) != 1 {
			certValue = fmt.Sprintf("%d certifications", len(certs))
		} else if certValue = certs[0].Text(); len(certValue) == 0 {
			certValue = "empty certification"
		}

		if mpaas := stat.Nfo.Doc.FindElementsPath(mpaaPath); len(mpaas) != 1 {
			mpaaValue = fmt.Sprintf("%d mpaa ratings", len(mpaas))
		} else if mpaaValue = mpaas[0].Text(); len(mpaaValue) == 0 {
			mpaaValue = "empty mpaa"
		}

		if strings.Contains(strings.ToLower(mpaaValue), "nc-17") {
			nc17 = append(nc17, &stat.Dir)
		}

		key := fmt.Sprintf("certification=%s, mpaa=%s", certValue, mpaaValue)
		counts[key]++
	}

	type kv struct {
		k string
		v int
	}
	kvs := make([]kv, 0, len(counts))
	for k, v := range counts {
		kvs = append(kvs, kv{k, v})
	}
	slices.SortFunc(kvs, func(a, b kv) int {
		return cmp.Compare(a.k, b.k)
	})
	fmt.Println("MPAA/Certification counts:")
	for _, kv := range kvs {
		fmt.Printf(" * %s : %d\n", kv.k, kv.v)
	}

	fmt.Println()
	fmt.Println("Examples of NC-17 ratings:")
	for _, dir := range nc17 {
		fmt.Printf(" * %s\n", dir.Name())
	}

	return nil
}