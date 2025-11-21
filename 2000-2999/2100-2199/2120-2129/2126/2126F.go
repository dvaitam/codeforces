package main

import (
	"bufio"
	"fmt"
	"os"
)

type edge struct {
	to int
	w  int64
}

type heavyEdge struct {
	idx int
	w   int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const threshold = 450

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, q int
		fmt.Fscan(in, &n, &q)

		color := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &color[i])
		}

		adj := make([][]edge, n)
		type rawEdge struct{ u, v int; w int64 }
		rawEdges := make([]rawEdge, 0, n-1)
		for i := 0; i < n-1; i++ {
			var u, v int
			var c int64
			fmt.Fscan(in, &u, &v, &c)
			u--
			v--
			adj[u] = append(adj[u], edge{to: v, w: c})
			adj[v] = append(adj[v], edge{to: u, w: c})
			rawEdges = append(rawEdges, rawEdge{u: u, v: v, w: c})
		}

		heavyID := make([]int, n)
		for i := range heavyID {
			heavyID[i] = -1
		}
		heavyNodes := make([]int, 0)
		for i := 0; i < n; i++ {
			if len(adj[i]) > threshold {
				heavyID[i] = len(heavyNodes)
				heavyNodes = append(heavyNodes, i)
			}
		}

		heavyMaps := make([]map[int]int64, len(heavyNodes))
		for i := range heavyMaps {
			heavyMaps[i] = make(map[int]int64)
		}

		heavyAdj := make([][]heavyEdge, n)
		for idx, h := range heavyNodes {
			for _, e := range adj[h] {
				col := color[e.to]
				heavyMaps[idx][col] += e.w
				heavyAdj[e.to] = append(heavyAdj[e.to], heavyEdge{idx: idx, w: e.w})
			}
		}

		var ans int64
		for _, e := range rawEdges {
			if color[e.u] != color[e.v] {
				ans += e.w
			}
		}

		computeSum := func(v int, c int) int64 {
			if heavyID[v] != -1 {
				return heavyMaps[heavyID[v]][c]
			}
			var sum int64
			for _, e := range adj[v] {
				if color[e.to] == c {
					sum += e.w
				}
			}
			return sum
		}

		for ; q > 0; q-- {
			var v, x int
			fmt.Fscan(in, &v, &x)
			v--
			if color[v] == x {
				fmt.Fprintln(out, ans)
				continue
			}
			oldColor := color[v]

			// sums of edges to neighbors with specified colors
			sumOld := computeSum(v, oldColor)
			sumNew := computeSum(v, x)

			ans += sumOld - sumNew

			color[v] = x

			for _, hEdge := range heavyAdj[v] {
				hMap := heavyMaps[hEdge.idx]
				hMap[oldColor] -= hEdge.w
				if hMap[oldColor] == 0 {
					delete(hMap, oldColor)
				}
				hMap[x] += hEdge.w
			}

			fmt.Fprintln(out, ans)
		}
	}
}

