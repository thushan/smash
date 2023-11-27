package report

import (
	"container/heap"
)

type TopFilesSummary struct {
	fileHeap *FileHeap
	size     int
}

type FileHeap []SmashFile

func (fh FileHeap) Len() int           { return len(fh) }
func (fh FileHeap) Less(i, j int) bool { return fh[i].FileSize < fh[j].FileSize }
func (fh FileHeap) Swap(i, j int)      { fh[i], fh[j] = fh[j], fh[i] }

func (fh *FileHeap) Push(x interface{}) {
	if f, ok := x.(SmashFile); ok {
		*fh = append(*fh, f)
	}
}

func (fh *FileHeap) Pop() interface{} {
	old := *fh
	n := len(old)
	file := old[n-1]
	*fh = old[0 : n-1]
	return file
}

func NewTopFilesSummary(size int) *TopFilesSummary {
	fileHeap := &FileHeap{}
	heap.Init(fileHeap)

	return &TopFilesSummary{fileHeap: fileHeap, size: size}
}

func (t *TopFilesSummary) Add(file SmashFile) {
	if t.fileHeap.Len() < t.size {
		heap.Push(t.fileHeap, file)
	} else if file.FileSize > (*t.fileHeap)[0].FileSize {
		heap.Pop(t.fileHeap)
		heap.Push(t.fileHeap, file)
	}
}

func (t *TopFilesSummary) Next() (SmashFile, bool) {
	if t.fileHeap.Len() != 0 {
		if f, ok := heap.Pop(t.fileHeap).(SmashFile); ok {
			return f, true
		}
	}
	return SmashFile{}, false
}
func (t *TopFilesSummary) All() []SmashFile {
	if t.fileHeap.Len() == 0 {
		return []SmashFile{}
	}
	files := make([]SmashFile, t.fileHeap.Len())
	copy(files, *t.fileHeap)
	return files
}
