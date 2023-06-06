// nnheap.go
// Christian Jordan
// Nearest neighbors' heap implementation

package config

// NeighborItem is a struct used for nearest neighbor heap
type NeighborItem struct {
	Neighbor *MileStone // Neighbor MileStone
	Dist     float32    //	Neighbor distance
	index    int        //	Index of item in heap
}

// NeighborHeap is a heap of NeighborItems
type NeighborHeap []*NeighborItem

// NewNeighborItem creates a new NeighborItem
func NewNeighborItem(neighbor *MileStone, dist float32) *NeighborItem {
	return &NeighborItem{
		Neighbor: neighbor,
		Dist:     dist,
	}
}

// Heap implementation
func (h NeighborHeap) Len() int { return len(h) }

func (h NeighborHeap) Less(i, j int) bool { return h[i].Dist < h[j].Dist }

// Swap swaps two items in the heap
func (h NeighborHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

// Push adds an item to the heap
func (h *NeighborHeap) Push(x interface{}) {
	item := x.(*NeighborItem)
	*h = append(*h, item)
}

// Pop removes the top item from the heap
func (h *NeighborHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	item.index = -1
	*h = old[0 : n-1]
	return item
}
