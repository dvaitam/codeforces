package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refSource = "./2125D.go"

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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
	for i, tc := range tests {
		exp, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}

		got, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%soutput:\n%s", i+1, err, tc.input, got)
			os.Exit(1)
		}

		if !equalTokens(exp, got) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, tc.input, exp, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2125D-ref-*")
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

func equalTokens(a, b string) bool {
	ta := strings.Fields(a)
	tb := strings.Fields(b)
	if len(ta) != 1 || len(tb) != 1 {
		return strings.TrimSpace(a) == strings.TrimSpace(b)
	}
	return ta[0] == tb[0]
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase

	// Statement samples
	tests = append(tests, testCase{input: "3 3\n1 2 1 3\n3 3 1 2\n1 3 2 3\n"})
	tests = append(tests, testCase{input: "2 3\n1 2 1 2\n2 3 1 2\n"})
	tests = append(tests, testCase{input: "8 5\n1 3 1 2\n1 5 1 6\n1 4 4 5\n5 5 1 7\n4 5 1 2\n4 5 2 5\n3 3 2 7\n1 2 1 3\n"})

	// Small deterministic
	tests = append(tests, testCase{input: "1 1\n1 1 1 2\n"})
	tests = append(tests, testCase{input: "1 4\n2 3 1 5\n"})

	// Random mid-size
	for i := 0; i < 40; i++ {
		tests = append(tests, randomCase(rng, rng.Intn(6)+1, rng.Intn(6)+1))
	}

	// Larger stress but still safe
	tests = append(tests, randomCase(rng, 2000, 2000))
	tests = append(tests, randomCase(rng, 50000, 50000))

	return tests
}

func randomCase(rng *rand.Rand, n, m int) testCase {
	if n < 1 {
		n = 1
	}
	if m < 1 {
		m = 1
	}
	segments := rng.Intn(n) + 1
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d\n", segments, m)
	for i := 0; i < segments; i++ {
		l := rng.Intn(m) + 1
		r := rng.Intn(m-l+1) + l
		p := rng.Intn(9) + 1
		q := p + rng.Intn(20) + 1
		if q >= 998244353 {
			q = 998244352
			if p >= q {
				p = q - 1
			}
		}
		fmt.Fprintf(&b, "%d %d %d %d\n", l, r, p, q)
	}
	return testCase{input: b.String()}
}
