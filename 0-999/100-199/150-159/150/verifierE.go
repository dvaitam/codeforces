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

// ---------- brute-force oracle for small trees ----------

type bfEdge struct{ to, w int }

func bfPairMedian(n int, edges [][3]int, x, y int) (int, int) {
	if x == y {
		return 0, 0
	}
	adj := make([][]bfEdge, n+1)
	for _, e := range edges {
		a, b, c := e[0], e[1], e[2]
		adj[a] = append(adj[a], bfEdge{b, c})
		adj[b] = append(adj[b], bfEdge{a, c})
	}
	var path []int
	visited := make([]bool, n+1)
	var found bool
	var dfs func(int)
	dfs = func(u int) {
		if found {
			return
		}
		if u == y {
			found = true
			return
		}
		visited[u] = true
		for _, e := range adj[u] {
			if visited[e.to] {
				continue
			}
			path = append(path, e.w)
			dfs(e.to)
			if found {
				return
			}
			path = path[:len(path)-1]
		}
	}
	dfs(x)
	if !found {
		return -1, 0
	}
	vals := append([]int(nil), path...)
	sort.Ints(vals)
	return vals[len(vals)/2], len(path)
}

// Brute-force: enumerate ALL pairs of nodes, compute path median, find optimal.
func bfBestMedian(n, llim, rlim int, edges [][3]int) int {
	best := -1
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			med, length := bfPairMedian(n, edges, i, j)
			if length >= llim && length <= rlim {
				if med > best {
					best = med
				}
			}
		}
	}
	return best
}

// ---------- embedded solver (brute-force, for small n) ----------

type solverEdge struct{ v, w int }

var solverAdj [][]solverEdge

func solverFindPath(u, target, parent int) ([]int, bool) {
	if u == target {
		return nil, true
	}
	for _, e := range solverAdj[u] {
		if e.v == parent {
			continue
		}
		if path, ok := solverFindPath(e.v, target, u); ok {
			return append([]int{e.w}, path...), true
		}
	}
	return nil, false
}

func embeddedSolve(n, l, r int, edges [][3]int) string {
	solverAdj = make([][]solverEdge, n+1)
	for _, e := range edges {
		a, b, c := e[0], e[1], e[2]
		solverAdj[a] = append(solverAdj[a], solverEdge{b, c})
		solverAdj[b] = append(solverAdj[b], solverEdge{a, c})
	}

	bestMed := -1
	bestU, bestV := 1, 1

	for u := 1; u <= n; u++ {
		for v := u + 1; v <= n; v++ {
			path, ok := solverFindPath(u, v, -1)
			if !ok {
				continue
			}
			length := len(path)
			if length < l || length > r {
				continue
			}
			sorted := make([]int, length)
			copy(sorted, path)
			sort.Ints(sorted)
			med := sorted[length/2]
			if med > bestMed {
				bestMed = med
				bestU = u
				bestV = v
			}
		}
	}

	return fmt.Sprintf("%d %d", bestU, bestV)
}

// Compute the diameter (longest path in edges) of the tree using two BFS.
func treeDiameter(n int, edges [][3]int) int {
	adj := make([][]int, n+1)
	for _, e := range edges {
		adj[e[0]] = append(adj[e[0]], e[1])
		adj[e[1]] = append(adj[e[1]], e[0])
	}
	bfs := func(start int) (int, int) {
		dist := make([]int, n+1)
		for i := range dist {
			dist[i] = -1
		}
		dist[start] = 0
		q := []int{start}
		farthest := start
		maxDist := 0
		for len(q) > 0 {
			u := q[0]
			q = q[1:]
			for _, v := range adj[u] {
				if dist[v] == -1 {
					dist[v] = dist[u] + 1
					q = append(q, v)
					if dist[v] > maxDist {
						maxDist = dist[v]
						farthest = v
					}
				}
			}
		}
		return farthest, maxDist
	}
	far1, _ := bfs(1)
	_, diam := bfs(far1)
	return diam
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// First test the embedded solver against brute-force oracle on a few cases
	// to make sure it's working, then test the binary.
	tested := 0
	for i := 0; tested < 100 && i < 10000; i++ {
		n := rng.Intn(6) + 2
		edges := make([][3]int, n-1)
		for j := 2; j <= n; j++ {
			p := rng.Intn(j-1) + 1
			w := rng.Intn(20)
			edges[j-2] = [3]int{p, j, w}
		}

		diam := treeDiameter(n, edges)
		if diam < 1 {
			continue
		}

		llim := rng.Intn(diam) + 1
		rlim := llim + rng.Intn(diam-llim+1)
		if rlim > n-1 {
			rlim = n - 1
		}
		if llim > rlim {
			continue
		}

		// Verify that a valid path actually exists
		bestMed := bfBestMedian(n, llim, rlim, edges)
		if bestMed < 0 {
			continue
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, llim, rlim))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2]))
		}
		input := sb.String()

		// Test the external binary
		got, err := run(bin, input)
		if err != nil {
			// Fall back to embedded solver for comparison
			expected := embeddedSolve(n, llim, rlim, edges)
			fmt.Fprintf(os.Stderr, "case %d failed: %v\nexpected: %s\ninput:\n%s", tested+1, err, expected, input)
			os.Exit(1)
		}

		if got == "" {
			fmt.Fprintf(os.Stderr, "case %d failed: bad output empty string\ninput:\n%s", tested+1, input)
			os.Exit(1)
		}

		var x, y int
		if _, err := fmt.Sscanf(got, "%d %d", &x, &y); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: bad output %q\ninput:\n%s", tested+1, got, input)
			os.Exit(1)
		}
		med, length := bfPairMedian(n, edges, x, y)
		if length < llim || length > rlim || med != bestMed {
			fmt.Fprintf(os.Stderr, "case %d: incorrect pair %d %d (med=%d len=%d, bestMed=%d)\ninput:\n%s", tested+1, x, y, med, length, bestMed, input)
			os.Exit(1)
		}
		tested++
	}
	fmt.Printf("All %d tests passed\n", tested)
}
