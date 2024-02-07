package indexer

import (
	"io/fs"
)

type Kind int

const (
	Local Kind = iota
)

type LocationFS struct {
	fs.FS        // embed the original fs.FS type
	Name  string // add a new field
	Kind  Kind
}

func NewLocationFS(kind Kind, name string, fsys fs.FS) *LocationFS {
	return &LocationFS{
		FS:   fsys,
		Name: name,
		Kind: kind,
	}
}
