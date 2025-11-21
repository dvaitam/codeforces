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
	refSource        = "699B.go"
	tempOraclePrefix = "oracle-699B-"
	randomTestsCount = 120
	maxRandomSize    = 35
)

type testCase struct {
	name string
	n    int
	m    int
	grid []string
}

type result struct {
	hasSolution bool
	row         int
	col         int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
	tests = append(tests, largeTests()...)

	for idx, tc := range tests {
		input := formatInput(tc)

		oracleOut, err := runProgram(oraclePath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		oracleRes, err := parseResult(oracleOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse oracle output on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, oracleOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		candRes, err := parseResult(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output parse error on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, candOut)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}

		if oracleRes.hasSolution != candRes.hasSolution {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected solution=%v, got %v\n", idx+1, tc.name, oracleRes.hasSolution, candRes.hasSolution)
			fmt.Println("Input:")
			fmt.Print(input)
			fmt.Println("Candidate output:")
			fmt.Print(candOut)
			os.Exit(1)
		}

		if candRes.hasSolution {
			if err := validateBomb(tc, candRes.row, candRes.col); err != nil {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\n", idx+1, tc.name, err)
				fmt.Println("Input:")
				fmt.Print(input)
				fmt.Println("Candidate output:")
				fmt.Print(candOut)
				os.Exit(1)
			}
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
	outPath := filepath.Join(tmpDir, "oracleB")
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

func parseResult(out string) (result, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return result{}, fmt.Errorf("empty output")
	}
	switch strings.ToUpper(fields[0]) {
	case "NO":
		if len(fields) != 1 {
			return result{}, fmt.Errorf("NO output must not contain extra tokens")
		}
		return result{hasSolution: false}, nil
	case "YES":
		if len(fields) != 3 {
			return result{}, fmt.Errorf("YES output must be followed by exactly two integers")
		}
		row, err := strconv.Atoi(fields[1])
		if err != nil {
			return result{}, fmt.Errorf("invalid row value %q", fields[1])
		}
		col, err := strconv.Atoi(fields[2])
		if err != nil {
			return result{}, fmt.Errorf("invalid column value %q", fields[2])
		}
		return result{hasSolution: true, row: row, col: col}, nil
	default:
		return result{}, fmt.Errorf("first token must be YES or NO")
	}
}

func validateBomb(tc testCase, row, col int) error {
	if row < 1 || row > tc.n || col < 1 || col > tc.m {
		return fmt.Errorf("coordinates (%d,%d) out of bounds", row, col)
	}
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.m; j++ {
			if tc.grid[i][j] == '*' {
				if i+1 != row && j+1 != col {
					return fmt.Errorf("cell (%d,%d) remains after bombing (%d,%d)", i+1, j+1, row, col)
				}
			}
		}
	}
	return nil
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
	for _, line := range tc.grid {
		sb.WriteString(line)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{name: "empty_1x1", n: 1, m: 1, grid: []string{"."}},
		{name: "single_star", n: 1, m: 1, grid: []string{"*"}},
		{name: "row_solution", n: 3, m: 4, grid: []string{
			"...*",
			"****",
			"..*.",
		}},
		{name: "column_solution", n: 4, m: 3, grid: []string{
			"*..",
			"*..",
			"*..",
			"...",
		}},
		{name: "impossible_small", n: 2, m: 2, grid: []string{
			"*.",
			".*",
		}},
		{name: "all_walls", n: 3, m: 3, grid: []string{
			"***",
			"***",
			"***",
		}},
	}
}

func randomTests(count int, rng *rand.Rand) []testCase {
	tests := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		n := rng.Intn(maxRandomSize-1) + 1
		m := rng.Intn(maxRandomSize-1) + 1
		grid := make([]string, n)
		for r := 0; r < n; r++ {
			row := make([]byte, m)
			for c := 0; c < m; c++ {
				if rng.Intn(3) == 0 {
					row[c] = '*'
				} else {
					row[c] = '.'
				}
			}
			grid[r] = string(row)
		}
		tests = append(tests, testCase{
			name: fmt.Sprintf("random_%d", i+1),
			n:    n,
			m:    m,
			grid: grid,
		})
	}
	return tests
}

func largeTests() []testCase {
	largeYes := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		row := make([]byte, 1000)
		for j := 0; j < 1000; j++ {
			if i == 500 || j == 400 {
				row[j] = '*'
			} else {
				row[j] = '.'
			}
		}
		largeYes[i] = string(row)
	}
	largeNo := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		row := make([]byte, 1000)
		for j := 0; j < 1000; j++ {
			if (i == j && i%2 == 0) || (i+j == 999 && j%2 == 1) {
				row[j] = '*'
			} else {
				row[j] = '.'
			}
		}
		largeNo[i] = string(row)
	}
	return []testCase{
		{name: "large_yes", n: 1000, m: 1000, grid: largeYes},
		{name: "large_no", n: 1000, m: 1000, grid: largeNo},
	}
}
