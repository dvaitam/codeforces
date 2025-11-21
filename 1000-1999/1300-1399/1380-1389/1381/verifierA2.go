package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type testCase struct {
	name  string
	input string
}

type caseData struct {
	n int
	a string
	b string
}

var problemDir string

func init() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("unable to determine verifier location")
	}
	problemDir = filepath.Dir(file)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA2.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()

	for idx, tc := range tests {
		cases, err := parseInput(tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse generated test %q: %v\n", tc.name, err)
			os.Exit(1)
		}

		if _, err := runProgram(refBin, tc.input); err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}

		candOut, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "solution runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}

		if err := verifyOutput(cases, candOut); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests)\n", len(tests))
}

func buildReferenceBinary() (string, error) {
	tmp, err := os.CreateTemp("", "cf-1381A2-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), "1381A2.go")
	cmd.Dir = problemDir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("go build error: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseInput(input string) ([]caseData, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, fmt.Errorf("read t: %w", err)
	}
	cases := make([]caseData, t)
	totalN := 0
	for i := 0; i < t; i++ {
		var n int
		var a, b string
		if _, err := fmt.Fscan(reader, &n); err != nil {
			return nil, fmt.Errorf("case %d: read n: %w", i+1, err)
		}
		if _, err := fmt.Fscan(reader, &a); err != nil {
			return nil, fmt.Errorf("case %d: read a: %w", i+1, err)
		}
		if _, err := fmt.Fscan(reader, &b); err != nil {
			return nil, fmt.Errorf("case %d: read b: %w", i+1, err)
		}
		if len(a) != n || len(b) != n {
			return nil, fmt.Errorf("case %d: string length mismatch", i+1)
		}
		cases[i] = caseData{n: n, a: a, b: b}
		totalN += n
	}
	if totalN > 100000 {
		return nil, fmt.Errorf("total n exceeds limits")
	}
	return cases, nil
}

func verifyOutput(cases []caseData, output string) error {
	reader := bufio.NewReader(strings.NewReader(output))
	for idx, tc := range cases {
		var k int
		if _, err := fmt.Fscan(reader, &k); err != nil {
			return fmt.Errorf("case %d: failed to read k: %v", idx+1, err)
		}
		if k < 0 || k > 2*tc.n {
			return fmt.Errorf("case %d: k=%d is outside [0, %d]", idx+1, k, 2*tc.n)
		}
		curr := []byte(tc.a)
		for i := 0; i < k; i++ {
			var p int
			if _, err := fmt.Fscan(reader, &p); err != nil {
				return fmt.Errorf("case %d: failed to read operation %d: %v", idx+1, i+1, err)
			}
			if p < 1 || p > tc.n {
				return fmt.Errorf("case %d: operation %d has invalid prefix %d", idx+1, i+1, p)
			}
			applyPrefix(curr, p)
		}
		if string(curr) != tc.b {
			return fmt.Errorf("case %d: final string %q does not match target %q", idx+1, string(curr), tc.b)
		}
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err != nil {
		if err != io.EOF {
			return fmt.Errorf("failed to parse trailing output: %v", err)
		}
	} else {
		return fmt.Errorf("unexpected extra token %q after processing all test cases", extra)
	}
	return nil
}

func applyPrefix(s []byte, p int) {
	for l, r := 0, p-1; l <= r; l, r = l+1, r-1 {
		left := invert(s[l])
		right := invert(s[r])
		s[l], s[r] = right, left
	}
}

func invert(b byte) byte {
	if b == '0' {
		return '1'
	}
	return '0'
}

func generateTests() []testCase {
	var tests []testCase

	tests = append(tests, buildTest("simple_cases", [][]string{
		{"1", "0", "1"},
		{"2", "00", "11"},
		{"3", "010", "101"},
	}))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTest("random_small", 20, 50, rng))
	tests = append(tests, randomTest("random_medium", 10, 5000, rng))
	tests = append(tests, randomTest("random_large", 5, 20000, rng))

	return tests
}

func buildTest(name string, entries [][]string) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(entries))
	for _, e := range entries {
		fmt.Fprintf(&b, "%s\n%s\n%s\n", e[0], e[1], e[2])
	}
	return testCase{name: name, input: b.String()}
}

func randomTest(name string, cases int, maxN int, rng *rand.Rand) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", cases)
	total := 0
	for i := 0; i < cases; i++ {
		n := rng.Intn(maxN) + 1
		if total+n > 100000 {
			n = 100000 - total
			if n <= 0 {
				n = 1
			}
		}
		total += n
		a := randomBinary(n, rng)
		bStr := randomBinary(n, rng)
		fmt.Fprintf(&b, "%d\n%s\n%s\n", n, a, bStr)
	}
	return testCase{name: name, input: b.String()}
}

func randomBinary(n int, rng *rand.Rand) string {
	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			buf[i] = '0'
		} else {
			buf[i] = '1'
		}
	}
	return string(buf)
}
