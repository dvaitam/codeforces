package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	n      int
	colors []int
	adj    [][]int
	st, en []int
	flat   []int
	sz     []int
	heavy  []int
	timer  int

	freqColor []int
	cntFreq   []int
	sumFreq   []int64
	maxFreq   int
	ans       []int64
)

func dfs1(u, p int) {
	timer++
	st[u] = timer
	flat[timer] = u
	sz[u] = 1
	heavy[u] = -1
	maxSize := 0
	for _, v := range adj[u] {
		if v == p {
			continue
		}
		dfs1(v, u)
		if sz[v] > maxSize {
			maxSize = sz[v]
			heavy[u] = v
		}
		sz[u] += sz[v]
	}
	en[u] = timer
}

func addNode(u int) {
	c := colors[u]
	old := freqColor[c]
	if old > 0 {
		cntFreq[old]--
		sumFreq[old] -= int64(c)
	}
	newf := old + 1
	freqColor[c] = newf
	cntFreq[newf]++
	sumFreq[newf] += int64(c)
	if newf > maxFreq {
		maxFreq = newf
	}
}

func removeNode(u int) {
	c := colors[u]
	f := freqColor[c]
	if f == 0 {
		return
	}
	cntFreq[f]--
	sumFreq[f] -= int64(c)
	freqColor[c] = 0
	if cntFreq[maxFreq] == 0 {
		for maxFreq > 0 && cntFreq[maxFreq] == 0 {
			maxFreq--
		}
	}
}

func dfs2(u, p int, keep bool) {
	for _, v := range adj[u] {
		if v == p || v == heavy[u] {
			continue
		}
		dfs2(v, u, false)
	}
	if heavy[u] != -1 {
		dfs2(heavy[u], u, true)
	}
	for _, v := range adj[u] {
		if v == p || v == heavy[u] {
			continue
		}
		for i := st[v]; i <= en[v]; i++ {
			addNode(flat[i])
		}
	}
	addNode(u)
	ans[u] = sumFreq[maxFreq]
	if !keep {
		for i := st[u]; i <= en[u]; i++ {
			removeNode(flat[i])
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscan(reader, &n)
	colors = make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &colors[i])
	}
	adj = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	st = make([]int, n+1)
	en = make([]int, n+1)
	flat = make([]int, n+1)
	sz = make([]int, n+1)
	heavy = make([]int, n+1)
	dfs1(1, 0)

	freqColor = make([]int, n+1)
	cntFreq = make([]int, n+1)
	sumFreq = make([]int64, n+1)
	ans = make([]int64, n+1)

	dfs2(1, 0, true)

	writer := bufio.NewWriter(os.Stdout)
	for i := 1; i <= n; i++ {
		if i > 1 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, ans[i])
	}
	writer.WriteByte('\n')
	writer.Flush()
}
