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

type testCase struct {
	input string
}

// Generates a random test case where a solution is guaranteed to exist.
func generateRandomCase(rng *rand.Rand) testCase {
	for {
		n := rng.Intn(10) + 3 // n in [3, 12]
		maxEdges := n * (n - 1) / 2
		if maxEdges < n-1 {
			continue // Should not happen for n>=3
		}

		// Ensure m is at least n-1 to allow for a connected graph
		m := rng.Intn(maxEdges-(n-1)+1) + (n - 1)

		edges := make([][3]int, 0, m)
		used := make(map[[2]int]bool)

		// 1. Create a connected graph by generating a random spanning tree.
		// A simple way is to create a random permutation of nodes and link them in a path.
		p := rand.Perm(n)
		for i := 0; i < n-1; i++ {
			u, v := p[i]+1, p[i+1]+1
			if u > v {
				u, v = v, u
			}
			if !used[[2]int{u, v}] {
				used[[2]int{u, v}] = true
				// By making all weights even, we satisfy the condition that the sum of weights for internal nodes is even.
				w := 2 * (rng.Intn(50) + 1)
				edges = append(edges, [3]int{u, v, w})
			}
		}

		// 2. Add more random edges until we have m edges.
		for len(edges) < m {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			if u == v {
				continue
			}
			if u > v {
				u, v = v, u
			}
			if used[[2]int{u, v}] {
				continue
			}
			used[[2]int{u, v}] = true
			w := 2 * (rng.Intn(50) + 1) // All weights are even
			edges = append(edges, [3]int{u, v, w})
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		rng.Shuffle(len(edges), func(i, j int) {
			edges[i], edges[j] = edges[j], edges[i]
		})
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2]))
		}
		return testCase{input: sb.String()}
	}
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

func verify(input, output string) error {
	scanner := strings.NewReader(input)
	var n, m int
	fmt.Fscan(scanner, &n, &m)

	edges := make([][3]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(scanner, &edges[i][0], &edges[i][1], &edges[i][2])
	}

	outFields := strings.Fields(output)
	if len(outFields) != m {
		return fmt.Errorf("incorrect number of orientations: need %d got %d", m, len(outFields))
	}

	orientations := make([]int, m)
	for i := 0; i < m; i++ {
		if _, err := fmt.Sscan(outFields[i], &orientations[i]); err != nil {
			return fmt.Errorf("bad orientation at index %d: %q", i, outFields[i])
		}
		if orientations[i] != 0 && orientations[i] != 1 {
			return fmt.Errorf("orientation must be 0 or 1, but got %d at index %d", orientations[i], i)
		}
	}

	outSum := make([]int64, n+1)
	inSum := make([]int64, n+1)

	for i := 0; i < m; i++ {
		u, v, w := edges[i][0], edges[i][1], int64(edges[i][2])

		if orientations[i] == 0 { // u -> v
			outSum[u] += w
			inSum[v] += w
		} else { // v -> u
			outSum[v] += w
			inSum[u] += w
		}
	}

	// For internal nodes (2 to n-1), in-flow must equal out-flow.
	for v := 2; v <= n-1; v++ {
		if inSum[v] != outSum[v] {
			return fmt.Errorf("flow conservation failed at node %d: in=%d, out=%d", v, inSum[v], outSum[v])
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateRandomCase(rng)
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if err := verify(tc.input, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\noutput:\n%s\n", i+1, err, tc.input, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
