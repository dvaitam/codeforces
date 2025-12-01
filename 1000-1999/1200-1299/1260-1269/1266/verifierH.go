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
		fmt.Println("usage: go run verifierH.go /path/to/binary")
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
		expected := normalizeOutput(refOut)

		candOut, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "solution runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		got := normalizeOutput(candOut)

		if expected != got {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %s, got %s\ninput:\n%s", idx+1, tc.name, expected, got, tc.input)
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests)\n", len(tests))
}

func buildReferenceBinary() (string, error) {
	tmp, err := os.CreateTemp("", "cf-1266H-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), "1266H.go")
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

func normalizeOutput(out string) string {
	var b strings.Builder
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			if b.Len() > 0 {
				b.WriteByte('\n')
			}
			b.WriteString(line)
		}
	}
	return b.String()
}

func generateTests() []testCase {
	var tests []testCase

	tests = append(tests, makeTestCase("simple_line", 4,
		[][2]int{{2, 2}, {3, 3}, {4, 4}},
		[][2]string{
			{fmt.Sprint(1), zeroString(3, 'B')},
			{fmt.Sprint(2), "RRB"},
			{fmt.Sprint(3), "RBR"},
		},
	))

	tests = append(tests, makeTestCase("branch", 5,
		[][2]int{{2, 3}, {4, 5}, {5, 5}, {5, 5}},
		[][2]string{
			{fmt.Sprint(1), "BBBB"},
			{fmt.Sprint(2), "RRRR"},
			{fmt.Sprint(3), "RBRB"},
			{fmt.Sprint(4), "BBRR"},
			{fmt.Sprint(2), "BRRB"},
		},
	))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomCase(rng, "random_small", 6, 50))
	tests = append(tests, randomCase(rng, "random_medium", 20, 500))
	tests = append(tests, randomCase(rng, "random_large", 58, 2000))
	tests = append(tests, randomCase(rng, "random_dense", 58, 4000))

	return tests
}

func makeTestCase(name string, n int, edges [][2]int, queries [][2]string) testCase {
	if len(edges) != n-1 {
		panic("edges length mismatch")
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	for i := 0; i < n-1; i++ {
		fmt.Fprintf(&b, "%d %d\n", edges[i][0], edges[i][1])
	}
	fmt.Fprintf(&b, "%d\n", len(queries))
	for _, q := range queries {
		fmt.Fprintf(&b, "%d %s\n", q[0], q[1])
	}
	return testCase{name: name, input: b.String()}
}

func zeroString(length int, fill byte) string {
	buf := make([]byte, length)
	for i := 0; i < length; i++ {
		buf[i] = fill
	}
	return string(buf)
}

func randomCase(rng *rand.Rand, name string, n int, q int) testCase {
	edges := make([][2]int, n-1)
	for i := 1; i <= n-1; i++ {
		blue := rng.Intn(n) + 1
		red := rng.Intn(n) + 1
		edges[i-1] = [2]int{blue, red}
	}
	queries := make([][2]string, q)
	for i := 0; i < q; i++ {
		v := rng.Intn(n-1) + 1
		s := randomState(rng, n-1)
		queries[i] = [2]string{fmt.Sprint(v), s}
	}
	return makeTestCase(name, n, edges, queries)
}

func randomState(rng *rand.Rand, length int) string {
	buf := make([]byte, length)
	for i := 0; i < length; i++ {
		if rng.Intn(2) == 0 {
			buf[i] = 'B'
		} else {
			buf[i] = 'R'
		}
	}
	return string(buf)
}
