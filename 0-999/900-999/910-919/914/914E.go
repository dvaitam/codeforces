package main

import (
	"bufio"
	"fmt"
	"os"
)

const ALPH = 20

var (
	n       int
	adj     [][]int
	label   []int
	parent  []int
	heavy   []int
	size    []int
	mask    []int
	freq    []int
	inc     []int64
	applied []int64
	diff    []int64
)

func dfsSize(v, p int) {
	parent[v] = p
	size[v] = 1
	heavy[v] = -1
	for _, to := range adj[v] {
		if to == p {
			continue
		}
		mask[to] = mask[v] ^ (1 << label[to])
		dfsSize(to, v)
		size[v] += size[to]
		if heavy[v] == -1 || size[to] > size[heavy[v]] {
			heavy[v] = to
		}
	}
}

func collect(v, p int, out *[]int) {
	*out = append(*out, v)
	for _, to := range adj[v] {
		if to == p {
			continue
		}
		collect(to, v, out)
	}
}

func addNode(x int) {
	m := mask[x]
	applied[x] = inc[m]
	freq[m]++
}

func removeNode(x int) {
	m := mask[x]
	freq[m]--
	diff[x] += inc[m] - applied[x]
}

func queryNode(x, lca int) {
	base := mask[x] ^ (1 << label[lca])
	pc := freq[base]
	diff[x] += int64(pc)
	inc[base]++
	for k := 0; k < ALPH; k++ {
		m := base ^ (1 << k)
		c := freq[m]
		diff[x] += int64(c)
		inc[m]++
		pc += c
	}
	diff[lca] -= int64(pc)
	if parent[lca] != 0 {
		diff[parent[lca]] -= int64(pc)
	}
}

func dfs(v, p int, keep bool) {
	for _, to := range adj[v] {
		if to == p || to == heavy[v] {
			continue
		}
		dfs(to, v, false)
	}
	if heavy[v] != -1 {
		dfs(heavy[v], v, true)
	}
	// process vertex v itself
	queryNode(v, v)
	addNode(v)
	for _, to := range adj[v] {
		if to == p || to == heavy[v] {
			continue
		}
		nodes := make([]int, 0, size[to])
		collect(to, v, &nodes)
		for _, x := range nodes {
			queryNode(x, v)
		}
		for _, x := range nodes {
			addNode(x)
		}
	}
	if !keep {
		nodes := make([]int, 0, size[v])
		collect(v, p, &nodes)
		for _, x := range nodes {
			removeNode(x)
		}
	}
}

func dfsAcc(u, p int) {
	for _, v := range adj[u] {
		if v == p {
			continue
		}
		dfsAcc(v, u)
		diff[u] += diff[v]
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscan(reader, &n)
	adj = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	var s string
	fmt.Fscan(reader, &s)
	label = make([]int, n+1)
	for i := 1; i <= n; i++ {
		label[i] = int(s[i-1] - 'a')
	}
	parent = make([]int, n+1)
	heavy = make([]int, n+1)
	size = make([]int, n+1)
	mask = make([]int, n+1)
	freq = make([]int, 1<<ALPH)
	inc = make([]int64, 1<<ALPH)
	applied = make([]int64, n+1)
	diff = make([]int64, n+1)

	mask[1] = 1 << label[1]
	dfsSize(1, 0)
	dfs(1, 0, false)
	dfsAcc(1, 0)

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for i := 1; i <= n; i++ {
		ans := diff[i] + 1
		if i > 1 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, ans)
	}
	fmt.Fprintln(writer)
}
