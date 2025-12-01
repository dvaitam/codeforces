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

const refSource = "./2041F.go"

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
		tmp, err := os.CreateTemp("", "verifier2041F-*")
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
		buildSingle("tiny", [][2]int64{{1, 2}, {2, 4}, {3, 6}}),
		buildSingle("mid_manual", [][2]int64{{1, 10}, {4, 12}, {10, 20}}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 6; i++ {
		tests = append(tests, randomBatch(fmt.Sprintf("random_small_%d", i+1), rng, 5, 50))
	}
	for i := 0; i < 4; i++ {
		tests = append(tests, randomBatch(fmt.Sprintf("random_mid_%d", i+1), rng, 10, 2000))
	}
	tests = append(tests, randomBatch("random_large", rng, 10, 100000))

	return tests
}

func sampleInput() string {
	return `3
1 30
16 18
142857 240135
`
}

func buildSingle(name string, pairs [][2]int64) testCase {
	var b strings.Builder
	fmt.Fprintln(&b, len(pairs))
	for _, p := range pairs {
		fmt.Fprintf(&b, "%d %d\n", p[0], p[1])
	}
	return testCase{name: name, input: b.String()}
}

func randomBatch(name string, rng *rand.Rand, tMax int, spanLimit int64) testCase {
	if tMax < 1 {
		tMax = 1
	}
	t := rng.Intn(tMax) + 1
	var b strings.Builder
	fmt.Fprintln(&b, t)
	for i := 0; i < t; i++ {
		length := rng.Int63n(spanLimit) + 1
		if length > 100000 {
			length = 100000
		}
		l := rng.Int63n(1_000_000_000_000-length) + 1
		r := l + length
		fmt.Fprintf(&b, "%d %d\n", l, r)
	}
	return testCase{name: name, input: b.String()}
}
