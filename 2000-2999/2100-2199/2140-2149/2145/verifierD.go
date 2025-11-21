package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const refSource = "2000-2999/2100-2199/2140-2149/2145/2145D.go"

type testCase struct {
	n, k int
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
	possible, err := parseReference(refOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\noutput:\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	if err := verifyCandidate(candOut, tests, possible); err != nil {
		fmt.Fprintf(os.Stderr, "verification failed: %v\ncandidate output:\n%s", err, candOut)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	outPath := "./ref_2145D.bin"
	cmd := exec.Command("go", "build", "-o", outPath, refSource)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return outPath, nil
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
	return stdout.String(), nil
}

func buildTests() []testCase {
	var tests []testCase
	add := func(n, k int) {
		tests = append(tests, testCase{n: n, k: k})
	}
	// small exhaustive combinations
	for n := 2; n <= 6; n++ {
		maxK := n * (n - 1) / 2
		for k := 0; k <= maxK; k++ {
			add(n, k)
		}
	}
	// specific edge cases
	add(7, 0)
	add(7, 21)
	add(8, 10)
	add(10, 0)
	add(10, 20)
	add(15, 50)
	add(20, 100)
	add(25, 200)
	add(30, 400)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 120 {
		n := rng.Intn(29) + 2
		maxK := n * (n - 1) / 2
		k := rng.Intn(maxK + 1)
		add(n, k)
	}
	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
	}
	return sb.String()
}

func parseReference(out string, tests []testCase) ([]bool, error) {
	tokens := strings.Fields(out)
	idx := 0
	possible := make([]bool, len(tests))
	for i, tc := range tests {
		if idx >= len(tokens) {
			return nil, fmt.Errorf("not enough tokens for test %d", i+1)
		}
		if tokens[idx] == "0" {
			possible[i] = false
			idx++
			continue
		}
		possible[i] = true
		if len(tokens)-idx < tc.n {
			return nil, fmt.Errorf("test %d: reference output truncated", i+1)
		}
		idx += tc.n
	}
	if idx != len(tokens) {
		return nil, fmt.Errorf("reference output has extra tokens")
	}
	return possible, nil
}

func verifyCandidate(out string, tests []testCase, possible []bool) error {
	tokens := strings.Fields(out)
	idx := 0
	for tIdx, tc := range tests {
		if idx >= len(tokens) {
			return fmt.Errorf("not enough output for test %d", tIdx+1)
		}
		if !possible[tIdx] {
			if tokens[idx] != "0" {
				return fmt.Errorf("test %d: expected 0, got %s", tIdx+1, tokens[idx])
			}
			idx++
			continue
		}
		if tokens[idx] == "0" {
			return fmt.Errorf("test %d: solution exists, but candidate printed 0", tIdx+1)
		}
		if len(tokens)-idx < tc.n {
			return fmt.Errorf("test %d: insufficient numbers for permutation", tIdx+1)
		}
		perm := make([]int, tc.n)
		for i := 0; i < tc.n; i++ {
			val, err := strconv.Atoi(tokens[idx+i])
			if err != nil {
				return fmt.Errorf("test %d: invalid integer %q", tIdx+1, tokens[idx+i])
			}
			perm[i] = val
		}
		idx += tc.n
		if err := checkPermutation(perm, tc); err != nil {
			return fmt.Errorf("test %d: %v", tIdx+1, err)
		}
	}
	if idx != len(tokens) {
		return fmt.Errorf("extra output detected")
	}
	return nil
}

func checkPermutation(perm []int, tc testCase) error {
	n := tc.n
	if len(perm) != n {
		return fmt.Errorf("expected length %d, got %d", n, len(perm))
	}
	seen := make([]bool, n+1)
	for _, v := range perm {
		if v < 1 || v > n {
			return fmt.Errorf("value %d out of range [1,%d]", v, n)
		}
		if seen[v] {
			return fmt.Errorf("value %d repeated", v)
		}
		seen[v] = true
	}
	inv := inversionValue(perm)
	if inv != tc.k {
		return fmt.Errorf("inversion value mismatch: expected %d, got %d", tc.k, inv)
	}
	return nil
}

func inversionValue(perm []int) int {
	n := len(perm)
	total := 0
	for l := 0; l < n; l++ {
		for r := l + 1; r < n; r++ {
			if hasInversion(perm[l : r+1]) {
				total++
			}
		}
	}
	return total
}

func hasInversion(sub []int) bool {
	for i := 0; i < len(sub); i++ {
		for j := i + 1; j < len(sub); j++ {
			if sub[i] > sub[j] {
				return true
			}
		}
	}
	return false
}
