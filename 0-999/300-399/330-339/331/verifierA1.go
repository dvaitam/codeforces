package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	refSourceA1 = "331A1.go"
	refBinaryA1 = "ref331A1.bin"
	totalTests  = 80
	maxN        = 60
)

type testCase struct {
	n int
	a []int64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA1.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	refPath, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		return
	}
	defer os.Remove(refPath)

	tests := generateTests()

	for i, tc := range tests {
		input := formatInput(tc)

		refOut, err := runProgram(refPath, input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			return
		}
		refSum, err := parseReference(refOut)
		if err != nil {
			fmt.Printf("failed to parse reference output on test %d: %v\noutput:\n%s", i+1, err, refOut)
			return
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			return
		}
		sumVal, cuts, err := parseCandidate(candOut)
		if err != nil {
			fmt.Printf("candidate output parse error on test %d: %v\noutput:\n%s", i+1, err, candOut)
			fmt.Println("Input used:")
			fmt.Println(string(input))
			return
		}

		if err := verifySolution(tc, refSum, sumVal, cuts); err != nil {
			fmt.Printf("test %d failed: %v\n", i+1, err)
			fmt.Println("Input used:")
			fmt.Println(string(input))
			return
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinaryA1, refSourceA1)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinaryA1), nil
}

func runProgram(path string, input []byte) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func formatInput(tc testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", tc.n)
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	return []byte(sb.String())
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := []testCase{
		{n: 2, a: []int64{1, 1}},
		{n: 3, a: []int64{5, -2, 5}},
		{n: 4, a: []int64{-3, -3, -3, -3}},
	}
	for len(tests) < totalTests {
		n := rnd.Intn(maxN-1) + 2
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			arr[i] = int64(rnd.Intn(2001) - 1000)
		}
		// ensure at least two equal values
		if !hasDuplicate(arr) {
			if n >= 2 {
				arr[n-1] = arr[0]
			}
		}
		tests = append(tests, testCase{n: n, a: arr})
	}
	return tests
}

func hasDuplicate(arr []int64) bool {
	seen := make(map[int64]struct{})
	for _, v := range arr {
		if _, ok := seen[v]; ok {
			return true
		}
		seen[v] = struct{}{}
	}
	return false
}

func parseReference(out string) (int64, error) {
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) == 0 || len(strings.TrimSpace(lines[0])) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	fields := strings.Fields(lines[0])
	if len(fields) < 1 {
		return 0, fmt.Errorf("missing sum")
	}
	sumVal, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, err
	}
	return sumVal, nil
}

func parseCandidate(out string) (int64, []int, error) {
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) == 0 || len(strings.TrimSpace(lines[0])) == 0 {
		return 0, nil, fmt.Errorf("empty output")
	}
	fields := strings.Fields(lines[0])
	if len(fields) < 2 {
		return 0, nil, fmt.Errorf("first line must contain sum and k")
	}
	sumVal, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, nil, fmt.Errorf("invalid sum: %v", err)
	}
	k, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, nil, fmt.Errorf("invalid k: %v", err)
	}
	var tokens []string
	if len(lines) > 1 {
		tokens = strings.Fields(strings.Join(lines[1:], " "))
	}
	if k == 0 {
		if len(tokens) > 0 {
			return 0, nil, fmt.Errorf("expected no indices for k=0")
		}
		return sumVal, nil, nil
	}
	if len(tokens) != k {
		return 0, nil, fmt.Errorf("expected %d indices, got %d", k, len(tokens))
	}
	cuts := make([]int, k)
	for i := 0; i < k; i++ {
		idx, err := strconv.Atoi(tokens[i])
		if err != nil {
			return 0, nil, fmt.Errorf("invalid index %q", tokens[i])
		}
		cuts[i] = idx
	}
	return sumVal, cuts, nil
}

func verifySolution(tc testCase, refSum, sumVal int64, cuts []int) error {
	if sumVal != refSum {
		return fmt.Errorf("sum mismatch: expected %d, got %d", refSum, sumVal)
	}
	n := tc.n
	removed := make([]bool, n)
	for _, idx := range cuts {
		if idx < 1 || idx > n {
			return fmt.Errorf("index %d out of range", idx)
		}
		if removed[idx-1] {
			return fmt.Errorf("index %d repeated", idx)
		}
		removed[idx-1] = true
	}
	left := make([]int64, 0, n-len(cuts))
	for i := 0; i < n; i++ {
		if !removed[i] {
			left = append(left, tc.a[i])
		}
	}
	if len(left) < 2 {
		return fmt.Errorf("less than two trees remain")
	}
	if left[0] != left[len(left)-1] {
		return fmt.Errorf("first and last values differ")
	}
	sum := int64(0)
	for _, v := range left {
		sum += v
	}
	if sum != sumVal {
		return fmt.Errorf("reported sum %d does not match actual %d", sumVal, sum)
	}
	return nil
}
