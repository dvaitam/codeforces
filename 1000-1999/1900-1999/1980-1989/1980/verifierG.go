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

// Brute-force reference solver for small trees
type treeEdge struct {
	u, v, w int
}

type query struct {
	typ int // 1 = xor all edges, 2 = query
	y   int // for type 1
	v   int // for type 2
	x   int // for type 2
}

// Compute XOR distance from root to all nodes using DFS
func computeDist(n int, adj [][]struct{ to, w int }) []int {
	dist := make([]int, n+1)
	visited := make([]bool, n+1)
	var dfs func(u int)
	dfs = func(u int) {
		visited[u] = true
		for _, e := range adj[u] {
			if !visited[e.to] {
				dist[e.to] = dist[u] ^ e.w
				dfs(e.to)
			}
		}
	}
	dfs(1)
	return dist
}

func solveCase(n, q int, edges []treeEdge, queries []query) []int {
	// Build adjacency list with mutable weights
	weights := make([]int, len(edges))
	for i, e := range edges {
		weights[i] = e.w
	}

	var answers []int
	xorAccum := 0 // accumulated XOR from type-1 queries

	for _, qr := range queries {
		if qr.typ == 1 {
			xorAccum ^= qr.y
		} else {
			// For query type 2: find max XOR of (dist(1,v) ^ dist(1,u) ^ x) over all u != v
			// where dist uses current edge weights (original XOR'd with xorAccum based on parity)
			// Since each edge weight w becomes w ^ xorAccum, and dist(1,u) = XOR of edges on path,
			// For a tree with n-1 edges on path from 1 to u, each edge gets XOR'd with xorAccum.
			// But the number of edges varies per path, so we need to recompute.

			// Build adjacency with current weights
			adj := make([][]struct{ to, w int }, n+1)
			for i, e := range edges {
				curW := weights[i] ^ xorAccum
				adj[e.u] = append(adj[e.u], struct{ to, w int }{e.v, curW})
				adj[e.v] = append(adj[e.v], struct{ to, w int }{e.u, curW})
			}

			dist := computeDist(n, adj)

			// The cycle from v to u via the added edge has XOR = dist(v) ^ dist(u) ^ x
			// (since dist(v,u) in tree = dist(1,v) ^ dist(1,u), and cycle = dist(v,u) ^ x)
			best := 0
			for u := 1; u <= n; u++ {
				if u == qr.v {
					continue
				}
				val := dist[qr.v] ^ dist[u] ^ qr.x
				if val > best {
					best = val
				}
			}
			answers = append(answers, best)
		}
	}
	return answers
}

func generateTree(rng *rand.Rand, n int) []treeEdge {
	edges := make([]treeEdge, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		w := rng.Intn(10) + 1
		edges[i-2] = treeEdge{p, i, w}
	}
	return edges
}

func generateCase(rng *rand.Rand) (string, []int) {
	n := rng.Intn(4) + 2
	q := rng.Intn(5) + 1
	edges := generateTree(rng, n)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e.u, e.v, e.w))
	}

	queries := make([]query, q)
	for i := 0; i < q; i++ {
		if rng.Intn(2) == 0 {
			y := rng.Intn(10) + 1
			sb.WriteString(fmt.Sprintf("^ %d\n", y))
			queries[i] = query{typ: 1, y: y}
		} else {
			v := rng.Intn(n) + 1
			x := rng.Intn(10) + 1
			sb.WriteString(fmt.Sprintf("? %d %d\n", v, x))
			queries[i] = query{typ: 2, v: v, x: x}
		}
	}

	expected := solveCase(n, q, edges, queries)
	return sb.String(), expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expected := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}

		// Parse candidate output
		outFields := strings.Fields(out)
		if len(outFields) != len(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d values, got %d\ninput:\n%soutput:\n%s\n", i+1, len(expected), len(outFields), in, out)
			os.Exit(1)
		}
		for j, exp := range expected {
			var got int
			if _, err := fmt.Sscan(outFields[j], &got); err != nil {
				fmt.Fprintf(os.Stderr, "case %d failed: cannot parse output value %q\n", i+1, outFields[j])
				os.Exit(1)
			}
			if got != exp {
				fmt.Fprintf(os.Stderr, "case %d query %d failed: expected %d got %d\ninput:\n%soutput:\n%s\n", i+1, j+1, exp, got, in, out)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
