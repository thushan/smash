package report

import (
	"encoding/hex"
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

func SummariseSmashedFile(stats slicer.SlicerStats, filename string, ms int64, duplicates *haxmap.Map[string, *DuplicateFiles], empty *EmptyFiles) {
	file := SmashFile{
		Hash:        hex.EncodeToString(stats.Hash),
		Filename:    filename,
		FileSize:    stats.FileSize,
		FullHash:    stats.HashedFullFile,
		EmptyFile:   stats.EmptyFile,
		FileSizeF:   humanize.Bytes(stats.FileSize),
		ElapsedTime: ms,
	}
	if file.EmptyFile {
		empty.Lock()
		empty.Files = append(empty.Files, file)
		empty.Unlock()
	} else {
		dupes, _ := duplicates.GetOrSet(file.Hash, &DuplicateFiles{})
		dupes.Lock()
		dupes.Files = append(dupes.Files, file)
		dupes.Unlock()
	}

}
