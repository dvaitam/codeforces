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
	"time"
)

const (
	refSource2033G = "2000-2999/2000-2099/2030-2039/2033/2033G.go"
)

type testCase struct {
	n       int
	edges   [][2]int
	queries [][2]int
}

type testSet struct {
	name  string
	cases []testCase
}

func main() {
	candPath, ok := parseBinaryArg()
	if !ok {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}

	refBin, cleanupRef, err := buildBinary(refSource2033G)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference solution: %v\n", err)
		os.Exit(1)
	}
	defer cleanupRef()

	candBin, cleanupCand, err := buildBinary(candPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare candidate binary: %v\n", err)
		os.Exit(1)
	}
	defer cleanupCand()

	tests := buildTests()
	for idx, ts := range tests {
		input := buildInput(ts.cases)

		refOut, err := runBinary(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference solution failed on test %d (%s): %v\ninput:\n%s", idx+1, ts.name, err, input)
			os.Exit(1)
		}
		refAns, err := parseOutput(refOut, ts.cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal error: failed to parse reference output on test %d (%s): %v\noutput:\n%s", idx+1, ts.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runBinary(candBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%s", idx+1, ts.name, err, input)
			os.Exit(1)
		}
		candAns, err := parseOutput(candOut, ts.cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: cannot parse output: %v\ninput:\n%soutput:\n%s", idx+1, ts.name, err, input, candOut)
			os.Exit(1)
		}

		if err := compareAnswers(refAns, candAns); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\ninput:\n%sreference:\n%soutput:\n%s", idx+1, ts.name, err, input, refOut, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func parseBinaryArg() (string, bool) {
	if len(os.Args) == 2 {
		return os.Args[1], true
	}
	if len(os.Args) == 3 && os.Args[1] == "--" {
		return os.Args[2], true
	}
	return "", false
}

func buildBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "verifier2033G-*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(path))
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("%v\n%s", err, out.String())
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	return abs, func() {}, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func buildTests() []testSet {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return []testSet{
		pathCases(),
		starCases(),
		randomTest("random_small", rng, 4, 2, 10, 3, 15, 4),
		randomTest("random_medium", rng, 3, 50, 200, 60, 200, 15),
		heavyCase(rng),
	}
}

func pathCases() testSet {
	edges := func(n int) [][2]int {
		e := make([][2]int, 0, n-1)
		for i := 2; i <= n; i++ {
			e = append(e, [2]int{i - 1, i})
		}
		return e
	}
	return testSet{
		name: "paths",
		cases: []testCase{
			{
				n:       5,
				edges:   edges(5),
				queries: [][2]int{{5, 0}, {5, 1}, {3, 2}},
			},
			{
				n:       8,
				edges:   edges(8),
				queries: [][2]int{{8, 7}, {4, 0}, {2, 3}, {6, 10}},
			},
		},
	}
}

func starCases() testSet {
	n := 9
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		edges = append(edges, [2]int{1, i})
	}
	return testSet{
		name: "stars",
		cases: []testCase{
			{
				n:       n,
				edges:   edges,
				queries: [][2]int{{1, 0}, {5, 0}, {7, 1}, {3, 4}},
			},
		},
	}
}

func randomTest(name string, rng *rand.Rand, t, nMin, nMax, qMin, qMax, maxK int) testSet {
	var cases []testCase
	for len(cases) < t {
		n := rng.Intn(nMax-nMin+1) + nMin
		q := rng.Intn(qMax-qMin+1) + qMin
		if n+q > 200_000 {
			continue
		}
		cases = append(cases, buildRandomCase(rng, n, q, maxK))
	}
	return testSet{name: name, cases: cases}
}

func heavyCase(rng *rand.Rand) testSet {
	// Keep within total constraints but stress performance.
	n := 90_000
	q := 90_000
	return testSet{
		name: "heavy",
		cases: []testCase{
			buildRandomCase(rng, n, q, n),
		},
	}
}

func buildRandomCase(rng *rand.Rand, n, q, maxK int) testCase {
	edges := make([][2]int, 0, n-1)
	for v := 2; v <= n; v++ {
		p := rng.Intn(v-1) + 1
		edges = append(edges, [2]int{p, v})
	}
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		v := rng.Intn(n) + 1
		k := rng.Intn(maxK + 1)
		queries[i] = [2]int{v, k}
	}
	return testCase{n: n, edges: edges, queries: queries}
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(cases))
	for idx, tc := range cases {
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for _, e := range tc.edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
		fmt.Fprintf(&sb, "%d\n", len(tc.queries))
		for i, q := range tc.queries {
			fmt.Fprintf(&sb, "%d %d", q[0], q[1])
			if i+1 != len(tc.queries) {
				sb.WriteByte('\n')
			}
		}
		if idx+1 != len(cases) {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func parseOutput(out string, cases []testCase) ([][]int, error) {
	fields := strings.Fields(out)
	want := 0
	for _, tc := range cases {
		want += len(tc.queries)
	}
	if len(fields) != want {
		return nil, fmt.Errorf("expected %d numbers, got %d", want, len(fields))
	}
	res := make([][]int, len(cases))
	pos := 0
	for i, tc := range cases {
		res[i] = make([]int, len(tc.queries))
		for j := range tc.queries {
			v, err := strconv.Atoi(fields[pos])
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q", fields[pos])
			}
			res[i][j] = v
			pos++
		}
	}
	return res, nil
}

func compareAnswers(expect, got [][]int) error {
	if len(expect) != len(got) {
		return fmt.Errorf("expected %d testcases, got %d", len(expect), len(got))
	}
	for i := range expect {
		if len(expect[i]) != len(got[i]) {
			return fmt.Errorf("test %d: expected %d answers, got %d", i+1, len(expect[i]), len(got[i]))
		}
		for j := range expect[i] {
			if expect[i][j] != got[i][j] {
				return fmt.Errorf("test %d query %d mismatch: expected %d got %d", i+1, j+1, expect[i][j], got[i][j])
			}
		}
	}
	return nil
}
