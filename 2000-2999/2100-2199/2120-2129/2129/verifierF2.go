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
)

const refSource2129F2 = "./2129F2.go"

type testCase struct {
	n    int
	perm []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF2.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	input := formatInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\noutput:\n%s", err, refOut)
		os.Exit(1)
	}
	expected, err := parseOutput(refOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference output parse error: %v\noutput:\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s", err, candOut)
		os.Exit(1)
	}
	got, err := parseOutput(candOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate output parse error: %v\noutput:\n%s", err, candOut)
		os.Exit(1)
	}

	for i := range tests {
		exp := expected[i]
		res := got[i]
		if len(exp) != len(res) {
			fmt.Fprintf(os.Stderr, "test %d: permutation length mismatch: expected %d, got %d\ninput:\n%s", i+1, len(exp), len(res), stringifyCase(tests[i]))
			os.Exit(1)
		}
		for j := range exp {
			if exp[j] != res[j] {
				fmt.Fprintf(os.Stderr, "test %d: mismatch at position %d: expected %d, got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
					i+1, j+1, exp[j], res[j], stringifyCase(tests[i]), refOut, candOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2129F2-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2129F2.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource2129F2)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		out.Write(errBuf.Bytes())
	}
	return out.String(), err
}

func parseOutput(out string, tests []testCase) ([][]int, error) {
	tokens := strings.Fields(out)
	res := make([][]int, 0, len(tests))
	idx := 0
	for len(res) < len(tests) && idx < len(tokens) {
		if tokens[idx] != "!" {
			idx++
			continue
		}
		tc := tests[len(res)]
		if idx+tc.n >= len(tokens) {
			return nil, fmt.Errorf("test %d: not enough numbers after '!'", len(res)+1)
		}
		perm := make([]int, tc.n)
		for i := 0; i < tc.n; i++ {
			v, err := strconv.Atoi(tokens[idx+1+i])
			if err != nil {
				return nil, fmt.Errorf("test %d: invalid integer at position %d: %v", len(res)+1, i+1, err)
			}
			perm[i] = v
		}
		res = append(res, perm)
		idx += tc.n + 1
	}
	if len(res) != len(tests) {
		return nil, fmt.Errorf("expected %d answers, parsed %d", len(tests), len(res))
	}
	return res, nil
}

func formatInput(tests []testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for i, v := range tc.perm {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return []byte(sb.String())
}

func stringifyCase(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", tc.n)
	for i, v := range tc.perm {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func buildTests() []testCase {
	tests := make([]testCase, 0, 50)

	// Small deterministic permutations.
	tests = append(tests,
		testCase{n: 2, perm: []int{1, 2}},
		testCase{n: 2, perm: []int{2, 1}},
		testCase{n: 3, perm: []int{3, 1, 2}},
		testCase{n: 4, perm: []int{4, 3, 2, 1}},
		testCase{n: 5, perm: []int{2, 4, 1, 5, 3}},
	)

	// Random permutations within n <= 50 to keep input compact.
	rng := rand.New(rand.NewSource(2129_2024))
	for len(tests) < 40 {
		n := rng.Intn(50) + 2 // 2..51
		perm := make([]int, n)
		for i := 0; i < n; i++ {
			perm[i] = i + 1
		}
		for i := n - 1; i > 0; i-- {
			j := rng.Intn(i + 1)
			perm[i], perm[j] = perm[j], perm[i]
		}
		tests = append(tests, testCase{n: n, perm: perm})
	}

	return tests
}
