package main

import (
	"net/url"
	"strings"
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
