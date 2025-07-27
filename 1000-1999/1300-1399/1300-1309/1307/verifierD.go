package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func runCandidate(bin, input string) (string, error) {
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

func bfs(start, n int, adj [][]int) []int {
	const INF = int(1e9)
	dist := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = INF
	}
	q := make([]int, 0, n)
	dist[start] = 0
	q = append(q, start)
	for i := 0; i < len(q); i++ {
		u := q[i]
		du := dist[u]
		for _, v := range adj[u] {
			if dist[v] > du+1 {
				dist[v] = du + 1
				q = append(q, v)
			}
		}
	}
	return dist
}

func solve(input string) string {
	in := strings.NewReader(input)
	var n, m, k int
	fmt.Fscan(in, &n, &m, &k)
	special := make([]int, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(in, &special[i])
	}
	adj := make([][]int, n+1)
	for i := 0; i < m; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		adj[x] = append(adj[x], y)
		adj[y] = append(adj[y], x)
	}
	d1 := bfs(1, n, adj)
	d2 := bfs(n, n, adj)
	L0 := d1[n]
	sort.Slice(special, func(i, j int) bool {
		return d1[special[i]]-d2[special[i]] < d1[special[j]]-d2[special[j]]
	})
	best := 0
	mxD1 := -1000000000
	for _, u := range special {
		if mxD1 > -1000000000 {
			cand := mxD1 + d2[u] + 1
			if cand > best {
				best = cand
			}
		}
		if d1[u] > mxD1 {
			mxD1 = d1[u]
		}
	}
	if best > L0 {
		best = L0
	}
	if best == 0 {
		best = L0
	}
	return fmt.Sprintf("%d", best)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 2
	m := n - 1 + rng.Intn(n)
	k := rng.Intn(n-1) + 2
	adjSet := make(map[[2]int]struct{})
	edges := make([][2]int, 0, m)
	// build tree
	for i := 2; i <= n; i++ {
		j := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{i, j})
		adjSet[[2]int{i, j}] = struct{}{}
		adjSet[[2]int{j, i}] = struct{}{}
	}
	// add extra edges
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if _, ok := adjSet[[2]int{u, v}]; ok {
			continue
		}
		edges = append(edges, [2]int{u, v})
		adjSet[[2]int{u, v}] = struct{}{}
		adjSet[[2]int{v, u}] = struct{}{}
	}
	specials := rng.Perm(n)[:k]
	for i := range specials {
		specials[i]++
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
	for i, v := range specials {
		if i+1 == k {
			fmt.Fprintf(&sb, "%d\n", v)
		} else {
			fmt.Fprintf(&sb, "%d ", v)
		}
	}
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	in := sb.String()
	return in, solve(in)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		got, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
