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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type edge struct{ u, v int }

func isBipartite(n int, edges []edge, skip int) bool {
	adj := make([][]int, n+1)
	for i, e := range edges {
		if i == skip {
			continue
		}
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
	}
	color := make([]int, n+1)
	for i := 1; i <= n; i++ {
		if color[i] != 0 {
			continue
		}
		color[i] = 1
		q := []int{i}
		for len(q) > 0 {
			v := q[0]
			q = q[1:]
			for _, to := range adj[v] {
				if color[to] == 0 {
					color[to] = -color[v]
					q = append(q, to)
				} else if color[to] == color[v] {
					return false
				}
			}
		}
	}
	return true
}

func solveCase(n int, edges []edge) []int {
	m := len(edges)
	if isBipartite(n, edges, -1) {
		res := make([]int, m)
		for i := range res {
			res[i] = i + 1
		}
		return res
	}
	var ans []int
	for i := 0; i < m; i++ {
		if isBipartite(n, edges, i) {
			ans = append(ans, i+1)
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges + 1)
	used := make(map[[2]int]bool)
	edges := make([]edge, m)
	for i := 0; i < m; i++ {
		var u, v int
		for {
			u = rng.Intn(n) + 1
			v = rng.Intn(n) + 1
			if u != v {
				if u > v {
					u, v = v, u
				}
				if !used[[2]int{u, v}] {
					break
				}
			}
		}
		used[[2]int{u, v}] = true
		edges[i] = edge{u, v}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
	}
	ans := solveCase(n, edges)
	sort.Ints(ans)
	var out strings.Builder
	fmt.Fprintf(&out, "%d\n", len(ans))
	for i, id := range ans {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprintf(&out, "%d", id)
	}
	out.WriteByte('\n')
	return sb.String(), out.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
