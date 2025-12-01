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

const refSource = "./2044C.go"

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
		tmp, err := os.CreateTemp("", "verifier2044C-*")
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
		buildSingle("edge_min", [][4]int64{
			{1, 1, 1, 1},
			{1, 2, 3, 4},
		}),
		buildSingle("balanced", [][4]int64{
			{10, 10, 10, 10},
			{5, 10, 10, 10},
		}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 6; i++ {
		tests = append(tests, randomBatch(fmt.Sprintf("random_small_%d", i+1), rng, 5, 50))
	}
	for i := 0; i < 4; i++ {
		tests = append(tests, randomBatch(fmt.Sprintf("random_mid_%d", i+1), rng, 20, 1000))
	}
	tests = append(tests, randomBatch("random_large", rng, 50, 100000000))

	return tests
}

func sampleInput() string {
	return `5
10 5 5 10
3 6 1 11
5 14 12 4
1 1 1 14
20 6 9 69
`
}

func buildSingle(name string, cases [][4]int64) testCase {
	var b strings.Builder
	fmt.Fprintln(&b, len(cases))
	for _, cs := range cases {
		fmt.Fprintf(&b, "%d %d %d %d\n", cs[0], cs[1], cs[2], cs[3])
	}
	return testCase{name: name, input: b.String()}
}

func randomBatch(name string, rng *rand.Rand, tMax int, maxVal int64) testCase {
	if tMax < 1 {
		tMax = 1
	}
	t := rng.Intn(tMax) + 1
	var b strings.Builder
	fmt.Fprintln(&b, t)
	for i := 0; i < t; i++ {
		m := randRange(rng, 1, maxVal)
		a := randRange(rng, 1, maxVal)
		bv := randRange(rng, 1, maxVal)
		c := randRange(rng, 1, maxVal)
		fmt.Fprintf(&b, "%d %d %d %d\n", m, a, bv, c)
	}
	return testCase{name: name, input: b.String()}
}

func randRange(rng *rand.Rand, lo, hi int64) int64 {
	if lo > hi {
		lo, hi = hi, lo
	}
	return rng.Int63n(hi-lo+1) + lo
}
