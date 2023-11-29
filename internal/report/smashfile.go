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
type DuplicateFiles struct {
	sync.RWMutex
	Files []SmashFile
}
type EmptyFiles struct {
	sync.RWMutex
	Files []SmashFile
}

func SummariseSmashedFile(stats slicer.SlicerStats, filename string, ms int64, dupes *haxmap.Map[string, *DuplicateFiles], empty *EmptyFiles) {
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
		empty.Lock()
		empty.Files = append(empty.Files, sf)
		empty.Unlock()
	} else {
		hash := sf.Hash
		files := &DuplicateFiles{}
		files.Files = []SmashFile{sf}
		files.Lock()
		if of, loaded := dupes.GetOrSet(hash, files); loaded {
			of.Lock()
			of.Files = append(of.Files, sf)
			of.Unlock()
			if ov, swapped := dupes.Swap(hash, of); !swapped {
				theme.Error.Println("Swap failed for ", hash, ". old: ", len(ov.Files), " | new: ", len(of.Files))
			}
		}
		files.Unlock()
	}

}
