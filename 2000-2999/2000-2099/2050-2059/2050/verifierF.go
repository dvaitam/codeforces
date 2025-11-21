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

const refSource = "2000-2999/2000-2099/2050-2059/2050/2050F.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
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
		tmp, err := os.CreateTemp("", "verifier2050F-*")
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
	tests := []testCase{
		{name: "sample", input: sampleInput()},
		buildSingle("single_elem", []int{1}, [][]int{{1, 1}}),
		buildSingle("two_values", []int{5, 5}, [][]int{{1, 2}, {2, 2}}),
		buildSingle("simple_diff", []int{5, 14, 2, 6, 3}, [][]int{{1, 5}, {2, 4}, {3, 3}}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 6; i++ {
		tests = append(tests, randomBatch(fmt.Sprintf("random_small_%d", i+1), rng, 5, 20, 50))
	}
	for i := 0; i < 4; i++ {
		tests = append(tests, randomBatch(fmt.Sprintf("random_mid_%d", i+1), rng, 50, 500, 1000))
	}
	tests = append(tests, randomBatch("random_large", rng, 200000, 200000, 1000000000))

	return tests
}

func sampleInput() string {
	return `3
5 5
5 14 2 6 3
1 4
2 4
3 5
1 1
1 5
4 5
1 4 2 4
3 5
1 1
1 1 7 1
1 3
2 1
`
}

func buildSingle(name string, arr []int, queries [][]int) testCase {
	var b strings.Builder
	fmt.Fprintln(&b, 1)
	fmt.Fprintf(&b, "%d %d\n", len(arr), len(queries))
	for i, v := range arr {
		if i > 0 {
			fmt.Fprint(&b, " ")
		}
		fmt.Fprint(&b, v)
	}
	fmt.Fprintln(&b)
	for _, q := range queries {
		fmt.Fprintf(&b, "%d %d\n", q[0], q[1])
	}
	return testCase{name: name, input: b.String()}
}

func randomBatch(name string, rng *rand.Rand, nLimit int, qLimit int, maxVal int) testCase {
	t := rng.Intn(3) + 1
	var b strings.Builder
	fmt.Fprintln(&b, t)
	totalN := 0
	totalQ := 0
	for i := 0; i < t; i++ {
		n := rng.Intn(nLimit) + 1
		q := rng.Intn(qLimit) + 1
		if totalN+n > 200000 {
			n = 1
		}
		if totalQ+q > 200000 {
			q = 1
		}
		totalN += n
		totalQ += q
		fmt.Fprintf(&b, "%d %d\n", n, q)
		for j := 0; j < n; j++ {
			if j > 0 {
				fmt.Fprint(&b, " ")
			}
			fmt.Fprint(&b, rng.Intn(maxVal)+1)
		}
		fmt.Fprintln(&b)
		for j := 0; j < q; j++ {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			fmt.Fprintf(&b, "%d %d\n", l, r)
		}
	}
	return testCase{name: name, input: b.String()}
}
