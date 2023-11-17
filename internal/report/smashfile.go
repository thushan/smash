package report

import (
	"encoding/hex"

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

func SummariseSmashedFile(stats slicer.SlicerStats, filename string, ms int64, duplicates *haxmap.Map[string, []SmashFile], emptyFiles *[]SmashFile) {
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
		if v, existing := duplicates.Get(hash); existing {
			v = append(v, sf)
			duplicates.Set(hash, v)
		} else {
			duplicates.Set(hash, []SmashFile{sf})
		}
	}

}
