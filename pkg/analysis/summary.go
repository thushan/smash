package analysis

import (
	"container/heap"
	"sort"
)

type Item struct {
	Key  string
	Size uint64
}
type Summary struct {
	itemHeap *ItemHeap
	maxSize  int
}

type ItemHeap []Item

func (ih ItemHeap) Len() int           { return len(ih) }
func (ih ItemHeap) Less(i, j int) bool { return ih[i].Size < ih[j].Size }
func (ih ItemHeap) Swap(i, j int)      { ih[i], ih[j] = ih[j], ih[i] }

func (ih *ItemHeap) Push(x interface{}) {
	if i, ok := x.(Item); ok {
		*ih = append(*ih, i)
	}
}

func (ih *ItemHeap) Pop() interface{} {
	old := *ih
	n := len(old)
	item := old[n-1]
	*ih = old[0 : n-1]
	return item
}

func NewSummary(size int) *Summary {
	fileHeap := &ItemHeap{}
	heap.Init(fileHeap)

	return &Summary{itemHeap: fileHeap, maxSize: size}
}

func (t *Summary) Add(item Item) {
	if t.itemHeap.Len() < t.maxSize {
		heap.Push(t.itemHeap, item)
	} else if item.Size > (*t.itemHeap)[0].Size {
		heap.Pop(t.itemHeap)
		heap.Push(t.itemHeap, item)
	}
}

func (t *Summary) Next() (Item, bool) {
	if t.itemHeap.Len() != 0 {
		if i, ok := heap.Pop(t.itemHeap).(Item); ok {
			return i, true
		}
	}
	return Item{}, false
}
func (t *Summary) All() []Item {
	if t.itemHeap.Len() == 0 {
		return []Item{}
	}
	items := make([]Item, t.itemHeap.Len())
	copy(items, *t.itemHeap)
	sort.Sort(ItemHeap(items))
	return items
}
