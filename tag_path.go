package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/beevik/etree"
)

type tagPath string

func newTagPath(parts ...string) tagPath {
	escaped := make([]string, len(parts))
	for i, p := range parts {
		escaped[i] = url.PathEscape(p)
	}
	return tagPath(strings.Join(escaped, "/"))
}

func (p tagPath) String() string {
	return string(p)
}

func (p tagPath) Parts() []string {
	raw := strings.Split(string(p), "/")
	parts := make([]string, len(raw))
	for i, r := range raw {
		var err error
		parts[i], err = url.PathUnescape(r)
		if err != nil {
			panic(err)
		}
	}
	if len(parts) == 1 && parts[0] == "" {
		return []string{}
	}
	return parts
}

func (p tagPath) Append(parts ...string) tagPath {
	return newTagPath(append(p.Parts(), parts...)...)
}

func (p tagPath) XmlPath() (etree.Path, error) {
	parts := p.Parts()
	if len(parts) == 0 {
		return etree.Path{}, fmt.Errorf("invalid tagPath: %s", p)
	}
	for _, part := range parts {
		if strings.Contains(part, "/") {
			return etree.Path{}, fmt.Errorf("cannot convert tagPath with slashes to etree.Path: %s", p)
		}
	}

	xmlPath, err := etree.CompilePath(string(p))
	if err != nil {
		return etree.Path{}, err
	}
	return xmlPath, nil
}
