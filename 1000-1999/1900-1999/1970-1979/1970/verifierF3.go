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
	"time"
)

type action struct {
	entity string
	act    string
	target string
}

type testCase struct {
	n, m  int
	grid  [][]string
	steps []action
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-1970F3-")
	if err != nil {
		return "", nil, err
	}
	path := filepath.Join(tmpDir, "oracleF3")
	cmd := exec.Command("go", "build", "-o", path, "1970F3.go")
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

func deterministicTests() []testCase {
	return []testCase{
		buildSimpleTest(),
		buildOwnGoalTest(),
		buildBludgerTest(),
		buildSnitchTest(),
	}
}

func buildSimpleTest() testCase {
	grid := [][]string{
		{"..", "RG", ".."},
		{"R0", ".Q", "B0"},
		{"..", "BG", ".."},
	}
	steps := []action{
		{"R0", "C", ".Q"},
		{"R0", "R", ""},
		{"R0", "T", ""},
		{"B0", "C", ".Q"},
		{"B0", "U", ""},
		{"B0", "U", ""},
		{"B0", "T", ""},
	}
	return testCase{n: 3, m: 3, grid: grid, steps: steps}
}

func buildOwnGoalTest() testCase {
	grid := [][]string{
		{"..", "RG", ".."},
		{"R0", ".Q", ".."},
		{"..", "BG", ".."},
	}
	steps := []action{
		{"R0", "C", ".Q"},
		{"R0", "U", ""},
		{"R0", "T", ""},
	}
	return testCase{n: 3, m: 3, grid: grid, steps: steps}
}

func buildBludgerTest() testCase {
	grid := [][]string{
		{"R0", "..", ".."},
		{"..", ".B", ".."},
		{"..", "..", "B0"},
	}
	steps := []action{
		{".B", "U", ""},
		{".B", "U", ""},
		{"B0", "L", ""},
		{"B0", "L", ""},
		{"B0", "L", ""},
	}
	return testCase{n: 3, m: 3, grid: grid, steps: steps}
}

func buildSnitchTest() testCase {
	grid := [][]string{
		{"R0", "..", ".."},
		{"..", ".S", ".."},
		{"..", "..", "BG"},
	}
	steps := []action{
		{"R0", "R", ""},
		{"R0", "R", ""},
		{"R0", "D", ""},
		{"R0", "C", ".S"},
	}
	return testCase{n: 3, m: 3, grid: grid, steps: steps}
}

func randomTest(rng *rand.Rand) testCase {
	n := rng.Intn(5) + 3
	m := rng.Intn(5) + 3
	grid := make([][]string, n)
	for i := range grid {
		grid[i] = make([]string, m)
		for j := range grid[i] {
			grid[i][j] = ".."
		}
	}
	grid[0][m/2] = "RG"
	grid[n-1][m/2] = "BG"
	grid[n/2][m/2] = ".Q"
	grid[n/2][0] = "R0"
	grid[n/2][m-1] = "B0"
	hasBl := rng.Intn(2) == 0
	if hasBl {
		grid[0][0] = ".B"
	}
	hasSn := rng.Intn(2) == 0
	if hasSn {
		grid[n-1][m-1] = ".S"
	}
	steps := []action{
		{"R0", "C", ".Q"},
		{"R0", "R", ""},
		{"R0", "R", ""},
		{"R0", "T", ""},
		{"B0", "L", ""},
		{"B0", "L", ""},
		{"B0", "T", ""},
	}
	if hasBl {
		steps = append(steps, action{".B", "D", ""})
		steps = append(steps, action{".B", "R", ""})
	}
	if hasSn {
		steps = append(steps, action{"B0", "R", ""})
		steps = append(steps, action{"B0", "C", ".S"})
	}
	return testCase{n: n, m: m, grid: grid, steps: steps}
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(tc.grid[i][j])
		}
		sb.WriteByte('\n')
	}
	sb.WriteString(strconv.Itoa(len(tc.steps)))
	sb.WriteByte('\n')
	for _, st := range tc.steps {
		sb.WriteString(st.entity)
		sb.WriteByte(' ')
		sb.WriteString(st.act)
		if st.target != "" {
			sb.WriteByte(' ')
			sb.WriteString(st.target)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF3.go /path/to/binary")
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
		gotOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		if expOut != gotOut {
			fmt.Fprintf(os.Stderr, "mismatch on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, input, expOut, gotOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
