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

type testCase struct {
	n   int
	adj [][]int
	s   int // zero-based
}

type parsedOutput struct {
	kind string
	path []int
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-936B-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleB")
	cmd := exec.Command("go", "build", "-o", outPath, "936B.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
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
	m := 0
	for _, edges := range tc.adj {
		m += len(edges)
	}
	sb.Grow(tc.n*8 + m*4 + 32)
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte(' ')
	sb.WriteString(strconv.Itoa(m))
	sb.WriteByte('\n')
	for i := 0; i < tc.n; i++ {
		sb.WriteString(strconv.Itoa(len(tc.adj[i])))
		for _, to := range tc.adj[i] {
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(to + 1))
		}
		sb.WriteByte('\n')
	}
	sb.WriteString(strconv.Itoa(tc.s + 1))
	sb.WriteByte('\n')
	return sb.String()
}

func parseOutput(out string) (parsedOutput, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return parsedOutput{}, fmt.Errorf("empty output")
	}
	kind := fields[0]
	switch kind {
	case "Win":
		if len(fields) < 2 {
			return parsedOutput{}, fmt.Errorf("expected path after Win")
		}
		path := make([]int, len(fields)-1)
		for i, f := range fields[1:] {
			val, err := strconv.Atoi(f)
			if err != nil {
				return parsedOutput{}, fmt.Errorf("invalid vertex %q: %v", f, err)
			}
			path[i] = val
		}
		return parsedOutput{kind: kind, path: path}, nil
	case "Draw", "Lose":
		if len(fields) != 1 {
			return parsedOutput{}, fmt.Errorf("%s output should be a single word, got %d tokens", kind, len(fields))
		}
		return parsedOutput{kind: kind}, nil
	default:
		return parsedOutput{}, fmt.Errorf("unknown verdict %q", kind)
	}
}

func verifyWinningPath(tc testCase, path []int) error {
	if len(path) == 0 {
		return fmt.Errorf("empty path for Win verdict")
	}
	if len(path) > 1_000_000 {
		return fmt.Errorf("path too long: %d vertices", len(path))
	}
	if path[0] != tc.s+1 {
		return fmt.Errorf("path must start at %d, got %d", tc.s+1, path[0])
	}
	for idx, v := range path {
		if v < 1 || v > tc.n {
			return fmt.Errorf("vertex %d at position %d out of range", v, idx+1)
		}
	}
	if (len(path)-1)%2 == 0 {
		return fmt.Errorf("number of moves must be odd, got %d", len(path)-1)
	}
	edgeSets := make([]map[int]struct{}, tc.n)
	for i := 0; i < tc.n; i++ {
		if len(tc.adj[i]) == 0 {
			continue
		}
		set := make(map[int]struct{}, len(tc.adj[i]))
		for _, to := range tc.adj[i] {
			set[to] = struct{}{}
		}
		edgeSets[i] = set
	}
	for i := 0; i+1 < len(path); i++ {
		u := path[i] - 1
		v := path[i+1] - 1
		set := edgeSets[u]
		if set == nil {
			return fmt.Errorf("no outgoing edges from vertex %d but path continues", path[i])
		}
		if _, ok := set[v]; !ok {
			return fmt.Errorf("edge %d -> %d does not exist", path[i], path[i+1])
		}
	}
	last := path[len(path)-1] - 1
	if len(tc.adj[last]) != 0 {
		return fmt.Errorf("last vertex %d is not terminal", last+1)
	}
	return nil
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 5,
			adj: [][]int{
				{1, 2},
				{3, 4},
				{3},
				{4},
				{},
			},
			s: 0,
		},
		{
			n: 3,
			adj: [][]int{
				{2},
				{0},
				{},
			},
			s: 1,
		},
		{
			n: 2,
			adj: [][]int{
				{1},
				{0},
			},
			s: 0,
		},
		{
			n: 2,
			adj: [][]int{
				{1},
				{},
			},
			s: 1,
		},
		{
			n: 4,
			adj: [][]int{
				{1},
				{2},
				{3},
				{},
			},
			s: 0,
		},
		{
			n: 3,
			adj: [][]int{
				{1},
				{2},
				{},
			},
			s: 0,
		},
		{
			n: 3,
			adj: [][]int{
				{1},
				{2},
				{1},
			},
			s: 0,
		},
		{
			n: 6,
			adj: [][]int{
				{1, 2},
				{3},
				{3, 4},
				{5},
				{5},
				{},
			},
			s: 0,
		},
		{
			n: 6,
			adj: [][]int{
				{1, 2},
				{3, 4},
				{4},
				{0},
				{5},
				{4},
			},
			s: 0,
		},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 200)
	for len(tests) < cap(tests) {
		n := rng.Intn(40) + 2
		adj := make([][]int, n)
		for i := 0; i < n; i++ {
			maxDeg := n - 1
			limit := 5
			if rng.Intn(4) == 0 {
				limit = 10
			}
			if limit > maxDeg {
				limit = maxDeg
			}
			deg := 0
			if limit > 0 {
				deg = rng.Intn(limit + 1)
			}
			if deg == 0 {
				continue
			}
			perm := rng.Perm(n)
			for _, to := range perm {
				if to == i {
					continue
				}
				adj[i] = append(adj[i], to)
				if len(adj[i]) == deg {
					break
				}
			}
		}
		s := rng.Intn(n)
		tests = append(tests, testCase{n: n, adj: adj, s: s})
	}
	// add some larger structured tests
	tests = append(tests,
		buildChainTest(60),
		buildCycleTest(50),
		buildAlternatingTest(80),
	)
	return tests
}

func buildChainTest(length int) testCase {
	n := length
	if n%2 == 1 {
		n++
	}
	adj := make([][]int, n)
	for i := 0; i+1 < n; i++ {
		adj[i] = []int{i + 1}
	}
	adj[n-1] = nil
	return testCase{n: n, adj: adj, s: 0}
}

func buildCycleTest(n int) testCase {
	adj := make([][]int, n)
	for i := 0; i < n; i++ {
		adj[i] = []int{(i + 1) % n}
	}
	return testCase{n: n, adj: adj, s: 0}
}

func buildAlternatingTest(n int) testCase {
	adj := make([][]int, n)
	for i := 0; i < n; i++ {
		if i+1 < n {
			adj[i] = append(adj[i], i+1)
		}
		if i+2 < n {
			adj[i] = append(adj[i], i+2)
		}
	}
	adj[n-1] = nil
	return testCase{n: n, adj: adj, s: 0}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)
	for idx, tc := range tests {
		input := buildInput(tc)
		expectedOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		actualOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		expectedParsed, err := parseOutput(expectedOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle produced invalid output on test %d: %v\noutput:\n%s", idx+1, err, expectedOut)
			os.Exit(1)
		}
		actualParsed, err := parseOutput(actualOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target produced invalid output on test %d: %v\noutput:\n%s", idx+1, err, actualOut)
			os.Exit(1)
		}
		if expectedParsed.kind != actualParsed.kind {
			fmt.Fprintf(os.Stderr, "test %d verdict mismatch: expected %s, got %s\ninput:\n%s", idx+1, expectedParsed.kind, actualParsed.kind, input)
			os.Exit(1)
		}
		if expectedParsed.kind == "Win" {
			if err := verifyWinningPath(tc, actualParsed.path); err != nil {
				fmt.Fprintf(os.Stderr, "test %d invalid winning path: %v\ninput:\n%s", idx+1, err, input)
				os.Exit(1)
			}
		}
	}

	fmt.Println("All tests passed.")
}
