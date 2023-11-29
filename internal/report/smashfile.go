package report

import (
	"encoding/hex"
	"github.com/thushan/smash/internal/theme"
	"sync"

	"github.com/alphadose/haxmap"
	"github.com/dustin/go-humanize"
	"github.com/thushan/smash/pkg/slicer"
)

type SmashFile struct {
	Filename    string
	Hash        string
	FileSizeF   string
	FileSize    uint64
	ElapsedTime int64
	FullHash    bool
	EmptyFile   bool
}
type SmashFiles struct {
	sync.RWMutex
	Duplicates []SmashFile
}

func SummariseSmashedFile(stats slicer.SlicerStats, filename string, ms int64, duplicates *haxmap.Map[string, *SmashFiles], emptyFiles *[]SmashFile) {
	sf := SmashFile{
		Hash:        hex.EncodeToString(stats.Hash),
		Filename:    filename,
		FileSize:    stats.FileSize,
		FullHash:    stats.HashedFullFile,
		EmptyFile:   stats.EmptyFile,
		FileSizeF:   humanize.Bytes(stats.FileSize),
		ElapsedTime: ms,
	}
	if sf.EmptyFile {
		*emptyFiles = append(*emptyFiles, sf)
	} else {
		hash := sf.Hash
		files := &SmashFiles{}
		files.Duplicates = []SmashFile{sf}
		files.Lock()
		if of, loaded := duplicates.GetOrSet(hash, files); loaded {
			of.Lock()
			of.Duplicates = append(of.Duplicates, sf)
			of.Unlock()
			if ov, swapped := duplicates.Swap(hash, of); !swapped {
				theme.Error.Println("Swap failed for ", hash, ". old: ", len(ov.Duplicates), " | new: ", len(of.Duplicates))
			}
		}
		files.Unlock()
	}

}
