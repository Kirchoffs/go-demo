package heap

import "sort"

// Default min heap

// Len & Less & Swap & Push & Pop
type Interface interface {
    sort.Interface
    Push(x interface{})
    Pop() interface{}
}

func Init(h Interface) {
    n := h.Len()
    // A heap is a complete binary tree, so there are ceil(n / 2) leaf nodes.
    for i := n / 2 - 1; i >= 0; i-- {
        down(h, i, n)
    }
}

func Push(h Interface, x any) {
    h.Push(x)
    up(h, h.Len() - 1)
}

// Pop removes and returns the minimum element (according to Less) from the heap.
// We should always pop the elmement with the largest index from h.
func Pop(h Interface) any {
    n := h.Len()
    h.Swap(0, n - 1)
    down(h, 0, n - 1)
    return h.Pop()
}

func Remove(h Interface, i int) any {
    n := h.Len()
    if i != n - 1 {
        h.Swap(i, n - 1)
        if !down(h, i, n - 1) {
            up(h, i)
        }
    }
    return h.Pop()
}

func Fix(h Interface, i int) {
    if !down(h, i, h.Len()) {
        up(h, i)
    }
}

func up(h Interface, i int) {
    for {
        j := (i - 1) / 2
        if i == 0 || !h.Less(i, j) {
            break
        }
        h.Swap(i, j)
        i = j
    }
}

func down(h Interface, i int, n int) bool {
    p := i
    for {
        l := 2 * i + 1
        if l >= n || l < 0 {
            break
        }
        c := l
        if r := l + 1; r < n && h.Less(r, l) {
            c = r
        }
        if !h.Less(c, i) {
            break
        }
        h.Swap(i, c)
        i = c
    }
    return i > p
}
