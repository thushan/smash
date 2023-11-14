package slicer

import (
	"io"
	"unicode/utf8"

	"golang.org/x/tools/godoc/util"
)

var bomPdf = []byte("%PDF")

func slicingSupported(sr *io.SectionReader, size uint64) bool {
	const MaxReadBytes = 8

	if size < utf8.UTFMax {
		// It's a tiny file, we should full hash this puppy
		return false
	}

	buf := make([]byte, MaxReadBytes)

	isTextFile := util.IsText(buf)

	if _, err := sr.Read(buf); err != nil ||
		isTextFile {
		return false
	} else {
		return true
	}
}
