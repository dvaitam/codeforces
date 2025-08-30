package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Edge struct {
	u, v int
	w    int
}

// Build directed adjacency: for each input edge (u,v,w), add u->v cost 0 and v->u cost w.
func buildAdj(n int, edges []Edge) [][]struct{ to, c int } {
	adj := make([][]struct{ to, c int }, n)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], struct{ to, c int }{e.v, 0})
		adj[e.v] = append(adj[e.v], struct{ to, c int }{e.u, e.w})
	}
	return adj
}

func dfsFrom(adj [][]struct{ to, c int }, k int, start int, vis []bool) int {
	stack := []int{start}
	vis[start] = true
	cnt := 0
	for len(stack) > 0 {
		x := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		cnt++
		for _, e := range adj[x] {
			if e.c <= k && !vis[e.to] {
				vis[e.to] = true
				stack = append(stack, e.to)
			}
		}
	}
	return cnt
}

// Check if the directed graph (with edges allowed when cost<=k) has a mother vertex.
func motherExists(n int, adj [][]struct{ to, c int }, k int) bool {
	vis := make([]bool, n)
	last := 0
	for i := 0; i < n; i++ {
		if !vis[i] {
			dfsFrom(adj, k, i, vis)
			last = i
		}
	}
	for i := range vis {
		vis[i] = false
	}
	return dfsFrom(adj, k, last, vis) == n
}

// Fast undirected connectivity check (ignoring weights and directions).
func isConnected(n int, edges []Edge) bool {
	if n == 0 {
		return true
	}
	ug := make([][]int, n)
	for _, e := range edges {
		ug[e.u] = append(ug[e.u], e.v)
		ug[e.v] = append(ug[e.v], e.u)
	}
	vis := make([]bool, n)
	q := []int{0}
	vis[0] = true
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		for _, v := range ug[u] {
			if !vis[v] {
				vis[v] = true
				q = append(q, v)
			}
		}
	}
	for i := 0; i < n; i++ {
		if !vis[i] {
			return false
		}
	}
	return true
}

func solveCase(n, m int, edges []Edge) int {
	// Early exit: if the underlying undirected graph is disconnected, impossible.
	if !isConnected(n, edges) {
		return -1
	}
	adj := buildAdj(n, edges)
	lo, hi := 0, int(1e9)+10
	for lo < hi {
		mid := (lo + hi) >> 1
		if motherExists(n, adj, mid) {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	if lo > int(1e9) {
		return -1
	}
	return lo
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + n - 1
	edges := make([]Edge, m)
	// random edges; allow zero weights too
	for i := 0; i < m; i++ {
		u := rng.Intn(n)
		v := rng.Intn(n)
		for u == v {
			v = rng.Intn(n)
		}
		edges[i] = Edge{u: u, v: v, w: rng.Intn(21)}
	}
	ans := solveCase(n, m, edges)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e.u+1, e.v+1, e.w))
	}
	return sb.String(), fmt.Sprintf("%d\n", ans)
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if outStr != exp {
		return fmt.Errorf("expected %s got %s", exp, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
