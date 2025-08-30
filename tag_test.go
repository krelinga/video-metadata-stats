package main

import (
	"testing"
)

func TestNewTagPath(t *testing.T) {
	p := newTagPath("a", "b", "c")
	if p.String() != "a/b/c" {
		t.Errorf("expected 'a/b/c', got %q", p.String())
	}
	parts := p.Parts()
	if len(parts) != 3 {
		t.Errorf("expected 3 parts, got %d", len(parts))
	}
	if parts[0] != "a" || parts[1] != "b" || parts[2] != "c" {
		t.Errorf("unexpected parts: %v", parts)
	}
}

func TestTagPathParts(t *testing.T) {
	p := newTagPath("a/b/c")
	parts := p.Parts()
	if len(parts) != 1 {
		t.Errorf("expected 1 part, got %d", len(parts))
	}
	if parts[0] != "a/b/c" {
		t.Errorf("unexpected parts: %v", parts)
	}
}

func TestTagPathAppend(t *testing.T) {
	p := newTagPath("a", "b")
	p = p.Append("c", "d")
	if p.String() != "a/b/c/d" {
		t.Errorf("expected 'a/b/c/d', got %q", p.String())
	}
}
