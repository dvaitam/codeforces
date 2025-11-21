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

const refSource = "2000-2999/2000-2099/2080-2089/2087/2087C.go"

type testCase struct {
	input string
}

type interval struct {
	l, r int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
		expected, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}

		got, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, tc.input, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, tc.input, expected, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2087C-ref-*")
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

func runCandidate(bin, input string) (string, error) {
	cmd := commandFor(bin)
	return runWithInput(cmd, input)
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(20872087))
	var tests []testCase

	tests = append(tests, sampleTest())
	tests = append(tests, makeCase("G", []interval{{1, 1}}))
	tests = append(tests, makeCase("BBBBB", []interval{{1, 5}, {2, 4}, {3, 3}}))
	tests = append(tests, makeCase("GS", []interval{{1, 1}, {2, 2}, {1, 2}}))

	for i := 0; i < 40; i++ {
		n := rng.Intn(60) + 1
		q := rng.Intn(70) + 1
		tests = append(tests, randomCase(rng, n, q))
	}

	tests = append(tests, singleTypeSweep("GSBGSBGSB", 'G'))
	tests = append(tests, singleTypeSweep("GSBGSBGSB", 'S'))
	tests = append(tests, singleTypeSweep("GSBGSBGSB", 'B'))

	tests = append(tests, randomCase(rng, 2000, 2000))
	tests = append(tests, periodicCase(20000, 20000))

	return tests
}

func sampleTest() testCase {
	return testCase{
		input: "BGSSBGB\n5\n1 7\n2 6\n1 5\n3 3\n4 7\n",
	}
}

func makeCase(s string, queries []interval) testCase {
	var b strings.Builder
	b.WriteString(s)
	b.WriteByte('\n')
	fmt.Fprintf(&b, "%d\n", len(queries))
	for _, q := range queries {
		fmt.Fprintf(&b, "%d %d\n", q.l, q.r)
	}
	return testCase{input: b.String()}
}

func randomCase(rng *rand.Rand, n, q int) testCase {
	var b strings.Builder
	coins := []byte{'G', 'S', 'B'}
	for i := 0; i < n; i++ {
		b.WriteByte(coins[rng.Intn(len(coins))])
	}
	b.WriteByte('\n')
	fmt.Fprintf(&b, "%d\n", q)
	for i := 0; i < q; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n) + 1
		if l > r {
			l, r = r, l
		}
		fmt.Fprintf(&b, "%d %d\n", l, r)
	}
	return testCase{input: b.String()}
}

func singleTypeSweep(s string, ch byte) testCase {
	var queries []interval
	for i := 0; i < len(s); i++ {
		if s[i] == ch {
			queries = append(queries, interval{i + 1, i + 1})
		}
	}
	if len(queries) == 0 {
		queries = append(queries, interval{1, len(s)})
	}
	return makeCase(s, queries)
}

func periodicCase(n, q int) testCase {
	var b strings.Builder
	pattern := []byte{'G', 'S', 'B'}
	for i := 0; i < n; i++ {
		idx := (i*37 + 13) % len(pattern)
		b.WriteByte(pattern[idx])
	}
	b.WriteByte('\n')
	fmt.Fprintf(&b, "%d\n", q)
	for i := 0; i < q; i++ {
		l := (i*7919)%n + 1
		r := (i*104729)%n + 1
		if l > r {
			l, r = r, l
		}
		fmt.Fprintf(&b, "%d %d\n", l, r)
	}
	return testCase{input: b.String()}
}
