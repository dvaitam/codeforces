package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type edge struct {
	i, j, b int
}

type testCase struct {
	input string
	n, m  int
	edges []edge
}

func buildRef() (string, error) {
	srcPath := os.Getenv("REFERENCE_SOURCE_PATH")
	if srcPath == "" {
		return "", fmt.Errorf("REFERENCE_SOURCE_PATH environment variable not set")
	}
	dir := filepath.Dir(srcPath)
	base := filepath.Base(srcPath)
	bin := filepath.Join(os.TempDir(), "1450E_ref.bin")
	cmd := exec.Command("go", "build", "-o", bin, base)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, out)
	}
	return bin, nil
}

func runProg(bin string, input string) (string, error) {
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

// generateConnectedGraph generates a random connected graph with n vertices and m edges
// where n-1 <= m and the graph is connected.
func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(6) + 2 // 2..7
	// ensure m >= n-1 for connectivity, and m <= min(2000, n*(n-1)/2)
	maxEdges := n * (n - 1) / 2
	if maxEdges > 20 {
		maxEdges = 20
	}
	minEdges := n - 1
	m := minEdges + rng.Intn(maxEdges-minEdges+1)

	// Build a random spanning tree first
	perm := rng.Perm(n)
	type rawEdge struct {
		u, v int
	}
	edgeSet := make(map[rawEdge]bool)
	var edges []edge

	for i := 1; i < n; i++ {
		u := perm[rng.Intn(i)]
		v := perm[i]
		a, b := u+1, v+1
		if a > b {
			a, b = b, a
		}
		edgeSet[rawEdge{a, b}] = true
		dir := rng.Intn(2)
		if dir == 0 {
			edges = append(edges, edge{a, b, rng.Intn(2)})
		} else {
			edges = append(edges, edge{b, a, rng.Intn(2)})
		}
	}

	// Add random extra edges
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		a, b := u, v
		if a > b {
			a, b = b, a
		}
		if edgeSet[rawEdge{a, b}] {
			continue
		}
		edgeSet[rawEdge{a, b}] = true
		dir := rng.Intn(2)
		if dir == 0 {
			edges = append(edges, edge{a, b, rng.Intn(2)})
		} else {
			edges = append(edges, edge{b, a, rng.Intn(2)})
		}
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, len(edges))
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d %d\n", e.i, e.j, e.b)
	}
	return testCase{input: sb.String(), n: n, m: len(edges), edges: edges}
}

// validateOutput checks if the candidate output is valid for the given test case.
// Returns nil if valid, error otherwise.
func validateOutput(tc testCase, refOutput, candOutput string) error {
	refLines := strings.Fields(refOutput)
	candLines := strings.Fields(candOutput)

	if len(candLines) == 0 {
		return fmt.Errorf("empty output")
	}

	refAnswer := strings.ToUpper(refLines[0])
	candAnswer := strings.ToUpper(candLines[0])

	if refAnswer == "NO" {
		if candAnswer != "NO" {
			return fmt.Errorf("expected NO, got %s", candAnswer)
		}
		return nil
	}

	// Reference says YES
	if candAnswer == "NO" {
		return fmt.Errorf("expected YES but got NO")
	}
	if candAnswer != "YES" {
		return fmt.Errorf("expected YES, got %s", candAnswer)
	}

	// Parse reference inequality
	if len(refLines) < 2 {
		return fmt.Errorf("reference output incomplete")
	}
	refInequality, err := strconv.Atoi(refLines[1])
	if err != nil {
		return fmt.Errorf("failed to parse reference inequality: %v", err)
	}

	// Parse candidate inequality and incomes
	if len(candLines) < 2+tc.n {
		return fmt.Errorf("candidate output too short: need inequality + %d incomes, got %d tokens", tc.n, len(candLines)-1)
	}
	candInequality, err := strconv.Atoi(candLines[1])
	if err != nil {
		return fmt.Errorf("failed to parse candidate inequality: %v", err)
	}

	if candInequality != refInequality {
		return fmt.Errorf("inequality mismatch: expected %d, got %d", refInequality, candInequality)
	}

	incomes := make([]int, tc.n+1)
	for i := 1; i <= tc.n; i++ {
		val, err := strconv.Atoi(candLines[1+i])
		if err != nil {
			return fmt.Errorf("failed to parse income %d: %v", i, err)
		}
		if val < 0 || val > 1000000 {
			return fmt.Errorf("income %d = %d out of range [0, 10^6]", i, val)
		}
		incomes[i] = val
	}

	// Check all edges: for each edge (i,j), |a_i - a_j| must be 1
	// If b=1, person i is envious of j, meaning a_j = a_i + 1
	for _, e := range tc.edges {
		diff := incomes[e.i] - incomes[e.j]
		if diff != 1 && diff != -1 {
			return fmt.Errorf("edge (%d,%d): incomes %d and %d differ by %d, not 1", e.i, e.j, incomes[e.i], incomes[e.j], int(math.Abs(float64(diff))))
		}
		if e.b == 1 {
			// i is envious of j: a_j = a_i + 1
			if incomes[e.j] != incomes[e.i]+1 {
				return fmt.Errorf("edge (%d,%d,b=1): expected a[%d]=a[%d]+1 but got %d and %d", e.i, e.j, e.j, e.i, incomes[e.j], incomes[e.i])
			}
		}
	}

	// Check inequality value
	minInc, maxInc := incomes[1], incomes[1]
	for i := 2; i <= tc.n; i++ {
		if incomes[i] < minInc {
			minInc = incomes[i]
		}
		if incomes[i] > maxInc {
			maxInc = incomes[i]
		}
	}
	if maxInc-minInc != candInequality {
		return fmt.Errorf("declared inequality %d but actual max-min = %d", candInequality, maxInc-minInc)
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}

	for i, tc := range cases {
		refOut, err := runProg(ref, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		candOut, err := runProg(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if err := validateOutput(tc, refOut, candOut); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%sref output:\n%s\ncand output:\n%s\n", i+1, err, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
