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

type cell struct {
	x int
	y int
}

type testCase struct {
	k    int
	grid [][]uint64
	fig  []cell
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-1270I-")
	if err != nil {
		return "", nil, err
	}
	path := filepath.Join(tmpDir, "oracleI")
	cmd := exec.Command("go", "build", "-o", path, "1270I.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
	return path, cleanup, nil
}

func runBinary(bin string, input string) (string, error) {
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
	return strings.TrimSpace(stdout.String()), nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.k))
	n := 1 << uint(tc.k)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatUint(tc.grid[i][j], 10))
		}
		sb.WriteByte('\n')
	}
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.fig)))
	for _, c := range tc.fig {
		sb.WriteString(fmt.Sprintf("%d %d\n", c.x+1, c.y+1))
	}
	return sb.String()
}

func parseOutput(out string) (int64, error) {
	if out == "-1" {
		return -1, nil
	}
	val, err := strconv.ParseInt(out, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer: %v", err)
	}
	if val < 0 {
		return 0, fmt.Errorf("operations count must be non-negative")
	}
	return val, nil
}

func deterministicTests() []testCase {
	return []testCase{
		genTest(1, [][]uint64{{0, 0}, {0, 0}}, []cell{{0, 0}}),
		genTest(2, [][]uint64{
			{1, 0, 0, 0},
			{0, 1, 0, 0},
			{0, 0, 1, 0},
			{0, 0, 0, 1},
		}, []cell{{0, 0}, {1, 1}, {2, 2}}),
	}
}

func genTest(k int, grid [][]uint64, fig []cell) testCase {
	n := 1 << uint(k)
	copyGrid := make([][]uint64, n)
	for i := 0; i < n; i++ {
		copyGrid[i] = make([]uint64, n)
		copy(copyGrid[i], grid[i])
	}
	return testCase{k: k, grid: copyGrid, fig: fig}
}

func randomTest(rng *rand.Rand) testCase {
	k := rng.Intn(4) + 1 // keep grid manageable
	n := 1 << uint(k)
	grid := make([][]uint64, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]uint64, n)
		for j := 0; j < n; j++ {
			if rng.Intn(4) == 0 {
				grid[i][j] = uint64(rng.Intn(16))
			} else {
				grid[i][j] = 0
			}
		}
	}
	limit := n * n
	if limit > 99 {
		limit = 99
	}
	t := rng.Intn(limit) + 1
	if t%2 == 0 {
		t++
		if t > limit {
			t -= 2
		}
	}
	used := make(map[cell]bool)
	fig := make([]cell, 0, t)
	for len(fig) < t {
		c := cell{x: rng.Intn(n), y: rng.Intn(n)}
		if used[c] {
			continue
		}
		used[c] = true
		fig = append(fig, c)
	}
	return testCase{k: k, grid: grid, fig: fig}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierI.go /path/to/binary")
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
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng))
	}

	for idx, tc := range tests {
		input := buildInput(tc)
		expOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		expVal, err := parseOutput(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\noutput:\n%s\n", idx+1, err, expOut)
			os.Exit(1)
		}

		gotOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		gotVal, err := parseOutput(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s\ninput:\n%s\n", idx+1, err, gotOut, input)
			os.Exit(1)
		}
		if gotVal != expVal {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %d got %d\ninput:\n%s\n", idx+1, expVal, gotVal, input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
