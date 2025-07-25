package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

var (
	n, m               int
	g                  [][]int
	fa, deep, son, top []int
	siz                []int
	sum, src           []float64
)

func dfs1(x int) {
	siz[x] = 1
	deep[x] = deep[fa[x]] + 1
	for _, to := range g[x] {
		if to != fa[x] {
			fa[to] = x
			dfs1(to)
			siz[x] += siz[to]
			sum[x] += sum[to] + float64(siz[to])
			if siz[son[x]] < siz[to] {
				son[x] = to
			}
		}
	}
}

func dfs2(x int) {
	if x == son[fa[x]] {
		top[x] = top[fa[x]]
	} else {
		top[x] = x
	}
	for _, to := range g[x] {
		if fa[to] == x {
			dfs2(to)
		}
	}
}

func dfs3(x int) {
	src[x] += sum[x]
	for _, to := range g[x] {
		if fa[to] == x {
			src[to] += src[x] - sum[to] - float64(siz[to]) + float64(n-siz[to])
			dfs3(to)
		}
	}
}

func lca(u, v int) int {
	for top[u] != top[v] {
		if deep[top[u]] > deep[top[v]] {
			u = fa[top[u]]
		} else {
			v = fa[top[v]]
		}
	}
	if deep[u] < deep[v] {
		return u
	}
	return v
}

func solveNode(u, f int) int {
	for {
		if deep[fa[top[u]]] > deep[f] {
			u = fa[top[u]]
		} else if fa[u] == f {
			return u
		} else if deep[top[u]] > deep[f] {
			u = top[u]
		} else {
			return son[f]
		}
	}
}

func ask(u, v int) float64 {
	l := lca(u, v)
	if u != l && v != l {
		return sum[u]/float64(siz[u]) + sum[v]/float64(siz[v]) + float64(deep[u]+deep[v]+1-2*deep[l])
	}
	if u == l {
		u, v = v, u
	}
	nw := solveNode(u, v)
	x := sum[u] / float64(siz[u])
	y := (src[v] - sum[nw] - float64(siz[nw])) / float64(n-siz[nw])
	return x + y + float64(deep[u]-deep[v]+1)
}

func solveE(nv int, edges [][2]int, queries [][2]int) []float64 {
	n = nv
	g = make([][]int, n+1)
	for _, e := range edges {
		u := e[0]
		v := e[1]
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	fa = make([]int, n+1)
	deep = make([]int, n+1)
	son = make([]int, n+1)
	top = make([]int, n+1)
	siz = make([]int, n+1)
	sum = make([]float64, n+1)
	src = make([]float64, n+1)
	fa[1] = 0
	deep[0] = 0
	dfs1(1)
	dfs2(1)
	dfs3(1)
	res := make([]float64, len(queries))
	for i, q := range queries {
		res[i] = ask(q[0], q[1])
	}
	return res
}

func genTree(n int) [][2]int {
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		p := rand.Intn(i-1) + 1
		edges[i-2] = [2]int{p, i}
	}
	return edges
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 0; t < 100; t++ {
		n := rand.Intn(10) + 2
		m := rand.Intn(10) + 1
		edges := genTree(n)
		queries := make([][2]int, m)
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d\n", n, m)
		for _, e := range edges {
			fmt.Fprintf(&input, "%d %d\n", e[0], e[1])
		}
		for i := 0; i < m; i++ {
			u := rand.Intn(n) + 1
			v := rand.Intn(n) + 1
			for v == u {
				v = rand.Intn(n) + 1
			}
			queries[i] = [2]int{u, v}
			fmt.Fprintf(&input, "%d %d\n", u, v)
		}
		expected := solveE(n, edges, queries)
		cmd := exec.Command(binary)
		cmd.Stdin = &input
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: binary error: %v\n", t, err)
			os.Exit(1)
		}
		scanner := bufio.NewScanner(&out)
		for i := 0; i < m; i++ {
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "test %d: no output line %d\n", t, i)
				os.Exit(1)
			}
			var got float64
			if _, err := fmt.Sscan(scanner.Text(), &got); err != nil {
				fmt.Fprintf(os.Stderr, "test %d: invalid output\n", t)
				os.Exit(1)
			}
			exp := expected[i]
			if math.Abs(got-exp) > 1e-6*math.Max(1, math.Abs(exp)) {
				fmt.Fprintf(os.Stderr, "test %d line %d: expected %.7f got %.7f\n", t, i, exp, got)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed.")
}
