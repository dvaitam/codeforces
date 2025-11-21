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

const refSource = "2065C2.go"

type testCase struct {
	n int
	m int
	a []int64
	b []int64
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierC2.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "n=%d m=%d\na=%v\nb=%v\n", tc.n, tc.m, tc.a, tc.b)
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
	tmp, err := os.CreateTemp("", "ref_2065C2_*.bin")
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

	// Fixed scenarios
	add(testCase{n: 1, m: 1, a: []int64{5}, b: []int64{10}})
	add(testCase{n: 2, m: 2, a: []int64{5, 4}, b: []int64{9, 1}})
	add(testCase{n: 3, m: 3, a: []int64{3, 2, 1}, b: []int64{4, 5, 6}})
	add(testCase{n: 4, m: 2, a: []int64{1, 100, 2, 99}, b: []int64{101, 3}})
	add(testCase{n: 5, m: 5, a: []int64{1, 2, 3, 4, 5}, b: []int64{1, 1, 1, 1, 1}})
	add(testCase{n: 5, m: 3, a: []int64{5, 4, 3, 2, 1}, b: []int64{3, 6, 9}})
	add(testCase{n: 6, m: 2, a: []int64{10, 9, 8, 7, 6, 5}, b: []int64{11, 12}})
	add(testCase{n: 3, m: 1, a: []int64{10, 1, 9}, b: []int64{10}})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	limit := int64(1_000_000_000)

	genArray := func(lenArr int) []int64 {
		arr := make([]int64, lenArr)
		for i := 0; i < lenArr; i++ {
			arr[i] = rng.Int63n(limit) + 1
		}
		return arr
	}

	for len(tests) < 120 {
		n := rng.Intn(50) + 1
		m := rng.Intn(50) + 1
		add(testCase{n: n, m: m, a: genArray(n), b: genArray(m)})
	}

	// Stress-sized but single case to keep runtime reasonable
	largeN := 200000
	largeM := 200000
	largeA := make([]int64, largeN)
	largeB := make([]int64, largeM)
	for i := 0; i < largeN; i++ {
		largeA[i] = int64((i % 1_000_000) + 1)
	}
	for i := 0; i < largeM; i++ {
		largeB[i] = int64(limit - int64(i%10))
	}
	add(testCase{n: largeN, m: largeM, a: largeA, b: largeB})

	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		if len(tc.a) != tc.n || len(tc.b) != tc.m {
			panic("length mismatch")
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		for i, v := range tc.b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, expected int) ([]string, error) {
	lines := strings.Fields(out)
	if len(lines) != expected {
		return nil, fmt.Errorf("expected %d outputs, got %d", expected, len(lines))
	}
	for i, s := range lines {
		lines[i] = normalizeAnswer(s)
	}
	return lines, nil
}

func normalizeAnswer(s string) string {
	if strings.EqualFold(s, "YES") {
		return "YES"
	}
	return "NO"
}
