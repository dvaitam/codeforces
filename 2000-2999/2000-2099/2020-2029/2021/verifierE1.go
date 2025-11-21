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
	n int
	m int
	p int

	special []int
	edges   [][3]int64
}

func main() {
	args := os.Args[1:]
	if len(args) >= 1 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierE1.go /path/to/solution")
		os.Exit(1)
	}
	target := args[0]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build oracle: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicCases(), randomCases()...)
	input := buildInput(tests)

	expectOut, err := runBinary(oracle, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle runtime error: %v\noutput:\n%s", err, expectOut)
		os.Exit(1)
	}

	candOut, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s", err, candOut)
		os.Exit(1)
	}

	expAns, err := parseOutput(expectOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse oracle output: %v\noutput:\n%s", err, expectOut)
		os.Exit(1)
	}
	candAns, err := parseOutput(candOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\noutput:\n%s", err, candOut)
		os.Exit(1)
	}

	for idx := range expAns {
		if len(expAns[idx]) != len(candAns[idx]) {
			fmt.Fprintf(os.Stderr, "test %d: expected %d values, got %d\ninput:\n%s", idx+1, len(expAns[idx]), len(candAns[idx]), input)
			os.Exit(1)
		}
		for k := range expAns[idx] {
			if expAns[idx][k] != candAns[idx][k] {
				fmt.Fprintf(os.Stderr, "test %d (k=%d): expected %d got %d\ninput:\n%s", idx+1, k+1, expAns[idx][k], candAns[idx][k], input)
				os.Exit(1)
			}
		}
	}

	fmt.Println("All tests passed.")
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier location")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2021E1-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleE1")
	cmd := exec.Command("go", "build", "-o", outPath, "2021E1.go")
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
		return stdout.String() + stderr.String(), fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.p))
		for i, v := range tc.special {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for _, e := range tc.edges {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2]))
		}
	}
	return sb.String()
}

func deterministicCases() []testCase {
	return []testCase{
		buildLineTree(2, 1, []int{1}),
		buildLineTree(3, 2, []int{1, 3}),
		buildLineTree(4, 3, []int{2, 4}),
		buildCycle(5, []int{1, 3, 5}),
		buildStar(6, []int{1, 2, 3, 4, 5, 6}),
		buildStar(6, []int{2, 4, 6}),
	}
}

func randomCases() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 80)
	for len(tests) < cap(tests) {
		n := rng.Intn(20) + 2
		maxM := n * (n - 1) / 2
		m := rng.Intn(maxM-(n-1)+1) + (n - 1)
		p := rng.Intn(n) + 1
		special := randSubset(rng, n, p)
		edges := randomGraph(rng, n, m)
		tests = append(tests, testCase{
			n:       n,
			m:       m,
			p:       p,
			special: special,
			edges:   edges,
		})
	}
	return tests
}

func randSubset(rng *rand.Rand, n, k int) []int {
	indices := rng.Perm(n)[:k]
	for i := range indices {
		indices[i]++
	}
	sortInts(indices)
	return indices
}

func randomGraph(rng *rand.Rand, n, m int) [][3]int64 {
	type pair struct{ u, v int }
	added := make(map[pair]bool)
	selected := make([]pair, 0, m)

	// ensure connectivity with a random spanning tree
	for v := 2; v <= n; v++ {
		u := rng.Intn(v-1) + 1
		a, b := u, v
		if a > b {
			a, b = b, a
		}
		p := pair{u: a, v: b}
		added[p] = true
		selected = append(selected, p)
	}

	// collect remaining edges
	remaining := make([]pair, 0, n*(n-1)/2)
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			p := pair{u: i, v: j}
			if !added[p] {
				remaining = append(remaining, p)
			}
		}
	}
	rng.Shuffle(len(remaining), func(i, j int) {
		remaining[i], remaining[j] = remaining[j], remaining[i]
	})
	for len(selected) < m && len(remaining) > 0 {
		p := remaining[len(remaining)-1]
		remaining = remaining[:len(remaining)-1]
		selected = append(selected, p)
	}

	res := make([][3]int64, len(selected))
	for i, e := range selected {
		w := rng.Intn(1000) + 1
		res[i] = [3]int64{int64(e.u), int64(e.v), int64(w)}
	}
	return res
}

func buildLineTree(n, baseWeight int64, special []int) testCase {
	edges := make([][3]int64, n-1)
	for i := 1; i < int(n); i++ {
		edges[i-1] = [3]int64{int64(i), int64(i + 1), baseWeight + int64(i)}
	}
	return testCase{
		n:       int(n),
		m:       int(n - 1),
		p:       len(special),
		special: append([]int(nil), special...),
		edges:   edges,
	}
}

func buildCycle(n int64, special []int) testCase {
	edges := make([][3]int64, n)
	for i := int64(1); i <= n; i++ {
		j := i%n + 1
		edges[i-1] = [3]int64{i, j, 1 + i}
	}
	return testCase{
		n:       int(n),
		m:       int(n),
		p:       len(special),
		special: append([]int(nil), special...),
		edges:   edges,
	}
}

func buildStar(n int64, special []int) testCase {
	edges := make([][3]int64, n-1)
	for i := int64(2); i <= n; i++ {
		edges[i-2] = [3]int64{1, i, i}
	}
	return testCase{
		n:       int(n),
		m:       int(n - 1),
		p:       len(special),
		special: append([]int(nil), special...),
		edges:   edges,
	}
}

func parseOutput(out string, tests []testCase) ([][]int64, error) {
	lines := strings.FieldsFunc(strings.TrimSpace(out), func(r rune) bool {
		return r == '\n' || r == '\r'
	})
	if len(lines) != len(tests) {
		return nil, fmt.Errorf("expected %d lines, got %d", len(tests), len(lines))
	}
	result := make([][]int64, len(tests))
	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != tests[i].n {
			return nil, fmt.Errorf("test %d: expected %d numbers, got %d", i+1, tests[i].n, len(fields))
		}
		row := make([]int64, len(fields))
		for j, f := range fields {
			val, err := strconv.ParseInt(f, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("test %d: invalid number %q", i+1, f)
			}
			if val < 0 {
				return nil, fmt.Errorf("test %d: negative answer %d", i+1, val)
			}
			row[j] = val
		}
		result[i] = row
	}
	return result, nil
}

func sortInts(a []int) {
	if len(a) <= 1 {
		return
	}
	for i := 1; i < len(a); i++ {
		j := i
		for j > 0 && a[j-1] > a[j] {
			a[j-1], a[j] = a[j], a[j-1]
			j--
		}
	}
}
