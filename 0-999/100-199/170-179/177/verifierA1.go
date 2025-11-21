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
	refSource        = "177A1.go"
	randomTestCount  = 120
	maxMatrixSize    = 101
	maxCellValue     = 100
	tempOraclePrefix = "oracle-177A1-"
)

type testCase struct {
	n      int
	values []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA1.go /path/to/binary")
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
	tests = append(tests, generateRandomTests(randomTestCount, rng)...)

	for idx, tc := range tests {
		input := formatInput(tc)
		exp, err := runProgram(oraclePath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d: %v\n", idx+1, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		got, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", idx+1, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\n", idx+1, exp, got)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("failed to determine verifier location")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", tempOraclePrefix)
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleA1")
	cmd := exec.Command("go", "build", "-o", outPath, refSource)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("build oracle failed: %v\n%s", err, string(out))
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

func formatInput(tc testCase) string {
	var b strings.Builder
	b.Grow(8 + tc.n*tc.n*4)
	fmt.Fprintf(&b, "%d\n", tc.n)
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.n; j++ {
			if j > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(strconv.Itoa(tc.values[i*tc.n+j]))
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func deterministicTests() []testCase {
	tests := []testCase{
		constantMatrix(1, 0),
		constantMatrix(1, maxCellValue),
		rowsToTestCase([][]int{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		}),
		rowsToTestCase([][]int{
			{0, 1, 0, 1, 0},
			{1, 2, 3, 4, 5},
			{6, 7, 8, 9, 10},
			{11, 12, 13, 14, 15},
			{16, 17, 18, 19, 20},
		}),
		checkerboardMatrix(9),
		sequentialMatrix(maxMatrixSize),
	}
	tests = append(tests, constantMatrix(101, 37))
	return tests
}

func generateRandomTests(count int, rng *rand.Rand) []testCase {
	tests := make([]testCase, 0, count)
	for len(tests) < count {
		n := rng.Intn((maxMatrixSize+1)/2)*2 + 1
		data := make([]int, n*n)
		for i := range data {
			data[i] = rng.Intn(maxCellValue + 1)
		}
		tests = append(tests, testCase{n: n, values: data})
	}
	return tests
}

func rowsToTestCase(rows [][]int) testCase {
	n := len(rows)
	if n == 0 {
		panic("empty matrix in rowsToTestCase")
	}
	if n%2 == 0 {
		panic("matrix size must be odd")
	}
	data := make([]int, 0, n*n)
	for _, row := range rows {
		if len(row) != n {
			panic("non-square matrix provided")
		}
		data = append(data, row...)
	}
	return testCase{n: n, values: data}
}

func constantMatrix(n, value int) testCase {
	if n%2 == 0 {
		panic("matrix size must be odd")
	}
	data := make([]int, n*n)
	for i := range data {
		data[i] = value
	}
	return testCase{n: n, values: data}
}

func checkerboardMatrix(n int) testCase {
	if n%2 == 0 {
		panic("matrix size must be odd")
	}
	data := make([]int, n*n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if (i+j)%2 == 0 {
				data[i*n+j] = 0
			} else {
				data[i*n+j] = maxCellValue
			}
		}
	}
	return testCase{n: n, values: data}
}

func sequentialMatrix(n int) testCase {
	if n%2 == 0 {
		panic("matrix size must be odd")
	}
	data := make([]int, n*n)
	for i := range data {
		data[i] = i % (maxCellValue + 1)
	}
	return testCase{n: n, values: data}
}
