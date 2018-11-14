package heap

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func TestMinBasic(t *testing.T) {
	h := make(MinHeap, 0, 10)
	for i := 0; i < 10; i++ {
		h = append(h, KV{
			K: i,
			V: float64(i),
		})
	}

	h.Init()

	if !reflect.DeepEqual(h[0].K, 0) {
		t.Fatalf("min heap first should be 0, but is %d", h[0].K)
	}

	h.Push(KV{
		K: -1,
		V: float64(-1),
	})
	if !reflect.DeepEqual(h[0].K, 0) {
		t.Fatalf("min heap first should be 0, but is %d", h[0].K)
	}

	h.Push(KV{
		K: 11,
		V: float64(11),
	})
	if !reflect.DeepEqual(h[0].K, 1) {
		t.Fatalf("min heap first should be 1, but is %d", h[0].K)
	}

	h.Push(KV{
		K: 12,
		V: float64(12),
	})
	h.Push(KV{
		K: 13,
		V: float64(13),
	})
	if !reflect.DeepEqual(h[0].K, 3) {
		t.Fatalf("min heap first should be 3, but is %d", h[0].K)
	}
	if len(h) != 10 {
		t.Fatalf("min heap len should be 10, but is %d", len(h))
	}
}

func BenchmarkMinHeap(b *testing.B) {
	b.StopTimer()
	h := make(MinHeap, 0, 1000)
	for i := 0; i < 1000; i++ {
		v := rand.Int63n(5000)
		h = append(h, KV{
			K: v,
			V: float64(v),
		})
	}
	b.StartTimer()

	h.Init()
	for i := 0; i < b.N; i++ {
		h.Push(KV{
			K: i,
			V: float64(i),
		})
	}
}

func TestSort(t *testing.T) {
	h := make(MinHeap, 10)
	for i := 0; i < 10; i++ {
		h[i] = KV{
			K: i,
			V: float64(i),
		}
	}
	sort.Sort(h)

	for i := 0; i < 10; i++ {
		k, ok := h[i].K.(int)
		if !ok {
			t.Fatalf("k should be int type")
		}
		if k != i {
			t.Fatalf("k should be %d but is %d", i, k)
		}
	}
}
