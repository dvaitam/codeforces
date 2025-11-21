package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type neighbor struct {
	to int
	id int
}

type triple struct {
	a, b, c int
}

type minHeap []int

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}
func (h *minHeap) Pop() interface{} {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

const mask32 uint64 = (1 << 32) - 1

func edgeKey(u, v int) uint64 {
	if u > v {
		u, v = v, u
	}
	return (uint64(u) << 32) | uint64(v)
}

func buildTargetTree(n int, parity []int) map[uint64]struct{} {
	target := make(map[uint64]struct{})
	oddCnt := 0
	for i := 1; i <= n; i++ {
		if parity[i] == 1 {
			oddCnt++
		}
	}
	if oddCnt == 0 {
		return target
	}

	deg := make([]int, n+1)
	for i := 1; i <= n; i++ {
		if parity[i] == 1 {
			deg[i] = 1
		} else {
			deg[i] = 2
		}
	}
	extra := oddCnt - 2
	idx := 1
	for extra > 0 {
		deg[idx] += 2
		extra -= 2
		idx++
		if idx > n {
			idx = 1
		}
	}

	if n == 1 {
		return target
	}

	prufer := make([]int, 0, n-2)
	for i := 1; i <= n; i++ {
		for c := deg[i] - 1; c > 0; c-- {
			prufer = append(prufer, i)
		}
	}

	h := &minHeap{}
	for i := 1; i <= n; i++ {
		if deg[i] == 1 {
			heap.Push(h, i)
		}
	}

	for _, v := range prufer {
		if h.Len() == 0 {
			break
		}
		leaf := heap.Pop(h).(int)
		target[edgeKey(leaf, v)] = struct{}{}
		deg[leaf]--
		deg[v]--
		if deg[v] == 1 {
			heap.Push(h, v)
		}
	}

	if h.Len() >= 2 {
		u := heap.Pop(h).(int)
		v := heap.Pop(h).(int)
		target[edgeKey(u, v)] = struct{}{}
	}

	return target
}

func appendCycleOps(cycle []int, ops *[]triple) {
	k := len(cycle)
	if k < 3 {
		return
	}
	for i := 1; i < k-1; i++ {
		*ops = append(*ops, triple{cycle[0], cycle[i], cycle[i+1]})
	}
}

func processPath(path []int, ops *[]triple) {
	if len(path) <= 1 {
		return
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	stack := []int{}
	pos := make(map[int]int)

	for _, v := range path {
		if idx, ok := pos[v]; !ok {
			pos[v] = len(stack)
			stack = append(stack, v)
			continue
		} else {
			cycle := append([]int(nil), stack[idx:]...)
			cycle = append(cycle, v)
			for j := idx; j < len(stack); j++ {
				delete(pos, stack[j])
			}
			stack = stack[:idx]
			pos[v] = len(stack)
			stack = append(stack, v)
			if len(cycle) >= 2 {
				unique := cycle[:len(cycle)-1]
				if len(unique) >= 3 {
					appendCycleOps(unique, ops)
				}
			}
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)

		parity := make([]int, n+1)
		initial := make(map[uint64]struct{}, m)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			if u == v {
				continue
			}
			key := edgeKey(u, v)
			initial[key] = struct{}{}
			parity[u] ^= 1
			parity[v] ^= 1
		}

		target := buildTargetTree(n, parity)

		diff := make(map[uint64]bool, len(initial)+len(target))
		for key := range initial {
			diff[key] = true
		}
		for key := range target {
			if diff[key] {
				delete(diff, key)
			} else {
				diff[key] = true
			}
		}

		edgeCount := len(diff)
		eu := make([]int, 0, edgeCount)
		ev := make([]int, 0, edgeCount)
		adj := make([][]neighbor, n+1)
		for key := range diff {
			u := int(key >> 32)
			v := int(key & mask32)
			id := len(eu)
			eu = append(eu, u)
			ev = append(ev, v)
			adj[u] = append(adj[u], neighbor{to: v, id: id})
			adj[v] = append(adj[v], neighbor{to: u, id: id})
		}

		used := make([]bool, len(eu))
		ptr := make([]int, n+1)
		degRem := make([]int, n+1)
		for i := 1; i <= n; i++ {
			degRem[i] = len(adj[i])
		}

		operations := make([]triple, 0)

		for start := 1; start <= n; start++ {
			if degRem[start] == 0 {
				continue
			}
			stack := []int{start}
			path := []int{}
			for len(stack) > 0 {
				u := stack[len(stack)-1]
				for ptr[u] < len(adj[u]) && used[adj[u][ptr[u]].id] {
					ptr[u]++
				}
				if ptr[u] == len(adj[u]) {
					path = append(path, u)
					stack = stack[:len(stack)-1]
				} else {
					e := adj[u][ptr[u]]
					ptr[u]++
					if used[e.id] {
						continue
					}
					used[e.id] = true
					degRem[eu[e.id]]--
					degRem[ev[e.id]]--
					stack = append(stack, e.to)
				}
			}
			processPath(path, &operations)
		}

		fmt.Fprintln(out, len(operations))
		for _, op := range operations {
			fmt.Fprintf(out, "%d %d %d\n", op.a, op.b, op.c)
		}
	}
}
