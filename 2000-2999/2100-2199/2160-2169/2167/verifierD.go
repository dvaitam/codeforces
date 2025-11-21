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

const (
	maxValue int64 = 1_000_000_000_000_000_000
	maxN           = 100000
)

type singleCase struct {
	arr []int64
}

type testCase struct {
	name  string
	cases []singleCase
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2167D-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleD")
	cmd := exec.Command("go", "build", "-o", outPath, "2167D.go")
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

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.cases)))
	for _, cs := range tc.cases {
		sb.WriteString(fmt.Sprintf("%d\n", len(cs.arr)))
		for i, v := range cs.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{
			name: "single_small",
			cases: []singleCase{
				{arr: []int64{1}},
			},
		},
		{
			name: "power_of_two",
			cases: []singleCase{
				{arr: []int64{2, 4, 8, 16}},
				{arr: []int64{3, 9, 27}},
			},
		},
		{
			name: "mixed_primes",
			cases: []singleCase{
				{arr: []int64{6, 10, 15}},
				{arr: []int64{35, 77, 143, 221}},
				{arr: []int64{2, 3, 5, 7, 11}},
			},
		},
		{
			name: "large_values",
			cases: []singleCase{
				{arr: []int64{maxValue, maxValue - 1}},
				{arr: []int64{9999999967, 9999999967}},
			},
		},
	}
}

func randomArray(n int, rng *rand.Rand) []int64 {
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Int63n(maxValue) + 1
	}
	return arr
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 15)
	for t := 0; t < 15; t++ {
		var cases []singleCase
		remaining := maxN
		for remaining > 0 {
			limit := remaining
			if limit > 2000 {
				limit = 2000
			}
			n := rng.Intn(limit) + 1
			cases = append(cases, singleCase{arr: randomArray(n, rng)})
			remaining -= n
			if len(cases) >= 50 && rng.Intn(3) == 0 {
				break
			}
		}
		tests = append(tests, testCase{
			name:  fmt.Sprintf("random_batch_%d", t+1),
			cases: cases,
		})
	}
	return tests
}

func stressTests() []testCase {
	// One test with n = 1e5 consisting of alternating primes and even numbers.
	arr := make([]int64, maxN)
	primes := []int64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}
	for i := 0; i < maxN; i++ {
		if i%2 == 0 {
			arr[i] = primes[i%len(primes)]
		} else {
			arr[i] = maxValue - int64(i)
			if arr[i] <= 0 {
				arr[i] = 1
			}
		}
	}
	return []testCase{
		{
			name:  "stress_large",
			cases: []singleCase{{arr: arr}},
		},
	}
}

func normalizeOutput(output string) []string {
	return strings.Fields(output)
}

func compareOutputs(expected, actual string) error {
	exp := normalizeOutput(expected)
	act := normalizeOutput(actual)
	if len(exp) != len(act) {
		return fmt.Errorf("token count mismatch: expected %d got %d", len(exp), len(act))
	}
	for i := range exp {
		if exp[i] != act[i] {
			return fmt.Errorf("mismatch at token %d: expected %s got %s", i+1, exp[i], act[i])
		}
	}
	return nil
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
	tests = append(tests, stressTests()...)

	for idx, tc := range tests {
		input := buildInput(tc)
		expected, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, input)
			os.Exit(1)
		}
		actual, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) runtime error: %v\ninput:\n%s", idx+1, tc.name, err, input)
			os.Exit(1)
		}
		if err := compareOutputs(strings.TrimSpace(expected), strings.TrimSpace(actual)); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", idx+1, tc.name, err, input, expected, actual)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
