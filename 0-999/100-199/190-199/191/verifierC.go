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
	to, id int
}

func runCandidate(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(n int, edges [][2]int, k int, queries [][2]int) string {
	adj := make([][]Edge, n+1)
	for i, e := range edges {
		u, v := e[0], e[1]
		id := i + 1
		adj[u] = append(adj[u], Edge{v, id})
		adj[v] = append(adj[v], Edge{u, id})
	}
	parent := make([]int, n+1)
	parentEdge := make([]int, n+1)
	depth := make([]int, n+1)
	order := make([]int, 0, n)
	queue := []int{1}
	parent[1] = 0
	depth[1] = 0
	for i := 0; i < len(queue); i++ {
		u := queue[i]
		order = append(order, u)
		for _, e := range adj[u] {
			v := e.to
			if v == parent[u] {
				continue
			}
			parent[v] = u
			parentEdge[v] = e.id
			depth[v] = depth[u] + 1
			queue = append(queue, v)
		}
	}
	const maxLog = 18
	p := make([][]int, maxLog)
	p[0] = make([]int, n+1)
	for v := 1; v <= n; v++ {
		p[0][v] = parent[v]
	}
	for l := 1; l < maxLog; l++ {
		p[l] = make([]int, n+1)
		for v := 1; v <= n; v++ {
			p[l][v] = p[l-1][p[l-1][v]]
		}
	}
	lca := func(u, v int) int {
		if depth[u] < depth[v] {
			u, v = v, u
		}
		diff := depth[u] - depth[v]
		for l := 0; l < maxLog; l++ {
			if diff&(1<<l) != 0 {
				u = p[l][u]
			}
		}
		if u == v {
			return u
		}
		for l := maxLog - 1; l >= 0; l-- {
			if p[l][u] != p[l][v] {
				u = p[l][u]
				v = p[l][v]
			}
		}
		return parent[u]
	}
	cnt := make([]int64, n+1)
	for _, q := range queries {
		a, b := q[0], q[1]
		cnt[a]++
		cnt[b]++
		l := lca(a, b)
		cnt[l] -= 2
	}
	ans := make([]int64, n)
	for i := len(order) - 1; i > 0; i-- {
		v := order[i]
		eid := parentEdge[v]
		ans[eid-1] = cnt[v]
		cnt[parent[v]] += cnt[v]
	}
	var sb strings.Builder
	for i := 0; i < n-1; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", ans[i]))
	}
	return sb.String()
}

func generateTree(rng *rand.Rand, n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
	}
	return edges
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(7) + 2
	edges := generateTree(rng, n)
	k := rng.Intn(n) + 1
	queries := make([][2]int, k)
	for i := 0; i < k; i++ {
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		for b == a {
			b = rng.Intn(n) + 1
		}
		queries[i] = [2]int{a, b}
	}
	var in strings.Builder
	fmt.Fprintf(&in, "%d\n", n)
	for _, e := range edges {
		fmt.Fprintf(&in, "%d %d\n", e[0], e[1])
	}
	fmt.Fprintf(&in, "%d\n", k)
	for _, q := range queries {
		fmt.Fprintf(&in, "%d %d\n", q[0], q[1])
	}
	exp := solve(n, edges, k, queries)
	return in.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		input, expect := generateCase(rng)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
