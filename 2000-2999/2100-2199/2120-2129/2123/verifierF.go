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
}

func callerDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("cannot determine caller")
	}
	return filepath.Dir(file), nil
}

func buildOracle() (string, func(), error) {
	dir, err := callerDir()
	if err != nil {
		return "", nil, err
	}
	tmpDir, err := os.MkdirTemp("", "oracle-2123F-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleF")
	cmd := exec.Command("go", "build", "-o", outPath, "2123F.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
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
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d\n", tc.n)
	}
	return sb.String()
}

func parseOutputs(out string, tests []testCase) ([][]int, error) {
	fields := strings.Fields(out)
	res := make([][]int, len(tests))
	pos := 0
	for idx, tc := range tests {
		if pos+tc.n > len(fields) {
			return nil, fmt.Errorf("not enough tokens for test %d: need %d, have %d", idx+1, tc.n, len(fields)-pos)
		}
		arr := make([]int, tc.n)
		for i := 0; i < tc.n; i++ {
			val, err := strconv.Atoi(fields[pos+i])
			if err != nil {
				return nil, fmt.Errorf("test %d token %d invalid int %q: %v", idx+1, i+1, fields[pos+i], err)
			}
			arr[i] = val
		}
		res[idx] = arr
		pos += tc.n
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("extra tokens after parsing: %d", len(fields)-pos)
	}
	return res, nil
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func checkPermutation(arr []int) error {
	n := len(arr)
	seen := make([]bool, n+1)
	for i, v := range arr {
		if v < 1 || v > n {
			return fmt.Errorf("position %d has value %d out of range [1,%d]", i+1, v, n)
		}
		if seen[v] {
			return fmt.Errorf("duplicate value %d", v)
		}
		seen[v] = true
	}
	return nil
}

func checkGood(arr []int) error {
	if err := checkPermutation(arr); err != nil {
		return err
	}
	for i := 1; i < len(arr); i++ { // 0-indexed, skip i=0 (1-based index 1)
		if gcd(arr[i], i+1) == 1 {
			return fmt.Errorf("gcd condition fails at position %d: p=%d gcd=1", i+1, arr[i])
		}
	}
	return nil
}

func fixedPoints(arr []int) int {
	cnt := 0
	for i, v := range arr {
		if v == i+1 {
			cnt++
		}
	}
	return cnt
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 2},
		{n: 3},
		{n: 4},
		{n: 5},
		{n: 6},
		{n: 7},
		{n: 12},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 60)
	sumN := 0
	limit := 100000

	add := func(n int) {
		if n <= 0 || sumN+n > limit {
			return
		}
		tests = append(tests, testCase{n: n})
		sumN += n
	}

	for i := 0; i < 20; i++ {
		add(rng.Intn(15) + 2)
	}
	for i := 0; i < 20; i++ {
		add(rng.Intn(200) + 50)
	}
	for i := 0; i < 10; i++ {
		add(rng.Intn(3000) + 1000)
	}
	// include one big case if space allows
	add(limit - sumN)
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)
	input := buildInput(tests)

	expOut, err := runBinary(oracle, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle runtime error: %v\noutput:\n%s\n", err, expOut)
		os.Exit(1)
	}
	gotOut, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\noutput:\n%s\n", err, gotOut)
		os.Exit(1)
	}

	expPerms, err := parseOutputs(expOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse oracle output: %v\n", err)
		os.Exit(1)
	}
	gotPerms, err := parseOutputs(gotOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse target output: %v\n", err)
		os.Exit(1)
	}

	for idx := range tests {
		exp := expPerms[idx]
		got := gotPerms[idx]

		if err := checkGood(exp); err != nil {
			fmt.Fprintf(os.Stderr, "oracle produced invalid permutation on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if err := checkGood(got); err != nil {
			fmt.Fprintf(os.Stderr, "test %d invalid permutation: %v\n", idx+1, err)
			os.Exit(1)
		}

		expFixed := fixedPoints(exp)
		gotFixed := fixedPoints(got)
		if gotFixed != expFixed {
			fmt.Fprintf(os.Stderr, "test %d: fixed points mismatch, expected %d got %d\n", idx+1, expFixed, gotFixed)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}
