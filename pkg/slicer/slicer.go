package slicer

import (
	"encoding/gob"
	"errors"
	"io"
	"io/fs"

	"github.com/thushan/smash/internal/algorithms"
)

type Slicer struct {
	defaultBytes []byte
	slices       int
	sliceSize    uint64
	threshold    uint64
	algorithm    algorithms.Algorithm
}

type SlicerStats struct {
	SliceOffsets   map[int]int64
	Filename       string
	Hash           []byte
	ReaderSize     int64
	SliceOffset    int64
	MidSize        uint64
	SliceSize      uint64
	FileSize       uint64
	Slices         int
	HashedFullFile bool
}

type MetaSlice struct {
	Size uint64
}

const DefaultSlices = 4
const DefaultSliceSize = 8 * 1024
const DefaultThreshold = 100 * 1024
const DefaultMinimumSize = (DefaultSlices + 2) * DefaultSliceSize

func New(algorithm algorithms.Algorithm) Slicer {
	return NewConfigured(algorithm, DefaultSlices, DefaultSliceSize, DefaultThreshold)
}
func NewConfigured(algorithm algorithms.Algorithm, slices int, size, maxSlice uint64) Slicer {
	return Slicer{
		slices:       slices,
		sliceSize:    size,
		threshold:    maxSlice,
		algorithm:    algorithm,
		defaultBytes: []byte{},
	}
}
func (slicer *Slicer) SliceFS(fs fs.FS, name string, disableSlicing bool) (SlicerStats, error) {

	stats := SlicerStats{Hash: slicer.defaultBytes, Filename: name}
	f, err := fs.Open(name)
	defer f.Close()

	if err != nil {
		return stats, err
	}

	fi, err := f.Stat()

	if err != nil {
		return stats, err
	}

	size := fi.Size()

	stats.FileSize = uint64(size)
	stats.Slices = slicer.slices
	stats.SliceSize = slicer.sliceSize

	if fr, ok := f.(io.ReaderAt); ok {
		sr := io.NewSectionReader(fr, 0, size)
		err := slicer.Slice(sr, disableSlicing, false, &stats)
		return stats, err
	} else {
		return stats, errors.New("the File System does not support readers")
	}
}
func (slicer *Slicer) Slice(sr *io.SectionReader, disableSlicing bool, disableMeta bool, stat *SlicerStats) error {

	/*
		Check the bytes are within the threshold for a full blob hash.
			OR blob is determined to be text OR slicing is disabled
				Hash the full blob
		Split the blob into n slices for the bytes between the head & tail of the blob (n + 2 = totalSlices)
			Read the head of the blob to n Bytes (slice1)
				Read the X slice to n Bytes (sliceX)
			Read the tail of the slice to n Bytes (slice3)

		Example (offset on the right):
			slices    := 4
			slice_size:= 8196
			file_size := 1024000
			file_head := 0
			slices[4]
				slice[0] := 251,904
				\_reader := 260,096
				slice[1] := 503,808
				\_reader := 512,000
				slice[2] := 755,712
				\_reader := 763,904
				slice[3] := 1,007,616
				\_reader := 1,015,808
			file_tail :=  1,015,808
	*/

	size := uint64(sr.Size())
	meta := MetaSlice{Size: size}

	algo := slicer.algorithm.New()
	algo.Reset()

	stat.ReaderSize = sr.Size()

	// checks
	greaterThanMinimumFileSize := uint64(slicer.slices+2)*slicer.sliceSize > size
	greaterThanMinimumThreshold := size < slicer.threshold
	invalidNumberOfSlices := slicer.slices <= 0
	fullFileHash := disableSlicing || greaterThanMinimumThreshold || greaterThanMinimumFileSize || invalidNumberOfSlices

	stat.HashedFullFile = fullFileHash

	// TODO: Detect text documents and force full hash
	if fullFileHash {
		if _, err := io.Copy(algo, sr); err != nil {
			return err
		}
	} else {
		slice := make([]byte, slicer.sliceSize)
		midSize := size - (slicer.sliceSize * 2)
		sliceOffset := int64((midSize / uint64(slicer.slices)) - slicer.sliceSize)

		stat.SliceOffset = sliceOffset
		stat.MidSize = midSize
		stat.SliceOffsets = make(map[int]int64)

		// head
		stat.SliceOffsets[0] = 0
		if _, err := sr.Seek(0, io.SeekStart); err != nil {
			return err
		}
		if _, err := sr.Read(slice); err != nil {
			return err
		}
		algo.Write(slice)

		// mid-slice crisis
		for i := 0; i < slicer.slices; i++ {
			if offset, err := sr.Seek(sliceOffset, io.SeekCurrent); err != nil {
				return err
			} else {
				stat.SliceOffsets[i+1] = offset
			}
			if _, err := sr.Read(slice); err != nil {
				return err
			}

			algo.Write(slice)
		}

		// tail
		tailOffset := int64(-slicer.sliceSize)
		if offset, err := sr.Seek(tailOffset, io.SeekEnd); err != nil {
			return err
		} else {
			stat.SliceOffsets[len(stat.SliceOffsets)] = offset
		}
		if _, err := sr.Read(slice); err != nil {
			return err
		}
		algo.Write(slice)

		// metadata
		if !disableMeta {
			enc := gob.NewEncoder(algo)
			if err := enc.Encode(meta); err != nil {
				return err
			}
		}
	}
	stat.Hash = algo.Sum(nil)
	return nil
}
