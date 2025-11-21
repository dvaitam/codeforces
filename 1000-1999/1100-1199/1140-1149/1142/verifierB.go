package main

import (
	"bytes"
	"fmt"
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
	q     int
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		expected, err := parseOutput(refOut, tc.q)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "solution runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		got, err := parseOutput(candOut, tc.q)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse solution output on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}
		if expected != got {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %s, got %s\ninput:\n%s", idx+1, tc.name, expected, got, tc.input)
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests)\n", len(tests))
}

func buildReferenceBinary() (string, error) {
	tmp, err := os.CreateTemp("", "cf-1142B-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), "1142B.go")
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

func parseOutput(output string, q int) (string, error) {
	var b strings.Builder
	for _, ch := range output {
		if ch == '0' || ch == '1' {
			b.WriteRune(ch)
		} else if ch == ' ' || ch == '\n' || ch == '\r' || ch == '\t' {
			continue
		} else {
			return "", fmt.Errorf("invalid character %q in output", ch)
		}
	}
	if b.Len() != q {
		return "", fmt.Errorf("expected %d digits, got %d", q, b.Len())
	}
	return b.String(), nil
}

func generateTests() []testCase {
	var tests []testCase

	tests = append(tests, makeTestCase(
		"sample_like",
		3, 6,
		[]int{2, 1, 3},
		[]int{1, 2, 3, 1, 2, 3},
		[][2]int{
			{1, 5},
			{2, 6},
			{3, 5},
		},
	))

	tests = append(tests, makeTestCase(
		"trivial_n1",
		1, 5,
		[]int{1},
		[]int{1, 1, 1, 1, 1},
		[][2]int{
			{1, 1},
			{1, 5},
			{2, 4},
		},
	))

	tests = append(tests, makeTestCase(
		"mixed_small",
		4, 7,
		[]int{3, 1, 4, 2},
		[]int{3, 4, 1, 2, 3, 1, 4},
		[][2]int{
			{1, 4},
			{2, 6},
			{3, 7},
			{4, 7},
		},
	))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	tests = append(tests, randomCase(rng, "random_mid", 5, 30, 25))
	tests = append(tests, randomCase(rng, "random_mid2", 8, 200, 150))
	tests = append(tests, randomCase(rng, "random_big", 20, 2000, 2000))
	tests = append(tests, randomCase(rng, "random_dense_queries", 50, 1000, 5000))
	tests = append(tests, randomCase(rng, "random_huge", 200000, 200000, 200000))

	return tests
}

func randomCase(rng *rand.Rand, name string, n, m, q int) testCase {
	p := randomPermutation(rng, n)
	a := make([]int, m)
	for i := range a {
		a[i] = rng.Intn(n) + 1
	}
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		l := rng.Intn(m) + 1
		r := l + rng.Intn(m-l+1)
		queries[i] = [2]int{l, r}
	}
	return makeTestCase(name, n, m, p, a, queries)
}

func randomPermutation(rng *rand.Rand, n int) []int {
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i + 1
	}
	for i := n - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		p[i], p[j] = p[j], p[i]
	}
	return p
}

func makeTestCase(name string, n, m int, p []int, a []int, queries [][2]int) testCase {
	if len(p) != n {
		panic("permutation length mismatch")
	}
	if len(a) != m {
		panic("array length mismatch")
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d %d\n", n, m, len(queries))
	for i, val := range p {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", val)
	}
	b.WriteByte('\n')
	for i, val := range a {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", val)
	}
	b.WriteByte('\n')
	for _, q := range queries {
		fmt.Fprintf(&b, "%d %d\n", q[0], q[1])
	}
	return testCase{
		name:  name,
		input: b.String(),
		q:     len(queries),
	}
}
