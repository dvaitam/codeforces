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

const refSource = "2000-2999/2000-2099/2000-2009/2006/2006E.go"

type singleCase struct {
	n       int
	parents []int
}

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
		tmp, err := os.CreateTemp("", "verifier2006E-*")
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

	tests = append(tests, buildTest("simple_chain", []singleCase{
		makeLinearCase(5),
	}))

	tests = append(tests, buildTest("high_degree_invalid", []singleCase{
		makeStarCase(10, 1, 5),
	}))

	tests = append(tests, buildTest("mixed_cases", []singleCase{
		randomCase(8, 2),
		randomCase(12, 3),
		randomCase(15, 1),
	}))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 6; i++ {
		tests = append(tests, buildTest(
			fmt.Sprintf("random_small_%d", i+1),
			randomCaseBatch(rng, 3, 4, 20),
		))
	}
	for i := 0; i < 4; i++ {
		tests = append(tests, buildTest(
			fmt.Sprintf("random_mid_%d", i+1),
			randomCaseBatch(rng, 4, 10, 500),
		))
	}
	tests = append(tests, buildTest("random_large", randomCaseBatch(rng, 5, 20, 5000)))
	tests = append(tests, buildTest("stress_limit", []singleCase{
		randomCaseWithSeed(200000, 1),
		randomCaseWithSeed(150000, 7),
		randomCaseWithSeed(150000, 11),
	}))

	return tests
}

func sampleInput() string {
	return `7
3
1 2
6
1 2 3 4 5
7
1 1 3 2 5 1
10
1 1 2 1 4 2 4 5 8
10
1 1 3 1 3 2 2 2 6
20
1 1 2 2 4 4 5 5 7 6 8 6 11 14 11 8 13 13 12
25
1 1 3 3 1 5 4 4 6 8 11 12 8 7 11 13 7 13 15 6 19 14 10 23
`
}

func buildTest(name string, cases []singleCase) testCase {
	var b strings.Builder
	fmt.Fprintln(&b, len(cases))
	for _, cs := range cases {
		fmt.Fprintln(&b, cs.n)
		for i, val := range cs.parents {
			if i > 0 {
				fmt.Fprint(&b, " ")
			}
			fmt.Fprint(&b, val)
		}
		if len(cs.parents) > 0 {
			fmt.Fprintln(&b)
		}
	}
	return testCase{name: name, input: b.String()}
}

func makeLinearCase(n int) singleCase {
	if n < 2 {
		n = 2
	}
	parents := make([]int, n-1)
	for i := range parents {
		parents[i] = i + 1
	}
	return singleCase{n: n, parents: parents}
}

func makeStarCase(n, center, extra int) singleCase {
	if n < 2 {
		n = 2
	}
	if center < 1 || center >= n {
		center = 1
	}
	parents := make([]int, n-1)
	for i := range parents {
		v := i + 2
		p := center
		if p >= v {
			p = 1
		}
		parents[i] = p
	}
	for i := 0; i < extra && i < len(parents); i++ {
		parents[i] = i + 1
	}
	return singleCase{n: n, parents: parents}
}

func randomCase(n int, branch int) singleCase {
	if n < 2 {
		n = 2
	}
	rng := rand.New(rand.NewSource(int64(n)*31 + int64(branch)*131))
	parents := make([]int, n-1)
	for i := range parents {
		bound := i + 1
		if branch > 1 && bound > branch {
			bound = branch
		}
		parents[i] = rng.Intn(bound) + 1
	}
	return singleCase{n: n, parents: parents}
}

func randomCaseWithSeed(n int, seed int64) singleCase {
	rng := rand.New(rand.NewSource(seed))
	if n < 2 {
		n = 2
	}
	parents := make([]int, n-1)
	for i := range parents {
		parents[i] = rng.Intn(i+1) + 1
	}
	return singleCase{n: n, parents: parents}
}

func randomCaseBatch(rng *rand.Rand, minCases, maxCases, maxN int) []singleCase {
	caseCnt := rng.Intn(maxCases-minCases+1) + minCases
	result := make([]singleCase, 0, caseCnt)
	for len(result) < caseCnt {
		n := rng.Intn(maxN-1) + 2
		parents := make([]int, n-1)
		for i := range parents {
			parents[i] = rng.Intn(i+1) + 1
		}
		result = append(result, singleCase{n: n, parents: parents})
	}
	return result
}
