package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "538F.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}

	candArg := os.Args[1]

	refBin, refCleanup, err := buildBinary(refSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer refCleanup()

	candBin, candCleanup, err := buildBinary(candArg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare candidate binary: %v\n", err)
		os.Exit(1)
	}
	defer candCleanup()

	tests := generateTests()
	for idx, tc := range tests {
		expect, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}

		got, err := runBinary(candBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}

		if !equalTokens(expect, got) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s)\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, tc.name, tc.input, expect, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildBinary(path string) (string, func(), error) {
	cleanPath := filepath.Clean(path)
	if strings.HasSuffix(cleanPath, ".go") {
		tmp, err := os.CreateTemp("", "538F-verifier-*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		source := filepath.Join(".", cleanPath)
		cmd := exec.Command("go", "build", "-o", tmp.Name(), source)
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
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
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
	rng := rand.New(rand.NewSource(538538538))
	var tests []testCase

	tests = append(tests,
		makeCase("sample1", []int64{1, 5, 4, 3, 2}),
		makeCase("sample2", []int64{2, 2, 2, 2, 2, 2}),
		makeCase("two_desc", []int64{5, -5}),
		makeCase("all_neg_desc", []int64{-1, -2, -3, -4, -5}),
		makeCase("zigzag", []int64{5, 1, 6, 2, 7, 3, 8, 4}),
		makeCase("plateau", []int64{0, 0, 0, 1, 1, 1, 1}),
		makeCase("alternating_small", []int64{-10, 10, -10, 10, -10, 10}),
		makeCase("increasing_small", []int64{-100, -80, -60, -40, -20, 0, 20, 40, 60}),
		makeCase("decreasing_small", []int64{100, 90, 80, 70, 60, 50, 40, 30, 20}),
	)

	tests = append(tests,
		alternatingCase("alternating_1e4", 10000, -1_000_000_000, 1_000_000_000),
		alternatingCase("alternating_border", 2000, 1_000_000_000, -1_000_000_000),
		gradientCase("gradient_up", 15000, -1_000_000_000, 1),
		gradientCase("gradient_down", 15000, 1_000_000_000, -1),
		blockCase("blocks_half", 20000, -500_000_000, 500_000_000),
		rampCase("ramp_mod", 50000),
	)

	for i := 0; i < 15; i++ {
		n := rng.Intn(20) + 2
		tests = append(tests, randomCase(fmt.Sprintf("random_small_%d", i+1), rng, n))
	}

	for i := 0; i < 10; i++ {
		n := rng.Intn(1000) + 50
		tests = append(tests, randomCase(fmt.Sprintf("random_mid_%d", i+1), rng, n))
	}

	for i := 0; i < 5; i++ {
		n := rng.Intn(5000) + 2000
		tests = append(tests, randomCase(fmt.Sprintf("random_large_%d", i+1), rng, n))
	}

	tests = append(tests,
		randomCase("random_50k", rng, 50000),
		randomCase("random_200k", rng, 200000),
		blockCase("blocks_thirds", 180000, -1_000_000_000, 1_000_000_000),
		gradientCase("gradient_full", 200000, -999_999_999, 10),
		rampCase("ramp_full", 200000),
	)

	return tests
}

func makeCase(name string, arr []int64) testCase {
	var b strings.Builder
	fmt.Fprintln(&b, len(arr))
	for i, v := range arr {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprint(&b, v)
	}
	b.WriteByte('\n')
	return testCase{name: name, input: b.String()}
}

func randomCase(name string, rng *rand.Rand, n int) testCase {
	if n < 2 {
		n = 2
	}
	var b strings.Builder
	b.Grow(n*12 + 20)
	fmt.Fprintln(&b, n)
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		val := rng.Int63n(2_000_000_001) - 1_000_000_000
		fmt.Fprint(&b, val)
	}
	b.WriteByte('\n')
	return testCase{name: name, input: b.String()}
}

func alternatingCase(name string, n int, a, b int64) testCase {
	if n < 2 {
		n = 2
	}
	var bld strings.Builder
	bld.Grow(n*12 + 20)
	fmt.Fprintln(&bld, n)
	for i := 0; i < n; i++ {
		if i > 0 {
			bld.WriteByte(' ')
		}
		if i%2 == 0 {
			fmt.Fprint(&bld, a)
		} else {
			fmt.Fprint(&bld, b)
		}
	}
	bld.WriteByte('\n')
	return testCase{name: name, input: bld.String()}
}

func gradientCase(name string, n int, start, step int64) testCase {
	if n < 2 {
		n = 2
	}
	var b strings.Builder
	b.Grow(n*12 + 20)
	fmt.Fprintln(&b, n)
	val := start
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(' ')
			val += step
		}
		if i == 0 {
			val = start
		}
		if val > 1_000_000_000 {
			val = 1_000_000_000
		}
		if val < -1_000_000_000 {
			val = -1_000_000_000
		}
		fmt.Fprint(&b, val)
	}
	b.WriteByte('\n')
	return testCase{name: name, input: b.String()}
}

func blockCase(name string, n int, low, high int64) testCase {
	if n < 2 {
		n = 2
	}
	var b strings.Builder
	b.Grow(n*12 + 20)
	fmt.Fprintln(&b, n)
	switchCount := n / 2
	if switchCount == 0 {
		switchCount = 1
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		if i < switchCount {
			fmt.Fprint(&b, low)
		} else {
			fmt.Fprint(&b, high)
		}
	}
	b.WriteByte('\n')
	return testCase{name: name, input: b.String()}
}

func rampCase(name string, n int) testCase {
	if n < 2 {
		n = 2
	}
	var b strings.Builder
	b.Grow(n*12 + 20)
	fmt.Fprintln(&b, n)
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		val := int64(((i*i + 17*i) % 2_000_000_001) - 1_000_000_000)
		fmt.Fprint(&b, val)
	}
	b.WriteByte('\n')
	return testCase{name: name, input: b.String()}
}
