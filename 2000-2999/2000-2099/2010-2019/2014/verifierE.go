package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const referenceSolutionRel = "2000-2999/2000-2099/2010-2019/2014/2014E.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "2014E.go")
		if _, err := os.Stat(candidate); err == nil {
			referenceSolutionPath = candidate
			return
		}
	}
	if abs, err := filepath.Abs(referenceSolutionRel); err == nil {
		if _, err := os.Stat(abs); err == nil {
			referenceSolutionPath = abs
		}
	}
}

type edge struct {
	u, v int
	w    int64
}

type testCase struct {
	n, m, h int
	horse   []int
	edges   []edge
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 2, m: 1, h: 1,
			horse: []int{1},
			edges: []edge{{1, 2, 10}},
		},
		{
			n: 3, m: 1, h: 2,
			horse: []int{1, 2},
			edges: []edge{{1, 3, 10}},
		},
		{
			n: 3, m: 3, h: 2,
			horse: []int{2, 3},
			edges: []edge{{1, 2, 4}, {1, 3, 10}, {2, 3, 6}},
		},
		{
			n: 4, m: 3, h: 2,
			horse: []int{2, 3},
			edges: []edge{{1, 2, 10}, {2, 3, 18}, {3, 4, 16}},
		},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(20250305))
	var tests []testCase
	totalN := 0
	totalM := 0
	for len(tests) < 80 && totalN < 200000 && totalM < 200000 {
		n := rng.Intn(30) + 2
		m := rng.Intn(n*(n-1)/2) + 1
		if totalN+n > 200000 {
			n = 200000 - totalN
		}
		if totalM+m > 200000 {
			m = 200000 - totalM
		}
		if n < 2 || m <= 0 {
			break
		}
		edgeSet := make(map[[2]int]struct{})
		edges := make([]edge, 0, m)
		for len(edges) < m {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			if u == v {
				continue
			}
			if u > v {
				u, v = v, u
			}
			key := [2]int{u, v}
			if _, ok := edgeSet[key]; ok {
				continue
			}
			edgeSet[key] = struct{}{}
			weight := int64(rng.Intn(500)+1) * 2
			edges = append(edges, edge{u, v, weight})
		}
		h := rng.Intn(n) + 1
		horses := make([]int, 0, h)
		perm := rng.Perm(n)
		for i := 0; i < h; i++ {
			horses = append(horses, perm[i]+1)
		}
		tests = append(tests, testCase{
			n:     n,
			m:     len(edges),
			h:     len(horses),
			horse: horses,
			edges: edges,
		})
		totalN += n
		totalM += len(edges)
	}
	return tests
}

func formatInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.h))
		for i, v := range tc.horse {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		for _, e := range tc.edges {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", e.u, e.v, e.w))
		}
	}
	return sb.String()
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildReferenceBinary() (string, func(), error) {
	if referenceSolutionPath == "" {
		return "", nil, fmt.Errorf("reference solution path not set")
	}
	if _, err := os.Stat(referenceSolutionPath); err != nil {
		return "", nil, fmt.Errorf("reference solution not found: %v", err)
	}
	tmpDir, err := os.MkdirTemp("", "2014E-ref")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_2014E")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSolutionPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return binPath, cleanup, nil
}

func parseOutputs(out string, count int) ([]int64, error) {
	scanner := bufio.NewScanner(strings.NewReader(out))
	scanner.Split(bufio.ScanWords)
	results := make([]int64, 0, count)
	for scanner.Scan() {
		var val int64
		if _, err := fmt.Sscan(scanner.Text(), &val); err != nil {
			return nil, fmt.Errorf("failed to parse %q: %v", scanner.Text(), err)
		}
		results = append(results, val)
	}
	if len(results) != count {
		return nil, fmt.Errorf("expected %d numbers, got %d", count, len(results))
	}
	return results, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := append(deterministicTests(), randomTests()...)
	input := formatInput(tests)

	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	expected, err := parseOutputs(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	userOut, err := runProgram(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}
	got, err := parseOutputs(userOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse participant output: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}

	for i := range tests {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %d, got %d\nn=%d m=%d h=%d\nhorse=%v\n", i+1, expected[i], got[i], tests[i].n, tests[i].m, tests[i].h, tests[i].horse)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
