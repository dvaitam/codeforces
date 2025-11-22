package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "2000-2999/2100-2199/2110-2119/2114/2114F.go"

type testCase struct {
	input string
}

type trio struct {
	x, y, k int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
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

		if !equalTokens(expect, got) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, tc.input, expect, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2114F-ref-*")
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

func equalTokens(expected, got string) bool {
	ta := strings.Fields(expected)
	tb := strings.Fields(got)
	if len(ta) != len(tb) {
		return false
	}
	for i := range ta {
		if ta[i] != tb[i] {
			return false
		}
	}
	return true
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(21142114))
	var tests []testCase

	// Single-case sanity checks.
	tests = append(tests, makeTest([]trio{{4, 6, 3}}))
	tests = append(tests, makeTest([]trio{{4, 5, 3}}))
	tests = append(tests, makeTest([]trio{{4, 6, 2}}))
	tests = append(tests, makeTest([]trio{{10, 45, 3}}))
	tests = append(tests, makeTest([]trio{{780, 23, 4}, {211, 270, 23}, {1, 98, 2}}))
	tests = append(tests, makeTest([]trio{{800, 131, 6}, {2, 2, 1}, {17, 17, 1}}))

	// Random small cases.
	for i := 0; i < 40; i++ {
		t := rng.Intn(6) + 1
		var arr []trio
		for j := 0; j < t; j++ {
			x := rng.Intn(100) + 1
			y := rng.Intn(100) + 1
			k := rng.Intn(10) + 1
			arr = append(arr, trio{x, y, k})
		}
		tests = append(tests, makeTest(arr))
	}

	// Cases with k = 1 (only possible when x == y).
	for i := 0; i < 10; i++ {
		x := rng.Intn(200) + 1
		tests = append(tests, makeTest([]trio{{x, x, 1}, {x, x + 1, 1}}))
	}

	// Cases with large k enabling many divisors.
	for i := 0; i < 20; i++ {
		t := rng.Intn(5) + 3
		var arr []trio
		for j := 0; j < t; j++ {
			x := rng.Intn(1_000_000) + 1
			y := rng.Intn(1_000_000) + 1
			k := rng.Intn(1_000_000) + 1
			arr = append(arr, trio{x, y, k})
		}
		tests = append(tests, makeTest(arr))
	}

	// Constructed cases with common gcd tweaks.
	tests = append(tests, makeTest([]trio{
		{12 * 7, 84, 6},
		{999983, 999983 * 6, 6},
		{720720, 99991, 30},
	}))

	// Max range stress.
	var big []trio
	for i := 0; i < 50; i++ {
		x := rng.Intn(1_000_000) + 1
		y := rng.Intn(1_000_000) + 1
		k := rng.Intn(1_000_000) + 1
		big = append(big, trio{x, y, k})
	}
	tests = append(tests, makeTest(big))

	return tests
}

func makeTest(arr []trio) testCase {
	return testCase{input: buildInput(arr)}
}

func buildInput(arr []trio) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(arr))
	for _, t := range arr {
		fmt.Fprintf(&b, "%d %d %d\n", t.x, t.y, t.k)
	}
	return b.String()
}
