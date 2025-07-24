package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
	"sort"
)

type Pair struct {
	cnt int
	val int
}
type PairHeap []Pair

func (h PairHeap) Len() int { return len(h) }
func (h PairHeap) Less(i, j int) bool {
	if h[i].cnt == h[j].cnt {
		return h[i].val < h[j].val
	}
	return h[i].cnt > h[j].cnt
}
func (h PairHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *PairHeap) Push(x interface{}) { *h = append(*h, x.(Pair)) }
func (h *PairHeap) Pop() interface{} {
	old := *h
	v := old[len(old)-1]
	*h = old[:len(old)-1]
	return v
}

var arr []int
var freq []int
var pq PairHeap

func addPos(pos int) {
	x := arr[pos]
	freq[x]++
	heap.Push(&pq, Pair{freq[x], x})
}
func removePos(pos int) {
	x := arr[pos]
	freq[x]--
	heap.Push(&pq, Pair{freq[x], x})
}
func currentAns() int {
	for {
		p := pq[0]
		if freq[p.val] == p.cnt {
			return p.val
		}
		heap.Pop(&pq)
	}
}

type Query struct {
	l1, r1, l2, r2 int
	idx            int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	vals := make([]int, n+1)
	maxVal := 0
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &vals[i])
		if vals[i] > maxVal {
			maxVal = vals[i]
		}
	}
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		adj[x] = append(adj[x], y)
		adj[y] = append(adj[y], x)
	}
	// Euler tour
	tin := make([]int, n+1)
	tout := make([]int, n+1)
	order := make([]int, 0, n)
	stack := []struct{ node, idx int }{{1, 0}}
	parent := make([]int, n+1)
	time := 0
	for len(stack) > 0 {
		top := &stack[len(stack)-1]
		if top.idx == 0 {
			tin[top.node] = time
			order = append(order, top.node)
			time++
		}
		if top.idx < len(adj[top.node]) {
			v := adj[top.node][top.idx]
			top.idx++
			if v == parent[top.node] {
				continue
			}
			parent[v] = top.node
			stack = append(stack, struct{ node, idx int }{v, 0})
		} else {
			tout[top.node] = time - 1
			stack = stack[:len(stack)-1]
		}
	}
	arr = make([]int, n)
	for i, node := range order {
		arr[i] = vals[node]
	}

	var q int
	fmt.Fscan(in, &q)
	queries := make([]Query, q)
	for i := 0; i < q; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		l1, r1 := tin[u], tout[u]
		l2, r2 := tin[v], tout[v]
		if l1 > r1 {
			l1, r1 = 0, -1
		}
		if l2 > r2 {
			l2, r2 = 0, -1
		}
		queries[i] = Query{l1, r1, l2, r2, i}
	}

	block := int(math.Sqrt(float64(n))) + 1
	sort.Slice(queries, func(i, j int) bool {
		qi, qj := queries[i], queries[j]
		bi := qi.l1 / block
		bj := qj.l1 / block
		if bi != bj {
			return bi < bj
		}
		bi = qi.r1 / block
		bj = qj.r1 / block
		if bi != bj {
			return bi < bj
		}
		bi = qi.l2 / block
		bj = qj.l2 / block
		if bi != bj {
			return bi < bj
		}
		if qi.r2 != qj.r2 {
			return qi.r2 < qj.r2
		}
		return qi.idx < qj.idx
	})

	freq = make([]int, maxVal+2)
	pq = PairHeap{}
	heap.Init(&pq)

	l1, r1, l2, r2 := 0, -1, 0, -1
	ans := make([]int, q)
	for _, qu := range queries {
		for l1 > qu.l1 {
			l1--
			addPos(l1)
		}
		for r1 < qu.r1 {
			r1++
			addPos(r1)
		}
		for l1 < qu.l1 {
			removePos(l1)
			l1++
		}
		for r1 > qu.r1 {
			removePos(r1)
			r1--
		}

		for l2 > qu.l2 {
			l2--
			addPos(l2)
		}
		for r2 < qu.r2 {
			r2++
			addPos(r2)
		}
		for l2 < qu.l2 {
			removePos(l2)
			l2++
		}
		for r2 > qu.r2 {
			removePos(r2)
			r2--
		}

		ans[qu.idx] = currentAns()
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for i := 0; i < q; i++ {
		if i > 0 {
			fmt.Fprint(out, "\n")
		}
		fmt.Fprint(out, ans[i])
	}
}
