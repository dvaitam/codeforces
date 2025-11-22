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
	n   int
	arr []int
}

func callerFile() (string, bool) {
	_, file, _, ok := runtime.Caller(0)
	return file, ok
}

func buildOracle() (string, func(), error) {
	_, file, ok := callerFile()
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2081D-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleD")
	cmd := exec.Command("go", "build", "-o", outPath, "2081D.go")
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

func buildInput(tcs []testCase) string {
	var sb strings.Builder
	sb.Grow(32 + len(tcs)*32)
	sb.WriteString(strconv.Itoa(len(tcs)))
	sb.WriteByte('\n')
	for _, tc := range tcs {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseAnswers(out string, t int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d answers, got %d", t, len(fields))
	}
	res := make([]int64, t)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		res[i] = val
	}
	return res, nil
}

func compareAnswers(exp, got []int64) error {
	if len(exp) != len(got) {
		return fmt.Errorf("answer count mismatch")
	}
	for i := range exp {
		if exp[i] != got[i] {
			return fmt.Errorf("test %d expected %d, got %d", i+1, exp[i], got[i])
		}
	}
	return nil
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 1, arr: []int{5}},
		{n: 4, arr: []int{4, 3, 3, 4}},
		{n: 5, arr: []int{2, 10, 3, 2, 9}},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 150)
	totalN := 0
	const limitN = 400000
	for len(tests) < cap(tests) && totalN < limitN {
		n := rng.Intn(8000) + 1
		if totalN+n > limitN {
			n = limitN - totalN
		}
		if rng.Intn(6) == 0 {
			n = rng.Intn(120) + 1
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			val := rng.Intn(500000) + 1
			if rng.Intn(6) == 0 {
				val = rng.Intn(50) + 1
			}
			arr[i] = val
		}
		tests = append(tests, testCase{n: n, arr: arr})
		totalN += n
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
		fmt.Fprintf(os.Stderr, "oracle failed: %v\n", err)
		os.Exit(1)
	}
	gotOut, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\n", err)
		os.Exit(1)
	}
	expAns, err := parseAnswers(expOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle output invalid: %v\n%s", err, expOut)
		os.Exit(1)
	}
	gotAns, err := parseAnswers(gotOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "target output invalid: %v\n%s", err, gotOut)
		os.Exit(1)
	}
	if err := compareAnswers(expAns, gotAns); err != nil {
		fmt.Fprintf(os.Stderr, "mismatch: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("All tests passed.")
}
