package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type node struct {
	sum int64
	i   int
	j   int
}

type maxHeap []node

func (h maxHeap) Len() int           { return len(h) }
func (h maxHeap) Less(i, j int) bool { return h[i].sum > h[j].sum }
func (h maxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) {
	*h = append(*h, x.(node))
}
func (h *maxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func pairKey(a, b int64) int64 {
	if a > b {
		a, b = b, a
	}
	return (a << 32) | b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		freq := make(map[int64]int)
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(in, &x)
			freq[x]++
		}
		groups := make(map[int][]int64)
		for v, c := range freq {
			groups[c] = append(groups[c], v)
		}
		counts := make([]int, 0, len(groups))
		for c := range groups {
			counts = append(counts, c)
			sort.Slice(groups[c], func(i, j int) bool { return groups[c][i] > groups[c][j] })
		}
		sort.Ints(counts)

		bad := make(map[int64]struct{})
		for i := 0; i < m; i++ {
			var x, y int64
			fmt.Fscan(in, &x, &y)
			bad[pairKey(x, y)] = struct{}{}
		}

		var ans int64
		for i := 0; i < len(counts); i++ {
			c1 := counts[i]
			arr1 := groups[c1]
			for j := i; j < len(counts); j++ {
				c2 := counts[j]
				arr2 := groups[c2]
				if c1 == c2 && len(arr1) < 2 {
					continue
				}
				h := &maxHeap{}
				visited := make(map[int64]struct{})
				if c1 == c2 {
					heap.Push(h, node{int64(arr1[0] + arr2[1]), 0, 1})
					visited[int64(0)<<32|1] = struct{}{}
				} else {
					heap.Push(h, node{int64(arr1[0] + arr2[0]), 0, 0})
					visited[int64(0)<<32|0] = struct{}{}
				}
				found := false
				for h.Len() > 0 && !found {
					cur := heap.Pop(h).(node)
					x := arr1[cur.i]
					y := arr2[cur.j]
					if x != y {
						if _, ok := bad[pairKey(x, y)]; !ok {
							val := int64(c1+c2) * (x + y)
							if val > ans {
								ans = val
							}
							found = true
							break
						}
					}
					ni, nj := cur.i+1, cur.j
					if ni < len(arr1) {
						key := int64(ni)<<32 | int64(nj)
						if _, ok := visited[key]; !ok {
							visited[key] = struct{}{}
							heap.Push(h, node{int64(arr1[ni] + arr2[nj]), ni, nj})
						}
					}
					ni, nj = cur.i, cur.j+1
					if nj < len(arr2) {
						key := int64(ni)<<32 | int64(nj)
						if _, ok := visited[key]; !ok {
							visited[key] = struct{}{}
							heap.Push(h, node{int64(arr1[ni] + arr2[nj]), ni, nj})
						}
					}
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}
