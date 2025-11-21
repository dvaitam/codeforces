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

const refSource2084A = "2000-2999/2000-2099/2080-2089/2084/2084A.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
		// Run reference just to ensure input is valid; we don't compare outputs because multiple solutions are allowed.
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
	dir, err := os.MkdirTemp("", "cf-2084A-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2084A.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource2084A)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() {
		_ = os.RemoveAll(dir)
	}
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

func validateOutput(input, output string) error {
	inFields := strings.Fields(input)
	if len(inFields) == 0 {
		return fmt.Errorf("empty input")
	}
	t, err := strconv.Atoi(inFields[0])
	if err != nil || t < 1 || t > 99 {
		return fmt.Errorf("invalid test count %q", inFields[0])
	}
	nVals := make([]int, t)
	for i := 0; i < t; i++ {
		n, err := strconv.Atoi(inFields[i+1])
		if err != nil || n < 2 || n > 100 {
			return fmt.Errorf("invalid n for case %d", i+1)
		}
		nVals[i] = n
	}

	outFields := strings.Fields(output)
	pos := 0
	for idx, n := range nVals {
		if pos >= len(outFields) {
			return fmt.Errorf("missing output for case %d", idx+1)
		}
		val, err := strconv.Atoi(outFields[pos])
		if err != nil {
			return fmt.Errorf("non-integer token %q for case %d", outFields[pos], idx+1)
		}
		pos++
		if val == -1 {
			if n%2 == 1 {
				return fmt.Errorf("case %d: permutation exists for odd n=%d, but got -1", idx+1, n)
			}
			continue
		}
		if n%2 == 0 {
			return fmt.Errorf("case %d: expected -1 for even n=%d, got permutation start", idx+1, n)
		}
		perm := make([]int, n)
		perm[0] = val
		for i := 1; i < n; i++ {
			if pos >= len(outFields) {
				return fmt.Errorf("case %d: missing permutation value %d", idx+1, i+1)
			}
			v, err := strconv.Atoi(outFields[pos])
			if err != nil {
				return fmt.Errorf("case %d: invalid integer %q", idx+1, outFields[pos])
			}
			perm[i] = v
			pos++
		}
		if err := checkPermutation(perm, n); err != nil {
			return fmt.Errorf("case %d: %v", idx+1, err)
		}
	}
	if pos != len(outFields) {
		return fmt.Errorf("extra output tokens: %v", outFields[pos:])
	}
	return nil
}

func checkPermutation(perm []int, n int) error {
	seen := make([]bool, n+1)
	for i, v := range perm {
		if v < 1 || v > n {
			return fmt.Errorf("value %d out of range at position %d", v, i+1)
		}
		if seen[v] {
			return fmt.Errorf("duplicate value %d", v)
		}
		seen[v] = true
	}
	for i := 1; i < n; i++ {
		mx := perm[i-1]
		if perm[i] > mx {
			mx = perm[i]
		}
		if mx%(i+1) != i {
			return fmt.Errorf("condition failed at i=%d: max=%d mod %d != %d", i+1, mx, i+1, i)
		}
	}
	return nil
}

func buildTests() []testCase {
	tests := []testCase{
		makeManual("small_even_impossible", []int{2, 4}),
		makeManual("small_odd", []int{3, 5}),
		makeManual("mix_sample", []int{2, 3, 4, 5}),
		makeManual("upper_even", []int{100}),
		makeManual("upper_odd", []int{99}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

func makeManual(name string, ns []int) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(ns)))
	for _, n := range ns {
		sb.WriteString(fmt.Sprintf("%d\n", n))
	}
	return testCase{name: name, input: sb.String()}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rng.Intn(99) + 2
		sb.WriteString(fmt.Sprintf("%d\n", n))
	}
	return testCase{name: fmt.Sprintf("random_%d", idx+1), input: sb.String()}
}
