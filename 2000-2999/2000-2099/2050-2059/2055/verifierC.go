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
	n, m int
	path string
	grid [][]int64 // with path cells set to 0 (tampered input)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	// Build reference once to ensure we can exercise the same environment.
	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)
	input := serializeInput(tests)

	// Run reference to ensure tests are solvable (not used for correctness comparison).
	if _, err := runProgram(refBin, input); err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}

	out, err := runProgram(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	if err := validateOutput(out, tests); err != nil {
		fmt.Fprintf(os.Stderr, "validation failed: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	fmt.Println("All tests passed")
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine current path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "ref-2055C-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "ref2055C")
	cmd := exec.Command("go", "build", "-o", outPath, "2055C.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}

func serializeInput(tests []testCase) string {
	var sb strings.Builder
	sb.Grow(64 + len(tests)*32)
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		sb.WriteString(tc.path)
		sb.WriteByte('\n')
		for i := 0; i < tc.n; i++ {
			for j := 0; j < tc.m; j++ {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.FormatInt(tc.grid[i][j], 10))
			}
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func validateOutput(out string, tests []testCase) error {
	tokens := strings.Fields(out)
	pos := 0

	for idx, tc := range tests {
		need := tc.n * tc.m
		if pos+need > len(tokens) {
			return fmt.Errorf("test %d: not enough numbers, expected %d more got %d", idx+1, need, len(tokens)-pos)
		}

		onPath := buildPath(tc.n, tc.m, tc.path)

		grid := make([][]int64, tc.n)
		for i := 0; i < tc.n; i++ {
			grid[i] = make([]int64, tc.m)
			for j := 0; j < tc.m; j++ {
				val, err := strconv.ParseInt(tokens[pos], 10, 64)
				if err != nil {
					return fmt.Errorf("test %d: invalid integer %q", idx+1, tokens[pos])
				}
				grid[i][j] = val
				pos++
			}
		}

		// Check unchanged off-path cells.
		for i := 0; i < tc.n; i++ {
			for j := 0; j < tc.m; j++ {
				if !onPath[i][j] && grid[i][j] != tc.grid[i][j] {
					return fmt.Errorf("test %d: cell (%d,%d) modified (expected %d got %d)", idx+1, i+1, j+1, tc.grid[i][j], grid[i][j])
				}
				if onPath[i][j] {
					if grid[i][j] < -1e15 || grid[i][j] > 1e15 {
						return fmt.Errorf("test %d: cell (%d,%d) out of allowed range", idx+1, i+1, j+1)
					}
				}
			}
		}

		// Check row/column sums equal.
		rowSum := make([]int64, tc.n)
		colSum := make([]int64, tc.m)
		for i := 0; i < tc.n; i++ {
			for j := 0; j < tc.m; j++ {
				rowSum[i] += grid[i][j]
				colSum[j] += grid[i][j]
			}
		}

		target := rowSum[0]
		for i := 0; i < tc.n; i++ {
			if rowSum[i] != target {
				return fmt.Errorf("test %d: row %d sum mismatch (got %d want %d)", idx+1, i+1, rowSum[i], target)
			}
		}
		for j := 0; j < tc.m; j++ {
			if colSum[j] != target {
				return fmt.Errorf("test %d: column %d sum mismatch (got %d want %d)", idx+1, j+1, colSum[j], target)
			}
		}
	}

	if pos != len(tokens) {
		return fmt.Errorf("extra output tokens: expected %d used %d", pos, len(tokens))
	}
	return nil
}

func buildPath(n, m int, s string) [][]bool {
	on := make([][]bool, n)
	for i := range on {
		on[i] = make([]bool, m)
	}
	x, y := 0, 0
	on[x][y] = true
	for _, ch := range s {
		if ch == 'D' {
			x++
		} else {
			y++
		}
		on[x][y] = true
	}
	return on
}

func deterministicTests() []testCase {
	tests := make([]testCase, 0, 4)
	// From samples (smaller subset due to size).
	tests = append(tests, buildTestCase(3, 3, "DRRD", [][]int64{
		{0, 2, 3},
		{0, 0, 0},
		{3, 1, 0},
	}))
	tests = append(tests, buildTestCase(2, 2, "DR", [][]int64{
		{0, 0},
		{0, 0},
	}))
	// Rectangle requiring zero sums.
	tests = append(tests, autoTestCase(3, 5, "DDRRR"))
	// Square with non-zero total.
	tests = append(tests, autoTestCase(4, 4, "DDRRDDRR"))
	return tests
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 80)
	for len(tests) < cap(tests) {
		n := rng.Intn(8) + 2
		m := rng.Intn(8) + 2
		path := randomPath(rng, n, m)
		tests = append(tests, autoTestCase(n, m, path))
	}
	return tests
}

func buildTestCase(n, m int, path string, grid [][]int64) testCase {
	// grid provided with path cells zero already.
	return testCase{n: n, m: m, path: path, grid: grid}
}

func autoTestCase(n, m int, path string) testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	orig := makeBalancedGrid(rng, n, m)
	on := buildPath(n, m, path)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if on[i][j] {
				orig[i][j] = 0
			}
		}
	}
	return testCase{n: n, m: m, path: path, grid: orig}
}

// Constructs a grid where all row sums and column sums are equal (zero if n != m).
func makeBalancedGrid(rng *rand.Rand, n, m int) [][]int64 {
	g := make([][]int64, n)
	rowSum := make([]int64, n)
	colSum := make([]int64, m)
	for i := 0; i < n; i++ {
		g[i] = make([]int64, m)
	}

	// Fill (n-1) x (m-1) block.
	for i := 0; i < n-1; i++ {
		for j := 0; j < m-1; j++ {
			val := int64(rng.Intn(201) - 100)
			g[i][j] = val
			rowSum[i] += val
			colSum[j] += val
		}
	}

	// Set last column to fix row sums.
	for i := 0; i < n-1; i++ {
		val := -rowSum[i]
		g[i][m-1] = val
		colSum[m-1] += val
		rowSum[i] = 0
	}

	// Set last row (except last cell) to fix column sums.
	var lastRowSum int64
	for j := 0; j < m-1; j++ {
		val := -colSum[j]
		g[n-1][j] = val
		lastRowSum += val
		colSum[j] = 0
	}

	// Set bottom-right cell to fix last column and last row.
	g[n-1][m-1] = -colSum[m-1]
	lastRowSum += g[n-1][m-1]
	colSum[m-1] = 0
	rowSum[n-1] = lastRowSum

	// If square, optionally shift by constant to get non-zero common sum.
	if n == m {
		shift := int64(rng.Intn(21) - 10)
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				g[i][j] += shift
			}
		}
	}
	return g
}

func randomPath(rng *rand.Rand, n, m int) string {
	steps := make([]byte, 0, n+m-2)
	for i := 0; i < n-1; i++ {
		steps = append(steps, 'D')
	}
	for j := 0; j < m-1; j++ {
		steps = append(steps, 'R')
	}
	for i := len(steps) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		steps[i], steps[j] = steps[j], steps[i]
	}
	return string(steps)
}
