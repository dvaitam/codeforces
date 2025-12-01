package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSourceL = "./1666L.go"

type testCase struct {
	input string
	n     int
	s     int
	adj   map[int]map[int]bool
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierL.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if err := evaluateTest(tc, refOut, candOut); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\nInput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1666L-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSourceL))
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, buf.String())
	}
	return tmp.Name(), nil
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func evaluateTest(tc testCase, refOut, candOut string) error {
	refPossible, err := parseDecision(refOut)
	if err != nil {
		return fmt.Errorf("reference output invalid: %v", err)
	}
	candPossible, err := parseDecision(candOut)
	if err != nil {
		return fmt.Errorf("candidate output invalid: %v", err)
	}

	if !refPossible {
		if candPossible {
			return fmt.Errorf("candidate claims Possible but reference Impossible")
		}
		if err := ensureNoExtraTokens(candOut); err != nil {
			return err
		}
		return nil
	}

	if !candPossible {
		return fmt.Errorf("candidate claims Impossible but solution exists")
	}

	path1, path2, err := parsePaths(candOut)
	if err != nil {
		return err
	}
	return validatePaths(tc, path1, path2)
}

func parseDecision(out string) (bool, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return false, fmt.Errorf("empty output")
	}
	if strings.EqualFold(fields[0], "Possible") {
		return true, nil
	}
	if strings.EqualFold(fields[0], "Impossible") {
		return false, nil
	}
	return false, fmt.Errorf("first token must be Possible or Impossible, got %q", fields[0])
}

func ensureNoExtraTokens(out string) error {
	reader := strings.NewReader(out)
	var token string
	if _, err := fmt.Fscan(reader, &token); err != nil {
		return fmt.Errorf("failed to read decision: %v", err)
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return fmt.Errorf("unexpected extra token %q after Impossible", extra)
	} else if err != io.EOF {
		return fmt.Errorf("error parsing output: %v", err)
	}
	return nil
}

func parsePaths(out string) ([]int, []int, error) {
	reader := strings.NewReader(out)
	var decision string
	if _, err := fmt.Fscan(reader, &decision); err != nil {
		return nil, nil, fmt.Errorf("failed to read decision: %v", err)
	}
	if !strings.EqualFold(decision, "Possible") {
		return nil, nil, fmt.Errorf("expected Possible, got %q", decision)
	}
	var len1 int
	if _, err := fmt.Fscan(reader, &len1); err != nil {
		return nil, nil, fmt.Errorf("failed to read first path length: %v", err)
	}
	path1 := make([]int, len1)
	for i := 0; i < len1; i++ {
		if _, err := fmt.Fscan(reader, &path1[i]); err != nil {
			return nil, nil, fmt.Errorf("failed to read first path vertex %d: %v", i+1, err)
		}
	}
	var len2 int
	if _, err := fmt.Fscan(reader, &len2); err != nil {
		return nil, nil, fmt.Errorf("failed to read second path length: %v", err)
	}
	path2 := make([]int, len2)
	for i := 0; i < len2; i++ {
		if _, err := fmt.Fscan(reader, &path2[i]); err != nil {
			return nil, nil, fmt.Errorf("failed to read second path vertex %d: %v", i+1, err)
		}
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, nil, fmt.Errorf("unexpected extra token %q after paths", extra)
	} else if err != io.EOF {
		return nil, nil, fmt.Errorf("error parsing trailing data: %v", err)
	}
	return path1, path2, nil
}

func validatePaths(tc testCase, path1, path2 []int) error {
	if err := validateSinglePath(tc, path1); err != nil {
		return fmt.Errorf("first path invalid: %v", err)
	}
	if err := validateSinglePath(tc, path2); err != nil {
		return fmt.Errorf("second path invalid: %v", err)
	}
	if len(path1) < 2 || len(path2) < 2 {
		return fmt.Errorf("paths must contain at least two vertices")
	}
	t1 := path1[len(path1)-1]
	t2 := path2[len(path2)-1]
	if t1 != t2 {
		return fmt.Errorf("paths end at different vertices (%d vs %d)", t1, t2)
	}
	if t1 == tc.s {
		return fmt.Errorf("destination must differ from start")
	}
	if samePath(path1, path2) {
		return fmt.Errorf("paths must be different")
	}
	used := make(map[int]bool)
	for i := 1; i < len(path1)-1; i++ {
		used[path1[i]] = true
	}
	for i := 1; i < len(path2)-1; i++ {
		if used[path2[i]] {
			return fmt.Errorf("paths share intermediate vertex %d", path2[i])
		}
	}
	return nil
}

func validateSinglePath(tc testCase, path []int) error {
	if len(path) < 2 {
		return fmt.Errorf("path too short")
	}
	if path[0] != tc.s {
		return fmt.Errorf("path must start at %d", tc.s)
	}
	seen := make(map[int]bool)
	for i, v := range path {
		if v < 1 || v > tc.n {
			return fmt.Errorf("vertex %d out of range", v)
		}
		if seen[v] {
			return fmt.Errorf("vertex %d appears multiple times", v)
		}
		seen[v] = true
		if i > 0 {
			u := path[i-1]
			if !tc.hasEdge(u, v) {
				return fmt.Errorf("no edge %d -> %d", u, v)
			}
		}
	}
	return nil
}

func samePath(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func (tc testCase) hasEdge(u, v int) bool {
	if m, ok := tc.adj[u]; ok {
		return m[v]
	}
	return false
}

func buildTests() []testCase {
	raw := []string{
		"5 5 1\n1 2\n2 3\n1 4\n4 3\n3 5\n",
		"5 5 1\n1 2\n2 3\n3 4\n2 5\n5 4\n",
		"3 3 2\n1 2\n2 3\n3 1\n",
		"6 6 1\n1 2\n1 3\n2 4\n3 5\n4 6\n5 6\n",
	}
	tests := make([]testCase, 0, len(raw))
	for idx, input := range raw {
		tc, err := newTestCase(input)
		if err != nil {
			panic(fmt.Sprintf("failed to parse test %d: %v", idx+1, err))
		}
		tests = append(tests, tc)
	}
	return tests
}

func newTestCase(input string) (testCase, error) {
	reader := strings.NewReader(input)
	var n, m, s int
	if _, err := fmt.Fscan(reader, &n, &m, &s); err != nil {
		return testCase{}, err
	}
	adj := make(map[int]map[int]bool, n+1)
	for i := 0; i < m; i++ {
		var u, v int
		if _, err := fmt.Fscan(reader, &u, &v); err != nil {
			return testCase{}, err
		}
		if adj[u] == nil {
			adj[u] = make(map[int]bool)
		}
		adj[u][v] = true
	}
	return testCase{
		input: input,
		n:     n,
		s:     s,
		adj:   adj,
	}, nil
}
