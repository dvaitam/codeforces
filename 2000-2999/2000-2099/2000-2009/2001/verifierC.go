package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type edge struct {
	u, v int
}

type testCase struct {
	n     int
	edges []edge
}

const refSource = "2000-2999/2000-2099/2000-2009/2001/2001C.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}

	candidate := os.Args[1]
	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	if _, err := parseOutput(refOut, tests); err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	candEdges, err := parseOutput(candOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse candidate output: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	for i, tc := range tests {
		if err := validateAnswer(tc, candEdges[i]); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2001C-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(output string, tests []testCase) ([][]edge, error) {
	tokens := strings.Fields(output)
	idx := 0
	res := make([][]edge, len(tests))
	for ti, tc := range tests {
		for idx < len(tokens) && tokens[idx] != "!" {
			idx++
		}
		if idx >= len(tokens) {
			return nil, fmt.Errorf("test %d: missing '!'", ti+1)
		}
		idx++
		need := 2 * (tc.n - 1)
		if len(tokens)-idx < need {
			return nil, fmt.Errorf("test %d: insufficient edge data", ti+1)
		}
		edges := make([]edge, tc.n-1)
		for j := 0; j < tc.n-1; j++ {
			u, err1 := strconv.Atoi(tokens[idx])
			idx++
			v, err2 := strconv.Atoi(tokens[idx])
			idx++
			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("test %d: invalid integer in edge", ti+1)
			}
			edges[j] = edge{u: u, v: v}
		}
		res[ti] = edges
	}
	return res, nil
}

func validateAnswer(tc testCase, ans []edge) error {
	if len(ans) != tc.n-1 {
		return fmt.Errorf("expected %d edges, got %d", tc.n-1, len(ans))
	}

	expected := make(map[[2]int]int, len(tc.edges))
	for _, e := range tc.edges {
		key := canonicalEdge(e)
		expected[key]++
	}

	seen := make(map[[2]int]int, len(ans))
	for i, e := range ans {
		if e.u < 1 || e.u > tc.n || e.v < 1 || e.v > tc.n {
			return fmt.Errorf("edge %d has invalid vertex (%d,%d)", i+1, e.u, e.v)
		}
		if e.u == e.v {
			return fmt.Errorf("edge %d forms a self-loop at %d", i+1, e.u)
		}
		key := canonicalEdge(e)
		seen[key]++
		if seen[key] > expected[key] {
			return fmt.Errorf("edge %d (%d,%d) not present in the secret tree", i+1, e.u, e.v)
		}
	}

	for key, cnt := range expected {
		if seen[key] != cnt {
			return fmt.Errorf("missing edge %d-%d", key[0], key[1])
		}
	}

	return nil
}

func canonicalEdge(e edge) [2]int {
	if e.u > e.v {
		return [2]int{e.v, e.u}
	}
	return [2]int{e.u, e.v}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for _, e := range tc.edges {
			fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
		}
	}
	return sb.String()
}

func generateTests() []testCase {
	var tests []testCase
	totalNodes := 0

	add := func(tc testCase) {
		tests = append(tests, tc)
		totalNodes += tc.n
	}

	add(testCase{n: 4, edges: []edge{{1, 2}, {1, 3}, {3, 4}}})
	add(testCase{n: 2, edges: []edge{{1, 2}}})
	add(buildStar(6))
	add(buildChain(7))

	rng := rand.New(rand.NewSource(2001))
	for totalNodes < 950 {
		n := rng.Intn(30) + 2
		if totalNodes+n > 1000 {
			n = 1000 - totalNodes
		}
		if n < 2 {
			break
		}
		add(randomTree(n, rng))
	}

	return tests
}

func buildStar(n int) testCase {
	edges := make([]edge, n-1)
	for i := 2; i <= n; i++ {
		edges[i-2] = edge{1, i}
	}
	return testCase{n: n, edges: edges}
}

func buildChain(n int) testCase {
	edges := make([]edge, n-1)
	for i := 1; i < n; i++ {
		edges[i-1] = edge{i, i + 1}
	}
	return testCase{n: n, edges: edges}
}

func randomTree(n int, rng *rand.Rand) testCase {
	edges := make([]edge, 0, n-1)
	for v := 2; v <= n; v++ {
		parent := rng.Intn(v-1) + 1
		edges = append(edges, edge{parent, v})
	}
	return testCase{n: n, edges: edges}
}
