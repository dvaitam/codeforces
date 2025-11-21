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

type edge struct {
	u, v int
	seq  []int
}

type testCase struct {
	n, m  int
	edges []edge
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-331E1-")
	if err != nil {
		return "", nil, err
	}
	path := filepath.Join(tmpDir, "oracleE1")
	cmd := exec.Command("go", "build", "-o", path, "331E1.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return path, cleanup, nil
}

func runBinary(bin, input string) (string, error) {
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
	return stdout.String(), nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d %d", e.u, e.v, len(e.seq)))
		for _, v := range e.seq {
			sb.WriteString(fmt.Sprintf(" %d", v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string) (int, []int, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, nil, fmt.Errorf("empty output")
	}
	k, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, nil, fmt.Errorf("invalid k: %v", err)
	}
	if k == 0 {
		if len(fields) != 1 {
			return 0, nil, fmt.Errorf("unexpected tokens after 0")
		}
		return 0, nil, nil
	}
	if len(fields) != k+1 {
		return 0, nil, fmt.Errorf("expected %d nodes, got %d", k, len(fields)-1)
	}
	path := make([]int, k)
	for i := 0; i < k; i++ {
		val, err := strconv.Atoi(fields[i+1])
		if err != nil {
			return 0, nil, fmt.Errorf("invalid node: %v", err)
		}
		path[i] = val
	}
	return k, path, nil
}

func validatePath(tc testCase, path []int) error {
	n := tc.n
	if len(path) < 2 {
		return fmt.Errorf("path length < 2")
	}
	if len(path) > 2*n {
		return fmt.Errorf("path length exceeds 2n")
	}
	for _, v := range path {
		if v < 1 || v > n {
			return fmt.Errorf("node %d out of range", v)
		}
	}
	edgeMap := make(map[[2]int][]int)
	for _, e := range tc.edges {
		edgeMap[[2]int{e.u, e.v}] = e.seq
	}
	vision := make([]int, 0)
	for i := 0; i < len(path)-1; i++ {
		key := [2]int{path[i], path[i+1]}
		seq, ok := edgeMap[key]
		if !ok {
			return fmt.Errorf("edge %d -> %d does not exist", path[i], path[i+1])
		}
		vision = append(vision, seq...)
	}
	if len(vision) != len(path) {
		return fmt.Errorf("vision length %d does not match path length %d", len(vision), len(path))
	}
	for i := range vision {
		if vision[i] != path[i] {
			return fmt.Errorf("vision mismatch at position %d: vision %d actual %d", i, vision[i], path[i])
		}
	}
	return nil
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 6, m: 6,
			edges: []edge{
				{1, 2, []int{1, 2}},
				{2, 3, []int{3}},
				{3, 4, []int{4, 5}},
				{4, 5, []int{}},
				{5, 3, []int{3}},
				{6, 1, []int{6}},
			},
		},
		{
			n: 2, m: 1,
			edges: []edge{{1, 2, []int{1}}},
		},
	}
}

func randomTest() testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := rng.Intn(6) + 2
	maxEdges := n * (n - 1)
	m := rng.Intn(maxEdges) + 1
	used := make(map[[2]int]bool)
	edges := make([]edge, 0, m)
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		key := [2]int{u, v}
		if used[key] {
			continue
		}
		used[key] = true
		k := rng.Intn(n)
		seq := make([]int, k)
		for i := 0; i < k; i++ {
			seq[i] = rng.Intn(n) + 1
		}
		edges = append(edges, edge{u, v, seq})
	}
	return testCase{n: n, m: len(edges), edges: edges}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	for i := 0; i < 30; i++ {
		tests = append(tests, randomTest())
	}

	for idx, tc := range tests {
		input := buildInput(tc)
		expOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		expK, _, err := parseOutput(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\n%s", idx+1, err, expOut)
			os.Exit(1)
		}

		actOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		actK, actPath, err := parseOutput(actOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\n%s", idx+1, err, actOut)
			os.Exit(1)
		}
		if expK == 0 {
			if actK != 0 {
				if err := validatePath(tc, actPath); err != nil {
					fmt.Fprintf(os.Stderr, "test %d: invalid path: %v\ninput:\n%s\noutput:\n%s\n", idx+1, err, input, actOut)
					os.Exit(1)
				}
			}
			continue
		}
		if actK == 0 {
			fmt.Fprintf(os.Stderr, "test %d: expected path exists but got 0\ninput:\n%s\n", idx+1, input)
			os.Exit(1)
		}
		if err := validatePath(tc, actPath); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: invalid path: %v\ninput:\n%s\noutput:\n%s\n", idx+1, err, input, actOut)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed.")
}
