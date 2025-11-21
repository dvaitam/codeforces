package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	refSource        = "600F.go"
	tempOraclePrefix = "oracle-600F-"
	randomTestsCount = 120
	maxRandomA       = 40
	maxRandomB       = 40
)

type edge struct {
	u int
	v int
}

type testCase struct {
	name      string
	a, b      int
	edges     []edge
	maxDegree int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	oraclePath, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTests(randomTestsCount, rng)...)

	for idx, tc := range tests {
		input := formatInput(tc)

		if _, err := runProgram(oraclePath, input); err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}

		colorCount, colors, err := parseAnswer(candOut, len(tc.edges))
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output parse error on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, candOut)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}

		if err := validateColors(tc, colorCount, colors); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\n", idx+1, tc.name, err)
			fmt.Println("Input:")
			fmt.Print(input)
			fmt.Println("Candidate output:")
			fmt.Print(candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("failed to determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", tempOraclePrefix)
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleF")
	cmd := exec.Command("go", "build", "-o", outPath, refSource)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", tc.a, tc.b, len(tc.edges))
	for _, e := range tc.edges {
		fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
	}
	return sb.String()
}

func parseAnswer(out string, m int) (int, []int, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, nil, fmt.Errorf("empty output")
	}
	c, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, nil, fmt.Errorf("invalid colour count %q", fields[0])
	}
	if c < 0 {
		return 0, nil, fmt.Errorf("colour count cannot be negative")
	}
	if len(fields)-1 != m {
		return 0, nil, fmt.Errorf("expected %d colour entries, got %d", m, len(fields)-1)
	}
	colours := make([]int, m)
	for i := 0; i < m; i++ {
		val, err := strconv.Atoi(fields[i+1])
		if err != nil {
			return 0, nil, fmt.Errorf("colour %q is not an integer", fields[i+1])
		}
		colours[i] = val
	}
	return c, colours, nil
}

func validateColors(tc testCase, colourCount int, colours []int) error {
	m := len(tc.edges)
	if m == 0 {
		if colourCount != 0 {
			return fmt.Errorf("expected 0 colours for empty graph, got %d", colourCount)
		}
		return nil
	}
	if colourCount != tc.maxDegree {
		return fmt.Errorf("colour count %d differs from optimal %d", colourCount, tc.maxDegree)
	}
	for idx, col := range colours {
		if col < 1 || col > colourCount {
			return fmt.Errorf("edge %d has colour %d outside [1,%d]", idx+1, col, colourCount)
		}
	}
	leftSeen := make([]map[int]struct{}, tc.a+1)
	rightSeen := make([]map[int]struct{}, tc.b+1)
	for i, e := range tc.edges {
		col := colours[i]
		if leftSeen[e.u] == nil {
			leftSeen[e.u] = make(map[int]struct{})
		}
		if _, exists := leftSeen[e.u][col]; exists {
			return fmt.Errorf("vertex %d in first part has multiple edges with colour %d", e.u, col)
		}
		leftSeen[e.u][col] = struct{}{}

		if rightSeen[e.v] == nil {
			rightSeen[e.v] = make(map[int]struct{})
		}
		if _, exists := rightSeen[e.v][col]; exists {
			return fmt.Errorf("vertex %d in second part has multiple edges with colour %d", e.v, col)
		}
		rightSeen[e.v][col] = struct{}{}
	}
	return nil
}

func deterministicTests() []testCase {
	return []testCase{
		newTestCase("empty_graph", 3, 4, nil),
		newTestCase("single_edge", 1, 1, []edge{{1, 1}}),
		newTestCase("simple_path", 3, 3, []edge{{1, 1}, {2, 1}, {2, 2}, {3, 2}, {3, 3}}),
		newTestCase("star_left", 5, 5, []edge{{1, 1}, {1, 2}, {1, 3}, {1, 4}, {1, 5}}),
		newTestCase("complete_small", 3, 3, completeBipartiteEdges(3, 3)),
		newTestCase("balanced_dense", 6, 6, []edge{
			{1, 1}, {1, 2}, {2, 2}, {2, 3}, {3, 3}, {3, 4},
			{4, 4}, {4, 5}, {5, 5}, {5, 6}, {6, 1}, {6, 6},
		}),
	}
}

func randomTests(count int, rng *rand.Rand) []testCase {
	tests := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		a := rng.Intn(maxRandomA-1) + 1
		b := rng.Intn(maxRandomB-1) + 1
		maxEdges := a * b
		eCount := 0
		if maxEdges > 0 {
			eCount = rng.Intn(maxEdges + 1)
		}
		edgeSet := make(map[int]struct{})
		edges := make([]edge, 0, eCount)
		for len(edges) < eCount {
			u := rng.Intn(a) + 1
			v := rng.Intn(b) + 1
			key := u*(b+1) + v
			if _, exists := edgeSet[key]; exists {
				continue
			}
			edgeSet[key] = struct{}{}
			edges = append(edges, edge{u: u, v: v})
		}
		tests = append(tests, newTestCase(
			fmt.Sprintf("random_%d", i+1),
			a, b, edges,
		))
	}
	return tests
}

func newTestCase(name string, a, b int, edges []edge) testCase {
	degLeft := make([]int, a+1)
	degRight := make([]int, b+1)
	maxDeg := 0
	for _, e := range edges {
		if e.u < 1 || e.u > a || e.v < 1 || e.v > b {
			panic("edge endpoints out of bounds")
		}
		degLeft[e.u]++
		degRight[e.v]++
		if degLeft[e.u] > maxDeg {
			maxDeg = degLeft[e.u]
		}
		if degRight[e.v] > maxDeg {
			maxDeg = degRight[e.v]
		}
	}
	return testCase{
		name:      name,
		a:         a,
		b:         b,
		edges:     append([]edge(nil), edges...),
		maxDegree: maxDeg,
	}
}

func completeBipartiteEdges(a, b int) []edge {
	edges := make([]edge, 0, a*b)
	for u := 1; u <= a; u++ {
		for v := 1; v <= b; v++ {
			edges = append(edges, edge{u: u, v: v})
		}
	}
	return edges
}
