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

// computeMinEdges computes the minimum number of edges to add so that the
// functional graph f becomes strongly connected. This mirrors the logic of the
// reference solver at /tmp/cf_t24_22_E.go:
//   1. Walk from every leaf (in-degree 0) to discover which weakly-connected
//      component it belongs to, recording each leaf.
//   2. Walk from remaining unvisited nodes (pure-cycle components) and record
//      the cycle root as a pseudo-leaf for that component.
//   3. If there is exactly one component with no real leaves, answer is 0.
//   4. Otherwise, for each component we need one edge to chain the cycles
//      together, plus one extra edge per additional leaf in that component.
//      Total = k + sum_i(max(0, len(leavesOfComp[i]) - 1))   where k = #components.
func computeMinEdges(n int, f []int) int {
	inDegree := make([]int, n+1)
	for i := 1; i <= n; i++ {
		inDegree[f[i]]++
	}

	compID := make([]int, n+1)
	var leavesOfComp [][]int
	currentComp := 0

	// Phase 1: walk from every leaf (in-degree 0 node).
	for i := 1; i <= n; i++ {
		if inDegree[i] == 0 {
			curr := i
			var path []int
			for compID[curr] == 0 {
				compID[curr] = -1
				path = append(path, curr)
				curr = f[curr]
			}
			c := compID[curr]
			if c == -1 {
				// We reached a cycle node that was visited in this walk but
				// not yet assigned a component.
				currentComp++
				c = currentComp
				leavesOfComp = append(leavesOfComp, []int{})
			}
			for _, v := range path {
				compID[v] = c
			}
			leavesOfComp[c-1] = append(leavesOfComp[c-1], i)
		}
	}

	// Phase 2: pure-cycle components (no leaves).
	for i := 1; i <= n; i++ {
		if compID[i] == 0 {
			curr := i
			var path []int
			for compID[curr] == 0 {
				compID[curr] = -1
				path = append(path, curr)
				curr = f[curr]
			}
			c := compID[curr]
			if c == -1 {
				currentComp++
				c = currentComp
				// For a pure-cycle component, add the cycle root as a
				// pseudo-leaf so it participates in the chaining.
				leavesOfComp = append(leavesOfComp, []int{curr})
			}
			for _, v := range path {
				compID[v] = c
			}
		}
	}

	numLeaves := 0
	for i := 1; i <= n; i++ {
		if inDegree[i] == 0 {
			numLeaves++
		}
	}

	if currentComp == 1 && numLeaves == 0 {
		return 0
	}

	// Count edges: one per component (to chain cycles) plus extras for leaves.
	k := currentComp
	total := 0
	for i := 0; i < k; i++ {
		// One edge to link cycle i -> leavesOfComp[(i+1)%k][0]
		total++
		// Extra edges for remaining leaves in component i
		if len(leavesOfComp[i]) > 1 {
			total += len(leavesOfComp[i]) - 1
		}
	}
	return total
}

// checkAnswer validates that after adding the given edges, the graph becomes
// strongly connected (everyone can reach everyone).
func checkAnswer(n int, f []int, edges [][2]int) bool {
	adj := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		adj[i] = append(adj[i], f[i])
	}
	for _, e := range edges {
		adj[e[0]] = append(adj[e[0]], e[1])
	}

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
