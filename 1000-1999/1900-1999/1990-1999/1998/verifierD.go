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

const refSource = "1000-1999/1900-1999/1990-1999/1998/1998D.go"

type testCase struct {
	n     int
	m     int
	edges [][2]int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]

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
		fmt.Fprintf(os.Stderr, "reference failed: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	expected, err := parseOutputs(refOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse reference output: %v\n", err)
		os.Exit(1)
	}
	got, err := parseOutputs(candOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i := range tests {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %s got %s\n", i+1, expected[i], got[i])
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1998D-ref-*")
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

func parseOutputs(output string, tests []testCase) ([]string, error) {
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	if len(lines) < len(tests) {
		return nil, fmt.Errorf("expected %d lines, got %d", len(tests), len(lines))
	}
	if len(lines) > len(tests) {
		return nil, fmt.Errorf("extra output detected starting at line %d", len(tests)+1)
	}
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) != tests[i].n-1 {
			return nil, fmt.Errorf("test %d: output length %d != %d", i+1, len(line), tests[i].n-1)
		}
		for j := 0; j < len(line); j++ {
			if line[j] != '0' && line[j] != '1' {
				return nil, fmt.Errorf("test %d: output contains invalid character %q", i+1, line[j])
			}
		}
	}
	return lines, nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(1998))
	var tests []testCase
	totalN := 0
	totalM := 0
	add := func(tc testCase) {
		if totalN+tc.n > 200000 || totalM+tc.m > 200000 {
			return
		}
		tests = append(tests, tc)
		totalN += tc.n
		totalM += tc.m
	}

	add(makeCase(6, 0, nil))
	add(makeCase(6, 1, [][2]int{{2, 6}}))
	add(makeCase(10, 4, [][2]int{{1, 3}, {1, 6}, {2, 7}, {3, 8}}))

	for totalN < 200000 && totalM < 200000 {
		n := rng.Intn(200) + 2
		maxExtra := min(2*n, 200)
		m := rng.Intn(maxExtra + 1)
		edges := randomEdges(rng, n, m)
		add(makeCase(n, len(edges), edges))
		if len(tests) > 400 {
			break
		}
	}
	return tests
}

func makeCase(n, m int, edges [][2]int) testCase {
	cp := make([][2]int, 0, m)
	for _, e := range edges {
		if e[0] < e[1] {
			cp = append(cp, e)
		}
	}
	return testCase{n: n, m: len(cp), edges: cp}
}

func randomEdges(rng *rand.Rand, n, m int) [][2]int {
	maxPossible := n*(n-1)/2 - (n - 1)
	if m > maxPossible {
		m = maxPossible
	}
	set := make(map[[2]int]struct{})
	for len(set) < m {
		u := rng.Intn(n-1) + 1
		v := rng.Intn(n-u) + u + 1
		if v == u+1 {
			continue
		}
		key := [2]int{u, v}
		if _, ok := set[key]; !ok {
			set[key] = struct{}{}
		}
	}
	res := make([][2]int, 0, len(set))
	for e := range set {
		res = append(res, e)
	}
	return res
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d %d\n", tc.n, tc.m)
		for _, e := range tc.edges {
			fmt.Fprintf(&b, "%d %d\n", e[0], e[1])
		}
	}
	return b.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
