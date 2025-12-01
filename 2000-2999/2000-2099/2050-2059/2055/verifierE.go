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
	a []int64
	b []int64
}

func callerFile() (string, bool) {
	_, file, _, ok := runtime.Caller(0)
	return file, ok
}

func buildOracle() (string, func(), error) {
	file, ok := callerFile()
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2055E-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleE")
	cmd := exec.Command("go", "build", "-o", outPath, "2055E.go")
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
	sb.Grow(32 + len(tcs)*24)
	sb.WriteString(strconv.Itoa(len(tcs)))
	sb.WriteByte('\n')
	for _, tc := range tcs {
		n := len(tc.a)
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for i := 0; i < n; i++ {
			sb.WriteString(strconv.FormatInt(tc.a[i], 10))
			sb.WriteByte(' ')
			sb.WriteString(strconv.FormatInt(tc.b[i], 10))
			sb.WriteByte('\n')
		}
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
		{
			a: []int64{3, 5},
			b: []int64{2, 4},
		},
		{
			a: []int64{1, 1, 1},
			b: []int64{1, 1, 1},
		},
		{
			a: []int64{10, 1, 10},
			b: []int64{1, 10, 3},
		},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 150)
	totalN := 0
	const limitN = 400000
	for len(tests) < cap(tests) && totalN < limitN {
		n := rng.Intn(5000) + 2
		if totalN+n > limitN {
			n = limitN - totalN
		}
		if rng.Intn(6) == 0 {
			n = rng.Intn(150) + 2
		}
		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			valA := rng.Int63n(1_000_000_000) + 1
			valB := rng.Int63n(1_000_000_000) + 1
			if rng.Intn(5) == 0 {
				valA = int64(rng.Intn(20) + 1)
			}
			if rng.Intn(5) == 0 {
				valB = int64(rng.Intn(20) + 1)
			}
			a[i] = valA
			b[i] = valB
		}
		tests = append(tests, testCase{a: a, b: b})
		totalN += n
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
