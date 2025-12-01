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
	refSource        = "2148B.go"
	tempOraclePrefix = "oracle-2148B-"
	randomTests      = 80
	maxN             = 20
	maxM             = 20
	maxCoord         = 1000
)

type testCase struct {
	n int
	m int
	x int64
	y int64
	a []int64
	b []int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build oracle: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTestsCases(rng, randomTests)...)

	for idx, tc := range tests {
		input := formatInput(tc)
		expOut, err := runProgram(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d: %v\n", idx+1, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", idx+1, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}

		expAns, err := parseSingleInt(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse oracle output on test %d: %v\noutput:\n%s", idx+1, err, expOut)
			os.Exit(1)
		}
		gotAns, err := parseSingleInt(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d: %v\noutput:\n%s", idx+1, err, candOut)
			os.Exit(1)
		}
		if expAns != gotAns {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %d got %d\n", idx+1, expAns, gotAns)
			fmt.Println("Input:")
			fmt.Print(input)
			fmt.Println("Expected:")
			fmt.Print(expOut)
			fmt.Println("Got:")
			fmt.Print(candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildOracle() (string, func(), error) {
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

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
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

func parseSingleInt(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d %d %d\n", tc.n, tc.m, tc.x, tc.y)
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	for i, v := range tc.b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 1, m: 1, x: 2, y: 2, a: []int64{1}, b: []int64{1}},
		{n: 2, m: 2, x: 5, y: 5, a: []int64{1, 4}, b: []int64{1, 4}},
		{n: 3, m: 1, x: 6, y: 6, a: []int64{1, 3, 5}, b: []int64{2}},
	}
}

func randomTestsCases(rng *rand.Rand, count int) []testCase {
	tests := make([]testCase, 0, count)
	for len(tests) < count {
		n := rng.Intn(maxN-1) + 1
		m := rng.Intn(maxM-1) + 1
		x := int64(rng.Intn(maxCoord-1) + 2)
		y := int64(rng.Intn(maxCoord-1) + 2)
		a := randomSorted(rng, n, y)
		b := randomSorted(rng, m, x)
		tests = append(tests, testCase{n: n, m: m, x: x, y: y, a: a, b: b})
	}
	return tests
}

func randomSorted(rng *rand.Rand, n int, limit int64) []int64 {
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = int64(rng.Int63n(limit-1) + 1)
	}
	sortInt64(arr)
	return arr
}

func sortInt64(a []int64) {
	if len(a) < 2 {
		return
	}
	for i := 1; i < len(a); i++ {
		key := a[i]
		j := i - 1
		for j >= 0 && a[j] > key {
			a[j+1] = a[j]
			j--
		}
		a[j+1] = key
	}
}
