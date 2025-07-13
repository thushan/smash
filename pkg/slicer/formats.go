package slicer

import (
	"io"
	"sync"
	"unicode/utf8"

	"golang.org/x/tools/godoc/util"
)

const MaxReadBytes = 1024

var detectBufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, MaxReadBytes)
	},
}

func slicingSupported(sr *io.SectionReader, size uint64) bool {

	if size < utf8.UTFMax {
		// It's a tiny file, we should full hash this puppy
		return false
	}

	bufInterface := detectBufferPool.Get()
	buf, ok := bufInterface.([]byte)
	if !ok {
		return false
	}
	defer detectBufferPool.Put(buf)

	if _, err := sr.Read(buf); err != nil && util.IsText(buf) {
		return false
	} else {
		return true
	}
}
