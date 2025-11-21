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
	n     int
	edges [][2]int
}

type pair struct {
	a, b int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build oracle: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)

	for idx, tc := range tests {
		input := tc.input()

		expOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, input, expOut)
			os.Exit(1)
		}
		expected, err := parseOutput(expOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse oracle output on test %d: %v\noutput:\n%s", idx+1, err, expOut)
			os.Exit(1)
		}

		gotOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, input, gotOut)
			os.Exit(1)
		}
		got, err := parseOutput(gotOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output parse error on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, input, gotOut)
			os.Exit(1)
		}

		if len(got) != len(expected) {
			fmt.Fprintf(os.Stderr, "test %d: pair count mismatch (expected %d, got %d)\ninput:\n%s", idx+1, len(expected), len(got), input)
			os.Exit(1)
		}
		for i := range expected {
			if expected[i] != got[i] {
				fmt.Fprintf(os.Stderr, "test %d: mismatch at pair %d (expected %d %d, got %d %d)\ninput:\n%s", idx+1, i+1, expected[i].a, expected[i].b, got[i].a, got[i].b, input)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier location")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-212E-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleE")
	cmd := exec.Command("go", "build", "-o", outPath, "212E.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, string(out))
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
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

func (tc testCase) input() string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for _, e := range tc.edges {
		sb.WriteString(strconv.Itoa(e[0]))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(e[1]))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func deterministicTests() []testCase {
	tests := []testCase{
		{n: 3, edges: [][2]int{{1, 2}, {2, 3}}},
		{n: 4, edges: [][2]int{{1, 2}, {2, 3}, {3, 4}}},
		{n: 5, edges: [][2]int{{1, 2}, {1, 3}, {3, 4}, {3, 5}}},
		{n: 6, edges: [][2]int{{1, 2}, {1, 3}, {1, 4}, {1, 5}, {1, 6}}},
		{n: 8, edges: [][2]int{{1, 2}, {1, 3}, {2, 4}, {2, 5}, {3, 6}, {6, 7}, {7, 8}}},
	}
	tests = append(tests, testCase{n: 5000, edges: lineTree(5000)})
	return tests
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 40)
	for len(tests) < 40 {
		var n int
		switch {
		case len(tests) < 10:
			n = rng.Intn(7) + 3 // 3..9
		case len(tests) < 25:
			n = rng.Intn(40) + 10 // 10..49
		default:
			n = rng.Intn(200-50) + 50 // 50..199
		}
		edges := randomTree(n, rng)
		tests = append(tests, testCase{n: n, edges: edges})
	}
	return tests
}

func randomTree(n int, rng *rand.Rand) [][2]int {
	if n < 2 {
		return nil
	}
	if n == 2 {
		return [][2]int{{1, 2}}
	}
	deg := make([]int, n)
	for i := range deg {
		deg[i] = 1
	}
	prufer := make([]int, n-2)
	for i := range prufer {
		v := rng.Intn(n)
		prufer[i] = v
		deg[v]++
	}
	edges := make([][2]int, 0, n-1)
	for _, v := range prufer {
		leaf := -1
		for i := 0; i < n; i++ {
			if deg[i] == 1 {
				leaf = i
				break
			}
		}
		edges = append(edges, [2]int{leaf + 1, v + 1})
		deg[leaf]--
		deg[v]--
	}
	first, second := -1, -1
	for i := 0; i < n; i++ {
		if deg[i] == 1 {
			if first == -1 {
				first = i + 1
			} else {
				second = i + 1
				break
			}
		}
	}
	edges = append(edges, [2]int{first, second})
	return edges
}

func lineTree(n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 1; i < n; i++ {
		edges = append(edges, [2]int{i, i + 1})
	}
	return edges
}

func parseOutput(out string, n int) ([]pair, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty output")
	}
	z, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("invalid pair count: %v", err)
	}
	if z <= 0 {
		return nil, fmt.Errorf("pair count must be positive, got %d", z)
	}
	if len(fields) != 1+2*z {
		return nil, fmt.Errorf("expected %d integers, got %d", 1+2*z, len(fields))
	}
	pairs := make([]pair, z)
	ptr := 1
	for i := 0; i < z; i++ {
		a, err := strconv.Atoi(fields[ptr])
		if err != nil {
			return nil, fmt.Errorf("invalid a at pair %d: %v", i+1, err)
		}
		ptr++
		b, err := strconv.Atoi(fields[ptr])
		if err != nil {
			return nil, fmt.Errorf("invalid b at pair %d: %v", i+1, err)
		}
		ptr++
		if a <= 0 || b <= 0 {
			return nil, fmt.Errorf("pair %d must be positive, got (%d,%d)", i+1, a, b)
		}
		if a+b != n-1 {
			return nil, fmt.Errorf("pair %d sum mismatch: %d + %d != %d", i+1, a, b, n-1)
		}
		pairs[i] = pair{a: a, b: b}
	}
	for i := 1; i < len(pairs); i++ {
		if pairs[i-1].a >= pairs[i].a {
			return nil, fmt.Errorf("pairs not strictly increasing at position %d", i+1)
		}
	}
	return pairs, nil
}
