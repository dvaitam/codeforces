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

type edge struct {
	u, v, w int
}

type testCase struct {
	n, m, T int
	edges   []edge
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.T))
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e.u, e.v, e.w))
	}
	return sb.String()
}

func expected(tc testCase) string {
	adj := make([][]edge, tc.n+1)
	for _, e := range tc.edges {
		adj[e.u] = append(adj[e.u], e)
	}
	dp := make([]map[int]int, tc.n+1)
	visited := make([]bool, tc.n+1)
	var dfs func(int)
	dfs = func(u int) {
		if visited[u] {
			return
		}
		visited[u] = true
		dp[u] = make(map[int]int)
		if u == tc.n {
			dp[u][1] = 0
			return
		}
		for _, ed := range adj[u] {
			dfs(ed.v)
			for k, val := range dp[ed.v] {
				wsum := val + ed.w
				if wsum > tc.T {
					continue
				}
				nk := k + 1
				if prev, ok := dp[u][nk]; !ok || wsum < prev {
					dp[u][nk] = wsum
				}
			}
		}
	}
	dfs(1)
	best := 0
	for k := range dp[1] {
		if k > best {
			best = k
		}
	}
	path := make([]int, 0, best)
	cur, k := 1, best
	for {
		path = append(path, cur)
		if cur == tc.n {
			break
		}
		for _, ed := range adj[cur] {
			if next, ok := dp[ed.v][k-1]; ok && next+ed.w == dp[cur][k] {
				cur = ed.v
				k--
				break
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintln(best))
	for i, x := range path {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(x))
	}
	return strings.TrimSpace(sb.String())
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(6) + 2 // 2..7
	edges := make([]edge, 0)
	sum := 0
	for i := 1; i < n; i++ {
		w := rng.Intn(5) + 1
		edges = append(edges, edge{i, i + 1, w})
		sum += w
	}
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges-(n-1)+1) + (n - 1)
	for len(edges) < m {
		u := rng.Intn(n-1) + 1
		v := rng.Intn(n-u) + u + 1
		w := rng.Intn(5) + 1
		edges = append(edges, edge{u, v, w})
	}
	T := sum + rng.Intn(10)
	return testCase{n: n, m: len(edges), T: T, edges: edges}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{}
	for len(cases) < 105 {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		input := buildInput(tc)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		exp := expected(tc)
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
