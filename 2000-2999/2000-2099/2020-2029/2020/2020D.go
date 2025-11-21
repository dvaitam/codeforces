package main

import (
	"bufio"
	"fmt"
	"os"
)

type DSU struct {
	parent []int
	size   []int
	comp   int
}

func NewDSU(n int) *DSU {
	parent := make([]int, n+1)
	size := make([]int, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = i
		size[i] = 1
	}
	return &DSU{parent: parent, size: size, comp: n}
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) union(a, b int) {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		return
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
	d.comp--
}

func getNext(edge []int, idx int) int {
	if idx >= len(edge) {
		return len(edge)
	}
	if edge[idx] == idx {
		return idx
	}
	edge[idx] = getNext(edge, edge[idx])
	return edge[idx]
}

func processRange(list []int, edges []int, idxL, idxR int, dsu *DSU) {
	if idxL >= idxR {
		return
	}
	i := getNext(edges, idxL)
	for i < idxR {
		dsu.union(list[i], list[i+1])
		if i < len(edges) {
			edges[i] = getNext(edges, i+1)
		}
		i = getNext(edges, i)
	}
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

		dsu := NewDSU(n)

		lists := make([][][]int, 11)
		pos := make([][]int, 11)
		for d := 1; d <= 10; d++ {
			lists[d] = make([][]int, d)
			pos[d] = make([]int, n+1)
			for r := 0; r < d; r++ {
				lists[d][r] = make([]int, 0)
			}
		}

		for node := 1; node <= n; node++ {
			for d := 1; d <= 10; d++ {
				r := node % d
				lst := &lists[d][r]
				pos[d][node] = len(*lst)
				*lst = append(*lst, node)
			}
		}

		edges := make([][][]int, 11)
		for d := 1; d <= 10; d++ {
			edges[d] = make([][]int, d)
			for r := 0; r < d; r++ {
				L := len(lists[d][r])
				if L >= 2 {
					arr := make([]int, L-1)
					for i := 0; i < L-1; i++ {
						arr[i] = i
					}
					edges[d][r] = arr
				} else {
					edges[d][r] = make([]int, 0)
				}
			}
		}

		for ; m > 0; m-- {
			var a, d, k int
			fmt.Fscan(in, &a, &d, &k)
			r := a % d
			idxL := pos[d][a]
			idxR := pos[d][a+k*d]
			processRange(lists[d][r], edges[d][r], idxL, idxR, dsu)
		}

		fmt.Fprintln(out, dsu.comp)
	}
}
