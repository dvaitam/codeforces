package main

import (
	"bytes"
	"fmt"
	"math"
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
	tmpDir, err := os.MkdirTemp("", "oracle-2131B-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleB")
	cmd := exec.Command("go", "build", "-o", outPath, "2131B.go")
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

func parseOutputs(out string, tests []testCase) ([][]int64, error) {
	fields := strings.Fields(out)
	res := make([][]int64, len(tests))
	pos := 0
	for idx, tc := range tests {
		if pos+tc.n > len(fields) {
			return nil, fmt.Errorf("not enough tokens for test %d: need %d, have %d", idx+1, tc.n, len(fields)-pos)
		}
		arr := make([]int64, tc.n)
		for i := 0; i < tc.n; i++ {
			val, err := strconv.ParseInt(fields[pos+i], 10, 64)
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

func isGood(arr []int64) error {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		if arr[i] == 0 || arr[i+1] == 0 || arr[i]*arr[i+1] >= 0 {
			return fmt.Errorf("adjacent product at %d not negative (%d, %d)", i+1, arr[i], arr[i+1])
		}
	}

	pref := make([]int64, n+1)
	minPref := pref[0]
	for i := 0; i < n; i++ {
		pref[i+1] = pref[i] + arr[i]
		if i >= 1 { // subarrays ending at i (length>=2) check later
			if pref[i+1]-minPref <= 0 {
				return fmt.Errorf("non-positive subarray sum ending at %d", i+1)
			}
		}
		if i >= 0 && pref[i] < minPref {
			minPref = pref[i]
		}
	}
	return nil
}

func absSeq(arr []int64) []int64 {
	res := make([]int64, len(arr))
	for i, v := range arr {
		res[i] = int64(math.Abs(float64(v)))
	}
	return res
}

func compareLex(a, b []int64) int {
	n := len(a)
	for i := 0; i < n; i++ {
		if a[i] < b[i] {
			return -1
		}
		if a[i] > b[i] {
			return 1
		}
	}
	return 0
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 2},
		{n: 3},
		{n: 4},
		{n: 5},
		{n: 10},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 40)
	total := 0
	limit := 200000

	add := func(n int) {
		if n <= 0 || total+n > limit {
			return
		}
		tests = append(tests, testCase{n: n})
		total += n
	}

	for i := 0; i < 15; i++ {
		add(rng.Intn(10) + 2)
	}
	for i := 0; i < 15; i++ {
		add(rng.Intn(200) + 50)
	}
	for i := 0; i < 5; i++ {
		add(rng.Intn(5000) + 1000)
	}
	add(limit - total)
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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

	expArrs, err := parseOutputs(expOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse oracle output: %v\n", err)
		os.Exit(1)
	}
	gotArrs, err := parseOutputs(gotOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse target output: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		exp := expArrs[idx]
		got := gotArrs[idx]
		if len(got) != tc.n {
			fmt.Fprintf(os.Stderr, "test %d length mismatch: expected %d got %d\n", idx+1, tc.n, len(got))
			os.Exit(1)
		}
		if err := isGood(exp); err != nil {
			fmt.Fprintf(os.Stderr, "oracle produced invalid array on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if err := isGood(got); err != nil {
			fmt.Fprintf(os.Stderr, "test %d invalid array: %v\n", idx+1, err)
			os.Exit(1)
		}
		expAbs := absSeq(exp)
		gotAbs := absSeq(got)
		if cmp := compareLex(gotAbs, expAbs); cmp != 0 {
			if cmp < 0 {
				fmt.Fprintf(os.Stderr, "test %d: candidate lexicographically smaller than oracle abs sequence, oracle may be wrong\n", idx+1)
			} else {
				fmt.Fprintf(os.Stderr, "test %d: abs sequence not minimal, expected %v got %v\n", idx+1, expAbs, gotAbs)
			}
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}
