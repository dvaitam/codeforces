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
}

type directedEdge struct {
	from, to int
}

func orientGraph(n int, edges []Edge) (bool, []directedEdge) {
	m := len(edges)
	adj := make([][]struct{ to, id int }, n+1)
	for i, e := range edges {
		id := i + 1
		adj[e.u] = append(adj[e.u], struct{ to, id int }{e.v, id})
		adj[e.v] = append(adj[e.v], struct{ to, id int }{e.u, id})
	}
	visited := make([]bool, n+1)
	dep := make([]int, n+1)
	low := make([]int, n+1)
	parent := make([]int, n+1)
	done := make([]bool, m+1)
	ans := make([]directedEdge, m+1)
	hasBridge := false

	var dfs func(x int)
	dfs = func(x int) {
		visited[x] = true
		low[x] = dep[x]
		for _, e := range adj[x] {
			j := e.to
			id := e.id
			if !visited[j] {
				if !done[id] {
					done[id] = true
					ans[id] = directedEdge{x, j}
				}
				parent[j] = x
				dep[j] = dep[x] + 1
				dfs(j)
				if low[j] < low[x] {
					low[x] = low[j]
				}
				if low[j] > dep[x] {
					hasBridge = true
				}
			} else if j != parent[x] {
				if dep[j] < low[x] {
					low[x] = dep[j]
				}
				if !done[id] {
					done[id] = true
					ans[id] = directedEdge{x, j}
				}
			}
		}
	}

	dfs(1)
	if hasBridge {
		return false, nil
	}
	out := make([]directedEdge, 0, m)
	for i := 1; i <= m; i++ {
		out = append(out, ans[i])
	}
	return true, out
}

func isStronglyConnected(n int, dirs []directedEdge) bool {
	adj := make([][]int, n+1)
	for _, e := range dirs {
		adj[e.from] = append(adj[e.from], e.to)
	}
	// bfs from 1
	vis := make([]bool, n+1)
	queue := []int{1}
	vis[1] = true
	for len(queue) > 0 {
		x := queue[0]
		queue = queue[1:]
		for _, y := range adj[x] {
			if !vis[y] {
				vis[y] = true
				queue = append(queue, y)
			}
		}
	}
	for i := 1; i <= n; i++ {
		if !vis[i] {
			return false
		}
	}
	// check reverse
	radj := make([][]int, n+1)
	for _, e := range dirs {
		radj[e.to] = append(radj[e.to], e.from)
	}
	vis = make([]bool, n+1)
	queue = []int{1}
	vis[1] = true
	for len(queue) > 0 {
		x := queue[0]
		queue = queue[1:]
		for _, y := range radj[x] {
			if !vis[y] {
				vis[y] = true
				queue = append(queue, y)
			}
		}
	}
	for i := 1; i <= n; i++ {
		if !vis[i] {
			return false
		}
	}
	return true
}

func generateCase(rng *rand.Rand) (string, int, bool, []directedEdge) {
	n := rng.Intn(8) + 2 // 2..9 nodes to keep tests small
	// Ensure connectivity: start with tree
	edges := make([]Edge, 0, n*n)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, Edge{p, i})
	}
	extra := rng.Intn(n) // add up to n-1 extra edges
	seen := make(map[[2]int]bool)
	for _, e := range edges {
		if e.u < e.v {
			seen[[2]int{e.u, e.v}] = true
		} else {
			seen[[2]int{e.v, e.u}] = true
		}
	}
	for i := 0; i < extra; i++ {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			i--
			continue
		}
		key := [2]int{u, v}
		if u > v {
			key = [2]int{v, u}
		}
		if seen[key] {
			i--
			continue
		}
		seen[key] = true
		edges = append(edges, Edge{u, v})
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
	}
	ok, dirs := orientGraph(n, edges)
	return sb.String(), n, ok, dirs
}

func expectedOutput(ok bool, dirs []directedEdge) string {
	if !ok {
		return "0"
	}
	var sb strings.Builder
	for _, d := range dirs {
		sb.WriteString(fmt.Sprintf("%d %d\n", d.from, d.to))
	}
	return strings.TrimRight(sb.String(), "\n")
}

func runCase(bin string, input string, n int, expectOK bool, expDirs []directedEdge) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expectedOutput(expectOK, expDirs) {
		return fmt.Errorf("expected:\n%s\n\ngot:\n%s", expectedOutput(expectOK, expDirs), got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, n, ok, dirs := generateCase(rng)
		if err := runCase(bin, in, n, ok, dirs); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
