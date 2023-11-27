package report

import (
	"reflect"
	"testing"
)

var files = []SmashFile{
	{
		Filename:    "file1.txt",
		Hash:        "",
		FileSizeF:   "",
		FileSize:    1024,
		ElapsedTime: 0,
		FullHash:    false,
		EmptyFile:   false,
	},
	{
		Filename:    "file2.txt",
		Hash:        "",
		FileSizeF:   "",
		FileSize:    2042,
		ElapsedTime: 0,
		FullHash:    false,
		EmptyFile:   false,
	},
	{
		Filename:    "file3.txt",
		Hash:        "",
		FileSizeF:   "",
		FileSize:    1942,
		ElapsedTime: 0,
		FullHash:    false,
		EmptyFile:   false,
	},
	{
		Filename:    "file4.txt",
		Hash:        "",
		FileSizeF:   "",
		FileSize:    1984,
		ElapsedTime: 0,
		FullHash:    false,
		EmptyFile:   false,
	},
	{
		Filename:    "file5.txt",
		Hash:        "",
		FileSizeF:   "",
		FileSize:    1992,
		ElapsedTime: 0,
		FullHash:    false,
		EmptyFile:   false,
	},
	{
		Filename:    "file6.txt",
		Hash:        "",
		FileSizeF:   "",
		FileSize:    2002,
		ElapsedTime: 0,
		FullHash:    false,
		EmptyFile:   false,
	},
	{
		Filename:    "file7.txt",
		Hash:        "",
		FileSizeF:   "",
		FileSize:    2007,
		ElapsedTime: 0,
		FullHash:    false,
		EmptyFile:   false,
	},
	{
		Filename:    "file8.txt",
		Hash:        "",
		FileSizeF:   "",
		FileSize:    2020,
		ElapsedTime: 0,
		FullHash:    false,
		EmptyFile:   false,
	},
	{
		Filename:    "file9.txt",
		Hash:        "",
		FileSizeF:   "",
		FileSize:    1957,
		ElapsedTime: 0,
		FullHash:    false,
		EmptyFile:   false,
	},
	{
		Filename:    "file10.txt",
		Hash:        "",
		FileSizeF:   "",
		FileSize:    1953,
		ElapsedTime: 0,
		FullHash:    false,
		EmptyFile:   false,
	},
}

func TestNextReturnsTop5Files(t *testing.T) {
	TopNumber := 5
	actual := []SmashFile{}
	expected := make([]SmashFile, TopNumber)
	tops := NewTopFilesSummary(TopNumber)

	for _, file := range files {
		tops.Add(file)
	}

	for {
		file, ok := tops.Next()
		if !ok {
			break
		}
		actual = append(actual, file)
	}
	expected[0] = files[4]
	expected[1] = files[5]
	expected[2] = files[6]
	expected[3] = files[7]
	expected[4] = files[1]

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v files", expected, actual)
	}
}

func TestAllReturnsTop5Files(t *testing.T) {
	TopNumber := 5
	expected := make([]SmashFile, TopNumber)
	tops := NewTopFilesSummary(TopNumber)

	for _, file := range files {
		tops.Add(file)
	}

	expected[0] = files[4]
	expected[1] = files[5]
	expected[2] = files[6]
	expected[3] = files[7]
	expected[4] = files[1]

	actual := tops.All()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v files", expected, actual)
	}
}
