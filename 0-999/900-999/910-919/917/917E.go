package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to   int
	char byte
}

var (
	g      [][]Edge
	parent []int
	depth  []int
	letter []byte
	up     [][]int
	maxLog int
)

func dfs(root int) {
	stack := []int{root}
	parent[root] = 0
	depth[root] = 0
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, e := range g[u] {
			if e.to == parent[u] {
				continue
			}
			parent[e.to] = u
			letter[e.to] = e.char
			depth[e.to] = depth[u] + 1
			stack = append(stack, e.to)
		}
	}
}

func buildLCA(n int) {
	maxLog = 0
	for (1 << uint(maxLog)) <= n {
		maxLog++
	}
	up = make([][]int, maxLog)
	up[0] = make([]int, n+1)
	copy(up[0], parent)
	for k := 1; k < maxLog; k++ {
		up[k] = make([]int, n+1)
		for v := 1; v <= n; v++ {
			if up[k-1][v] > 0 {
				up[k][v] = up[k-1][up[k-1][v]]
			}
		}
	}
}

func lca(a, b int) int {
	if depth[a] < depth[b] {
		a, b = b, a
	}
	diff := depth[a] - depth[b]
	for k := 0; diff > 0; k++ {
		if diff&1 == 1 {
			a = up[k][a]
		}
		diff >>= 1
	}
	if a == b {
		return a
	}
	for k := maxLog - 1; k >= 0; k-- {
		if up[k][a] != up[k][b] {
			a = up[k][a]
			b = up[k][b]
		}
	}
	return up[0][a]
}

func pathString(u, v int) string {
	l := lca(u, v)
	buf := make([]byte, 0, depth[u]+depth[v]-2*depth[l])
	for x := u; x != l; x = parent[x] {
		buf = append(buf, letter[x])
	}
	tmp := make([]byte, 0)
	for x := v; x != l; x = parent[x] {
		tmp = append(tmp, letter[x])
	}
	for i := len(tmp) - 1; i >= 0; i-- {
		buf = append(buf, tmp[i])
	}
	return string(buf)
}

func kmpTable(pat string) []int {
	n := len(pat)
	table := make([]int, n)
	j := 0
	for i := 1; i < n; i++ {
		for j > 0 && pat[i] != pat[j] {
			j = table[j-1]
		}
		if pat[i] == pat[j] {
			j++
			table[i] = j
		} else {
			table[i] = j
		}
	}
	return table
}

func countOccurrences(text, pat string) int {
	if len(pat) == 0 {
		return 0
	}
	table := kmpTable(pat)
	j := 0
	cnt := 0
	for i := 0; i < len(text); i++ {
		for j > 0 && text[i] != pat[j] {
			j = table[j-1]
		}
		if text[i] == pat[j] {
			j++
			if j == len(pat) {
				cnt++
				j = table[j-1]
			}
		}
	}
	return cnt
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, q int
	fmt.Fscan(in, &n, &m, &q)

	g = make([][]Edge, n+1)
	parent = make([]int, n+1)
	depth = make([]int, n+1)
	letter = make([]byte, n+1)

	for i := 0; i < n-1; i++ {
		var v, u int
		var c string
		fmt.Fscan(in, &v, &u, &c)
		ch := c[0]
		g[v] = append(g[v], Edge{u, ch})
		g[u] = append(g[u], Edge{v, ch})
	}

	dfs(1)
	buildLCA(n)

	words := make([]string, m+1)
	for i := 1; i <= m; i++ {
		fmt.Fscan(in, &words[i])
	}

	for i := 0; i < q; i++ {
		var a, b, k int
		fmt.Fscan(in, &a, &b, &k)
		txt := pathString(a, b)
		pat := words[k]
		fmt.Fprintln(out, countOccurrences(txt, pat))
	}
}
