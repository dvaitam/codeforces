package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
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
	p []int
	q []int
}

type op struct {
	i int
	j int
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier location")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2062G-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleG")
	cmd := exec.Command("go", "build", "-o", outPath, "2062G.go")
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

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.Grow(len(tests) * 512)
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for i, v := range tc.p {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for i, v := range tc.q {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, tests []testCase) ([][]op, error) {
	res := make([][]op, len(tests))
	r := bufio.NewReader(strings.NewReader(out))
	for idx, tc := range tests {
		var k int
		if _, err := fmt.Fscan(r, &k); err != nil {
			return nil, fmt.Errorf("test %d: cannot read operation count: %v", idx+1, err)
		}
		if k < 0 {
			return nil, fmt.Errorf("test %d: negative operation count %d", idx+1, k)
		}
		if k > tc.n*tc.n {
			return nil, fmt.Errorf("test %d: too many operations %d > %d", idx+1, k, tc.n*tc.n)
		}
		ops := make([]op, k)
		for i := 0; i < k; i++ {
			if _, err := fmt.Fscan(r, &ops[i].i, &ops[i].j); err != nil {
				return nil, fmt.Errorf("test %d: cannot read op %d: %v", idx+1, i+1, err)
			}
		}
		res[idx] = ops
	}
	return res, nil
}

func computeCost(tc testCase, ops []op) (int64, error) {
	p := append([]int(nil), tc.p...)
	n := tc.n
	var total int64
	for idx, op := range ops {
		i, j := op.i, op.j
		if i < 1 || i > n || j < 1 || j > n || i == j {
			return 0, fmt.Errorf("invalid swap (%d,%d) in op %d", i, j, idx+1)
		}
		pi := p[i-1]
		pj := p[j-1]
		c := abs(i - j)
		v := abs(pi - pj)
		if v < c {
			c = v
		}
		total += int64(c)
		p[i-1], p[j-1] = p[j-1], p[i-1]
	}
	for idx, v := range p {
		if v != tc.q[idx] {
			return 0, fmt.Errorf("final permutation mismatch at position %d: have %d, want %d", idx+1, v, tc.q[idx])
		}
	}
	return total, nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func compareCosts(expected, actual []int64, tests []testCase) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("cost length mismatch")
	}
	for i := range expected {
		if expected[i] != actual[i] {
			return fmt.Errorf("test %d: cost mismatch, expected %d, got %d", i+1, expected[i], actual[i])
		}
	}
	return nil
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 2,
			p: []int{1, 2},
			q: []int{1, 2},
		},
		{
			n: 2,
			p: []int{2, 1},
			q: []int{1, 2},
		},
		{
			n: 4,
			p: []int{4, 3, 2, 1},
			q: []int{1, 2, 3, 4},
		},
		{
			n: 5,
			p: []int{2, 1, 5, 3, 4},
			q: []int{5, 4, 3, 2, 1},
		},
	}
}

func randomPermutation(n int, rng *rand.Rand) []int {
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		perm[i] = i + 1
	}
	for i := n - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		perm[i], perm[j] = perm[j], perm[i]
	}
	return perm
}

func randomTests(limit int) []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 50)
	used := 0
	for used < limit {
		remaining := limit - used
		n := rng.Intn(99) + 2 // 2..100
		if n*n*n > remaining {
			n = max(2, int(math.Cbrt(float64(remaining))))
		}
		if n < 2 {
			n = 2
		}
		if n > 100 {
			n = 100
		}
		usage := n * n * n
		if usage > remaining {
			usage = remaining
			// shrink n to fit cube limit
			for n > 2 && n*n*n > remaining {
				n--
			}
		}
		p := randomPermutation(n, rng)
		q := randomPermutation(n, rng)
		tests = append(tests, testCase{n: n, p: p, q: q})
		used += n * n * n
	}
	return tests
}

func costsFromOutput(out string, tests []testCase) ([]int64, error) {
	ops, err := parseOutput(out, tests)
	if err != nil {
		return nil, err
	}
	res := make([]int64, len(tests))
	for i, tc := range tests {
		cost, err := computeCost(tc, ops[i])
		if err != nil {
			return nil, fmt.Errorf("test %d: %v", i+1, err)
		}
		res[i] = cost
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
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
	const cubeLimit = 900_000 // within problem bound 1e6
	used := 0
	for _, tc := range tests {
		used += tc.n * tc.n * tc.n
	}
	if used < cubeLimit {
		tests = append(tests, randomTests(cubeLimit-used)...)
	}

	input := buildInput(tests)

	expOut, err := runBinary(oracle, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle failed: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	actOut, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	expectedCosts, err := costsFromOutput(expOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle output invalid: %v\n%s", err, expOut)
		os.Exit(1)
	}
	actualCosts, err := costsFromOutput(actOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target output invalid: %v\n%s", err, actOut)
		os.Exit(1)
	}

	if err := compareCosts(expectedCosts, actualCosts, tests); err != nil {
		fmt.Fprintf(os.Stderr, "%v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	fmt.Println("All tests passed.")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
