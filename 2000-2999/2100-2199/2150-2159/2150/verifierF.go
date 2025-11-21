package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSource = "2000-2999/2100-2199/2150-2159/2150/2150F.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for i, input := range tests {
		wantOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		want, err := parseOutput(input, wantOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\ninput:\n%soutput:\n%s", i+1, err, input, wantOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := parseOutput(input, gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d: %v\ninput:\n%soutput:\n%s", i+1, err, input, gotOut)
			os.Exit(1)
		}
		if err := compareResults(want, got); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%sreference output:\n%s\ncandidate output:\n%s", i+1, err, input, wantOut, gotOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2150F-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("go build failed: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		abs, err := filepath.Abs(bin)
		if err != nil {
			return "", err
		}
		cmd = exec.Command("go", "run", abs)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

type testCase struct {
	n     int
	m     int
	edges [][2]int
}

type testInput struct {
	cases []testCase
}

func parseInput(input string) (testInput, error) {
	sc := bufio.NewScanner(strings.NewReader(input))
	sc.Split(bufio.ScanWords)
	nextInt := func() (int, error) {
		if !sc.Scan() {
			return 0, fmt.Errorf("unexpected EOF")
		}
		return strconv.Atoi(sc.Text())
	}
	t, err := nextInt()
	if err != nil {
		return testInput{}, err
	}
	cases := make([]testCase, t)
	for i := 0; i < t; i++ {
		n, err := nextInt()
		if err != nil {
			return testInput{}, err
		}
		m, err := nextInt()
		if err != nil {
			return testInput{}, err
		}
		edges := make([][2]int, m)
		for j := 0; j < m; j++ {
			u, err := nextInt()
			if err != nil {
				return testInput{}, err
			}
			v, err := nextInt()
			if err != nil {
				return testInput{}, err
			}
			edges[j] = [2]int{u, v}
		}
		cases[i] = testCase{n: n, m: m, edges: edges}
	}
	return testInput{cases: cases}, nil
}

type operation struct {
	k     int
	paths [][]int
}

type outputCase struct {
	ops []operation
}

type outputData struct {
	cases []outputCase
}

func parseOutput(input string, out string) (outputData, error) {
	inp, err := parseInput(input)
	if err != nil {
		return outputData{}, err
	}
	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Split(bufio.ScanWords)
	nextInt := func() (int, error) {
		if !sc.Scan() {
			return 0, fmt.Errorf("unexpected EOF")
		}
		return strconv.Atoi(sc.Text())
	}
	cases := make([]outputCase, len(inp.cases))
	for i := range inp.cases {
		opCnt, err := nextInt()
		if err != nil {
			return outputData{}, err
		}
		if opCnt < 0 || opCnt > 2 {
			return outputData{}, fmt.Errorf("case %d: operation count must be between 0 and 2", i+1)
		}
		ops := make([]operation, opCnt)
		for j := 0; j < opCnt; j++ {
			k, err := nextInt()
			if err != nil {
				return outputData{}, err
			}
			if k < 1 || k > inp.cases[i].n {
				return outputData{}, fmt.Errorf("case %d: operation %d has invalid k=%d", i+1, j+1, k)
			}
			s, err := nextInt()
			if err != nil {
				return outputData{}, err
			}
			if s < 0 || s > (inp.cases[i].n*(inp.cases[i].n-1))/2 {
				return outputData{}, fmt.Errorf("case %d: operation %d invalid path count s=%d", i+1, j+1, s)
			}
			paths := make([][]int, s)
			for p := 0; p < s; p++ {
				path := make([]int, k)
				for idx := 0; idx < k; idx++ {
					val, err := nextInt()
					if err != nil {
						return outputData{}, err
					}
					path[idx] = val
				}
				paths[p] = path
			}
			ops[j] = operation{k: k, paths: paths}
		}
		cases[i] = outputCase{ops: ops}
	}
	return outputData{cases: cases}, nil
}

func compareResults(want outputData, got outputData) error {
	if len(want.cases) != len(got.cases) {
		return fmt.Errorf("number of cases mismatch")
	}
	for i := range want.cases {
		if err := compareCase(want.cases[i], got.cases[i]); err != nil {
			return fmt.Errorf("case %d: %w", i+1, err)
		}
	}
	return nil
}

func compareCase(ref outputCase, cand outputCase) error {
	if len(ref.ops) != len(cand.ops) {
		return fmt.Errorf("operation count mismatch: expected %d got %d", len(ref.ops), len(cand.ops))
	}
	for i := range ref.ops {
		if ref.ops[i].k != cand.ops[i].k {
			return fmt.Errorf("operation %d: expected k=%d got k=%d", i+1, ref.ops[i].k, cand.ops[i].k)
		}
		if len(ref.ops[i].paths) != len(cand.ops[i].paths) {
			return fmt.Errorf("operation %d: path count mismatch (expected %d got %d)", i+1, len(ref.ops[i].paths), len(cand.ops[i].paths))
		}
		for j := range ref.ops[i].paths {
			if !equalPath(ref.ops[i].paths[j], cand.ops[i].paths[j]) {
				return fmt.Errorf("operation %d path %d mismatch", i+1, j+1)
			}
		}
	}
	return nil
}

func equalPath(a, b []int) bool {
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

func generateTests() []string {
	var tests []string
	tests = append(tests, buildTest([]testCase{
		{
			n:     3,
			m:     2,
			edges: [][2]int{{1, 2}, {2, 3}},
		},
	}))

	tests = append(tests, buildTest([]testCase{
		{
			n:     4,
			m:     4,
			edges: [][2]int{{1, 2}, {2, 3}, {3, 4}, {1, 4}},
		},
	}))

	tests = append(tests, buildTest([]testCase{
		randomGraph(5, 6),
		randomGraph(6, 8),
	}))

	for i := 0; i < 5; i++ {
		tests = append(tests, buildTest([]testCase{
			randomGraph(8+i, 12+i*2),
		}))
	}

	return tests
}

func buildTest(cases []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, cs := range cases {
		sb.WriteString(fmt.Sprintf("%d %d\n", cs.n, cs.m))
		for _, e := range cs.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
	}
	return sb.String()
}

var rng = rand.New(rand.NewSource(42))

func randomGraph(n int, m int) testCase {
	if n < 3 {
		n = 3
	}
	maxEdges := n * (n - 1) / 2
	if m < n-1 {
		m = n - 1
	}
	if m > maxEdges {
		m = maxEdges
	}
	edges := make([][2]int, 0, m)
	// Start with tree
	for i := 2; i <= n; i++ {
		u := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{u, i})
	}
	existing := make(map[[2]int]bool)
	for _, e := range edges {
		if e[0] > e[1] {
			e[0], e[1] = e[1], e[0]
		}
		existing[e] = true
	}
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		if existing[[2]int{u, v}] {
			continue
		}
		edges = append(edges, [2]int{u, v})
		existing[[2]int{u, v}] = true
	}
	return testCase{
		n:     n,
		m:     len(edges),
		edges: edges,
	}
}
