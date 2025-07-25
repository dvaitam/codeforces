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

type edge struct{ a, b int }

type testCase struct {
	n     int
	edges []edge
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(3) + 4 // 4..6
	m := n
	edges := make([]edge, 0, m*2)
	for i := 0; i < n; i++ { // cycle
		edges = append(edges, edge{i + 1, (i+1)%n + 1})
	}
	for len(edges) < n*3 {
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		if a == b {
			continue
		}
		edges = append(edges, edge{a, b})
	}
	return testCase{n: n, edges: edges}
}

func bfs(n int, adj [][]int, start int) []int {
	dist := make([]int, n)
	for i := 0; i < n; i++ {
		dist[i] = -1
	}
	q := []int{start}
	dist[start] = 0
	for h := 0; h < len(q); h++ {
		v := q[h]
		for _, to := range adj[v] {
			if dist[to] == -1 {
				dist[to] = dist[v] + 1
				q = append(q, to)
			}
		}
	}
	return dist
}

func allDist(tc testCase) [][]int {
	adj := make([][]int, tc.n)
	for _, e := range tc.edges {
		adj[e.a-1] = append(adj[e.a-1], e.b-1)
	}
	d := make([][]int, tc.n)
	for i := 0; i < tc.n; i++ {
		d[i] = bfs(tc.n, adj, i)
	}
	return d
}

func solve(tc testCase) (int, [4]int) {
	n := tc.n
	adj := make([][]int, n)
	for _, e := range tc.edges {
		adj[e.a-1] = append(adj[e.a-1], e.b-1)
	}
	dist := make([][]int, n)
	for i := 0; i < n; i++ {
		dist[i] = bfs(n, adj, i)
	}
	bestVal := -1
	res := [4]int{}
	for a := 0; a < n; a++ {
		for b := 0; b < n; b++ {
			if a == b || dist[a][b] == -1 {
				continue
			}
			for c := 0; c < n; c++ {
				if c == a || c == b || dist[b][c] == -1 {
					continue
				}
				for d := 0; d < n; d++ {
					if d == a || d == b || d == c || dist[c][d] == -1 {
						continue
					}
					val := dist[a][b] + dist[b][c] + dist[c][d]
					if val > bestVal {
						bestVal = val
						res = [4]int{a + 1, b + 1, c + 1, d + 1}
					}
				}
			}
		}
	}
	return bestVal, res
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
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, len(tc.edges)))
		for _, e := range tc.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e.a, e.b))
		}
		best, _ := solve(tc)
		gotStr, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		var a, b, c, d int
		if _, err := fmt.Sscan(gotStr, &a, &b, &c, &d); err != nil {
			fmt.Fprintf(os.Stderr, "case %d invalid output: %s\n", i+1, gotStr)
			os.Exit(1)
		}
		dist := allDist(tc)
		if a < 1 || a > tc.n || b < 1 || b > tc.n || c < 1 || c > tc.n || d < 1 || d > tc.n {
			fmt.Fprintf(os.Stderr, "case %d invalid city index\n", i+1)
			os.Exit(1)
		}
		if dist[a-1][b-1] == -1 || dist[b-1][c-1] == -1 || dist[c-1][d-1] == -1 {
			fmt.Fprintf(os.Stderr, "case %d impossible path\n", i+1)
			os.Exit(1)
		}
		val := dist[a-1][b-1] + dist[b-1][c-1] + dist[c-1][d-1]
		if val != best {
			fmt.Fprintf(os.Stderr, "case %d failed expected %d got %d\n", i+1, best, val)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
