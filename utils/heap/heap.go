package heap

// Kv 代表小堆中的元素，其中K相当于redis zset中的值，
// V相当于值所对应的score
type KV struct {
	K interface{}
	V float64
}

// 小根堆
type MinHeap []KV

func (h MinHeap) Init() {
	l := len(h)

	for i := (l - 1) / 2; i >= 0; i-- {
		h.down(i, l)
	}
}

func (h MinHeap) down(s, l int) {
	r := s
	for {
		c := 2*r + 1
		if c >= l {
			return
		}
		if c+1 < l && h[c+1].V < h[c].V {
			c = c + 1
		}
		if h[r].V < h[c].V {
			return
		}
		h[r], h[c] = h[c], h[r]
		r = c
	}
}

func (h MinHeap) Push(kv KV) bool {
	if h[0].V >= kv.V {
		return false
	}
	h[0] = kv
	h.down(0, len(h))
	return true
}

func (h MinHeap) Len() int {
	return len(h)
}

func (h MinHeap) Less(i, j int) bool {
	return h[i].V < h[j].V
}

func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
