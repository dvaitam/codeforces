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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
		t, err := countCases(tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse generated test %q: %v\n", tc.name, err)
			os.Exit(1)
		}

		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refVals, err := parseOutputs(refOut, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "solution runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		candVals, err := parseOutputs(candOut, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse solution output on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}

		for i := 0; i < t; i++ {
			if refVals[i] != candVals[i] {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed on case %d: expected %d, got %d\ninput:\n%soutput:\n%s", idx+1, tc.name, i+1, refVals[i], candVals[i], tc.input, candOut)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("Accepted (%d tests)\n", len(tests))
}

func buildReferenceBinary() (string, error) {
	tmp, err := os.CreateTemp("", "cf-1592A-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), "1592A.go")
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

func countCases(input string) (int, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return 0, err
	}
	return t, nil
}

func parseOutputs(output string, expected int) ([]int64, error) {
	reader := bufio.NewReader(strings.NewReader(output))
	res := make([]int64, expected)
	for i := 0; i < expected; i++ {
		if _, err := fmt.Fscan(reader, &res[i]); err != nil {
			return nil, fmt.Errorf("reading answer %d: %v", i+1, err)
		}
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err != nil {
		if err != io.EOF {
			return nil, fmt.Errorf("failed to parse trailing output: %v", err)
		}
	} else {
		return nil, fmt.Errorf("unexpected extra token %q", extra)
	}
	return res, nil
}

func generateTests() []testCase {
	var tests []testCase

	// manual cases
	tests = append(tests, buildManualTest("manual_small", [][]int{
		{2, 4, 7, 2},
		{3, 6, 4, 2, 1},
		{3, 11, 2, 7, 3},
	}))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTest("random_small", 20, 10, rng))
	tests = append(tests, randomTest("random_medium", 50, 200, rng))
	tests = append(tests, randomTest("random_large", 100, 1000, rng))

	return tests
}

// manual tests: for each entry: [n, H, a1, a2,...]
func buildManualTest(name string, entries [][]int) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(entries))
	for _, e := range entries {
		n := e[0]
		H := e[1]
		fmt.Fprintf(&b, "%d %d\n", n, H)
		for i := 0; i < n; i++ {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", e[2+i])
		}
		b.WriteByte('\n')
	}
	return testCase{name: name, input: b.String()}
}

func randomTest(name string, caseCount int, maxN int, rng *rand.Rand) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", caseCount)
	totalN := 0
	for i := 0; i < caseCount; i++ {
		remaining := 200000 - totalN
		if remaining < 2 {
			break
		}
		n := rng.Intn(min(maxN, remaining-1)-1) + 2
		H := rng.Int63n(1_000_000_000) + 1
		fmt.Fprintf(&b, "%d %d\n", n, H)
		for j := 0; j < n; j++ {
			if j > 0 {
				b.WriteByte(' ')
			}
			val := rng.Intn(1_000_000_000) + 1
			fmt.Fprintf(&b, "%d", val)
		}
		b.WriteByte('\n')
		totalN += n
	}
	return testCase{name: name, input: b.String()}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
