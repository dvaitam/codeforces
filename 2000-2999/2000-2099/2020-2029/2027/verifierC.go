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

const refSource = "2000-2999/2000-2099/2020-2029/2027/2027C.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}

	refBin, refCleanup, err := buildBinary(refSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer refCleanup()

	candBin, candCleanup, err := buildBinary(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare candidate binary: %v\n", err)
		os.Exit(1)
	}
	defer candCleanup()

	tests := generateTests()
	for idx, tc := range tests {
		expect, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}

		got, err := runBinary(candBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, got)
			os.Exit(1)
		}

		if !equalTokens(expect, got) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d (%s)\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, tc.name, tc.input, expect, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildBinary(path string) (string, func(), error) {
	cleanPath := filepath.Clean(path)
	if strings.HasSuffix(cleanPath, ".go") {
		tmp, err := os.CreateTemp("", "verifier2027C-*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), cleanPath)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("%v\n%s", err, out.String())
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	abs, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", nil, err
	}
	return abs, func() {}, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), fmt.Errorf("%v", err)
	}
	return stdout.String(), nil
}

func equalTokens(a, b string) bool {
	ta := strings.Fields(a)
	tb := strings.Fields(b)
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
	var tests []testCase
	tests = append(tests, testCase{
		name:  "sample",
		input: sampleInput(),
	})
	tests = append(tests, buildSingleCase("single", [][]int64{{1}, {42}, {7}}))
	tests = append(tests, buildSingleCase("monotonic", [][]int64{
		{2, 4, 6, 2, 5},
		{5, 4, 4, 5, 1},
		{6, 8, 2, 3, 11},
	}))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 6; i++ {
		tests = append(tests, randomBatch(fmt.Sprintf("random_small_%d", i+1), rng, 5, 10))
	}
	for i := 0; i < 4; i++ {
		tests = append(tests, randomBatch(fmt.Sprintf("random_mid_%d", i+1), rng, 8, 1000))
	}
	tests = append(tests, randomBatch("random_large", rng, 5, 300000))

	return tests
}

func sampleInput() string {
	return `4
5
2 4 6 2 5
5
5 4 4 5 1
4
6 8 2 3
1
11
`
}

func buildSingleCase(name string, arrays [][]int64) testCase {
	var b strings.Builder
	fmt.Fprintln(&b, len(arrays))
	for _, arr := range arrays {
		fmt.Fprintln(&b, len(arr))
		for i, val := range arr {
			if i > 0 {
				fmt.Fprint(&b, " ")
			}
			fmt.Fprint(&b, val)
		}
		fmt.Fprintln(&b)
	}
	return testCase{name: name, input: b.String()}
}

func randomBatch(name string, rng *rand.Rand, cases int, maxN int) testCase {
	var b strings.Builder
	fmt.Fprintln(&b, cases)
	for i := 0; i < cases; i++ {
		n := rng.Intn(maxN) + 1
		fmt.Fprintln(&b, n)
		for j := 0; j < n; j++ {
			if j > 0 {
				fmt.Fprint(&b, " ")
			}
			val := rng.Int63n(1_000_000_000_000) + 1
			fmt.Fprint(&b, val)
		}
		fmt.Fprintln(&b)
	}
	return testCase{name: name, input: b.String()}
}
