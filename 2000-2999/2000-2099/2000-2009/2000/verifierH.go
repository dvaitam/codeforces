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

type operation struct {
	typ string
	val int
}

type testCase struct {
	initial []int
	ops     []operation
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2000H-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleH")
	cmd := exec.Command("go", "build", "-o", outPath, "2000H.go")
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
	sb.Grow(len(tc.initial)*6 + len(tc.ops)*8 + 64)
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(len(tc.initial)))
	sb.WriteByte('\n')
	for i, v := range tc.initial {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	sb.WriteString(strconv.Itoa(len(tc.ops)))
	sb.WriteByte('\n')
	for _, op := range tc.ops {
		sb.WriteString(op.typ)
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(op.val))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string) ([]int, error) {
	fields := strings.Fields(out)
	res := make([]int, len(fields))
	for i, f := range fields {
		val, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", f, err)
		}
		res[i] = val
	}
	return res, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{
			initial: []int{1, 2, 5, 9},
			ops: []operation{
				{typ: "-", val: 2},
				{typ: "?", val: 2},
				{typ: "?", val: 1},
				{typ: "-", val: 1},
				{typ: "?", val: 1},
				{typ: "+", val: 4},
				{typ: "+", val: 2},
				{typ: "?", val: 2},
			},
		},
		{
			initial: []int{3, 4, 5, 6, 8},
			ops: []operation{
				{typ: "?", val: 5},
				{typ: "-", val: 5},
				{typ: "?", val: 5},
				{typ: "+", val: 1},
				{typ: "?", val: 2},
				{typ: "-", val: 6},
				{typ: "-", val: 8},
				{typ: "+", val: 6},
				{typ: "?", val: 5},
			},
		},
		{
			initial: []int{6, 7, 8, 9, 10},
			ops: []operation{
				{typ: "?", val: 5},
				{typ: "-", val: 6},
				{typ: "?", val: 4},
				{typ: "-", val: 10},
				{typ: "+", val: 5},
				{typ: "-", val: 8},
				{typ: "+", val: 3},
				{typ: "+", val: 2},
				{typ: "-", val: 3},
				{typ: "+", val: 10},
				{typ: "?", val: 5},
			},
		},
	}
}

func randomOperations(rng *rand.Rand, size int) []operation {
	ops := make([]operation, 0, size)
	current := make(map[int]bool)
	for len(current) < 3 {
		val := rng.Intn(30) + 1
		current[val] = true
	}
	keys := make([]int, 0, len(current))
	for k := range current {
		keys = append(keys, k)
	}
	for _, k := range keys {
		ops = append(ops, operation{typ: "+", val: k})
	}
	for len(ops) < size {
		if len(current) == 0 || rng.Intn(2) == 0 {
			val := rng.Intn(50) + 1
			for current[val] {
				val = rng.Intn(50) + 1
			}
			current[val] = true
			ops = append(ops, operation{typ: "+", val: val})
		} else {
			idx := rng.Intn(len(current))
			var val int
			i := 0
			for k := range current {
				if i == idx {
					val = k
					break
				}
				i++
			}
			delete(current, val)
			ops = append(ops, operation{typ: "-", val: val})
		}
		if rng.Intn(3) == 0 {
			k := rng.Intn(10) + 1
			ops = append(ops, operation{typ: "?", val: k})
		}
	}
	return ops
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 60)
	for len(tests) < cap(tests) {
		n := rng.Intn(10) + 1
		initialSet := make(map[int]bool)
		for len(initialSet) < n {
			initialSet[rng.Intn(50)+1] = true
		}
		initial := make([]int, 0, len(initialSet))
		for val := range initialSet {
			initial = append(initial, val)
		}
		ops := randomOperations(rng, rng.Intn(50)+10)
		tests = append(tests, testCase{initial: initial, ops: ops})
	}
	return tests
}

func filterQueries(ops []operation) int {
	count := 0
	for _, op := range ops {
		if op.typ == "?" {
			count++
		}
	}
	return count
}

func compareAnswers(expected, actual []int) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("query count mismatch: expected %d answers, got %d", len(expected), len(actual))
	}
	for i := range expected {
		if expected[i] != actual[i] {
			return fmt.Errorf("answer %d mismatch: expected %d, got %d", i+1, expected[i], actual[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
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
	for idx, tc := range tests {
		input := buildInput(tc)
		expectedOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		actualOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		expectedVals, err := parseOutput(expectedOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\noutput:\n%s", idx+1, err, expectedOut)
			os.Exit(1)
		}
		actualVals, err := parseOutput(actualOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s", idx+1, err, actualOut)
			os.Exit(1)
		}
		queryCount := filterQueries(tc.ops)
		if len(expectedVals) != queryCount {
			fmt.Fprintf(os.Stderr, "oracle output mismatch on test %d: expected %d answers, got %d\n", idx+1, queryCount, len(expectedVals))
			os.Exit(1)
		}
		if err := compareAnswers(expectedVals, actualVals); err != nil {
			fmt.Fprintf(os.Stderr, "test %d mismatch: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed.")
}
