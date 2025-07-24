package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	adj := make([][]int, n)
	deg := make([]int, n)
	for i := 2; i <= n; i++ {
		var p int
		fmt.Fscan(in, &p)
		p--
		adj[i-1] = append(adj[i-1], p)
		adj[p] = append(adj[p], i-1)
		deg[i-1]++
		deg[p]++
	}

	leaves := []int{}
	for i, d := range deg {
		if d == 1 {
			leaves = append(leaves, i)
		}
	}

	leafDist := make([]int, n)
	for i := range leafDist {
		leafDist[i] = -1
	}
	q := make([]int, len(leaves))
	copy(q, leaves)
	for _, v := range leaves {
		leafDist[v] = 0
	}
	head := 0
	for head < len(q) {
		u := q[head]
		head++
		for _, v := range adj[u] {
			if leafDist[v] == -1 {
				leafDist[v] = leafDist[u] + 1
				q = append(q, v)
			}
		}
	}

	bfs := func(start int) []int {
		dist := make([]int, n)
		for i := range dist {
			dist[i] = -1
		}
		q := []int{start}
		dist[start] = 0
		for h := 0; h < len(q); h++ {
			u := q[h]
			for _, v := range adj[u] {
				if dist[v] == -1 {
					dist[v] = dist[u] + 1
					q = append(q, v)
				}
			}
		}
		return dist
	}

	var qnum int
	fmt.Fscan(in, &qnum)
	xs := make([]int, qnum)
	for i := range xs {
		fmt.Fscan(in, &xs[i])
	}

	res := make([]int, qnum)
	for idx, x := range xs {
		dist0 := bfs(0)
		minLeaf := int(1e9)
		for _, l := range leaves {
			if dist0[l] < minLeaf {
				minLeaf = dist0[l]
			}
		}
		farthest := 0
		maxVal := -1
		for i := 0; i < n; i++ {
			val := dist0[i]
			alt := minLeaf + x + leafDist[i]
			if alt < val {
				val = alt
			}
			if val > maxVal {
				maxVal = val
				farthest = i
			}
		}
		dist1 := bfs(farthest)
		minLeaf = int(1e9)
		for _, l := range leaves {
			if dist1[l] < minLeaf {
				minLeaf = dist1[l]
			}
		}
		maxVal = 0
		for i := 0; i < n; i++ {
			val := dist1[i]
			alt := minLeaf + x + leafDist[i]
			if alt < val {
				val = alt
			}
			if val > maxVal {
				maxVal = val
			}
		}
		res[idx] = maxVal
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for i, v := range res {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	fmt.Fprintln(out)
}
