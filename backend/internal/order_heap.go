package internal

type OrderHeap []OrderItem

type OrderItem struct {
	Price     int32
	Volume    int32
	AccountId string
}

func (h OrderHeap) Len() int           { return len(h) }
func (h OrderHeap) Less(i, j int) bool { return h[i].Price < h[j].Price }
func (h OrderHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *OrderHeap) Push(x OrderItem) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x)
}

func (h *OrderHeap) Pop() OrderItem {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h *OrderHeap) Top() OrderItem {
	if len(*h) == 0 {
		return OrderItem{Price: 0, Volume: 0}
	}
	return (*h)[0]
}
