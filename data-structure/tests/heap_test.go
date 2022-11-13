package tests

import (
    "testing"
    "data-structure/heap"
)

type IntHeap []int

func (h IntHeap) Len() int {
    return len(h)
}

func (h IntHeap) Less(i, j int) bool {
    return h[i] < h[j]
}

func (h IntHeap) Swap(i, j int) { 
    h[i], h[j] = h[j], h[i] 
}

func (h *IntHeap) Push(x interface{}) {
    *h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}

func TestHeap(t *testing.T) {
    h := &IntHeap{2, 1, 5}
    heap.Init(h)

    heap.Push(h, 3)
    if (*h)[0] != 1 {
        t.Errorf("heap: %v", *h)
    }

    heap.Remove(h, 2)
    if (*h)[2] != 3 {
        t.Errorf("heap: %v", *h)
    }

    heap.Push(h, 5)
    heap.Push(h, 1)
    if (*h)[1] != 1 {
        t.Errorf("heap: %v", *h)
    }

    (*h)[1] = 4
    heap.Fix(h, 1)
    if (*h)[1] != 2 {
        t.Errorf("heap: %v", *h)
    }

    if heap.Pop(h) != 1 {
        t.Errorf("heap: %v", *h)
    }
    if heap.Pop(h) != 2 {
        t.Errorf("heap: %v", *h)
    }
    if heap.Pop(h) != 3 {
        t.Errorf("heap: %v", *h)
    }
}
