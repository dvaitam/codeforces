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
	w    int64
}

func bfs(adj [][]int, start int) []bool {
	n := len(adj)
	vis := make([]bool, n)
	q := []int{start}
	vis[start] = true
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		for _, to := range adj[v] {
			if !vis[to] {
				vis[to] = true
				q = append(q, to)
			}
		}
	}
	return vis
}

func solveCase(n, m, k int, specials []int, c []int64, w []int64, edges []Edge) []int64 {
	best := make([]int64, n)
	for i := range best {
		best[i] = int64(-1 << 60)
	}
	pow := 1
	for i := 0; i < m; i++ {
		pow *= 3
	}
	orient := make([]int, m)
	for mask := 0; mask < pow; mask++ {
		tmp := mask
		cost := int64(0)
		for i := 0; i < m; i++ {
			orient[i] = tmp % 3
			tmp /= 3
			if orient[i] == 2 {
				cost += w[i]
			}
		}
		adj := make([][]int, n)
		for i, e := range edges {
			switch orient[i] {
			case 0:
				adj[e.u] = append(adj[e.u], e.v)
			case 1:
				adj[e.v] = append(adj[e.v], e.u)
			case 2:
				adj[e.u] = append(adj[e.u], e.v)
				adj[e.v] = append(adj[e.v], e.u)
			}
		}
		reachAll := make([]bool, n)
		for i := range reachAll {
			reachAll[i] = true
		}
		for _, s := range specials {
			vis := bfs(adj, s)
			for i := 0; i < n; i++ {
				if !vis[i] {
					reachAll[i] = false
				}
			}
		}
		profit := -cost
		for i := 0; i < n; i++ {
			if reachAll[i] {
				profit += c[i]
			}
		}
		for i := 0; i < n; i++ {
			if reachAll[i] && profit > best[i] {
				best[i] = profit
			}
		}
	}
	return best
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genGraph(rng *rand.Rand, n, m int) []Edge {
	edges := make([]Edge, 0, m)
	// ensure tree edges for connectivity
	for i := 1; i < n; i++ {
		u := i - 1
		v := i
		w := int64(rng.Intn(5))
		edges = append(edges, Edge{u, v, w})
	}
	for len(edges) < m {
		u := rng.Intn(n)
		v := rng.Intn(n)
		if u == v {
			continue
		}
		dup := false
		for _, e := range edges {
			if (e.u == u && e.v == v) || (e.u == v && e.v == u) {
				dup = true
				break
			}
		}
		if dup {
			continue
		}
		w := int64(rng.Intn(5))
		edges = append(edges, Edge{u, v, w})
	}
	return edges
}

func genTest(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 2
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(min(maxEdges, 4-n+1)) + n - 1
	k := rng.Intn(n) + 1
	specials := make([]int, k)
	used := make([]bool, n)
	for i := 0; i < k; i++ {
		for {
			v := rng.Intn(n)
			if !used[v] {
				specials[i] = v
				used[v] = true
				break
			}
		}
	}
	c := make([]int64, n)
	for i := 0; i < n; i++ {
		c[i] = int64(rng.Intn(6))
	}
	edges := genGraph(rng, n, m)
	w := make([]int64, m)
	for i := 0; i < m; i++ {
		w[i] = int64(rng.Intn(6))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("1\n%d %d %d\n", n, m, k))
	for i := 0; i < k; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", specials[i]+1))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", c[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < m; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", w[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < m; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", edges[i].u+1, edges[i].v+1))
	}
	best := solveCase(n, m, k, specials, c, w, edges)
	var exp strings.Builder
	for i := 0; i < n; i++ {
		if i > 0 {
			exp.WriteByte(' ')
		}
		exp.WriteString(fmt.Sprintf("%d", best[i]))
	}
	return sb.String(), exp.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input, exp := genTest(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i, input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
