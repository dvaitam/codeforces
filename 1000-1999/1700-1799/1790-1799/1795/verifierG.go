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

// Brute-force solver for 1795G (small n).
// Build graph, find all valid removal orderings (permutations where each
// vertex is removed when its current degree equals a[i]).
// Count nice pairs: pairs (x,y) with x<y where x appears before y in some
// valid ordering and y appears before x in another.

func solveCase(n, m int, a []int, edgeList [][2]int) int64 {
	adj := make([][]int, n+1)
	for _, e := range edgeList {
		adj[e[0]] = append(adj[e[0]], e[1])
		adj[e[1]] = append(adj[e[1]], e[0])
	}

	// Find all valid removal sequences via backtracking
	removed := make([]bool, n+1)
	deg := make([]int, n+1)
	for _, e := range edgeList {
		deg[e[0]]++
		deg[e[1]]++
	}

	// For each pair, track if we've seen x before y and y before x
	type pairInfo struct {
		ab bool // x before y in some ordering
		ba bool // y before x in some ordering
	}
	pairs := make(map[[2]int]*pairInfo)
	for x := 1; x <= n; x++ {
		for y := x + 1; y <= n; y++ {
			pairs[[2]int{x, y}] = &pairInfo{}
		}
	}

	order := make([]int, 0, n)

	var dfs func()
	dfs = func() {
		if len(order) == n {
			// Record ordering
			pos := make([]int, n+1)
			for i, v := range order {
				pos[v] = i
			}
			for x := 1; x <= n; x++ {
				for y := x + 1; y <= n; y++ {
					p := pairs[[2]int{x, y}]
					if pos[x] < pos[y] {
						p.ab = true
					} else {
						p.ba = true
					}
				}
			}
			return
		}
		for v := 1; v <= n; v++ {
			if removed[v] {
				continue
			}
			// Current degree of v among non-removed neighbors
			curDeg := 0
			for _, u := range adj[v] {
				if !removed[u] {
					curDeg++
				}
			}
			if curDeg != a[v] {
				continue
			}
			// Remove v
			removed[v] = true
			deg[v] = -1
			order = append(order, v)
			dfs()
			order = order[:len(order)-1]
			removed[v] = false
		}
	}
	dfs()

	count := int64(0)
	for _, p := range pairs {
		if p.ab && p.ba {
			count++
		}
	}
	return count
}

func genTest(rng *rand.Rand) (string, int, int, []int, [][2]int) {
	n := rng.Intn(5) + 1
	maxM := n * (n - 1) / 2
	m := rng.Intn(maxM + 1)

	// Generate random graph with valid removal sequence
	// First generate edges
	edgeSet := make(map[[2]int]bool)
	edges := make([][2]int, 0, m)
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		if edgeSet[[2]int{u, v}] {
			continue
		}
		edgeSet[[2]int{u, v}] = true
		edges = append(edges, [2]int{u, v})
	}

	// Compute a valid removal sequence to determine valid a[] values
	deg := make([]int, n+1)
	adjList := make([][]int, n+1)
	for _, e := range edges {
		deg[e[0]]++
		deg[e[1]]++
		adjList[e[0]] = append(adjList[e[0]], e[1])
		adjList[e[1]] = append(adjList[e[1]], e[0])
	}

	// Choose a random valid removal order
	removed := make([]bool, n+1)
	a := make([]int, n+1)
	perm := rng.Perm(n)
	for idx := 0; idx < n; idx++ {
		// Find vertices that can be removed (degree == some value we'll set)
		// Strategy: pick a random non-removed vertex, set a[v] = current degree
		v := perm[idx] + 1
		// But we need to find one that we haven't removed yet
		// Simple: just remove in order of perm, setting a[v] to current degree
		if removed[v] {
			// This shouldn't happen with a permutation
			continue
		}
		curDeg := 0
		for _, u := range adjList[v] {
			if !removed[u] {
				curDeg++
			}
		}
		a[v] = curDeg
		removed[v] = true
	}

	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", a[i]))
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String(), n, len(edges), a[1:], edges
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for tc := 0; tc < 200; tc++ {
		input, n, m, aSlice, edges := genTest(rng)
		a := make([]int, n+1)
		for i := 0; i < n; i++ {
			a[i+1] = aSlice[i]
		}
		expectedVal := solveCase(n, m, a, edges)
		expectedOut := strconv.FormatInt(expectedVal, 10)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", tc+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expectedOut {
			fmt.Printf("case %d failed: expected %s got %s\ninput:\n%s", tc+1, expectedOut, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
