package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	valueLimit           = 100000000
	referenceSolutionRel = "0-999/100-199/170-179/178/178D1.go"
)

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "178D1.go")
		if _, err := os.Stat(candidate); err == nil {
			referenceSolutionPath = candidate
			return
		}
	}
	if _, err := os.Stat(referenceSolutionRel); err == nil {
		abs, err := filepath.Abs(referenceSolutionRel)
		if err == nil {
			referenceSolutionPath = abs
		}
	}
}

type testCase struct {
	n    int
	nums []int64
}

var baseSquares = map[int][][]int64{
	1: {
		{0},
	},
	2: {
		{0, 0},
		{0, 0},
	},
	3: {
		{8, 1, 6},
		{3, 5, 7},
		{4, 9, 2},
	},
	4: {
		{16, 2, 3, 13},
		{5, 11, 10, 8},
		{9, 7, 6, 12},
		{4, 14, 15, 1},
	},
}

func cloneGrid(src [][]int64) [][]int64 {
	n := len(src)
	dst := make([][]int64, n)
	for i := range src {
		dst[i] = append([]int64(nil), src[i]...)
	}
	return dst
}

func randomValue(rng *rand.Rand) int64 {
	return rng.Int63n(2*valueLimit+1) - valueLimit
}

func randomGrid(n int, rng *rand.Rand) [][]int64 {
	base := cloneGrid(baseSquares[n])
	grid := make([][]int64, n)
	switch n {
	case 1:
		grid[0] = []int64{randomValue(rng)}
		return grid
	case 2:
		val := randomValue(rng)
		for i := 0; i < n; i++ {
			grid[i] = make([]int64, n)
			for j := range grid[i] {
				grid[i][j] = val
			}
		}
		return grid
	default:
		absAdd := 84000000
		mulRange := 2001
		add := rng.Int63n(2*int64(absAdd)+1) - int64(absAdd)
		var mul int64
		for mul == 0 {
			mul = int64(rng.Intn(mulRange)) - int64(mulRange/2)
		}
		for i := 0; i < n; i++ {
			grid[i] = make([]int64, n)
			for j := 0; j < n; j++ {
				val := base[i][j]*mul + add
				if val < -valueLimit || val > valueLimit {
					return randomGrid(n, rng)
				}
				grid[i][j] = val
			}
		}
		return grid
	}
}

func flattenAndShuffle(grid [][]int64, rng *rand.Rand) []int64 {
	var vals []int64
	for i := range grid {
		vals = append(vals, grid[i]...)
	}
	rng.Shuffle(len(vals), func(i, j int) { vals[i], vals[j] = vals[j], vals[i] })
	return vals
}

func genTests() []testCase {
	rng := rand.New(rand.NewSource(42))
	var tests []testCase
	for n := 1; n <= 4; n++ {
		grid := cloneGrid(baseSquares[n])
		nums := flattenAndShuffle(grid, rng)
		tests = append(tests, testCase{n: n, nums: nums})
	}
	for i := 0; i < 200; i++ {
		n := rng.Intn(4) + 1
		grid := randomGrid(n, rng)
		tests = append(tests, testCase{
			n:    n,
			nums: flattenAndShuffle(grid, rng),
		})
	}
	return tests
}

func inputString(t testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t.n))
	for i, v := range t.nums {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildReferenceBinary() (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "178D1-ref")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_178D1")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSolutionPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return binPath, cleanup, nil
}

func checkOutput(t testCase, out string) error {
	reader := strings.NewReader(out)
	var s int64
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return fmt.Errorf("failed to read magic sum: %v", err)
	}
	total := int64(0)
	for _, v := range t.nums {
		total += v
	}
	if total%int64(t.n) != 0 {
		return fmt.Errorf("invalid test: total %d not divisible by %d", total, t.n)
	}
	expectedSum := total / int64(t.n)
	if s != expectedSum {
		return fmt.Errorf("reported sum %d differs from expected %d", s, expectedSum)
	}
	counts := make(map[int64]int)
	for _, v := range t.nums {
		counts[v]++
	}
	grid := make([][]int64, t.n)
	for i := 0; i < t.n; i++ {
		row := make([]int64, t.n)
		for j := 0; j < t.n; j++ {
			if _, err := fmt.Fscan(reader, &row[j]); err != nil {
				return fmt.Errorf("failed to read cell (%d,%d): %v", i+1, j+1, err)
			}
			counts[row[j]]--
		}
		grid[i] = row
	}
	for val, c := range counts {
		if c != 0 {
			return fmt.Errorf("value %d used incorrect number of times (delta %d)", val, c)
		}
	}
	for i := 0; i < t.n; i++ {
		sum := int64(0)
		for j := 0; j < t.n; j++ {
			sum += grid[i][j]
		}
		if sum != s {
			return fmt.Errorf("row %d sums to %d instead of %d", i+1, sum, s)
		}
	}
	for j := 0; j < t.n; j++ {
		sum := int64(0)
		for i := 0; i < t.n; i++ {
			sum += grid[i][j]
		}
		if sum != s {
			return fmt.Errorf("column %d sums to %d instead of %d", j+1, sum, s)
		}
	}
	diag := int64(0)
	for i := 0; i < t.n; i++ {
		diag += grid[i][i]
	}
	if diag != s {
		return fmt.Errorf("main diagonal sums to %d instead of %d", diag, s)
	}
	diag = 0
	for i := 0; i < t.n; i++ {
		diag += grid[i][t.n-1-i]
	}
	if diag != s {
		return fmt.Errorf("secondary diagonal sums to %d instead of %d", diag, s)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD1.go /path/to/binary")
		os.Exit(1)
	}
	args := os.Args[1:]
	bin := args[len(args)-1]
	if bin == "--" {
		fmt.Println("usage: go run verifierD1.go /path/to/binary")
		os.Exit(1)
	}
	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()
	tests := genTests()
	for i, t := range tests {
		in := inputString(t)
		// Run the reference implementation to ensure the testcase stays solvable,
		// but ignore its exact output because it does not enforce the diagonal constraints rigorously.
		if refOut, err := runProgram(refBin, in); err != nil {
			fmt.Fprintf(os.Stderr, "reference solution failed on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, in, refOut)
			os.Exit(1)
		}
		out, runErr := runProgram(bin, in)
		if runErr != nil {
			fmt.Fprintf(os.Stderr, "test %d runtime error: %v\ninput:\n%soutput:\n%s\n", i+1, runErr, in, out)
			os.Exit(1)
		}
		if err := checkOutput(t, out); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%soutput:\n%s\n", i+1, err, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
