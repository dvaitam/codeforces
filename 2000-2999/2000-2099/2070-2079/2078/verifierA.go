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

const refSource = "2078A.go"

type testCase struct {
	n   int
	x   int
	arr []int
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %s got %s\n", i+1, refAns[i], candAns[i])
			fmt.Fprintf(os.Stderr, "n=%d x=%d arr=%v\n", tc.n, tc.x, tc.arr)
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
	tmp, err := os.CreateTemp("", "ref_2078A_*.bin")
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

	// Fixed small cases
	add(testCase{n: 1, x: 3, arr: []int{3}})
	add(testCase{n: 2, x: 5, arr: []int{5, 5}})
	add(testCase{n: 3, x: 2, arr: []int{1, 2, 3}})
	add(testCase{n: 4, x: 1, arr: []int{1, 2, 1, 0}})
	add(testCase{n: 5, x: 9, arr: []int{9, 9, 9, 9, 9}})
	add(testCase{n: 6, x: 7, arr: []int{1, 2, 3, 4, 5, 6}})
	add(testCase{n: 10, x: 10, arr: []int{10, 10, 10, 10, 10, 10, 10, 10, 10, 10}})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	makeArr := func(n, x int, sumMatch bool) []int {
		arr := make([]int, n)
		curSum := 0
		for i := 0; i < n; i++ {
			v := rng.Intn(100) + 1
			arr[i] = v
			curSum += v
		}
		if sumMatch {
			diff := n*x - curSum
			// adjust by shifting first element within bounds when possible
			arr[0] += diff
			if arr[0] <= 0 {
				arr[0] = 1
			}
		}
		return arr
	}

	for len(tests) < 200 {
		n := rng.Intn(100) + 1
		x := rng.Intn(100) + 1
		sumMatch := rng.Intn(2) == 0
		arr := makeArr(n, x, sumMatch)
		tests = append(tests, testCase{n: n, x: x, arr: arr})
	}

	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		if len(tc.arr) != tc.n {
			panic("array length mismatch")
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.x))
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

func parseOutput(out string, expected int) ([]string, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d outputs, got %d", expected, len(fields))
	}
	ans := make([]string, expected)
	for i, s := range fields {
		ans[i] = normalizeAnswer(s)
	}
	return ans, nil
}

func normalizeAnswer(s string) string {
	if strings.EqualFold(s, "YES") {
		return "YES"
	}
	return "NO"
}
