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

const refSource = "2123B.go"

type testCase struct {
	n   int
	j   int
	k   int
	arr []int
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "n=%d j=%d k=%d arr=%v\n", tc.n, tc.j, tc.k, tc.arr)
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
	tmp, err := os.CreateTemp("", "ref_2123B_*.bin")
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

	// Example-like and edge cases
	add(testCase{n: 5, j: 2, k: 3, arr: []int{3, 2, 4, 4, 1}})
	add(testCase{n: 5, j: 4, k: 1, arr: []int{5, 3, 4, 5, 2}})
	add(testCase{n: 6, j: 1, k: 1, arr: []int{1, 2, 3, 4, 5, 6}})
	add(testCase{n: 2, j: 1, k: 1, arr: []int{2, 1}})
	add(testCase{n: 2, j: 2, k: 2, arr: []int{1, 1}})
	add(testCase{n: 4, j: 3, k: 2, arr: []int{1, 4, 2, 4}})
	add(testCase{n: 4, j: 1, k: 4, arr: []int{2, 2, 2, 2}})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	totalN := 0
	for _, tc := range tests {
		totalN += tc.n
	}

	for len(tests) < 200 && totalN < 200000 {
		n := rng.Intn(3000) + 2 // 2..3001
		if totalN+n > 200000 {
			n = 200000 - totalN
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(n) + 1
		}
		j := rng.Intn(n) + 1
		k := rng.Intn(n) + 1
		add(testCase{n: n, j: j, k: k, arr: arr})
		totalN += n
	}

	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.j, tc.k))
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
	res := make([]string, expected)
	for i, s := range fields {
		res[i] = normalizeAnswer(s)
	}
	return res, nil
}

func normalizeAnswer(s string) string {
	if strings.EqualFold(s, "YES") {
		return "YES"
	}
	return "NO"
}
