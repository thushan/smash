package slicer

import (
	"encoding/gob"
	"errors"
	"io"
	"io/fs"
	"os"
	"sync"

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
	EmptyFile      bool
	IgnoredFile    bool
	HashedFullFile bool
}

type MetaSlice struct {
	Size uint64
}
type Options struct {
	MinSize         uint64
	MaxSize         uint64
	DisableSlicing  bool
	DisableMeta     bool
	DisableAutoText bool
}

const MaxSlices = 128
const DefaultSlices = 4
const DefaultSliceSize = 8 * 1024
const DefaultThreshold = 100 * 1024
const DefaultMinimumSize = (DefaultSlices + 2) * DefaultSliceSize
const DefaultMinSize = 0
const DefaultMaxSize = 0

var sliceBufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, DefaultSliceSize)
	},
}

func getSliceBuffer(size uint64) []byte {
	if size == DefaultSliceSize {
		bufInterface := sliceBufferPool.Get()
		if buf, ok := bufInterface.([]byte); ok {
			return buf
		}
	}
	// Prevent panic from unreasonable sizes
	const maxReasonableSize = 1 << 30 // 1GB
	if size > maxReasonableSize {
		return nil
	}
	return make([]byte, size)
}

func putSliceBuffer(buf []byte) {
	if cap(buf) == DefaultSliceSize {
		sliceBufferPool.Put(buf)
	}
}

func New(algorithm algorithms.Algorithm) Slicer {
	return NewConfigured(algorithm, DefaultSlices, DefaultSliceSize, DefaultThreshold)
}
func NewConfigured(algorithm algorithms.Algorithm, slices int, size, threshold uint64) Slicer {
	return Slicer{
		slices:       slices,
		sliceSize:    size,
		threshold:    threshold,
		algorithm:    algorithm,
		defaultBytes: []byte{},
	}
}
func (slicer *Slicer) SliceFS(fileSystem fs.FS, name string, options *Options) (SlicerStats, error) {

	stats := SlicerStats{Hash: slicer.defaultBytes, Filename: name}
	fio, ferr := fs.Stat(fileSystem, name)

	if ferr != nil {
		return stats, ferr
	}

	fileSize := fio.Size()
	if fileSize < 0 {
		return stats, errors.New("file size cannot be negative")
	}
	size := uint64(fileSize)
	isEmptyFile := size == 0

	if !shouldAnalyseBasedOnSize(size, options.MinSize, options.MaxSize) ||
		shouldIgnoreFileMode(fio) ||
		isEmptyFile {
		stats.IgnoredFile = true
		stats.EmptyFile = isEmptyFile
		stats.Hash = nil
		return stats, nil
	}

	f, err := fileSystem.Open(name)
	if err != nil {
		return stats, err
	}
	defer func(fs io.Closer) {
		if fs != nil {
			_ = fs.Close()
		}
	}(f)

	stats.FileSize = size
	stats.Slices = slicer.slices
	stats.SliceSize = slicer.sliceSize

	if fr, ok := f.(io.ReaderAt); ok {
		sr := io.NewSectionReader(fr, 0, fileSize)
		err := slicer.Slice(sr, options, &stats)
		return stats, err
	} else {
		return stats, errors.New("the File System does not support readers")
	}
}
func (slicer *Slicer) Slice(sr *io.SectionReader, options *Options, stats *SlicerStats) error {

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

	srSize := sr.Size()
	if srSize < 0 {
		return errors.New("section reader size cannot be negative")
	}
	size := uint64(srSize)
	isEmptyFile := size == 0

	if !shouldAnalyseBasedOnSize(size, options.MinSize, options.MaxSize) ||
		isEmptyFile {
		stats.IgnoredFile = true
		stats.EmptyFile = isEmptyFile
		stats.Hash = nil
		return nil
	}

	algo := slicer.algorithm.New()
	algo.Reset()

	stats.ReaderSize = sr.Size()

	// checks
	canSliceFile := !options.DisableAutoText && slicingSupported(sr, size)
	slicesPlus2 := slicer.slices + 2
	if slicesPlus2 < 0 {
		return errors.New("slices overflow")
	}
	greaterThanMinimumFileSize := uint64(slicesPlus2)*slicer.sliceSize > size
	greaterThanMinimumThreshold := size < slicer.threshold
	invalidNumberOfSlices := slicer.slices <= 0
	// fullHash only those times we have to
	fullHash := options.DisableSlicing ||
		greaterThanMinimumThreshold ||
		greaterThanMinimumFileSize ||
		invalidNumberOfSlices ||
		!canSliceFile

	stats.HashedFullFile = fullHash

	// Reset after text detection
	_, _ = sr.Seek(0, io.SeekStart)

	if fullHash {
		if _, err := io.Copy(algo, sr); err != nil {
			return err
		}
	} else {
		slice := getSliceBuffer(slicer.sliceSize)
		if slice == nil {
			return errors.New("slice size too large to allocate buffer")
		}
		defer putSliceBuffer(slice)

		midSize := size - (slicer.sliceSize * 2)
		if slicer.slices <= 0 {
			return errors.New("invalid number of slices")
		}
		divResult := midSize / uint64(slicer.slices)
		if divResult < slicer.sliceSize {
			return errors.New("slice offset would be negative")
		}
		sliceOffsetCalc := divResult - slicer.sliceSize
		const maxInt64 = 1<<63 - 1
		if sliceOffsetCalc > maxInt64 {
			return errors.New("slice offset overflow")
		}
		sliceOffset := int64(sliceOffsetCalc)

		stats.SliceOffset = sliceOffset
		stats.MidSize = midSize
		stats.SliceOffsets = make(map[int]int64, slicer.slices+2)

		// head
		stats.SliceOffsets[0] = 0
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
				stats.SliceOffsets[i+1] = offset
			}
			if _, err := sr.Read(slice); err != nil {
				return err
			}

			algo.Write(slice)
		}

		// tail
		if slicer.sliceSize > maxInt64 {
			return errors.New("slice size too large for tail offset")
		}
		tailOffset := -int64(slicer.sliceSize)
		if offset, err := sr.Seek(tailOffset, io.SeekEnd); err != nil {
			return err
		} else {
			stats.SliceOffsets[len(stats.SliceOffsets)] = offset
		}
		if _, err := sr.Read(slice); err != nil {
			return err
		}
		algo.Write(slice)

		// metadata
		if !options.DisableMeta {
			enc := gob.NewEncoder(algo)
			meta := MetaSlice{Size: size}
			if err := enc.Encode(meta); err != nil {
				return err
			}
		}
	}
	stats.Hash = algo.Sum(nil)
	return nil
}
func shouldAnalyseBasedOnSize(fileSize, minSize, maxSize uint64) bool {
	if minSize == DefaultMinSize && maxSize == DefaultMaxSize {
		return true
	}
	if minSize != DefaultMinSize && fileSize < minSize {
		return false
	}
	if maxSize != DefaultMaxSize && fileSize > maxSize {
		return false
	}
	return true
}
func shouldIgnoreFileMode(fio os.FileInfo) bool {
	return fio.Mode()&os.ModeNamedPipe != 0 ||
		fio.Mode()&os.ModeSocket != 0 ||
		fio.Mode()&os.ModeDevice != 0 ||
		fio.Mode()&os.ModeSymlink != 0
}
