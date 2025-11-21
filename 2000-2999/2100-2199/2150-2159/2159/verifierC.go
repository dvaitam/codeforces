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

const (
	numRandomTests = 200
	maxDegree      = 400000
)

type testCase struct {
	n      int
	coeffs []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}

	candidate, candCleanup, err := prepareBinary(os.Args[1])
	if err != nil {
		fmt.Println("failed to prepare contestant binary:", err)
		return
	}
	if candCleanup != nil {
		defer candCleanup()
	}

	oraclePath, oracleCleanup, err := prepareOracle()
	if err != nil {
		fmt.Println("failed to build reference solution:", err)
		return
	}
	defer oracleCleanup()

	tests := deterministicTests()
	total := 0
	for idx, input := range tests {
		if err := runCase(candidate, oraclePath, input); err != nil {
			fmt.Printf("deterministic case %d failed: %v\ninput:\n%s", idx+1, err, input)
			return
		}
		total++
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < numRandomTests; i++ {
		input := generateRandomInput(rng)
		if err := runCase(candidate, oraclePath, input); err != nil {
			fmt.Printf("random case %d failed: %v\ninput:\n%s", i+1, err, input)
			return
		}
		total++
	}

	fmt.Printf("All %d tests passed.\n", total)
}

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		abs, err := filepath.Abs(path)
		if err != nil {
			return "", nil, err
		}
		tmp := filepath.Join(os.TempDir(), fmt.Sprintf("candidateC_%d", time.Now().UnixNano()))
		cmd := exec.Command("go", "build", "-o", tmp, abs)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	return abs, nil, nil
}

func prepareOracle() (string, func(), error) {
	dir := sourceDir()
	src := filepath.Join(dir, "2159C.go")
	tmp := filepath.Join(os.TempDir(), fmt.Sprintf("oracleC_%d", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", tmp, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", nil, fmt.Errorf("go build oracle failed: %v: %s", err, out)
	}
	return tmp, func() { os.Remove(tmp) }, nil
}

func runCase(candidate, oracle, input string) error {
	expect, err := runBinary(oracle, input)
	if err != nil {
		return fmt.Errorf("oracle runtime error: %v", err)
	}
	got, err := runBinary(candidate, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	if expect != got {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func deterministicTests() []string {
	return []string{
		sampleInput(),
		conflictInput(),
		largeInput(),
	}
}

func sampleInput() string {
	cases := []testCase{
		{n: 1, coeffs: []int{-1, -1}},
		{n: 2, coeffs: []int{-1, 2, -1}},
		{n: 2, coeffs: []int{-1, -1, -1}},
		{n: 3, coeffs: []int{-1, -1, 3, -1}},
		{n: 3, coeffs: []int{-1, 2, 3, -1}},
		{n: 5, coeffs: []int{-1, -1, -1, 1, 0, -1}},
	}
	return buildInput(cases)
}

func conflictInput() string {
	cases := []testCase{
		{n: 3, coeffs: []int{-1, 0, -1, -1}},
		{n: 3, coeffs: []int{-1, 5, -1, -1}},
	}
	return buildInput(cases)
}

func largeInput() string {
	n := maxDegree
	coeffs := make([]int, n+1)
	for i := range coeffs {
		coeffs[i] = -1
	}
	coeffs[1] = 0
	if n >= 2 {
		coeffs[2] = 2
	}
	if n >= 3 {
		coeffs[n/2] = n / 2
		coeffs[n/2+1] = 1
	}
	return buildInput([]testCase{{n: n, coeffs: coeffs}})
}

func generateRandomInput(rng *rand.Rand) string {
	t := rng.Intn(4) + 1
	cases := make([]testCase, t)
	totalN := 0
	for i := 0; i < t; i++ {
		remain := 500 - totalN
		if remain < 1 {
			break
		}
		n := rng.Intn(min(50, remain)) + 1
		totalN += n
		cases[i] = randomCase(rng, n)
	}
	cases = compactCases(cases)
	return buildInput(cases)
}

func compactCases(cases []testCase) []testCase {
	out := cases[:0]
	for _, c := range cases {
		if c.n > 0 {
			out = append(out, c)
		}
	}
	if len(out) == 0 {
		return []testCase{{n: 1, coeffs: []int{-1, -1}}}
	}
	return out
}

func randomCase(rng *rand.Rand, n int) testCase {
	coeffs := make([]int, n+1)
	for i := range coeffs {
		coeffs[i] = -1
	}
	for i := 1; i < n; i++ {
		roll := rng.Intn(100)
		switch {
		case roll < 40:
			coeffs[i] = -1
		case roll < 65:
			coeffs[i] = 0
		case roll < 80:
			coeffs[i] = i
		case roll < 95:
			coeffs[i] = rng.Intn(n + 1)
		default:
			coeffs[i] = n + 1 + rng.Intn(5)
		}
	}
	coeffs[0] = -1
	coeffs[n] = -1
	return testCase{n: n, coeffs: coeffs}
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(cases))
	for _, c := range cases {
		fmt.Fprintf(&sb, "%d\n", c.n)
		for i, val := range c.coeffs {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", val)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func sourceDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Dir(file)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
