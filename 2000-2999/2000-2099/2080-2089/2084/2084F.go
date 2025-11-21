package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type Fenwick struct {
	tree []int
	n    int
}

func newFenwick(n int) *Fenwick {
	return &Fenwick{tree: make([]int, n+2), n: n}
}

func (f *Fenwick) Reset() {
	for i := range f.tree {
		f.tree[i] = 0
	}
}

func (f *Fenwick) Add(idx, delta int) {
	for idx <= f.n {
		f.tree[idx] += delta
		idx += idx & -idx
	}
}

func (f *Fenwick) Sum(idx int) int {
	res := 0
	for idx > 0 {
		res += f.tree[idx]
		idx -= idx & -idx
	}
	return res
}

type Item struct {
	r   int
	val int
}

type MinHeap []Item

func (h MinHeap) Len() int { return len(h) }
func (h MinHeap) Less(i, j int) bool {
	return h[i].r < h[j].r || (h[i].r == h[j].r && h[i].val < h[j].val)
}
func (h MinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.(Item)) }
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

func solveOne(n int, a, c []int, out *bufio.Writer) {
	pos := make([]int, n+1)
	for i := 0; i < n; i++ {
		pos[a[i]] = i + 1
	}

	L := make([]int, n+1)
	R := make([]int, n+1)

	fen := newFenwick(n)
	for i := 1; i <= n; i++ {
		val := a[i-1]
		L[val] = 1 + fen.Sum(val-1)
		fen.Add(val, 1)
	}

	fen = newFenwick(n)
	for i := n; i >= 1; i-- {
		val := a[i-1]
		greater := fen.Sum(n) - fen.Sum(val)
		R[val] = n - greater
		fen.Add(val, 1)
	}

	fixedPos := make([]int, n+1)
	for i := 0; i < n; i++ {
		if c[i] != 0 {
			val := c[i]
			if fixedPos[val] != 0 {
				fmt.Fprintln(out, -1)
				return
			}
			fixedPos[val] = i + 1
		}
	}

	for val := 1; val <= n; val++ {
		if fixedPos[val] != 0 {
			posn := fixedPos[val]
			if posn < L[val] || posn > R[val] {
				fmt.Fprintln(out, -1)
				return
		}
	}
}

	addAt := make([][]int, n+2)
	for val := 1; val <= n; val++ {
		if fixedPos[val] == 0 {
			addAt[L[val]] = append(addAt[L[val]], val)
		}
	}

	assigned := make([]bool, n+1)
	res := make([]int, n)
	pq := MinHeap{}
	for posIdx := 1; posIdx <= n; posIdx++ {
		for _, v := range addAt[posIdx] {
			heap.Push(&pq, Item{r: R[v], val: v})
		}

		for pq.Len() > 0 && pq[0].r < posIdx {
			fmt.Fprintln(out, -1)
			return
		}

		if c[posIdx-1] != 0 {
			val := c[posIdx-1]
			if assigned[val] {
				fmt.Fprintln(out, -1)
				return
			}
			assigned[val] = true
			res[posIdx-1] = val
			continue
		}

		if pq.Len() == 0 {
			fmt.Fprintln(out, -1)
			return
		}
		item := heap.Pop(&pq).(Item)
		val := item.val
		if assigned[val] {
			fmt.Fprintln(out, -1)
			return
		}
		assigned[val] = true
		res[posIdx-1] = val
	}

	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, res[i])
	}
	fmt.Fprintln(out)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		c := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &c[i])
		}
		solveOne(n, a, c, out)
	}
}
