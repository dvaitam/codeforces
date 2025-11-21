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
	refSource   = "2000-2999/2100-2199/2100-2109/2101/2101E.go"
	maxTests    = 400
	totalNLimit = 6000
)

type testCase struct {
	n int
	s string
	e [][2]int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	refAns, err := parseOutput(refOut, tests)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse reference output:", err)
		os.Exit(1)
	}

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}
	candAns, err := parseOutput(candOut, tests)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse candidate output:", err)
		os.Exit(1)
	}

	for i := range tests {
		for j := 0; j < tests[i].n; j++ {
			if candAns[i][j] != refAns[i][j] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d node %d: expected %d, got %d\n", i+1, j+1, refAns[i][j], candAns[i][j])
				os.Exit(1)
			}
		}
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2101E-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		out.WriteString(errBuf.String())
		return out.String(), err
	}
	if errBuf.Len() > 0 {
		out.WriteString(errBuf.String())
	}
	return out.String(), nil
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d\n", tc.n)
		fmt.Fprintf(&b, "%s\n", tc.s)
		for _, e := range tc.e {
			fmt.Fprintf(&b, "%d %d\n", e[0], e[1])
		}
	}
	return b.String()
}

func parseOutput(out string, tests []testCase) ([][]int64, error) {
	fields := strings.Fields(out)
	totalNeed := 0
	for _, tc := range tests {
		totalNeed += tc.n
	}
	if len(fields) != totalNeed {
		return nil, fmt.Errorf("expected %d integers, got %d", totalNeed, len(fields))
	}
	res := make([][]int64, len(tests))
	idx := 0
	for i, tc := range tests {
		res[i] = make([]int64, tc.n)
		for j := 0; j < tc.n; j++ {
			val, err := strconv.ParseInt(fields[idx], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed parsing number %q at position %d: %v", fields[idx], idx+1, err)
			}
			res[i][j] = val
			idx++
		}
	}
	return res, nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	sumN := 0

	add := func(tc testCase) {
		if len(tests) >= maxTests || sumN+tc.n > totalNLimit {
			return
		}
		tests = append(tests, tc)
		sumN += tc.n
	}

	// Edge cases.
	add(makeLineCase(1, "0"))
	add(makeLineCase(1, "1"))
	add(makeLineCase(2, "10"))
	add(makeLineCase(2, "11"))

	for len(tests) < maxTests && sumN < totalNLimit {
		n := rng.Intn(400) + 1
		if sumN+n > totalNLimit {
			n = totalNLimit - sumN
		}
		s := make([]byte, n)
		hasOne := false
		for i := 0; i < n; i++ {
			if rng.Intn(5) == 0 {
				s[i] = '1'
				hasOne = true
			} else {
				s[i] = '0'
			}
		}
		if !hasOne {
			s[rng.Intn(n)] = '1'
		}
		edges := randomTree(n, rng)
		add(testCase{n: n, s: string(s), e: edges})
	}

	if len(tests) == 0 {
		tests = append(tests, makeLineCase(1, "1"))
	}
	return tests
}

func makeLineCase(n int, s string) testCase {
	edges := make([][2]int, 0, n-1)
	for i := 1; i < n; i++ {
		edges = append(edges, [2]int{i, i + 1})
	}
	return testCase{n: n, s: s, e: edges}
}

func randomTree(n int, rng *rand.Rand) [][2]int {
	edges := make([][2]int, 0, n-1)
	for v := 2; v <= n; v++ {
		parent := rng.Intn(v-1) + 1
		edges = append(edges, [2]int{parent, v})
	}
	return edges
}
