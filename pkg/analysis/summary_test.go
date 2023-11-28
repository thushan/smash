package analysis

import (
	"reflect"
	"testing"
)

var files = []Item{
	{Key: "file1.txt", Size: 1024},
	{Key: "file2.txt", Size: 2024},
	{Key: "file3.txt", Size: 3024},
	{Key: "file4.txt", Size: 4024},
	{Key: "file5.txt", Size: 5024},
	{Key: "file6.txt", Size: 6024},
	{Key: "file7.txt", Size: 7024},
	{Key: "file8.txt", Size: 8024},
	{Key: "file9.txt", Size: 9024},
	{Key: "file10.txt", Size: 10024},
	{Key: "file11.txt", Size: 11024},
	{Key: "file12.txt", Size: 12024},
}

func TestNextIteratorReturnsTop5Files(t *testing.T) {
	TopNumber := 5
	actual := []Item{}
	expected := make([]Item, TopNumber)
	tops := NewSummary(TopNumber)

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
	expected[0] = files[7]
	expected[1] = files[8]
	expected[2] = files[9]
	expected[3] = files[10]
	expected[4] = files[11]

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v files", expected, actual)
	}
}

func TestAllReturnsTop5Files(t *testing.T) {
	TopNumber := 5
	expected := make([]Item, TopNumber)
	tops := NewSummary(TopNumber)

	for _, file := range files {
		tops.Add(file)
	}

	expected[0] = files[7]
	expected[1] = files[8]
	expected[2] = files[9]
	expected[3] = files[10]
	expected[4] = files[11]

	actual := tops.All()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v files", expected, actual)
	}
}
