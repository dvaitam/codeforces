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
	a []int64
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-324A1-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(tmpDir, "oracleA1")
	cmd := exec.Command("go", "build", "-o", bin, "324A1.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return bin, cleanup, nil
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

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i, val := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(val, 10))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func parseSubmission(out string, n int) (int64, []int, error) {
	fields := strings.Fields(out)
	if len(fields) < 2 {
		return 0, nil, fmt.Errorf("expected at least two integers, got %d", len(fields))
	}
	sum, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, nil, fmt.Errorf("invalid sum %q: %v", fields[0], err)
	}
	k, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, nil, fmt.Errorf("invalid k %q: %v", fields[1], err)
	}
	if k < 0 || k > n {
		return 0, nil, fmt.Errorf("k=%d out of range", k)
	}
	if len(fields) != 2+k {
		return 0, nil, fmt.Errorf("expected %d indices, got %d tokens", k, len(fields)-2)
	}
	removed := make([]int, k)
	for i := 0; i < k; i++ {
		idx, err := strconv.Atoi(fields[2+i])
		if err != nil {
			return 0, nil, fmt.Errorf("invalid index %q: %v", fields[2+i], err)
		}
		removed[i] = idx
	}
	return sum, removed, nil
}

func validateAnswer(tc testCase, sum int64, removed []int, bestSum int64) error {
	n := tc.n
	if len(removed) > n {
		return fmt.Errorf("too many removals: %d", len(removed))
	}
	used := make([]bool, n)
	for _, idx := range removed {
		if idx < 1 || idx > n {
			return fmt.Errorf("index %d out of range", idx)
		}
		if used[idx-1] {
			return fmt.Errorf("duplicate removal index %d", idx)
		}
		used[idx-1] = true
	}
	remain := make([]int64, 0, n-len(removed))
	for i, val := range tc.a {
		if !used[i] {
			remain = append(remain, val)
		}
	}
	if len(remain) < 2 {
		return fmt.Errorf("less than two trees remain (%d)", len(remain))
	}
	if remain[0] != remain[len(remain)-1] {
		return fmt.Errorf("first and last values differ: %d vs %d", remain[0], remain[len(remain)-1])
	}
	actualSum := int64(0)
	for _, val := range remain {
		actualSum += val
	}
	if actualSum != sum {
		return fmt.Errorf("reported sum %d but actual sum %d", sum, actualSum)
	}
	if sum != bestSum {
		return fmt.Errorf("sum %d is not optimal (expected %d)", sum, bestSum)
	}
	if len(removed) != n-len(remain) {
		return fmt.Errorf("reported number of removals inconsistent")
	}
	return nil
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 2, a: []int64{1, 1}},
		{n: 2, a: []int64{-5, -5}},
		{n: 3, a: []int64{5, -1, 5}},
		{n: 4, a: []int64{1, 2, 2, 1}},
		{n: 5, a: []int64{3, -2, 3, -2, 3}},
		{n: 6, a: []int64{10, -5, -6, 10, -1, 10}},
	}
}

func randomTest(rng *rand.Rand) testCase {
	n := rng.Intn(200) + 2
	if rng.Intn(5) == 0 {
		n = rng.Intn(1000) + 2
	}
	a := make([]int64, n)
	for i := range a {
		switch rng.Intn(3) {
		case 0:
			a[i] = int64(rng.Intn(21) - 10)
		case 1:
			a[i] = int64(rng.Intn(2001) - 1000)
		default:
			a[i] = int64(rng.Intn(1_000_001)) - 500_000
		}
	}
	// ensure at least one duplicate value
	i1 := rng.Intn(n)
	i2 := rng.Intn(n - 1)
	if i2 >= i1 {
		i2++
	}
	a[i2] = a[i1]
	return testCase{n: n, a: a}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA1.go /path/to/binary")
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
		bestSum, _, err := parseSubmission(expOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse oracle output on test %d: %v\noutput:\n%s\n", idx+1, err, expOut)
			os.Exit(1)
		}

		actOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		sum, removed, err := parseSubmission(actOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s\n", idx+1, err, actOut)
			os.Exit(1)
		}
		if err := validateAnswer(tc, sum, removed, bestSum); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%s\noutput:\n%s\n", idx+1, err, input, actOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
