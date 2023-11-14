package slicer

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestSlicingSupported_OnRandom1MbFile(t *testing.T) {
	fsSize := 1024000 // 1Mb
	binary := randomBytes(fsSize)
	reader := bytes.NewReader(binary)

	expected := true
	runSlicingSupportTest(reader, int64(fsSize), expected, t)
}

func TestSlicingSupported_OnRandomTiny3byteFile(t *testing.T) {
	fsSize := 3 // 3 bytes and that's all!
	binary := randomBytes(fsSize)
	reader := bytes.NewReader(binary)

	expected := false

	runSlicingSupportTest(reader, int64(fsSize), expected, t)
}

func TestSlicingSupported_OnTextBytesFile(t *testing.T) {
	binary := []byte("OMG THIS IS TEXT!")
	fsSize := len(binary)
	reader := bytes.NewReader(binary)

	expected := true

	runSlicingSupportTest(reader, int64(fsSize), expected, t)
}

func TestSlicingSupported_OnTextStringFile(t *testing.T) {
	stexty := "OMG THIS IS TEXT!"
	fsSize := len(stexty)
	reader := strings.NewReader(stexty)

	expected := true

	runSlicingSupportTest(reader, int64(fsSize), expected, t)
}

func TestSlicingSupported_OnInvalidBinaryFile(t *testing.T) {
	filename := "./artefacts/test.1mb"
	if binary, err := os.ReadFile(filename); err != nil {
		t.Errorf("Unexpected io error %v", err)
	} else {
		fsSize := len(binary)
		reader := bytes.NewReader(binary)
		expected := true
		runSlicingSupportTest(reader, int64(fsSize), expected, t)
	}
}
func runSlicingSupportTest(reader io.ReaderAt, size int64, expected bool, t *testing.T) {

	sr := io.NewSectionReader(reader, 0, size)

	actual := slicingSupported(sr, uint64(size))

	if actual != expected {
		t.Errorf("expected slicing supported %t, got %t", expected, actual)
	}

}
