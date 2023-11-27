package report

import "container/heap"

type TopFiles []SmashFile

var fileHeap = &TopFiles{}

func (fh TopFiles) Len() int           { return len(fh) }
func (fh TopFiles) Less(i, j int) bool { return fh[i].FileSize < fh[j].FileSize }
func (fh TopFiles) Swap(i, j int)      { fh[i], fh[j] = fh[j], fh[i] }

func (fh *TopFiles) Push(x interface{}) {
	*fh = append(*fh, x.(SmashFile))
}
func init() {
	heap.Init(fileHeap)
}
func (fh *TopFiles) Pop() interface{} {
	old := *fh
	n := len(old)
	file := old[n-1]
	*fh = old[0 : n-1]
	return file
}

func (fh *TopFiles) Index(file SmashFile) {
	if fileHeap.Len() < 5 {
		heap.Push(fileHeap, file)
	} else if file.FileSize > (*fileHeap)[0].FileSize {
		heap.Pop(fileHeap)
		heap.Push(fileHeap, file)
	}
}

func (fh *TopFiles) Iterator() func() (SmashFile, bool) {

	index := 0

	return func() (SmashFile, bool) {
		if fileHeap.Len() == 0 || index >= fileHeap.Len() {
			return SmashFile{}, false
		}
		file := (*fileHeap)[index]
		index++
		return file, true
	}
}
