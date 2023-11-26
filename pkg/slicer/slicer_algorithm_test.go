// nolint
package slicer

import (
	"bytes"
	"encoding/hex"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/thushan/smash/internal/algorithms"
)

// fieldalignment: struct with 40 pointer bytes could be 24 (govet)
// but this is nicer to see / read :)
var algoData = []struct {
	algorithm      algorithms.Algorithm
	disableSlicing bool
	filename       string
	expectHash     string
}{
	{algorithms.Xxhash, false, "./artefacts/test-manipulated.1mb", "4f595576799edcd9"},
	{algorithms.Xxhash, true, "./artefacts/test-manipulated.1mb", "4a1960f16a88960c"},
	{algorithms.Xxhash, false, "./artefacts/test.1mb", "bb83f43630ee546f"},
	{algorithms.Xxhash, true, "./artefacts/test.1mb", "6b6255ee515dcc04"},
	{algorithms.Murmur3_128, false, "./artefacts/test-manipulated.1mb", "daa0b57d39ab077f56bcdf855753d8dd"},
	{algorithms.Murmur3_128, true, "./artefacts/test-manipulated.1mb", "7b49601fb19613cfa36cc032910228b7"},
	{algorithms.Murmur3_128, false, "./artefacts/test.1mb", "92d0c527266ec9151a6a9239c105df84"},
	{algorithms.Murmur3_128, true, "./artefacts/test.1mb", "35ec8ac6041a7e9b70c61cc30d40b592"},
	{algorithms.Murmur3_64, false, "./artefacts/test-manipulated.1mb", "daa0b57d39ab077f"},
	{algorithms.Murmur3_64, true, "./artefacts/test-manipulated.1mb", "7b49601fb19613cf"},
	{algorithms.Murmur3_64, false, "./artefacts/test.1mb", "92d0c527266ec915"},
	{algorithms.Murmur3_64, true, "./artefacts/test.1mb", "35ec8ac6041a7e9b"},
	{algorithms.Murmur3_32, false, "./artefacts/test-manipulated.1mb", "eb6482f3"},
	{algorithms.Murmur3_32, true, "./artefacts/test-manipulated.1mb", "e0fa6869"},
	{algorithms.Murmur3_32, false, "./artefacts/test.1mb", "5ca146ee"},
	{algorithms.Murmur3_32, true, "./artefacts/test.1mb", "3a3133fa"},
	{algorithms.Fnv128, false, "./artefacts/test-manipulated.1mb", "e91da5b6fb6c3df866d19794bcc031a2"},
	{algorithms.Fnv128, true, "./artefacts/test-manipulated.1mb", "8808e2a6d269deb5bce97f110f60e8dc"},
	{algorithms.Fnv128, false, "./artefacts/test.1mb", "af25513dbbfb8ebf847829a2cd6e76f2"},
	{algorithms.Fnv128, true, "./artefacts/test.1mb", "e55b683eca015645afc7316f7df9993b"},
	{algorithms.Fnv128a, false, "./artefacts/test-manipulated.1mb", "04721f877b7be5ad3e487b87ad486f30"},
	{algorithms.Fnv128a, true, "./artefacts/test-manipulated.1mb", "998f1046fb1e726b7dedd1eecd453c1a"},
	{algorithms.Fnv128a, false, "./artefacts/test.1mb", "f80ebc069329ec8a59e2c444c300f218"},
	{algorithms.Fnv128a, true, "./artefacts/test.1mb", "ebc231b45eb5b9c7be1c936829047f1f"},
	{algorithms.Md5, false, "./artefacts/test-manipulated.1mb", "040ca2ff5e59e6b0870b0f68a92a3968"},
	{algorithms.Md5, true, "./artefacts/test-manipulated.1mb", "df221ae4955e4b77f50ade6ab70c5210"},
	{algorithms.Md5, false, "./artefacts/test.1mb", "546b9508c9650e5d2e0c1c15f63c342c"},
	{algorithms.Md5, true, "./artefacts/test.1mb", "4c18efb7e70ac81f341ce3f5ef3684a4"},
	{algorithms.Sha512, false, "./artefacts/test-manipulated.1mb", "b8b783b66d20b280709522abd2478f0f7e599a31d62d9f876d8d91e7ad3874e75964f5bbb2e35ca1380e4d28d9135c40b12d3cee7c7b1f89c29b5d2ef38d0cc7"},
	{algorithms.Sha512, true, "./artefacts/test-manipulated.1mb", "dd69b1afbcb92135421574297fa47f612a23b386721b8562cd7852a0eebe0f4d8436d02b6773b7c072c18c67027d53eeedc9d18cc6171dfc82a907bfa570ae03"},
	{algorithms.Sha512, false, "./artefacts/test.1mb", "88402b9df2f2dd06597f0a1db9c6257645acb6ddb949d4daa00a7f28dfd681b5a46cef809774e9c0e5f0f581d8a240eac62bde89d99220055342dae8d6e680cf"},
	{algorithms.Sha512, true, "./artefacts/test.1mb", "8cedef8fa8d1ab8bdee1a9441165fe2af8ee37c9672e06f15ca30f5a3f840096585e474c2b800760bd66db96239f3c67761303ec1d87553f27afc7d8c9e7ea9f"},
	{algorithms.Sha256, false, "./artefacts/test-manipulated.1mb", "9539725bbdda1bfb410c51d9ebc0ba72391e7ba2145e74422028253a30672506"},
	{algorithms.Sha256, true, "./artefacts/test-manipulated.1mb", "aae139d218d16eb32cd63dc6f842f77c89a773fc26a8e7ef3b9023600fad3f17"},
	{algorithms.Sha256, false, "./artefacts/test.1mb", "11cfdec95e731151953ab8dbe24de8b3c1a029731740ca649bc82f95338e0540"},
	{algorithms.Sha256, true, "./artefacts/test.1mb", "e9403adc74d6a890a0db579ab217e2c4b0490b43e5a87552d3a239f1bdde91b8"},
}

func TestSlice_New_HashingAlgorithms_WithFileSystemFiles(t *testing.T) {

	options := SlicerOptions{
		DisableSlicing:       true,
		DisableMeta:          false,
		DisableFileDetection: false,
	}

	for _, item := range algoData {
		options.DisableSlicing = item.disableSlicing
		runHashCheckTestsForFileSystemFile(item.filename, item.algorithm, &options, item.expectHash, t)
	}
}

func runHashCheckTestsForFileSystemFile(filename string, algorithm algorithms.Algorithm, options *SlicerOptions, expected string, t *testing.T) {
	if binary, err := os.ReadFile(filename); err != nil {
		t.Errorf("Unexpected io error %v", err)
	} else {

		fsSize := len(binary)
		reader := bytes.NewReader(binary)
		sr := io.NewSectionReader(reader, 0, int64(fsSize))

		stats := SlicerStats{}

		slicer := New(algorithm)

		if err := slicer.Slice(sr, options, &stats); err != nil {
			t.Errorf("Unexpected Slicer error %v", err)
		}

		actual := hex.EncodeToString(stats.Hash)

		if len(expected) != len(actual) {
			t.Errorf("hash length expected %d, got %d", len(expected), len(actual))
		}

		if !strings.EqualFold(actual, expected) {
			t.Errorf("expected hash %s, got %s", expected, actual)
		}
	}
}
