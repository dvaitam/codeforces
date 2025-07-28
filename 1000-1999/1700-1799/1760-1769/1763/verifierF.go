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

type pair struct{ u, v int }

func key(u, v int) pair {
	if u > v {
		u, v = v, u
	}
	return pair{u, v}
}

func edgesOnPaths(n int, edges []pair, a, b int) map[pair]bool {
	adj := make([][]int, n+1)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
	}
	res := map[pair]bool{}
	visited := make([]bool, n+1)
	var dfs func(int, []pair)
	dfs = func(u int, path []pair) {
		if u == b {
			for _, e := range path {
				res[e] = true
			}
			return
		}
		visited[u] = true
		for _, v := range adj[u] {
			if !visited[v] {
				dfs(v, append(path, key(u, v)))
			}
		}
		visited[u] = false
	}
	dfs(a, nil)
	return res
}

func reachable(n int, edges []pair, a, b int, rem pair) bool {
	adj := make([][]int, n+1)
	for _, e := range edges {
		if e == rem {
			continue
		}
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
	}
	q := []int{a}
	vis := make([]bool, n+1)
	vis[a] = true
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		if u == b {
			return true
		}
		for _, v := range adj[u] {
			if !vis[v] {
				vis[v] = true
				q = append(q, v)
			}
		}
	}
	return false
}

func queryAns(n int, edges []pair, a, b int) int {
	use := edgesOnPaths(n, edges, a, b)
	cnt := 0
	for _, e := range edges {
		if use[e] && reachable(n, edges, a, b, e) {
			cnt++
		}
	}
	return cnt
}

func expected(n int, edges []pair, queries [][2]int) string {
	var sb strings.Builder
	for i, q := range queries {
		ans := queryAns(n, edges, q[0], q[1])
		sb.WriteString(fmt.Sprintf("%d", ans))
		if i+1 < len(queries) {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges-(n-1)+1) + (n - 1)
	edgesMap := map[pair]bool{}
	edges := make([]pair, 0, m)
	// ensure connected by building a tree first
	for i := 2; i <= n; i++ {
		v := rng.Intn(i-1) + 1
		p := key(i, v)
		edges = append(edges, p)
		edgesMap[p] = true
	}
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		p := key(u, v)
		if edgesMap[p] {
			continue
		}
		edgesMap[p] = true
		edges = append(edges, p)
	}
	q := rng.Intn(5) + 1
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		queries[i] = [2]int{a, b}
	}
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
	for _, e := range edges {
		input.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
	}
	input.WriteString(fmt.Sprintf("%d\n", q))
	for _, qu := range queries {
		input.WriteString(fmt.Sprintf("%d %d\n", qu[0], qu[1]))
	}
	return input.String(), expected(n, edges, queries)
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
