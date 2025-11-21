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

const (
	refSource        = "424D.go"
	tempOraclePrefix = "oracle-424D-"
	randomTestsCount = 80
	maxRandomSize    = 40
)

type testCase struct {
	name       string
	n, m       int
	target     int64
	tp, tu, td int
	height     [][]int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	oraclePath, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTests(randomTestsCount, rng)...)

	for idx, tc := range tests {
		input := formatInput(tc)
		oracleOut, err := runProgram(oraclePath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		candidateOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}

		oracleRect, err := parseRectangles(oracleOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse oracle output on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, oracleOut)
			os.Exit(1)
		}
		candidateRect, err := parseRectangles(candidateOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output parse error on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, candidateOut)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}

		refDiff, err := evaluateRectangle(tc, oracleRect)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle rectangle invalid on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		candDiff, err := evaluateRectangle(tc, candidateRect)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate rectangle invalid on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}

		if candDiff != refDiff {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: time difference %d differs from optimal %d\n", idx+1, tc.name, candDiff, refDiff)
			fmt.Println("Input:")
			fmt.Print(input)
			fmt.Println("Candidate rectangle:", candidateRect)
			fmt.Println("Oracle rectangle:", oracleRect)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("failed to determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", tempOraclePrefix)
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleD")
	cmd := exec.Command("go", "build", "-o", outPath, refSource)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
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

func parseRectangles(out string) ([4]int, error) {
	fields := strings.Fields(out)
	if len(fields) != 4 {
		return [4]int{}, fmt.Errorf("expected 4 integers, got %d tokens", len(fields))
	}
	var rect [4]int
	for i := 0; i < 4; i++ {
		val, err := strconv.Atoi(fields[i])
		if err != nil {
			return [4]int{}, fmt.Errorf("token %q is not an integer", fields[i])
		}
		rect[i] = val
	}
	return rect, nil
}

func evaluateRectangle(tc testCase, rect [4]int) (int64, error) {
	r1, c1, r2, c2 := rect[0], rect[1], rect[2], rect[3]
	if r1 < 1 || r1 > tc.n || r2 < 1 || r2 > tc.n || c1 < 1 || c1 > tc.m || c2 < 1 || c2 > tc.m {
		return 0, fmt.Errorf("rectangle coordinates out of bounds")
	}
	if r1 >= r2 || c1 >= c2 {
		return 0, fmt.Errorf("rectangle must have positive area")
	}
	if r2-r1+1 < 3 || c2-c1+1 < 3 {
		return 0, fmt.Errorf("rectangle sides must have at least three cells")
	}
	perimeterTime, err := perimeterCost(tc, r1, c1, r2, c2)
	if err != nil {
		return 0, err
	}
	diff := perimeterTime - tc.target
	if diff < 0 {
		diff = -diff
	}
	return diff, nil
}

func perimeterCost(tc testCase, r1, c1, r2, c2 int) (int64, error) {
	var total int64
	// top edge
	for c := c1; c < c2; c++ {
		cost, err := moveCost(tc, r1, c, r1, c+1)
		if err != nil {
			return 0, err
		}
		total += cost
	}
	// right edge
	for r := r1; r < r2; r++ {
		cost, err := moveCost(tc, r, c2, r+1, c2)
		if err != nil {
			return 0, err
		}
		total += cost
	}
	// bottom edge
	for c := c2; c > c1; c-- {
		cost, err := moveCost(tc, r2, c, r2, c-1)
		if err != nil {
			return 0, err
		}
		total += cost
	}
	// left edge
	for r := r2; r > r1; r-- {
		cost, err := moveCost(tc, r, c1, r-1, c1)
		if err != nil {
			return 0, err
		}
		total += cost
	}
	return total, nil
}

func moveCost(tc testCase, r1, c1, r2, c2 int) (int64, error) {
	if r1 < 1 || r1 > tc.n || r2 < 1 || r2 > tc.n || c1 < 1 || c1 > tc.m || c2 < 1 || c2 > tc.m {
		return 0, fmt.Errorf("edge goes out of bounds")
	}
	if (abs(r1-r2) + abs(c1-c2)) != 1 {
		return 0, fmt.Errorf("edge endpoints are not adjacent")
	}
	h1 := tc.height[r1-1][c1-1]
	h2 := tc.height[r2-1][c2-1]
	switch {
	case h1 == h2:
		return int64(tc.tp), nil
	case h1 < h2:
		return int64(tc.tu), nil
	default:
		return int64(tc.td), nil
	}
}

func deterministicTests() []testCase {
	return []testCase{
		makeManualTest("all_flat_3x3", 3, 3, 10, 3, 4, 5, constantGrid(3, 3, 7)),
		makeManualTest("ascending_rows", 4, 5, 50, 2, 5, 1, incrementalRows(4, 5)),
		makeManualTest("descending_cols", 6, 6, 80, 3, 6, 2, incrementalCols(6, 6)),
		makeManualTest("mixed_medium", 8, 9, 120, 2, 4, 3, checkerMix(8, 9)),
		makeManualTest("largeish", 30, 30, 400, 1, 6, 2, randomGridSeed(30, 30, 17)),
	}
}

func randomTests(count int, rng *rand.Rand) []testCase {
	tests := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		n := rng.Intn(maxRandomSize-2) + 3
		m := rng.Intn(maxRandomSize-2) + 3
		grid := make([][]int, n)
		for r := 0; r < n; r++ {
			grid[r] = make([]int, m)
			for c := 0; c < m; c++ {
				grid[r][c] = rng.Intn(1000000) + 1
			}
		}
		tp := rng.Intn(10) + 1
		tu := rng.Intn(10) + 1
		td := rng.Intn(10) + 1
		target := rng.Int63n(200000) + 1
		tests = append(tests, makeManualTest(
			fmt.Sprintf("random_%d", i+1),
			n, m, int(target), tp, tu, td, grid,
		))
	}
	return tests
}

