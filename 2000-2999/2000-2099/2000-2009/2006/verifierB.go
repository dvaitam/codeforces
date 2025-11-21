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

type event struct {
	x int
	y int64
}

type testCase struct {
	n      int
	w      int64
	parent []int
	events []event
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2006B-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleB")
	cmd := exec.Command("go", "build", "-o", outPath, "2006B.go")
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

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.Grow(len(tests) * 128)
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.w))
		for i := 2; i <= tc.n; i++ {
			if i > 2 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(tc.parent[i]))
		}
		sb.WriteByte('\n')
		for _, ev := range tc.events {
			sb.WriteString(fmt.Sprintf("%d %d\n", ev.x, ev.y))
		}
	}
	return sb.String()
}

func parseOutput(out string, tests []testCase) ([][]int64, error) {
	fields := strings.Fields(out)
	idx := 0
	res := make([][]int64, len(tests))
	for i, tc := range tests {
		need := tc.n - 1
		if idx+need > len(fields) {
			return nil, fmt.Errorf("not enough numbers for test %d", i+1)
		}
		ans := make([]int64, need)
		for j := 0; j < need; j++ {
			val, err := strconv.ParseInt(fields[idx+j], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q", fields[idx+j])
			}
			ans[j] = val
		}
		res[i] = ans
		idx += need
	}
	if idx != len(fields) {
		return nil, fmt.Errorf("extra numbers in output")
	}
	return res, nil
}

func compareAnswers(expected, actual [][]int64) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("test case count mismatch")
	}
	for i := range expected {
		if len(expected[i]) != len(actual[i]) {
			return fmt.Errorf("test %d answer length mismatch: expected %d, got %d", i+1, len(expected[i]), len(actual[i]))
		}
		for j := range expected[i] {
			if expected[i][j] != actual[i][j] {
				return fmt.Errorf("test %d event %d mismatch: expected %d, got %d", i+1, j+1, expected[i][j], actual[i][j])
			}
		}
	}
	return nil
}

func newTestCase(parent []int, weights []int64, order []int) testCase {
	n := len(parent) - 1
	if len(weights) != n+1 {
		panic("weights length mismatch")
	}
	if len(order) != n-1 {
		panic("order length mismatch")
	}
	used := make([]bool, n+1)
	for _, v := range order {
		if v < 2 || v > n {
			panic("order value out of range")
		}
		if used[v] {
			panic("duplicate order value")
		}
		used[v] = true
	}
	var w int64
	for i := 2; i <= n; i++ {
		w += weights[i]
	}
	events := make([]event, 0, n-1)
	for _, x := range order {
		events = append(events, event{x: x, y: weights[x]})
	}
	return testCase{
		n:      n,
		w:      w,
		parent: parent,
		events: events,
	}
}

func deterministicTests() []testCase {
	// Case 1: n=2
	parent1 := []int{0, 1, 1}
	weights1 := []int64{0, 0, 5}
	order1 := []int{2}

	// Case 2: n=4
	parent2 := []int{0, 1, 1, 2, 2}
	weights2 := []int64{0, 0, 2, 3, 4}
	order2 := []int{2, 4, 3}

	// Case 3: n=5
	parent3 := []int{0, 1, 1, 1, 3, 4}
	weights3 := []int64{0, 0, 7, 1, 5, 2}
	order3 := []int{5, 2, 4, 3}

	return []testCase{
		newTestCase(parent1, weights1, order1),
		newTestCase(parent2, weights2, order2),
		newTestCase(parent3, weights3, order3),
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 120)
	for len(tests) < cap(tests) {
		n := rng.Intn(30) + 2
		parent := make([]int, n+1)
		parent[1] = 1
		for i := 2; i <= n; i++ {
			parent[i] = rng.Intn(i-1) + 1
		}
		weights := make([]int64, n+1)
		for i := 2; i <= n; i++ {
			weights[i] = rng.Int63n(1_000_000)
		}
		order := make([]int, n-1)
		permutation := randPerm(rng, n-1)
		for i, idx := range permutation {
			order[i] = idx + 2
		}
		tests = append(tests, newTestCase(parent, weights, order))
	}
	return tests
}

func randPerm(rng *rand.Rand, n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = i
	}
	for i := n - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
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

	expectedOut, err := runBinary(oracle, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle failed: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}
	actualOut, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	expectedAns, err := parseOutput(expectedOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle output invalid: %v\n%s", err, expectedOut)
		os.Exit(1)
	}
	actualAns, err := parseOutput(actualOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target output invalid: %v\n%s", err, actualOut)
		os.Exit(1)
	}

	if err := compareAnswers(expectedAns, actualAns); err != nil {
		fmt.Fprintf(os.Stderr, "%v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	fmt.Println("All tests passed.")
}
