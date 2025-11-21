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
		fmt.Println("usage: go run verifierE2.go /path/to/binary")
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
		expected, err := parseSingleInt(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "solution runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		got, err := parseSingleInt(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse solution output on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}

		if expected != got {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %d, got %d\ninput:\n%soutput:\n%s", idx+1, tc.name, expected, got, tc.input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests)\n", len(tests))
}

func buildReferenceBinary() (string, error) {
	tmp, err := os.CreateTemp("", "cf-1970E2-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), "1970E2.go")
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

func parseSingleInt(out string) (int64, error) {
	reader := bufio.NewReader(strings.NewReader(strings.TrimSpace(out)))
	var x int64
	if _, err := fmt.Fscan(reader, &x); err != nil {
		return 0, err
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err != nil {
		if err != io.EOF {
			return 0, fmt.Errorf("failed to parse trailing output: %v", err)
		}
	} else {
		return 0, fmt.Errorf("unexpected extra token %q", extra)
	}
	return x, nil
}

func generateTests() []testCase {
	var tests []testCase

	tests = append(tests, makeTestCase("basic", 1, 1, []int64{1}, []int64{0}))
	tests = append(tests, makeTestCase("simple2", 2, 5, []int64{1, 2}, []int64{3, 4}))
	tests = append(tests, makeTestCase("all_zero", 3, 10, []int64{0, 0, 0}, []int64{0, 0, 0}))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTest("random_small", 5, 3, rng))
	tests = append(tests, randomTest("random_medium", 50, 20, rng))
	tests = append(tests, randomTest("random_large", 100, 100, rng))

	return tests
}

func makeTestCase(name string, m int, n int64, s, l []int64) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d\n", m, n)
	for i, val := range s {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", val)
	}
	b.WriteByte('\n')
	for i, val := range l {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", val)
	}
	b.WriteByte('\n')
	return testCase{name: name, input: b.String()}
}

func randomTest(name string, maxM int, maxTrail int64, rng *rand.Rand) testCase {
	m := rng.Intn(maxM-1) + 1
	n := rng.Int63n(1_000_000_000) + 1
	s := make([]int64, m)
	l := make([]int64, m)
	for i := 0; i < m; i++ {
		s[i] = rng.Int63n(maxTrail + 1)
		l[i] = rng.Int63n(maxTrail + 1)
	}
	return makeTestCase(name, m, n, s, l)
}
