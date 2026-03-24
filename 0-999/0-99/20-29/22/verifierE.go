package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// computeMinEdges computes the minimum number of edges to add.
// The functional graph has components; each node i -> f[i].
// We need to count the number of weakly connected components in the
// underlying undirected graph of the functional graph, then handle
// leaves (in-degree 0 nodes) and cycles.
func computeMinEdges(n int, f []int) int {
	// Find all weakly connected components and count cycles
	inDeg := make([]int, n+1)
	for i := 1; i <= n; i++ {
		inDeg[f[i]]++
	}

	compID := make([]int, n+1)
	numCycles := 0
	numLeaves := 0

	// Process trees (nodes with in-degree 0)
	for i := 1; i <= n; i++ {
		if inDeg[i] == 0 {
			numLeaves++
		}
	}

	// Find cycles
	vis := make([]int, n+1)
	for i := 1; i <= n; i++ {
		if vis[i] != 0 {
			continue
		}
		cur := i
		for vis[cur] == 0 {
			vis[cur] = i
			cur = f[cur]
		}
		if vis[cur] == i {
			// Found a new cycle
			numCycles++
			compID[cur] = numCycles
			next := f[cur]
			for next != cur {
				compID[next] = numCycles
				next = f[next]
			}
		}
	}

	if numCycles == 1 && numLeaves == 0 {
		return 0
	}

	// The answer is max(numLeaves, numCycles) when there are leaves,
	// or numCycles when there are no leaves but multiple cycles.
	// Actually for this functional graph problem:
	// - Each component has exactly one cycle
	// - numLeaves nodes have in-degree 0
	// - We need to connect all components and redirect leaves
	// The minimum edges = max(numLeaves, numCycles) if numLeaves > 0, else numCycles
	// Wait, let me reconsider. Each weakly connected component has one cycle.
	// Trees hang off cycles. Leaves are nodes with in-degree 0.

	// Count weakly connected components
	parent := make([]int, n+1)
	for i := 0; i <= n; i++ {
		parent[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		for parent[x] != x {
			parent[x] = parent[parent[x]]
			x = parent[x]
		}
		return x
	}
	union := func(a, b int) {
		a, b = find(a), find(b)
		if a != b {
			parent[a] = b
		}
	}
	for i := 1; i <= n; i++ {
		union(i, f[i])
	}
	compSet := make(map[int]bool)
	for i := 1; i <= n; i++ {
		compSet[find(i)] = true
	}
	numComponents := len(compSet)

	if numComponents == 1 && numLeaves == 0 {
		return 0
	}

	// For each component, we need at least one edge to connect them (numComponents - 1 for connectivity).
	// Plus we need to redirect each leaf.
	// But each added edge can both connect components and redirect a leaf.
	// The answer is max(numLeaves, numComponents).
	// If numLeaves == 0, answer = numComponents (each cycle-only component needs one edge to break into the whole).
	if numLeaves == 0 {
		return numComponents
	}
	if numLeaves > numComponents {
		return numLeaves
	}
	return numComponents
}

// checkAnswer validates that after adding the given edges, the graph becomes
// strongly connected (everyone can reach everyone).
func checkAnswer(n int, f []int, edges [][2]int) bool {
	// Build adjacency list with original edges + added edges
	adj := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		adj[i] = append(adj[i], f[i])
	}
	for _, e := range edges {
		adj[e[0]] = append(adj[e[0]], e[1])
	}

	// Check strong connectivity: BFS/DFS from node 1 should reach all,
	// and in the reverse graph, BFS/DFS from node 1 should reach all.
	reachable := func(adjList [][]int) bool {
		vis := make([]bool, n+1)
		queue := []int{1}
		vis[1] = true
		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			for _, v := range adjList[u] {
				if !vis[v] {
					vis[v] = true
					queue = append(queue, v)
				}
			}
		}
		for i := 1; i <= n; i++ {
			if !vis[i] {
				return false
			}
		}
		return true
	}

	if !reachable(adj) {
		return false
	}

	// Build reverse graph
	radj := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		radj[f[i]] = append(radj[f[i]], i)
	}
	for _, e := range edges {
		radj[e[1]] = append(radj[e[1]], e[0])
	}
	return reachable(radj)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n := rng.Intn(8) + 2
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			v := rng.Intn(n) + 1
			if v == i {
				if v < n {
					v++
				} else {
					v--
				}
			}
			a[i] = v
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 1; i <= n; i++ {
			if i > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", a[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()

		expectedCount := computeMinEdges(n, a)

		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t+1, err, input)
			os.Exit(1)
		}

		lines := strings.Split(got, "\n")
		if len(lines) == 0 {
			fmt.Fprintf(os.Stderr, "case %d failed: empty output\ninput:\n%s", t+1, input)
			os.Exit(1)
		}

		gotCount, err := strconv.Atoi(strings.TrimSpace(lines[0]))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: cannot parse count: %v\ninput:\n%s", t+1, err, input)
			os.Exit(1)
		}

		if gotCount != expectedCount {
			fmt.Fprintf(os.Stderr, "case %d failed: expected count %d got %d\ninput:\n%s", t+1, expectedCount, gotCount, input)
			os.Exit(1)
		}

		if gotCount == 0 {
			continue
		}

		if len(lines) < gotCount+1 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d edge lines, got %d lines total\ninput:\n%s", t+1, gotCount, len(lines), input)
			os.Exit(1)
		}

		edges := make([][2]int, gotCount)
		for i := 0; i < gotCount; i++ {
			fields := strings.Fields(lines[i+1])
			if len(fields) != 2 {
				fmt.Fprintf(os.Stderr, "case %d failed: edge line %d has %d fields\ninput:\n%s", t+1, i+1, len(fields), input)
				os.Exit(1)
			}
			u, err1 := strconv.Atoi(fields[0])
			v, err2 := strconv.Atoi(fields[1])
			if err1 != nil || err2 != nil {
				fmt.Fprintf(os.Stderr, "case %d failed: cannot parse edge %d\ninput:\n%s", t+1, i+1, input)
				os.Exit(1)
			}
			if u < 1 || u > n || v < 1 || v > n {
				fmt.Fprintf(os.Stderr, "case %d failed: edge %d out of range: %d %d\ninput:\n%s", t+1, i+1, u, v, input)
				os.Exit(1)
			}
			edges[i] = [2]int{u, v}
		}

		if !checkAnswer(n, a, edges) {
			fmt.Fprintf(os.Stderr, "case %d failed: graph not strongly connected after adding edges\ninput:\n%sgot:\n%s\n", t+1, input, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
