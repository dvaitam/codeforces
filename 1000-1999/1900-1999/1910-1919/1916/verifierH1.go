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
		fmt.Println("usage: go run verifierH1.go /path/to/binary")
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
		expected, err := parseOutput(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "solution runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		got, err := parseOutput(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse solution output on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}
		if len(expected) != len(got) {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %d values, got %d\ninput:\n%soutput:\n%s", idx+1, tc.name, len(expected), len(got), tc.input, candOut)
			os.Exit(1)
		}
		for i := range expected {
			if expected[i] != got[i] {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed at index %d: expected %d, got %d\ninput:\n%soutput:\n%s", idx+1, tc.name, i, expected[i], got[i], tc.input, candOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("Accepted (%d tests)\n", len(tests))
}

func buildReferenceBinary() (string, error) {
	tmp, err := os.CreateTemp("", "cf-1916H1-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), "1916H1.go")
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

func parseOutput(output string) ([]int64, error) {
	reader := bufio.NewReader(strings.NewReader(strings.TrimSpace(output)))
	var values []int64
	for {
		var x int64
		_, err := fmt.Fscan(reader, &x)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		values = append(values, x)
	}
	return values, nil
}

func generateTests() []testCase {
	var tests []testCase

	tests = append(tests, makeTestCase("small_manual", 1, 2, 3))
	tests = append(tests, makeTestCase("medium_manual", 5, 7, 5))
	tests = append(tests, makeTestCase("k_gt_n", 3, 11, 6))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTest("random_small", rng, 10))
	tests = append(tests, randomTest("random_medium", rng, 100))
	tests = append(tests, randomTest("random_large", rng, 5000))

	return tests
}

func makeTestCase(name string, n int64, p int64, k int) testCase {
	return testCase{
		name:  name,
		input: fmt.Sprintf("%d %d %d\n", n, p, k),
	}
}

func randomPrime(upper int64, rng *rand.Rand) int64 {
	primes := []int64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97}
	return primes[rng.Intn(len(primes))]
}

func randomTest(name string, rng *rand.Rand, maxK int) testCase {
	n := rng.Int63n(1_000_000_000_000) + 1
	p := randomPrime(100, rng)
	k := rng.Intn(maxK + 1)
	return makeTestCase(name, n, p, k)
}
