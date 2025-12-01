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

// refSource points to the local reference solution to avoid GOPATH resolution.
const refSource = "2134D.go"

type edge struct{ u, v int }

type testCase struct {
	n      int
	edges  []edge
	name   string
	isPath bool
	adj    [][]int
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := args[0]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := buildInput(tests)

	// Run reference only to ensure it works; we won't use its exact output for validation
	if _, err := runProgram(refBin, input); err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	if err := validateCandidate(candOut, tests); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	outPath := "./ref_2134D.bin"
	cmd := exec.Command("go", "build", "-o", outPath, refSource)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return outPath, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func buildTests() []testCase {
	var tests []testCase
	add := func(name string, n int, edges []edge) {
		tc := testCase{name: name, n: n, edges: edges}
		tc.buildAdj()
		tests = append(tests, tc)
	}

	// Simple cases
	add("single", 1, nil)
	add("path2", 2, []edge{{1, 2}})
	add("path4", 4, []edge{{1, 2}, {2, 3}, {3, 4}})
	add("star4", 4, []edge{{1, 2}, {1, 3}, {1, 4}})
	add("star5", 5, []edge{{2, 1}, {2, 3}, {2, 4}, {2, 5}})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	totalN := 0
	const maxTotalN = 180000
	for len(tests) < 180 && totalN < maxTotalN {
		n := rng.Intn(200_000) + 1
		if totalN+n > maxTotalN {
			n = maxTotalN - totalN
			if n < 1 {
				break
			}
		}
		var edges []edge
		if rng.Intn(4) == 0 && n >= 4 {
			// Force a high-degree center to ensure non-path
			center := rng.Intn(n) + 1
			edges = make([]edge, 0, n-1)
			for v := 1; v <= n; v++ {
				if v == center {
					continue
				}
				edges = append(edges, edge{center, v})
			}
		} else {
			edges = make([]edge, 0, n-1)
			for v := 2; v <= n; v++ {
				p := rng.Intn(v-1) + 1
				edges = append(edges, edge{v, p})
			}
		}
		add(fmt.Sprintf("random_%d", len(tests)), n, edges)
		totalN += n
	}

	return tests
}

func (tc *testCase) buildAdj() {
	tc.adj = make([][]int, tc.n+1)
	for _, e := range tc.edges {
		tc.adj[e.u] = append(tc.adj[e.u], e.v)
		tc.adj[e.v] = append(tc.adj[e.v], e.u)
	}
	tc.isPath = true
	for i := 1; i <= tc.n; i++ {
		if len(tc.adj[i]) > 2 {
			tc.isPath = false
			return
		}
	}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for _, e := range tc.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
		}
	}
	return sb.String()
}

func validateCandidate(out string, tests []testCase) error {
	tokens := strings.Fields(out)
	pos := 0
	for i, tc := range tests {
		if pos >= len(tokens) {
			return fmt.Errorf("candidate output ended early at case %d", i+1)
		}
		if tokens[pos] == "-1" {
			pos++
			if !tc.isPath {
				return fmt.Errorf("case %d (%s): tree not a path but candidate output -1", i+1, tc.name)
			}
			continue
		}
		if tc.isPath {
			return fmt.Errorf("case %d (%s): tree is already a path, expected -1", i+1, tc.name)
		}
		if pos+3 > len(tokens) {
			return fmt.Errorf("case %d (%s): expected three integers, got fewer", i+1, tc.name)
		}
		a, err1 := strconv.Atoi(tokens[pos])
		b, err2 := strconv.Atoi(tokens[pos+1])
		c, err3 := strconv.Atoi(tokens[pos+2])
		pos += 3
		if err1 != nil || err2 != nil || err3 != nil {
			return fmt.Errorf("case %d (%s): non-integer tokens in operation", i+1, tc.name)
		}
		if a == b || b == c || a == c {
			return fmt.Errorf("case %d (%s): vertices must be distinct", i+1, tc.name)
		}
		if a < 1 || a > tc.n || b < 1 || b > tc.n || c < 1 || c > tc.n {
			return fmt.Errorf("case %d (%s): vertex out of range", i+1, tc.name)
		}
		if !connected(tc.adj, a, b) || !connected(tc.adj, b, c) {
			return fmt.Errorf("case %d (%s): b must be adjacent to both a and c", i+1, tc.name)
		}
		if len(tc.adj[b]) < 3 {
			return fmt.Errorf("case %d (%s): chosen center b=%d has degree %d, need at least 3", i+1, tc.name, b, len(tc.adj[b]))
		}
	}
	if pos != len(tokens) {
		return fmt.Errorf("candidate output has %d extra tokens", len(tokens)-pos)
	}
	return nil
}

func connected(adj [][]int, u, v int) bool {
	for _, x := range adj[u] {
		if x == v {
			return true
		}
	}
	return false
}
