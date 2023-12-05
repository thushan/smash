package slicer

import (
	"io"
	"unicode/utf8"

	"golang.org/x/tools/godoc/util"
)

const MaxReadBytes = 1024

var detectBuffer []byte

func init() {
	detectBuffer = make([]byte, MaxReadBytes)
}
func slicingSupported(sr *io.SectionReader, size uint64) bool {

	if size < utf8.UTFMax {
		// It's a tiny file, we should full hash this puppy
		return false
	}

	if _, err := sr.Read(detectBuffer); err != nil && util.IsText(detectBuffer) {
		return false
	} else {
		return true
	}
}
