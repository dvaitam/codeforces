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

const refSource2141F = "2000-2999/2100-2199/2140-2149/2141/2141F.go"

type testCase struct {
	n   int
	arr []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
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
		if len(expected[i]) != len(got[i]) {
			fmt.Fprintf(os.Stderr, "test %d: expected %d numbers, got %d\ninput:\n%s", i+1, len(expected[i]), len(got[i]), stringifyCase(tests[i]))
			os.Exit(1)
		}
		for j := range expected[i] {
			if expected[i][j] != got[i][j] {
				fmt.Fprintf(os.Stderr, "test %d, position %d: expected %d, got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
					i+1, j, expected[i][j], got[i][j], stringifyCase(tests[i]), refOut, candOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2141F-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2141F.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource2141F)
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

func parseOutput(out string, tests []testCase) ([][]int64, error) {
	scan := strings.Fields(out)
	res := make([][]int64, 0, len(tests))
	pos := 0
	for idx, tc := range tests {
		if pos+tc.n > len(scan) {
			return nil, fmt.Errorf("test %d: expected %d numbers, got %d", idx+1, tc.n, len(scan)-pos)
		}
		cur := make([]int64, tc.n)
		for i := 0; i < tc.n; i++ {
			v, err := strconv.ParseInt(scan[pos+i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("test %d: invalid integer at position %d: %v", idx+1, i+1, err)
			}
			cur[i] = v
		}
		res = append(res, cur)
		pos += tc.n
	}
	return res, nil
}

func formatInput(tests []testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for i, v := range tc.arr {
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
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func buildTests() []testCase {
	tests := make([]testCase, 0, 40)

	// Deterministic small arrays.
	tests = append(tests,
		testCase{n: 1, arr: []int{1}},
		testCase{n: 2, arr: []int{1, 1}},
		testCase{n: 3, arr: []int{1, 2, 3}},
		testCase{n: 4, arr: []int{2, 2, 1, 1}},
		testCase{n: 5, arr: []int{1, 2, 2, 3, 3}},
		testCase{n: 6, arr: []int{4, 4, 4, 4, 4, 4}},
		testCase{n: 7, arr: []int{1, 1, 2, 2, 3, 3, 4}},
	)

	// Random tests with moderate sizes to keep total n small.
	rng := rand.New(rand.NewSource(2141_2024))
	totalN := 0
	for len(tests) < 35 {
		n := rng.Intn(50) + 2 // 2..51
		if totalN+n > 2000 {  // keep total small for verifier speed
			break
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(n) + 1
		}
		tests = append(tests, testCase{n: n, arr: arr})
		totalN += n
	}

	return tests
}
