package report

import (
	"fmt"
	"testing"
)

func TestTop5FilesReturned(t *testing.T) {
	files := []SmashFile{
		{
			Filename:    "file1.txt",
			Hash:        "",
			FileSizeF:   "",
			FileSize:    1,
			ElapsedTime: 0,
			FullHash:    false,
			EmptyFile:   false,
		},
		{
			Filename:    "file2.txt",
			Hash:        "",
			FileSizeF:   "",
			FileSize:    2,
			ElapsedTime: 0,
			FullHash:    false,
			EmptyFile:   false,
		},
		{
			Filename:    "file3.txt",
			Hash:        "",
			FileSizeF:   "",
			FileSize:    3,
			ElapsedTime: 0,
			FullHash:    false,
			EmptyFile:   false,
		},
		{
			Filename:    "file4.txt",
			Hash:        "",
			FileSizeF:   "",
			FileSize:    4,
			ElapsedTime: 0,
			FullHash:    false,
			EmptyFile:   false,
		},
		{
			Filename:    "file5.txt",
			Hash:        "",
			FileSizeF:   "",
			FileSize:    5,
			ElapsedTime: 0,
			FullHash:    false,
			EmptyFile:   false,
		},
		{
			Filename:    "file6.txt",
			Hash:        "",
			FileSizeF:   "",
			FileSize:    6,
			ElapsedTime: 0,
			FullHash:    false,
			EmptyFile:   false,
		},
		{
			Filename:    "file7.txt",
			Hash:        "",
			FileSizeF:   "",
			FileSize:    7,
			ElapsedTime: 0,
			FullHash:    false,
			EmptyFile:   false,
		},
		{
			Filename:    "file8.txt",
			Hash:        "",
			FileSizeF:   "",
			FileSize:    8,
			ElapsedTime: 0,
			FullHash:    false,
			EmptyFile:   false,
		},
		{
			Filename:    "file9.txt",
			Hash:        "",
			FileSizeF:   "",
			FileSize:    9,
			ElapsedTime: 0,
			FullHash:    false,
			EmptyFile:   false,
		},
		{
			Filename:    "file10.txt",
			Hash:        "",
			FileSizeF:   "",
			FileSize:    10,
			ElapsedTime: 0,
			FullHash:    false,
			EmptyFile:   false,
		},
		{
			Filename:    "file11.txt",
			Hash:        "",
			FileSizeF:   "",
			FileSize:    11,
			ElapsedTime: 0,
			FullHash:    false,
			EmptyFile:   false,
		},
	}
	tops := TopFiles{}

	for _, file := range files {
		tops.Index(file)
	}

	it := tops.Iterator()

	for {
		file, ok := it()
		if !ok {
			break
		}
		fmt.Printf("%s - Size: %d\n", file.Filename, file.FileSize)
	}

}
