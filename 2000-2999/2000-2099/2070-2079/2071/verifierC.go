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

const (
	refSource2071C = "2071C.go"
	refBinary2071C = "ref2071C.bin"
	maxTests       = 150
	maxTotalN      = 100000
)

type edge struct {
	u int
	v int
}

type testCase struct {
	n  int
	st int
	en int
	es []edge
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		return
	}
	defer os.Remove(ref)

	tests := generateTests()
	input := formatInput(tests)

	refOut, err := runProgram(ref, input)
	if err != nil {
		fmt.Printf("reference runtime error: %v\n", err)
		return
	}
	expected, err := parseOutput(refOut, tests)
	if err != nil {
		fmt.Printf("reference output parse error: %v\noutput:\n%s", err, refOut)
		return
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Printf("candidate runtime error: %v\n", err)
		return
	}
	got, err := parseOutput(candOut, tests)
	if err != nil {
		fmt.Printf("candidate output parse error: %v\noutput:\n%s", err, candOut)
		return
	}

	for i := range tests {
		if expected[i] != got[i] {
			fmt.Printf("Mismatch on test %d\nexpected: %s\ngot: %s\n", i+1, expected[i], got[i])
			fmt.Println("Input used:")
			fmt.Println(string(input))
			return
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary2071C, refSource2071C)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary2071C), nil
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

func parseOutput(out string, tests []testCase) ([]string, error) {
	lines := strings.FieldsFunc(strings.TrimSpace(out), func(r rune) bool { return r == '\n' || r == '\r' })
	if len(lines) != len(tests) {
		return nil, fmt.Errorf("expected %d lines, got %d", len(tests), len(lines))
	}
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			return nil, fmt.Errorf("empty line for test %d", i+1)
		}
		if line == "-1" {
			lines[i] = "-1"
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != tests[i].n {
			return nil, fmt.Errorf("test %d: expected %d numbers, got %d", i+1, tests[i].n, len(parts))
		}
		seen := make([]bool, tests[i].n+1)
		for _, p := range parts {
			v, err := strconv.Atoi(p)
			if err != nil || v < 1 || v > tests[i].n {
				return nil, fmt.Errorf("test %d: invalid vertex %q", i+1, p)
			}
			if seen[v] {
				return nil, fmt.Errorf("test %d: duplicate vertex %d", i+1, v)
			}
			seen[v] = true
		}
		lines[i] = strings.Join(parts, " ")
	}
	return lines, nil
}

func formatInput(tests []testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d %d\n", tc.n, tc.st, tc.en)
		for _, e := range tc.es {
			fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
		}
	}
	return []byte(sb.String())
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(2071))
	var tests []testCase
	totalN := 0

	addCase := func(tc testCase) {
		tests = append(tests, tc)
		totalN += tc.n
	}

	addCase(testCase{n: 1, st: 1, en: 1, es: nil})
	addCase(testCase{n: 2, st: 1, en: 2, es: []edge{{1, 2}}})
	addCase(testCase{n: 3, st: 2, en: 2, es: []edge{{1, 2}, {2, 3}}})

	for len(tests) < maxTests && totalN < maxTotalN {
		remain := maxTotalN - totalN
		if remain <= 0 {
			break
		}
		n := rnd.Intn(min(remain, 200)) + 1
		st := rnd.Intn(n) + 1
		en := rnd.Intn(n) + 1
		for en == st && rnd.Intn(3) == 0 {
			en = rnd.Intn(n) + 1
		}
		es := make([]edge, 0, n-1)
		for v := 2; v <= n; v++ {
			u := rnd.Intn(v-1) + 1
			es = append(es, edge{u: u, v: v})
		}
		addCase(testCase{n: n, st: st, en: en, es: es})
	}
	return tests
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
