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

const refSource = "2097E.go"

type testCase struct {
	n int
	d int
	a []int64
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := args[0]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}
	refAns, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	candAns, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n%s", err, candOut)
		os.Exit(1)
	}

	for i := range tests {
		if refAns[i] != candAns[i] {
			tc := tests[i]
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %d got %d\n", i+1, refAns[i], candAns[i])
			fmt.Fprintf(os.Stderr, "n=%d d=%d a=%v\n", tc.n, tc.d, tc.a)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	refPath, err := referencePath()
	if err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp("", "ref_2097E_*.bin")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), refPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func referencePath() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to locate verifier path")
	}
	dir := filepath.Dir(file)
	return filepath.Join(dir, refSource), nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func buildTests() []testCase {
	var tests []testCase
	add := func(tc testCase) {
		tests = append(tests, tc)
	}

	// Deterministic small cases
	add(testCase{n: 1, d: 1, a: []int64{5}})
	add(testCase{n: 2, d: 1, a: []int64{3, 4}})
	add(testCase{n: 3, d: 2, a: []int64{1, 5, 2}})
	add(testCase{n: 5, d: 2, a: []int64{1, 5, 2, 1, 2}})
	add(testCase{n: 4, d: 4, a: []int64{10, 10, 10, 10}})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	totalN := 0

	makeRandom := func(n, d int, maxVal int64) testCase {
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Int63n(maxVal) + 1
		}
		return testCase{n: n, d: d, a: arr}
	}

	for len(tests) < 25 && totalN < 60000 {
		n := rng.Intn(3000) + 1
		d := rng.Intn(n) + 1
		tc := makeRandom(n, d, 1_000_000_000)
		totalN += n
		add(tc)
	}

	// Larger structured cases
	n1 := 50000
	arr1 := make([]int64, n1)
	for i := range arr1 {
		arr1[i] = 1
	}
	add(testCase{n: n1, d: 1, a: arr1})

	n2 := 40000
	arr2 := make([]int64, n2)
	for i := range arr2 {
		arr2[i] = int64((i % 100) + 1)
	}
	add(testCase{n: n2, d: 200, a: arr2})

	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		if len(tc.a) != tc.n {
			panic("array length mismatch")
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.d))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, expected int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d outputs, got %d", expected, len(fields))
	}
	ans := make([]int64, expected)
	for i, s := range fields {
		val, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", s)
		}
		ans[i] = val
	}
	return ans, nil
}
