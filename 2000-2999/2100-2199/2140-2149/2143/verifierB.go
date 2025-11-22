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

const refSource = "2000-2999/2100-2199/2140-2149/2143/2143B.go"

type testCase struct {
	input string
}

type instance struct {
	n int
	k int
	a []int64
	b []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
	tmp, err := os.CreateTemp("", "2143B-ref-*")
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
	rng := rand.New(rand.NewSource(21432143))
	var tests []testCase

	// Example-inspired small case.
	tests = append(tests, makeTest([]instance{
		{
			n: 5, k: 3,
			a: []int64{18, 3, 7, 2, 9},
			b: []int{3, 1, 1},
		},
	}))

	// Single voucher, minimal n.
	tests = append(tests, makeTest([]instance{
		{
			n: 3, k: 1,
			a: []int64{5, 1, 4},
			b: []int{2},
		},
	}))

	// All vouchers of size 1.
	tests = append(tests, makeTest([]instance{
		{
			n: 6, k: 4,
			a: []int64{10, 20, 30, 40, 50, 60},
			b: []int{1, 1, 1, 1},
		},
	}))

	// Random moderate instances with controlled sums.
	for i := 0; i < 20; i++ {
		t := rng.Intn(4) + 1
		var insts []instance
		totalN := 0
		for j := 0; j < t; j++ {
			n := rng.Intn(50) + 5
			k := rng.Intn(n) + 1
			insts = append(insts, randomInstance(rng, n, k))
			totalN += n
			if totalN > 3000 {
				break
			}
		}
		tests = append(tests, makeTest(insts))
	}

	// Larger stress cases.
	tests = append(tests, makeTest([]instance{
		randomInstance(rng, 2000, 1500),
	}))
	tests = append(tests, makeTest([]instance{
		randomInstance(rng, 5000, 4000),
	}))

	return tests
}

func randomInstance(rng *rand.Rand, n, k int) instance {
	// Ensure k <= n and sum(b) <= n to keep reference program safe.
	if k > n {
		k = n
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Int63n(1_000_000_000) + 1
	}
	target := rng.Intn(n-k+1) + k // between k and n inclusive
	parts := make([]int, k)
	remaining := target
	for i := 0; i < k; i++ {
		maxVal := remaining - (k - i - 1)
		val := 1
		if maxVal > 1 {
			val = rng.Intn(maxVal) + 1
		}
		parts[i] = val
		remaining -= val
	}
	// Distribute any leftover (remaining should be 0).
	parts[0] += remaining

	return instance{n: n, k: k, a: a, b: parts}
}

func makeTest(insts []instance) testCase {
	return testCase{input: buildInput(insts)}
}

func buildInput(insts []instance) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(insts))
	for _, in := range insts {
		fmt.Fprintf(&b, "%d %d\n", in.n, in.k)
		for i, v := range in.a {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
		for i, v := range in.b {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
	}
	return b.String()
}
