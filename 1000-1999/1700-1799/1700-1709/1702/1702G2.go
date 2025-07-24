package main

import (
	"bufio"
	"fmt"
	"os"
)

const MAXN = 200005
const LOGN = 20

var (
	g     [MAXN][]int
	up    [LOGN][MAXN]int
	depth [MAXN]int
)

func buildLCA(n int) {
	parent := make([]int, n+1)
	stack := []int{1}
	parent[1] = 1
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		up[0][v] = parent[v]
		for i := 1; i < LOGN; i++ {
			up[i][v] = up[i-1][up[i-1][v]]
		}
		for _, to := range g[v] {
			if to == parent[v] {
				continue
			}
			parent[to] = v
			depth[to] = depth[v] + 1
			stack = append(stack, to)
		}
	}
	for j := 1; j < LOGN; j++ {
		for i := 1; i <= n; i++ {
			up[j][i] = up[j-1][up[j-1][i]]
		}
	}
}

func lca(a, b int) int {
	if depth[a] < depth[b] {
		a, b = b, a
	}
	diff := depth[a] - depth[b]
	for i := 0; i < LOGN; i++ {
		if diff>>i&1 == 1 {
			a = up[i][a]
		}
	}
	if a == b {
		return a
	}
	for i := LOGN - 1; i >= 0; i-- {
		if up[i][a] != up[i][b] {
			a = up[i][a]
			b = up[i][b]
		}
	}
	return up[0][a]
}

func dist(a, b int) int {
	l := lca(a, b)
	return depth[a] + depth[b] - 2*depth[l]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	buildLCA(n)

	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var k int
		fmt.Fscan(reader, &k)
		arr := make([]int, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		a := arr[0]
		b := a
		maxd := -1
		for _, v := range arr {
			d := dist(a, v)
			if d > maxd {
				maxd = d
				b = v
			}
		}
		c := b
		maxd = -1
		for _, v := range arr {
			d := dist(b, v)
			if d > maxd {
				maxd = d
				c = v
			}
		}
		dia := maxd
		pass := true
		for _, v := range arr {
			if dist(b, v)+dist(v, c) != dia {
				pass = false
				break
			}
		}
		if pass {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
