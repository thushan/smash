package slicer

import (
	"encoding/gob"
	"io"
	"io/fs"
	"log"

	"github.com/logrusorgru/aurora/v3"
	"github.com/thushan/smash/internal/algorithms"
)

type Slicer struct {
	buffer       []byte
	defaultBytes []byte
	slices       int
	sliceSize    uint64
	threshold    uint64
	algorithm    algorithms.Algorithm
}

type MetaSlice struct {
	Size uint64
}

const DefaultSlices = 8
const DefaultSliceSize = 8 * 1024
const DefaultThreshold = 100 * 1024

func New(algorithm algorithms.Algorithm) Slicer {
	return NewConfigured(algorithm, DefaultSlices, DefaultSliceSize, DefaultThreshold)
}
func NewConfigured(algorithm algorithms.Algorithm, slices int, size, maxSlice uint64) Slicer {
	return Slicer{
		slices:       slices,
		sliceSize:    size,
		threshold:    maxSlice,
		algorithm:    algorithm,
		buffer:       make([]byte, size),
		defaultBytes: []byte{},
	}
}
func (slicer *Slicer) SliceFS(fs fs.FS, name string, disableSlicing bool) ([]byte, uint64, error) {

	f, err := fs.Open(name)

	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(name, " failed to close because of ", aurora.Red(err))
		}
	}()

	if err != nil {
		return slicer.defaultBytes, 0, err
	}

	fi, err := f.Stat()

	if err != nil {
		return slicer.defaultBytes, 0, err
	}
	size := fi.Size()

	if fr, ok := f.(io.ReaderAt); ok {
		sr := io.NewSectionReader(fr, 0, size)
		slice, err := slicer.Slice(sr, disableSlicing)
		return slice, uint64(size), err
	} else {
		return slicer.defaultBytes, uint64(size), nil
	}
}
func (slicer *Slicer) Slice(sr *io.SectionReader, disableSlicing bool) ([]byte, error) {

	/*
		Check the bytes are within the threshold for a full blob hash.
			OR blob is determined to be text OR slicing is disabled
				Hash the full blob
		Split the blob into n slices for the bytes between the head & tail of the blob (n + 2 = totalSlices)
			Read the head of the blob to n Bytes (slice1)
				Read the X slice to n Bytes (sliceX)
			Read the tail of the slice to n Bytes (slice3)

		Example (offset on the right):
			slice_size:= 8196
			file_size := 1024000
			file_head := 0
			slices[3]
				slice[0] := 251904
				slice[1] := 503808
				slice[2] := 755712
			file_tail := 1015804
	*/

	size := uint64(sr.Size())
	meta := MetaSlice{Size: size}
	algo := slicer.algorithm.New()
	algo.Reset()

	// TODO: Detect text documents and force full hash
	if size < slicer.threshold || slicer.slices <= 0 || disableSlicing /* force full hashes */ {
		if _, err := io.Copy(algo, sr); err != nil {
			return slicer.defaultBytes, err
		}
	} else {
		slice := slicer.buffer
		midSize := size - (slicer.sliceSize * 2)
		sliceFirstSize := int64(midSize / uint64(slicer.slices+1))

		offset := int64(0)

		// head
		sr.Seek(offset, io.SeekStart)
		sr.Read(slice)
		algo.Write(slice)

		// mid-slice crisis
		for i := 0; i < slicer.slices; i++ {
			offset += sliceFirstSize
			sr.Seek(offset, io.SeekCurrent)
			sr.Read(slice)
			algo.Write(slice)
		}

		// tail
		tailOffset := int64(-slicer.sliceSize)
		sr.Seek(tailOffset, io.SeekEnd)
		sr.Read(slice)
		algo.Write(slice)

		// meta
		enc := gob.NewEncoder(algo)
		enc.Encode(meta)
	}
	return algo.Sum(nil), nil
}
