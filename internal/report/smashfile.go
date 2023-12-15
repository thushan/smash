package report

import (
	"encoding/hex"
	"sync"

	"github.com/thushan/smash/pkg/indexer"

	"github.com/puzpuzpuz/xsync/v3"

	"github.com/dustin/go-humanize"
	"github.com/thushan/smash/pkg/slicer"
)

type SmashFile struct {
	Filename    string
	Location    string
	Path        string
	Base        string
	Hash        string
	FileSizeF   string
	FileSize    uint64
	ElapsedTime int64
	FullHash    bool
	EmptyFile   bool
}
type DuplicateFiles struct {
	Files []SmashFile
	sync.RWMutex
}
type EmptyFiles struct {
	Files []SmashFile
	sync.RWMutex
}

func SummariseSmashedFile(stats slicer.SlicerStats, ffs *indexer.FileFS, ms int64, duplicates *xsync.MapOf[string, *DuplicateFiles], empty *EmptyFiles) {
	file := SmashFile{
		Hash:        hex.EncodeToString(stats.Hash),
		Filename:    ffs.Name,
		Location:    ffs.Location,
		Path:        ffs.Path,
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
		dupes, _ := duplicates.LoadOrStore(file.Hash, &DuplicateFiles{
			Files:   []SmashFile{},
			RWMutex: sync.RWMutex{},
		})
		dupes.Lock()
		dupes.Files = append(dupes.Files, file)
		dupes.Unlock()
	}

}