func makeManualTest(name string, n, m, target, tp, tu, td int, grid [][]int) testCase {
	return testCase{
		name:   name,
		n:      n,
		m:      m,
		target: int64(target),
		tp:     tp,
		tu:     tu,
		td:     td,
		height: grid,
	}
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", tc.n, tc.m, tc.target)
	fmt.Fprintf(&sb, "%d %d %d\n", tc.tp, tc.tu, tc.td)
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(tc.height[i][j]))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func constantGrid(n, m, val int) [][]int {
	grid := make([][]int, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]int, m)
		for j := 0; j < m; j++ {
			grid[i][j] = val
		}
	}
	return grid
}

func incrementalRows(n, m int) [][]int {
	grid := make([][]int, n)
	val := 1
	for i := 0; i < n; i++ {
		grid[i] = make([]int, m)
		for j := 0; j < m; j++ {
			grid[i][j] = val
			val += 2
		}
	}
	return grid
}

func incrementalCols(n, m int) [][]int {
	grid := make([][]int, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]int, m)
		for j := 0; j < m; j++ {
			grid[i][j] = (j + 1) * 10
		}
	}
	return grid
}

func checkerMix(n, m int) [][]int {
	grid := make([][]int, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]int, m)
		for j := 0; j < m; j++ {
			if (i+j)%2 == 0 {
				grid[i][j] = 100 + i
			} else {
				grid[i][j] = 200 + j
			}
		}
	}
	return grid
}

func randomGridSeed(n, m int, seed int64) [][]int {
	grid := make([][]int, n)
	rng := rand.New(rand.NewSource(seed))
	for i := 0; i < n; i++ {
		grid[i] = make([]int, m)
		for j := 0; j < m; j++ {
			grid[i][j] = rng.Intn(1000000) + 1
		}
	}
	return grid
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
