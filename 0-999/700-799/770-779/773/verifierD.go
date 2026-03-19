package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func run(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// normalizeTokens extracts all whitespace-separated tokens and joins them with single spaces.
func normalizeTokens(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// refSolve is the correct embedded reference solver for 773D.
// Given an upper-triangular weight matrix for a complete graph, it computes
// for each root r: sum over all other vertices v of the max edge weight
// on the path from r to v in the MST (Prim's algorithm).
func refSolve(input string) string {
	tokens := strings.Fields(input)
	idx := 0
	nextInt := func() int {
		v := 0
		s := tokens[idx]
		idx++
		for _, c := range s {
			v = v*10 + int(c-'0')
		}
		return v
	}

	n := nextInt()
	if n == 0 {
		return ""
	}

	adj := make([][]int, n)
	for i := 0; i < n; i++ {
		adj[i] = make([]int, n)
	}
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			w := nextInt()
			adj[i][j] = w
			adj[j][i] = w
		}
	}

	// Prim's MST
	inTree := make([]bool, n)
	minEdge := make([]int, n)
	parent := make([]int, n)
	for i := 0; i < n; i++ {
		minEdge[i] = int(2e9)
		parent[i] = -1
	}
	minEdge[0] = 0

	type Edge struct{ to, w int }
	mst := make([][]Edge, n)

	for i := 0; i < n; i++ {
		u := -1
		for j := 0; j < n; j++ {
			if !inTree[j] && (u == -1 || minEdge[j] < minEdge[u]) {
				u = j
			}
		}
		inTree[u] = true
		if parent[u] != -1 {
			mst[u] = append(mst[u], Edge{parent[u], minEdge[u]})
			mst[parent[u]] = append(mst[parent[u]], Edge{u, minEdge[u]})
		}
		for v := 0; v < n; v++ {
			if !inTree[v] && adj[u][v] < minEdge[v] {
				minEdge[v] = adj[u][v]
				parent[v] = u
			}
		}
	}

	var parts []string
	for r := 0; r < n; r++ {
		var sum int64
		var dfs func(u, p, maxW int)
		dfs = func(u, p, maxW int) {
			if u != r {
				sum += int64(maxW)
			}
			for _, edge := range mst[u] {
				if edge.to != p {
					newMax := maxW
					if edge.w > newMax {
						newMax = edge.w
					}
					dfs(edge.to, u, newMax)
				}
			}
		}
		dfs(r, -1, 0)
		parts = append(parts, fmt.Sprintf("%d", sum))
	}
	return strings.Join(parts, " ")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	rand.Seed(42)
	for t := 0; t < 100; t++ {
		n := rand.Intn(4) + 2
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n-1; i++ {
			for j := i + 1; j < n; j++ {
				w := rand.Intn(10) + 1
				sb.WriteString(fmt.Sprintf("%d ", w))
			}
			if i < n-2 {
				sb.WriteByte('\n')
			}
		}
		sb.WriteByte('\n')
		input := sb.String()
		expect := refSolve(input)
		got, err := run(userBin, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "program failed on test", t+1, ":", err)
			os.Exit(1)
		}
		if normalizeTokens(expect) != normalizeTokens(got) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %s got %s\n", t+1, normalizeTokens(expect), normalizeTokens(got))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
