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

const refSource2137B = "2000-2999/2100-2199/2130-2139/2137/2137B.go"

type testCase struct {
	name  string
	input string
}

type caseData struct {
	n int
	p []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
	for idx, tc := range tests {
		// Ensure reference runs correctly.
		if _, err := runProgram(refBin, tc.input); err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		if err := validateOutput(tc.input, candOut); err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2137B-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2137B.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource2137B)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseInput(input string) ([]caseData, error) {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil || t < 1 || t > 10000 {
		return nil, fmt.Errorf("invalid test count %q", fields[pos])
	}
	pos++
	cases := make([]caseData, 0, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if pos >= len(fields) {
			return nil, fmt.Errorf("missing n for case %d", caseIdx+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil || n < 2 || n > 200000 {
			return nil, fmt.Errorf("invalid n %q for case %d", fields[pos], caseIdx+1)
		}
		pos++
		if pos+n > len(fields) {
			return nil, fmt.Errorf("not enough permutation values for case %d", caseIdx+1)
		}
		p := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(fields[pos+i])
			if err != nil {
				return nil, fmt.Errorf("invalid p value %q for case %d", fields[pos+i], caseIdx+1)
			}
			p[i] = val
		}
		pos += n
		cases = append(cases, caseData{n: n, p: p})
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("extra input tokens: %v", fields[pos:])
	}
	return cases, nil
}

func parseOutput(input, output string) ([][]int, error) {
	cases, err := parseInput(input)
	if err != nil {
		return nil, err
	}
	outFields := strings.Fields(output)
	pos := 0
	results := make([][]int, len(cases))
	for idx, cs := range cases {
		n := cs.n
		if pos+n > len(outFields) {
			return nil, fmt.Errorf("case %d: expected %d values, got %d", idx+1, n, len(outFields)-pos)
		}
		q := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(outFields[pos+i])
			if err != nil {
				return nil, fmt.Errorf("case %d: invalid integer %q", idx+1, outFields[pos+i])
			}
			q[i] = val
		}
		pos += n
		results[idx] = q
	}
	if pos != len(outFields) {
		return nil, fmt.Errorf("extra output tokens: %v", outFields[pos:])
	}
	return results, nil
}

func validateOutput(input, output string) error {
	cases, err := parseInput(input)
	if err != nil {
		return err
	}
	perms, err := parseOutput(input, output)
	if err != nil {
		return err
	}
	for idx, cs := range cases {
		q := perms[idx]
		if err := validatePermutation(cs, q); err != nil {
			return fmt.Errorf("case %d: %v", idx+1, err)
		}
	}
	return nil
}

func validatePermutation(cs caseData, q []int) error {
	n := cs.n
	if len(q) != n {
		return fmt.Errorf("expected %d values, got %d", n, len(q))
	}
	seen := make([]bool, n+1)
	for i, v := range q {
		if v < 1 || v > n {
			return fmt.Errorf("value %d out of range at position %d", v, i+1)
		}
		if seen[v] {
			return fmt.Errorf("duplicate value %d", v)
		}
		seen[v] = true
	}
	for i := 0; i+1 < n; i++ {
		if gcd(cs.p[i]+q[i], cs.p[i+1]+q[i+1]) < 3 {
			return fmt.Errorf("gcd condition failed at positions %d and %d: gcd(%d,%d)= %d",
				i+1, i+2, cs.p[i]+q[i], cs.p[i+1]+q[i+1], gcd(cs.p[i]+q[i], cs.p[i+1]+q[i+1]))
		}
	}
	return nil
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func buildTests() []testCase {
	tests := []testCase{
		buildSample(),
		makeManual("small_fixed", []int{2, 3, 4}),
		makeManualPermutation("reverse_perm", [][]int{{2, 1}, {3, 2, 1, 4}}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 160; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	tests = append(tests, largeTest())
	return tests
}

func buildSample() testCase {
	const sampleInput = `3
3
1 3 2
5
5 1 2 4 3
7
6 7 1 5 4 3 2
`
	return testCase{name: "statement_sample", input: sampleInput}
}

func makeManual(name string, ns []int) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(ns)))
	for _, n := range ns {
		sb.WriteString(fmt.Sprintf("%d\n", n))
		perm := make([]int, n)
		for i := 0; i < n; i++ {
			perm[i] = i + 1
		}
		writeArray(&sb, perm)
		sb.WriteByte('\n')
	}
	return testCase{name: name, input: sb.String()}
}

func makeManualPermutation(name string, perms [][]int) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(perms)))
	for _, p := range perms {
		sb.WriteString(fmt.Sprintf("%d\n", len(p)))
		writeArray(&sb, p)
		sb.WriteByte('\n')
	}
	return testCase{name: name, input: sb.String()}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rng.Intn(50) + 2
		perm := make([]int, n)
		for j := 0; j < n; j++ {
			perm[j] = j + 1
		}
		rng.Shuffle(n, func(a, b int) {
			perm[a], perm[b] = perm[b], perm[a]
		})
		sb.WriteString(fmt.Sprintf("%d\n", n))
		writeArray(&sb, perm)
		sb.WriteByte('\n')
	}
	return testCase{name: fmt.Sprintf("random_%d", idx+1), input: sb.String()}
}

func largeTest() testCase {
	n := 200000
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		perm[i] = i + 1
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	writeArray(&sb, perm)
	sb.WriteByte('\n')
	return testCase{name: "large_identity", input: sb.String()}
}

func writeArray(sb *strings.Builder, arr []int) {
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
}
