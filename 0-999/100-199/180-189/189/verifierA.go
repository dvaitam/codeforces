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

const refSource = "./189A.go"

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	candidate := os.Args[1]
	tests := generateTests()

	for i, tc := range tests {
		expect, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}

		got, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, tc.input, got)
			os.Exit(1)
		}

		if !equalIntegers(expect, got) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, tc.input, expect, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "189A-ref-*")
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

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func equalIntegers(a, b string) bool {
	va, err := strconv.Atoi(strings.TrimSpace(a))
	if err != nil {
		return false
	}
	vb, err := strconv.Atoi(strings.TrimSpace(b))
	if err != nil {
		return false
	}
	return va == vb
}

func generateTests() []testCase {
	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 200; i++ {
		tests = append(tests, randomCase(rng))
	}
	return tests
}

func deterministicTests() []testCase {
	type params struct {
		n, a, b, c int
	}
	makeCounts := func(a, b, c, ca, cb, cc int) params {
		return params{
			n: ca*a + cb*b + cc*c,
			a: a, b: b, c: c,
		}
	}
	var cases []params
	cases = append(cases,
		params{n: 5, a: 5, b: 3, c: 2},
		params{n: 7, a: 5, b: 2, c: 2},
		params{n: 4000, a: 1, b: 1, c: 1},
		params{n: 4000, a: 4000, b: 4000, c: 4000},
		makeCounts(5, 7, 9, 50, 20, 10),
		makeCounts(33, 44, 55, 30, 10, 5),
		makeCounts(37, 11, 13, 10, 5, 0),
		makeCounts(14, 22, 30, 10, 0, 0),
		makeCounts(2000, 1000, 4000, 1, 2, 0),
		makeCounts(3999, 4000, 1, 1, 0, 1),
		makeCounts(7, 7, 7, 3, 0, 0),
	)

	var tests []testCase
	for _, tc := range cases {
		if tc.n == 0 {
			tc.n = tc.a
		}
		tests = append(tests, testCase{
			input: fmt.Sprintf("%d %d %d %d\n", tc.n, tc.a, tc.b, tc.c),
		})
	}
	return tests
}

func randomCase(rng *rand.Rand) testCase {
	a := rng.Intn(4000) + 1
	b := rng.Intn(4000) + 1
	c := rng.Intn(4000) + 1
	lengths := []int{a, b, c}
	total := 0

	for i := 0; i < len(lengths); i++ {
		if total == 4000 {
			break
		}
		maxAdd := (4000 - total) / lengths[i]
		if maxAdd <= 0 {
			continue
		}
		cnt := rng.Intn(maxAdd + 1)
		total += cnt * lengths[i]
	}

	if total == 0 {
		total = lengths[rng.Intn(3)]
	}

	return testCase{
		input: fmt.Sprintf("%d %d %d %d\n", total, a, b, c),
	}
}
